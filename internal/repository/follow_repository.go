package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// FollowRepository handles database operations for user follows
type FollowRepository struct {
	DB *gorm.DB
}

// NewFollowRepository creates a new FollowRepository
func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

// Create creates a new follow record
func (r *FollowRepository) Create(follow *domain.Follow) error {
	return r.DB.Create(follow).Error
}

// Delete deletes a follow record
func (r *FollowRepository) Delete(userID, followingID uint) error {
	return r.DB.Unscoped().Where("user_id = ? AND following_id = ?", userID, followingID).Delete(&domain.Follow{}).Error
}

// Exists checks if a follow relationship exists
func (r *FollowRepository) Exists(userID, followingID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.Follow{}).Where("user_id = ? AND following_id = ?", userID, followingID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountFollowing counts users that a user is following
func (r *FollowRepository) CountFollowing(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Follow{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountFollowers counts followers of a user
func (r *FollowRepository) CountFollowers(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindFollowing returns list of users that a user is following
func (r *FollowRepository) FindFollowing(userID uint, page, pageSize int) ([]domain.User, error) {
	var users []domain.User
	offset := (page - 1) * pageSize
	err := r.DB.
		Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.user_id = ?", userID).
		Select("users.*").
		Order("follows.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindFollowers returns list of users following a user
func (r *FollowRepository) FindFollowers(userID uint, page, pageSize int) ([]domain.User, error) {
	var users []domain.User
	offset := (page - 1) * pageSize
	err := r.DB.
		Joins("JOIN follows ON follows.user_id = users.id").
		Where("follows.following_id = ?", userID).
		Select("users.*").
		Order("follows.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// BatchCheckFollowing checks follow status for multiple users
func (r *FollowRepository) BatchCheckFollowing(userID uint, targetIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if len(targetIDs) == 0 {
		return result, nil
	}

	var follows []domain.Follow
	err := r.DB.Where("user_id = ? AND following_id IN ?", userID, targetIDs).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	for _, f := range follows {
		result[f.FollowingID] = true
	}
	return result, nil
}
