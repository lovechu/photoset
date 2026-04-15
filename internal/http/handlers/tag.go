package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service *service.PhotoSetService
}

func NewTagHandler(service *service.PhotoSetService) *TagHandler {
	return &TagHandler{service: service}
}

// List 标签列表（公开）
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.service.GetAllTags()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取标签列表失败")
		return
	}

	response.Success(c, tags)
}

// AdminList 标签管理列表（管理员，支持分页和搜索）
func (h *TagHandler) AdminList(c *gin.Context) {
	var req struct {
		PageNumber int    `form:"page,default=1"`
		PageSize   int    `form:"size,default=20"`
		Keyword    string `form:"keyword"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var tags []domain.Tag
	db := database.GetMySQL()
	query := db.Model(&domain.Tag{})

	if req.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	var total int64
	query.Count(&total)

	offset := (req.PageNumber - 1) * req.PageSize
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&tags).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "获取标签列表失败")
		return
	}

	response.Success(c, gin.H{
		"total": total,
		"page":  req.PageNumber,
		"size":  req.PageSize,
		"data":  tags,
	})
}

// Create 创建标签
func (h *TagHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=20"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	tag := domain.Tag{
		Name: req.Name,
	}

	db := database.GetMySQL()
	if err := db.Create(&tag).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建标签失败")
		return
	}

	response.Success(c, tag)
}

// Update 更新标签
func (h *TagHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required,min=1,max=20"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := database.GetMySQL()
	if err := db.Model(&domain.Tag{}).Where("id = ?", id).Update("name", req.Name).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "更新标签失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// Delete 删除标签
func (h *TagHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	db := database.GetMySQL()
	if err := db.Delete(&domain.Tag{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "删除标签失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}


