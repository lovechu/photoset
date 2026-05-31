package service

import (
	"errors"
	"photoset/internal/domain"
	"photoset/internal/repository"
)

// MessageService provides message business logic
type MessageService struct {
	messageRepo *repository.MessageRepository
	userRepo    repository.UserRepository
}

// NewMessageService creates a new MessageService
func NewMessageService(messageRepo *repository.MessageRepository, userRepo repository.UserRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

// SendMessage sends a private message
func (s *MessageService) SendMessage(fromUserID, toUserID uint, content string) (*domain.Message, error) {
	if fromUserID == toUserID {
		return nil, errors.New("不能给自己发送消息")
	}

	if content == "" {
		return nil, errors.New("消息内容不能为空")
	}

	// Verify target user exists
	_, err := s.userRepo.FindByID(toUserID)
	if err != nil {
		return nil, errors.New("目标用户不存在")
	}

	message := &domain.Message{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		IsRead:     false,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	// Reload with user info
	return s.messageRepo.FindByID(message.ID)
}

// GetConversation returns messages between two users
func (s *MessageService) GetConversation(userID, otherUserID uint, page, pageSize int) ([]domain.Message, int64, error) {
	return s.messageRepo.GetConversation(userID, otherUserID, page, pageSize)
}

// GetConversations returns conversation list for a user
func (s *MessageService) GetConversations(userID uint) ([]domain.Conversation, error) {
	return s.messageRepo.GetConversations(userID)
}

// GetUnreadCount returns unread message count
func (s *MessageService) GetUnreadCount(userID uint) (int64, error) {
	return s.messageRepo.CountUnread(userID)
}

// MarkAsRead marks a message as read
func (s *MessageService) MarkAsRead(messageID, userID uint) error {
	return s.messageRepo.MarkAsRead(messageID, userID)
}

// MarkConversationAsRead marks all messages from a user as read
func (s *MessageService) MarkConversationAsRead(userID, fromUserID uint) error {
	return s.messageRepo.MarkConversationAsRead(userID, fromUserID)
}

// DeleteMessage deletes a message
func (s *MessageService) DeleteMessage(messageID, userID uint) error {
	return s.messageRepo.Delete(messageID, userID)
}

// FormatMessage formats a message for display
func (s *MessageService) FormatMessage(m domain.Message) map[string]interface{} {
	result := map[string]interface{}{
		"id":         m.ID,
		"from_user_id": m.FromUserID,
		"to_user_id": m.ToUserID,
		"content":    m.Content,
		"is_read":    m.IsRead,
		"created_at": m.CreatedAt,
	}

	if m.FromUser != nil {
		result["from_user"] = map[string]interface{}{
			"id":       m.FromUser.ID,
			"username": m.FromUser.Nickname,
			"avatar":   m.FromUser.Avatar,
		}
	}

	if m.ToUser != nil {
		result["to_user"] = map[string]interface{}{
			"id":       m.ToUser.ID,
			"username": m.ToUser.Nickname,
			"avatar":   m.ToUser.Avatar,
		}
	}

	return result
}

// FormatMessages formats a list of messages
func (s *MessageService) FormatMessages(messages []domain.Message) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, m := range messages {
		result[i] = s.FormatMessage(m)
	}
	return result
}

// FormatConversations formats a conversation list
func (s *MessageService) FormatConversations(conversations []domain.Conversation) []map[string]interface{} {
	result := make([]map[string]interface{}, len(conversations))
	for i, c := range conversations {
		result[i] = map[string]interface{}{
			"user_id":       c.UserID,
			"username":      c.Username,
			"avatar":        c.Avatar,
			"last_message":  c.LastMessage,
			"last_msg_time": c.LastMsgTime,
			"unread_count":  c.UnreadCount,
		}
	}
	return result
}
