# Aithink 部署指南 (Docker Host 模式)

本文档介绍如何在生产服务器上使用 **Docker Host 网络模式** 部署 Aithink 系统。此模式在提升性能的同时，能通过本地回环地址（127.0.0.1）实现最高级别的安全性。

---

## 1. 基础设施要求
- **操作系统**: 建议使用 Debian 12 或 Ubuntu 22.04+。
- **软件依赖**:
  - Docker & Docker Compose v2.0+
  - MariaDB (宿主机安装)
  - Redis (宿主机安装)
  - Nginx (宿主机安装，作为反向代理)

---

## 2. 安全性加固 (宿主机)

为了防止数据库和 Redis 暴露在公网，请修改它们的配置文件，使其仅监听本地回环地址。

### 2.1 MariaDB 配置
编辑 `/etc/mysql/mariadb.conf.d/50-server.cnf`:
```ini
bind-address = 127.0.0.1
```
重启服务: `sudo systemctl restart mariadb`

### 2.2 Redis 配置
编辑 `/etc/redis/redis.conf`:
```ini
bind 127.0.0.1
requirepass your_secure_password
```
重启服务: `sudo systemctl restart redis-server`

---

## 3. 数据库初始化

1. 创建数据库并导入结构：
   ```bash
   sudo mysql -e "CREATE DATABASE IF NOT EXISTS blog DEFAULT CHARSET utf8mb4;"
   sudo mysql blog < backend/models/schema.sql
   sudo mysql blog < backend/models/seed.sql
   ```
2. **授权 (关键)**：允许本地 TCP 连接。
   ```sql
   -- 在 mysql 终端执行
   GRANT ALL PRIVILEGES ON blog.* TO 'bloguser'@'127.0.0.1' IDENTIFIED BY 'your_password';
   FLUSH PRIVILEGES;
   ```

---

## 4. 部署步骤

### 4.1 准备目录与配置
将项目同步到服务器 `/opt/aithink` 后执行：

```bash
cd /opt/aithink

# 1. 创建数据挂载目录
mkdir -p deploy/data/logs deploy/data/uploads deploy/data/conf

# 2. 准备 Docker 专用配置文件 (手动拷贝以确保同步)
cp backend/conf/config.docker.yaml deploy/data/conf/
```

### 4.2 启动容器
```bash
cd deploy
# 首次构建并后台启动
docker compose up -d --build
```

---

## 5. Nginx 反向代理配置

将 `deploy/nginx.conf` 复制到宿主机 Nginx 环境。主要转发逻辑如下：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # 前端页面 (Docker Host 模式监听 8081)
    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        # ... 其它 proxy 参数
    }

    # 后端接口 (Docker Host 模式监听 8080)
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
    }
}
```

---

## 6. 常见问题排查

### 6.1 容器启动报错 `OCI runtime create failed`
- **原因**: Docker 尝试将一个不存在的文件挂载为文件，结果自动创建了一个同名目录。
- **解决**: 停止容器，删除宿主机上错误的目录：`rm -rf deploy/data/conf/config.docker.yaml`，确保该路径下是真实的文件后再启动。

### 6.2 数据库连接失败 `Access denied for user bloguser@localhost`
- **原因**: 即使在 Host 模式下，Go 驱动也会通过 TCP 连接 `127.0.0.1`，而 MySQL 默认可能只开了 Socket (`localhost`) 权限。
- **解决**: 按照第 3 节中的授权步骤，为 `'bloguser'@'127.0.0.1'` 授权。

### 6.3 配置文件未生效
- **原因**: 由于 `rsync` 排除规则，某些挂载文件可能未被更新。
- **解决**: 手动执行 `cp backend/conf/config.docker.yaml deploy/data/conf/` 并重启容器。

---

## 7. 维护命令
- **查看后端日志**: `docker compose logs -f backend`
- **更新代码后重启**: `docker compose down && docker compose up -d --build`
- **只重启某个服务**: `docker compose restart backend`

## 测试