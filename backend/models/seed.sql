INSERT INTO users (id, username, nickname, password_hash, role, status, totp_secret, totp_enabled, created_at, updated_at)
VALUES
(1, 'admin', '管理员', '$2a$10$unCvAZPbJE34rvHT/Mbv8OyCom3CvPvgE3U.vQ7dQGYCE4wJHfd8a', 'admin', 'approved', '', 0, NOW(), NOW());

INSERT INTO posts (id, title, summary, content_md, cover_url, category, tags, author_id, status, like_count, created_at, updated_at)
VALUES
(10001, '第一篇：搭建一个简约博客', '从目录结构到接口，记录搭建过程。', '# 搭建一个简约博客\n\n这里是正文内容，支持 **Markdown**。', 'https://images.unsplash.com/photo-1451187580459-43490279c0fa', '工程', 'Gin,Vue3', 1, 'approved', 12, NOW(), NOW()),
(10002, '第二篇：写作工作流', '如何用 Markdown 高效写作与发布。', '# 写作工作流\n\n- 写作\n- 预览\n- 发布', 'https://images.unsplash.com/photo-1455390582262-044cdead277a', '写作', 'Markdown,效率', 1, 'approved', 23, NOW(), NOW()),
(10003, '第三篇：样式与布局', '侧边栏+卡片区域的布局细节。', '# 样式与布局\n\n强调简约大气。', 'https://images.unsplash.com/photo-1498050108023-c5249f4df085', '设计', 'CSS,布局', 1, 'approved', 7, NOW(), NOW());
