package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// LoginHistoryRepository handles login history data operations
type LoginHistoryRepository struct {
	db *gorm.DB
}

// NewLoginHistoryRepository creates a new LoginHistoryRepository
func NewLoginHistoryRepository(db *gorm.DB) *LoginHistoryRepository {
	return &LoginHistoryRepository{db: db}
}

// Create adds a new login history record
func (r *LoginHistoryRepository) Create(history *domain.LoginHistory) error {
	return r.db.Create(history).Error
}

// GetByUserID gets login history for a user
func (r *LoginHistoryRepository) GetByUserID(userID uint, page, pageSize int) ([]domain.LoginHistory, int64, error) {
	var histories []domain.LoginHistory
	var total int64

	query := r.db.Model(&domain.LoginHistory{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// DeleteOldRecords deletes login history records older than days
func (r *LoginHistoryRepository) DeleteOldRecords(days int) error {
	return r.db.Where("created_at < NOW() - INTERVAL ? DAY", days).
		Delete(&domain.LoginHistory{}).Error
}
