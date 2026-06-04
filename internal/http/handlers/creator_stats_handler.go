package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/http/middleware"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// CreatorStatsHandler 创作者数据统计处理器
type CreatorStatsHandler struct {
	statsService *service.CreatorStatsService
}

// NewCreatorStatsHandler 创建创作者数据统计处理器
func NewCreatorStatsHandler(statsService *service.CreatorStatsService) *CreatorStatsHandler {
	return &CreatorStatsHandler{statsService: statsService}
}

// GetCreatorStats 获取创作者统计数据
func (h *CreatorStatsHandler) GetCreatorStats(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	stats, err := h.statsService.GetCreatorStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计数据失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}

// GetDailyStats 获取每日统计数据
func (h *CreatorStatsHandler) GetDailyStats(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	days := 30
	if d := c.Query("days"); d != "" {
		if v, err := strconv.Atoi(d); err == nil && v > 0 && v <= 90 {
			days = v
		}
	}

	stats, err := h.statsService.GetDailyStats(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取每日统计失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}

// GetPhotoSetStats 获取套图统计数据
func (h *CreatorStatsHandler) GetPhotoSetStats(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "请先登录",
		})
		return
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 50 {
			limit = v
		}
	}

	stats, err := h.statsService.GetPhotoSetStats(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取套图统计失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}
