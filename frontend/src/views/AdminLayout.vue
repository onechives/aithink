<template>
  <section class="admin-layout">
    <nav class="admin-tabs" aria-label="Admin tabs">
      <RouterLink class="tab" :class="{ active: isActive('admin-posts') }" to="/admin/posts">
        {{ $t('admin.postManage') }}
      </RouterLink>
      <RouterLink class="tab" :class="{ active: isActive('admin-review-posts') }" to="/admin/review/posts">
        {{ $t('admin.postReview') }}
      </RouterLink>
      <RouterLink class="tab" :class="{ active: isActive('admin-review-users') }" to="/admin/review/users">
        {{ $t('admin.userReview') }}
      </RouterLink>
      <RouterLink class="tab" :class="{ active: isActive('admin-review-nicknames') }" to="/admin/review/nicknames">
        {{ $t('admin.nicknameReview') }}
      </RouterLink>
    </nav>
    <RouterView />
  </section>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router";

const route = useRoute();

function isActive(name: string) {
  return route.name === name;
}
</script>

<style scoped lang="scss">
.admin-layout {
  display: block;
}

.admin-tabs {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 24px 36px 0;
  flex-wrap: nowrap;
  overflow-x: auto;
}

.admin-tabs::-webkit-scrollbar {
  display: none;
}

.tab {
  text-decoration: none;
  color: var(--muted);
  font-weight: 600;
  background: var(--surface-alt);
  padding: 10px 18px;
  border-radius: 999px;
  border: 1px solid transparent;
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease;
  white-space: nowrap;
}

.tab.active {
  color: var(--ink);
  background: var(--surface);
  border-color: var(--border);
  box-shadow: var(--shadow);
}

@media (max-width: 960px) {
  .admin-tabs {
    padding: 24px 24px 0;
  }
}
</style>
