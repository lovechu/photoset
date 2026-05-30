package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type PostFavoriteRepository struct {
	db *gorm.DB
}

func NewPostFavoriteRepository(db *gorm.DB) *PostFavoriteRepository {
	return &PostFavoriteRepository{db: db}
}

// Add 添加收藏（用 FirstOrCreate 实现幂等，重复收藏不报错）
func (r *PostFavoriteRepository) Add(userID, postID uint) error {
	fav := domain.PostFavorite{UserID: userID, PostID: postID}
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).FirstOrCreate(&fav).Error
}

// Remove 取消收藏（不存在也返回 nil）
func (r *PostFavoriteRepository) Remove(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&domain.PostFavorite{}).Error
}

// IsFavorited 检查是否已收藏
func (r *PostFavoriteRepository) IsFavorited(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.PostFavorite{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

// Toggle 切换收藏状态（已收藏则取消，未收藏则添加）
func (r *PostFavoriteRepository) Toggle(userID, postID uint) (bool, error) {
	isFavorited, err := r.IsFavorited(userID, postID)
	if err != nil {
		return false, err
	}

	if isFavorited {
		err = r.Remove(userID, postID)
		return false, err
	}

	err = r.Add(userID, postID)
	return true, err
}

// List 用户收藏的帖子列表（分页）
func (r *PostFavoriteRepository) List(userID uint, page, pageSize int) ([]domain.PostFavorite, int64, error) {
	var favorites []domain.PostFavorite
	var total int64

	if err := r.db.Model(&domain.PostFavorite{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Preload("Post.User").
		Preload("Post.Likes").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&favorites).Error

	if err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// BatchCheck 批量检查收藏状态
func (r *PostFavoriteRepository) BatchCheck(userID uint, postIDs []uint) (map[uint]bool, error) {
	if len(postIDs) == 0 {
		return make(map[uint]bool), nil
	}

	var favorites []domain.PostFavorite
	err := r.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]bool)
	for _, postID := range postIDs {
		result[postID] = false
	}
	for _, fav := range favorites {
		result[fav.PostID] = true
	}

	return result, nil
}