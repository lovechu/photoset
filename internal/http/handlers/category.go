package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.PhotoSetService
}

func NewCategoryHandler(svc *service.PhotoSetService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

// List 公开的分类列表接口（供高级搜索下拉框使用，无需登录）
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类列表失败")
		return
	}
	response.Success(c, categories)
}

// AdminList 管理后台分类列表（带分页和搜索）
func (h *CategoryHandler) AdminList(c *gin.Context) {
	var req struct {
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
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

	categories, err := h.service.GetAllCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	// 如果有关键词过滤（内存中简单过滤）
	if req.Keyword != "" {
		var filtered []domain.Category
		for _, cat := range categories {
			if containsStr(cat.Name, req.Keyword) || containsStr(cat.Slug, req.Keyword) {
				filtered = append(filtered, cat)
			}
		}
		categories = filtered
	}

	total := len(categories)
	// 手动分页
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if start >= total {
		categories = []domain.Category{}
	} else if end > total {
		categories = categories[start:]
	} else {
		categories = categories[start:end]
	}

	response.Success(c, gin.H{
		"list":      categories,
		"total":     int64(total),
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// Create 创建分类
func (h *CategoryHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,max=50"`
		Slug        string `json:"slug" binding:"required,max=50"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	cat := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	if err := h.service.CreateCategory(cat); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建分类失败: "+err.Error())
		return
	}

	response.Success(c, cat)
}

// Update 更新分类
func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"omitempty,max=50"`
		Slug        string `json:"slug" binding:"omitempty,max=50"`
		Description string `json:"description"`
		SortOrder   *int   `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Slug != "" {
		updates["slug"] = req.Slug
	}
	if req.Description != "" {
		updates["description"] = req.Description
	} else {
		updates["description"] = nil
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) == 0 {
		response.Error(c, http.StatusBadRequest, "没有要更新的字段")
		return
	}

	if err := h.service.UpdateCategory(uint(id), updates); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新分类失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// Delete 删除分类
func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除分类失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// 辅助函数：字符串包含（大小写不敏感）
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(containsAtIndex(s, substr) || containsAtIndex(lowerStr(s), lowerStr(substr)))))
}

func containsAtIndex(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func lowerStr(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result[i] = c
	}
	return string(result)
}
