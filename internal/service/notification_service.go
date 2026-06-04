package service

import (
	"fmt"
	"photoset/internal/domain"
	"photoset/internal/repository"
)

// NotificationService provides notification business logic
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	hub              *Hub
	pushService      *PushNotificationService
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(notificationRepo *repository.NotificationRepository, hub *Hub) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		hub:              hub,
	}
}

// SetPushService 设置推送通知服务（延迟注入）
func (s *NotificationService) SetPushService(pushService *PushNotificationService) {
	s.pushService = pushService
}

// sendRealTimeNotification sends a WebSocket notification to the user
func (s *NotificationService) sendRealTimeNotification(userID uint, notification *domain.Notification) {
	if s.hub != nil {
		s.hub.SendToUser(userID, WSTypeNotification, map[string]interface{}{
			"id":          notification.ID,
			"type":        notification.Type,
			"title":       notification.Title,
			"content":     notification.Content,
			"sender_id":   notification.SenderID,
			"target_id":   notification.TargetID,
			"target_type": notification.TargetType,
			"created_at":  notification.CreatedAt,
		})
		// Also send updated unread count
		unreadCount, _ := s.notificationRepo.CountUnread(userID)
		s.hub.SendToUser(userID, WSTypeUnreadCount, map[string]interface{}{
			"notification_unread": unreadCount,
		})
	}

	// 发送推送通知（离线用户）
	if s.pushService != nil && s.pushService.IsEnabled() {
		data := map[string]string{
			"notification_id": fmt.Sprintf("%d", notification.ID),
			"type":            string(notification.Type),
		}
		if notification.TargetID != nil {
			data["target_id"] = fmt.Sprintf("%d", *notification.TargetID)
			data["target_type"] = notification.TargetType
		}
		
		// 异步发送推送，不阻塞主流程
		go func() {
			if err := s.pushService.SendPushNotification(userID, notification.Title, notification.Content, data); err != nil {
				// 推送失败不影响主流程，仅记录日志
				_ = err
			}
		}()
	}
}

// CreateLikeNotification creates a notification when someone likes a post
func (s *NotificationService) CreateLikeNotification(userID, senderID, postID uint) error {
	if userID == senderID {
		return nil // Don't notify self
	}

	notification := &domain.Notification{
		UserID:     userID,
		Type:       domain.NotificationTypeLike,
		Title:      "收到新点赞",
		Content:    "有人点赞了你的帖子",
		SenderID:   &senderID,
		TargetID:   &postID,
		TargetType: "post",
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	// Send real-time notification
	s.sendRealTimeNotification(userID, notification)
	return nil
}

// CreateReplyNotification creates a notification when someone replies to a post
func (s *NotificationService) CreateReplyNotification(userID, senderID, postID uint) error {
	if userID == senderID {
		return nil // Don't notify self
	}

	notification := &domain.Notification{
		UserID:     userID,
		Type:       domain.NotificationTypeReply,
		Title:      "收到新回复",
		Content:    "有人回复了你的帖子",
		SenderID:   &senderID,
		TargetID:   &postID,
		TargetType: "post",
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	// Send real-time notification
	s.sendRealTimeNotification(userID, notification)
	return nil
}

// CreateFollowNotification creates a notification when someone follows a user
func (s *NotificationService) CreateFollowNotification(userID, senderID uint) error {
	if userID == senderID {
		return nil // Don't notify self
	}

	notification := &domain.Notification{
		UserID:   userID,
		Type:     domain.NotificationTypeFollow,
		Title:    "收到新关注",
		Content:  "有人关注了你",
		SenderID: &senderID,
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	// Send real-time notification
	s.sendRealTimeNotification(userID, notification)
	return nil
}

// CreateMentionNotification creates a notification when someone mentions a user
func (s *NotificationService) CreateMentionNotification(userID, senderID, postID uint) error {
	if userID == senderID {
		return nil // Don't notify self
	}

	notification := &domain.Notification{
		UserID:     userID,
		Type:       domain.NotificationTypeMention,
		Title:      "被提及",
		Content:    "有人在帖子中提到了你",
		SenderID:   &senderID,
		TargetID:   &postID,
		TargetType: "post",
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	// Send real-time notification
	s.sendRealTimeNotification(userID, notification)
	return nil
}

// CreateSystemNotification creates a system notification
func (s *NotificationService) CreateSystemNotification(userID uint, title, content string) error {
	notification := &domain.Notification{
		UserID:  userID,
		Type:    domain.NotificationTypeSystem,
		Title:   title,
		Content: content,
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	// Send real-time notification
	s.sendRealTimeNotification(userID, notification)
	return nil
}

// GetNotifications returns notifications for a user
func (s *NotificationService) GetNotifications(userID uint, page, pageSize int, unreadOnly bool) ([]domain.Notification, int64, error) {
	return s.notificationRepo.ListByUserID(userID, page, pageSize, unreadOnly)
}

// GetUnreadCount returns the count of unread notifications
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.notificationRepo.CountUnread(userID)
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(notificationID, userID uint) error {
	return s.notificationRepo.MarkAsRead(notificationID, userID)
}

// MarkAllAsRead marks all notifications as read
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(notificationID, userID uint) error {
	return s.notificationRepo.Delete(notificationID, userID)
}

// DeleteAllNotifications deletes all notifications for a user
func (s *NotificationService) DeleteAllNotifications(userID uint) error {
	return s.notificationRepo.DeleteAll(userID)
}

// FormatNotification formats a notification for display
func (s *NotificationService) FormatNotification(n domain.Notification) map[string]interface{} {
	result := map[string]interface{}{
		"id":         n.ID,
		"type":       n.Type,
		"title":      n.Title,
		"content":    n.Content,
		"is_read":    n.IsRead,
		"created_at": n.CreatedAt,
	}

	if n.SenderID != nil {
		result["sender_id"] = *n.SenderID
		if n.Sender != nil {
			result["sender_name"] = n.Sender.Nickname
			result["sender_avatar"] = n.Sender.Avatar
		}
	}

	if n.TargetID != nil {
		result["target_id"] = *n.TargetID
		result["target_type"] = n.TargetType
	}

	return result
}

// FormatNotifications formats a list of notifications
func (s *NotificationService) FormatNotifications(notifications []domain.Notification) []map[string]interface{} {
	result := make([]map[string]interface{}, len(notifications))
	for i, n := range notifications {
		result[i] = s.FormatNotification(n)
	}
	return result
}

// GetNotificationMessage returns a human-readable message for a notification
func GetNotificationMessage(n domain.Notification) string {
	switch n.Type {
	case domain.NotificationTypeLike:
		return "点赞了你的帖子"
	case domain.NotificationTypeReply:
		return "回复了你的帖子"
	case domain.NotificationTypeFollow:
		return "关注了你"
	case domain.NotificationTypeMention:
		return "在帖子中提到了你"
	case domain.NotificationTypeSystem:
		return n.Content
	default:
		return fmt.Sprintf("未知通知类型: %s", n.Type)
	}
}
