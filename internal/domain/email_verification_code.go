package domain

import "time"

// EmailVerificationCode 邮箱验证码
type EmailVerificationCode struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Email    string    `gorm:"type:varchar(100);not null;index" json:"email"`
	Code     string    `gorm:"type:varchar(6);not null" json:"code"`
	Used     bool      `gorm:"default:false" json:"used"`
	Expire   time.Time `gorm:"not null;index" json:"expire"`
	Purpose  string    `gorm:"type:varchar(20);default:'bind';comment:verify-注册验证,bind-绑定邮箱" json:"purpose"`
	Attempts int       `gorm:"default:0" json:"attempts"`
}

func (EmailVerificationCode) TableName() string {
	return "email_verification_codes"
}
