package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"photoset/internal/config"
	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"photoset/internal/service"
	"photoset/internal/storage"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	photosetRepo *repository.PhotoSetRepository
	orderRepo    *repository.OrderRepository
	orderService *service.OrderService
	settingRepo  *repository.SiteSettingRepository
	logRepo      *repository.AdminLogRepository
	userRepo     repository.UserRepository
	userService  service.UserService
	mailService  *service.MailService
}

func NewAdminHandler(photosetRepo *repository.PhotoSetRepository, orderRepo *repository.OrderRepository, orderService *service.OrderService) *AdminHandler {
	userRepo := repository.NewUserRepository()
	return &AdminHandler{
		photosetRepo: photosetRepo,
		orderRepo:    orderRepo,
		orderService: orderService,
		settingRepo:  repository.NewSiteSettingRepository(),
		logRepo:      repository.NewAdminLogRepository(),
		userRepo:     userRepo,
		userService:  service.NewUserService(userRepo),
	}
}

// recordLog 记录操作日志
func (h *AdminHandler) recordLog(c *gin.Context, action, target, detail string) {
	adminID, exists := c.Get("user_id")
	if !exists {
		return
	}
	adminName, _ := c.Get("username")
	var uid uint
	switch v := adminID.(type) {
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	case int:
		uid = uint(v)
	default:
		return
	}
	nameStr, _ := adminName.(string)
	log := &domain.AdminLog{
		AdminID:   uid,
		AdminName: nameStr,
		Action:    action,
		Target:    target,
		Detail:    detail,
		IP:        c.ClientIP(),
	}
	go h.logRepo.Create(log) // 异步记录，不阻塞
}

// GetAdminLogs 获取操作日志列表
func (h *AdminHandler) GetAdminLogs(c *gin.Context) {
	var req struct {
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
		Action   string `form:"action"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	logs, total, err := h.logRepo.List(req.Page, req.PageSize, req.Action)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取日志失败")
		return
	}

	response.Success(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetPhotoSetsByStatus 获取指定状态的套图列表
func (h *AdminHandler) GetPhotoSetsByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		status = "pending"
	}

	photosets, err := h.photosetRepo.ListByStatus(status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取套图列表失败")
		return
	}

	response.Success(c, photosets)
}

// ApprovePhotoSet 审核通过套图
func (h *AdminHandler) ApprovePhotoSet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	if err := h.photosetRepo.UpdateStatus(uint(id), "published"); err != nil {
		response.Error(c, http.StatusInternalServerError, "审核通过失败")
		return
	}

	response.Success(c, gin.H{"message": "审核通过"})
	h.recordLog(c, "approve", "套图#"+idStr, "审核通过")
}

// RejectPhotoSet 审核拒绝套图
func (h *AdminHandler) RejectPhotoSet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.photosetRepo.UpdateStatus(uint(id), "draft"); err != nil {
		response.Error(c, http.StatusInternalServerError, "审核拒绝失败")
		return
	}

	response.Success(c, gin.H{
		"message": "已拒绝",
		"reason":  req.Reason,
	})
	h.recordLog(c, "reject", "套图#"+idStr, "拒绝原因: "+req.Reason)
}

// GetUsers 用户列表（不含密码，支持分页、角色筛选、关键字搜索）
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var req struct {
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
		Role     string `form:"role"`
		Status   int    `form:"status"`
		Keyword  string `form:"keyword"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	users, total, err := h.userRepo.List(req.Page, req.PageSize, req.Role, req.Keyword, req.Status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetUserDetail 获取用户详情
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, photoSetCount, orderCount, totalSpent, favoriteCount, err := h.userRepo.FindByIDWithStats(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"user":             user,
		"photoset_count":   photoSetCount,
		"order_count":      orderCount,
		"total_spent":      totalSpent,
		"favorite_count":   favoriteCount,
	})
}

// BanUser 封号/解封
func (h *AdminHandler) BanUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 手动解析 body，避免 ShouldBindJSON 读取空 body 的问题
	var body struct {
		Status int `json:"status"`
	}
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	log.Printf("[BanUser] userID=%d, rawBody=%s, contentLength=%d", id, string(bodyBytes), c.Request.ContentLength)
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			log.Printf("[BanUser] JSON parse error: %v", err)
			response.Error(c, http.StatusBadRequest, "参数格式错误")
			return
		}
	}

	if body.Status != 0 && body.Status != 1 {
		log.Printf("[BanUser] Invalid status: %d", body.Status)
		response.Error(c, http.StatusBadRequest, "参数错误，status 只能为 0 或 1")
		return
	}

	log.Printf("[BanUser] Request received: userID=%d, status=%d", id, body.Status)

	if err := h.userRepo.UpdateStatus(uint(id), body.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, "操作失败")
		return
	}

	msg := "已解封"
	if body.Status == 0 {
		msg = "已封号"
	}
	response.Success(c, gin.H{"message": msg})
	actionName := "unban_user"
	if body.Status == 0 {
		actionName = "ban_user"
	}
	h.recordLog(c, actionName, "用户#"+idStr, "状态改为 "+strconv.Itoa(body.Status))
}

// UpdateUserRole 更新用户角色
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=guest user member creator admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，role 只能为 guest, user, member, creator, admin 其中之一")
		return
	}

	if err := h.userRepo.UpdateRole(uint(id), req.Role); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新角色失败")
		return
	}

	response.Success(c, gin.H{"message": "角色更新成功"})
	h.recordLog(c, "role_change", "用户#"+idStr, "角色改为 "+req.Role)
}

// ResetUserPassword 管理员重置用户密码
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，新密码长度不能少于6位")
		return
	}

	if err := h.userService.ResetPassword(uint(id), req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "密码重置成功"})
	h.recordLog(c, "reset_password", "用户#"+idStr, "管理员重置密码")
}

// Stats 平台统计
func (h *AdminHandler) Stats(c *gin.Context) {
	totalUsers, err := h.userRepo.CountAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	totalPhotoSets, err := h.photosetRepo.CountAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	pendingReviews, err := h.photosetRepo.CountByStatus("pending")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	totalOrders, totalRevenue, err := h.orderRepo.CountStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, gin.H{
		"total_users":      totalUsers,
		"total_photosets":  totalPhotoSets,
		"total_orders":     totalOrders,
		"total_revenue":    totalRevenue,
		"pending_reviews":  pendingReviews,
	})
}

// StatsTrend 获取趋势数据（近 N 天）
func (h *AdminHandler) StatsTrend(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if n, err := strconv.Atoi(d); err == nil && n >= 1 && n <= 30 {
			days = n
		}
	}

	type TrendItem struct {
		Date      string  `json:"date"`
		NewUsers  int64   `json:"new_users"`
		NewOrders int64   `json:"new_orders"`
		Revenue   float64 `json:"revenue"`
		NewSets   int64   `json:"new_photosets"`
	}

	// 生成日期范围
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(days - 1))

	// 使用 repository 获取趋势数据
	stats, err := h.photosetRepo.GetTrendStats(startDate)
	if err != nil {
		response.ServerError(c, "查询趋势数据失败")
		return
	}

	// 构建 map 用于合并
	dataMap := make(map[string]TrendItem)
	for _, stat := range stats {
		dateStr, _ := stat["date"].(string)
		if dateStr == "" {
			continue
		}
		item := TrendItem{
			Date:      dateStr[5:], // 格式: 2026-04-20 → 04-20
			NewUsers:  toInt64(stat["new_users"]),
			NewOrders: toInt64(stat["new_orders"]),
			Revenue:   toFloat64(stat["revenue"]),
			NewSets:   toInt64(stat["new_sets"]),
		}
		dataMap[dateStr] = item
	}

	// 补齐缺失的日期（某天没有任何数据也要显示）
	var items []TrendItem
	for i := days - 1; i >= 0; i-- {
		dayTime := now.AddDate(0, 0, -i)
		key := dayTime.Format("2006-01-02")
		if item, ok := dataMap[key]; ok {
			items = append(items, item)
		} else {
			items = append(items, TrendItem{Date: dayTime.Format("01-02")})
		}
	}

	response.Success(c, gin.H{
		"days":  days,
		"trend": items,
	})
}

// toInt64 安全转换 interface{} 为 int64
func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case float64:
		return int64(val)
	default:
		return 0
	}
}

// toFloat64 安全转换 interface{} 为 float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	default:
		return 0
	}
}

// AdminRefund 管理员退款
func (h *AdminHandler) AdminRefund(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}

	if err := h.orderService.AdminRefundOrder(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "退款成功"})
	h.recordLog(c, "refund", "订单#"+idStr, "管理员退款")
}

// GetOrders 获取订单列表（管理员）
func (h *AdminHandler) GetOrders(c *gin.Context) {
	var req struct {
		PageNumber int    `form:"page,default=1"`
		PageSize   int    `form:"size,default=20"`
		Status     string `form:"status"`
		UserID     string `form:"user_id"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	orders, total, err := h.orderRepo.List(req.PageNumber, req.PageSize, req.Status, req.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  req.PageNumber,
		"size":  req.PageSize,
	})
}

// BatchApprovePhotoSets 批量审核通过套图
func (h *AdminHandler) BatchApprovePhotoSets(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var count int
	for _, id := range req.IDs {
		if err := h.photosetRepo.UpdateStatus(id, "published"); err == nil {
			count++
			h.recordLog(c, "batch_approve", "套图#"+strconv.Itoa(int(id)), "批量审核通过")
		}
	}

	response.Success(c, gin.H{
		"message": fmt.Sprintf("成功通过 %d 个套图", count),
		"count":   count,
	})
}

// BatchRejectPhotoSets 批量审核拒绝套图
func (h *AdminHandler) BatchRejectPhotoSets(c *gin.Context) {
	var req struct {
		IDs    []uint `json:"ids" binding:"required"`
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var count int
	for _, id := range req.IDs {
		if err := h.photosetRepo.UpdateStatus(id, "draft"); err == nil {
			count++
			h.recordLog(c, "batch_reject", "套图#"+strconv.Itoa(int(id)), "批量拒绝原因: "+req.Reason)
		}
	}

	response.Success(c, gin.H{
		"message": fmt.Sprintf("成功拒绝 %d 个套图", count),
		"count":   count,
	})
}

// BatchDeletePhotoSets 批量删除套图
func (h *AdminHandler) BatchDeletePhotoSets(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var count int
	for _, id := range req.IDs {
		if err := h.photosetRepo.Delete(id); err == nil {
			count++
			h.recordLog(c, "batch_delete", "套图#"+strconv.Itoa(int(id)), "批量删除")
		}
	}

	response.Success(c, gin.H{
		"message": fmt.Sprintf("成功删除 %d 个套图", count),
		"count":   count,
	})
}

// ExportUsers 导出用户列表为 CSV
func (h *AdminHandler) ExportUsers(c *gin.Context) {
	role := c.Query("role")
	keyword := c.Query("keyword")
	status, _ := strconv.Atoi(c.Query("status"))

	users, _, err := h.userRepo.List(1, 10000, role, keyword, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "导出用户列表失败")
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=users.csv")
	// 写入 UTF-8 BOM 以兼容 Excel
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	c.Writer.WriteString("ID,昵称,邮箱,角色,状态,注册时间,最后登录\n")
	for _, u := range users {
		statusStr := "正常"
		if u.Status == 0 {
			statusStr = "已封禁"
		}
		line := fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s\n",
			u.ID, escapeCSV(u.Nickname), escapeCSV(u.Email), u.Role, statusStr,
			u.CreatedAt.Format("2006-01-02 15:04:05"), u.LastLoginAt.Format("2006-01-02 15:04:05"))
		c.Writer.WriteString(line)
	}
}

// ExportOrders 导出订单列表为 CSV
func (h *AdminHandler) ExportOrders(c *gin.Context) {
	status := c.Query("status")
	userID := c.Query("user_id")

	orders, _, err := h.orderRepo.List(1, 10000, status, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "导出订单列表失败")
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=orders.csv")
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	c.Writer.WriteString("订单ID,订单号,用户ID,套图ID,会员ID,类型,金额,状态,创建时间,支付时间\n")
	for _, o := range orders {
		paidAt := ""
		if o.PaidAt != nil {
			paidAt = o.PaidAt.Format("2006-01-02 15:04:05")
		}
		photosetID := uint(0)
		if o.PhotoSetID != nil {
			photosetID = *o.PhotoSetID
		}
		membershipID := uint(0)
		if o.MembershipID != nil {
			membershipID = *o.MembershipID
		}
		line := fmt.Sprintf("%d,%s,%d,%d,%d,%s,%.2f,%s,%s,%s\n",
			o.ID, escapeCSV(o.OrderNo), o.UserID, photosetID, membershipID, o.Type,
			o.Amount, o.Status, o.CreatedAt.Format("2006-01-02 15:04:05"), paidAt)
		c.Writer.WriteString(line)
	}
}

// escapeCSV 转义 CSV 字段中的逗号和双引号
func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

// RestartServer 重启后端服务（依赖外部进程管理器，如 Docker restart policy 或 supervisord）
func (h *AdminHandler) RestartServer(c *gin.Context) {
	response.Success(c, gin.H{"message": "正在重启后端服务"})
	// 在另一个 goroutine 中退出进程，确保响应先发送给客户端
	go func() {
		// 给响应一点时间发送出去
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}

// TestStorageConnection 测试存储连接
func (h *AdminHandler) TestStorageConnection(c *gin.Context) {
	var req struct {
		StorageType   string `json:"storage_type"`
		S3Endpoint    string `json:"s3_endpoint"`
		S3Region      string `json:"s3_region"`
		S3AccessKey   string `json:"s3_access_key"`
		S3SecretKey   string `json:"s3_secret_key"`
		S3Bucket      string `json:"s3_bucket"`
		CDNDomain     string `json:"cdn_domain"`
		R2AccountID   string `json:"r2_account_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	settings := map[string]interface{}{
		"storage_type":   req.StorageType,
		"s3_endpoint":    req.S3Endpoint,
		"s3_region":      req.S3Region,
		"s3_access_key":  req.S3AccessKey,
		"s3_secret_key":  req.S3SecretKey,
		"s3_bucket":      req.S3Bucket,
		"cdn_domain":     req.CDNDomain,
		"r2_account_id":  req.R2AccountID,
	}

	stor, err := storage.NewStorageFromSettings(settings)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := stor.TestConnection(); err != nil {
		response.Error(c, http.StatusBadRequest, "连接失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":      "连接成功",
		"storage_type": req.StorageType,
	})
}

// GetStorageStatus 获取当前存储状态
func (h *AdminHandler) GetStorageStatus(c *gin.Context) {
	cfg := config.Load()
	stor, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "存储初始化失败")
		return
	}

	status := gin.H{
		"type":       cfg.Storage.Type,
		"configurable": true,
	}

	switch stor.Type() {
	case storage.StorageLocal:
		status["label"] = "本地存储"
		status["path"] = cfg.Storage.LocalPath
	case storage.StorageS3:
		status["label"] = "云存储 (S3 兼容)"
		if cfg.Storage.R2PublicURL != "" {
			status["cdn_domain"] = cfg.Storage.R2PublicURL
		}
		if cfg.Storage.R2AccountID != "" {
			status["provider"] = "Cloudflare R2"
		} else if cfg.Storage.S3Endpoint != "" {
			status["provider"] = cfg.Storage.S3Endpoint
		}
	}

	// 隐藏敏感信息
	status["s3_access_key_set"] = cfg.Storage.S3AccessKey != ""
	status["s3_secret_key_set"] = cfg.Storage.S3SecretKey != ""
	status["s3_bucket_set"] = cfg.Storage.S3Bucket != ""

	response.Success(c, status)
}

// GetSettings 获取所有站点设置
func (h *AdminHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingRepo.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取配置失败")
		return
	}
	response.Success(c, settings)
}

// GetPublicSettings 获取公开的站点设置（不需要认证，供移动端使用）
func (h *AdminHandler) GetPublicSettings(c *gin.Context) {
	settings, err := h.settingRepo.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取配置失败")
		return
	}

	// 过滤敏感字段，只返回允许公开的配置
	publicSettings := make(map[string]interface{})
	publicKeys := []string{
		"site_title", "site_description", "site_keywords", "about_me",
		"logo_url", "favicon_url", "site_icp", "copyright_year", "about_content",
		"terms_content", "privacy_content", "help_content", "contact_content",
		// 域名配置（供移动端使用）
		"site_url", "api_url", "dev_api_url",
		// 导航菜单（前端使用）
		"nav_menu",
	}
	for key, value := range settings {
		for _, allowed := range publicKeys {
			if key == allowed {
				publicSettings[key] = value
				break
			}
		}
		// SMTP设置、水印设置、邮件密码等敏感信息不对外暴露
	}

	response.Success(c, publicSettings)
}

// UpdateSettings 批量更新站点设置
func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.settingRepo.BatchUpsert(data); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存配置失败")
		return
	}

	response.Success(c, gin.H{"message": "保存成功"})
}

// TestMailConnection 测试邮件 SMTP 连接
func (h *AdminHandler) TestMailConnection(c *gin.Context) {
	if h.mailService == nil {
		h.mailService = service.NewMailService()
	}

	success, message := h.mailService.TestConnection()
	if success {
		response.Success(c, gin.H{"message": message})
	} else {
		response.Error(c, http.StatusBadRequest, message)
	}
}

// GetMailConfig 获取邮件配置信息（不含密码）
func (h *AdminHandler) GetMailConfig(c *gin.Context) {
	if h.mailService == nil {
		h.mailService = service.NewMailService()
	}

	info := h.mailService.GetConfigInfo()
	response.Success(c, info)
}

// SendTestMail 发送测试邮件
func (h *AdminHandler) SendTestMail(c *gin.Context) {
	var req struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，请提供有效的邮箱地址")
		return
	}

	if h.mailService == nil {
		h.mailService = service.NewMailService()
	}

	if req.Subject == "" {
		req.Subject = "PhotoSet 测试邮件"
	}
	if req.Body == "" {
		req.Body = "<h1>测试成功！</h1><p>这是一封来自 PhotoSet 的测试邮件。</p>"
	}

	if err := h.mailService.Send(req.To, req.Subject, req.Body); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("发送失败: %v", err))
		return
	}

	response.Success(c, gin.H{"message": fmt.Sprintf("测试邮件已发送到 %s", req.To)})
}

// GetWatermarkInfo 获取水印配置信息
func (h *AdminHandler) GetWatermarkInfo(c *gin.Context) {
	watermarkService := service.NewWatermarkService()
	info := watermarkService.GetWatermarkInfo()
	response.Success(c, info)
}

// ==================== 开发者中心 API ====================

// ListApiKeys 获取 API 密钥列表
func (h *AdminHandler) ListApiKeys(c *gin.Context) {
	repo := repository.NewApiKeyRepository()
	keys, err := repo.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取 API 密钥列表失败")
		return
	}
	response.Success(c, keys)
}

// CreateApiKey 创建新的 API 密钥
func (h *AdminHandler) CreateApiKey(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=2,max=50"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请提供密钥名称（2-50字符）")
		return
	}

	// 获取当前管理员 ID
	adminID, _ := c.Get("user_id")
	var uid uint
	switch v := adminID.(type) {
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	case int:
		uid = uint(v)
	default:
		uid = 0
	}

	repo := repository.NewApiKeyRepository()
	apiKey, err := repo.Create(req.Name, uid)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建 API 密钥失败")
		return
	}

	h.recordLog(c, "create_api_key", "API密钥", "创建: "+req.Name)
	response.Success(c, gin.H{
		"message": "API 密钥创建成功，请妥善保存以下信息（仅显示一次）：",
		"key":     apiKey,
	})
}

// DeleteApiKey 删除 API 密钥
func (h *AdminHandler) DeleteApiKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的密钥 ID")
		return
	}

	repo := repository.NewApiKeyRepository()
	if err := repo.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	h.recordLog(c, "delete_api_key", "API密钥", "删除 ID: "+idStr)
	response.Success(c, gin.H{"message": "API 密钥已删除"})
}

// GetApiDocs 获取 API 文档
func (h *AdminHandler) GetApiDocs(c *gin.Context) {
	docs := []map[string]interface{}{
		{
			"category": "套图",
			"endpoints": []map[string]string{
				{"method": "GET", "path": "/api/photosets", "desc": "获取套图列表"},
				{"method": "GET", "path": "/api/photosets/:id", "desc": "获取套图详情"},
				{"method": "GET", "path": "/api/photosets/advanced", "desc": "高级搜索套图"},
				{"method": "POST", "path": "/api/photosets", "desc": "创建套图（需认证）"},
				{"method": "PUT", "path": "/api/photosets/:id", "desc": "更新套图（需认证）"},
				{"method": "DELETE", "path": "/api/photosets/:id", "desc": "删除套图（需认证）"},
			},
		},
		{
			"category": "用户",
			"endpoints": []map[string]string{
				{"method": "GET", "path": "/api/users/profile", "desc": "获取当前用户信息（需认证）"},
				{"method": "PUT", "path": "/api/auth/password", "desc": "修改密码（需认证）"},
			},
		},
		{
			"category": "收藏",
			"endpoints": []map[string]string{
				{"method": "GET", "path": "/api/favorites", "desc": "获取收藏列表（需认证）"},
				{"method": "POST", "path": "/api/favorites/:photosetId", "desc": "添加收藏（需认证）"},
				{"method": "DELETE", "path": "/api/favorites/:photosetId", "desc": "取消收藏（需认证）"},
			},
		},
		{
			"category": "订单",
			"endpoints": []map[string]string{
				{"method": "GET", "path": "/api/orders", "desc": "获取订单列表（需认证）"},
				{"method": "POST", "path": "/api/orders", "desc": "创建订单（需认证）"},
				{"method": "POST", "path": "/api/orders/:id/refund", "desc": "申请退款（需认证）"},
			},
		},
		{
			"category": "会员套餐",
			"endpoints": []map[string]string{
				{"method": "GET", "path": "/api/memberships", "desc": "获取会员套餐列表"},
			},
		},
		{
			"category": "公开信息",
			"endpoints": []map[string]string{
				{"method": "GET", "path": "/api/tags", "desc": "获取标签列表"},
				{"method": "GET", "path": "/api/categories", "desc": "获取分类列表"},
				{"method": "GET", "path": "/api/pages/:slug", "desc": "获取页面内容"},
				{"method": "GET", "path": "/api/settings", "desc": "获取站点公开设置"},
				{"method": "GET", "path": "/api/health", "desc": "健康检查"},
			},
		},
	}

	response.Success(c, gin.H{
		"docs":         docs,
		"auth_header":  "Authorization",
		"auth_format":  "Bearer <token>",
		"content_type": "application/json",
	})
}

// GetSignUrlDocs 获取图片签名 URL 文档
func (h *AdminHandler) GetSignUrlDocs(c *gin.Context) {
	docs := gin.H{
		"description": "付费图片使用签名 URL 进行访问验证，防止盗链",
		"signature_required": true,
		"query_params": []map[string]string{
			{"name": "sign", "desc": "HMAC-SHA256 签名"},
			{"name": "expires", "desc": "签名过期时间戳（Unix）"},
		},
		"signature_algorithm": "HMAC-SHA256",
		"signature_example": gin.H{
			"message": "path?expires=<timestamp>",
			"key": "<your_secret_key>",
			"output": "hex encoded hmac",
		},
		"code_example": gin.H{
			"python": `import hmac
import hashlib
import time

def generate_sign_url(path, secret_key, expires=3600):
    expires_at = int(time.time()) + expires
    message = "%s?expires=%d" % (path, expires_at)
    sign = hmac.new(secret_key.encode(), message.encode(), hashlib.sha256).hexdigest()
    return "%s&sign=%s" % (message, sign)`,
			"javascript": `// generateSignUrl(path, secretKey, expires=3600)
const expiresAt = Math.floor(Date.now() / 1000) + expires;
const message = path + "?expires=" + expiresAt;
const sign = crypto.createHmac('sha256', secretKey).update(message).digest('hex');
return message + "&sign=" + sign;`,
		},
	}

	response.Success(c, docs)
}
