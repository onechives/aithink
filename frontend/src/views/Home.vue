<template>
  <BlogLayout>
    <template #sidebar>
      <SidebarList />
    </template>
    <section class="content">
      <div class="toolbar">
        <div>
          <h2>{{ $t('home.title') }}</h2>
          <p>{{ $t('home.subtitle') }}</p>
        </div>
        <div class="sort-buttons">
          <button :class="{ active: sort === 'latest' }" @click="setSort('latest')">{{ $t('home.latest') }}</button>
          <button :class="{ active: sort === 'likes' }" @click="setSort('likes')">{{ $t('home.mostLiked') }}</button>
        </div>
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
  </BlogLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import BlogLayout from "../components/BlogLayout.vue";
import SidebarList from "../components/SidebarList.vue";
import PostCard from "../components/PostCard.vue";
import { likePost, listPosts, type PostSummary } from "../api/posts";

const posts = ref<PostSummary[]>([]);
const total = ref(0);
const page = ref(1);
const size = 6;
const sort = ref<"latest" | "likes">("latest");

async function fetchPosts() {
  // 拉取首页文章列表
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
  // 上一页
  if (page.value > 1) {
    page.value -= 1;
    fetchPosts();
  }
}

function nextPage() {
  // 下一页
  const maxPage = Math.ceil(total.value / size);
  if (page.value < maxPage) {
    page.value += 1;
    fetchPosts();
  }
}

async function handleLike(id: string) {
  // 点赞后本地更新数值
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
.content {
  display: flex;
  flex-direction: column;
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

.sort-buttons {
  display: flex;
  gap: 10px;
}

.sort-buttons button {
  border-radius: 10px;
  border: none;
  background: var(--surface);
  padding: 10px 18px;
  cursor: pointer;
  color: var(--ink);
}

.sort-buttons button.active {
  background: var(--accent);
  color: var(--accent-contrast);
}

.cards {
  display: grid;
  gap: 22px;
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


@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: flex-start;
  }

  .sort-buttons {
    width: 100%;
  }

  .sort-buttons button {
    flex: 1;
  }

}

</style>
