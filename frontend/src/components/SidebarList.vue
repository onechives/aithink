<template>
  <section class="sidebar-card">
    <div class="controls">
      <input v-model="keyword" type="text" :placeholder="`${$t('common.search')}...`" @input="emitSearch" />
      <div class="sort">
        <button :class="{ active: sort === 'likes' }" @click="setSort('likes')">{{ $t('home.mostLiked') }}</button>
        <button :class="{ active: sort === 'latest' }" @click="setSort('latest')">{{ $t('home.latest') }}</button>
      </div>
    </div>
    <ul class="title-list">
      <li v-for="item in items" :key="item.id">
        <RouterLink :to="`/posts/${item.id}`">{{ item.title }}</RouterLink>
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { listPostTitles, type PostTitle } from "../api/posts";

const items = ref<PostTitle[]>([]);
const keyword = ref("");
const sort = ref<"likes" | "latest">("latest");

async function fetchTitles() {
  // 根据关键词与排序拉取标题列表
  items.value = await listPostTitles({
    sort: sort.value,
    keyword: keyword.value,
    size: 50,
  });
}

function setSort(value: "likes" | "latest") {
  // 切换排序并刷新
  sort.value = value;
  fetchTitles();
}

let typingTimer: number | undefined;
function emitSearch() {
  // 简单防抖，避免频繁请求
  window.clearTimeout(typingTimer);
  typingTimer = window.setTimeout(() => {
    fetchTitles();
  }, 300);
}

onMounted(() => {
  fetchTitles();
});
</script>

<style scoped lang="scss">
.sidebar-card {
  padding: 20px 20px;
  border-radius: 20px;
  background: var(--surface-alt);
  border: 1px solid var(--border);
}

.controls {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.sort {
  display: flex;
  gap: 10px;
}

.sort button {
  flex: 1;
  border: none;
  background: var(--surface);
  padding: 10px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  color: var(--ink);
}

.sort button.active {
  background: var(--accent);
  color: var(--accent-contrast);
}

input {
  border-radius: 8px;
  border: 1px solid var(--border);
  padding: 10px 12px;
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
}

.title-list {
  list-style: none;
  padding: 0;
  margin: 18px 0 0;
  display: grid;
  gap: 12px;
  max-width: 100%;
}

.title-list a {
  text-decoration: none;
  color: var(--ink);
  font-weight: 600;
  display: block;
  padding: 6px 0;
  border-radius: 0;
  background: transparent;
  border: none;
}

.title-list a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .sidebar-card {
    padding: var(--mobile-space-3);
  }

  .title-list {
    grid-auto-flow: row;
    overflow-x: visible;
    padding-bottom: 0;
  }
}

</style>
