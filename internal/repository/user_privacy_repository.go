package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// UserPrivacyRepository handles user privacy settings data operations
type UserPrivacyRepository struct {
	db *gorm.DB
}

// NewUserPrivacyRepository creates a new UserPrivacyRepository
func NewUserPrivacyRepository(db *gorm.DB) *UserPrivacyRepository {
	return &UserPrivacyRepository{db: db}
}

// GetByUserID gets privacy settings for a user
func (r *UserPrivacyRepository) GetByUserID(userID uint) (*domain.UserPrivacySetting, error) {
	var setting domain.UserPrivacySetting
	err := r.db.Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认设置
			return &domain.UserPrivacySetting{
				UserID:       userID,
				ShowProfile:  true,
				ShowPosts:    true,
				ShowFavorites: true,
				AllowSearch:  true,
				AllowMessage: true,
			}, nil
		}
		return nil, err
	}
	return &setting, nil
}

// CreateOrUpdate creates or updates privacy settings
func (r *UserPrivacyRepository) CreateOrUpdate(setting *domain.UserPrivacySetting) error {
	// 先尝试查找
	var existing domain.UserPrivacySetting
	err := r.db.Where("user_id = ?", setting.UserID).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 不存在则创建
			return r.db.Create(setting).Error
		}
		return err
	}
	// 存在则更新
	return r.db.Model(&existing).Updates(setting).Error
}

// Delete deletes privacy settings for a user
func (r *UserPrivacyRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.UserPrivacySetting{}).Error
}
