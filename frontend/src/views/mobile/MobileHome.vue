<template>
  <section class="mobile-home">
    <header class="hero">
      <h1>{{ $t('home.title') }}</h1>
      <p>{{ $t('home.subtitle') }}</p>
    </header>

    <MobileSearch />

    <div class="sort-buttons">
      <button :class="{ active: sort === 'latest' }" @click="setSort('latest')">{{ $t('home.latest') }}</button>
      <button :class="{ active: sort === 'likes' }" @click="setSort('likes')">{{ $t('home.mostLiked') }}</button>
    </div>

    <div class="cards">
      <PostCard v-for="post in posts" :key="post.id" :post="post" @like="handleLike" />
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
import MobileSearch from "../../components/mobile/MobileSearch.vue";
import PostCard from "../../components/PostCard.vue";
import { likePost, listPosts, type PostSummary } from "../../api/posts";

const posts = ref<PostSummary[]>([]);
const total = ref(0);
const page = ref(1);
const size = 6;
const sort = ref<"latest" | "likes">("latest");

async function fetchPosts() {
  // 拉取移动端首页文章列表
  const data = await listPosts({
    sort: sort.value,
    keyword: "",
    page: page.value,
    size,
  });
  posts.value = data.items;
  total.value = data.total;
}

function setSort(value: "latest" | "likes") {
  // 切换排序并重置分页
  sort.value = value;
  page.value = 1;
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

async function handleLike(id: string) {
  // 点赞并更新本地计数
  const result = await likePost(id);
  const target = posts.value.find((item) => item.id === id);
  if (target) {
    target.likeCount = result.likeCount;
  }
}

onMounted(() => {
  fetchPosts();
});
</script>

<style scoped lang="scss">
.mobile-home {
  display: grid;
  gap: var(--mobile-space-4);
  padding: var(--mobile-space-4) var(--mobile-space-3);
}

.hero h1 {
  margin: 0;
  font-size: var(--mobile-font-lg);
}

.hero p {
  margin: var(--mobile-space-1) 0 0;
  color: var(--muted);
}

.sort-buttons {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--mobile-space-2);
}

.sort-buttons button {
  border: none;
  background: var(--surface);
  color: var(--ink);
  border-radius: var(--mobile-radius);
  padding: var(--mobile-space-2);
  font-weight: 600;
  cursor: pointer;
}

.sort-buttons button.active {
  background: var(--accent);
  color: var(--accent-contrast);
}

.cards {
  display: grid;
  gap: var(--mobile-space-3);
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
