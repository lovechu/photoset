package domain

import "time"

// OAuthToken 访问令牌
type OAuthToken struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	UserID       uint       `json:"user_id" gorm:"index"`         // 用户ID
	ClientID     string     `json:"client_id" gorm:"size:64;index"` // 客户端ID
	AccessToken  string     `json:"access_token" gorm:"size:256;uniqueIndex"` // 访问令牌
	RefreshToken string     `json:"refresh_token" gorm:"size:256"`  // 刷新令牌
	Scopes       string     `json:"scopes" gorm:"size:500"`        // 权限范围
	ExpiresAt    time.Time  `json:"expires_at"`                    // 过期时间
	Revoked      bool       `json:"revoked" gorm:"default:false"`  // 是否已撤销
	CreatedAt    time.Time  `json:"created_at"`
}

func (OAuthToken) TableName() string {
	return "oauth_tokens"
}