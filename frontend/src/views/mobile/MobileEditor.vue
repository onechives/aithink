<template>
  <section class="mobile-editor">
    <div class="header">
      <div>
        <h2>{{ isEditing ? $t('editor.title') : $t('editor.newTitle') }}</h2>
        <p>{{ $t('editor.description') }}</p>
      </div>
      <div class="actions">
        <button class="secondary" v-if="isEditing" @click="handleDelete">{{ $t('common.delete') }}</button>
        <button class="secondary" @click="handleSaveAndReturn">{{ $t('editor.saveAndReturn') }}</button>
        <button @click="handleSave">{{ $t('common.save') }}</button>
      </div>
    </div>
    <div class="form">
      <label>
        {{ $t('editor.postTitle') }}
        <input v-model="form.title" type="text" />
      </label>
      <label>
        {{ $t('editor.summary') }}
        <textarea v-model="form.summary" rows="3" />
      </label>
      <label>
        {{ $t('editor.category') }}
        <input v-model="form.category" type="text" />
      </label>
      <label>
        {{ $t('editor.tagsHint') }}
        <input v-model="form.tags" type="text" />
      </label>
      <label>
        {{ $t('editor.coverUrl') }}
        <input v-model="form.coverUrl" type="text" />
      </label>
    </div>
    <div class="editor-area">
      <MdEditor v-model="form.content" :toolbars="toolbars" :theme="theme" :language="editorLanguage"
        @upload-img="handleUpload" />
    </div>
    <span class="error" v-if="error">{{ error }}</span>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import { createPost, deletePost, getPostDetail, updatePost } from "../../api/posts";
import { uploadImage } from "../../api/uploads";
import { useTheme } from "../../composables/useTheme";
import { useI18n } from "vue-i18n";

const route = useRoute();
const router = useRouter();
const error = ref("");
const isEditing = ref(false);
const { theme } = useTheme();
const { locale } = useI18n();
const editorLanguage = computed(() => (locale.value === "en-US" ? "en-US" : "zh-CN"));
// 编辑器工具栏配置
const toolbars = [
  "bold",
  "italic",
  "underline",
  "strikeThrough",
  "title",
  "quote",
  "unorderedList",
  "orderedList",
  "codeRow",
  "code",
  "link",
  "image",
  "table",
  "preview",
];

async function handleUpload(files: File[], callback: (urls: string[]) => void) {
  // 上传图片并插入编辑器
  error.value = "";
  try {
    const urls = await uploadImage(files);
    callback(urls);
  } catch (err) {
    error.value = (err as Error).message;
  }
}

const form = reactive({
  // 编辑表单
  title: "",
  summary: "",
  content: "",
  coverUrl: "",
  category: "",
  tags: "",
});

async function loadPost() {
  // 编辑态加载文章详情
  const id = route.params.id ? String(route.params.id) : "";
  if (!id) return;
  const data = await getPostDetail(id);
  form.title = data.title;
  form.summary = data.summary;
  form.content = data.content;
  form.coverUrl = data.coverUrl;
  form.category = data.category;
  form.tags = data.tags;
  isEditing.value = true;
}

async function handleSave() {
  // 保存文章（新建/更新）
  error.value = "";
  try {
    if (isEditing.value && route.params.id) {
      await updatePost(String(route.params.id), form);
    } else {
      const result = await createPost(form);
      router.replace({ name: "m-editor", params: { id: result.id } });
    }
  } catch (err) {
    error.value = (err as Error).message;
  }
}

async function handleSaveAndReturn() {
  // 保存后回到移动端首页
  error.value = "";
  try {
    if (isEditing.value && route.params.id) {
      await updatePost(String(route.params.id), form);
    } else {
      await createPost(form);
    }
    router.push({ name: "m-home" });
  } catch (err) {
    error.value = (err as Error).message;
  }
}

async function handleDelete() {
  // 删除文章
  if (!route.params.id) return;
  try {
    await deletePost(String(route.params.id));
    router.push({ name: "m-home" });
  } catch (err) {
    error.value = (err as Error).message;
  }
}

onMounted(() => {
  loadPost();
});
</script>

<style scoped lang="scss">
.mobile-editor {
  padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
  display: grid;
  gap: var(--mobile-space-4);
}

.header h2 {
  margin: 0;
  font-size: var(--mobile-font-md);
}

.header p {
  margin: var(--mobile-space-1) 0 0;
  color: var(--muted);
}

.actions {
  display: grid;
  gap: var(--mobile-space-2);
  margin-top: var(--mobile-space-3);
}

button {
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: var(--mobile-radius);
  padding: 10px 16px;
  cursor: pointer;
  font-weight: 600;
}

button.secondary {
  background: var(--surface);
  color: var(--ink);
  border: 1px solid var(--border);
}

.form {
  display: grid;
  gap: var(--mobile-space-3);
  background: var(--surface);
  padding: var(--mobile-space-3);
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
}

label {
  display: grid;
  gap: 8px;
  font-weight: 600;
}

input,
textarea {
  border-radius: 8px;
  border: 1px solid var(--border);
  padding: 10px 12px;
  font-size: 14px;
  font-family: inherit;
  background: var(--surface);
  color: var(--ink);
}

.editor-area {
  border-radius: var(--mobile-radius);
  overflow: hidden;
  border: 1px solid var(--border);
  background: var(--surface);
}

.error {
  color: #b00020;
}

.editor-area :deep(.md-editor-toolbar-wrapper) {
  overflow-x: hidden;
}

.editor-area :deep(.md-editor-toolbar) {
  flex-wrap: wrap;
  justify-content: flex-start;
  row-gap: 6px;
}

.editor-area :deep(.md-editor-toolbar-left),
.editor-area :deep(.md-editor-toolbar-right) {
  flex-wrap: wrap;
  width: 100%;
  justify-content: flex-start;
}

.editor-area :deep(.md-editor-content) {
  flex-direction: column;
}

.editor-area :deep(.md-editor-resize-operate) {
  display: none;
}

.editor-area :deep(.md-editor-input-wrapper),
.editor-area :deep(.md-editor-preview-wrapper) {
  width: 100%;
  flex: 1 1 0;
}
</style>
