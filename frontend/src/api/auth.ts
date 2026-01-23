import { request } from "./http";

export type LoginResponse = {
  // 登录响应：可能需要 2FA
  userId: number;
  role: string;
  token?: string;
  need2fa: boolean;
  tempToken?: string;
};

export type MeResponse = {
  // 当前用户信息
  userId: number;
  username: string;
  nickname: string;
  role: string;
  status: string;
  totpEnabled: boolean;
};

// register 注册账号（注册后需审核）
export function register(payload: { username: string; password: string }) {
  return request<{ status: string }>("/api/v1/register", {
    method: "POST",
    body: payload,
  });
}

// login 登录第一步（账号密码）
export function login(payload: { username: string; password: string }) {
  return request<LoginResponse>("/api/v1/login", {
    method: "POST",
    body: payload,
  });
}

// loginVerify 登录第二步（2FA 验证）
export function loginVerify(payload: { tempToken: string; code: string }) {
  return request<LoginResponse>("/api/v1/login/verify", {
    method: "POST",
    body: payload,
  });
}

// getMe 获取当前用户信息
export function getMe() {
  return request<MeResponse>("/api/v1/me", { auth: true });
}

// init2FA 初始化 2FA，返回密钥与 otpauth URL
export function init2FA() {
  return request<{ secret: string; url: string }>("/api/v1/me/2fa/init", {
    method: "POST",
    auth: true,
  });
}

// enable2FA 启用 2FA
export function enable2FA(payload: { code: string }) {
  return request<{ enabled: boolean }>("/api/v1/me/2fa/enable", {
    method: "POST",
    body: payload,
    auth: true,
  });
}

// disable2FA 关闭 2FA
export function disable2FA(payload: { code: string }) {
  return request<{ enabled: boolean }>("/api/v1/me/2fa/disable", {
    method: "POST",
    body: payload,
    auth: true,
  });
}

// requestNickname 申请昵称变更
export function requestNickname(payload: { nickname: string }) {
  return request<{ status: string }>("/api/v1/me/nickname", {
    method: "POST",
    body: payload,
    auth: true,
  });
}
