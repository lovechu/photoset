package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 发表评论
func (r *CommentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

// GetByPhotosetID 获取套图的评论列表（分页，支持当前用户点赞状态）
func (r *CommentRepository) GetByPhotosetID(photosetID uint, userID uint, page, pageSize int) ([]domain.Comment, int64, error) {
	var comments []domain.Comment
	var total int64

	// 只统计顶级评论数量
	if err := r.db.Model(&domain.Comment{}).
		Where("photoset_id = ? AND parent_id IS NULL", photosetID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("photoset_id = ? AND parent_id IS NULL", photosetID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	// 加载每个顶级评论的回复（最多显示前3条）
	for i := range comments {
		var replies []domain.Comment
		r.db.Where("parent_id = ?", comments[i].ID).
			Preload("User").
			Order("created_at ASC").
			Limit(3).
			Find(&replies)
		// 用关联字段传递回复
		comments[i].Parent = &domain.Comment{} // marker
		// 回复通过 API 单独加载
		_ = replies
	}

	return comments, total, nil
}

// GetReplies 获取评论的回复列表
func (r *CommentRepository) GetReplies(parentID uint, userID uint, page, pageSize int) ([]domain.Comment, int64, error) {
	var replies []domain.Comment
	var total int64

	if err := r.db.Model(&domain.Comment{}).
		Where("parent_id = ?", parentID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("parent_id = ?", parentID).
		Preload("User").
		Order("created_at ASC").
		Offset(offset).Limit(pageSize).
		Find(&replies).Error

	return replies, total, err
}

// GetByID 获取单个评论
func (r *CommentRepository) GetByID(id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	return &comment, err
}

// Delete 删除评论（软删除，只有作者本人或管理员可以删）
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Comment{}, id).Error
}

// ToggleLike 点赞/取消点赞评论
func (r *CommentRepository) ToggleLike(commentID, userID uint) (bool, error) {
	var like domain.CommentLike
	err := r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&like).Error

	if err == gorm.ErrRecordNotFound {
		// 点赞
		like = domain.CommentLike{CommentID: commentID, UserID: userID}
		if err := r.db.Create(&like).Error; err != nil {
			return false, err
		}
		// 更新评论点赞数
		r.db.Model(&domain.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1"))
		return true, nil
	}

	if err != nil {
		return false, err
	}

	// 取消点赞
	r.db.Delete(&like)
	r.db.Model(&domain.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	return false, nil
}

// IsLiked 检查用户是否已点赞
func (r *CommentRepository) IsLiked(commentID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetCommentCount 获取套图的评论总数
func (r *CommentRepository) GetCommentCount(photosetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Comment{}).
		Where("photoset_id = ?", photosetID).
		Count(&count).Error
	return count, err
}

// GetLikedCommentIDs 获取用户已点赞的评论ID列表
func (r *CommentRepository) GetLikedCommentIDs(userID uint, commentIDs []uint) (map[uint]bool, error) {
	likedMap := make(map[uint]bool)
	if len(commentIDs) == 0 || userID == 0 {
		return likedMap, nil
	}

	var likes []domain.CommentLike
	err := r.db.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).
		Find(&likes).Error
	if err != nil {
		return nil, err
	}

	for _, like := range likes {
		likedMap[like.CommentID] = true
	}
	return likedMap, nil
}
