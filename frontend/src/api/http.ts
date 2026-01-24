import { useAuthStore } from "../stores/auth";
import router from "../router";

// API 根地址（未配置时使用同域）
const baseUrl = import.meta.env.VITE_API_BASE || "";

type RequestOptions = {
  method?: string;
  body?: unknown;
  auth?: boolean;
};

export async function request<T>(path: string, options: RequestOptions = {}) {
  // 统一 JSON 请求头
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (options.auth) {
    // 需要鉴权时注入 token 与 user_id
    const auth = useAuthStore();
    if (auth.token) {
      headers.Authorization = `Bearer ${auth.token}`;
      headers["user_id"] = auth.userId;
    }
  }

  // 统一的请求入口
  const response = await fetch(`${baseUrl}${path}`, {
    method: options.method || "GET",
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const data = await response.json();
  // 后端约定：code === 1000 为成功
  if (!response.ok || data.code !== 1000) {
    if (options.auth && (data.code === 1002 || data.code === 1003)) {
      const auth = useAuthStore();
      auth.clearAuth();
      const currentPath = router.currentRoute.value.path;
      const target = currentPath.startsWith("/m") ? "/m/login" : "/login";
      if (currentPath !== target) {
        router.replace({ path: target, query: { redirect: currentPath } });
      }
    }
    throw new Error(data.message || "request failed");
  }
  return data.data as T;
}
