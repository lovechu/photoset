package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type CollectionHandler struct {
	service *service.CollectionService
}

func NewCollectionHandler(service *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}

// @Summary      创建合集
// @Description  创建新的套图合集
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        body  body  object  true  "合集信息 {name, description, is_public}"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Security     BearerAuth
// @Router       /api/collections [post]
// Create 创建合集
func (h *CollectionHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,max=100"`
		Description string `json:"description" binding:"max=500"`
		IsPublic    bool   `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	collection, err := h.service.CreateCollection(userID.(uint), req.Name, req.Description, req.IsPublic)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, collection)
}

// @Summary      更新合集
// @Description  修改合集的名称、描述或可见性
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        id    path  int     true  "合集ID"
// @Param        body  body  object  true  "更新信息 {name, description, is_public}"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Security     BearerAuth
// @Router       /api/collections/{id} [put]
// Update 更新合集
func (h *CollectionHandler) Update(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的合集ID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"max=100"`
		Description string `json:"description" binding:"max=500"`
		IsPublic    *bool  `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	collection, err := h.service.UpdateCollection(userID.(uint), uint(collectionID), req.Name, req.Description, req.IsPublic)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, collection)
}

// @Summary      删除合集
// @Description  删除指定的合集及其关联数据
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "合集ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的合集ID"
// @Security     BearerAuth
// @Router       /api/collections/{id} [delete]
// Delete 删除合集
func (h *CollectionHandler) Delete(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的合集ID")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.service.DeleteCollection(userID.(uint), uint(collectionID)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// @Summary      获取合集详情
// @Description  获取指定合集的详细信息和包含的套图列表
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "合集ID"
// @Success      200  {object}  response.Response  "合集详情"
// @Failure      400  {object}  response.Response  "无效的合集ID"
// @Security     BearerAuth
// @Router       /api/collections/{id} [get]
// Get 获取合集详情
func (h *CollectionHandler) Get(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的合集ID")
		return
	}

	collection, err := h.service.GetCollection(uint(collectionID))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, collection)
}

// @Summary      获取合集列表
// @Description  获取当前用户的所有合集列表
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        page      query int  false "页码"
// @Param        page_size query int  false "每页数量"
// @Success      200  {object}  response.Response  "合集列表"
// @Security     BearerAuth
// @Router       /api/collections [get]
// List 获取用户合集列表
func (h *CollectionHandler) List(c *gin.Context) {
	var req struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
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

	userID, _ := c.Get("user_id")
	collections, total, err := h.service.ListCollections(userID.(uint), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取合集列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":      collections,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// @Summary      添加套图到合集
// @Description  将指定套图添加到合集
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        id    path  int     true  "合集ID"
// @Param        body  body  object  true  "套图信息 {photoset_id}"
// @Success      200  {object}  response.Response  "添加成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Security     BearerAuth
// @Router       /api/collections/{id}/items [post]
// AddItem 添加套图到合集
func (h *CollectionHandler) AddItem(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的合集ID")
		return
	}

	var req struct {
		PhotoSetID uint `json:"photoset_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.service.AddItem(userID.(uint), uint(collectionID), req.PhotoSetID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// @Summary      从合集移除套图
// @Description  从指定合集中移除套图
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        id          path  int  true  "合集ID"
// @Param        photosetId  path  int  true  "套图ID"
// @Success      200  {object}  response.Response  "移除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Security     BearerAuth
// @Router       /api/collections/{id}/items/{photosetId} [delete]
// RemoveItem 从合集移除套图
func (h *CollectionHandler) RemoveItem(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的合集ID")
		return
	}

	photosetID, err := strconv.ParseUint(c.Param("photosetId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.service.RemoveItem(userID.(uint), uint(collectionID), uint(photosetID)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// @Summary      批量添加套图到合集
// @Description  批量将多个套图添加到指定合集
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        id    path  int     true   "合集ID"
// @Param        body  body  object  true   "套图ID数组 {photoset_ids}"
// @Success      200  {object}  response.Response  "批量添加成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Security     BearerAuth
// @Router       /api/collections/{id}/items/batch [post]
// BatchAddItems 批量添加套图到合集
func (h *CollectionHandler) BatchAddItems(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的合集ID")
		return
	}

	var req struct {
		PhotoSetIDs []uint `json:"photoset_ids" binding:"required,min=1,max=50"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误：需要 photoset_ids 数组（1-50个）")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.service.BatchAddItems(userID.(uint), uint(collectionID), req.PhotoSetIDs); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// @Summary      获取套图所在合集
// @Description  获取指定套图所属的合集列表
// @Tags         Collections
// @Accept       json
// @Produce      json
// @Param        photosetId  path  int  true  "套图ID"
// @Success      200  {object}  response.Response  "合集列表"
// @Failure      400  {object}  response.Response  "参数错误"
// @Security     BearerAuth
// @Router       /api/collections/by-photoset/{photosetId} [get]
// GetCollectionsForPhotoset 获取包含指定套图的用户合集
func (h *CollectionHandler) GetCollectionsForPhotoset(c *gin.Context) {
	photosetID, err := strconv.ParseUint(c.Param("photosetId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	userID, _ := c.Get("user_id")
	collections, err := h.service.GetCollectionsContaining(userID.(uint), uint(photosetID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取合集失败")
		return
	}
	response.Success(c, collections)
}
