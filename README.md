# PhotoSet 摄影套图浏览平台

> 专业的摄影套图内容浏览与会员订阅平台，支持社区互动、消息通知、等级成长体系

---

## 项目概述

| 项目 | 说明 |
|------|------|
| **项目名称** | PhotoSet |
| **技术栈** | Go + Gin / Vue 3 / Flutter / MySQL + Redis / Docker |
| **项目状态** | ✅ Phase 1~6 完成 + 社区功能上线 |
| **访问地址** | https://tt.cy.mk |
| **管理后台** | https://admin.tt.cy.mk |
| **App（Flutter）** | 支持 Android / iOS |

---

## 技术架构

```
┌─────────────────────────────────────────────────────┐
│                        用户端 (Vue 3)                        │
│                     端口 3000 / Nginx                       │
├─────────────────────────────────────────────────────┤
│                      管理后台 (Vue 3)                        │
│                     端口 3001 / Nginx                        │
├─────────────────────────────────────────────────────┤
│                      Flutter App                           │
│                  Android / iOS / 移动端                       │
├─────────────────────────────────────────────────────┤
│                       后端 API (Go)                         │
│                       端口 8080                              │
│  ┌──────────┬───────────┬──────────┬───────────┐           │
│  │ Handlers │ Middleware │ Services  │ Repos     │           │
│  └──────────┴───────────┴──────────┴───────────┘           │
├─────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │    MySQL     │    │    Redis     │    │  Cloudflare  │  │
│  │   photoset   │    │    缓存      │    │      R2      │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
└─────────────────────────────────────────────────────┘
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
| 日志 | Zap | 结构化日志 |

### 前端技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 框架 | Vue 3 | Composition API |
| UI 库 | Element Plus | 管理后台 |
| 状态 | Pinia | 状态管理 |
| 构建 | Vite | 快速构建 |
| 图表 | ECharts | Dashboard 统计 |

### Flutter App 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 框架 | Flutter 3 | 跨平台 |
| 状态 | BLoC / GetIt | 状态管理 + 依赖注入 |
| 网络 | Retrofit + Dio | HTTP 客户端 |
| 本地 | Hive | 本地存储 |
| 路由 | go_router | 路由管理 |
| UI | flutter_screenutil | 屏幕适配 |

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
│   │   ├── admin_log.go       # 操作日志模型
│   │   ├── post.go             # 社区帖子模型
│   │   ├── post_reply.go      # 帖子回复模型
│   │   ├── post_like.go       # 帖子点赞模型
│   │   ├── post_reply_like.go # 回复点赞模型
│   │   ├── user_point.go      # 用户积分模型
│   │   ├── sensitive_word.go   # 敏感词模型
│   │   ├── post_report.go     # 帖子举报模型
│   │   ├── message.go          # 消息模型
│   │   ├── notification.go     # 通知模型
│   │   ├── achievement.go      # 成就模型
│   │   ├── draft.go            # 草稿模型
│   │   ├── post_share.go      # 帖子分享模型
│   │   ├── topic.go            # 话题模型
│   │   ├── tag.go              # 标签模型（社区）
│   │   ├── points_mall.go      # 积分商城模型
│   │   └── user_level_config.go # 用户等级配置
│   │
│   ├── http/
│   │   ├── handlers/           # HTTP 处理器
│   │   │   ├── auth.go         # 认证相关
│   │   │   ├── photoset.go     # 套图 CRUD
│   │   │   ├── order_handler.go # 订单管理
│   │   │   ├── admin_handler.go # 管理后台
│   │   │   ├── upload_handler.go # 文件上传
│   │   │   ├── community_handler.go # 社区帖子
│   │   │   ├── message_handler.go   # 消息通知
│   │   │   ├── notification_handler.go # 通知
│   │   │   └── user_level_handler.go   # 等级/成就/积分商城
│   │   │
│   │   ├── middleware/         # 中间件
│   │   │   ├── auth.go         # JWT 认证
│   │   │   ├── cors.go         # 跨域
│   │   │   ├── sign.go         # URL 签名验证
│   │   │   └── logger.go       # 请求日志
│   │   │
│   │   └── routes/
│   │       └── routes.go       # 路由配置
│   │
│   ├── repository/             # 数据访问层
│   │   ├── user_repository.go
│   │   ├── photoset_repository.go
│   │   ├── order_repository.go
│   │   ├── post_repository.go
│   │   ├── post_reply_repository.go
│   │   ├── post_like_repository.go
│   │   ├── user_point_repository.go
│   │   ├── sensitive_word_repository.go
│   │   ├── post_report_repository.go
│   │   ├── message_repository.go
│   │   ├── notification_repository.go
│   │   ├── achievement_repository.go
│   │   ├── draft_repository.go
│   │   ├── post_share_repository.go
│   │   ├── post_tag_repository.go
│   │   ├── post_topic_repository.go
│   │   ├── topic_repository.go
│   │   ├── tag_repository.go
│   │   ├── points_mall_repository.go
│   │   └── user_level_repository.go
│   │
│   ├── service/                # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── photoset_service.go
│   │   ├── cache_service.go    # Redis 缓存
│   │   ├── mail.go            # 邮件服务
│   │   ├── watermark.go        # 水印处理
│   │   ├── community_service.go  # 社区业务
│   │   ├── mention_service.go   # @提及服务
│   │   ├── message_service.go   # 消息服务
│   │   ├── notification_service.go # 通知服务
│   │   ├── recommendation_service.go # 推荐服务
│   │   └── user_level_service.go   # 等级/成就/积分商城
│   │
│   ├── storage/                # 存储抽象层
│   │   ├── factory.go          # 存储工厂
│   │   ├── storage.go          # 存储接口
│   │   └── local.go            # 本地存储实现
│   │
│   └── logger/                # 日志模块
│       └── logger.go           # Zap 日志配置
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
│   │   │   ├── Pages.vue       # 页面管理
│   │   │   ├── PostManage.vue  # 社区帖子管理
│   │   │   ├── ReplyManage.vue # 社区回帖管理
│   │   │   ├── SensitiveWordManage.vue # 敏感词管理
│   │   │   ├── ReportManage.vue # 举报处理
│   │   │   └── PointsManage.vue # 用户积分管理
│   │   └── ...
│   └── package.json
│
├── photoset_mobile/            # Flutter App
│   ├── lib/
│   │   ├── blocs/            # BLoC 状态管理
│   │   ├── models/           # 数据模型
│   │   ├── repositories/      # 数据仓储
│   │   ├── screens/          # 页面
│   │   └── widgets/         # 组件
│   └── pubspec.yaml
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
| `post_favorites` | 帖子收藏表 | user_id, post_id |
| `follows` | 用户关注表 | follower_id, following_id |

### 社区功能表

| 表名 | 说明 | 关联 |
|------|------|------|
| `posts` | 社区帖子表 | author_id, category |
| `post_replies` | 帖子回复表 | post_id, author_id, reply_to |
| `post_likes` | 帖子点赞表 | user_id, post_id |
| `post_reply_likes` | 回复点赞表 | user_id, reply_id |
| `user_points` | 用户积分表 | user_id |
| `sensitive_words` | 敏感词表 | - |
| `post_reports` | 帖子举报表 | post_id, reporter_id |
| `messages` | 私信消息表 | from_user_id, to_user_id |
| `notifications` | 通知表 | user_id, actor_id |
| `achievements` | 成就配置表 | - |
| `user_achievements` | 用户成就记录 | user_id, achievement_id |
| `drafts` | 草稿表 | user_id |
| `post_shares` | 帖子分享表 | post_id, user_id |
| `post_tags` | 帖子标签关联表 | post_id, tag_id |
| `post_topics` | 帖子话题关联表 | post_id, topic_id |
| `topics` | 话题表 | creator_id |
| `tags` | 标签表（社区） | - |
| `points_mall_items` | 积分商城物品表 | - |
| `points_mall_exchanges` | 积分兑换记录 | user_id, item_id |
| `user_level_configs` | 用户等级配置 | - |

### 用户角色权限

| 角色 | 说明 | 权限 |
|------|------|------|
| `guest` | 游客 | 浏览公开内容 |
| `user` | 普通用户 | 浏览 + 购买 + 社区互动 |
| `member` | 会员 | 浏览 + 购买 + 会员权益 + 社区 |
| `creator` | 创作者 | 发布/编辑套图 + 社区 |
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

### 社区帖子接口 `/api/community/posts`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/community/posts` | 帖子列表 | 公开 |
| GET | `/api/community/posts/:id` | 帖子详情 | 公开 |
| POST | `/api/community/posts` | 发帖 | 登录 |
| PUT | `/api/community/posts/:id` | 编辑帖子 | 作者/Admin |
| DELETE | `/api/community/posts/:id` | 删除帖子 | 作者/Admin |
| POST | `/api/community/posts/:id/like` | 点赞/取消 | 登录 |
| POST | `/api/community/posts/:id/favorite` | 收藏/取消 | 登录 |
| POST | `/api/community/posts/:id/share` | 分享帖子 | 登录 |
| GET | `/api/community/posts/user/:id` | 指定用户帖子 | 公开 |

### 社区回帖接口 `/api/community/posts/:id/replies`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/community/posts/:id/replies` | 回复列表（楼中楼） | 公开 |
| POST | `/api/community/posts/:id/replies` | 发表回复 | 登录 |
| POST | `/api/community/replies/:id/like` | 点赞回复 | 登录 |

### 社区分类/标签/话题接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/community/categories` | 分类列表 | 公开 |
| GET | `/api/community/tags` | 标签列表 | 公开 |
| GET | `/api/community/topics` | 话题列表 | 公开 |
| GET | `/api/community/topics/:id` | 话题详情+帖子 | 公开 |

### 关注接口 `/api/community/follows`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/community/follows/:id` | 关注/取关 | 登录 |
| GET | `/api/community/follows/following` | 我的关注 | 登录 |
| GET | `/api/community/follows/followers` | 我的粉丝 | 登录 |

### 消息接口 `/api/messages`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/messages/conversations` | 会话列表 | 登录 |
| GET | `/api/messages/conversations/:id` | 与指定用户对话 | 登录 |
| POST | `/api/messages` | 发送私信 | 登录 |
| PUT | `/api/messages/:id/read` | 标记已读 | 登录 |

### 通知接口 `/api/notifications`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/notifications` | 通知列表 | 登录 |
| PUT | `/api/notifications/:id/read` | 标记已读 | 登录 |
| PUT | `/api/notifications/read-all` | 全部已读 | 登录 |

### 用户等级/成就/积分商城接口 `/api/user`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/user/level` | 我的等级信息 | 登录 |
| GET | `/api/user/achievements` | 我的成就 | 登录 |
| GET | `/api/user/points/history` | 积分明细 | 登录 |
| GET | `/api/user/exchange-history` | 兑换记录 | 登录 |
| GET | `/api/mall/items` | 积分商城列表 | 登录 |
| POST | `/api/mall/exchange` | 兑换物品 | 登录 |

### 管理后台 `/api/admin`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/stats` | 统计概览 |
| GET | `/api/admin/stats/trend` | 趋势数据（7/14/30天） |
| GET | `/api/admin/users` | 用户列表 |
| PUT | `/api/admin/users/:id/ban` | 封号/解封 |
| PUT | `/api/admin/users/:id/role` | 修改角色 |
| GET | `/api/admin/photosets` | 套图列表（按状态） |
| POST | `/api/admin/photosets/:id/approve` | 审核通过 |
| POST | `/api/admin/photosets/:id/reject` | 审核拒绝 |
| DELETE | `/api/admin/photosets/:id` | 删除套图 |
| GET | `/api/admin/orders` | 订单列表 |
| POST | `/api/admin/orders/:id/refund` | 管理员退款 |
| GET/POST/PUT/DELETE | `/api/admin/tags` | 标签管理 |
| GET/POST/PUT/DELETE | `/api/admin/categories` | 分类管理 |
| GET/PUT | `/api/admin/settings` | 站点设置 |
| GET/POST/PUT/DELETE | `/api/admin/pages` | 页面管理 |
| GET | `/api/admin/logs` | 操作日志 |
| GET | `/api/admin/community/posts` | 社区帖子管理 |
| GET | `/api/admin/community/replies` | 社区回帖管理 |
| GET/POST/PUT/DELETE | `/api/admin/community/categories` | 社区分类管理 |
| GET/POST/PUT/DELETE | `/api/admin/sensitive-words` | 敏感词管理 |
| GET | `/api/admin/community/reports` | 举报管理 |
| GET | `/api/admin/points` | 用户积分管理 |
| PUT | `/api/admin/points/:id/adjust` | 调整用户积分 |

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

### Phase 7: 社区功能 ✅
- ✅ 发帖/回帖（楼中楼）
- ✅ 点赞 Toggle
- ✅ 举报功能
- ✅ 敏感词过滤（***替换）
- ✅ 积分等级系统
- ✅ 帖子收藏
- ✅ 用户关注/取关
- ✅ 私信消息
- ✅ 通知系统
- ✅ 成就系统
- ✅ 积分商城
- ✅ 草稿功能
- ✅ 帖子分享
- ✅ 话题/标签
- ✅ 推荐服务

### Flutter App ✅
- ✅ 图片浏览/个人中心/收藏功能
- ✅ 搜索页面
- ✅ 购买支付流程
- ✅ 我的订单
- ✅ 图片加载优化（CachedNetworkImage）
- ✅ 社区模块（帖子浏览/发帖/回复/点赞）
- ✅ 消息通知
- ✅ 等级/成就展示

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

2. **配置 Nginx 反向代理**（网站 → tt.cy.mk → 配置文件）：

> ⚠️ **关键注意事项（踩坑记录）**
> - `location /uploads/` 和 `location /api/` **必须放在 `location ^~ /` 之前**，否则 `^~` 优先级最高会先匹配，导致 `/uploads/` 和 `/api/` 被错误转发到前端 (:3000)
> - 宝塔面板默认配置里有 `location ^~ /` 且 `proxy_pass` 指向前端，新增的 location block 必须排在它前面
> - 如果配置了 SSL (443)，需要确保这些 location block 在 443 的 server block 里也存在（或写在 80/443 共用的同一个 server block 里）
> - `proxy_set_header` 中的变量名必须正确：`X-Forwarded-For`（不是 `X-Forwarded-For`），`X-Real-IP`，`X-Forwarded-Proto`

```nginx
# ===== 必须放在 location ^~ / 之前 =====

# 上传文件代理（本地存储时必需）
location /uploads/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

# 后端 API 代理
location /api/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

# 管理后台静态资源
location /admin/assets/ {
    proxy_pass http://127.0.0.1:3001;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}

# 管理后台
location /admin/ {
    proxy_pass http://127.0.0.1:3001/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

# 消息/通知 API（如果独立部署）
location /messages/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}

# ===== 原有 location ^~ / 保持不变，放在后面 =====
location ^~ / {
    proxy_pass http://127.0.0.1:3000;
    # ... 原有配置 ...
}
```

**排查 `/uploads/` 返回 404 的步骤：**
```bash
# 1. 确认 location 是否加载
nginx -T 2>&1 | grep -A 5 "location /uploads"

# 2. 确认是 HTTP 还是 HTTPS 的 server block 生效
#    如果访问 https://domain/uploads/...，检查 443 server block 里是否有 /uploads/ location

# 3. 直接访问后端端口验证后端是否正常
curl http://127.0.0.1:8080/api/health

# 4. 检查 nginx error log
tail -50 /www/wwwlogs/tt.cy.mk.error.log
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

### Flutter App

```bash
cd photoset_mobile/
flutter pub get
flutter run
# 连接设备后自动部署
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
| Phase 7 | 社区功能（帖子/消息/等级/成就/积分商城） | ✅ |
| Flutter App | 移动端完整实现 | ✅ |

---

## 相关文档

- [后端 API 文档](./docker/README.DOCKER.md) — Docker 部署指南
- [前端 README](./frontend/README.md) — 用户端详细文档
- [管理后台 README](./frontend-admin/README.md) — 后台详细文档

---

*最后更新: 2026-05-31*
