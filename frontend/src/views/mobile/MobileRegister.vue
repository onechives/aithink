<template>
  <section class="login">
    <div class="card">
      <h2>{{ $t('register.title') }}</h2>
      <p>{{ $t('register.reviewNotice') }}</p>
      <form @submit.prevent="handleRegister">
        <label>
          {{ $t('register.username') }}
          <input v-model="username" type="text" required />
        </label>
        <label>
          {{ $t('register.password') }}
          <input v-model="password" type="password" required />
        </label>
        <label>
          {{ $t('register.confirmPassword') }}
          <input v-model="confirmPassword" type="password" required />
        </label>
        <button type="submit">{{ $t('register.submit') }}</button>
      </form>

      <div class="hint">
        <RouterLink to="/m/login">{{ $t('register.loginLink') }}</RouterLink>
      </div>
      <span class="error" v-if="error">{{ error }}</span>
      <span class="success" v-if="success">{{ success }}</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { register } from "../../api/auth";

const username = ref("");
const password = ref("");
const confirmPassword = ref("");
const error = ref("");
const success = ref("");
const router = useRouter();
const { t } = useI18n();

async function handleRegister() {
  error.value = "";
  success.value = "";
  if (password.value !== confirmPassword.value) {
    error.value = t("register.passwordMismatch");
    return;
  }
  try {
    await register({ username: username.value, password: password.value });
    success.value = t("register.pendingApprovalSuccess");
    setTimeout(() => {
      router.push({ name: "m-login" });
    }, 1000);
  } catch (err) {
    error.value = (err as Error).message;
  }
}
</script>

<style scoped lang="scss">
.login {
  min-height: calc(100vh - 120px);
  display: grid;
  place-items: center;
  padding: var(--mobile-space-5) var(--mobile-space-3) var(--mobile-space-6);
}

.card {
  max-width: 420px;
  width: 100%;
  background: var(--surface);
  border-radius: var(--mobile-radius);
  padding: var(--mobile-space-5);
  box-shadow: var(--shadow);
}

h2 {
  margin: 0 0 10px;
}

p {
  margin: 0 0 20px;
  color: var(--muted);
}

form {
  display: grid;
  gap: 14px;
}

label {
  display: grid;
  gap: 6px;
  font-weight: 600;
}

input {
  border-radius: 8px;
  border: 1px solid var(--border);
  padding: 10px 12px;
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
}

button {
  margin-top: 6px;
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: var(--mobile-radius);
  padding: 12px 16px;
  cursor: pointer;
  font-weight: 600;
}

.hint {
  margin-top: 16px;
}

.hint a {
  color: var(--ink);
  text-decoration: none;
  font-weight: 600;
}

.error {
  margin-top: 12px;
  color: #b00020;
  display: block;
}

.success {
  margin-top: 12px;
  color: #1b7b3d;
  display: block;
}
</style>
