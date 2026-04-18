package repository

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"photoset/internal/database"
	"photoset/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ApiKeyRepository struct{}

func NewApiKeyRepository() *ApiKeyRepository {
	return &ApiKeyRepository{}
}

// List 获取所有 API 密钥
func (r *ApiKeyRepository) List() ([]domain.ApiKey, error) {
	var keys []domain.ApiKey
	err := database.GetMySQL().Order("created_at DESC").Find(&keys).Error
	return keys, err
}

// GetByID 根据 ID 获取
func (r *ApiKeyRepository) GetByID(id uint) (*domain.ApiKey, error) {
	var key domain.ApiKey
	err := database.GetMySQL().First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetByKey 根据 key 获取
func (r *ApiKeyRepository) GetByKey(key string) (*domain.ApiKey, error) {
	var apiKey domain.ApiKey
	err := database.GetMySQL().Where("`key` = ? AND status = 1", key).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// Create 创建 API 密钥
func (r *ApiKeyRepository) Create(name string, createdBy uint) (*domain.ApiKey, error) {
	// 生成随机 Key 和 Secret
	keyStr, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}
	secretStr, err := generateRandomString(64)
	if err != nil {
		return nil, err
	}

	// 存储时对 secret 进行加密
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secretStr), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	apiKey := &domain.ApiKey{
		Name:      name,
		Key:       "psk_" + keyStr, // 前缀标识
		Secret:    string(hashedSecret),
		Status:    1,
		CreatedBy: createdBy,
	}

	if err := database.GetMySQL().Create(apiKey).Error; err != nil {
		return nil, err
	}

	// 返回时附带明文 secret（只返回这一次）
	apiKey.Secret = secretStr
	return apiKey, nil
}

// Delete 删除 API 密钥
func (r *ApiKeyRepository) Delete(id uint) error {
	return database.GetMySQL().Delete(&domain.ApiKey{}, id).Error
}

// UpdateLastUsed 更新最后使用时间
func (r *ApiKeyRepository) UpdateLastUsed(id uint) error {
	now := time.Now()
	return database.GetMySQL().Model(&domain.ApiKey{}).Where("id = ?", id).Update("last_used", now).Error
}

// ValidateSecret 验证 secret
func (r *ApiKeyRepository) ValidateSecret(apiKey *domain.ApiKey, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(apiKey.Secret), []byte(secret))
	return err == nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AutoMigrate 自动迁移表
func (r *ApiKeyRepository) AutoMigrate() error {
	return database.GetMySQL().AutoMigrate(&domain.ApiKey{})
}

// GetDB 获取数据库连接，用于手动迁移
func GetDB() *gorm.DB {
	return database.GetMySQL()
}
