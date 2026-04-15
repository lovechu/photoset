# 🐋 PhotoSet Docker 部署指南

本目录包含 PhotoSet 项目的完整 Docker 编排配置，支持一键式部署开发和生产环境。

## 📁 文件结构

```
docker/
├── backend.Dockerfile         # Go 后端服务镜像
├── frontend.Dockerfile        # 用户端 Vue 前端镜像
├── admin.Dockerfile           # 管理后台 Vue 前端镜像
├── nginx-frontend.conf        # 用户端 Nginx 配置
├── nginx-admin.conf           # 管理后台 Nginx 配置
├── README.DOCKER.md           # 本文件
└── (docker-compose.yml 在项目根目录)
```

## 🚀 快速开始

### 1. 环境准备
确保已安装：
- Docker Engine 20.10+
- Docker Compose v2
- Git

### 2. 配置环境变量
```bash
# 复制 Docker 环境模板
cp .env.docker.example .env.docker

# 编辑重要安全配置
nano .env.docker
```

**必须修改的配置**：
- `JWT_SECRET` - 生产环境 JWT 密钥
- `SIGN_SECRET` - URL 签名密钥
- `CORS_ALLOW_ORIGINS` - 前端域名

### 3. 启动全部服务
```bash
# 构建并启动所有容器（第一次使用）
docker-compose up -d --build

# 查看启动日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

### 4. 访问服务
| 服务 | 访问地址 | 端口映射 | 说明 |
|------|----------|----------|------|
| **MySQL** | `mysql:3306` | - | 仅容器内访问 |
| **Redis** | `redis:6379` | - | 仅容器内访问 |
| **后端API** | `http://localhost:8080` | 8080 | Go Gin 服务 |
| **用户端前端** | `http://localhost:3000` | 3000 | Vue 3 用户端 |
| **管理后台** | `http://localhost:3001` | 3001 | Vue 3 管理后台 |

### 5. 数据库初始化
后端启动时会自动执行 GORM AutoMigrate，创建：
- users, photosets, photos, orders, order_items
- tags, favorites, categories (新加的分类系统)

如果需要手动初始化数据：
```bash
# 进入 MySQL 容器
docker-compose exec mysql mysql -u root -p

# 查看数据库
SHOW DATABASES;
USE photoset;
SHOW TABLES;
```

## 🔧 生产部署建议

### 1. 使用 Docker Swarm 或 Kubernetes
```bash
# Docker Swarm 初始化
docker swarm init
docker stack deploy -c docker-compose.yml photoset
```

### 2. HTTPS 配置
在 Nginx 配置中加入 SSL 证书：
```nginx
# docker/nginx-frontend.conf 和 docker/nginx-admin.conf 中添加
listen 443 ssl;
ssl_certificate /etc/ssl/certs/your-cert.pem;
ssl_certificate_key /etc/ssl/private/your-key.pem;
```

### 3. 数据持久化
- `mysql_data` 卷保存 MySQL 数据
- `redis_data` 卷保存 Redis 数据
- `uploads_volume` 卷保存上传的图片

备份数据：
```bash
docker-compose exec mysql sh -c 'exec mysqldump --all-databases -uroot -p"$MYSQL_ROOT_PASSWORD"' > backup.sql

docker-compose exec redis redis-cli SAVE
```

### 4. 监控与日志
```bash
# 查看实时日志
docker-compose logs -f backend
docker-compose logs -f frontend

# 容器资源使用
docker stats

# 进入容器调试
docker-compose exec backend sh
```

## 🛠️ 开发模式

### 1. 使用热重载（开发环境）
创建 `docker-compose.dev.yml`：
```yaml
version: '3.8'

services:
  backend:
    build:
      context: .
      dockerfile: ./docker/backend.Dockerfile.dev
    volumes:
      - .:/app
      - go-mod-cache:/go/pkg/mod
    environment:
      APP_ENV: development
      DEBUG: true
    command: air -c .air.toml

  frontend:
    build:
      context: .
      dockerfile: ./docker/frontend.Dockerfile.dev
    volumes:
      - ./frontend:/app
      - /app/node_modules
    environment:
      NODE_ENV: development
    command: npm run dev

  admin:
    build:
      context: .
      dockerfile: ./docker/admin.Dockerfile.dev
    volumes:
      - ./frontend-admin:/app
      - /app/node_modules
    environment:
      NODE_ENV: development
    command: npm run dev

volumes:
  go-mod-cache:
```

运行开发环境：
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

## ⚙️ 容器管理命令

```bash
# 停止服务
docker-compose down

# 停止并清理数据卷
docker-compose down -v

# 重建单个服务
docker-compose up -d --no-deps --build backend

# 执行数据库迁移
docker-compose exec backend /app/main -migrate

# 查看容器日志
docker-compose logs --tail=100 backend

# 进入容器
docker-compose exec backend sh
docker-compose exec mysql mysql -u photoset -p photoset123
```

## 🐛 常见问题

### Q1: 数据库连接失败
检查 MySQL 容器是否启动：
```bash
docker-compose logs mysql
docker-compose exec mysql mysqladmin ping
```

### Q2: 前端无法访问后端 API
检查 CORS 配置：
```bash
# 确保 .env.docker 中 CORS_ALLOW_ORIGINS 包含前端地址
CORS_ALLOW_ORIGINS=http://frontend:80,http://admin:80
```

### Q3: 图片上传失败
检查上传目录权限：
```bash
docker-compose exec backend ls -la /app/uploads
docker-compose exec backend chmod 777 /app/uploads
```

### Q4: 内存不足
调整容器资源限制：
```yaml
# docker-compose.yml 中为每个服务添加
resources:
  limits:
    memory: 1G
  reservations:
    memory: 256M
```

## 📊 性能调优

### 1. 启用多阶段构建
所有 Dockerfile 已使用多阶段构建，减小镜像体积：
- Go: 从 ~1.2GB 减少到 ~30MB
- Node: 从 ~1.5GB 减少到 ~20MB

### 2. 使用健康检查
各服务都有健康检查，确保容器正常工作后才开始服务依赖。

### 3. 合理配置缓存
- Redis 缓存数据
- Nginx 缓存静态资源
- 浏览器缓存静态文件

## 🔒 安全建议

1. **生产环境必须修改**：
   - JWT_SECRET
   - SIGN_SECRET  
   - MySQL 密码
   - Redis 密码（如果需要）

2. **使用 Docker Secrets**：
```bash
echo "your-jwt-secret" | docker secret create jwt_secret -
```

3. **限制网络暴露**：
   - MySQL、Redis 不要对外暴露端口
   - 使用内部网络通信

4. **定期更新镜像**：
```bash
docker-compose pull
docker-compose up -d
```

## 📈 监控集成

### 1. Prometheus 监控
为后端添加指标端点：
```go
// internal/http/routes/routes.go 中添加
r.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### 2. 日志聚合
使用 ELK 或 Loki：
```yaml
logging:
  driver: loki
  options:
    loki-url: "http://loki:3100/api/prom/push"
```

## 📚 相关文档

- [Docker Compose 文档](https://docs.docker.com/compose/)
- [Go 官方多阶段构建指南](https://docs.docker.com/language/golang/)
- [Nginx 容器化配置](https://hub.docker.com/_/nginx)
- [MySQL Docker 配置](https://hub.docker.com/_/mysql)

---

**维护者**: PhotoSet 团队  
**最后更新**: 2026-04-15  
**状态**: ✅ 生产就绪 — Phase 1~6 全部完成