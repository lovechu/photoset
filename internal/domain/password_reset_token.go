package domain

import "time"

// PasswordResetToken 密码重置令牌
type PasswordResetToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint   `gorm:"not null;index" json:"user_id"`
	Email  string `gorm:"type:varchar(100);not null;index" json:"email"`
	Token  string `gorm:"type:varchar(64);uniqueIndex;not null" json:"token"`
	Used   bool   `gorm:"default:false" json:"used"`
	Expire time.Time `gorm:"not null;index" json:"expire"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
