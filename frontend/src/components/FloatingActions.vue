<template>
  <div class="floating-actions">
    <button class="fab" type="button" @click="scrollToTop" :aria-label="$t('home.backToTop')">
      <svg class="fab-icon fab-icon-top" viewBox="0 0 1024 1024" aria-hidden="true">
        <path
          d="M512 128c-211.7 0-384 172.3-384 384s172.3 384 384 384 384-172.3 384-384-172.3-384-384-384z m0 717.4c-183.8 0-333.4-149.6-333.4-333.4S328.2 178.6 512 178.6 845.4 328.2 845.4 512 695.8 845.4 512 845.4z"
          fill="currentColor" />
        <path
          d="M540 414.8c-7.5-7.5-17.4-11.6-28-11.6s-20.5 4.1-28 11.6L321.2 577.6c-9.9 9.9-9.9 25.9 0 35.8 9.9 9.9 25.9 9.9 35.8 0l155-155 155 155c4.9 4.9 11.4 7.4 17.9 7.4s13-2.5 17.9-7.4c9.9-9.9 9.9-25.9 0-35.8L540 414.8z"
          fill="currentColor" />
      </svg>
    </button>
    <RouterLink v-if="auth.token" class="fab secondary fab-message" to="/me/messages" :aria-label="$t('nav.myMessages')">
      <svg class="fab-icon fab-icon-message" viewBox="0 0 1024 1024" aria-hidden="true">
        <path
          d="M523.946667 85.333333C802.773333 85.333333 982.826667 307.541333 982.826667 511.338667c0 242.837333-197.397333 448-458.837334 448-84.16 0-150.464-16.597333-221.610666-55.402667l-88.618667 50.197333c-27.797333 8.426667-50.474667 5.162667-68.010667-9.834666-17.557333-14.997333-25.130667-34.730667-22.72-59.2a19570.688 19570.688 0 0 0 29.653334-106.368C123.008 741.76 64 656.426667 64 511.338667 64 307.541333 245.141333 85.333333 523.946667 85.333333z m-1.28 64C304.064 149.333333 128 317.12 128 522.666667c0 77.354667 24.874667 151.125333 70.634667 213.184l5.397333 7.125333 18.218667 23.509333-36.970667 128.746667a0.32 0.32 0 0 0 0.490667 0.384l113.408-63.829333 26.752 14.592C385.237333 878.72 452.544 896 522.666667 896 741.269333 896 917.333333 728.213333 917.333333 522.666667S741.269333 149.333333 522.666667 149.333333z m-192 320a53.333333 53.333333 0 1 1 0 106.666667 53.333333 53.333333 0 0 1 0-106.666667z m182.848 0a53.333333 53.333333 0 1 1 0 106.666667 53.333333 53.333333 0 0 1 0-106.666667z m182.869333 0a53.333333 53.333333 0 1 1 0 106.666667 53.333333 53.333333 0 0 1 0-106.666667z"
          fill="currentColor" />
      </svg>
      <span v-if="unreadCount > 0" class="fab-badge">{{ unreadCount }}</span>
    </RouterLink>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useAuthStore } from "../stores/auth";
import { getUnreadCount } from "../api/messages";

const auth = useAuthStore();
const unreadCount = ref(0);
const unreadEvent = "messages-updated";

async function fetchUnreadCount() {
  // 未登录时不请求
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

function scrollToTop() {
  // 平滑回到顶部
  window.scrollTo({ top: 0, behavior: "smooth" });
}

onMounted(() => {
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
  window.removeEventListener(unreadEvent, fetchUnreadCount);
});
</script>

<style scoped lang="scss">
.floating-actions {
  position: fixed;
  right: 24px;
  bottom: 24px;
  display: grid;
  gap: 10px;
  z-index: 1200;
}

.fab {
  width: 44px;
  height: 44px;
  border-radius: 0;
  background: transparent;
  color: var(--ink);
  border: none;
  cursor: pointer;
  display: grid;
  place-items: center;
  text-decoration: none;
  box-shadow: none;
  padding: 0;
  box-sizing: border-box;
}

.fab.secondary {
  background: transparent;
  color: var(--ink);
  border: none;
}

.fab-icon {
  width: 39px;
  height: 39px;
  transition: color 0.2s ease;
}

.fab-icon-top {
  width: 45px;
  height: 45px;
}

.fab-icon-message {
  width: 39px;
  height: 39px;
}

.fab-message {
  position: relative;
}

.fab-badge {
  position: absolute;
  top: -2px;
  right: -2px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  border-radius: 999px;
  background: #e24b2b;
  color: #fff;
  font-size: 10px;
  line-height: 16px;
  text-align: center;
}

@media (max-width: 768px) {
  .floating-actions {
    right: var(--mobile-space-3);
    bottom: var(--mobile-space-3);
  }

  .fab {
    width: calc(var(--mobile-unit) * 14);
    height: calc(var(--mobile-unit) * 14);
    font-size: var(--mobile-font-xs);
  }
}
</style>
