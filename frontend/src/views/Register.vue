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
        <button type="submit">{{ $t('register.submit') }}</button>
      </form>
      <div class="hint">
        <RouterLink to="/login">{{ $t('register.loginLink') }}</RouterLink>
      </div>
      <span class="success" v-if="success">{{ success }}</span>
      <span class="error" v-if="error">{{ error }}</span>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { register } from "../api/auth";

const username = ref("");
const password = ref("");
const error = ref("");
const success = ref("");
const { t } = useI18n();

async function handleRegister() {
  // 提交注册，成功后提示“待审核”
  error.value = "";
  success.value = "";
  try {
    await register({ username: username.value, password: password.value });
    success.value = t("register.pendingApprovalSuccess");
    username.value = "";
    password.value = "";
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
  padding: 40px;
}

.card {
  max-width: 420px;
  width: 100%;
  background: var(--surface);
  border-radius: 24px;
  padding: 32px;
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
  border-radius: 10px;
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

.success {
  margin-top: 12px;
  color: #1b7b3d;
  display: block;
}

.error {
  margin-top: 12px;
  color: #b00020;
  display: block;
}

@media (max-width: 768px) {
  .login {
    padding: var(--mobile-space-5) var(--mobile-space-3) var(--mobile-space-6);
  }

  .card {
    border-radius: var(--mobile-radius);
    padding: var(--mobile-space-5);
  }
}

</style>
