package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	photosetRepo *repository.PhotoSetRepository
	orderRepo    *repository.OrderRepository
	orderService *service.OrderService
}

func NewAdminHandler(photosetRepo *repository.PhotoSetRepository, orderRepo *repository.OrderRepository, orderService *service.OrderService) *AdminHandler {
	return &AdminHandler{
		photosetRepo: photosetRepo,
		orderRepo:    orderRepo,
		orderService: orderService,
	}
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
}

// GetUsers 用户列表（不含密码，支持分页和角色筛选）
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var req struct {
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
		Role     string `form:"role"`
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

// BanUser 封号/解封
func (h *AdminHandler) BanUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，status 只能为 0 或 1")
		return
	}

	db := database.GetMySQL()
	if err := db.Model(&domain.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "操作失败")
		return
	}

	msg := "已解封"
	if req.Status == 0 {
		msg = "已封号"
	}
	response.Success(c, gin.H{"message": msg})
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
}
