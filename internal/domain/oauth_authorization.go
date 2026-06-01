package domain

import "time"

// OAuthAuthorization 用户授权记录
type OAuthAuthorization struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"index"`         // 用户ID
	ClientID      string    `json:"client_id" gorm:"size:64;index"` // 客户端ID
	Scopes        string    `json:"scopes" gorm:"size:500"`        // 授权的权限范围
	Code          string    `json:"code" gorm:"size:128;uniqueIndex"` // 授权码
	CodeExpiresAt time.Time `json:"code_expires_at"`               // 授权码过期时间
	CreatedAt     time.Time `json:"created_at"`
}

func (OAuthAuthorization) TableName() string {
	return "oauth_authorizations"
}