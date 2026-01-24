<template>
  <section class="settings">
    <header v-if="isAuthed">
      <h2>{{ $t('settings.title') }}</h2>
      <p>{{ $t('settings.subtitle') }}</p>
    </header>
    <div class="pc-settings">
      <div class="card link-card">
        <button type="button" class="link-button" @click="router.push({ name: 'my-posts' })">
          {{ $t('myPosts.title') }}
        </button>
      </div>
      <div class="card link-card">
        <button type="button" class="link-button" @click="router.push({ name: 'messages' })">
          {{ $t('messages.title') }}
        </button>
      </div>
      <div v-if="isAuthed" class="card">
        <h3>{{ $t('settings.nicknameTitle') }}</h3>
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
      <div v-if="isAuthed" class="card">
        <button type="button" @click="handleLogout">{{ $t('nav.logout') }}</button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "../stores/auth";
import { disable2FA, enable2FA, getMe, init2FA, requestNickname } from "../api/auth";
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
const { t } = useI18n();
// 是否已登录（用于展示安全设置与退出入口）
const isAuthed = computed(() => Boolean(auth.token));

function handleLogout() {
  // 退出登录并返回首页
  auth.clearAuth();
  router.push({ name: "home" });
}

async function fetchMe() {
  // 获取用户资料与 2FA 状态
  if (!auth.token) return;
  const me = await getMe();
  totpEnabled.value = me.totpEnabled;
  nickname.value = me.nickname || me.username;
}

async function handleInit() {
  // 初始化 2FA，生成密钥与二维码
  error.value = "";
  success.value = "";
  loading.value = true;
  try {
    const data = await init2FA();
    secret.value = data.secret;
    url.value = data.url;
    qrCode.value = await QRCode.toDataURL(data.url, {
      width: 180,
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
.settings {
  padding: 32px 36px 60px;
  display: grid;
  gap: 24px;
}

header h2 {
  margin: 0;
  font-size: 28px;
}

header p {
  margin: 8px 0 0;
  color: var(--muted);
}

.card {
  background: var(--surface);
  border-radius: 20px;
  padding: 16px;
  border: 1px solid var(--border);
  display: grid;
  gap: 12px;
  font-size: 14px;
}

.card h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.card p {
  margin: 0;
  line-height: 1.5;
  font-size: 13px;
  color: var(--muted);
}

.pc-settings {
  display: grid;
  gap: 12px;
}

.link-card {
  padding: 0;
}

.link-button {
  border: none;
  background: transparent;
  color: var(--ink);
  border-radius: 12px;
  padding: 14px 16px;
  cursor: pointer;
  font-weight: 600;
  text-align: left;
  width: 100%;
}




button {
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: 10px;
  padding: 10px 16px;
  cursor: pointer;
  font-weight: 600;
}

.link-card .link-button {
  background: transparent;
  color: var(--ink);
  border: none;
}

label {
  display: grid;
  gap: 8px;
  font-weight: 600;
  font-size: 14px;
}

input {
  border-radius: 8px;
  border: 1px solid var(--border);
  padding: 10px 12px;
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
}

.secret {
  display: grid;
  gap: 10px;
}

code {
  display: block;
  background: var(--md-inline-code-bg);
  padding: 6px 8px;
  border-radius: 8px;
  word-break: break-all;
}

.qr {
  margin-top: 20px;
  width: 200px;
  height: 200px;
  border-radius: 5px;
  border: 1px solid var(--border);
}

.success {
  color: #1b7b3d;
}

.error {
  color: #b00020;
}

strong {
  display: inline-block;
  font-size: 16px;
  margin-top: 6px;
}

.miyao {
  width: 100%;
}

.nickname {
  display: grid;
  gap: 8px;
}

.nickname label {
  line-height: 1.2;
}

.nickname span,
.nickname strong {
  line-height: 1.2;
}

@media (max-width: 768px) {
  .settings {
    padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
  }

  .card {
    border-radius: var(--mobile-radius);
  }

  .settings header {
    display: none;
  }
}

</style>
