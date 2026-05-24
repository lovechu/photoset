package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostReplyRepository handles database operations for post replies
type PostReplyRepository struct {
	DB *gorm.DB
}

// NewPostReplyRepository creates a new PostReplyRepository
func NewPostReplyRepository(db *gorm.DB) *PostReplyRepository {
	return &PostReplyRepository{DB: db}
}

// Create creates a new reply
func (r *PostReplyRepository) Create(reply *domain.PostReply) error {
	return r.DB.Create(reply).Error
}

// FindByID finds a reply by ID
func (r *PostReplyRepository) FindByID(id uint) (*domain.PostReply, error) {
	var reply domain.PostReply
	err := r.DB.Preload("User").First(&reply, id).Error
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// FindByPostID finds replies by post ID (top-level only)
func (r *PostReplyRepository) FindByPostID(postID uint, page, pageSize int) ([]domain.PostReply, error) {
	var replies []domain.PostReply

	query := r.DB.Model(&domain.PostReply{}).
		Where("post_id = ? AND parent_reply_id IS NULL", postID).
		Preload("User").
		Preload("Children.User").
		Order("created_at ASC")

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&replies).Error
	if err != nil {
		return nil, err
	}

	return replies, nil
}

// FindChildren finds child replies (nested replies)
func (r *PostReplyRepository) FindChildren(parentReplyID uint) ([]domain.PostReply, error) {
	var replies []domain.PostReply
	err := r.DB.Where("parent_reply_id = ?", parentReplyID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error
	if err != nil {
		return nil, err
	}
	return replies, nil
}

// FindByUserID finds replies by user ID
func (r *PostReplyRepository) FindByUserID(userID uint, page, pageSize int) ([]domain.PostReply, int64, error) {
	var replies []domain.PostReply
	var total int64

	query := r.DB.Model(&domain.PostReply{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Post").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&replies).Error
	if err != nil {
		return nil, 0, err
	}

	return replies, total, nil
}

// Update updates a reply
func (r *PostReplyRepository) Update(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.PostReply{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a reply (hard delete)
func (r *PostReplyRepository) Delete(id uint) error {
	return r.DB.Unscoped().Delete(&domain.PostReply{}, id).Error
}

// IncrementLikeCount increments the like count
func (r *PostReplyRepository) IncrementLikeCount(id uint) error {
	return r.DB.Model(&domain.PostReply{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount decrements the like count
func (r *PostReplyRepository) DecrementLikeCount(id uint) error {
	return r.DB.Model(&domain.PostReply{}).Where("id = ?", id).Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// DeleteByPostID deletes all replies for a post
func (r *PostReplyRepository) DeleteByPostID(postID uint) error {
	return r.DB.Unscoped().Where("post_id = ?", postID).Delete(&domain.PostReply{}).Error
}

// CountByPostID counts replies for a post
func (r *PostReplyRepository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.PostReply{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}
