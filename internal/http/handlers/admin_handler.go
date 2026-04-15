package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"photoset/internal/service"
	"photoset/internal/storage"
	"photoset/internal/config"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	photosetRepo *repository.PhotoSetRepository
	orderRepo    *repository.OrderRepository
	orderService *service.OrderService
	settingRepo  *repository.SiteSettingRepository
	logRepo      *repository.AdminLogRepository
}

func NewAdminHandler(photosetRepo *repository.PhotoSetRepository, orderRepo *repository.OrderRepository, orderService *service.OrderService) *AdminHandler {
	return &AdminHandler{
		photosetRepo: photosetRepo,
		orderRepo:    orderRepo,
		orderService: orderService,
		settingRepo:  repository.NewSiteSettingRepository(),
		logRepo:      repository.NewAdminLogRepository(),
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

	var photosets []domain.PhotoSet
	db := database.GetMySQL()
	query := db.Model(&domain.PhotoSet{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Preload("User").Order("created_at DESC").Find(&photosets).Error; err != nil {
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

	db := database.GetMySQL()
	if err := db.Model(&domain.PhotoSet{}).Where("id = ?", id).Update("status", "published").Error; err != nil {
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

	db := database.GetMySQL()
	if err := db.Model(&domain.PhotoSet{}).Where("id = ?", id).Update("status", "draft").Error; err != nil {
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

	db := database.GetMySQL()
	query := db.Model(&domain.User{})

	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}
	if req.Status > 0 || (req.Status == 0 && c.Query("status") != "") {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		query = query.Where("nickname LIKE ? OR email LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	var users []domain.User
	offset := (req.Page - 1) * req.PageSize
	if err := query.Select("id, nickname, email, role, status, created_at, last_login_at, membership_expires").
		Order("created_at DESC").
		Offset(offset).Limit(req.PageSize).
		Find(&users).Error; err != nil {
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

	db := database.GetMySQL()
	var user domain.User
	if err := db.Select("id, nickname, email, role, status, created_at, last_login_at, membership_expires").
		Where("id = ?", id).First(&user).Error; err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 获取用户相关统计
	var photoSetCount int64
	db.Model(&domain.PhotoSet{}).Where("user_id = ?", id).Count(&photoSetCount)

	var orderCount int64
	var totalSpent float64
	db.Model(&domain.Order{}).Where("user_id = ? AND status = ?", id, "completed").Count(&orderCount)
	db.Model(&domain.Order{}).Where("user_id = ? AND status = ?", id, "completed").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalSpent)

	response.Success(c, gin.H{
		"user":             user,
		"photoset_count":   photoSetCount,
		"order_count":      orderCount,
		"total_spent":      totalSpent,
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

	db := database.GetMySQL()
	if err := db.Model(&domain.User{}).Where("id = ?", id).Update("status", body.Status).Error; err != nil {
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

	db := database.GetMySQL()
	if err := db.Model(&domain.User{}).Where("id = ?", id).Update("role", req.Role).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新角色失败")
		return
	}

	response.Success(c, gin.H{"message": "角色更新成功"})
	h.recordLog(c, "role_change", "用户#"+idStr, "角色改为 "+req.Role)
}

// Stats 平台统计
func (h *AdminHandler) Stats(c *gin.Context) {
	db := database.GetMySQL()

	var totalUsers int64
	db.Model(&domain.User{}).Count(&totalUsers)

	var totalPhotoSets int64
	db.Model(&domain.PhotoSet{}).Count(&totalPhotoSets)

	var pendingReviews int64
	db.Model(&domain.PhotoSet{}).Where("status = ?", "pending").Count(&pendingReviews)

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

	db := database.GetMySQL()

	type TrendItem struct {
		Date       string  `json:"date"`
		NewUsers   int64   `json:"new_users"`
		NewOrders  int64   `json:"new_orders"`
		Revenue    float64 `json:"revenue"`
		NewSets    int64   `json:"new_photosets"`
	}

	var items []TrendItem
	for i := days - 1; i >= 0; i-- {
		dayTime := time.Now().AddDate(0, 0, -i)
		dayStart := time.Date(dayTime.Year(), dayTime.Month(), dayTime.Day(), 0, 0, 0, 0, dayTime.Location()).Unix()
		dayEnd := dayStart + 86400
		dateStr := dayTime.Format("01-02")

		var newUsers int64
		db.Model(&domain.User{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&newUsers)

		var newOrders int64
		db.Model(&domain.Order{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&newOrders)

		var revenue float64
		db.Model(&domain.Order{}).Where("created_at >= ? AND created_at < ? AND status = ?", dayStart, dayEnd, "completed").
			Select("COALESCE(SUM(amount), 0)").Scan(&revenue)

		var newSets int64
		db.Model(&domain.PhotoSet{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&newSets)

		items = append(items, TrendItem{
			Date:      dateStr,
			NewUsers:  newUsers,
			NewOrders: newOrders,
			Revenue:   revenue / 100, // 分转元
			NewSets:   newSets,
		})
	}

	response.Success(c, gin.H{
		"days":   days,
		"trend":  items,
	})
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

	var orders []domain.Order
	db := database.GetMySQL()
	query := db.Model(&domain.Order{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.UserID != "" {
		if userID, err := strconv.ParseUint(req.UserID, 10, 32); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	var total int64
	query.Count(&total)

	offset := (req.PageNumber - 1) * req.PageSize
	if err := query.
		Preload("User").
		Preload("PhotoSet").
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&orders).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单列表失败")
		return
	}

	response.Success(c, gin.H{
		"total": total,
		"page":  req.PageNumber,
		"size":  req.PageSize,
		"data":  orders,
	})
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

// GetPublicSettings 获取公开的站点设置（不需要认证）
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
