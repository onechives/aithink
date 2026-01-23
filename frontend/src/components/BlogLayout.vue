<template>
  <section class="blog-layout" :class="{ 'blog-layout--narrow': sidebarNarrow }">
    <aside class="sidebar">
      <slot name="sidebar" />
    </aside>
    <main class="main">
      <slot />
    </main>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  // 控制左侧栏宽度是否收窄（文章详情页使用）
  sidebarNarrow?: boolean;
}>();
</script>

<style scoped lang="scss">
.blog-layout {
  display: grid;
  width: 100%;
  min-width: 100%;
  grid-template-columns: minmax(220px, 22%) 1fr;
  gap: 28px;
  padding: 28px 36px 48px;
}

.blog-layout--narrow {
  grid-template-columns: minmax(176px, 17.6%) 1fr;
}

.sidebar {
  position: sticky;
  top: 24px;
  align-self: start;
  min-width: 0;
}

.main {
  min-height: 60vh;
  min-width: 0;
}

@media (max-width: 960px) {
  .blog-layout {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: static;
  }
}

@media (max-width: 768px) {
  .blog-layout {
    padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
    gap: var(--mobile-space-4);
  }
}

</style>
