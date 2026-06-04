package handlers

import (
	"net/http"
	"photoset/internal/http/middleware"
	"photoset/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LoginHistoryHandler handles login history HTTP requests
type LoginHistoryHandler struct {
	historyService *service.LoginHistoryService
}

// NewLoginHistoryHandler creates a new LoginHistoryHandler
func NewLoginHistoryHandler(historyService *service.LoginHistoryService) *LoginHistoryHandler {
	return &LoginHistoryHandler{historyService: historyService}
}

// @Summary      登录历史
// @Description  获取当前用户的登录历史记录
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        page      query int  false  "页码"
// @Param        page_size query int  false  "每页数量"
// @Success      200  {object}  object  "登录历史列表"
// @Failure      401  {object}  object  "未登录"
// @Security     BearerAuth
// @Router       /api/user/login-history [get]
// GetLoginHistory handles getting user login history
func (h *LoginHistoryHandler) GetLoginHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	histories, total, err := h.historyService.GetLoginHistory(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取登录历史失败"})
		return
	}

	// Format response
	type HistoryItem struct {
		ID        uint   `json:"id"`
		IP        string `json:"ip"`
		IPLocation string `json:"ip_location"`
		Device    string `json:"device"`
		Browser   string `json:"browser"`
		OS        string `json:"os"`
		LoginType string `json:"login_type"`
		Success   bool   `json:"success"`
		LoginAt   string `json:"login_at"`
	}

	items := make([]HistoryItem, 0, len(histories))
	for _, history := range histories {
		item := HistoryItem{
			ID:         history.ID,
			IP:         history.IP,
			IPLocation: history.IPLocation,
			Device:     history.Device,
			Browser:    history.Browser,
			OS:         history.OS,
			LoginType:  history.LoginType,
			Success:    history.Success,
			LoginAt:    history.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  items,
			"total": total,
		},
	})
}