package handlers

import (
	"strconv"

	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// MessageHandler handles message-related requests
type MessageHandler struct {
	messageService *service.MessageService
	hub            *service.Hub
}

// NewMessageHandler creates a new MessageHandler
func NewMessageHandler(messageService *service.MessageService, hub *service.Hub) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		hub:            hub,
	}
}

// SendMessage sends a private message
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req struct {
		ToUserID uint   `json:"to_user_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	message, err := h.messageService.SendMessage(userID, req.ToUserID, req.Content)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// Send real-time WebSocket notification to recipient
	if h.hub != nil {
		h.hub.SendToUser(req.ToUserID, service.WSTypeMessage, map[string]interface{}{
			"message": h.messageService.FormatMessage(*message),
		})
		// Also send unread count update
		unreadCount, _ := h.messageService.GetUnreadCount(req.ToUserID)
		h.hub.SendToUser(req.ToUserID, service.WSTypeUnreadCount, map[string]interface{}{
			"message_unread": unreadCount,
		})
	}

	response.Success(c, gin.H{
		"message": h.messageService.FormatMessage(*message),
	})
}

// GetConversations returns conversation list
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	conversations, err := h.messageService.GetConversations(userID)
	if err != nil {
		response.ServerError(c, "获取会话列表失败")
		return
	}

	response.Success(c, gin.H{
		"conversations": h.messageService.FormatConversations(conversations),
	})
}

// GetConversation returns messages with a specific user
func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	otherUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	messages, total, err := h.messageService.GetConversation(userID, uint(otherUserID), page, pageSize)
	if err != nil {
		response.ServerError(c, "获取消息失败")
		return
	}

	response.Success(c, gin.H{
		"messages": h.messageService.FormatMessages(messages),
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetUnreadCount returns unread message count
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	count, err := h.messageService.GetUnreadCount(userID)
	if err != nil {
		response.ServerError(c, "获取未读数量失败")
		return
	}

	response.Success(c, gin.H{
		"unread_count": count,
	})
}

// MarkAsRead marks a message as read
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的消息ID")
		return
	}

	if err := h.messageService.MarkAsRead(uint(messageID), userID); err != nil {
		response.ServerError(c, "标记已读失败")
		return
	}

	response.Success(c, gin.H{
		"message": "已标记为已读",
	})
}

// MarkConversationAsRead marks all messages from a user as read
func (h *MessageHandler) MarkConversationAsRead(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	fromUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.messageService.MarkConversationAsRead(userID, uint(fromUserID)); err != nil {
		response.ServerError(c, "标记已读失败")
		return
	}

	response.Success(c, gin.H{
		"message": "已标记为已读",
	})
}

// DeleteMessage deletes a message
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的消息ID")
		return
	}

	if err := h.messageService.DeleteMessage(uint(messageID), userID); err != nil {
		response.ServerError(c, "删除消息失败")
		return
	}

	response.Success(c, gin.H{
		"message": "消息已删除",
	})
}
