<template>
  <section class="mobile-settings">
    <div class="mobile-section">
      <h3 class="section-title">{{ $t('settings.sectionContent') }}</h3>
      <div class="section-card">
        <button class="mobile-item" type="button" @click="router.push({ name: 'm-messages' })">
          {{ $t('messages.title') }}
        </button>
      </div>
      <div class="section-card">
        <button class="mobile-item" type="button" @click="router.push({ name: 'm-my-posts' })">
          {{ $t('myPosts.title') }}
        </button>
      </div>
    </div>

    <div class="mobile-section">
      <h3 class="section-title">{{ $t('settings.sectionPreferences') }}</h3>
      <div class="section-card">
        <button class="mobile-item no-chevron" type="button" @click="toggleTheme">
          <span>{{ $t('topbar.theme') }}</span>
          <span class="value">{{ theme === 'dark' ? $t('common.dark') : $t('common.light') }}</span>
        </button>
      </div>
      <div class="section-inline">
        <span class="inline-label">{{ $t('topbar.language') }}</span>
        <div class="lang-switch">
          <button v-for="(label, code) in languages" :key="code" type="button"
            :class="{ active: locale === code }" @click="changeLanguage(code)">
            {{ label }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="auth.role === 'admin'" class="mobile-section">
      <h3 class="section-title">{{ $t('settings.sectionAdmin') }}</h3>
      <div class="section-card">
        <button class="mobile-item" type="button" @click="router.push({ name: 'admin-posts' })">
          {{ $t('nav.manage') }}
        </button>
      </div>
    </div>

    <div class="account-section">
      <h3 class="section-title">{{ $t('settings.sectionAccount') }}</h3>
      <div class="section-card">
        <button v-if="!isAuthed" class="mobile-item" type="button" @click="router.push({ name: 'm-login' })">
          {{ $t('nav.login') }}
        </button>
        <button v-if="!isAuthed" class="mobile-item" type="button" @click="router.push({ name: 'm-register' })">
          {{ $t('nav.register') }}
        </button>
        <button v-if="isAuthed" class="mobile-item danger" type="button" @click="handleLogout">
          {{ $t('nav.logout') }}
        </button>
      </div>
      <div v-if="isAuthed" class="card">
        <h3>{{ $t('settings.nicknameTitle') }}</h3>
        <p>{{ $t('settings.nicknameDesc') }}</p>
        <div class="nickname">
          <div>
            <span>{{ $t('settings.currentNickname') }}</span>
            <strong>{{ nickname }}</strong>
          </div>
          <label>
            {{ $t('settings.newNicknameHint') }}
            <input v-model="newNickname" type="text" maxlength="10" />
          </label>
          <button @click="handleNickname">{{ $t('settings.submitReview') }}</button>
        </div>
      </div>

      <div v-if="isAuthed" class="card">
        <h3>{{ $t('settings.twoFactorTitle') }}</h3>
        <p v-if="!totpEnabled">{{ $t('settings.twoFactorIntroDisabled') }}</p>
        <p v-else>{{ $t('settings.twoFactorIntroEnabled') }}</p>

        <div v-if="!totpEnabled">
          <button @click="handleInit" :disabled="loading" class="miyao">{{ $t('settings.generateKey') }}</button>
          <div v-if="secret" class="secret">
            <img v-if="qrCode" class="qr" :src="qrCode" alt="2FA QR Code" />
            <div>
              <span>{{ $t('settings.key') }}</span>
              <code>{{ secret }}</code>
            </div>
            <div>
              <span>{{ $t('settings.url') }}</span>
              <code>{{ url }}</code>
            </div>
            <label>
              {{ $t('settings.complete2fa') }}
              <input v-model="code" type="text" maxlength="6" />
            </label>
            <button @click="handleEnable" :disabled="loading">{{ $t('settings.enable') }}</button>
          </div>
        </div>

        <div v-else class="disable">
          <label>
            {{ $t('settings.disable2fa') }}
            <input v-model="code" type="text" maxlength="6" />
          </label>
          <button @click="handleDisable" :disabled="loading">{{ $t('settings.disable') }}</button>
        </div>

        <span class="success" v-if="success">{{ success }}</span>
        <span class="error" v-if="error">{{ error }}</span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useTheme } from "../../composables/useTheme";
import { useAuthStore } from "../../stores/auth";
import { disable2FA, enable2FA, getMe, init2FA, requestNickname } from "../../api/auth";
import QRCode from "qrcode";

const totpEnabled = ref(false);
const nickname = ref("");
const newNickname = ref("");
const secret = ref("");
const url = ref("");
const qrCode = ref("");
const code = ref("");
const error = ref("");
const success = ref("");
const loading = ref(false);
const router = useRouter();
const auth = useAuthStore();
const { locale, t } = useI18n();
const { theme, toggleTheme } = useTheme();
const isAuthed = computed(() => Boolean(auth.token));

const languages = {
  "zh-CN": "中文",
  "en-US": "EN",
};

function changeLanguage(code: string) {
  locale.value = code;
  localStorage.setItem("language", code);
}

function handleLogout() {
  // 退出登录后回到移动端首页
  auth.clearAuth();
  router.push({ name: "m-home" });
}

async function fetchMe() {
  // 获取用户信息与 2FA 状态
  if (!auth.token) return;
  const me = await getMe();
  totpEnabled.value = me.totpEnabled;
  nickname.value = me.nickname || me.username;
}

async function handleInit() {
  // 初始化 2FA，生成二维码
  error.value = "";
  success.value = "";
  loading.value = true;
  try {
    const data = await init2FA();
    secret.value = data.secret;
    url.value = data.url;
    qrCode.value = await QRCode.toDataURL(data.url, {
      width: 160,
      margin: 1,
      color: {
        dark: "#111111",
        light: "#ffffff",
      },
    });
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    loading.value = false;
  }
}

async function handleNickname() {
  // 提交昵称审核申请
  error.value = "";
  success.value = "";
  if (!newNickname.value.trim()) {
    error.value = t("settings.nicknameRequired");
    return;
  }
  loading.value = true;
  try {
    await requestNickname({ nickname: newNickname.value.trim() });
    success.value = t("settings.nicknameSubmitted");
    newNickname.value = "";
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    loading.value = false;
  }
}

async function handleEnable() {
  // 启用 2FA
  error.value = "";
  success.value = "";
  loading.value = true;
  try {
    await enable2FA({ code: code.value });
    totpEnabled.value = true;
    success.value = t("settings.enable2faSuccess");
    code.value = "";
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    loading.value = false;
  }
}

async function handleDisable() {
  // 关闭 2FA
  error.value = "";
  success.value = "";
  loading.value = true;
  try {
    await disable2FA({ code: code.value });
    totpEnabled.value = false;
    success.value = t("settings.disable2faSuccess");
    code.value = "";
    secret.value = "";
    url.value = "";
    qrCode.value = "";
  } catch (err) {
    error.value = (err as Error).message;
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchMe();
});
</script>

<style scoped lang="scss">
.mobile-settings {
  display: grid;
  gap: var(--mobile-space-3);
  padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
}

.mobile-section {
  display: grid;
  gap: var(--mobile-space-2);
}

.section-title {
  margin: 0;
  font-size: var(--mobile-font-sm);
  color: var(--muted);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.section-card {
  background: var(--surface);
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
  overflow: hidden;
  display: grid;
}

.mobile-item {
  width: 100%;
  text-align: left;
  background: transparent;
  color: var(--ink);
  border: none;
  padding: var(--mobile-space-3);
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--mobile-space-2);
}

.mobile-item::after {
  content: ">";
  color: var(--muted);
  font-weight: 700;
}

.mobile-item + .mobile-item {
  border-top: 1px solid var(--border);
}

.mobile-item.danger {
  color: #b00020;
}

.mobile-item.danger::after,
.mobile-item.no-chevron::after {
  content: "";
}

.mobile-item .value {
  color: var(--muted);
  font-weight: 500;
}

.section-inline {
  background: var(--surface);
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
  padding: var(--mobile-space-3);
  display: grid;
  gap: var(--mobile-space-2);
}

.inline-label {
  font-weight: 600;
  color: var(--ink);
}

.lang-switch {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--mobile-space-2);
}

.lang-switch button {
  background: var(--surface-alt);
  border: none;
  border-radius: var(--mobile-radius);
  padding: var(--mobile-space-2);
  font-weight: 600;
  cursor: pointer;
  color: var(--ink);
}

.lang-switch button.active {
  background: var(--accent);
  color: var(--accent-contrast);
}

.account-section {
  display: grid;
  gap: var(--mobile-space-2);
}

.card {
  background: var(--surface);
  border-radius: var(--mobile-radius);
  padding: var(--mobile-space-3);
  border: 1px solid var(--border);
  display: grid;
  gap: var(--mobile-space-2);
  font-size: var(--mobile-font-sm);
}

.card h3 {
  margin: 0;
  font-size: var(--mobile-font-sm);
}

.card p {
  margin: 0;
  font-size: var(--mobile-font-xs);
  color: var(--muted);
}

.card button {
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: var(--mobile-radius);
  padding: var(--mobile-space-2);
  cursor: pointer;
  font-weight: 600;
  width: 100%;
}

label {
  display: grid;
  gap: var(--mobile-space-1);
  font-weight: 600;
  font-size: var(--mobile-font-xs);
}

input {
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
  padding: var(--mobile-space-2);
  font-size: var(--mobile-font-sm);
  background: var(--surface);
  color: var(--ink);
}

.secret {
  display: grid;
  gap: var(--mobile-space-2);
}

.nickname {
  display: grid;
  gap: var(--mobile-space-2);
}

.nickname label,
.nickname span,
.nickname strong {
  line-height: 1.2;
}

code {
  display: block;
  background: var(--md-inline-code-bg);
  padding: var(--mobile-space-1) var(--mobile-space-2);
  border-radius: var(--mobile-radius);
  word-break: break-all;
}

.qr {
  width: calc(var(--mobile-unit) * 36);
  height: calc(var(--mobile-unit) * 36);
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
}

.success {
  color: #1b7b3d;
}

.error {
  color: #b00020;
}

.miyao {
  width: 100%;
}
</style>
