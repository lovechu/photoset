package domain

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a private message between users
type Message struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	FromUserID uint           `gorm:"not null;index" json:"from_user_id"`
	ToUserID   uint           `gorm:"not null;index" json:"to_user_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	IsRead     bool           `gorm:"not null;default:false;index" json:"is_read"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	FromUser *User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   *User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
}

// TableName returns the table name for Message
func (Message) TableName() string {
	return "messages"
}

// Conversation represents a conversation summary between two users
type Conversation struct {
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Avatar       string    `json:"avatar"`
	LastMessage  string    `json:"last_message"`
	LastMsgTime  time.Time `json:"last_msg_time"`
	UnreadCount  int64     `json:"unread_count"`
}
