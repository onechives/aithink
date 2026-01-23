<template>
  <header class="topbar">
    <div class="brand">
      <div>
        <h1>{{ $t('topbar.title') }}</h1>
      </div>
    </div>
    <nav class="actions">
      <button class="theme-toggle" type="button" @click="toggleTheme"
        :aria-label="theme === 'dark' ? $t('topbar.switchToLight') : $t('topbar.switchToDark')">
        <svg class="icon" viewBox="0 0 1024 1024" aria-hidden="true">
          <path
            d="M512.006138 0C232.124098 0 5.125033 229.233211 5.125033 512s226.937687 512 506.881105 512 506.868829-229.233211 506.868829-512S791.93728 0 512.006138 0z m0 919.252919c-229.049078 0-414.728716-185.679637-414.728716-414.71644S282.95706 89.820039 512.006138 89.820039s414.71644 185.667362 414.71644 414.71644S741.04294 919.252919 512.006138 919.252919z" />
          <path
            d="M490.352106 148.902155c192.271596 0 348.158527 155.899206 348.158527 348.158527s-155.899206 348.158527-348.158527 348.158527" />
        </svg>
      </button>

      <div class="language-selector">
        <button class="lang-button" type="button" @click="showLanguageMenu = !showLanguageMenu"
          :aria-label="$t('topbar.language')">
          <svg class="lang-icon" viewBox="0 0 1024 1024" aria-hidden="true">
            <path
              d="M277.4528 328.2944h-12.9024L240.64 427.264h60.8256l-24.0128-98.9696h0.0512z m485.2736 197.888c13.0048 28.672 30.208 51.9168 49.3056 71.68 19.0464-19.7632 38.2976-43.008 51.2-71.68h-100.5056z" />
            <path
              d="M934.0416 252.416h-372.6336l77.7216 513.9968c1.6896 21.1968-6.5024 41.984-22.2208 56.32l-131.584 123.9552h448.7168c49.6128 0 89.9584-33.28 89.9584-74.24V328.2944c0-40.96-40.3968-75.8784-89.9584-75.8784z m0 273.7664h-8.0384c-17.0496 45.1584-44.1344 80.4864-72.2432 108.3392 22.016 16.64 45.568 30.208 68.9664 45.5168 12.9536 8.4992 15.0528 24.064 4.6592 34.7648-10.24 10.6496-29.2352 12.3904-42.1376 3.84-25.4464-16.5888-49.3056-30.464-73.216-48.5376-23.9104 18.0736-45.824 31.9488-71.2704 48.5376-12.9024 8.5504-31.8464 6.8096-42.1888-3.84-10.3424-10.7008-8.192-26.2656 4.6592-34.816 23.4496-15.2064 45.056-28.8256 67.072-45.4656a300.1856 300.1856 0 0 1-70.2976-108.3392h-7.9872c-16.64 0-30.0544-11.0592-30.0544-24.7296s13.4656-24.7296 30.0544-24.7296h89.9584v-24.7808c0-13.6192 13.4656-24.6784 30.0544-24.6784 16.5376 0 29.952 11.0592 29.952 24.6784v24.7808h92.0064c16.5888 0 30.0032 11.0592 30.0032 24.7296s-13.4144 24.7296-29.952 24.7296z" />
            <path
              d="M488.2944 167.424C482.7136 130.2528 444.3136 102.4 399.0528 102.4H90.0096C40.3456 102.4 0 135.6288 0 176.5888V723.968c0 40.96 40.3968 74.24 90.0096 74.24h473.088c8.7552-8.2944 16.128-13.4656 16.384-24.2176 0.1024-2.7136-90.7264-604.0064-91.136-606.72v0.0512z m-120.4224 407.7056c-15.872 2.7648-32.0512-5.7856-35.328-19.3536l-19.1488-79.0528H228.5568l-19.1488 79.0528c-3.2256 13.3632-18.8416 22.1696-35.328 19.3536-16.1792-2.6112-26.7264-15.6672-23.5008-29.0816l60.0064-247.296c2.816-11.5712 15.104-19.9168 29.44-19.9168h61.952c14.336 0 26.624 8.3456 29.44 19.9168l60.0064 247.296c3.2256 13.4656-7.3216 26.4704-23.552 29.184v-0.1024z m-17.3568 272.5888l5.12 33.9456c3.4816 22.7328 21.8112 45.8752 51.9168 57.4976l100.7104-91.4432H350.5152z" />
          </svg>
        </button>
        <div v-if="showLanguageMenu" class="language-menu">
          <button v-for="(label, code) in languages" :key="code" type="button" :class="{ active: locale === code }"
            @click="changeLanguage(code)">
            {{ label }}
          </button>
        </div>
      </div>

      <RouterLink to="/">{{ $t('nav.articles') }}</RouterLink>
      <RouterLink v-if="auth.role === 'admin'" to="/admin/posts">{{ $t('nav.manage') }}</RouterLink>
      <RouterLink v-if="auth.token" to="/write">{{ $t('nav.write') }}</RouterLink>
      <RouterLink v-if="auth.token" to="/me/settings">{{ $t('nav.settings') }}</RouterLink>
      <RouterLink v-if="!auth.token" to="/login">{{ $t('nav.login') }}</RouterLink>
      <RouterLink v-if="!auth.token" to="/register">{{ $t('nav.register') }}</RouterLink>
      <button v-if="auth.token" class="logout" type="button" @click="handleLogout" :aria-label="$t('nav.logout')">
        {{ $t('nav.logout') }}
      </button>

    </nav>
  </header>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref } from "vue";
import { useTheme } from "../composables/useTheme";
import { useAuthStore } from "../stores/auth";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";

const { theme, toggleTheme } = useTheme();
const auth = useAuthStore();
const router = useRouter();
const { locale } = useI18n();
const showLanguageMenu = ref(false);

const languages = {
  "zh-CN": "中文",
  "en-US": "EN",
};

function changeLanguage(code: string) {
  // 切换语言并写入本地存储
  locale.value = code;
  localStorage.setItem("language", code);
  showLanguageMenu.value = false;
}

function handleLogout() {
  // 清理登录态后返回首页
  auth.clearAuth();
  router.push({ name: "home" });
}

onBeforeUnmount(() => {
  showLanguageMenu.value = false;
});
</script>

<style scoped lang="scss">
.topbar {

  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-left: 10px;
  padding-right: 20px;
  padding-top: 3px;
  padding-bottom: 3px;

  border-bottom: 1px solid var(--border);
  background: var(--topbar-bg);
  position: sticky;
  top: 0;
  z-index: 1200;
}

.brand {
  display: flex;
  align-items: center;
}

.logo {
  width: 44px;
  height: 44px;
  border-radius: 16px;
  background: var(--accent);
  color: var(--accent-contrast);
  display: grid;
  place-items: center;
  font-weight: 700;
  font-size: 20px;
  letter-spacing: 0.04em;
}

h1 {
  font-size: 20px;
  margin-left: 40px;
  line-height: 1;
}

p {
  margin: 4px 0 0;
  color: var(--muted);
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 18px;
  font-weight: 600;
  font-family: inherit;
  text-transform: none;
  letter-spacing: 0;
  font-size: 12px;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
  margin-left: auto;
}

.actions a,
.actions button {
  color: var(--ink);
  text-decoration: none;
  padding: 6px 4px;
  border-radius: 0;
  background: transparent;
  border: none;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  white-space: nowrap;
}

.actions a.router-link-active {
  color: var(--accent);
  border-bottom: 2px solid var(--accent);
}

.nav-item {
  position: relative;
}

.badge {
  position: absolute;
  top: -6px;
  right: -6px;
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

.theme-toggle {
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.theme-toggle .icon {
  width: 16px;
  height: 16px;
  fill: currentColor;
}

.language-selector {
  position: relative;
}

.lang-button {
  width: auto;
  padding: 6px 8px;
  font-size: 12px;
}

.lang-icon {
  width: 18px;
  height: 18px;
  display: block;
  fill: currentColor;
}

.language-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 4px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
  min-width: 120px;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.language-menu button {
  width: 100%;
  padding: 10px 16px;
  border: none;
  background: transparent;
  color: var(--ink);
  cursor: pointer;
  text-align: left;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.2s;

  &:hover {
    background: var(--surface-alt);
  }

  &.active {
    color: var(--accent);
  }
}

@media (max-width: 768px) {
  .topbar {
    padding: var(--mobile-space-2) var(--mobile-space-3);
    flex-direction: row;
    align-items: center;
    gap: var(--mobile-space-2);
  }

  .actions {
    width: 100%;
    justify-content: flex-end;
    flex-wrap: nowrap;
    gap: var(--mobile-space-2);
  }

  h1 {
    font-size: var(--mobile-font-lg);
  }

  .actions a {
    display: none;
  }
}
</style>
