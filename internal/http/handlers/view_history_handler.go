package handlers

import (
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ViewHistoryHandler struct {
	service *service.ViewHistoryService
}

func NewViewHistoryHandler(service *service.ViewHistoryService) *ViewHistoryHandler {
	return &ViewHistoryHandler{service: service}
}

// Record 记录浏览历史
func (h *ViewHistoryHandler) Record(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		PhotosetID uint    `json:"photoset_id" binding:"required"`
		Title      string  `json:"title" binding:"required"`
		CoverImage string  `json:"cover_image"`
		PhotoCount int     `json:"photo_count"`
		IsVip      bool    `json:"is_vip"`
		Price      float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.Record(userID, req.PhotosetID, req.Title, req.CoverImage, req.PhotoCount, req.IsVip, req.Price); err != nil {
		response.ServerError(c, "记录浏览历史失败")
		return
	}

	response.Success(c, nil)
}

// List 查询浏览历史列表
func (h *ViewHistoryHandler) List(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 支持按日期范围筛选
	var histories interface{}
	var total int64
	var err error

	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	if startStr != "" && endStr != "" {
		start, err1 := time.Parse("2006-01-02", startStr)
		end, err2 := time.Parse("2006-01-02", endStr)
		if err1 == nil && err2 == nil {
			// end 日期加一天，包含当天
			end = end.AddDate(0, 0, 1)
			histories, total, err = h.service.ListByDateRange(userID, start, end, page, pageSize)
		} else {
			histories, total, err = h.service.List(userID, page, pageSize)
		}
	} else {
		histories, total, err = h.service.List(userID, page, pageSize)
	}

	if err != nil {
		response.ServerError(c, "获取浏览历史失败")
		return
	}

	response.Success(c, gin.H{
		"list":  histories,
		"total": total,
	})
}

// Delete 删除单条浏览历史
func (h *ViewHistoryHandler) Delete(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.Delete(userID, uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// ClearAll 清空所有浏览历史
func (h *ViewHistoryHandler) ClearAll(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	if err := h.service.ClearAll(userID); err != nil {
		response.ServerError(c, "清空失败")
		return
	}

	response.Success(c, nil)
}

// BatchDelete 批量删除浏览历史
func (h *ViewHistoryHandler) BatchDelete(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if len(req.IDs) > 50 {
		response.BadRequest(c, "单次最多删除50条")
		return
	}

	if err := h.service.BatchDelete(userID, req.IDs); err != nil {
		response.ServerError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}
