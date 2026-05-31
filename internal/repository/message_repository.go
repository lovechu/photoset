package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// MessageRepository handles message data operations
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create creates a new message
func (r *MessageRepository) Create(message *domain.Message) error {
	return r.db.Create(message).Error
}

// FindByID finds a message by ID
func (r *MessageRepository) FindByID(id uint) (*domain.Message, error) {
	var message domain.Message
	err := r.db.Preload("FromUser").Preload("ToUser").
		First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetConversation returns messages between two users
func (r *MessageRepository) GetConversation(userID, otherUserID uint, page, pageSize int) ([]domain.Message, int64, error) {
	var messages []domain.Message
	var total int64

	query := r.db.Model(&domain.Message{}).
		Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
			userID, otherUserID, otherUserID, userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("FromUser").
		Preload("ToUser").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error

	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetConversations returns conversation list for a user
func (r *MessageRepository) GetConversations(userID uint) ([]domain.Conversation, error) {
	var conversations []domain.Conversation

	// Use a single raw SQL query to get conversation partners with their latest message and unread count.
	// The "?" for the subquery is replaced by inlining the subquery SQL directly to avoid
	// GORM parameter binding issues when mixing raw placeholders with GORM query objects.
	subQuerySQL := `SELECT CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END as user_id
		FROM messages WHERE from_user_id = ? OR to_user_id = ? GROUP BY user_id`

	querySQL := `SELECT 
		u.id as user_id,
		u.nickname as username,
		u.avatar,
		m.content as last_message,
		m.created_at as last_msg_time,
		COALESCE(unread.cnt, 0) as unread_count
	FROM users u
	INNER JOIN (` + subQuerySQL + `) AS partners ON partners.user_id = u.id
	INNER JOIN messages m ON m.id = (
		SELECT m2.id FROM messages m2 
		WHERE (m2.from_user_id = ? AND m2.to_user_id = u.id) 
		   OR (m2.from_user_id = u.id AND m2.to_user_id = ?)
		ORDER BY m2.created_at DESC LIMIT 1
	)
	LEFT JOIN (
		SELECT from_user_id, COUNT(*) as cnt 
		FROM messages 
		WHERE to_user_id = ? AND is_read = false AND deleted_at IS NULL
		GROUP BY from_user_id
	) unread ON unread.from_user_id = u.id
	ORDER BY m.created_at DESC`

	// Parameters in order of appearance:
	// subQuery: ?, ?, ? = userID x3
	// inner select: ?, ? = userID x2
	// left join: ? = userID x1
	// Total: 6 parameters
	rows, err := r.db.Raw(querySQL,
		userID, userID, userID, // subQuery params
		userID, userID,          // inner message subquery params
		userID,                  // unread count param
	).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.Conversation
		if err := rows.Scan(&c.UserID, &c.Username, &c.Avatar, &c.LastMessage, &c.LastMsgTime, &c.UnreadCount); err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}

	return conversations, nil
}

// CountUnread returns unread message count for a user
func (r *MessageRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Message{}).
		Where("to_user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// CountUnreadFromUser returns unread message count from a specific user
func (r *MessageRepository) CountUnreadFromUser(userID, fromUserID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Message{}).
		Where("to_user_id = ? AND from_user_id = ? AND is_read = ?", userID, fromUserID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead marks a message as read
func (r *MessageRepository) MarkAsRead(messageID, userID uint) error {
	return r.db.Model(&domain.Message{}).
		Where("id = ? AND to_user_id = ?", messageID, userID).
		Update("is_read", true).Error
}

// MarkConversationAsRead marks all messages from a user as read
func (r *MessageRepository) MarkConversationAsRead(userID, fromUserID uint) error {
	return r.db.Model(&domain.Message{}).
		Where("to_user_id = ? AND from_user_id = ? AND is_read = ?", userID, fromUserID, false).
		Update("is_read", true).Error
}

// Delete deletes a message (soft delete)
func (r *MessageRepository) Delete(messageID, userID uint) error {
	return r.db.Where("(from_user_id = ? OR to_user_id = ?) AND id = ?", userID, userID, messageID).
		Delete(&domain.Message{}).Error
}
