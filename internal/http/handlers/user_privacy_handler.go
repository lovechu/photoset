package handlers

import (
	"net/http"
	"photoset/internal/domain"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// UserPrivacyHandler handles user privacy settings HTTP requests
type UserPrivacyHandler struct {
	privacyService *service.UserPrivacyService
}

// NewUserPrivacyHandler creates a new UserPrivacyHandler
func NewUserPrivacyHandler(privacyService *service.UserPrivacyService) *UserPrivacyHandler {
	return &UserPrivacyHandler{privacyService: privacyService}
}

// GetPrivacySettings handles getting user privacy settings
func (h *UserPrivacyHandler) GetPrivacySettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	settings, err := h.privacyService.GetPrivacySettings(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取隐私设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": settings,
	})
}

// UpdatePrivacySettings handles updating user privacy settings
func (h *UserPrivacyHandler) UpdatePrivacySettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	var settings domain.UserPrivacySetting
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	// 确保设置属于当前用户
	settings.UserID = userID.(uint)

	if err := h.privacyService.UpdatePrivacySettings(userID.(uint), &settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新隐私设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "隐私设置更新成功",
	})
}