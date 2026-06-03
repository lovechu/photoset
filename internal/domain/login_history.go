package domain

import (
	"time"
)

// LoginHistory represents a login history record
type LoginHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID uint `gorm:"not null;index" json:"user_id"`

	// 登录信息
	IP        string `gorm:"type:varchar(50);not null" json:"ip"`
	IPLocation string `gorm:"type:varchar(100)" json:"ip_location"`
	UserAgent string `gorm:"type:text" json:"user_agent"`
	Device    string `gorm:"type:varchar(100)" json:"device"`     // 设备类型：iOS, Android, Web
	Browser   string `gorm:"type:varchar(100)" json:"browser"`   // 浏览器类型
	OS        string `gorm:"type:varchar(100)" json:"os"`        // 操作系统

	// 登录结果
	LoginType string `gorm:"type:varchar(20);not null;default:password" json:"login_type"` // 登录方式：password, oauth, sms
	Success   bool   `gorm:"not null;default:true" json:"success"`                          // 是否成功
	FailReason string `gorm:"type:varchar(200)" json:"fail_reason"`                        // 失败原因

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name
func (LoginHistory) TableName() string {
	return "login_history"
}
