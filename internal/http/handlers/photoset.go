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

type PhotoSetHandler struct {
	service *service.PhotoSetService
}

func NewPhotoSetHandler(service *service.PhotoSetService) *PhotoSetHandler {
	return &PhotoSetHandler{service: service}
}

// ListRequest 套图列表请求
type ListRequest struct {
	Page      int     `form:"page" binding:"min=1"`
	PageSize  int     `form:"page_size" binding:"min=1,max=100"`
	Tag       string  `form:"tag"`
	Mine      bool    `form:"mine"`
	Keyword   string  `form:"keyword"`
	
	// 高级筛选参数
	Category   string   `form:"category"`
	PriceMin   float64  `form:"price_min"`
	PriceMax   float64  `form:"price_max"`
	IsFree     *bool    `form:"is_free"`
	SortBy     string   `form:"sort_by"`
	TimeRange  string   `form:"time_range"`
	UserID     uint     `form:"user_id"`
}

// CreateRequest 创建套图请求
type CreateRequest struct {
	Title       string   `json:"title" binding:"required,max=200"`
	Cover       string   `json:"cover" binding:"required,max=500"`
	Description string   `json:"description"`
	IsFree      int8     `json:"is_free" binding:"oneof=0 1"`
	Price       float64  `json:"price"`
	Tags        []string `json:"tags"`
	Photos      []Photo  `json:"photos"`
	Status      string   `json:"status" binding:"oneof=draft published pending"`
	Category    string   `json:"category"` // 分类 slug
}

// Photo 图片信息
type Photo struct {
	URL       string `json:"url" binding:"required,max=500"`
	SortOrder int    `json:"sort_order"`
}

// List 套图列表（基础版本，向后兼容）
func (h *PhotoSetHandler) List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 获取当前用户ID（可选）
	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	// 调用基础搜索（只使用基础的tag和keyword参数）
	photosets, total, err := h.service.GetPhotoSetList(
		req.Page, req.PageSize, 
		req.Tag, req.Keyword,
		userID, req.Mine,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取套图列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":      photosets,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// Detail 套图详情
func (h *PhotoSetHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	// 获取当前用户信息(可选鉴权)
	var userRole string
	var userID uint
	var isLoggedIn bool

	if role, exists := c.Get("user_role"); exists {
		userRole = role.(string)
		isLoggedIn = true
	}
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	// 获取套图基础信息
	photoset, err := h.service.GetPhotoSetDetailBasic(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "套图不存在")
		return
	}

	// 判断是否可以查看完整图片列表
	canViewFull := h.service.CanViewFullPhotos(photoset, userRole, userID, isLoggedIn)

	if canViewFull {
		// 可以查看完整图片列表
		photoset, err = h.service.GetPhotoSetDetail(uint(id))
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "获取套图详情失败")
			return
		}
	} else {
		// 不能查看完整图片列表，只返回封面和基础信息
		photoset.Photos = []domain.Photo{}
	}

	// 如果已登录，查询收藏状态
	if isLoggedIn {
		favRepo := repository.NewFavoriteRepository(database.GetMySQL())
		isFav, _ := favRepo.IsFavorited(userID, uint(id))
		photoset.IsFavorited = isFav
	}

	response.Success(c, photoset)
}

// UpdateRequest 更新套图请求
type UpdateRequest struct {
	Title       string   `json:"title" binding:"required,max=200"`
	Cover       string   `json:"cover" binding:"required,max=500"`
	Description string   `json:"description"`
	IsFree      int8     `json:"is_free" binding:"oneof=0 1"`
	Price       float64  `json:"price"`
	Tags        []string `json:"tags"`
	Photos      []Photo  `json:"photos"`
	Status      string   `json:"status" binding:"oneof=draft published pending"`
	Category    string   `json:"category"` // 分类 slug
}

// Update 更新套图（creator 更新自己的 / admin 更新任意）
func (h *PhotoSetHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	// 查找套图
	photoset, err := h.service.GetPhotoSetDetailBasic(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "套图不存在")
		return
	}

	// 权限校验：creator 只能改自己的，admin 无限制
	if userRole.(string) != "admin" && photoset.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权编辑此套图")
		return
	}

	// 更新基础字段
	updates := map[string]interface{}{
		"title":       req.Title,
		"cover":       req.Cover,
		"description": req.Description,
		"is_free":     req.IsFree,
		"price":       req.Price,
		"status":      req.Status,
		"category":    req.Category, // <-- 新增
	}
	if err := h.service.UpdatePhotoSet(uint(id), updates, req.Tags, toPhotos(req.Photos, uint(id))); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// Delete 删除套图（creator 删除自己的 / admin 删除任意）
func (h *PhotoSetHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	// 查找套图
	photoset, err := h.service.GetPhotoSetDetailBasic(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "套图不存在")
		return
	}

	// 权限校验：creator 只能删自己的，admin 无限制
	if userRole.(string) != "admin" && photoset.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权删除此套图")
		return
	}

	if err := h.service.DeletePhotoSet(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

func toPhotos(ps []Photo, photosetID uint) []domain.Photo {
	var result []domain.Photo
	for _, p := range ps {
		result = append(result, domain.Photo{
			PhotoSetID: photosetID,
			URL:        p.URL,
			SortOrder:  p.SortOrder,
		})
	}
	return result
}

// AdvancedList 高级搜索套图列表
func (h *PhotoSetHandler) AdvancedList(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 获取当前用户ID（可选）
	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	// 对于"仅我的作品"筛选，需要使用当前用户ID
	var filterUserID uint
	var onlyMine bool
	
	if req.Mine && userID > 0 {
		onlyMine = true
	} else if req.UserID > 0 {
		filterUserID = req.UserID
		onlyMine = true
	}

	photosets, total, err := h.service.GetPhotoSetListAdvanced(
		req.Page, req.PageSize, 
		req.Tag, req.Keyword, 
		userID, onlyMine,
		req.Category, req.PriceMin, req.PriceMax, req.IsFree, 
		req.SortBy, req.TimeRange, filterUserID,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取套图列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":      photosets,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
		"has_advanced_filters": hasAdvancedFilters(req),
	})
}

// 检查是否有高级筛选参数
func hasAdvancedFilters(req ListRequest) bool {
	return req.Category != "" ||
		req.PriceMin > 0 ||
		req.PriceMax > 0 ||
		req.IsFree != nil ||
		(req.SortBy != "" && req.SortBy != "latest") ||
		req.TimeRange != "" ||
		req.UserID > 0
}

// Create 创建套图
func (h *PhotoSetHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	// 创建套图
	photoset := &domain.PhotoSet{
		Title:       req.Title,
		Cover:       req.Cover,
		Description: req.Description,
		IsFree:      req.IsFree,
		Price:       req.Price,
		UserID:      userID.(uint),
		Status:      req.Status,
		Category:    req.Category, // <-- 新增
	}

	// 转换图片
	var photos []domain.Photo
	for _, p := range req.Photos {
		photos = append(photos, domain.Photo{
			URL:       p.URL,
			SortOrder: p.SortOrder,
		})
	}

	// 创建套图
	if err := h.service.CreatePhotoSet(photoset, req.Tags, photos); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建套图失败: "+err.Error())
		return
	}

	response.Success(c, photoset)
}

