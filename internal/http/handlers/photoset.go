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
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Tag      string `form:"tag"`
	Mine     bool   `form:"mine"`
	Keyword  string `form:"keyword"`
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
}

// Photo 图片信息
type Photo struct {
	URL       string `json:"url" binding:"required,max=500"`
	SortOrder int    `json:"sort_order"`
}

// List 套图列表
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

	photosets, total, err := h.service.GetPhotoSetList(req.Page, req.PageSize, req.Tag, req.Keyword, userID, req.Mine)
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
