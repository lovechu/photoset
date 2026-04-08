# PhotoSet Backend

基于 Go + Gin 的摄影套图平台后端服务。

## 项目结构

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # 程序入口
├── internal/
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库连接
│   ├── domain/                  # 领域模型
│   ├── http/
│   │   ├── handlers/            # HTTP 处理器
│   │   ├── middleware/          # 中间件
│   │   └── routes/              # 路由配置
│   ├── pkg/                     # 工具包
│   │   ├── jwt/                 # JWT 工具
│   │   ├── password/            # 密码加密
│   │   └── response/            # 统一响应格式
│   ├── repository/              # 数据仓储层
│   └── service/                 # 业务逻辑层
├── scripts/
│   └── init.sql                 # 数据库初始化脚本
├── .env.example                 # 环境变量示例
├── go.mod                       # Go 模块文件
└── README.md                    # 项目说明
```

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **配置**: godotenv
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

修改数据库连接信息和 JWT 密钥。

### 3. 初始化数据库

```bash
mysql -u root -p < scripts/init.sql
```

### 4. 启动服务

```bash
go run cmd/api/main.go
```

服务将在 `http://localhost:8080` 启动。

## 已实现接口

### 健康检查

```bash
GET /api/health
```

返回示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "status": "ok",
    "time": "2026-04-08 13:00:00"
  }
}
```

### 用户注册

```bash
POST /api/auth/register
Content-Type: application/json

{
  "nickname": "张三",
  "email": "zhangsan@example.com",
  "password": "password123"
}
```

返回示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user": {
      "id": 1,
      "nickname": "张三",
      "email": "zhangsan@example.com",
      "role": "user",
      "status": 1,
      "created_at": "2026-04-08T13:00:00Z",
      "updated_at": "2026-04-08T13:00:00Z"
    }
  }
}
```

### 用户登录

```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "zhangsan@example.com",
  "password": "password123"
}
```

返回示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nickname": "张三",
      "email": "zhangsan@example.com",
      "role": "user"
    }
  }
}
```

### 获取当前用户信息

```bash
GET /api/auth/me
Authorization: Bearer <token>
```

返回示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user": {
      "id": 1,
      "nickname": "张三",
      "email": "zhangsan@example.com",
      "role": "user",
      "status": 1,
      "created_at": "2026-04-08T13:00:00Z",
      "updated_at": "2026-04-08T13:00:00Z"
    }
  }
}
```

### 套图列表

```bash
GET /api/photosets?page=1&page_size=20&tag=风景
```

查询参数：
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 20，最大 100）
- `tag`: 标签名（可选，用于筛选）

返回示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "西湖美景",
        "cover": "https://example.com/cover.jpg",
        "description": "西湖风景摄影",
        "is_free": 1,
        "price": 0,
        "user_id": 1,
        "status": "published",
        "created_at": "2026-04-08T13:00:00Z",
        "updated_at": "2026-04-08T13:00:00Z",
        "user": {
          "id": 1,
          "nickname": "张三"
        },
        "tags": [
          {
            "id": 1,
            "name": "风景",
            "created_at": "2026-04-08T13:00:00Z"
          }
        ]
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 套图详情

```bash
GET /api/photosets/:id
Authorization: Bearer <token> (可选)
```

付费内容访问规则：
- **免费套图**: 任何人都可查看完整图片列表
- **付费套图**:
  - 未登录游客、普通用户：只返回基础信息和封面，不返回完整图片列表
  - 作者本人、管理员、会员：返回完整图片列表

返回示例（免费套图）：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "西湖美景",
    "cover": "https://example.com/cover.jpg",
    "description": "西湖风景摄影",
    "is_free": 1,
    "price": 0,
    "user_id": 1,
    "status": "published",
    "created_at": "2026-04-08T13:00:00Z",
    "updated_at": "2026-04-08T13:00:00Z",
    "user": {
      "id": 1,
      "nickname": "张三",
      "email": "zhangsan@example.com"
    },
    "tags": [
      {
        "id": 1,
        "name": "风景",
        "created_at": "2026-04-08T13:00:00Z"
      }
    ],
    "photos": [
      {
        "id": 1,
        "photoset_id": 1,
        "url": "https://example.com/photo1.jpg",
        "sort_order": 1,
        "created_at": "2026-04-08T13:00:00Z"
      }
    ]
  }
}
```

返回示例（付费套图，无权限用户）：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "西湖美景",
    "cover": "https://example.com/cover.jpg",
    "description": "西湖风景摄影",
    "is_free": 0,
    "price": 99.00,
    "user_id": 1,
    "status": "published",
    "created_at": "2026-04-08T13:00:00Z",
    "updated_at": "2026-04-08T13:00:00Z",
    "user": {
      "id": 1,
      "nickname": "张三",
      "email": "zhangsan@example.com"
    },
    "tags": [
      {
        "id": 1,
        "name": "风景",
        "created_at": "2026-04-08T13:00:00Z"
      }
    ],
    "photos": []
  }
}
```

### 创建套图

```bash
POST /api/photosets
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "西湖美景",
  "cover": "https://example.com/cover.jpg",
  "description": "西湖风景摄影",
  "is_free": 0,
  "price": 99.00,
  "tags": ["风景", "摄影"],
  "photos": [
    {
      "url": "https://example.com/photo1.jpg",
      "sort_order": 1
    },
    {
      "url": "https://example.com/photo2.jpg",
      "sort_order": 2
    }
  ],
  "status": "published"
}
```

请求参数说明：
- `title`: 标题（必填，最大 200 字符）
- `cover`: 封面 URL（必填，最大 500 字符）
- `description`: 描述（可选）
- `is_free`: 是否免费（必填，1-免费，0-付费）
- `price`: 价格（当 is_free=1 时应为 0）
- `tags`: 标签数组（可选，标签不存在时自动创建）
- `photos`: 图片数组（可选，每个图片包含 url 和 sort_order）
- `status`: 状态（必填，可选值：draft/published/pending）

**权限要求**：
- 必须登录
- 只有 `creator` 和 `admin` 角色可以创建套图
- 普通用户访问返回 403 权限不足

返回示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "西湖美景",
    "cover": "https://example.com/cover.jpg",
    "description": "西湖风景摄影",
    "is_free": 0,
    "price": 99.00,
    "user_id": 1,
    "status": "published",
    "created_at": "2026-04-08T13:00:00Z",
    "updated_at": "2026-04-08T13:00:00Z"
  }
}
```

### 标签列表

```bash
GET /api/tags
```

返回示例：

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "风景",
      "created_at": "2026-04-08T13:00:00Z"
    },
    {
      "id": 2,
      "name": "摄影",
      "created_at": "2026-04-08T13:00:00Z"
    }
  ]
}
```

## 用户角色

系统支持五类用户角色：

- **guest**: 游客
- **user**: 普通用户
- **member**: 会员
- **creator**: 上传者
- **admin**: 管理员

## API 响应格式

所有接口统一返回格式：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误响应：

```json
{
  "code": -1,
  "message": "error message",
  "data": null
}
```

## 环境变量说明

参考 `backend/.env.example`：

- 服务器配置（端口、模式）
- 数据库配置（MySQL）
- Redis 配置
- JWT 配置（密钥、过期时间）
- 文件存储配置（本地或海外对象存储）

## 用户模型字段

- `id`: 用户 ID
- `nickname`: 昵称
- `email`: 邮箱（唯一）
- `password_hash`: 密码哈希（bcrypt 加密）
- `role`: 角色（guest/user/member/creator/admin）
- `status`: 状态（1-正常，0-禁用）
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `last_login_at`: 最后登录时间

## 数据表结构

### photosets（套图表）
- `id`: 套图 ID（主键）
- `title`: 标题
- `cover`: 封面 URL
- `description`: 描述
- `is_free`: 是否免费（1-免费，0-付费）
- `price`: 价格
- `user_id`: 作者 ID（外键）
- `status`: 状态（draft/published/pending）
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `deleted_at`: 删除时间（软删除）

### photos（图片表）
- `id`: 图片 ID（主键）
- `photoset_id`: 套图 ID（外键）
- `url`: 图片 URL
- `sort_order`: 排序顺序
- `created_at`: 创建时间
- `deleted_at`: 删除时间（软删除）

### tags（标签表）
- `id`: 标签 ID（主键）
- `name`: 标签名（唯一）
- `created_at`: 创建时间

### photoset_tags（套图标签关联表）
- `photoset_id`: 套图 ID（外键）
- `tag_id`: 标签 ID（外键）
- 主键：`(photoset_id, tag_id)`

## 已实现功能

### 第一阶段：后端基础骨架 + 认证基础 (2026-04-08)
- ✅ Go + Gin 项目结构
- ✅ 配置管理（环境变量）
- ✅ MySQL + Redis 连接初始化
- ✅ 统一 JSON 响应格式
- ✅ 基础中间件（CORS、日志、异常恢复）
- ✅ JWT Token 生成、解析、校验
- ✅ 密码 bcrypt 加密
- ✅ 用户注册、登录、获取当前用户接口
- ✅ 鉴权中间件
- ✅ 用户模型与仓储、服务层实现
- ✅ 数据库初始化脚本

### 第二阶段：套图模块 MVP (2026-04-08)
- ✅ 套图列表接口（分页、标签筛选）
- ✅ 套图详情接口（付费内容裁剪）
- ✅ 创建套图接口（creator/admin 权限）
- ✅ 标签列表接口
- ✅ 角色权限中间件
- ✅ 套图领域模型、仓储层、服务层
- ✅ 图片和标签关联

## 本轮未实现的功能

本轮为套图模块最小可用版本（MVP），以下功能暂未实现：
- ❌ 图片上传接口（创建接口只接收已存在的图片 URL）
- ❌ 套图更新、删除接口
- ❌ 收藏功能
- ❌ 订单与支付
- ❌ 会员订阅逻辑
- ❌ 后台审核流
- ❌ 完整搜索系统（关键词搜索）
- ❌ 管理后台接口
- ❌ 文件存储（本地存储 / S3 / R2 / B2）

## 下一轮建议

建议优先实现：
1. 图片上传模块（本地存储或海外对象存储）
2. 套图更新、删除接口
3. 会员订阅与订单支付基础流程

