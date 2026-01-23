<template>
  <section class="admin">
    <header class="toolbar">
      <div>
        <h2>{{ $t('admin.postManage') }}</h2>
        <p>{{ $t('editor.description') }}</p>
      </div>
      <div class="actions">
        <RouterLink class="nav" to="/admin/review/posts">{{ $t('admin.postReview') }}</RouterLink>
        <RouterLink class="nav" to="/admin/review/users">{{ $t('admin.userReview') }}</RouterLink>
        <RouterLink class="nav" to="/admin/review/nicknames">{{ $t('admin.nicknameReview') }}</RouterLink>
        <input v-model="keyword" type="text" :placeholder="$t('common.search') + '...'" @input="handleSearch" />
        <button @click="router.push({ name: 'editor' })">{{ $t('editor.newTitle') }}</button>
      </div>
    </header>

    <div class="table">
      <div class="row header">
        <span>{{ $t('admin.title') }}</span>
        <span>{{ $t('editor.category') }}</span>
        <span>{{ $t('home.likes') }}</span>
        <span>{{ $t('common.save') }}</span>
        <span>{{ $t('admin.action') }}</span>
      </div>
      <div class="row" v-for="post in posts" :key="post.id">
        <span class="title" :data-label="$t('admin.title')">{{ post.title }}</span>
        <span :data-label="$t('editor.category')">{{ post.category }}</span>
        <span :data-label="$t('home.likes')">{{ post.likeCount }}</span>
        <span :data-label="$t('common.save')">{{ formatDate(post.updatedAt) }}</span>
        <span class="row-actions" :data-label="$t('admin.action')">
          <RouterLink :to="`/write/${post.id}`">{{ $t('admin.edit') }}</RouterLink>
          <RouterLink class="link" :to="`/posts/${post.id}`">{{ $t('common.preview') }}</RouterLink>
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
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { listPosts, type PostSummary } from "../api/posts";

const router = useRouter();
const posts = ref<PostSummary[]>([]);
const total = ref(0);
const page = ref(1);
const size = 10;
const sort = "latest";
const keyword = ref("");

async function fetchPosts() {
  // 按关键词与分页获取文章列表
  const data = await listPosts({
    sort,
    keyword: keyword.value,
    page: page.value,
    size,
  });
  posts.value = data.items;
  total.value = data.total;
}

function formatDate(value: string) {
  // 日期展示格式
  return new Date(value).toLocaleDateString();
}

let typingTimer: number | undefined;
function handleSearch() {
  // 简单防抖搜索
  window.clearTimeout(typingTimer);
  typingTimer = window.setTimeout(() => {
    page.value = 1;
    fetchPosts();
  }, 300);
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
  gap: 12px;
}

input {
  border-radius: 8px;
  border: 1px solid var(--border);
  padding: 10px 12px;
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
}

.nav {
  text-decoration: none;
  color: var(--ink);
  font-weight: 600;
  background: var(--surface-alt);
  padding: 10px 14px;
  border-radius: 10px;
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

.table {
  background: var(--surface);
  border-radius: 20px;
  border: 1px solid var(--border);
  overflow: hidden;
}

.row {
  display: grid;
  grid-template-columns: 2fr 0.8fr 0.6fr 0.8fr 0.8fr;
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

.row a {
  text-decoration: none;
  color: var(--ink);
  font-weight: 600;
  margin-right: 12px;
}

.row a.link {
  color: var(--muted);
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
  cursor: pointer;
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
    align-items: stretch;
  }

  .actions input {
    width: 100%;
  }

  .actions button,
  .actions .nav {
    width: 100%;
    text-align: center;
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
