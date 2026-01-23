<template>
  <header class="mobile-topbar">
    <button class="back" type="button" @click="goBack" :aria-label="$t('common.back')">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path
          d="M14.7 5.3a1 1 0 0 1 0 1.4L9.4 12l5.3 5.3a1 1 0 1 1-1.4 1.4l-6-6a1 1 0 0 1 0-1.4l6-6a1 1 0 0 1 1.4 0z"
          fill="currentColor" />
      </svg>
    </button>
    <span class="title">{{ title }}</span>
    <span class="spacer" aria-hidden="true" />
  </header>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";

const route = useRoute();
const router = useRouter();
const { t } = useI18n();

const title = computed(() => {
  // 根据路由 meta.titleKey 渲染标题
  const key = route.meta?.titleKey as string | undefined;
  return key ? t(key) : "";
});

function goBack() {
  // 优先浏览器历史返回，否则回到移动端首页
  if (window.history.length > 1) {
    router.back();
    return;
  }
  const target = route.path.startsWith("/m") ? "m-home" : "home";
  router.push({ name: target });
}
</script>

<style scoped lang="scss">
.mobile-topbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1200;
  height: calc(var(--mobile-unit) * 9 + env(safe-area-inset-top));
  display: grid;
  grid-template-columns: calc(var(--mobile-unit) * 8) 1fr calc(var(--mobile-unit) * 8);
  align-items: center;
  padding: env(safe-area-inset-top) var(--mobile-space-2) 0;
  background: var(--topbar-bg);
  border-bottom: 1px solid var(--border);
}

.back {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: calc(var(--mobile-unit) * 7);
  height: calc(var(--mobile-unit) * 7);
  border: none;
  background: transparent;
  color: var(--ink);
  padding: 0;
  cursor: pointer;
}

.back svg {
  width: calc(var(--mobile-unit) * 4);
  height: calc(var(--mobile-unit) * 4);
}

.title {
  text-align: center;
  font-weight: 600;
  font-size: var(--mobile-font-md);
  letter-spacing: calc(var(--mobile-unit) * 0.1);
}

.spacer {
  display: block;
  width: calc(var(--mobile-unit) * 7);
  height: calc(var(--mobile-unit) * 7);
}
</style>
