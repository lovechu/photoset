package handlers

import (
	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
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
		response.Unauthorized(c, "请先登录")
		return
	}

	var req RegisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.pushService.RegisterToken(userID, req.Token, req.Platform, req.DeviceID, req.DeviceName); err != nil {
		response.ServerError(c, "注册推送令牌失败")
		return
	}

	response.SuccessWithMessage(c, "推送令牌注册成功", nil)
}

// UnregisterToken 注销推送令牌
func (h *PushNotificationHandler) UnregisterToken(c *gin.Context) {
	var req UnregisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.pushService.UnregisterToken(req.Token); err != nil {
		response.ServerError(c, "注销推送令牌失败")
		return
	}

	response.SuccessWithMessage(c, "推送令牌已注销", nil)
}
