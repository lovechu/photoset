# PhotoSet 模块修复自查清单

## 修复内容

### 1. 统一角色上下文键名
- ✅ `auth.go` 使用 `RoleKey = "user_role"` 写入上下文
- ✅ `role.go` 改为从 `RoleKey` 读取角色
- ✅ `photoset.go` Detail 和 Create 统一使用 `"user_role"` 和 `"user_id"`

### 2. 详情接口支持可选鉴权
- ✅ 路由配置: `photosets.GET("/:id", photosetHandler.Detail)` 无强制鉴权
- ✅ Handler 逻辑:
  - 不带 token: `userRole=""`, `isLoggedIn=false` → 按游客处理
  - 带 token: `userRole` 和 `userID` 正常读取 → 按角色权限处理

## 自查验证场景

### 场景 1: creator 创建套图成功
**前提**:
- 用户注册并登录,角色为 creator
- JWT token 包含 `user_role=creator` 和 `user_id`

**请求**:
```bash
POST /api/photosets
Authorization: Bearer <creator_token>
{
  "title": "测试套图",
  "cover": "https://example.com/cover.jpg",
  "is_free": 1,
  "status": "published",
  "photos": [
    {"url": "https://example.com/photo1.jpg", "sort_order": 1}
  ]
}
```

**预期结果**:
- `RequireRoles("creator", "admin")` 从 `RoleKey="user_role"` 读取到 `creator`
- 角色匹配,通过中间件
- Handler 从 `"user_id"` 获取 userID
- 创建成功,返回套图信息

### 场景 2: 普通用户创建套图被拒绝
**前提**:
- 用户注册并登录,角色为 user
- JWT token 包含 `user_role=user` 和 `user_id`

**请求**:
```bash
POST /api/photosets
Authorization: Bearer <user_token>
{
  "title": "测试套图",
  "cover": "https://example.com/cover.jpg",
  "is_free": 1,
  "status": "published"
}
```

**预期结果**:
- `RequireRoles("creator", "admin")` 从 `RoleKey="user_role"` 读取到 `user`
- 角色不在允许列表中
- 返回 403 权限不足

### 场景 3: 游客访问付费套图详情,看不到完整图片列表
**前提**:
- 数据库存在付费套图 (is_free=0)
- 游客未登录,无 token

**请求**:
```bash
GET /api/photosets/1
```

**预期结果**:
- 无 token,`c.Get("user_role")` 返回 false,`userRole=""`, `isLoggedIn=false`
- 调用 `CanViewFullPhotos(photoset, "", 0, false)`
- 逻辑判断: 不是作者、不是 admin、不是 member、未登录
- `canViewFull = false`
- 返回套图基础信息,`photos.Photos = []` (空数组)

### 场景 4: admin 访问付费套图详情,能拿到完整图片列表
**前提**:
- 数据库存在付费套图 (is_free=0)
- admin 用户登录,角色为 admin

**请求**:
```bash
GET /api/photosets/1
Authorization: Bearer <admin_token>
```

**预期结果**:
- 有 token,`c.Get("user_role")` 返回 `admin`, `userRole="admin"`, `isLoggedIn=true`
- 调用 `CanViewFullPhotos(photoset, "admin", <admin_id>, true)`
- 逻辑判断: 角色 `admin` 在允许列表中
- `canViewFull = true`
- 调用 `GetPhotoSetDetail(id)` 获取完整信息
- 返回套图完整信息,包含 `photos` 数组

### 场景 5: 作者本人访问付费套图详情,能拿到完整图片列表
**前提**:
- 数据库存在付费套图 (is_free=0),user_id=123
- 用户 123 (creator) 登录,是该套图的作者

**请求**:
```bash
GET /api/photosets/1
Authorization: Bearer <author_token>
```

**预期结果**:
- 有 token,`userRole="creator"`, `userID=123`, `isLoggedIn=true`
- 调用 `CanViewFullPhotos(photoset, "creator", 123, true)`
- 逻辑判断: `userID == photoset.UserID` (123 == 123)
- `canViewFull = true`
- 返回套图完整信息,包含 `photos` 数组

### 场景 6: 会员访问付费套图详情,能拿到完整图片列表
**前提**:
- 数据库存在付费套图 (is_free=0),user_id=456
- 会员用户 (member) 登录,不是该套图的作者

**请求**:
```bash
GET /api/photosets/1
Authorization: Bearer <member_token>
```

**预期结果**:
- 有 token,`userRole="member"`, `userID=789`, `isLoggedIn=true`
- 调用 `CanViewFullPhotos(photoset, "member", 789, true)`
- 逻辑判断: 不是作者,但角色 `member` 在允许列表中
- `canViewFull = true`
- 返回套图完整信息,包含 `photos` 数组

## 代码检查要点

### 键名统一性
```go
// auth.go
const RoleKey = "user_role"
c.Set(RoleKey, claims.Role)  // 写入 "user_role"

// role.go
role, exists := c.Get(RoleKey)  // 读取 "user_role"

// photoset.go
if role, exists := c.Get("user_role"); exists {  // 读取 "user_role"
    userRole = role.(string)
}
if uid, exists := c.Get("user_id"); exists {  // 读取 "user_id"
    userID = uid.(uint)
}
```

### 可选鉴权逻辑
```go
// photoset.go - Detail 方法
var userRole string
var userID uint
var isLoggedIn bool

if role, exists := c.Get("user_role"); exists {  // 可选读取
    userRole = role.(string)
    isLoggedIn = true
}
if uid, exists := c.Get("user_id"); exists {  // 可选读取
    userID = uid.(uint)
}

// isLoggedIn=false 时,userRole 和 userID 为空值,按游客处理
```

## 修复验证结果

- ✅ 键名已统一: 全项目使用 `"user_role"` 和 `"user_id"`
- ✅ RequireRoles 中间件正确读取统一键名
- ✅ 详情接口支持可选鉴权: 不带 token 游客访问,带 token 按角色权限
- ✅ 无 linter 错误
- ✅ 路由配置正确: Detail 接口无强制鉴权中间件
