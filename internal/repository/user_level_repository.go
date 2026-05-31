package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// UserLevelRepository handles database operations for user level configs
type UserLevelRepository struct {
	DB *gorm.DB
}

// NewUserLevelRepository creates a new UserLevelRepository
func NewUserLevelRepository(db *gorm.DB) *UserLevelRepository {
	return &UserLevelRepository{DB: db}
}

// GetAllLevelConfigs returns all level configurations
func (r *UserLevelRepository) GetAllLevelConfigs() ([]domain.UserLevelConfig, error) {
	var configs []domain.UserLevelConfig
	err := r.DB.Order("level ASC").Find(&configs).Error
	return configs, err
}

// GetLevelConfig returns config for a specific level
func (r *UserLevelRepository) GetLevelConfig(level int) (*domain.UserLevelConfig, error) {
	var config domain.UserLevelConfig
	err := r.DB.Where("level = ?", level).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateLevelConfig creates a new level config
func (r *UserLevelRepository) CreateLevelConfig(config *domain.UserLevelConfig) error {
	return r.DB.Create(config).Error
}

// UpdateLevelConfig updates a level config
func (r *UserLevelRepository) UpdateLevelConfig(config *domain.UserLevelConfig) error {
	return r.DB.Save(config).Error
}

// InitDefaultLevelConfigs initializes default level configs if not exists
func (r *UserLevelRepository) InitDefaultLevelConfigs() error {
	var count int64
	r.DB.Model(&domain.UserLevelConfig{}).Count(&count)
	if count > 0 {
		return nil // Already initialized
	}

	configs := domain.DefaultLevelConfigs()
	for i := range configs {
		if err := r.DB.Create(&configs[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
