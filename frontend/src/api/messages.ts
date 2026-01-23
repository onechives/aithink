import { request } from "./http";

export type MessageItem = {
  // 站内信结构（含 i18n 扩展字段）
  id: number;
  title: string;
  content: string;
  status: string;
  createdAt: string;
  messageType?: string;
  i18nKey?: string;
  params?: Record<string, string>;
};

// listMessages 获取站内信列表（可按状态过滤）
export function listMessages(status = "") {
  const query = new URLSearchParams(status ? { status } : {}).toString();
  const path = query ? `/api/v1/me/messages?${query}` : "/api/v1/me/messages";
  return request<{ items: MessageItem[] }>(path, { auth: true });
}

// markMessageRead 标记消息已读
export function markMessageRead(id: number) {
  return request<{ id: number; status: string }>(`/api/v1/me/messages/${id}/read`, {
    method: "POST",
    auth: true,
  });
}

// getUnreadCount 获取未读数量（用于红点）
export function getUnreadCount() {
  return request<{ count: number }>("/api/v1/me/messages/unread-count", { auth: true });
}
