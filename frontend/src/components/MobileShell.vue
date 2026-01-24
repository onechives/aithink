<template>
  <div class="mobile-shell">
    <MobileTopBar />
    <main class="mobile-body">
      <div class="mobile-body__inner">
        <RouterView />
      </div>
    </main>
    <div class="mobile-fabs">
      <RouterLink v-if="auth.token" class="message-fab" to="/m/messages" :aria-label="$t('nav.myMessages')">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path
            d="M4.5 7.2C4.5 4.6 7 2.5 10.2 2.5h3.6c3.2 0 5.7 2.1 5.7 4.7v4.7c0 2.6-2.5 4.7-5.7 4.7h-1.2l-3.6 2.2a0.9 0.9 0 0 1-1.4-0.8v-1.7H10.2c-3.2 0-5.7-2.1-5.7-4.7V7.2z"
            fill="currentColor" />
        </svg>
        <span v-if="unreadCount > 0" class="fab-badge">{{ unreadCount }}</span>
      </RouterLink>
      <button v-if="showBackToTop" class="back-to-top" type="button" @click="scrollToTop"
        :aria-label="$t('home.backToTop')">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path
            d="M12 5.5a1 1 0 0 1 .7.3l6 6a1 1 0 1 1-1.4 1.4L13 8.9V19a1 1 0 1 1-2 0V8.9L6.7 13.2a1 1 0 1 1-1.4-1.4l6-6a1 1 0 0 1 .7-.3z"
            fill="currentColor" />
        </svg>
      </button>
    </div>
    <MobileNav v-if="!hideMobileNav" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import MobileNav from "./MobileNav.vue";
import MobileTopBar from "./MobileTopBar.vue";
import { useAuthStore } from "../stores/auth";
import { getUnreadCount } from "../api/messages";

const route = useRoute();
const auth = useAuthStore();
// meta.hideMobileNav 用于登录/注册等不展示底部栏的页面
const hideMobileNav = computed(() => route.meta?.hideMobileNav === true);
const showBackToTop = ref(false);
const unreadCount = ref(0);
const unreadEvent = "messages-updated";

function handleScroll() {
  // 滚动一定距离才显示返回顶部
  showBackToTop.value = window.scrollY > 200;
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: "smooth" });
}

async function fetchUnreadCount() {
  if (!auth.token) {
    unreadCount.value = 0;
    return;
  }
  try {
    const data = await getUnreadCount();
    unreadCount.value = data.count;
  } catch {
    unreadCount.value = 0;
  }
}

onMounted(() => {
  handleScroll();
  window.addEventListener("scroll", handleScroll, { passive: true });
  fetchUnreadCount();
  window.addEventListener(unreadEvent, fetchUnreadCount);
});

watch(
  () => auth.token,
  () => {
    fetchUnreadCount();
  }
);

onBeforeUnmount(() => {
  window.removeEventListener("scroll", handleScroll);
  window.removeEventListener(unreadEvent, fetchUnreadCount);
});
</script>

<style scoped lang="scss">
.mobile-shell {
  min-height: 100vh;
  background: var(--page-bg);
  width: 100%;
  overflow-x: hidden;
}

.mobile-body {
  padding-top: calc(var(--mobile-unit) * 9 + env(safe-area-inset-top));
  padding-bottom: calc(var(--mobile-space-14) + env(safe-area-inset-bottom));
  width: 100%;
  min-width: 0;
  overflow-x: hidden;
}

.mobile-body__inner {
  width: 100%;
  max-width: 720px;
  margin: 0 auto;
  min-width: 0;
}

.mobile-fabs {
  position: fixed;
  right: var(--mobile-space-3);
  bottom: calc(var(--mobile-space-14) + var(--mobile-space-3) + env(safe-area-inset-bottom));
  display: grid;
  gap: var(--mobile-space-2);
  z-index: 1200;
}

.back-to-top,
.message-fab {
  width: calc(var(--mobile-unit) * 10);
  height: calc(var(--mobile-unit) * 10);
  border-radius: 999px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--ink);
  display: grid;
  place-items: center;
  box-shadow: var(--shadow);
  cursor: pointer;
  padding: 0;
  text-decoration: none;
}

.back-to-top svg,
.message-fab svg {
  width: calc(var(--mobile-unit) * 6);
  height: calc(var(--mobile-unit) * 6);
}

.message-fab {
  position: relative;
}

.fab-badge {
  position: absolute;
  top: -3px;
  right: -3px;
  min-width: 14px;
  height: 14px;
  padding: 0 4px;
  border-radius: 999px;
  background: #e24b2b;
  color: #fff;
  font-size: 9px;
  line-height: 14px;
  text-align: center;
}
</style>
