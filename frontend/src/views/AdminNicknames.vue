<template>
  <section class="admin">
    <header class="toolbar">
      <div>
        <h2>{{ $t('admin.nicknameReview') }}</h2>
        <p>{{ $t('admin.nicknameReviewDesc') }}</p>
      </div>
      <div class="actions">
        <button class="secondary" @click="setStatus('pending')">{{ $t('admin.pending') }}</button>
        <button class="secondary" @click="setStatus('approved')">{{ $t('admin.approved') }}</button>
        <button class="secondary" @click="setStatus('rejected')">{{ $t('admin.rejected2') }}</button>
      </div>
    </header>

    <div class="table">
      <div class="row header">
        <span>{{ $t('admin.userid') }}</span>
        <span>{{ $t('admin.nickname') }}</span>
        <span>{{ $t('admin.status') }}</span>
        <span>{{ $t('admin.action') }}</span>
      </div>
      <div class="row" v-for="item in items" :key="item.id">
        <span :data-label="$t('admin.userid')">#{{ item.userId }}</span>
        <span class="title" :data-label="$t('admin.nickname')">{{ item.nickname }}</span>
        <span class="status" :data-status="item.status" :data-label="$t('admin.status')">{{ statusLabel(item.status) }}</span>
        <span class="row-actions" :data-label="$t('admin.action')">
          <button v-if="item.status !== 'approved'" @click="handleApprove(item.id)">{{ $t('admin.approve') }}</button>
          <button v-if="item.status !== 'rejected'" class="secondary" @click="handleReject(item.id)">{{
            $t('admin.reject') }}</button>
        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { approveNickname, listNicknameRequests, rejectNickname, type NicknameRequest } from "../api/admin";

const items = ref<NicknameRequest[]>([]);
const status = ref("pending");
const { t } = useI18n();

async function fetchItems() {
  // 按状态获取昵称申请
  const data = await listNicknameRequests(status.value);
  items.value = data.items;
}

function setStatus(value: string) {
  // 切换状态并刷新
  status.value = value;
  fetchItems();
}

function statusLabel(value: string) {
  // 审核状态文案
  if (value === "approved") return t("admin.approved");
  if (value === "rejected") return t("admin.rejected2");
  return t("admin.pending");
}

async function handleApprove(id: number) {
  // 通过昵称申请
  await approveNickname(id);
  fetchItems();
}

async function handleReject(id: number) {
  // 驳回昵称申请并填写原因
  const reason = window.prompt(t("admin.rejectReasonPrompt"), "") || "";
  await rejectNickname(id, reason.trim());
  fetchItems();
}

onMounted(() => {
  fetchItems();
});
</script>

<style scoped lang="scss">
.admin {
  padding: 32px 36px 60px;
  display: grid;
  gap: 24px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
}

.toolbar h2 {
  margin: 0;
  font-size: 28px;
}

.toolbar p {
  margin: 8px 0 0;
  color: var(--muted);
}

.actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

button {
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: 10px;
  padding: 8px 14px;
  cursor: pointer;
  font-weight: 600;
}

button.secondary {
  background: var(--surface-alt);
  color: var(--ink);
}

.table {
  background: var(--surface);
  border-radius: 20px;
  border: 1px solid var(--border);
  overflow: hidden;
}

.row {
  display: grid;
  grid-template-columns: 1fr 1.4fr 0.8fr 1fr;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border);
  align-items: center;
  gap: 12px;
}

.row.header {
  background: var(--surface-alt);
  font-weight: 600;
}

.row:last-child {
  border-bottom: none;
}

.title {
  font-weight: 600;
}

.status[data-status="approved"] {
  color: #1b7b3d;
}

.status[data-status="rejected"] {
  color: #b00020;
}

@media (max-width: 960px) {
  .toolbar {
    flex-direction: column;
    align-items: flex-start;
  }

  .row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .admin {
    padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
  }

  .actions {
    flex-wrap: wrap;
  }

  .actions button {
    flex: 1 1 auto;
  }

  .table {
    border: none;
    background: transparent;
  }

  .row {
    grid-template-columns: 1fr;
    border: 1px solid var(--border);
    border-radius: var(--mobile-radius);
    margin-bottom: var(--mobile-space-3);
  }

  .row.header {
    display: none;
  }

  .row > span {
    display: grid;
    grid-template-columns: minmax(calc(var(--mobile-unit) * 26), 1fr) 2fr;
    gap: var(--mobile-space-2);
    align-items: center;
  }

  .row > span::before {
    content: attr(data-label);
    color: var(--muted);
    font-size: var(--mobile-font-xs);
    text-transform: uppercase;
    letter-spacing: calc(var(--mobile-unit) * 0.1);
  }

  .row .row-actions {
    display: grid;
    gap: var(--mobile-space-2);
  }
}

</style>
