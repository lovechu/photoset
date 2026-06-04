package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/service"
)

type PushNotificationHandler struct {
	pushService *service.PushNotificationService
}

func NewPushNotificationHandler(pushService *service.PushNotificationService) *PushNotificationHandler {
	return &PushNotificationHandler{
		pushService: pushService,
	}
}

type RegisterTokenRequest struct {
	Token      string `json:"token" binding:"required"`
	Platform   string `json:"platform" binding:"required,oneof=ios android web"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

type UnregisterTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RegisterToken 注册推送令牌
func (h *PushNotificationHandler) RegisterToken(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请先登录"})
		return
	}

	var req RegisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := h.pushService.RegisterToken(userID, req.Token, req.Platform, req.DeviceID, req.DeviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册推送令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "推送令牌注册成功",
	})
}

// UnregisterToken 注销推送令牌
func (h *PushNotificationHandler) UnregisterToken(c *gin.Context) {
	var req UnregisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := h.pushService.UnregisterToken(req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注销推送令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "推送令牌已注销",
	})
}