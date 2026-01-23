import { request } from "./http";

export type PostId = string | number;

export type PostSummary = {
  // 后端返回的文章摘要结构
  id: string;
  title: string;
  summary: string;
  coverUrl: string;
  category: string;
  tags: string;
  likeCount: number;
  author: string;
  authorId: number;
  createdAt: string;
  updatedAt: string;
  status: string;
};

export type PostDetail = PostSummary & {
  // 文章正文（Markdown）
  content: string;
};

export type PostTitle = {
  // 仅包含 ID 与标题的轻量结构
  id: string;
  title: string;
};

export type PostListResponse = {
  // 列表分页响应
  items: PostSummary[];
  total: number;
};

export type PostPayload = {
  // 写作/编辑表单字段
  title: string;
  summary: string;
  content: string;
  coverUrl: string;
  category: string;
  tags: string;
};

// normalizeId 保证 ID 以字符串形式传输，避免大整数精度丢失
const normalizeId = (id: PostId) => String(id);

export function listPosts(params: {
  sort: string;
  keyword: string;
  page: number;
  size: number;
}) {
  // 文章列表（分页 + 搜索 + 排序）
  const query = new URLSearchParams({
    sort: params.sort,
    keyword: params.keyword,
    page: String(params.page),
    size: String(params.size),
  }).toString();
  return request<PostListResponse>(`/api/v1/posts?${query}`);
}

export function listPostTitles(params: { sort: string; keyword: string; size: number }) {
  // 文章标题列表（侧边栏用）
  const query = new URLSearchParams({
    sort: params.sort,
    keyword: params.keyword,
    size: String(params.size),
  }).toString();
  return request<PostTitle[]>(`/api/v1/post-titles?${query}`);
}

export function getPostDetail(id: PostId) {
  // 先尝试携带鉴权（作者可看未审核），失败则降级为匿名访问
  const normalizedId = normalizeId(id);
  return request<PostDetail>(`/api/v1/posts/${normalizedId}`, { auth: true }).catch(() => {
    return request<PostDetail>(`/api/v1/posts/${normalizedId}`);
  });
}

export function likePost(id: PostId) {
  // 点赞文章
  return request<{ likeCount: number }>(`/api/v1/posts/${normalizeId(id)}/like`, { method: "POST" });
}

export function createPost(payload: PostPayload) {
  // 新建文章（需登录）
  return request<{ id: string }>("/api/v1/posts", { method: "POST", body: payload, auth: true });
}

export function updatePost(id: PostId, payload: PostPayload) {
  // 更新文章（需登录）
  return request<{ id: string }>(`/api/v1/posts/${normalizeId(id)}`, {
    method: "PUT",
    body: payload,
    auth: true,
  });
}

export function deletePost(id: PostId) {
  // 删除文章（需登录）
  return request<{ id: string }>(`/api/v1/posts/${normalizeId(id)}`, {
    method: "DELETE",
    auth: true,
  });
}

export function listMyPosts(params: { page: number; size: number }) {
  // 获取我的文章列表
  const query = new URLSearchParams({
    page: String(params.page),
    size: String(params.size),
  }).toString();
  return request<PostListResponse>(`/api/v1/me/posts?${query}`, { auth: true });
}

export function listAdminPosts(params: { status: string; page: number; size: number }) {
  // 管理端获取文章列表
  const query = new URLSearchParams({
    status: params.status,
    page: String(params.page),
    size: String(params.size),
  }).toString();
  return request<PostListResponse>(`/api/v1/admin/posts?${query}`, { auth: true });
}

export function approvePost(id: PostId) {
  // 管理端通过审核
  return request<{ id: string; status: string }>(`/api/v1/admin/posts/${normalizeId(id)}/approve`, {
    method: "POST",
    auth: true,
  });
}

export function rejectPost(id: PostId, reason: string) {
  // 管理端驳回审核
  return request<{ id: string; status: string }>(`/api/v1/admin/posts/${normalizeId(id)}/reject`, {
    method: "POST",
    body: { reason },
    auth: true,
  });
}
