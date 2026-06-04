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
	Status    string  `form:"status"` // draft/published/pending，空默认 published
	
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
//
//	@Summary		获取套图列表
//	@Description	获取套图列表，支持分页、标签、关键词等筛选条件（可选鉴权）
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"页码，默认1"
//	@Param			page_size	query	int		false	"每页数量，默认20，最大100"
//	@Param			tag			query	string	false	"标签名称"
//	@Param			mine		query	bool	false	"仅查看自己的作品"
//	@Param			keyword		query	string	false	"搜索关键词"
//	@Param			status		query	string	false	"状态筛选(draft/published/pending)"
//	@Success		200			{object}	response.Response{data=object}	"成功"
//	@Failure		400			{object}	response.Response				"参数错误"
//	@Router			/api/photosets [get]
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

	// 调用高级搜索，支持分类等参数
	photosets, total, err := h.service.GetPhotoSetListAdvanced(
		req.Page, req.PageSize,
		req.Tag, req.Keyword,
		userID, req.Mine,
		req.Category, req.PriceMin, req.PriceMax, req.IsFree,
		req.SortBy, req.TimeRange, 0,
		req.Status,
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
//
//	@Summary		获取套图详情
//	@Description	根据套图ID获取套图详情，支持可选鉴权（登录用户可查看更多信息）
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"套图ID"
//	@Success		200	{object}	response.Response{data=domain.PhotoSet}	"成功"
//	@Failure		400	{object}	response.Response							"无效的套图ID"
//	@Failure		404	{object}	response.Response							"套图不存在"
//	@Router			/api/photosets/{id} [get]
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

// Download 下载套图（获取图片 URL 列表）
//
//	@Summary		下载套图
//	@Description	获取套图的图片URL列表，需要登录认证
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"套图ID"
//	@Success		200	{object}	response.Response{data=object}	"成功"
//	@Failure		400	{object}	response.Response				"无效的套图ID"
//	@Failure		403	{object}	response.Response				"无权下载"
//	@Security		BearerAuth
//	@Router			/api/photosets/{id}/download [get]
func (h *PhotoSetHandler) Download(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	// 获取当前用户信息
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

	// 调用服务获取下载信息
	photoURLs, err := h.service.GetPhotoSetDownload(uint(id), userRole, userID, isLoggedIn)
	if err != nil {
		response.Error(c, http.StatusForbidden, err.Error())
		return
	}

	// 返回图片 URL 列表
	response.Success(c, gin.H{
		"photoset_id": uint(id),
		"photos":      photoURLs,
		"total":       len(photoURLs),
	})
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
//
//	@Summary		更新套图
//	@Description	更新套图信息，仅创建者或管理员可操作
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"套图ID"
//	@Param			request	body		UpdateRequest	true	"更新套图请求体"
//	@Success		200		{object}	response.Response{data=object}	"更新成功"
//	@Failure		400		{object}	response.Response				"参数错误"
//	@Failure		403		{object}	response.Response				"无权编辑"
//	@Failure		404		{object}	response.Response				"套图不存在"
//	@Security		BearerAuth
//	@Router			/api/photosets/{id} [put]
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
//
//	@Summary		删除套图
//	@Description	软删除套图，仅创建者或管理员可操作
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"套图ID"
//	@Success		200	{object}	response.Response{data=object}	"删除成功"
//	@Failure		400	{object}	response.Response				"无效的套图ID"
//	@Failure		403	{object}	response.Response				"无权删除"
//	@Failure		404	{object}	response.Response				"套图不存在"
//	@Security		BearerAuth
//	@Router			/api/photosets/{id} [delete]
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

// GetTrash 获取回收站列表
//
//	@Summary		获取回收站列表
//	@Description	获取当前用户已删除的套图列表
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]domain.PhotoSet}	"成功"
//	@Failure		401	{object}	response.Response								"未登录"
//	@Security		BearerAuth
//	@Router			/api/photosets/trash [get]
func (h *PhotoSetHandler) GetTrash(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	trash, err := h.service.GetTrashList(uid)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取回收站失败")
		return
	}

	response.Success(c, trash)
}

// Restore 恢复已删除的套图
//
//	@Summary		恢复套图
//	@Description	从回收站恢复已删除的套图
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"套图ID"
//	@Success		200	{object}	response.Response{data=object}	"恢复成功"
//	@Failure		400	{object}	response.Response				"无效的套图ID"
//	@Failure		401	{object}	response.Response				"未登录"
//	@Security		BearerAuth
//	@Router			/api/photosets/{id}/restore [post]
func (h *PhotoSetHandler) Restore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	userID, _ := c.Get("user_id")
	uid, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	if err := h.service.RestorePhotoSet(uint(id), uid); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "恢复成功"})
}

// PermanentDelete 永久删除套图
//
//	@Summary		永久删除套图
//	@Description	从数据库永久删除套图，不可恢复
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"套图ID"
//	@Success		200	{object}	response.Response{data=object}	"永久删除成功"
//	@Failure		400	{object}	response.Response				"无效的套图ID"
//	@Failure		401	{object}	response.Response				"未登录"
//	@Security		BearerAuth
//	@Router			/api/photosets/{id}/permanent [delete]
func (h *PhotoSetHandler) PermanentDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	userID, _ := c.Get("user_id")
	uid, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	if err := h.service.PermanentDeletePhotoSet(uint(id), uid); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "永久删除成功"})
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
//
//	@Summary		高级搜索套图
//	@Description	高级搜索套图列表，支持分类、价格、排序等多维度筛选（可选鉴权）
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"页码，默认1"
//	@Param			page_size	query	int		false	"每页数量，默认20，最大100"
//	@Param			tag			query	string	false	"标签名称"
//	@Param			mine		query	bool	false	"仅查看自己的作品"
//	@Param			keyword		query	string	false	"搜索关键词"
//	@Param			status		query	string	false	"状态筛选(draft/published/pending)"
//	@Param			category	query	string	false	"分类slug"
//	@Param			price_min	query	number	false	"最低价格"
//	@Param			price_max	query	number	false	"最高价格"
//	@Param			is_free		query	bool	false	"是否免费"
//	@Param			sort_by		query	string	false	"排序方式"
//	@Param			time_range	query	string	false	"时间范围"
//	@Param			user_id		query	int		false	"用户ID筛选"
//	@Success		200			{object}	response.Response{data=object}	"成功"
//	@Failure		400			{object}	response.Response				"参数错误"
//	@Router			/api/photosets/advanced [get]
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
		req.Status,
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
//
//	@Summary		创建套图
//	@Description	创建新套图，需要登录（creator/admin角色）
//	@Tags			PhotoSet
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateRequest	true	"创建套图请求体"
//	@Success		200		{object}	response.Response{data=domain.PhotoSet}	"创建成功"
//	@Failure		400		{object}	response.Response							"参数错误"
//	@Failure		401		{object}	response.Response							"未登录"
//	@Security		BearerAuth
//	@Router			/api/photosets [post]
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

