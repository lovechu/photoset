# PhotoSet 后端第三轮开发指令

## 本轮覆盖范围

| 任务 | 内容 | 新增文件 | 修改文件 |
|------|------|----------|----------|
| 任务 8a | 会员套餐模型 + API | 2 新增 | 2 修改 |
| 任务 8b | 订单模型 + API（含模拟支付） | 3 新增 | 2 修改 |
| 任务 9 | 管理后台 API | 3 新增 | 1 修改 |

**合计：8 新增文件 + 3 修改文件**（3 个修改文件存在重叠，实际修改 `main.go`、`routes.go`、`user.go`）

---

## 一、任务 8a：会员套餐模型 + API

### 1.1 新增领域模型

**文件：`internal/domain/membership.go`**

```go
package domain

import "time"

// Membership 会员套餐模型
type Membership struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Name        string    `gorm:"size:100;not null" json:"name"`              // 套餐名称，如 "月度会员"、"年度会员"
    Duration    int       `gorm:"not null;comment:会员有效天数" json:"duration"` // 30 / 90 / 365
    Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`    // 套餐价格
    Description string    `gorm:"type:text" json:"description"`                 // 套餐描述
    Status      int8      `gorm:"default:1;comment:1-上架,0-下架" json:"status"`
    SortOrder   int       `gorm:"default:0" json:"sort_order"`                 // 排序权重，越大越靠前
}

func (Membership) TableName() string {
    return "memberships"
}
```

### 1.2 新增仓储

**文件：`internal/repository/membership_repository.go`**

```go
package repository

import (
    "photoset/internal/domain"
    "gorm.io/gorm"
)

type MembershipRepository interface {
    List() ([]domain.Membership, error)   // 返回所有上架套餐，按 sort_order DESC
}

type membershipRepository struct {
    db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
    return &membershipRepository{db: db}
}

func (r *membershipRepository) List() ([]domain.Membership, error) {
    var list []domain.Membership
    err := r.db.Where("status = 1").Order("sort_order DESC, id ASC").Find(&list).Error
    return list, err
}
```

### 1.3 新增处理器

**文件：`internal/http/handlers/membership_handler.go`**

```
MembershipHandler
  └─ List(c *gin.Context)     // GET /api/memberships
```

**接口行为：**
- 查询所有 `status=1` 的套餐，按 `sort_order DESC, id ASC` 排序
- 返回数组，无需分页（套餐数量极少）
- 无需鉴权，任何人可查看

**返回格式：**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "年度会员",
        "duration": 365,
        "price": 199.00,
        "description": "全年无限查看所有付费套图",
        "sort_order": 100
      }
    ]
  }
}
```

---

## 二、任务 8b：订单模型 + API（含模拟支付）

### 2.1 新增领域模型

**文件：`internal/domain/order.go`**

```go
package domain

import "time"

// Order 订单模型
type Order struct {
    ID              uint       `gorm:"primaryKey" json:"id"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
    OrderNo         string     `gorm:"size:64;uniqueIndex;not null" json:"order_no"`         // 订单号
    UserID          uint       `gorm:"not null;index" json:"user_id"`                        // 下单用户
    MembershipID    uint       `gorm:"index" json:"membership_id"`                           // 套餐 ID（会员购买时）
    MembershipName  string     `gorm:"size:100" json:"membership_name"`                      // 套餐名称（冗余快照）
    PhotoSetID      uint       `gorm:"index" json:"photoset_id"`                             // 套图 ID（单套图购买时，可为 0）
    PhotoSetTitle   string     `gorm:"size:200" json:"photoset_title"`                       // 套图标题（冗余快照）
    Amount          float64    `gorm:"type:decimal(10,2);not null" json:"amount"`            // 支付金额
    Status          string     `gorm:"size:20;default:pending;not null" json:"status"`       // pending / paid / failed / cancelled
    PaidAt          *time.Time `json:"paid_at"`                                              // 支付时间
    TransactionNo   string     `gorm:"size:128" json:"transaction_no"`                       // 第三方交易号（mock 时自动生成）
    ExpiredAt       *time.Time `json:"expired_at"`                                           // 订单过期时间
}

func (Order) TableName() string {
    return "orders"
}
```

**订单号生成规则：** `PS` + 时间戳(yyyyMMddHHmmss) + 6位随机数字 = 共20位，如 `PS20260408210000123456`

**订单状态流转：**
```
pending → paid       (模拟支付成功)
pending → failed     (模拟支付失败)
pending → cancelled  (用户取消)
```

### 2.2 新增仓储

**文件：`internal/repository/order_repository.go`**

```go
package repository

interface OrderRepository {
    Create(order *domain.Order) error
    FindByID(id uint) (*domain.Order, error)
    FindByOrderNo(orderNo string) (*domain.Order, error)
    Update(order *domain.Order) error
    ListByUserID(userID uint, page, pageSize int) ([]domain.Order, int64, error)  // 分页
}
```

### 2.3 新增服务

**文件：`internal/service/order_service.go`**

```go
package service

type OrderService struct {
    orderRepo      repository.OrderRepository
    membershipRepo repository.MembershipRepository
    userRepo       repository.UserRepository
}

// 核心方法：
// CreateMembershipOrder(userID uint, membershipID uint) (*domain.Order, error)
//   - 校验套餐存在且上架
//   - 生成订单号
//   - 冗余存储套餐名称
//   - 设置过期时间（30 分钟后）
//   - 写入 orders 表，status=pending

// MockPay(orderNo string) (*domain.Order, string, error)
//   - 根据 orderNo 查找订单
//   - 校验状态必须为 pending，且未过期
//   - 更新 status=paid, paid_at=now, transaction_no=mock_xxx
//   - 如果是会员订单（membership_id > 0）：
//       更新 users 表：role='member', 同时在 user 上设置 member_expires_at = now + duration天
//   - 调用 jwt.GenerateToken(userID, 新role) 生成新 Token
//   - 返回 (order, newToken, nil)

// CancelOrder(userID uint, orderNo string) error
//   - 只有订单所有者可以取消
//   - 只有 pending 状态可以取消
//   - 更新 status=cancelled

// GetOrderList(userID uint, page, pageSize int) ([]domain.Order, int64, error)
//   - 查询指定用户的订单列表，按 created_at DESC
```

### 2.4 新增处理器

**文件：`internal/http/handlers/order_handler.go`**

```
OrderHandler
  ├─ CreateMembershipOrder(c *gin.Context)   // POST /api/orders/membership
  ├─ MockPay(c *gin.Context)                 // POST /api/orders/:orderNo/pay
  ├─ CancelOrder(c *gin.Context)             // POST /api/orders/:orderNo/cancel
  └─ List(c *gin.Context)                    // GET  /api/orders
```

**接口详细行为：**

#### POST /api/orders/membership — 创建会员订单
- 鉴权：`Auth()` + `RequireRoles("user", "member")`（admin/creator 不需要买）
  - 实际上不限制也行，让所有人都能创建，简单起见只要求 `Auth()`
- 请求体：
```json
{
  "membership_id": 1
}
```
- 响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "order": {
      "id": 1,
      "order_no": "PS20260408210000123456",
      "amount": 199.00,
      "status": "pending",
      "expired_at": "2026-04-08T21:30:00Z"
    }
  }
}
```

#### POST /api/orders/:orderNo/pay — 模拟支付
- 鉴权：`Auth()`
- 只有订单所有者可以支付
- 成功响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "order": {
      "id": 1,
      "order_no": "PS20260408210000123456",
      "status": "paid",
      "paid_at": "2026-04-08T21:05:00Z",
      "transaction_no": "mock_20260408210500001"
    },
    "token": "新的JWT Token（角色已更新为member）"
  }
}
```
- **关键行为**：支付成功后自动将用户的 `role` 更新为 `member`，返回新 Token
- 新 Token 中 `role` 字段为 `"member"`，前端需要用新 Token 替换旧 Token
- 如果购买的是会员套餐，同时更新 `member_expires_at` 字段（见下方 user.go 修改）
- 如果用户已经是 admin/creator，不降级其角色，仅记录订单

#### POST /api/orders/:orderNo/cancel — 取消订单
- 鉴权：`Auth()`
- 只有订单所有者可以取消
- 只有 `pending` 状态可以取消

#### GET /api/orders — 订单列表
- 鉴权：`Auth()`
- 查询当前用户的订单，分页
- 请求参数：`page`, `page_size`
- 响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [...],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 三、任务 9：管理后台 API

### 3.1 新增处理器

**文件：`internal/http/handlers/admin_handler.go`**

```
AdminHandler
  ├─ GetUsers(c *gin.Context)              // GET    /api/admin/users        用户列表
  ├─ UpdateUserStatus(c *gin.Context)      // PUT    /api/admin/users/:id/status  启用/禁用用户
  ├─ GetPhotosets(c *gin.Context)          // GET    /api/admin/photosets    套图列表（含所有状态）
  ├─ ReviewPhotoset(c *gin.Context)        // POST   /api/admin/photosets/:id/review  审核套图
  ├─ GetOrders(c *gin.Context)             // GET    /api/admin/orders       订单列表
  └─ GetDashboard(c *gin.Context)          // GET    /api/admin/dashboard    数据统计
```

#### GET /api/admin/users — 用户管理列表
- 鉴权：`Auth()` + `RequireRoles("admin")`
- 请求参数：`page`, `page_size`, `role`（可选，按角色筛选）, `keyword`（可选，搜索昵称/邮箱）
- 响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "nickname": "admin",
        "email": "admin@photoset.com",
        "role": "admin",
        "status": 1,
        "created_at": "...",
        "last_login_at": "..."
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

#### PUT /api/admin/users/:id/status — 启用/禁用用户
- 鉴权：`Auth()` + `RequireRoles("admin")`
- 请求体：
```json
{
  "status": 0
}
```
- 不能禁用自己
- 不能操作其他 admin

#### GET /api/admin/photosets — 套图管理列表
- 鉴权：`Auth()` + `RequireRoles("admin")`
- 请求参数：`page`, `page_size`, `status`（可选，筛选 draft/published/pending）
- 返回所有状态的套图，包含作者信息
- 响应格式同用户列表

#### POST /api/admin/photosets/:id/review — 审核套图
- 鉴权：`Auth()` + `RequireRoles("admin")`
- 请求体：
```json
{
  "action": "approve"
}
```
- `action` 可选值：`approve`（通过，状态改为 published）、`reject`（拒绝，状态改为 draft）
- 只能审核 `pending` 状态的套图

#### GET /api/admin/orders — 订单管理列表
- 鉴权：`Auth()` + `RequireRoles("admin")`
- 请求参数：`page`, `page_size`, `status`（可选）
- 返回所有用户的订单，包含用户信息
- 响应格式同用户列表

#### GET /api/admin/dashboard — 数据统计
- 鉴权：`Auth()` + `RequireRoles("admin")`
- 无请求参数
- 响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "total_users": 1000,
    "total_members": 200,
    "total_creators": 50,
    "total_photosets": 500,
    "total_published_photosets": 400,
    "total_pending_photosets": 10,
    "total_orders": 800,
    "total_revenue": 50000.00,
    "today_orders": 10,
    "today_revenue": 500.00
  }
}
```

### 3.2 新增仓储

**文件：`internal/repository/admin_repository.go`**

```go
package repository

// AdminRepository 管理后台数据查询仓储
// 直接使用 *gorm.DB 进行复杂查询
// 方法：
//   - CountUsers(role string) (int64, error)
//   - CountPhotoSets(status string) (int64, error)
//   - CountOrders(status string) (int64, error)
//   - SumRevenue(status string) (float64, error)
//   - SumTodayRevenue() (float64, error)
//   - CountTodayOrders() (int64, error)
//   - ListUsers(page, pageSize int, role, keyword string) ([]domain.User, int64, error)
//   - ListAllPhotoSets(page, pageSize int, status string) ([]domain.PhotoSet, int64, error)
//   - ListAllOrders(page, pageSize int, status string) ([]domain.Order, int64, error)
```

### 3.3 新增服务

**文件：`internal/service/admin_service.go`**

```go
package service

// AdminService 管理后台服务
// 方法：
//   - GetDashboard() (map[string]interface{}, error)
//   - GetUserList(page, pageSize int, role, keyword string) ([]domain.User, int64, error)
//   - UpdateUserStatus(adminID, targetID uint, status int) error
//   - GetPhotoSetList(page, pageSize int, status string) ([]domain.PhotoSet, int64, error)
//   - ReviewPhotoset(photosetID uint, action string) error
//   - GetOrderList(page, pageSize int, status string) ([]domain.Order, int64, error)
```

---

## 四、需要修改的现有文件

### 4.1 `internal/domain/user.go` — 新增会员过期时间字段

```go
type User struct {
    // ... 现有字段保持不变 ...
    MemberExpiresAt *time.Time `gorm:"index" json:"member_expires_at"` // 会员过期时间，nil 表示非会员或永不过期
}
```

**注意：** 只新增这一个字段，不要改动其他任何字段。

### 4.2 `cmd/api/main.go` — AutoMigrate 新增表

在 `database.GetMySQL().AutoMigrate(...)` 中追加：
```go
&domain.Membership{},
&domain.Order{},
```

### 4.3 `internal/http/routes/routes.go` — 注册新路由

追加以下路由组：

```go
// 会员套餐路由
api.GET("/memberships", membershipHandler.List)

// 订单路由
orderGroup := api.Group("/orders")
{
    orderGroup.Use(middleware.Auth())
    orderGroup.POST("/membership", orderHandler.CreateMembershipOrder)
    orderGroup.GET("", orderHandler.List)
    orderGroup.POST("/:orderNo/pay", orderHandler.MockPay)
    orderGroup.POST("/:orderNo/cancel", orderHandler.CancelOrder)
}

// 管理后台路由
adminGroup := api.Group("/admin")
{
    adminGroup.Use(middleware.Auth(), middleware.RequireRoles("admin"))
    adminGroup.GET("/dashboard", adminHandler.GetDashboard)
    adminGroup.GET("/users", adminHandler.GetUsers)
    adminGroup.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
    adminGroup.GET("/photosets", adminHandler.GetPhotosets)
    adminGroup.POST("/photosets/:id/review", adminHandler.ReviewPhotoset)
    adminGroup.GET("/orders", adminHandler.GetOrders)
}
```

在 `Setup()` 函数顶部初始化新的服务和处理器：
```go
// 会员套餐
membershipRepo := repository.NewMembershipRepository(database.GetMySQL())

// 订单
orderRepo := repository.NewOrderRepository(database.GetMySQL())
orderService := service.NewOrderService(orderRepo, membershipRepo, userRepo)

// 管理后台
adminRepo := repository.NewAdminRepository(database.GetMySQL())
adminService := service.NewAdminService(adminRepo)
```

---

## 五、文件清单

### 新增 8 个文件
| # | 文件路径 | 说明 |
|---|----------|------|
| 1 | `internal/domain/membership.go` | 会员套餐领域模型 |
| 2 | `internal/domain/order.go` | 订单领域模型 |
| 3 | `internal/repository/membership_repository.go` | 会员套餐仓储 |
| 4 | `internal/repository/order_repository.go` | 订单仓储 |
| 5 | `internal/repository/admin_repository.go` | 管理后台数据查询仓储 |
| 6 | `internal/service/order_service.go` | 订单服务（含模拟支付逻辑） |
| 7 | `internal/service/admin_service.go` | 管理后台服务 |
| 8 | `internal/http/handlers/membership_handler.go` | 会员套餐处理器 |
| 9 | `internal/http/handlers/order_handler.go` | 订单处理器 |
| 10 | `internal/http/handlers/admin_handler.go` | 管理后台处理器 |

### 修改 3 个文件
| # | 文件路径 | 改动内容 |
|---|----------|----------|
| 1 | `internal/domain/user.go` | 新增 `MemberExpiresAt` 字段 |
| 2 | `cmd/api/main.go` | AutoMigrate 追加 `Membership{}`、`Order{}` |
| 3 | `internal/http/routes/routes.go` | 注册套餐/订单/管理后台路由及初始化 |

---

## 六、关键约束

### 6.1 模拟支付流程（核心）

```
用户点击购买 → 创建订单(pending) → 前端跳转支付页
→ 用户点击"模拟支付"按钮 → POST /api/orders/:orderNo/pay
→ 后端校验(pending + 未过期) → 更新订单为paid
→ 如果是会员订单：更新用户role为member + 设置member_expires_at
→ 生成新JWT Token(包含新role) → 返回给前端
→ 前端用新Token替换旧Token → 后续请求自动携带member身份
```

**注意事项：**
- `transaction_no` 使用 `mock_` 前缀 + 时间戳，如 `mock_1712581500001`
- 支付成功后，`users.role` 更新为 `member`，同时设置 `member_expires_at = now + duration天`
- 如果用户已经是 `admin` 或 `creator` 角色，**不要降级**其角色，仅记录订单为已支付
- 新 Token 通过响应体中的 `token` 字段返回，前端需要自行替换

### 6.2 角色上下文键名

- 全项目统一使用 `middleware.UserKey = "user_id"` 和 `middleware.RoleKey = "user_role"`
- 不要出现 `"role"` 或其他自定义键名

### 6.3 API 响应格式

统一使用现有 `response` 包：
```go
response.Success(c, data)          // 成功
response.Error(c, code, msg)       // 业务错误
response.BadRequest(c, msg)        // 参数错误
response.Unauthorized(c, msg)      // 未授权
response.Forbidden(c, msg)         // 禁止访问
response.ServerError(c, msg)       // 服务器错误
```

### 6.4 编码规范

- 仓储层只做数据访问，不含业务逻辑
- 服务层封装业务逻辑，通过仓储层访问数据
- 处理器层做参数校验、调用服务、格式化响应
- 使用 `middleware.GetUserID(c)` 和 `middleware.GetUserRole(c)` 获取上下文中的用户信息

### 6.5 不需要实现的功能

- ❌ 真实支付网关（只用 mock）
- ❌ 单套图购买（当前只做会员套餐购买）
- ❌ 退款功能
- ❌ 会员到期自动降级（后续实现）
- ❌ 管理后台的前端页面
- ❌ 权限更细粒度控制（如 creator 只能管理自己的套图）

---

## 七、数据统计 Dashboard 说明

`GetDashboard` 需要的统计指标及 SQL 思路：

| 指标 | 说明 |
|------|------|
| `total_users` | `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL` |
| `total_members` | `SELECT COUNT(*) FROM users WHERE role='member' AND deleted_at IS NULL` |
| `total_creators` | `SELECT COUNT(*) FROM users WHERE role='creator' AND deleted_at IS NULL` |
| `total_photosets` | `SELECT COUNT(*) FROM photosets WHERE deleted_at IS NULL` |
| `total_published_photosets` | `... WHERE status='published' AND deleted_at IS NULL` |
| `total_pending_photosets` | `... WHERE status='pending' AND deleted_at IS NULL` |
| `total_orders` | `SELECT COUNT(*) FROM orders` |
| `total_revenue` | `SELECT COALESCE(SUM(amount), 0) FROM orders WHERE status='paid'` |
| `today_orders` | `... WHERE status='paid' AND paid_at >= 今天0点` |
| `today_revenue` | `SELECT COALESCE(SUM(amount), 0) FROM orders WHERE status='paid' AND paid_at >= 今天0点` |

---

## 八、自查清单

完成后请自查以下场景：

1. **GET /api/memberships** — 无需 token，返回上架套餐列表
2. **POST /api/orders/membership** — 普通 user 带有效 token 创建订单成功
3. **POST /api/orders/:orderNo/pay** — 模拟支付成功，返回新 Token，用户角色变为 member
4. **POST /api/orders/:orderNo/pay** — 重复支付已支付订单返回错误
5. **POST /api/orders/:orderNo/pay** — 支付过期订单返回错误
6. **GET /api/orders** — 查看自己的订单列表
7. **GET /api/admin/dashboard** — admin 查看统计数据
8. **GET /api/admin/dashboard** — 非 admin 访问返回 403
9. **POST /api/admin/photosets/:id/review** — admin 审核通过套图
10. **PUT /api/admin/users/:id/status** — admin 禁用用户（不能禁用自己和其他 admin）
11. **admin 购买会员后角色不变** — admin/creator 支付后不会被降级为 member
