package domain

import "time"

// OAuthClient 第三方应用
type OAuthClient struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:100;not null"`          // 应用名称
	Description  string    `json:"description" gorm:"type:text"`           // 应用描述
	ClientID     string    `json:"client_id" gorm:"size:64;uniqueIndex;not null"` // 客户端ID
	ClientSecret string    `json:"-" gorm:"size:128;not null"`             // 客户端密钥（不返回给前端）
	RedirectURIs string    `json:"redirect_uris" gorm:"type:text;not null"` // 重定向URI（JSON数组）
	Scopes       string    `json:"scopes" gorm:"size:500;not null"`        // 权限范围（逗号分隔）
	LogoURL      string    `json:"logo_url" gorm:"size:500"`               // 应用Logo
	Status       int       `json:"status" gorm:"default:1"`                // 1=启用, 0=禁用
	CreatedBy    uint      `json:"created_by"`                             // 创建者ID
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (OAuthClient) TableName() string {
	return "oauth_clients"
}