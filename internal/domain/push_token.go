package domain

import (
	"time"

	"gorm.io/gorm"
)

// PushToken 存储用户的推送令牌
type PushToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID   uint   `gorm:"not null;index:idx_user_platform" json:"user_id"`
	Token    string `gorm:"type:varchar(500);not null;uniqueIndex" json:"token"`
	Platform string `gorm:"type:varchar(20);not null;index:idx_user_platform" json:"platform"` // ios, android, web
	DeviceID string `gorm:"type:varchar(100);default:''" json:"device_id"`
	DeviceName string `gorm:"type:varchar(100);default:''" json:"device_name"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (PushToken) TableName() string {
	return "push_tokens"
}