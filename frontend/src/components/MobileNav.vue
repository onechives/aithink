<template>
  <nav class="mobile-nav" :style="{ '--nav-count': navItems.length }">
    <RouterLink v-for="item in navItems" :key="item.key" :to="item.to" custom v-slot="{ href, navigate }">
      <a :href="href" class="nav-item" :class="{ 'nav-item--active': isNavActive(item.to) }" @click="navigate">
        <span>{{ item.label }}</span>
      </a>
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";

const { t } = useI18n();
const route = useRoute();

const navItems = computed(() => {
  // 移动端底部导航入口
  return [
    { key: "home", label: t("nav.articles"), to: "/m" },
    { key: "write", label: t("nav.write"), to: "/m/write" },
    { key: "settings", label: t("nav.settings"), to: "/m/settings" },
  ];
});

const isNavActive = (to: string) => {
  // /m 首页需要严格匹配，其他按前缀匹配
  if (to === "/m") return route.path === "/m";
  return route.path.startsWith(to);
};
</script>

<style scoped lang="scss">
.mobile-nav {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  --nav-count: 5;
  padding: var(--mobile-space-2) var(--mobile-space-2) calc(var(--mobile-space-2) + env(safe-area-inset-bottom));
  background: var(--surface);
  border-top: 1px solid var(--border);
  gap: var(--mobile-space-1);
  align-items: center;
  justify-content: space-between;
  z-index: 1300;
}

.nav-item {
  color: var(--muted);
  text-decoration: none;
  font-weight: 600;
  font-size: calc(var(--viewport-w) / (var(--nav-count) * 10));
  text-transform: uppercase;
  letter-spacing: calc(var(--mobile-unit) * 0.05);
  padding: calc(var(--mobile-space-1) * 0.6) calc(var(--mobile-space-1) * 0.3);
  text-align: center;
  white-space: nowrap;
  line-height: 1.1;
  min-width: 0;
  flex: 1 1 0;
  max-width: calc((var(--viewport-w) - (var(--mobile-space-2) * 2)) / var(--nav-count));
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-item--active {
  color: var(--accent);
}

</style>
