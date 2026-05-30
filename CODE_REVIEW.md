# PhotoSet Backend 代码审查报告

> 项目：PhotoSet 照片套图付费平台后端
> 技术栈：Go 1.22 + Gin + GORM + MySQL + Redis
> 审查范围：`internal/` 全部 Go 源码（70 个文件）
> 审查日期：2026-04-20

---

## 严重问题（🔴 必须立即修复）

### #1 RestartBackend 命令执行接口 — 安全风险

- **文件**：`internal/http/handlers/admin_handler.go`
- **位置**：`RestartBackend` 方法
- **问题**：通过 HTTP 接口直接调用 `exec.Command("docker", "restart", "photoset-backend")`，暴露了进程执行能力。虽然当前参数硬编码，但这个接口本身就危险。如果未来改为动态参数，即构成 RCE。
- **修复方案**：
  - 方案 A（推荐）：删除此接口，通过运维工具管理进程重启
  - 方案 B：保留但限制只允许内网 IP 访问，增加额外的管理员密码/签名验证

---

### #2 config.Load() 重复调用 — 配置不一致

- **文件**：`internal/http/routes/routes.go` 第 14 行、`internal/http/middleware/signverify.go` 第 38 行
- **问题**：`routes.Setup()` 和 `signverify.go` 中各自调用了 `config.Load()`，而 `main.go` 已经加载过一次。`Load()` 每次返回新实例并重新读取 `.env`。如果部署期间 `.env` 被修改，中间件用的 `SignSecret` 和 main 里的不一致，会导致签名验证失败。
- **修复方案**：
  - 将 `cfg` 作为参数从 `main.go` 传递到 `routes.Setup(r, cfg)`、`middleware.SignVerify(cfg)` 等
  - 或使用单例模式（`sync.Once`）确保 config 只加载一次

---

### #3 CORS 配置 Allow-Origin: * 与 Allow-Credentials: true 冲突

- **文件**：`internal/http/middleware/cors.go`
- **问题**：`Access-Control-Allow-Origin: *` 和 `Access-Control-Allow-Credentials: true` 不能同时设置。浏览器规范禁止通配符 origin 搭配 credentials，部分浏览器会直接拒绝跨域请求。
- **修复方案**：
  - 如果前端需要带 cookie/token：将 `*` 改为具体域名白名单，从请求的 `Origin` 头匹配
  - 如果不需要 credentials：删除 `Allow-Credentials: true`
  - 示例：
    ```go
    origin := c.GetHeader("Origin")
    allowedOrigins := []string{"https://yourdomain.com", "http://localhost:3000"}
    for _, o := range allowedOrigins {
        if origin == o {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            break
        }
    }
    ```

---

### #4 JWT / 签名密钥使用不安全默认值

- **文件**：`internal/config/config.go`
- **位置**：`JWT_SECRET` 默认 `"default-secret-key"`，`SIGN_SECRET` 默认 `"default-sign-secret-change-me"`
- **问题**：部署时如果忘记配置环境变量，密钥等于公开，JWT 可被伪造，签名 URL 可被绕过。
- **修复方案**：启动时检测到默认值直接 panic：
  ```go
  if cfg.JWT.Secret == "default-secret-key" {
      log.Fatal("FATAL: JWT_SECRET is not configured. Set a strong random secret in .env")
  }
  if cfg.Storage.SignSecret == "default-sign-secret-change-me" {
      log.Fatal("FATAL: SIGN_SECRET is not configured. Set a strong random secret in .env")
  }
  ```

---

## 中等问题（🟡 应尽快修复）

### #5 StatsTrend N+1 查询

- **文件**：`internal/http/handlers/admin_handler.go` — `StatsTrend` 方法
- **问题**：循环 N 天，每天执行 4 条独立 SQL（新用户、新订单、营收、新套图），7 天 = 28 次 SQL，30 天 = 120 次。
- **修复方案**：用一条聚合 SQL 按 `DATE(created_at)` 分组：
  ```sql
  SELECT
    DATE(created_at) as date,
    SUM(CASE WHEN type='user' THEN 1 ELSE 0 END) as new_users,
    ...
  FROM (
    SELECT created_at, 'user' as type FROM users WHERE created_at >= ?
    UNION ALL
    SELECT created_at, 'order' as type FROM orders WHERE created_at >= ?
    UNION ALL
    SELECT created_at, 'photoset' as type FROM photosets WHERE created_at >= ?
  ) t
  GROUP BY DATE(created_at)
  ```

---

### #6 photoset_repository.go List() 方法 — GORM 语法错误

- **文件**：`internal/repository/photoset_repository.go` — `List` 方法，约第 75 行
- **问题**：`.Where(query)` 传入了一个 `*gorm.DB` 对象。GORM 的 `Where()` 接受 string / struct / map，不接受 `*gorm.DB`。这会导致运行时 panic 或错误结果。
- **修复方案**：直接在 `query` 上链式调用，不要新建查询：
  ```go
  // 错误写法：
  err := r.db.Table("photosets").Select(...).Where(query).Scan(&photosets).Error

  // 正确写法：直接用 query 继续链式
  err := query.Select(`photosets.*, (SELECT COUNT(*) ...) AS photo_count, ...`).
      Preload("User").Preload("Tags").
      Offset(offset).Limit(pageSize).
      Order("created_at DESC").
      Scan(&photosets).Error
  ```
  注意：如果 `query` 已经是 `r.db.Model(&domain.PhotoSet{})`，需要改成 `r.db.Table("photosets")` 并重建条件。

---

### #7 paidStatus() 在 middleware 中重复实现

- **文件**：`internal/http/middleware/signverify.go` — `paidStatus()` 函数
- **问题**：和 `internal/service/paid_status_cache.go` 的 `IsPaid()` 功能完全重复。middleware 层不应直接操作数据库和 Redis，绕过了 service 层的统一缓存管理。
- **修复方案**：
  - 删除 `signverify.go` 中的 `paidStatus()` 函数
  - 改为调用 `service.IsPaid(photosetID)`（需要将 service 注入到 middleware，或通过依赖传递）

---

### #8 Handler 层直接操作 database.GetMySQL() — 破坏分层

- **文件**：`internal/http/handlers/admin_handler.go`
- **涉及方法**：`GetPhotoSetsByStatus`、`ApprovePhotoSet`、`RejectPhotoSet`、`GetUsers`、`GetUserDetail`、`BanUser`、`UpdateUserRole`、`Stats`、`StatsTrend`、`GetOrders`、`AdminRefund`
- **问题**：Handler 层直接调用 `database.GetMySQL()` 做 CRUD，绕过了 repository 和 service 层。这使得业务逻辑分散在 handler 中，无法复用和单独测试。
- **修复方案**：将所有数据库操作下沉到 service/repository 层，handler 只做请求解析和响应构造。

---

### #9 订单号生成使用 math/rand — 不安全

- **文件**：`internal/service/order_service.go` — `CreateOrder` 方法
- **位置**：`fmt.Sprintf("PS%s%06d", time.Now().Format("20060102150405"), rand.Intn(1000000))`
- **问题**：Go 1.22 的 `math/rand` 不是密码学安全的，高并发下可能产生重复订单号。
- **修复方案**：
  ```go
  import "crypto/rand"
  import "encoding/hex"

  func generateOrderNo() string {
      b := make([]byte, 3)
      crypto_rand.Read(b)
      return fmt.Sprintf("PS%s%s", time.Now().Format("20060102150405"), hex.EncodeToString(b))
  }
  ```

---

### #10 Error() 响应全部返回 HTTP 200

- **文件**：`internal/pkg/response/response.go` — `Error()` 方法
- **问题**：业务错误（如参数错误、权限不足）全部返回 HTTP 200，通过 body 的 `code` 字段区分。这会导致 HTTP 监控（如 Prometheus、APM）无法正确统计错误率，中间件（如 rate limiter、WAF）也无法基于状态码做决策。
- **修复方案**：调用方已有 `BadRequest`、`Unauthorized`、`Forbidden` 等方法返回正确状态码，但大量 handler 在调用 `response.Error(c, http.StatusBadRequest, "...)"` 时实际返回了 200。建议修改 `Error()` 使其也设置 HTTP 状态码：
  ```go
  func Error(c *gin.Context, httpCode int, message string) {
      c.JSON(httpCode, Response{Code: -1, Message: message, Data: nil})
  }
  ```
  或者将所有调用 `response.Error(c, 400, ...)` 的地方改为 `response.BadRequest(c, ...)`。

---

## 小问题（🟢 建议修复）

### #11 CreatePhotoSetTags 逐条 INSERT

- **文件**：`internal/repository/photoset_repository.go` — `CreatePhotoSetTags` 方法
- **问题**：循环逐条 INSERT 标签关联，N 个标签 N 次数据库往返。
- **修复方案**：使用 `db.Table("photoset_tags").Create(&batch)` 批量插入。

---

### #12 CreatePhotoSet 缺少事务保护

- **文件**：`internal/service/photoset_service.go` — `CreatePhotoSet` 方法
- **问题**：先创建套图、再创建标签关联、再创建图片，三步没有包事务。中间任何一步失败会导致数据不一致（如套图已创建但图片丢失）。
- **修复方案**：
  ```go
  return r.db.Transaction(func(tx *gorm.DB) error {
      if err := tx.Create(photoset).Error; err != nil { return err }
      // ... 创建标签关联
      // ... 创建图片
      return nil
  })
  ```

---

### #13 遗留的调试日志

- **文件**：`internal/http/handlers/admin_handler.go` — `BanUser` 方法
- **位置**：
  ```go
  log.Printf("[BanUser] userID=%d, rawBody=%s, contentLength=%d", id, string(bodyBytes), c.Request.ContentLength)
  log.Printf("[BanUser] JSON parse error: %v", err)
  log.Printf("[BanUser] Invalid status: %d", body.Status)
  log.Printf("[BanUser] Request received: userID=%d, status=%d", id, body.Status)
  ```
- **修复方案**：生产环境删除或改为 `DEBUG` 级别日志（配合日志分级框架）。

---

### #14 GetUsers 的 status=0 判断逻辑绕

- **文件**：`internal/http/handlers/admin_handler.go` — `GetUsers` 方法
- **位置**：
  ```go
  if req.Status > 0 || (req.Status == 0 && c.Query("status") != "") {
      query = query.Where("status = ?", req.Status)
  }
  ```
- **问题**：为了让 `status=0`（封禁）也能作为筛选条件，逻辑绕了一圈。如果前端传 `{"status": 0}` 但没有 query string `status=0`，条件不成立。
- **修复方案**：用指针类型：
  ```go
  type GetUsersRequest struct {
      Status *int `form:"status"`
  }
  if req.Status != nil {
      query = query.Where("status = ?", *req.Status)
  }
  ```

---

### #15 AdminHandler 中未使用的字段

- **文件**：`internal/http/handlers/admin_handler.go`
- **问题**：`AdminHandler` 结构体定义了 `mailService *service.MailService`，但 `NewAdminHandler` 中没有初始化（为 nil），实际是在 `TestMailConnection`、`GetMailConfig`、`SendTestMail` 中懒加载。虽然能工作，但结构不清晰。
- **修复方案**：在 `NewAdminHandler` 中统一初始化，或直接在方法内部创建（去掉结构体字段）。

---

## 修复优先级

| 优先级 | 编号 | 问题 | 影响 |
|--------|------|------|------|
| P0 | #1 | RestartBackend 命令执行 | 安全漏洞 |
| P0 | #4 | 默认密钥检测 | 安全漏洞 |
| P0 | #6 | List() GORM 语法错误 | 运行时报错 |
| P1 | #2 | config 重复加载 | 签名验证偶发失败 |
| P1 | #3 | CORS 配置错误 | 跨域请求被浏览器拒绝 |
| P1 | #5 | N+1 查询 | 性能问题 |
| P1 | #10 | HTTP 200 返回错误 | 监控失效 |
| P2 | #7 | paidStatus 重复实现 | 代码冗余 |
| P2 | #8 | Handler 绕过分层 | 架构问题 |
| P2 | #9 | rand 不安全 | 订单号可能重复 |
| P3 | #11-#15 | 其他小问题 | 代码质量 |
