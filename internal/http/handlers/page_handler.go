package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	service *service.PageService
}

func NewPageHandler(service *service.PageService) *PageHandler {
	return &PageHandler{service: service}
}

// GetBySlug 公开获取页面内容
func (h *PageHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	page, err := h.service.GetPublishedPage(slug)
	if err != nil {
		response.Error(c, http.StatusNotFound, "页面不存在")
		return
	}
	response.Success(c, page)
}

// ListPublished 公开列表（用于站点地图等）
func (h *PageHandler) ListPublished(c *gin.Context) {
	pages, err := h.service.ListPublishedPages()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取页面列表失败")
		return
	}
	response.Success(c, pages)
}

// AdminList 后台列表
func (h *PageHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.DefaultQuery("keyword", "")

	pages, total, err := h.service.AdminListPages(page, pageSize, keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取页面列表失败")
		return
	}
	response.Success(c, gin.H{
		"list":      pages,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// AdminCreate 创建页面
func (h *PageHandler) AdminCreate(c *gin.Context) {
	var req CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	page := &domain.Page{
		Slug:      req.Slug,
		Title:     req.Title,
		ContentMD: req.ContentMD,
		UserID:    userID.(uint),
		Status:    req.Status,
	}

	if err := h.service.CreatePage(page); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建页面失败")
		return
	}
	response.Success(c, page)
}

// AdminUpdate 更新页面
func (h *PageHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID 格式错误")
		return
	}

	var req UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	page, err := h.service.GetPageByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "页面不存在")
		return
	}

	page.Slug = req.Slug
	page.Title = req.Title
	page.ContentMD = req.ContentMD
	page.Status = req.Status

	if err := h.service.UpdatePage(page); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新页面失败")
		return
	}
	response.Success(c, page)
}

// AdminGet 获取单个页面（后台）
func (h *PageHandler) AdminGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID 格式错误")
		return
	}

	page, err := h.service.GetPageByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "页面不存在")
		return
	}
	response.Success(c, page)
}

// AdminDelete 删除页面
func (h *PageHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID 格式错误")
		return
	}

	if err := h.service.DeletePage(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除页面失败")
		return
	}
	response.Success(c, nil)
}

// 请求体结构
type CreatePageRequest struct {
	Slug      string `json:"slug" binding:"required,max=50"`
	Title     string `json:"title" binding:"required,max=200"`
	ContentMD string `json:"content_md" binding:"required"`
	Status    string `json:"status" binding:"oneof=published draft"`
}

type UpdatePageRequest struct {
	Slug      string `json:"slug" binding:"required,max=50"`
	Title     string `json:"title" binding:"required,max=200"`
	ContentMD string `json:"content_md" binding:"required"`
	Status    string `json:"status" binding:"oneof=published draft"`
}