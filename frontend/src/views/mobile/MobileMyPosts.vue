<template>
  <section class="mobile-my-posts">
    <div v-if="!isAuthed" class="guest-card">
      <p>{{ $t('settings.loginHint') }}</p>
      <div class="guest-actions">
        <RouterLink to="/m/login">{{ $t('nav.login') }}</RouterLink>
        <RouterLink to="/m/register">{{ $t('nav.register') }}</RouterLink>
      </div>
    </div>

    <div v-else class="list">
      <article v-for="post in posts" :key="post.id" class="item">
        <div class="meta">
          <h3>{{ post.title }}</h3>
          <span class="status" :data-status="post.status">{{ statusLabel(post.status) }}</span>
        </div>
        <p class="date">{{ formatDate(post.updatedAt) }}</p>
        <div class="actions">
          <RouterLink :to="`/m/write/${post.id}`">{{ $t('admin.edit') }}</RouterLink>
          <RouterLink :to="`/m/posts/${post.id}`">{{ $t('common.preview') }}</RouterLink>
        </div>
      </article>
    </div>

    <div class="pager" v-if="isAuthed && total > size">
      <button :disabled="page === 1" @click="prevPage">{{ $t('home.prevPage') }}</button>
      <span>{{ page }} / {{ Math.ceil(total / size) }}</span>
      <button :disabled="page >= Math.ceil(total / size)" @click="nextPage">{{ $t('home.nextPage') }}</button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { listMyPosts, type PostSummary } from "../../api/posts";
import { useAuthStore } from "../../stores/auth";

const posts = ref<PostSummary[]>([]);
const total = ref(0);
const page = ref(1);
const size = 10;
const { t } = useI18n();
const auth = useAuthStore();
const isAuthed = computed(() => Boolean(auth.token));

async function fetchPosts() {
  if (!auth.token) {
    posts.value = [];
    total.value = 0;
    return;
  }
  const data = await listMyPosts({ page: page.value, size });
  posts.value = data.items;
  total.value = data.total;
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString();
}

function statusLabel(status: string) {
  if (status === "approved") return t("admin.approved");
  if (status === "rejected") return t("admin.rejected2");
  return t("admin.pending");
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
.mobile-my-posts {
  padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
  display: grid;
  gap: var(--mobile-space-3);
}

.guest-card {
  background: var(--surface);
  border-radius: var(--mobile-radius);
  padding: var(--mobile-space-3);
  border: 1px solid var(--border);
  display: grid;
  gap: var(--mobile-space-2);
}

.guest-actions {
  display: flex;
  gap: var(--mobile-space-2);
}

.guest-actions a {
  text-decoration: none;
  font-weight: 600;
  color: var(--ink);
  padding: 10px 14px;
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
  background: var(--surface-alt);
}

.list {
  display: grid;
  gap: var(--mobile-space-2);
}

.item {
  background: var(--surface);
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
  padding: var(--mobile-space-3);
  display: grid;
  gap: 8px;
}

.meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--mobile-space-2);
}

.meta h3 {
  margin: 0;
  font-size: 16px;
}

.status {
  font-size: 12px;
  font-weight: 600;
  color: var(--muted);
}

.status[data-status="approved"] {
  color: #1b7b3d;
}

.status[data-status="rejected"] {
  color: #b00020;
}

.date {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
}

.actions {
  display: flex;
  gap: var(--mobile-space-2);
}

.actions a {
  text-decoration: none;
  font-weight: 600;
  color: var(--ink);
}

.pager {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--mobile-space-2);
}

.pager button {
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
  background: var(--surface);
  padding: 8px 12px;
  color: var(--ink);
}
</style>
