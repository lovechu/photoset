package handlers

import (
	"strconv"

	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetNotifications gets notifications for the current user
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	unreadOnly := c.Query("unread_only") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	notifications, total, err := h.notificationService.GetNotifications(userID, page, pageSize, unreadOnly)
	if err != nil {
		response.ServerError(c, "获取通知失败")
		return
	}

	response.Success(c, gin.H{
		"notifications": h.notificationService.FormatNotifications(notifications),
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetUnreadCount gets the count of unread notifications
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		response.ServerError(c, "获取未读数量失败")
		return
	}

	response.Success(c, gin.H{
		"unread_count": count,
	})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的通知ID")
		return
	}

	if err := h.notificationService.MarkAsRead(uint(notificationID), userID); err != nil {
		response.ServerError(c, "标记已读失败")
		return
	}

	response.Success(c, gin.H{
		"message": "已标记为已读",
	})
}

// MarkAllAsRead marks all notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		response.ServerError(c, "标记全部已读失败")
		return
	}

	response.Success(c, gin.H{
		"message": "已全部标记为已读",
	})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的通知ID")
		return
	}

	if err := h.notificationService.DeleteNotification(uint(notificationID), userID); err != nil {
		response.ServerError(c, "删除通知失败")
		return
	}

	response.Success(c, gin.H{
		"message": "通知已删除",
	})
}

// DeleteAllNotifications deletes all notifications
func (h *NotificationHandler) DeleteAllNotifications(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	if err := h.notificationService.DeleteAllNotifications(userID); err != nil {
		response.ServerError(c, "删除所有通知失败")
		return
	}

	response.Success(c, gin.H{
		"message": "已删除所有通知",
	})
}
