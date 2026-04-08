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

## 后续开发计划

- [ ] 套图模块（CRUD）
- [ ] 图片上传模块（本地存储 + 海外对象存储）
- [ ] 订单支付模块
- [ ] 管理后台接口
- [ ] 权限中间件（角色权限控制）

