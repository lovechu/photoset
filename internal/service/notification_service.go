package service

import (
	"fmt"
	"photoset/internal/domain"
	"photoset/internal/repository"
)

// NotificationService provides notification business logic
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
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

	return s.notificationRepo.Create(notification)
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

	return s.notificationRepo.Create(notification)
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

	return s.notificationRepo.Create(notification)
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

	return s.notificationRepo.Create(notification)
}

// CreateSystemNotification creates a system notification
func (s *NotificationService) CreateSystemNotification(userID uint, title, content string) error {
	notification := &domain.Notification{
		UserID:  userID,
		Type:    domain.NotificationTypeSystem,
		Title:   title,
		Content: content,
	}

	return s.notificationRepo.Create(notification)
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
