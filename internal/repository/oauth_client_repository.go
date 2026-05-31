package repository

import (
	"encoding/json"

	"photoset/internal/database"
	"photoset/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type OAuthClientRepository struct{}

func NewOAuthClientRepository() *OAuthClientRepository {
	return &OAuthClientRepository{}
}

// List 获取所有 OAuth 应用
func (r *OAuthClientRepository) List() ([]domain.OAuthClient, error) {
	var clients []domain.OAuthClient
	err := database.GetMySQL().Order("created_at DESC").Find(&clients).Error
	return clients, err
}

// GetByID 根据 ID 获取
func (r *OAuthClientRepository) GetByID(id uint) (*domain.OAuthClient, error) {
	var client domain.OAuthClient
	err := database.GetMySQL().First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetByClientID 根据 ClientID 获取
func (r *OAuthClientRepository) GetByClientID(clientID string) (*domain.OAuthClient, error) {
	var client domain.OAuthClient
	err := database.GetMySQL().Where("client_id = ? AND status = 1", clientID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// Create 创建 OAuth 应用
func (r *OAuthClientRepository) Create(name, description string, redirectURIs []string, scopes []string, logoURL string, createdBy uint) (*domain.OAuthClient, error) {
	// 生成随机 ClientID 和 ClientSecret
	clientIDStr, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}
	clientSecretStr, err := generateRandomString(64)
	if err != nil {
		return nil, err
	}

	// 对 ClientSecret 进行加密
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecretStr), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 将 redirectURIs 和 scopes 转换为 JSON 字符串
	redirectURIsJSON, err := json.Marshal(redirectURIs)
	if err != nil {
		return nil, err
	}
	scopesStr := ""
	for i, scope := range scopes {
		if i > 0 {
			scopesStr += ","
		}
		scopesStr += scope
	}

	client := &domain.OAuthClient{
		Name:         name,
		Description:  description,
		ClientID:     "pc_" + clientIDStr, // 前缀标识
		ClientSecret: string(hashedSecret),
		RedirectURIs: string(redirectURIsJSON),
		Scopes:       scopesStr,
		LogoURL:      logoURL,
		Status:       1,
		CreatedBy:    createdBy,
	}

	if err := database.GetMySQL().Create(client).Error; err != nil {
		return nil, err
	}

	// 返回时附带明文 ClientSecret（只返回这一次）
	client.ClientSecret = clientSecretStr
	return client, nil
}

// Update 更新 OAuth 应用
func (r *OAuthClientRepository) Update(id uint, name, description string, redirectURIs []string, scopes []string, logoURL string, status int) error {
	// 将 redirectURIs 和 scopes 转换为 JSON 字符串
	redirectURIsJSON, err := json.Marshal(redirectURIs)
	if err != nil {
		return err
	}
	scopesStr := ""
	for i, scope := range scopes {
		if i > 0 {
			scopesStr += ","
		}
		scopesStr += scope
	}

	return database.GetMySQL().Model(&domain.OAuthClient{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":          name,
		"description":   description,
		"redirect_uris": string(redirectURIsJSON),
		"scopes":        scopesStr,
		"logo_url":      logoURL,
		"status":        status,
	}).Error
}

// Delete 删除 OAuth 应用
func (r *OAuthClientRepository) Delete(id uint) error {
	return database.GetMySQL().Delete(&domain.OAuthClient{}, id).Error
}

// ValidateSecret 验证 ClientSecret
func (r *OAuthClientRepository) ValidateSecret(client *domain.OAuthClient, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(secret))
	return err == nil
}

// AutoMigrate 自动迁移表
func (r *OAuthClientRepository) AutoMigrate() error {
	return database.GetMySQL().AutoMigrate(&domain.OAuthClient{})
}