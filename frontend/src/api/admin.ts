import { request } from "./http";

export type AdminUser = {
  // 管理端用户结构
  id: number;
  username: string;
  nickname: string;
  role: string;
  status: string;
  totpEnabled: boolean;
  createdAt: string;
};

// listUsers 获取用户列表（按状态过滤）
export function listUsers(status: string) {
  const query = new URLSearchParams({ status }).toString();
  return request<{ items: AdminUser[] }>(`/api/v1/admin/users?${query}`, { auth: true });
}

// approveUser 通过用户审核
export function approveUser(id: number) {
  return request<{ id: number; status: string }>(`/api/v1/admin/users/${id}/approve`, {
    method: "POST",
    auth: true,
  });
}

// rejectUser 驳回用户审核
export function rejectUser(id: number, reason: string) {
  return request<{ id: number; status: string }>(`/api/v1/admin/users/${id}/reject`, {
    method: "POST",
    body: { reason },
    auth: true,
  });
}

export type NicknameRequest = {
  // 昵称审核申请结构
  id: number;
  userId: number;
  nickname: string;
  status: string;
  createdAt: string;
};

// listNicknameRequests 获取昵称申请列表
export function listNicknameRequests(status: string) {
  const query = new URLSearchParams(status ? { status } : {}).toString();
  const path = query ? `/api/v1/admin/nicknames?${query}` : "/api/v1/admin/nicknames";
  return request<{ items: NicknameRequest[] }>(path, { auth: true });
}

// approveNickname 通过昵称申请
export function approveNickname(id: number) {
  return request<{ id: number; status: string }>(`/api/v1/admin/nicknames/${id}/approve`, {
    method: "POST",
    auth: true,
  });
}

// rejectNickname 驳回昵称申请
export function rejectNickname(id: number, reason: string) {
  return request<{ id: number; status: string }>(`/api/v1/admin/nicknames/${id}/reject`, {
    method: "POST",
    body: { reason },
    auth: true,
  });
}
