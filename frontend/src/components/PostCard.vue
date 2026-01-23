<template>
  <article class="card" @click="goDetail">
    <div class="cover" v-if="post.coverUrl" :style="{ backgroundImage: `url(${post.coverUrl})` }" />
    <div class="content">
      <div class="meta">
        <span>{{ post.author }}</span>
        <span class="dot">•</span>
        <span>{{ post.category }}</span>
        <span class="dot">•</span>
        <span>{{ formatDate(post.createdAt) }}</span>
      </div>
      <h3>{{ post.title }}</h3>
      <p>{{ post.summary }}</p>
      <div class="tags" v-if="post.tags">
        <span v-for="tag in tags" :key="tag">{{ tag }}</span>
      </div>
      <div class="actions">
        <button class="like" @click.stop="$emit('like', post.id)">
          <svg class="icon" viewBox="0 0 1024 1024" aria-hidden="true">
            <path
              d="M621.674667 408.021333c16.618667-74.24 28.224-127.936 34.837333-161.194666C673.152 163.093333 629.941333 85.333333 544.298667 85.333333c-77.226667 0-116.010667 38.378667-138.88 115.093334l-0.586667 2.24c-13.728 62.058667-34.72 110.165333-62.506667 144.586666a158.261333 158.261333 0 0 1-119.733333 58.965334l-21.909333 0.469333C148.437333 407.808 106.666667 450.816 106.666667 503.498667V821.333333c0 64.8 52.106667 117.333333 116.394666 117.333334h412.522667c84.736 0 160.373333-53.568 189.12-133.92l85.696-239.584c21.802667-60.96-9.536-128.202667-70.005333-150.186667a115.552 115.552 0 0 0-39.488-6.954667H621.674667zM544.256 149.333333c39.253333 0 59.498667 36.48 49.888 84.928-7.573333 38.144-21.984 104.426667-43.221333 198.666667-4.512 20.021333 10.56 39.093333 30.912 39.093333h218.666666c6.101333 0 12.16 1.066667 17.909334 3.168 27.445333 9.984 41.674667 40.554667 31.776 68.266667l-85.568 239.573333C744.981333 838.026667 693.301333 874.666667 635.402667 874.666667H223.498667C194.314667 874.666667 170.666667 850.784 170.666667 821.333333V503.498667c0-17.866667 14.144-32.448 31.829333-32.821334l21.866667-0.469333a221.12 221.12 0 0 0 167.381333-82.56c34.346667-42.602667 59.146667-99.306667 74.869333-169.877333C482.101333 166.336 499.552 149.333333 544.266667 149.333333z" />
          </svg>
          <span>{{ post.likeCount }}</span>
        </button>
        <RouterLink class="read" :to="`/posts/${post.id}`" @click.stop>{{ $t('home.readMore') }}</RouterLink>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { PostSummary } from "../api/posts";
import { useRouter } from "vue-router";

const props = defineProps<{ post: PostSummary }>();
const router = useRouter();

const tags = props.post.tags ? props.post.tags.split(",").map((t) => t.trim()).filter(Boolean) : [];

function formatDate(value: string) {
  // 统一日期展示格式
  return new Date(value).toLocaleDateString();
}

function goDetail() {
  // 卡片整体点击进入详情页
  router.push(`/posts/${props.post.id}`);
}
</script>

<style scoped lang="scss">
.card {
  display: grid;
  grid-template-columns: minmax(140px, 220px) 1fr;
  gap: 20px;
  padding: 20px;
  border-radius: 22px;
  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: var(--shadow);
  cursor: pointer;
}

.cover {
  border-radius: 18px;
  background-size: cover;
  background-position: center;
  min-height: 160px;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--muted);
}

.dot {
  font-size: 14px;
}

h3 {
  margin: 0;
  font-size: 22px;
}

p {
  margin: 0;
  color: var(--muted);
  line-height: 1.6;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tags span {
  background: var(--chip);
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
}

.actions {
  margin-top: auto;
  display: flex;
  align-items: center;
  gap: 12px;
}

.actions a {
  text-decoration: none;
  font-weight: 600;
  color: var(--ink);
}

.actions .read {
  margin-left: auto;
}

.actions .read {
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-size: 11px;
}

.actions button {
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

.actions .icon {
  width: 16px;
  height: 16px;
  fill: currentColor;
}

@media (max-width: 768px) {
  .card {
    grid-template-columns: 1fr;
  }

  .cover {
    min-height: calc(var(--mobile-unit) * 50);
  }
}

</style>
