package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostLikeRepository handles database operations for post likes
type PostLikeRepository struct {
	DB *gorm.DB
}

// NewPostLikeRepository creates a new PostLikeRepository
func NewPostLikeRepository(db *gorm.DB) *PostLikeRepository {
	return &PostLikeRepository{DB: db}
}

// Create creates a new like record
func (r *PostLikeRepository) Create(like *domain.PostLike) error {
	return r.DB.Create(like).Error
}

// Delete deletes a like record
func (r *PostLikeRepository) Delete(userID, postID uint) error {
	return r.DB.Unscoped().Where("user_id = ? AND post_id = ?", userID, postID).Delete(&domain.PostLike{}).Error
}

// FindByUserAndPost finds a like by user ID and post ID
func (r *PostLikeRepository) FindByUserAndPost(userID, postID uint) (*domain.PostLike, error) {
	var like domain.PostLike
	err := r.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// Exists checks if a like exists
func (r *PostLikeRepository) Exists(userID, postID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.PostLike{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByPostID counts likes for a post
func (r *PostLikeRepository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.PostLike{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindByUserID finds likes by user ID
func (r *PostLikeRepository) FindByUserID(userID uint, page, pageSize int) ([]domain.PostLike, error) {
	var likes []domain.PostLike
	offset := (page - 1) * pageSize
	err := r.DB.Where("user_id = ?", userID).
		Preload("Post").
		Preload("Post.User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}
