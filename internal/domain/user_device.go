package domain

import (
	"time"
)

// UserDevice represents a user's device
type UserDevice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint `gorm:"not null;index" json:"user_id"`

	// 设备信息
	DeviceID   string `gorm:"type:varchar(100);not null" json:"device_id"`   // 设备唯一标识
	DeviceName string `gorm:"type:varchar(100)" json:"device_name"`          // 设备名称
	DeviceType string `gorm:"type:varchar(20);not null" json:"device_type"`  // 设备类型：ios, android, web
	OS         string `gorm:"type:varchar(50)" json:"os"`                    // 操作系统
	Browser    string `gorm:"type:varchar(50)" json:"browser"`               // 浏览器

	// 登录信息
	IP         string    `gorm:"type:varchar(50)" json:"ip"`
	IPLocation string    `gorm:"type:varchar(100)" json:"ip_location"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	LastActiveAt time.Time `json:"last_active_at"`                           // 最后活跃时间

	// 状态
	IsActive bool `gorm:"not null;default:true" json:"is_active"` // 是否活跃

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name
func (UserDevice) TableName() string {
	return "user_devices"
}
