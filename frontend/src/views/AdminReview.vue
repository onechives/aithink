<template>
  <section class="admin">
    <header class="toolbar">
      <div>
        <h2>{{ $t('admin.postReview') }}</h2>
        <p>{{ $t('admin.postReviewDesc') }}</p>
      </div>
      <div class="actions">
        <button class="secondary" @click="setStatus('pending')">{{ $t('admin.pending') }}</button>
        <button class="secondary" @click="setStatus('approved')">{{ $t('admin.approved') }}</button>
        <button class="secondary" @click="setStatus('rejected')">{{ $t('admin.rejected2') }}</button>
      </div>
    </header>

    <div class="table">
      <div class="row header">
        <span>{{ $t('admin.title') }}</span>
        <span>{{ $t('admin.author') }}</span>
        <span>{{ $t('admin.status') }}</span>
        <span>{{ $t('admin.action') }}</span>
      </div>
      <div class="row" v-for="post in posts" :key="post.id">
        <span class="title" :data-label="$t('admin.title')">{{ post.title }}</span>
        <span :data-label="$t('admin.author')">{{ post.author }}</span>
        <span class="status" :data-status="post.status" :data-label="$t('admin.status')">{{ statusLabel(post.status) }}</span>
        <span class="row-actions" :data-label="$t('admin.action')">
          <button v-if="post.status !== 'approved'" @click="handleApprove(post.id)">{{ $t('admin.approve') }}</button>
          <button v-if="post.status !== 'rejected'" class="secondary" @click="handleReject(post.id)">{{
            $t('admin.reject') }}</button>
          <RouterLink class="link" :to="`${detailBase}/${post.id}`">{{ $t('common.preview') }}</RouterLink>
        </span>
      </div>
    </div>

    <div class="pager" v-if="total > size">
      <button :disabled="page === 1" @click="prevPage">{{ $t('home.prevPage') }}</button>
      <span>{{ page }} / {{ Math.ceil(total / size) }}</span>
      <button :disabled="page >= Math.ceil(total / size)" @click="nextPage">{{ $t('home.nextPage') }}</button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { useI18n } from "vue-i18n";
import { approvePost, listAdminPosts, rejectPost, type PostSummary } from "../api/posts";

const posts = ref<PostSummary[]>([]);
const total = ref(0);
const page = ref(1);
const size = 10;
const status = ref("pending");
const { t } = useI18n();
const route = useRoute();
const isMobileRoute = computed(() => route.path.startsWith("/m"));
const detailBase = computed(() => (isMobileRoute.value ? "/m/posts" : "/posts"));

async function fetchPosts() {
  // 按审核状态拉取文章
  const data = await listAdminPosts({ status: status.value, page: page.value, size });
  posts.value = data.items;
  total.value = data.total;
}

function setStatus(value: string) {
  // 切换状态并刷新列表
  status.value = value;
  page.value = 1;
  fetchPosts();
}

function statusLabel(value: string) {
  // 状态展示为多语言文案
  if (value === "approved") return t("admin.approved");
  if (value === "rejected") return t("admin.rejected2");
  return t("admin.pending");
}

async function handleApprove(id: string) {
  // 通过审核
  await approvePost(id);
  fetchPosts();
}

async function handleReject(id: string) {
  // 驳回审核并填写原因
  const reason = window.prompt(t("admin.rejectReasonPrompt"), "") || "";
  await rejectPost(id, reason.trim());
  fetchPosts();
}

function prevPage() {
  if (page.value > 1) {
    page.value -= 1;
    fetchPosts();
  }
}

function nextPage() {
  const maxPage = Math.ceil(total.value / size);
  if (page.value < maxPage) {
    page.value += 1;
    fetchPosts();
  }
}

onMounted(() => {
  fetchPosts();
});
</script>

<style scoped lang="scss">
.admin {
  padding: 32px 36px 60px;
  display: grid;
  gap: 24px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
}

.toolbar h2 {
  margin: 0;
  font-size: 28px;
}

.toolbar p {
  margin: 8px 0 0;
  color: var(--muted);
}

.actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

button {
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: 10px;
  padding: 8px 14px;
  cursor: pointer;
  font-weight: 600;
}

button.secondary {
  background: var(--surface-alt);
  color: var(--ink);
}

.table {
  background: var(--surface);
  border-radius: 20px;
  border: 1px solid var(--border);
  overflow: hidden;
}

.row {
  display: grid;
  grid-template-columns: 2fr 0.8fr 0.8fr 1.2fr;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border);
  align-items: center;
  gap: 12px;
}

.row.header {
  background: var(--surface-alt);
  font-weight: 600;
}

.row:last-child {
  border-bottom: none;
}

.title {
  font-weight: 600;
}

.status[data-status="approved"] {
  color: #1b7b3d;
}

.status[data-status="rejected"] {
  color: #b00020;
}

.row .link {
  color: var(--muted);
  text-decoration: none;
  font-weight: 600;
}

.pager {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.pager button {
  border-radius: 10px;
  border: 1px solid var(--border);
  background: var(--surface);
  padding: 8px 16px;
  color: var(--ink);
}

@media (max-width: 960px) {
  .toolbar {
    flex-direction: column;
    align-items: flex-start;
  }

  .row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .admin {
    padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
  }

  .actions {
    flex-wrap: wrap;
  }

  .actions button {
    flex: 1 1 auto;
  }

  .table {
    border: none;
    background: transparent;
  }

  .row {
    grid-template-columns: 1fr;
    border: 1px solid var(--border);
    border-radius: var(--mobile-radius);
    margin-bottom: var(--mobile-space-3);
  }

  .row.header {
    display: none;
  }

  .row > span {
    display: grid;
    grid-template-columns: minmax(calc(var(--mobile-unit) * 26), 1fr) 2fr;
    gap: var(--mobile-space-2);
    align-items: center;
  }

  .row > span::before {
    content: attr(data-label);
    color: var(--muted);
    font-size: var(--mobile-font-xs);
    text-transform: uppercase;
    letter-spacing: calc(var(--mobile-unit) * 0.1);
  }

  .row .row-actions {
    display: grid;
    gap: var(--mobile-space-2);
  }
}

</style>
