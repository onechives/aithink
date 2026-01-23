<template>
  <div class="mobile-shell">
    <MobileTopBar />
    <main class="mobile-body">
      <RouterView />
    </main>
    <button v-if="showBackToTop" class="back-to-top" type="button" @click="scrollToTop"
      :aria-label="$t('home.backToTop')">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path
          d="M12 5.5a1 1 0 0 1 .7.3l6 6a1 1 0 1 1-1.4 1.4L13 8.9V19a1 1 0 1 1-2 0V8.9L6.7 13.2a1 1 0 1 1-1.4-1.4l6-6a1 1 0 0 1 .7-.3z"
          fill="currentColor" />
      </svg>
    </button>
    <MobileNav v-if="!hideMobileNav" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import MobileNav from "./MobileNav.vue";
import MobileTopBar from "./MobileTopBar.vue";

const route = useRoute();
// meta.hideMobileNav 用于登录/注册等不展示底部栏的页面
const hideMobileNav = computed(() => route.meta?.hideMobileNav === true);
const showBackToTop = ref(false);

function handleScroll() {
  // 滚动一定距离才显示返回顶部
  showBackToTop.value = window.scrollY > 200;
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: "smooth" });
}

onMounted(() => {
  handleScroll();
  window.addEventListener("scroll", handleScroll, { passive: true });
});

onBeforeUnmount(() => {
  window.removeEventListener("scroll", handleScroll);
});
</script>

<style scoped lang="scss">
.mobile-shell {
  min-height: 100vh;
  background: var(--page-bg);
}

.mobile-body {
  padding-top: calc(var(--mobile-unit) * 9 + env(safe-area-inset-top));
  padding-bottom: calc(var(--mobile-space-14) + env(safe-area-inset-bottom));
}

.back-to-top {
  position: fixed;
  right: var(--mobile-space-3);
  bottom: calc(var(--mobile-space-14) + var(--mobile-space-3) + env(safe-area-inset-bottom));
  width: calc(var(--mobile-unit) * 12);
  height: calc(var(--mobile-unit) * 12);
  border-radius: 999px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--ink);
  display: grid;
  place-items: center;
  box-shadow: var(--shadow);
  cursor: pointer;
  padding: 0;
  z-index: 1200;
}

.back-to-top svg {
  width: calc(var(--mobile-unit) * 5);
  height: calc(var(--mobile-unit) * 5);
}
</style>
