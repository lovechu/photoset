package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// UserBlockRepository handles user block data operations
type UserBlockRepository struct {
	db *gorm.DB
}

// NewUserBlockRepository creates a new UserBlockRepository
func NewUserBlockRepository(db *gorm.DB) *UserBlockRepository {
	return &UserBlockRepository{db: db}
}

// Create adds a new block record
func (r *UserBlockRepository) Create(block *domain.UserBlock) error {
	return r.db.Create(block).Error
}

// Delete removes a block record
func (r *UserBlockRepository) Delete(userID, blockedID uint) error {
	return r.db.Where("user_id = ? AND blocked_id = ?", userID, blockedID).
		Delete(&domain.UserBlock{}).Error
}

// FindByUserAndBlocked finds a block record
func (r *UserBlockRepository) FindByUserAndBlocked(userID, blockedID uint) (*domain.UserBlock, error) {
	var block domain.UserBlock
	err := r.db.Where("user_id = ? AND blocked_id = ?", userID, blockedID).
		First(&block).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

// GetBlockedUsers gets list of blocked users
func (r *UserBlockRepository) GetBlockedUsers(userID uint, page, pageSize int) ([]domain.UserBlock, int64, error) {
	var blocks []domain.UserBlock
	var total int64

	query := r.db.Model(&domain.UserBlock{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Blocked").
		Order("created_at DESC").
		Find(&blocks).Error; err != nil {
		return nil, 0, err
	}

	return blocks, total, nil
}

// GetBlockedUserIDs gets list of blocked user IDs for filtering
func (r *UserBlockRepository) GetBlockedUserIDs(userID uint) ([]uint, error) {
	var ids []uint
	if err := r.db.Model(&domain.UserBlock{}).
		Where("user_id = ?", userID).
		Pluck("blocked_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// IsBlocked checks if a user is blocked by another user
func (r *UserBlockRepository) IsBlocked(userID, targetID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.UserBlock{}).
		Where("user_id = ? AND blocked_id = ?", userID, targetID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsMuted checks if a user is muted by another user
func (r *UserBlockRepository) IsMuted(userID, targetID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.UserBlock{}).
		Where("user_id = ? AND blocked_id = ? AND block_type = ?", userID, targetID, "mute").
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// AutoMigrate creates or updates the table
func (r *UserBlockRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&domain.UserBlock{})
}
