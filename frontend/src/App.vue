<template>
  <div class="app" :class="{ 'app--mobile': isMobileRoute }">
    <TopBar v-if="!isMobileRoute" />
    <RouterView />
    <FloatingActions v-if="!isMobileRoute" />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import TopBar from "./components/TopBar.vue";
import FloatingActions from "./components/FloatingActions.vue";

const route = useRoute();
// 仅 /m 前缀走移动端壳，其他走 PC 布局
const isMobileRoute = computed(() => route.path === "/m" || route.path.startsWith("/m/"));
</script>

<style scoped lang="scss">
.app {
  min-height: 100vh;
  width: 100%;
  min-width: 100%;
  background: var(--page-bg);
}

.app--mobile {
  background: var(--page-bg);
}
</style>
