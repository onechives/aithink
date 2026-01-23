<template>
  <section class="login">
    <div class="card">
      <h2>{{ $t('login.title') }}</h2>
      <p v-if="step === 'password'">{{ $t('login.welcome') }}</p>
      <p v-else>{{ $t('login.step2fa') }}</p>
      <form v-if="step === 'password'" @submit.prevent="handleLogin">
        <label>
          {{ $t('login.username') }}
          <input v-model="username" type="text" required />
        </label>
        <label>
          {{ $t('login.password') }}
          <input v-model="password" type="password" required />
        </label>
        <button type="submit">{{ $t('login.submit') }}</button>
      </form>

      <form v-else @submit.prevent="handleVerify">
        <label>
          {{ $t('login.code') }}
          <input v-model="code" type="text" maxlength="6" required />
        </label>
        <div class="actions">
          <button class="secondary" type="button" @click="reset">{{ $t('login.back') }}</button>
          <button type="submit">{{ $t('login.verify') }}</button>
        </div>
      </form>

      <div class="hint">
        <RouterLink to="/m/register">{{ $t('login.noAccount') }}{{ $t('login.register') }}</RouterLink>
      </div>
      <span class="error" v-if="error">{{ error }}</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { login, loginVerify } from "../../api/auth";
import { useAuthStore } from "../../stores/auth";

const username = ref("");
const password = ref("");
const code = ref("");
const tempToken = ref("");
const error = ref("");
const step = ref<"password" | "2fa">("password");
const router = useRouter();
const auth = useAuthStore();

async function handleLogin() {
  error.value = "";
  try {
    const data = await login({ username: username.value, password: password.value });
    if (data.need2fa) {
      tempToken.value = data.tempToken || "";
      step.value = "2fa";
      return;
    }
    if (data.token) {
      auth.setAuth(data.token, String(data.userId), data.role, username.value);
      router.push({ name: "m-home" });
    }
  } catch (err) {
    error.value = (err as Error).message;
  }
}

async function handleVerify() {
  error.value = "";
  try {
    const data = await loginVerify({ tempToken: tempToken.value, code: code.value });
    if (data.token) {
      auth.setAuth(data.token, String(data.userId), data.role, username.value);
      router.push({ name: "m-home" });
    }
  } catch (err) {
    error.value = (err as Error).message;
  }
}

function reset() {
  step.value = "password";
  code.value = "";
  tempToken.value = "";
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

.actions {
  display: flex;
  gap: 10px;
}

.actions .secondary {
  background: var(--surface-alt);
  color: var(--ink);
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
</style>
