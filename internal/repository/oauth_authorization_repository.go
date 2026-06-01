package repository

import (
	"time"

	"photoset/internal/database"
	"photoset/internal/domain"
)

type OAuthAuthorizationRepository struct{}

func NewOAuthAuthorizationRepository() *OAuthAuthorizationRepository {
	return &OAuthAuthorizationRepository{}
}

// GetByCode 根据授权码获取
func (r *OAuthAuthorizationRepository) GetByCode(code string) (*domain.OAuthAuthorization, error) {
	var auth domain.OAuthAuthorization
	err := database.GetMySQL().Where("code = ? AND code_expires_at > ?", code, time.Now()).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// Create 创建授权记录
func (r *OAuthAuthorizationRepository) Create(userID uint, clientID, scopes string) (*domain.OAuthAuthorization, error) {
	// 生成随机授权码
	codeStr, err := generateRandomString(64)
	if err != nil {
		return nil, err
	}

	// 授权码有效期 5 分钟
	expiresAt := time.Now().Add(5 * time.Minute)

	auth := &domain.OAuthAuthorization{
		UserID:        userID,
		ClientID:      clientID,
		Scopes:        scopes,
		Code:          "oc_" + codeStr, // 前缀标识
		CodeExpiresAt: expiresAt,
	}

	if err := database.GetMySQL().Create(auth).Error; err != nil {
		return nil, err
	}

	return auth, nil
}

// Delete 删除授权记录（授权码使用后删除）
func (r *OAuthAuthorizationRepository) Delete(id uint) error {
	return database.GetMySQL().Delete(&domain.OAuthAuthorization{}, id).Error
}

// CleanExpired 清理过期的授权码
func (r *OAuthAuthorizationRepository) CleanExpired() error {
	return database.GetMySQL().Where("code_expires_at < ?", time.Now()).Delete(&domain.OAuthAuthorization{}).Error
}

// AutoMigrate 自动迁移表
func (r *OAuthAuthorizationRepository) AutoMigrate() error {
	return database.GetMySQL().AutoMigrate(&domain.OAuthAuthorization{})
}