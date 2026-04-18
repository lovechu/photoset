package domain

import "time"

type ApiKey struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`           // 密钥名称
	Key       string    `json:"key" gorm:"size:64;uniqueIndex;not null"` // API Key
	Secret    string    `json:"-" gorm:"size:128;not null"`              // 密钥（不返回给前端）
	Status    int       `json:"status" gorm:"default:1"`               // 1=启用, 0=禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint      `json:"created_by"` // 创建者 ID
	LastUsed  *time.Time `json:"last_used"` // 最后使用时间
}

func (ApiKey) TableName() string {
	return "api_keys"
}
