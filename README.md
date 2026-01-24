# Aithink 博客系统

## 项目简介
Aithink 是一个全栈博客平台，支持文章发布与审核、消息通知、多语言切换、移动端与 PC 端分离布局，并提供昵称审核与两步验证（TOTP/Google Authenticator）能力。

## 主要功能
- 文章浏览与详情展示（Markdown 渲染）
- 文章点赞、文章列表与搜索
- 写作编辑器（md-editor-v3）
- 管理后台：文章/用户/昵称审核
- 站内信与未读提醒
- 账户昵称申请与审核
- 两步验证（TOTP/Google Authenticator）
- 中英双语（vue-i18n）
- 移动端与 PC 端独立路由与布局

## 技术栈
- 前端：Vue 3、Vite、Pinia、Vue Router、vue-i18n、md-editor-v3、marked、Sass
- 后端：Go + Gin、MariaDB/MySQL、Redis、JWT、TOTP
- 部署：Docker + Nginx（反向代理 /api 与 /uploads）

## 目录结构
```
aithink/
  backend/           Go 后端
  frontend/          Vue 前端
  deploy/            Docker Compose + Nginx 配置
```

## 本地开发

### 1) 准备依赖
- Node.js 20+
- Go 1.25.5
- MariaDB 15.1
- Redis 7

### 2) 初始化数据库
创建数据库并导入表结构与种子数据：
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS blog DEFAULT CHARSET utf8mb4;"
mysql -u root -p blog < backend/models/schema.sql
mysql -u root -p blog < backend/models/seed.sql
```
注意：`seed.sql` 会插入示例文章与管理员账号（密码为加密值，如需登录请自行重置密码）。

### 3) 启动后端
```bash
cd backend
go run . ./conf/config.yaml
```
默认端口：`8080`  
配置文件：`backend/conf/config.yaml`

### 4) 启动前端
```bash
cd frontend
npm install
npm run dev
```
默认端口：`5173`  
Vite 已配置 `/api` 与 `/uploads` 代理至 `http://localhost:8080`。

## 生产部署（Docker + 宿主机 MariaDB/Redis + 宿主机 Nginx）

### 1) 安装并配置 MariaDB（Debian）
```bash
sudo apt update
sudo apt install -y mariadb-server
sudo systemctl enable --now mariadb
```

创建数据库并导入结构与种子数据：
```bash
sudo mysql -e "CREATE DATABASE IF NOT EXISTS blog DEFAULT CHARSET utf8mb4;"
sudo mysql blog < backend/models/schema.sql
sudo mysql blog < backend/models/seed.sql
```

允许容器访问（示例：绑定到 0.0.0.0 或 docker0）：
```text
/etc/mysql/mariadb.conf.d/50-server.cnf
bind-address = 0.0.0.0
```
修改后重启：
```bash
sudo systemctl restart mariadb
```
建议使用防火墙限制 3306 端口仅对本机或 Docker 网桥开放。

### 2) 安装并配置 Redis（Debian）
```bash
sudo apt install -y redis-server
```

编辑配置并设置密码（示例）：
```text
/etc/redis/redis.conf
bind 0.0.0.0
requirepass grdxjpqs
```
修改后重启：
```bash
sudo systemctl restart redis-server
```
建议使用防火墙限制 6379 端口仅对本机或 Docker 网桥开放。

### 3) 配置后端 Docker 连接宿主机服务
`backend/conf/config.docker.yaml` 已默认使用 `host.docker.internal` 访问宿主机 MariaDB/Redis。
如需自定义账号、密码或端口，请修改该文件中的 `mysql` 与 `redis` 配置。

> 注意：Linux 上 `host.docker.internal` 需要 Docker 20.10+ 并通过 `extra_hosts` 映射（已在 `deploy/docker-compose.yml` 中配置）。  
> 若宿主机未支持，请改为宿主机在 docker0 网桥上的 IP（例如 `172.17.0.1`）。

### 4) 准备宿主机日志目录与上传目录
```bash
mkdir -p deploy/data/logs deploy/data/uploads
```
`deploy/data/logs` 将映射到容器内 `/app/logfile`，日志文件路径为 `/app/logfile/log.log`。

### 5) 构建并启动容器
```bash
cd deploy
docker compose up -d --build
```

容器端口（仅本机回环）：
- 前端容器：`127.0.0.1:8081`
- 后端容器：`127.0.0.1:8080`

### 6) 配置宿主机 Nginx 反向代理
将 `deploy/nginx.conf` 复制到 Nginx 站点配置（示例路径 `/etc/nginx/sites-available/aithink.conf`）：
```bash
sudo cp deploy/nginx.conf /etc/nginx/sites-available/aithink.conf
sudo ln -s /etc/nginx/sites-available/aithink.conf /etc/nginx/sites-enabled/aithink.conf
sudo nginx -t
sudo systemctl reload nginx
```

如需 HTTPS，可按需安装 Certbot 并申请证书：
```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## 配置说明

### 后端配置
默认读取 `backend/conf/config.yaml`，也可在启动时传参指定：
```bash
go run . /path/to/config.yaml
```

Docker 环境使用 `backend/conf/config.docker.yaml`，并由 `backend/Dockerfile` 启动：
```bash
/app/server /app/conf/config.docker.yaml
```

### 前端 API 地址
前端默认使用同域 `/api`。如需指定其他后端地址：
```
VITE_API_BASE=http://your-api-host
```

## 其他说明
- 上传文件默认保存在 `backend/uploads/`，Docker 部署时映射到 `deploy/data/uploads/`
- 日志输出在 `deploy/data/logs/`（容器内路径 `/app/logfile`）


- 详细部署教程在博客 https://www.ai-code.net/
- （博客就是用这个项目搭建的 ）