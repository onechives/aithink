<template>
  <section class="mobile-search">
    <input v-model="keyword" type="text" :placeholder="`${$t('common.search')}...`" @input="emitSearch" />
    <ul v-if="keyword.trim()" class="results">
      <li v-for="item in items" :key="item.id">
        <RouterLink :to="`/m/posts/${item.id}`">{{ item.title }}</RouterLink>
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { listPostTitles, type PostTitle } from "../../api/posts";

const items = ref<PostTitle[]>([]);
const keyword = ref("");

async function fetchTitles() {
  items.value = await listPostTitles({
    sort: "latest",
    keyword: keyword.value,
    size: 50,
  });
}

let typingTimer: number | undefined;
function emitSearch() {
  window.clearTimeout(typingTimer);
  typingTimer = window.setTimeout(() => {
    if (!keyword.value.trim()) {
      items.value = [];
      return;
    }
    fetchTitles();
  }, 300);
}

onMounted(() => {
  items.value = [];
});
</script>

<style scoped lang="scss">
.mobile-search {
  display: grid;
  gap: var(--mobile-space-2);
}

input {
  border-radius: 10px;
  border: 1px solid var(--border);
  padding: 12px 14px;
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
}

.results {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 10px;
}

.results a {
  text-decoration: none;
  color: var(--ink);
  font-weight: 600;
  display: block;
  padding: 10px 12px;
  border-radius: 12px;
  background: var(--surface);
  border: 1px solid var(--border);
}

.results a:hover {
  background: var(--surface-alt);
}
</style>
