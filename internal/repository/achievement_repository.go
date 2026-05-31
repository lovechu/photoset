package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// AchievementRepository handles database operations for achievements
type AchievementRepository struct {
	DB *gorm.DB
}

// NewAchievementRepository creates a new AchievementRepository
func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{DB: db}
}

// GetAllAchievements returns all achievements
func (r *AchievementRepository) GetAllAchievements() ([]domain.Achievement, error) {
	var achievements []domain.Achievement
	err := r.DB.Where("is_hidden = ?", false).Order("sort_order ASC").Find(&achievements).Error
	return achievements, err
}

// GetAllAchievementsIncludingHidden returns all achievements including hidden ones
func (r *AchievementRepository) GetAllAchievementsIncludingHidden() ([]domain.Achievement, error) {
	var achievements []domain.Achievement
	err := r.DB.Order("sort_order ASC").Find(&achievements).Error
	return achievements, err
}

// GetAchievementByID returns an achievement by ID
func (r *AchievementRepository) GetAchievementByID(id uint) (*domain.Achievement, error) {
	var achievement domain.Achievement
	err := r.DB.First(&achievement, id).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// GetAchievementByName returns an achievement by name
func (r *AchievementRepository) GetAchievementByName(name string) (*domain.Achievement, error) {
	var achievement domain.Achievement
	err := r.DB.Where("name = ?", name).First(&achievement).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// GetUserAchievements returns achievements unlocked by a user
func (r *AchievementRepository) GetUserAchievements(userID uint) ([]domain.UserAchievement, error) {
	var userAchievements []domain.UserAchievement
	err := r.DB.Where("user_id = ?", userID).
		Preload("Achievement").
		Order("created_at DESC").
		Find(&userAchievements).Error
	return userAchievements, err
}

// HasAchievement checks if user has already unlocked an achievement
func (r *AchievementRepository) HasAchievement(userID uint, achievementID uint) bool {
	var count int64
	r.DB.Model(&domain.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count)
	return count > 0
}

// UnlockAchievement unlocks an achievement for a user
func (r *AchievementRepository) UnlockAchievement(userID uint, achievementID uint) error {
	userAchievement := domain.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
	}
	return r.DB.Create(&userAchievement).Error
}

// CreateAchievement creates a new achievement
func (r *AchievementRepository) CreateAchievement(achievement *domain.Achievement) error {
	return r.DB.Create(achievement).Error
}

// UpdateAchievement updates an achievement
func (r *AchievementRepository) UpdateAchievement(achievement *domain.Achievement) error {
	return r.DB.Save(achievement).Error
}

// DeleteAchievement soft-deletes an achievement
func (r *AchievementRepository) DeleteAchievement(id uint) error {
	return r.DB.Delete(&domain.Achievement{}, id).Error
}

// InitDefaultAchievements initializes default achievements if not exists
func (r *AchievementRepository) InitDefaultAchievements() error {
	var count int64
	r.DB.Model(&domain.Achievement{}).Count(&count)
	if count > 0 {
		return nil // Already initialized
	}

	achievements := domain.DefaultAchievements()
	for i := range achievements {
		if err := r.DB.Create(&achievements[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetUserAchievementCount returns count of achievements unlocked by user
func (r *AchievementRepository) GetUserAchievementCount(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.UserAchievement{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
