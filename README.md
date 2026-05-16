# PhotoSet 摄影套图浏览平台

> 专业的摄影套图内容浏览与会员订阅平台

---

## 项目概述

| 项目 | 说明 |
|------|------|
| **项目名称** | PhotoSet |
| **技术栈** | Go + Gin / Vue 3 / MySQL + Redis / Docker |
| **项目状态** | ✅ Phase 1~6 全部完成 |
| **访问地址** | https://tt.cy.mk |
| **管理后台** | https://admin.tt.cy.mk |

---

## 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                        用户端 (Vue 3)                        │
│                     端口 3000 / Nginx                       │
├─────────────────────────────────────────────────────────────┤
│                      管理后台 (Vue 3)                        │
│                     端口 3001 / Nginx                        │
├─────────────────────────────────────────────────────────────┤
│                       后端 API (Go)                         │
│                       端口 8080                              │
│  ┌──────────┬───────────┬──────────┬───────────┐           │
│  │ Handlers │ Middleware │ Services  │ Repos     │           │
│  └──────────┴───────────┴──────────┴───────────┘           │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │    MySQL     │    │    Redis     │    │  Cloudflare  │  │
│  │   photoset   │    │    缓存      │    │      R2      │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 后端技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 框架 | Go + Gin | 高性能 Web 框架 |
| ORM | GORM | MySQL 数据库操作 |
| 认证 | JWT | HS256 签名 |
| 缓存 | Redis | 多级缓存策略 |
| 存储 | S3/R2 SDK | 云存储支持 |
| 搜索 | MySQL FULLTEXT | ngram 中文分词 |

### 前端技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 框架 | Vue 3 | Composition API |
| UI 库 | Element Plus | 管理后台 |
| 状态 | Pinia | 状态管理 |
| 构建 | Vite | 快速构建 |
| 图表 | ECharts | Dashboard 统计 |

---

## 项目目录结构

```
backend/
├── cmd/
│   └── main.go                 # 程序入口
│
├── internal/
│   ├── config/
│   │   └── config.go           # 配置加载 (.env)
│   │
│   ├── database/
│   │   └── database.go         # MySQL/Redis 连接初始化
│   │
│   ├── domain/                 # 数据模型（实体层）
│   │   ├── user.go             # 用户模型
│   │   ├── photoset.go         # 套图模型
│   │   ├── photo.go            # 图片模型
│   │   ├── order.go            # 订单模型
│   │   ├── membership.go       # 会员套餐模型
│   │   ├── favorite.go         # 收藏模型
│   │   ├── tag.go              # 标签模型
│   │   ├── category.go         # 分类模型
│   │   ├── page.go             # 静态页面模型
│   │   ├── site_setting.go     # 站点设置模型
│   │   └── admin_log.go       # 操作日志模型
│   │
│   ├── http/
│   │   ├── handlers/           # HTTP 处理器
│   │   │   ├── auth.go         # 认证相关
│   │   │   ├── photoset.go     # 套图 CRUD
│   │   │   ├── order_handler.go # 订单管理
│   │   │   ├── admin_handler.go # 管理后台
│   │   │   ├── upload_handler.go # 文件上传
│   │   │   └── ...
│   │   ├── middleware/         # 中间件
│   │   │   ├── auth.go         # JWT 认证
│   │   │   ├── cors.go         # 跨域
│   │   │   └── sign.go         # URL 签名验证
│   │   └── routes/
│   │       └── routes.go       # 路由配置
│   │
│   ├── repository/             # 数据访问层
│   │   ├── user_repository.go
│   │   ├── photoset_repository.go
│   │   ├── order_repository.go
│   │   └── ...
│   │
│   ├── service/                # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── photoset_service.go
│   │   ├── cache_service.go    # Redis 缓存
│   │   ├── mail.go            # 邮件服务
│   │   └── watermark.go        # 水印处理
│   │
│   └── storage/                # 存储抽象层
│       ├── factory.go          # 存储工厂
│       ├── storage.go          # 存储接口
│       └── local.go            # 本地存储实现
│
├── frontend/                   # Vue 3 用户端
│   ├── src/
│   │   ├── api/               # API 封装
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 公共组件
│   │   ├── stores/            # Pinia 状态
│   │   └── router/            # 路由配置
│   └── package.json
│
├── frontend-admin/             # Vue 3 管理后台
│   ├── src/
│   │   ├── views/
│   │   │   ├── Dashboard.vue   # 统计面板
│   │   │   ├── ContentReview.vue # 内容审核
│   │   │   ├── UserManage.vue  # 用户管理
│   │   │   ├── OrderManage.vue # 订单管理
│   │   │   ├── TagManage.vue   # 标签管理
│   │   │   ├── SiteSettings.vue # 站点设置
│   │   │   └── Pages.vue       # 页面管理
│   │   └── ...
│   └── package.json
│
├── docker/                     # Docker 配置
│   ├── backend.Dockerfile
│   ├── frontend.Dockerfile
│   ├── admin.Dockerfile
│   ├── nginx-frontend.conf
│   ├── nginx-admin.conf
│   └── README.DOCKER.md
│
├── migrations/                 # 数据库迁移
│   └── *.sql
│
├── nginx/                      # Nginx 配置
├── logs/                       # 日志目录
├── uploads/                    # 上传文件目录
│
├── docker-compose.yml          # 开发环境
├── docker-compose.prod.yml     # 生产环境
├── docker-compose.server.yml   # 服务器部署
├── go.mod
└── go.sum
```

---

## 数据库设计

### 核心数据表

| 表名 | 说明 | 关联 |
|------|------|------|
| `users` | 用户表 | - |
| `photosets` | 套图表 | user_id, category |
| `photos` | 图片表 | photoset_id |
| `photoset_tags` | 套图-标签关联表 | photoset_id, tag_id |
| `tags` | 标签表 | - |
| `categories` | 分类表 | - |
| `orders` | 订单表 | user_id, membership_id, photoset_id |
| `memberships` | 会员套餐表 | - |
| `favorites` | 收藏表 | user_id, photoset_id |
| `pages` | 静态页面表 | - |
| `site_settings` | 站点设置表 | key-value |
| `admin_logs` | 操作日志表 | - |
| `api_keys` | API 密钥表 | - |

### 用户角色权限

| 角色 | 说明 | 权限 |
|------|------|------|
| `guest` | 游客 | 浏览公开内容 |
| `user` | 普通用户 | 浏览 + 购买 |
| `member` | 会员 | 浏览 + 购买 + 会员权益 |
| `creator` | 创作者 | 发布/编辑套图 |
| `admin` | 管理员 | 全站管理 |

---

## API 接口文档

### 认证接口 `/api/auth`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/auth/captcha` | 获取图形验证码 | 公开 |
| POST | `/api/auth/register` | 用户注册 | 公开 |
| POST | `/api/auth/login` | 用户登录 | 公开 |
| GET | `/api/auth/me` | 获取当前用户 | 可选认证 |
| PUT | `/api/auth/password` | 修改密码 | 登录 |
| POST | `/api/auth/forgot-password` | 忘记密码 | 公开 |
| POST | `/api/auth/reset-password` | 重置密码 | 公开 |

### 套图接口 `/api/photosets`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/photosets` | 套图列表（基础搜索） | 公开 |
| GET | `/api/photosets/advanced` | 高级搜索 | 公开 |
| GET | `/api/photosets/:id` | 套图详情 | 公开 |
| POST | `/api/photosets` | 创建套图 | 创作者+ |
| PUT | `/api/photosets/:id` | 更新套图 | 创作者（本人） |
| DELETE | `/api/photosets/:id` | 删除套图 | 创作者/Admin |

**高级搜索参数：**
```
GET /api/photosets/advanced?
  keyword=关键词&
  category=分类&
  min_price=0&
  max_price=100&
  creator_id=1&
  sort=newest|oldest|price_asc|price_desc&
  page=1&
  page_size=20
```

### 收藏接口 `/api/favorites`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/favorites` | 我的收藏列表 | 登录 |
| POST | `/api/favorites/:photosetId` | 添加收藏 | 登录 |
| DELETE | `/api/favorites/:photosetId` | 取消收藏 | 登录 |

### 订单接口 `/api/orders`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/orders` | 我的订单 | 登录 |
| POST | `/api/orders` | 创建订单 | 登录 |
| POST | `/api/orders/:id/pay` | 模拟支付 | 登录 |
| POST | `/api/orders/:id/refund` | 用户退款（48h内） | 登录 |

### 管理后台 `/api/admin`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/stats` | 统计概览 |
| GET | `/admin/stats/trend` | 趋势数据（7/14/30天） |
| GET | `/admin/users` | 用户列表 |
| PUT | `/admin/users/:id/ban` | 封号/解封 |
| PUT | `/admin/users/:id/role` | 修改角色 |
| GET | `/admin/photosets` | 套图列表（按状态） |
| POST | `/admin/photosets/:id/approve` | 审核通过 |
| POST | `/admin/photosets/:id/reject` | 审核拒绝 |
| DELETE | `/admin/photosets/:id` | 删除套图 |
| GET | `/admin/orders` | 订单列表 |
| POST | `/admin/orders/:id/refund` | 管理员退款 |
| GET/POST/PUT/DELETE | `/admin/tags` | 标签管理 |
| GET/POST/PUT/DELETE | `/admin/categories` | 分类管理 |
| GET/PUT | `/admin/settings` | 站点设置 |
| GET/POST/PUT/DELETE | `/admin/pages` | 页面管理 |
| GET | `/admin/logs` | 操作日志 |

---

## 功能特性

### Phase 1: 基础架构
- ✅ 用户注册/登录/JWT认证
- ✅ 密码重置（邮件）
- ✅ 图形验证码
- ✅ 角色权限控制

### Phase 2: 核心功能
- ✅ 套图 CRUD
- ✅ 图片上传（本地/S3/R2）
- ✅ 标签系统
- ✅ 收藏功能
- ✅ 搜索功能

### Phase 3: 前端实现
- ✅ Vue 3 用户端
- ✅ 响应式设计
- ✅ 高级搜索过滤器
- ✅ 会员套餐展示

### Phase 4: 支付系统
- ✅ 订单系统
- ✅ 模拟支付
- ✅ 用户退款（48h）
- ✅ 管理后台 Dashboard
- ✅ 内容审核
- ✅ 用户管理

### Phase 5: 高级功能
- ✅ Redis 多级缓存（列表5min/详情10min/标签30min）
- ✅ MySQL FULLTEXT 中文全文搜索（ngram）
- ✅ URL 签名防盗链（HMAC-SHA256）
- ✅ 管理员无限期退款

### Phase 6: 站点设置
- ✅ 基本信息管理
- ✅ SEO 关键词设置
- ✅ 关于页面管理
- ✅ 邮件配置（SMTP）
- ✅ 水印设置（文字/图片）
- ✅ 存储配置（本地/S3/R2）
- ✅ 静态页面管理
- ✅ 操作日志

---

## 环境变量配置

### 后端 `.env`

```bash
# 服务配置
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=photoset

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRE_HOURS=24

# 存储类型: local / s3 / r2
STORAGE_TYPE=local
LOCAL_STORAGE_PATH=./uploads

# S3/R2 配置
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=
S3_REGION=
R2_ACCOUNT_ID=
R2_PUBLIC_URL=

# URL 签名
SIGN_SECRET=your-sign-secret
SIGN_EXPIRE=7200
```

---

## 部署指南

### Docker 一键部署

```bash
# 1. 克隆项目
git clone https://github.com/lovechu/photoset.git
cd photoset/backend

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 填入配置

# 3. 构建并启动
docker-compose up -d --build

# 4. 查看服务状态
docker-compose ps
```

### 宝塔面板部署

1. **上传代码** 到 `/www/dk_project/wwwroot/tt.cy.mk/`

2. **添加伪静态规则**（网站 → tt.cy.mk → 伪静态）：

```nginx
# 上传文件代理
location /uploads/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

# 管理后台静态资源
location /admin/assets/ {
    proxy_pass http://127.0.0.1:3001/assets/;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
}

# 管理后台
location /admin/ {
    proxy_pass http://127.0.0.1:3001/;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

# 后端 API
location ~ ^/api/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

3. **启动 Docker 容器**：
```bash
cd /opt/photoset
docker-compose -f docker-compose.server.yml up -d
```

### Docker 运维命令

```bash
# 查看容器状态
docker ps

# 查看日志
docker logs -f photoset-backend
docker logs -f photoset-frontend
docker logs -f photoset-admin

# 重启容器
docker restart photoset-backend
docker restart photoset-frontend
docker restart photoset-admin

# 进入容器
docker exec -it photoset-backend sh

# 备份数据库
docker exec photoset-mysql mysqldump -u root -p<password> photoset > backup_$(date +%Y%m%d).sql

# 重建后端
docker-compose -f docker-compose.server.yml up -d --no-deps --build backend
```

---

## 开发环境

### 后端 (WSL)

```bash
# 启动 MySQL 和 Redis
sudo service mysql start
sudo service redis-server start

# 启动后端
cd backend/
go run cmd/main.go
# 监听端口: 8080
```

### 前端用户端

```bash
cd frontend/
npm install
npm run dev
# 访问: http://localhost:3000
```

### 管理后台

```bash
cd frontend-admin/
npm install
npm run dev
# 访问: http://localhost:3001
```

---

## 项目进度

| Phase | 内容 | 状态 |
|-------|------|------|
| Phase 1 | 后端基础架构 + 用户认证 | ✅ |
| Phase 2 | 套图核心功能 API | ✅ |
| Phase 3 | Web 前端完整实现 | ✅ |
| Phase 4 | 会员支付系统 + 管理后台 | ✅ |
| Phase 4 补齐 | Cloudflare R2 + 套图编辑 | ✅ |
| Phase 5 | Redis缓存/FULLTEXT/退款 | ✅ |
| Phase 6 | 站点设置 + 页面管理 | ✅ |

---

## 相关文档

- [后端 API 文档](./docker/) — Docker 部署指南
- [前端 README](./frontend/README.md) — 用户端详细文档
- [管理后台 README](./frontend-admin/README.md) — 后台详细文档

---

*最后更新: 2026-05-16*
