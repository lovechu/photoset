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
//
//	@Summary		获取标签列表
//	@Description	获取所有标签列表，无需登录
//	@Tags			Tags
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]domain.Tag}	"成功"
//	@Router			/api/tags [get]
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
//
//	@Summary		更新标签
//	@Description	管理员更新标签名称
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"标签ID"
//	@Param			request	body		object{name=string}		true	"标签名称"
//	@Success		200		{object}	response.Response{data=object}	"更新成功"
//	@Failure		400		{object}	response.Response				"参数错误"
//	@Security		BearerAuth
//	@Router			/api/admin/tags/{id} [put]
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
//
//	@Summary		删除标签
//	@Description	管理员删除标签
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"标签ID"
//	@Success		200	{object}	response.Response{data=object}	"删除成功"
//	@Failure		400	{object}	response.Response				"无效的标签ID"
//	@Security		BearerAuth
//	@Router			/api/admin/tags/{id} [delete]
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


