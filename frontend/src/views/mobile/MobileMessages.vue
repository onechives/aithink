<template>
  <section class="mobile-messages">
    <div class="filters">
      <button class="secondary" @click="setFilter('')">{{ $t('messages.all') }}</button>
      <button class="secondary" @click="setFilter('unread')">{{ $t('messages.unread') }}</button>
      <button class="secondary" @click="setFilter('read')">{{ $t('messages.read') }}</button>
    </div>

    <div class="list">
      <p v-if="messages.length === 0" class="empty">{{ $t('messages.noMessages') }}</p>
      <article v-for="item in messages" :key="item.id" class="item" :data-status="item.status">
        <div class="meta">
          <span class="title">{{ resolveTitle(item) }}</span>
          <span class="time">{{ formatDate(item.createdAt) }}</span>
        </div>
        <p>{{ resolveContent(item) }}</p>
        <button v-if="item.status === 'unread'" @click="handleRead(item.id)">{{ $t('messages.markAsRead') }}</button>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { listMessages, markMessageRead, type MessageItem } from "../../api/messages";

const messages = ref<MessageItem[]>([]);
const filter = ref("");
const { locale, t } = useI18n();

async function fetchMessages() {
  // 拉取消息列表并触发未读刷新
  const data = await listMessages(filter.value);
  messages.value = data.items;
  window.dispatchEvent(new Event("messages-updated"));
}

function setFilter(value: string) {
  // 切换筛选条件
  filter.value = value;
  fetchMessages();
}

async function handleRead(id: number) {
  // 标记已读
  await markMessageRead(id);
  fetchMessages();
}

function formatDate(value: string) {
  // 根据语言格式化时间
  return new Date(value).toLocaleString(locale.value);
}

function resolveTitle(item: MessageItem) {
  // 优先使用 i18nKey
  if (item.i18nKey) {
    return t(`${item.i18nKey}.title`, item.params || {});
  }
  return item.title;
}

function resolveContent(item: MessageItem) {
  // 根据是否有原因决定内容模板
  if (item.i18nKey) {
    const params = item.params || {};
    if (params.reason) {
      return t(`${item.i18nKey}.contentWithReason`, params);
    }
    return t(`${item.i18nKey}.content`, params);
  }
  return item.content;
}

onMounted(() => {
  fetchMessages();
});
</script>

<style scoped lang="scss">
.mobile-messages {
  padding: var(--mobile-space-4) var(--mobile-space-3) var(--mobile-space-6);
  display: grid;
  gap: var(--mobile-space-3);
}

.filters {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--mobile-space-2);
}

button {
  border: none;
  background: var(--accent);
  color: var(--accent-contrast);
  border-radius: var(--mobile-radius);
  padding: 8px 10px;
  cursor: pointer;
  font-weight: 600;
}

button.secondary {
  background: var(--surface-alt);
  color: var(--ink);
}

.list {
  display: grid;
  gap: var(--mobile-space-2);
}

.empty {
  margin: 0;
  color: var(--muted);
}

.item {
  background: var(--surface);
  border-radius: var(--mobile-radius);
  border: 1px solid var(--border);
  padding: var(--mobile-space-3);
  display: grid;
  gap: var(--mobile-space-2);
}

.item[data-status="unread"] {
  border-color: var(--accent);
}

.meta {
  display: grid;
  gap: 4px;
}

.title {
  font-weight: 600;
}

.time {
  color: var(--muted);
  font-size: 12px;
}

p {
  margin: 0;
  color: var(--muted);
}
</style>
