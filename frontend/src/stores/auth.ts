import { defineStore } from "pinia";

type AuthState = {
  // 登录态缓存
  token: string;
  userId: string;
  role: string;
  username: string;
};

// storageKey 本地缓存 key
const storageKey = "blog-admin-auth";

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => {
    const saved = localStorage.getItem(storageKey);
    if (saved) {
      try {
        const parsed = JSON.parse(saved) as Partial<AuthState>;
        return {
          token: parsed.token || "",
          userId: parsed.userId || "",
          role: parsed.role || "",
          username: parsed.username || "",
        };
      } catch {
        return { token: "", userId: "", role: "", username: "" };
      }
    }
    return { token: "", userId: "", role: "", username: "" };
  },
  actions: {
    setAuth(token: string, userId: string, role: string, username = "") {
      // 写入内存与本地缓存
      this.token = token;
      this.userId = userId;
      this.role = role;
      this.username = username;
      localStorage.setItem(storageKey, JSON.stringify({ token, userId, role, username }));
    },
    clearAuth() {
      // 清理登录态
      this.token = "";
      this.userId = "";
      this.role = "";
      this.username = "";
      localStorage.removeItem(storageKey);
    },
  },
});
