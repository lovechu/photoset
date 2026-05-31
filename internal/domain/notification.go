package domain

import (
	"time"

	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeReply   NotificationType = "reply"
	NotificationTypeFollow  NotificationType = "follow"
	NotificationTypeMention NotificationType = "mention"
	NotificationTypeSystem  NotificationType = "system"
)

// Notification represents a user notification
type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID   uint             `gorm:"not null;index" json:"user_id"`
	Type     NotificationType `gorm:"type:varchar(20);not null;index" json:"type"`
	Title    string           `gorm:"type:varchar(200);not null" json:"title"`
	Content  string           `gorm:"type:text" json:"content"`
	IsRead   bool             `gorm:"not null;default:false;index" json:"is_read"`
	SenderID *uint            `gorm:"index" json:"sender_id"`
	TargetID *uint            `json:"target_id"` // Post ID or Reply ID
	TargetType string         `gorm:"type:varchar(20)" json:"target_type"` // "post" or "reply"

	// Associations
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Sender *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

// TableName specifies the table name
func (Notification) TableName() string {
	return "notifications"
}
