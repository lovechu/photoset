package repository

import (
	"time"

	"photoset/internal/database"
	"photoset/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type OAuthTokenRepository struct{}

func NewOAuthTokenRepository() *OAuthTokenRepository {
	return &OAuthTokenRepository{}
}

// GetByAccessToken 根据访问令牌获取
func (r *OAuthTokenRepository) GetByAccessToken(accessToken string) (*domain.OAuthToken, error) {
	var token domain.OAuthToken
	err := database.GetMySQL().Where("access_token = ? AND revoked = false AND expires_at > ?", accessToken, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByRefreshToken 根据刷新令牌获取
func (r *OAuthTokenRepository) GetByRefreshToken(refreshToken string) (*domain.OAuthToken, error) {
	var token domain.OAuthToken
	err := database.GetMySQL().Where("refresh_token = ? AND revoked = false", refreshToken).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Create 创建访问令牌
func (r *OAuthTokenRepository) Create(userID uint, clientID, scopes string) (*domain.OAuthToken, error) {
	// 生成随机 Access Token 和 Refresh Token
	accessTokenStr, err := generateRandomString(64)
	if err != nil {
		return nil, err
	}
	refreshTokenStr, err := generateRandomString(64)
	if err != nil {
		return nil, err
	}

	// 对 Access Token 和 Refresh Token 进行加密存储
	hashedAccessToken, err := bcrypt.GenerateFromPassword([]byte(accessTokenStr), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshTokenStr), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Access Token 有效期 2 小时
	expiresAt := time.Now().Add(2 * time.Hour)

	token := &domain.OAuthToken{
		UserID:       userID,
		ClientID:     clientID,
		AccessToken:  string(hashedAccessToken),
		RefreshToken: string(hashedRefreshToken),
		Scopes:       scopes,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}

	if err := database.GetMySQL().Create(token).Error; err != nil {
		return nil, err
	}

	// 返回时附带明文令牌（只返回这一次）
	token.AccessToken = accessTokenStr
	token.RefreshToken = refreshTokenStr
	return token, nil
}

// Refresh 刷新令牌
func (r *OAuthTokenRepository) Refresh(oldToken *domain.OAuthToken) (*domain.OAuthToken, error) {
	// 删除旧令牌
	if err := r.Revoke(oldToken.ID); err != nil {
		return nil, err
	}

	// 创建新令牌
	return r.Create(oldToken.UserID, oldToken.ClientID, oldToken.Scopes)
}

// Revoke 撤销令牌
func (r *OAuthTokenRepository) Revoke(id uint) error {
	return database.GetMySQL().Model(&domain.OAuthToken{}).Where("id = ?", id).Update("revoked", true).Error
}

// RevokeByClientID 撤销指定应用的所有令牌
func (r *OAuthTokenRepository) RevokeByClientID(clientID string) error {
	return database.GetMySQL().Model(&domain.OAuthToken{}).Where("client_id = ?", clientID).Update("revoked", true).Error
}

// RevokeByUserIDAndClientID 撤销指定用户和应用的所有令牌
func (r *OAuthTokenRepository) RevokeByUserIDAndClientID(userID uint, clientID string) error {
	return database.GetMySQL().Model(&domain.OAuthToken{}).Where("user_id = ? AND client_id = ?", userID, clientID).Update("revoked", true).Error
}

// ValidateAccessToken 验证访问令牌
func (r *OAuthTokenRepository) ValidateAccessToken(token *domain.OAuthToken, accessToken string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(token.AccessToken), []byte(accessToken))
	return err == nil
}

// ValidateRefreshToken 验证刷新令牌
func (r *OAuthTokenRepository) ValidateRefreshToken(token *domain.OAuthToken, refreshToken string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(token.RefreshToken), []byte(refreshToken))
	return err == nil
}

// AutoMigrate 自动迁移表
func (r *OAuthTokenRepository) AutoMigrate() error {
	return database.GetMySQL().AutoMigrate(&domain.OAuthToken{})
}