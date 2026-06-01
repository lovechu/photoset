package domain

import (
	"time"
)

// UserBlock represents a user blocking another user
type UserBlock struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID    uint `gorm:"not null;uniqueIndex:idx_user_block" json:"user_id"`     // 拉黑者
	BlockedID uint `gorm:"not null;uniqueIndex:idx_user_block;index:idx_blocked_id" json:"blocked_id"` // 被拉黑者

	// 拉黑类型：block(完全拉黑), mute(静音，不看其内容但不通知)
	BlockType string `gorm:"type:varchar(20);not null;default:block" json:"block_type"`

	// Associations
	User    User `gorm:"foreignKey:UserID" json:"-"`
	Blocked User `gorm:"foreignKey:BlockedID" json:"-"`
}

// TableName specifies the table name
func (UserBlock) TableName() string {
	return "user_blocks"
}
