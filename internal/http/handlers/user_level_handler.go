package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/domain"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// UserLevelHandler handles user level, achievements, and points mall requests
type UserLevelHandler struct {
	userLevelService *service.UserLevelService
}

// NewUserLevelHandler creates a new UserLevelHandler
func NewUserLevelHandler(userLevelService *service.UserLevelService) *UserLevelHandler {
	return &UserLevelHandler{
		userLevelService: userLevelService,
	}
}

// ===== Level System =====

// GetUserLevelInfo returns current user's level information
// @Summary Get current user's level information
// @Tags UserLevel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/community/user/level [get]
func (h *UserLevelHandler) GetUserLevelInfo(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "login required")
		return
	}

	info, err := h.userLevelService.GetUserLevelInfo(userID)
	if err != nil {
		response.ServerError(c, "failed to get level info")
		return
	}

	response.Success(c, info)
}

// GetUserLevelInfoByID returns a specific user's level information
// @Summary Get a specific user's level information
// @Tags UserLevel
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/community/users/:id/level [get]
func (h *UserLevelHandler) GetUserLevelInfoByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	info, err := h.userLevelService.GetUserLevelInfo(uint(id))
	if err != nil {
		response.ServerError(c, "failed to get level info")
		return
	}

	response.Success(c, info)
}

// GetAllLevelConfigs returns all level configurations
// @Summary Get all level configurations
// @Tags UserLevel
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/community/levels [get]
func (h *UserLevelHandler) GetAllLevelConfigs(c *gin.Context) {
	configs, err := h.userLevelService.GetAllLevelConfigs()
	if err != nil {
		response.ServerError(c, "failed to get level configs")
		return
	}

	response.Success(c, configs)
}

// ===== Achievement System =====

// GetUserAchievements returns current user's achievements
// @Summary Get current user's achievements
// @Tags UserLevel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/community/user/achievements [get]
func (h *UserLevelHandler) GetUserAchievements(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "login required")
		return
	}

	achievements, err := h.userLevelService.GetUserAchievements(userID)
	if err != nil {
		response.ServerError(c, "failed to get achievements")
		return
	}

	response.Success(c, achievements)
}

// GetUserAchievementsByID returns a specific user's achievements
// @Summary Get a specific user's achievements
// @Tags UserLevel
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/community/users/:id/achievements [get]
func (h *UserLevelHandler) GetUserAchievementsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	achievements, err := h.userLevelService.GetUserAchievements(uint(id))
	if err != nil {
		response.ServerError(c, "failed to get achievements")
		return
	}

	response.Success(c, achievements)
}

// ===== Points Mall =====

// GetPointsMallItems returns all points mall items
// @Summary Get all points mall items
// @Tags UserLevel
// @Accept json
// @Produce json
// @Param category query string false "Category filter"
// @Success 200 {object} response.Response
// @Router /api/community/points-mall/items [get]
func (h *UserLevelHandler) GetPointsMallItems(c *gin.Context) {
	category := c.Query("category")

	items, err := h.userLevelService.GetPointsMallItems(category)
	if err != nil {
		response.ServerError(c, "failed to get items")
		return
	}

	response.Success(c, items)
}

// GetPointsMallCategories returns available categories
// GET /api/community/points-mall/categories
func (h *UserLevelHandler) GetPointsMallCategories(c *gin.Context) {
	categories := h.userLevelService.GetPointsMallCategories()
	response.Success(c, categories)
}

// ExchangeItem exchanges points for an item
// POST /api/community/points-mall/exchange
func (h *UserLevelHandler) ExchangeItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "login required")
		return
	}

	var req struct {
		ItemID uint `json:"item_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	exchange, err := h.userLevelService.ExchangeItem(userID, req.ItemID)
	if err != nil {
		switch err {
		case domain.ErrItemNotFound:
			response.Error(c, http.StatusNotFound, "item not found")
		case domain.ErrItemNotActive:
			response.Error(c, http.StatusBadRequest, "item is not active")
		case domain.ErrOutOfStock:
			response.Error(c, http.StatusBadRequest, "item is out of stock")
		case domain.ErrInsufficientPoints:
			response.Error(c, http.StatusBadRequest, "insufficient points")
		case domain.ErrLevelTooLow:
			response.Error(c, http.StatusBadRequest, "level too low")
		default:
			response.ServerError(c, "failed to exchange item")
		}
		return
	}

	response.Success(c, exchange)
}

// GetUserExchangeHistory returns user's exchange history
// GET /api/community/points-mall/history
func (h *UserLevelHandler) GetUserExchangeHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "login required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	exchanges, total, err := h.userLevelService.GetUserExchangeHistory(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get history")
		return
	}

	response.Success(c, gin.H{
		"exchanges": exchanges,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetPointsLeaderboard returns points leaderboard
// GET /api/community/points/leaderboard
func (h *UserLevelHandler) GetPointsLeaderboard(c *gin.Context) {
	// This would need a repository method to get top users by points
	// For now, return a placeholder
	response.Success(c, []interface{}{})
}
