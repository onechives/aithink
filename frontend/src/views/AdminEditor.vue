<template>
  <section class="editor">
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
      <div class="row">
        <label>
          {{ $t('editor.category') }}
          <input v-model="form.category" type="text" />
        </label>
        <label>
          {{ $t('editor.tagsHint') }}
          <input v-model="form.tags" type="text" />
        </label>
      </div>
      <div class="cover-field">
        <span>{{ $t('editor.coverUrl') }}</span>
        <div class="cover-actions">
          <label class="cover-upload">
            <input type="file" accept="image/*" @change="handleCoverSelect" />
            <span>{{ form.coverUrl ? $t('editor.coverChange') : $t('editor.coverUpload') }}</span>
          </label>
          <button class="secondary" type="button" v-if="form.coverUrl" @click="handleCoverRemove">
            {{ $t('editor.coverRemove') }}
          </button>
        </div>
        <p class="cover-hint">{{ $t('editor.coverHint') }}</p>
        <div v-if="form.coverUrl" class="cover-preview">
          <img :src="form.coverUrl" alt="" />
        </div>
      </div>
    </div>
    <div v-if="coverCropVisible" class="cover-cropper">
      <div class="cover-cropper__backdrop" @click="handleCoverCropCancel" />
      <div class="cover-cropper__panel">
        <img ref="coverImage" :src="coverPreview" alt="" />
        <div class="cover-cropper__actions">
          <button class="secondary" type="button" @click="handleCoverCropCancel">{{ $t('common.cancel') }}</button>
          <button type="button" @click="handleCoverCropConfirm" :disabled="coverUploading">
            {{ coverUploading ? $t('common.loading') : $t('editor.coverCrop') }}
          </button>
        </div>
      </div>
    </div>
    <div class="editor-area">
      <MdEditor v-model="form.content" :toolbars="toolbars" :theme="theme" :language="editorLanguage"
        @upload-img="handleUpload" />
    </div>
    <span class="error" v-if="error">{{ error }}</span>
  </section>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { MdEditor } from "md-editor-v3";
import "md-editor-v3/lib/style.css";
import Cropper from "cropperjs";
import "cropperjs/dist/cropper.css";
import { createPost, deletePost, getPostDetail, updatePost } from "../api/posts";
import { uploadImage } from "../api/uploads";
import { useTheme } from "../composables/useTheme";
import { useI18n } from "vue-i18n";

const route = useRoute();
const router = useRouter();
const error = ref("");
const isEditing = ref(false);
const { theme } = useTheme();
const { locale } = useI18n();
const editorLanguage = computed(() => (locale.value === "en-US" ? "en-US" : "zh-CN"));
// 编辑器工具栏配置（控制按钮顺序与显示）
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
  // 上传图片并回填到编辑器
  error.value = "";
  try {
    const urls = await uploadImage(files);
    callback(urls);
  } catch (err) {
    error.value = (err as Error).message;
  }
}

const form = reactive({
  // 文章编辑表单
  title: "",
  summary: "",
  content: "",
  coverUrl: "",
  category: "",
  tags: "",
});

const coverCropVisible = ref(false);
const coverPreview = ref("");
const coverUploading = ref(false);
const coverImage = ref<HTMLImageElement | null>(null);
let coverCropper: Cropper | null = null;
let coverObjectUrl = "";

async function loadPost() {
  // 进入编辑态时加载文章详情
  const id = route.params.id ? String(route.params.id) : "";
  if (!id) return;
  const data = await getPostDetail(id);
  form.title = data.title;
  form.summary = data.summary || deriveSummary(data.content, 100);
  form.content = data.content;
  form.coverUrl = data.coverUrl;
  form.category = data.category;
  form.tags = data.tags;
  isEditing.value = true;
}

function applySummary() {
  if (!form.summary.trim()) {
    form.summary = deriveSummary(form.content, 100);
  }
}

async function handleSave() {
  // 保存文章（新建或更新）
  error.value = "";
  try {
    applySummary();
    if (isEditing.value && route.params.id) {
      await updatePost(String(route.params.id), form);
    } else {
      const result = await createPost(form);
      router.replace({ name: "editor", params: { id: result.id } });
    }
  } catch (err) {
    error.value = (err as Error).message;
  }
}

async function handleSaveAndReturn() {
  // 保存后返回“我的文章”
  error.value = "";
  try {
    applySummary();
    if (isEditing.value && route.params.id) {
      await updatePost(String(route.params.id), form);
    } else {
      await createPost(form);
    }
    router.push({ name: "my-posts" });
  } catch (err) {
    error.value = (err as Error).message;
  }
}

async function handleDelete() {
  // 删除文章后回到首页
  if (!route.params.id) return;
  try {
    await deletePost(String(route.params.id));
    router.push({ name: "home" });
  } catch (err) {
    error.value = (err as Error).message;
  }
}

function handleCoverSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  openCoverCropper(file);
  input.value = "";
}

function openCoverCropper(file: File) {
  cleanupCoverCropper();
  coverObjectUrl = URL.createObjectURL(file);
  coverPreview.value = coverObjectUrl;
  coverCropVisible.value = true;
  nextTick(() => {
    if (!coverImage.value) return;
    coverCropper = new Cropper(coverImage.value, {
      aspectRatio: 16 / 9,
      viewMode: 1,
      autoCropArea: 1,
      background: false,
      responsive: true,
    });
  });
}

async function handleCoverCropConfirm() {
  if (!coverCropper) return;
  coverUploading.value = true;
  try {
    const canvas = coverCropper.getCroppedCanvas({
      width: 1280,
      height: 720,
      imageSmoothingQuality: "high",
    });
    if (!canvas) {
      error.value = "cover crop failed";
      return;
    }
    const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, "image/jpeg", 0.9));
    if (!blob) {
      error.value = "cover crop failed";
      return;
    }
    const file = new File([blob], `cover-${Date.now()}.jpg`, { type: blob.type });
    const urls = await uploadImage([file]);
    form.coverUrl = urls[0] || "";
    coverCropVisible.value = false;
    cleanupCoverCropper();
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    coverUploading.value = false;
  }
}

function handleCoverCropCancel() {
  coverCropVisible.value = false;
  cleanupCoverCropper();
}

function handleCoverRemove() {
  form.coverUrl = "";
}

function cleanupCoverCropper() {
  if (coverCropper) {
    coverCropper.destroy();
    coverCropper = null;
  }
  if (coverObjectUrl) {
    URL.revokeObjectURL(coverObjectUrl);
    coverObjectUrl = "";
  }
  coverPreview.value = "";
}

function deriveSummary(content: string, maxLength: number) {
  const plain = content
    .replace(/```[\s\S]*?```/g, " ")
    .replace(/`[^`]*`/g, " ")
    .replace(/!\[[^\]]*]\([^)]*\)/g, " ")
    .replace(/\[([^\]]+)\]\([^)]*\)/g, "$1")
    .replace(/[#>*_-]+/g, " ")
    .replace(/\s+/g, " ")
    .trim();
  if (!plain) return "";
  return Array.from(plain).slice(0, maxLength).join("");
}

onMounted(() => {
  loadPost();
});

onBeforeUnmount(() => {
  cleanupCoverCropper();
});
</script>

<style scoped lang="scss">
.editor {
  padding: 32px 36px 60px;
  display: grid;
  gap: 24px;
  position: relative;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
  position: relative;
  z-index: 2;
}

.header h2 {
  margin: 0;
  font-size: 28px;
}

.header p {
  margin: 8px 0 0;
  color: var(--muted);
}

.actions {
  display: flex;
  gap: 12px;
}

button {
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: 10px;
  padding: 10px 18px;
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
  gap: 16px;
  background: var(--surface);
  padding: 20px;
  border-radius: 20px;
  border: 1px solid var(--border);
}

label {
  display: grid;
  gap: 8px;
  font-weight: 600;
}

.cover-field {
  display: grid;
  gap: 10px;
  font-weight: 600;
}

.cover-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.cover-upload {
  position: relative;
  overflow: hidden;
  border-radius: 999px;
  border: 1px dashed var(--border);
  padding: 8px 14px;
  font-size: 14px;
  cursor: pointer;
  background: var(--surface);
}

.cover-upload input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}

.cover-hint {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
  font-weight: 400;
}

.cover-preview {
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--border);
  max-width: 360px;
}

.cover-preview img {
  display: block;
  width: 100%;
  height: auto;
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
  resize: none;
}

.row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.editor-area {
  border-radius: 20px;
  overflow: hidden;
  position: relative;
  z-index: 1;
  border: 1px solid var(--border);
  background: var(--surface);
}

.editor-area :deep(.md-editor-footer) {
  order: -1;
  border-top: none;
  border-bottom: 1px solid var(--border);
  background: var(--surface);
}

.editor-area :deep(.md-editor-footer-item) {
  padding: 0 10px;
}

.error {
  color: #b00020;
}

.cover-cropper {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: grid;
  place-items: center;
}

.cover-cropper__backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
}

.cover-cropper__panel {
  position: relative;
  z-index: 1;
  width: min(720px, 92vw);
  background: var(--surface);
  border-radius: 20px;
  padding: 16px;
  display: grid;
  gap: 12px;
  box-shadow: var(--shadow);
  max-height: 80vh;
  overflow: hidden;
}

.cover-cropper__panel img {
  width: 100%;
  max-height: 60vh;
  display: block;
  object-fit: contain;
}

.cover-cropper__actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 960px) {
  .header {
    flex-direction: column;
    align-items: flex-start;
  }

  .row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .editor {
    padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
  }

  .actions {
    width: 100%;
    flex-wrap: wrap;
  }

  .actions button {
    width: 100%;
  }

  .header h2 {
    font-size: var(--mobile-font-md);
  }

  .editor-area {
    border-radius: var(--mobile-radius);
  }

  .cover-actions {
    width: 100%;
  }

  .cover-preview {
    max-width: 100%;
  }

  .cover-cropper__panel {
    width: min(600px, 92vw);
    border-radius: var(--mobile-radius);
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
}

</style>
