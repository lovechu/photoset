package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostReplyLikeRepository handles database operations for post reply likes
type PostReplyLikeRepository struct {
	DB *gorm.DB
}

// NewPostReplyLikeRepository creates a new PostReplyLikeRepository
func NewPostReplyLikeRepository(db *gorm.DB) *PostReplyLikeRepository {
	return &PostReplyLikeRepository{DB: db}
}

// Create creates a new reply like record
func (r *PostReplyLikeRepository) Create(like *domain.PostReplyLike) error {
	return r.DB.Create(like).Error
}

// Delete deletes a reply like record
func (r *PostReplyLikeRepository) Delete(userID, replyID uint) error {
	return r.DB.Unscoped().Where("user_id = ? AND reply_id = ?", userID, replyID).Delete(&domain.PostReplyLike{}).Error
}

// FindByUserAndReply finds a like by user ID and reply ID
func (r *PostReplyLikeRepository) FindByUserAndReply(userID, replyID uint) (*domain.PostReplyLike, error) {
	var like domain.PostReplyLike
	err := r.DB.Where("user_id = ? AND reply_id = ?", userID, replyID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// Exists checks if a reply like exists
func (r *PostReplyLikeRepository) Exists(userID, replyID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.PostReplyLike{}).Where("user_id = ? AND reply_id = ?", userID, replyID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByReplyID counts likes for a reply
func (r *PostReplyLikeRepository) CountByReplyID(replyID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.PostReplyLike{}).Where("reply_id = ?", replyID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
