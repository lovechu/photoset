package repository

import (
	"time"

	"photoset/internal/domain"

	"gorm.io/gorm"
)

type PushTokenRepository struct {
	db *gorm.DB
}

func NewPushTokenRepository(db *gorm.DB) *PushTokenRepository {
	return &PushTokenRepository{db: db}
}

// Save 保存或更新推送令牌
func (r *PushTokenRepository) Save(token *domain.PushToken) error {
	// 检查是否已存在相同的 token
	var existing domain.PushToken
	err := r.db.Where("token = ?", token.Token).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// 新令牌，直接创建
		return r.db.Create(token).Error
	}
	if err != nil {
		return err
	}
	
	// 已存在，更新信息
	now := time.Now()
	updates := map[string]interface{}{
		"user_id":     token.UserID,
		"platform":    token.Platform,
		"device_id":   token.DeviceID,
		"device_name": token.DeviceName,
		"is_active":   true,
		"last_used_at": &now,
	}
	return r.db.Model(&existing).Updates(updates).Error
}

// GetActiveTokensByUserID 获取用户的所有活跃推送令牌
func (r *PushTokenRepository) GetActiveTokensByUserID(userID uint) ([]domain.PushToken, error) {
	var tokens []domain.PushToken
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&tokens).Error
	return tokens, err
}

// DeactivateToken 停用推送令牌
func (r *PushTokenRepository) DeactivateToken(token string) error {
	return r.db.Model(&domain.PushToken{}).
		Where("token = ?", token).
		Update("is_active", false).Error
}

// DeactivateUserTokens 停用用户的所有推送令牌
func (r *PushTokenRepository) DeactivateUserTokens(userID uint) error {
	return r.db.Model(&domain.PushToken{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
}

// DeleteToken 删除推送令牌
func (r *PushTokenRepository) DeleteToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.PushToken{}).Error
}

// CleanInactiveTokens 清理长时间不活跃的令牌（超过90天）
func (r *PushTokenRepository) CleanInactiveTokens() (int64, error) {
	threshold := time.Now().AddDate(0, -3, 0)
	result := r.db.Where("last_used_at < ? OR (last_used_at IS NULL AND created_at < ?)", threshold, threshold).
		Delete(&domain.PushToken{})
	return result.RowsAffected, result.Error
}