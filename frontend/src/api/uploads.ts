import { useAuthStore } from "../stores/auth";

type UploadResponse = {
  // 上传接口统一响应结构
  code: number;
  message: string;
  data?: {
    url: string;
  };
};

// baseUrl 允许通过环境变量配置 API 域名
const baseUrl = import.meta.env.VITE_API_BASE || "";

// uploadImage 上传图片，返回可访问的 URL 列表
export async function uploadImage(files: File[]) {
  const auth = useAuthStore();
  const urls: string[] = [];

  for (const file of files) {
    const form = new FormData();
    form.append("file", file);

    const response = await fetch(`${baseUrl}/api/v1/upload`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${auth.token}`,
        user_id: auth.userId,
      },
      body: form,
    });

    // 兼容后端返回的 JSON 文本
    const raw = await response.text();
    let data: UploadResponse;
    try {
      data = JSON.parse(raw) as UploadResponse;
    } catch (err) {
      throw new Error(`upload failed: ${raw}`);
    }
    if (!response.ok || data.code !== 1000 || !data.data?.url) {
      throw new Error(data.message || "upload failed");
    }
    urls.push(data.data.url);
  }

  return urls;
}
