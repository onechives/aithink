<template>
  <BlogLayout sidebarNarrow>
    <template #sidebar>
      <SidebarList />
    </template>
    <article class="detail" v-if="post">
      <header class="hero">
        <div>
          <p class="meta">
            <span>{{ post.author }}</span>
            <span class="dot">•</span>
            <span>{{ post.category }}</span>
            <span class="dot">•</span>
            <span>{{ formatDate(post.createdAt) }}</span>
          </p>
          <h2>{{ post.title }}</h2>
          <p class="summary">{{ post.summary }}</p>
          <div class="tags" v-if="post.tags">
            <span v-for="tag in tags" :key="tag">{{ tag }}</span>
          </div>
        </div>
        <div class="cover" v-if="post.coverUrl" :style="{ backgroundImage: `url(${post.coverUrl})` }" />
      </header>
      <section class="content" v-html="htmlContent" @click="handleContentClick" />
      <button class="like" @click="handleLike">
        <svg class="icon" viewBox="0 0 1024 1024" aria-hidden="true">
          <path
            d="M621.674667 408.021333c16.618667-74.24 28.224-127.936 34.837333-161.194666C673.152 163.093333 629.941333 85.333333 544.298667 85.333333c-77.226667 0-116.010667 38.378667-138.88 115.093334l-0.586667 2.24c-13.728 62.058667-34.72 110.165333-62.506667 144.586666a158.261333 158.261333 0 0 1-119.733333 58.965334l-21.909333 0.469333C148.437333 407.808 106.666667 450.816 106.666667 503.498667V821.333333c0 64.8 52.106667 117.333333 116.394666 117.333334h412.522667c84.736 0 160.373333-53.568 189.12-133.92l85.696-239.584c21.802667-60.96-9.536-128.202667-70.005333-150.186667a115.552 115.552 0 0 0-39.488-6.954667H621.674667zM544.256 149.333333c39.253333 0 59.498667 36.48 49.888 84.928-7.573333 38.144-21.984 104.426667-43.221333 198.666667-4.512 20.021333 10.56 39.093333 30.912 39.093333h218.666666c6.101333 0 12.16 1.066667 17.909334 3.168 27.445333 9.984 41.674667 40.554667 31.776 68.266667l-85.568 239.573333C744.981333 838.026667 693.301333 874.666667 635.402667 874.666667H223.498667C194.314667 874.666667 170.666667 850.784 170.666667 821.333333V503.498667c0-17.866667 14.144-32.448 31.829333-32.821334l21.866667-0.469333a221.12 221.12 0 0 0 167.381333-82.56c34.346667-42.602667 59.146667-99.306667 74.869333-169.877333C482.101333 166.336 499.552 149.333333 544.266667 149.333333z"
          />
        </svg>
        <span>{{ post.likeCount }}</span>
      </button>
    </article>
  </BlogLayout>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { marked } from "marked";
import hljs from "highlight.js";
import { useI18n } from "vue-i18n";
import BlogLayout from "../components/BlogLayout.vue";
import SidebarList from "../components/SidebarList.vue";
import { getPostDetail, likePost, type PostDetail } from "../api/posts";

const route = useRoute();
const post = ref<PostDetail | null>(null);
const { t } = useI18n();

const copyLabel = computed(() => t("post.code.copy"));
const copiedLabel = computed(() => t("post.code.copied"));
const copyFailedLabel = computed(() => t("post.code.copyFailed"));

const tags = computed(() => {
  // 逗号分隔的标签字符串转数组
  if (!post.value?.tags) return [];
  return post.value.tags.split(",").map((t) => t.trim()).filter(Boolean);
});

const htmlContent = computed(() => {
  // Markdown 渲染为 HTML
  if (!post.value) return "";
  return renderMarkdown(post.value.content, copyLabel.value);
});

function escapeHtml(value: string) {
  return value
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function renderMarkdown(content: string, copyText: string) {
  const renderer = new marked.Renderer();
  renderer.code = (code, infostring, escaped) => {
    const isToken = typeof code === "object" && code !== null;
    const rawCode = isToken ? String((code as { text?: unknown }).text ?? "") : String(code ?? "");
    const rawLangInput = isToken ? String((code as { lang?: unknown }).lang ?? "") : String(infostring ?? "");
    const isEscaped = isToken ? Boolean((code as { escaped?: unknown }).escaped) : Boolean(escaped);
    const rawLang = rawLangInput.match(/\S*/)?.[0] || "";
    const safeLang = rawLang.replace(/[^a-zA-Z0-9_-]/g, "");
    const langLabel = safeLang || "text";
    const escapedLang = escapeHtml(langLabel);
    const escapedCopy = escapeHtml(copyText);
    const langClass = safeLang ? `language-${safeLang}` : "language-text";

    let highlighted = "";
    if (rawCode.trim()) {
      try {
        highlighted = safeLang && hljs.getLanguage(safeLang)
          ? hljs.highlight(rawCode, { language: safeLang, ignoreIllegals: true }).value
          : hljs.highlightAuto(rawCode).value;
      } catch {
        highlighted = "";
      }
    }
    const escapedCode = highlighted || (isEscaped ? rawCode : escapeHtml(rawCode));

    return `
      <div class="code-block" data-lang="${escapedLang}">
        <div class="code-toolbar">
          <span class="code-dots" aria-hidden="true">
            <i></i><i></i><i></i>
          </span>
          <span class="code-lang">${escapedLang}</span>
          <button class="code-copy" type="button" aria-label="${escapedCopy}">${escapedCopy}</button>
        </div>
        <pre><code class="hljs ${langClass}">${escapedCode}</code></pre>
      </div>
    `;
  };
  return marked.parse(content, { renderer });
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString();
}

async function fetchDetail() {
  // 根据路由参数加载文章详情
  const id = route.params.id ? String(route.params.id) : "";
  if (!id) return;
  post.value = await getPostDetail(id);
}

async function handleLike() {
  // 点赞后更新本地计数
  if (!post.value) return;
  const result = await likePost(post.value.id);
  post.value.likeCount = result.likeCount;
}

async function handleContentClick(event: MouseEvent) {
  const target = event.target as HTMLElement | null;
  const button = target?.closest(".code-copy") as HTMLButtonElement | null;
  if (!button) return;
  const wrapper = button.closest(".code-block");
  const code = wrapper?.querySelector("code")?.textContent ?? "";
  if (!code) return;

  const existingTimeout = button.dataset.resetTimeout;
  if (existingTimeout) {
    window.clearTimeout(Number(existingTimeout));
  }

  let success = false;
  if (navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(code);
      success = true;
    } catch {
      success = false;
    }
  }

  if (!success) {
    const textarea = document.createElement("textarea");
    textarea.value = code;
    textarea.setAttribute("readonly", "");
    textarea.style.position = "fixed";
    textarea.style.opacity = "0";
    document.body.appendChild(textarea);
    textarea.focus();
    textarea.select();
    try {
      success = document.execCommand("copy");
    } catch {
      success = false;
    }
    document.body.removeChild(textarea);
  }

  const originalLabel = button.dataset.label || button.textContent || copyLabel.value;
  button.dataset.label = originalLabel;
  button.textContent = success ? copiedLabel.value : copyFailedLabel.value;
  button.classList.toggle("is-copied", success);

  const resetTimeout = window.setTimeout(() => {
    button.textContent = button.dataset.label || copyLabel.value;
    button.classList.remove("is-copied");
  }, 1600);
  button.dataset.resetTimeout = String(resetTimeout);
}

watch(
  () => route.params.id,
  () => {
    fetchDetail();
  },
  { immediate: true }
);
</script>

<style scoped lang="scss">
.detail {
  background: var(--surface);
  border-radius: 24px;
  padding: 28px;
  box-shadow: var(--shadow);
}

.hero {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 22px;
  align-items: center;
  border-bottom: 1px solid var(--border);
  padding-bottom: 20px;
}

.meta {
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 12px;
  color: var(--muted);
}

.dot {
  margin: 0 8px;
}

.summary {
  color: var(--muted);
}

.cover {
  border-radius: 20px;
  min-height: 200px;
  background-size: cover;
  background-position: center;
}

.content {
  margin-top: 24px;
  line-height: 1.8;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.content :deep(a) {
  color: var(--md-link);
  text-decoration: underline;
  text-decoration-thickness: 2px;
  text-underline-offset: 3px;
}

.content :deep(blockquote) {
  margin: 20px 0;
  padding: 12px 18px;
  border-left: 3px solid var(--md-border);
  background: var(--md-blockquote);
}

.content :deep(.code-block) {
  margin: 22px 0;
  border-radius: 14px;
  overflow: hidden;
  border: 1px solid var(--md-code-toolbar-border);
  background: var(--md-code-bg);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.06);
}

.content :deep(.code-toolbar) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  background: var(--md-code-toolbar);
  border-bottom: 1px solid var(--md-code-toolbar-border);
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--muted);
}

.content :deep(.code-dots) {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.content :deep(.code-dots i) {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  display: inline-block;
}

.content :deep(.code-dots i:nth-child(1)) {
  background: #ff5f56;
}

.content :deep(.code-dots i:nth-child(2)) {
  background: #ffbd2e;
}

.content :deep(.code-dots i:nth-child(3)) {
  background: #27c93f;
}

.content :deep(.code-lang) {
  font-weight: 600;
  font-size: 11px;
  letter-spacing: 0.06em;
}

.content :deep(.code-copy) {
  border: 1px solid var(--md-code-copy-border);
  background: var(--md-code-copy-bg);
  color: var(--md-code-copy-text);
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  cursor: pointer;
  transition: transform 0.15s ease, background 0.15s ease, border-color 0.15s ease;
}

.content :deep(.code-copy:hover) {
  background: var(--md-code-copy-bg-hover);
  border-color: var(--md-code-copy-border);
  transform: translateY(-1px);
}

.content :deep(.code-copy.is-copied) {
  background: #27c93f;
  color: #111111;
  border-color: #27c93f;
}

.content :deep(pre) {
  background: var(--md-code-bg);
  color: var(--md-code-text);
  padding: 16px;
  border-radius: 14px;
  overflow: auto;
  max-width: 100%;
}

.content :deep(.code-block pre) {
  margin: 0;
  border-radius: 0;
  background: var(--md-code-bg);
}

.content :deep(code) {
  font-family: "SFMono-Regular", Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

.content :deep(p code),
.content :deep(li code) {
  background: var(--md-inline-code-bg);
  padding: 2px 6px;
  border-radius: 6px;
}

.content :deep(table) {
  display: block;
  max-width: 100%;
  overflow-x: auto;
  border-collapse: collapse;
  margin-top: 16px;
  -webkit-overflow-scrolling: touch;
}

.content :deep(img) {
  max-width: 100%;
  height: auto;
  display: block;
}

.content :deep(th),
.content :deep(td) {
  border: 1px solid var(--md-border);
  padding: 10px;
  text-align: left;
}

.content :deep(h1),
.content :deep(h2),
.content :deep(h3) {
  margin-top: 28px;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.tags span {
  background: var(--chip);
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
}

.like {
  margin-top: 24px;
  border: none;
  background: transparent;
  color: var(--ink);
  padding: 0;
  cursor: pointer;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.like .icon {
  width: 18px;
  height: 18px;
  fill: currentColor;
}

@media (max-width: 768px) {
  .detail {
    padding: var(--mobile-space-4);
  }

  .hero {
    grid-template-columns: 1fr;
  }

  .cover {
    min-height: calc(var(--mobile-unit) * 45);
  }
}

</style>
