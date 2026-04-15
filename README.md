# PhotoSet - 摄影套图浏览平台

基于 Go + Gin + Vue 3 的全栈摄影套图平台，前后端一体仓库。

## GitHub 仓库
https://github.com/lovechu/photoset

## 项目状态
**全部功能开发完成，生产就绪**

| Phase | 内容 | 状态 |
|-------|------|------|
| Phase 1 | 基础架构 | ✅ |
| Phase 2 | 核心功能（收藏/个人中心/搜索/上传） | ✅ |
| Phase 3 | 会员订单系统 | ✅ |
| Phase 4 | 管理后台 + 套图编辑 + Cloudflare R2 | ✅ |
| Phase 5 | 分类/Redis缓存/URL签名/全文搜索/退款 | ✅ |
| Phase 6 | 站点设置 + 页面管理 + 管理后台增强 | ✅ |

## 核心功能

| 模块 | 功能 |
|------|------|
| **后端 (Go)** | 用户/套图/订单/标签/分类/搜索/验证码/JWT/缓存/站点设置/页面管理/操作日志 |
| **前端用户端 (Vue 3)** | 首页/详情/上传/编辑/个人中心/收藏/订单/搜索/分类/静态页面/会员 |
| **管理后台 (Vue 3)** | 登录/Dashboard/内容审核/用户管理/订单管理/标签管理/分类管理/站点设置/页面管理/操作日志 |
| **存储** | 本地存储 + S3 兼容（Cloudflare R2/AWS S3/MinIO/阿里 OSS） |
| **CDN** | 支持自定义 CDN 域名，后台可配置 |
| **API** | 完整 RESTful + 安全中间件 + 签名防盗链 |

## 项目结构

```
backend/                  (本仓库根目录)
├── cmd/main.go          # 程序入口
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # MySQL 连接
│   ├── domain/          # 领域模型
│   │   ├── user.go / photoset.go / order.go / tag.go / category.go
│   │   ├── admin_log.go / page.go / site_setting.go
│   │   └── membership.go / photo.go / favorite.go
│   ├── http/
│   │   ├── handlers/    # 请求处理器
│   │   ├── middleware/  # 中间件（JWT/CORS/限流/签名验证）
│   │   └── routes/      # 路由注册
│   ├── pkg/             # 工具包 (jwt/password/response/signurl)
│   ├── repository/      # 数据仓储层
│   ├── service/         # 业务逻辑层
│   └── storage/         # 存储抽象（本地/S3兼容）
│       ├── storage.go   # Storage 接口
│       ├── local.go     # 本地存储
│       ├── s3.go        # S3 兼容存储
│       └── factory.go   # 存储工厂（从配置动态创建）
├── frontend/            # 用户端 Web (Vue 3, 端口 3000)
├── frontend-admin/      # 管理后台 (Vue 3, 端口 3001)
├── docker/              # Docker 部署配置
├── scripts/             # 数据库脚本
├── migrations/          # 数据库迁移
└── .env.example         # 环境变量模板
```

## 技术栈

| 层面 | 技术 |
|------|------|
| 后端框架 | Go + Gin |
| ORM | GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis（限流/验证码/业务缓存，静默降级） |
| 对象存储 | AWS SDK v2（S3 兼容，支持 R2/OSS/S3/MinIO） |
| 全文搜索 | MySQL FULLTEXT（ngram 中文分词） |
| 图片防盗链 | HMAC-SHA256 URL 签名 |
| 认证 | JWT + bcrypt |
| 前端 | Vue 3 + Vite + Pinia + Vue Router + Element Plus |
| 管理后台 | Vue 3 + Vite + Pinia + Element Plus + ECharts |
| 容器化 | Docker + Docker Compose |

## 快速开始

### 环境要求
- Go 1.22+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+（可选，无 Redis 时自动降级）

### 后端
```bash
cp .env.example .env          # 复制并编辑环境变量
go mod download                # 下载依赖
go run cmd/main.go             # 启动后端 (端口 8080)
```

### 前端用户端
```bash
cd frontend
npm install
npm run dev                    # 启动 (端口 3000)
```

### 管理后台
```bash
cd frontend-admin
npm install
npm run dev                    # 启动 (端口 3001)
```

### Docker 部署
```bash
cp .env.docker.example .env.docker  # 编辑环境变量
docker-compose up -d --build        # 一键启动全部服务
```

详细 Docker 部署文档见 `docker/README.DOCKER.md`

## 环境变量说明

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASS` / `DB_NAME` | MySQL 连接 | - |
| `REDIS_HOST` / `REDIS_PORT` / `REDIS_PASS` | Redis 连接（可选） | localhost:6379 |
| `JWT_SECRET` | JWT 签名密钥 | - |
| `SIGN_SECRET` | URL 签名密钥 | - |
| `STORAGE_TYPE` | 存储类型：local / s3 | local |
| `S3_ENDPOINT` / `S3_REGION` / `S3_BUCKET` / `S3_ACCESS_KEY` / `S3_SECRET_KEY` | S3 配置 | - |
| `CDN_DOMAIN` | CDN 公开域名 | - |
| `SIGN_EXPIRE` | 签名 URL 过期时间（秒） | 3600 |

## Phase 6 新增功能 (2026-04-15)

### 站点设置系统
- `SiteSetting` 模型 + key-value 存储结构
- 管理后台"站点设置"页面（5 个 Tab）：
  - **基本信息**：站点名称、描述、Logo、备案号
  - **SEO 设置**：关键词、描述
  - **关于我**：个人简介、联系方式
  - **邮件设置**：SMTP 配置
  - **水印设置**：文字/图片水印配置

### 页面管理系统
- `Page` 模型 + CRUD API
- 管理后台"页面管理"：自定义静态页面
- 前端 `/p/:slug` 路由渲染静态页面

### 管理后台增强
- **用户管理**：搜索/筛选、角色修改、用户详情抽屉（统计数据）
- **Dashboard**：ECharts 趋势折线图（7/14/30 天）、真实收入数据
- **操作日志**：自动记录封号/审核/退款等操作，日志列表 + 筛选

### 存储与 CDN 配置
- 后台设置中新增"存储与 CDN"Tab
- 支持 S3 兼容存储（R2/OSS/S3/MinIO）配置
- CDN 域名配置
- 测试连接功能

## 安全特性

- JWT Token 认证 + 过期控制
- 多级权限系统（guest/user/member/creator/admin）
- IP 限流防暴力破解
- 图形验证码防机器人
- HMAC-SHA256 URL 签名防盗链
- 输入验证 + XSS/CSRF 防护

## 相关文档

- `frontend/README.md` — 前端用户端文档
- `frontend-admin/README.md` — 管理后台文档
- `docker/README.DOCKER.md` — Docker 部署指南
- `.env.example` — 环境变量模板

---

**最后更新**: 2026-04-15
**项目状态**: ✅ 生产就绪 — Phase 1~6 全部完成
**维护者**: PhotoSet 开发团队
