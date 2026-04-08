package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

// Add 添加收藏（用 FirstOrCreate 实现幂等，重复收藏不报错）
func (r *FavoriteRepository) Add(userID, photosetID uint) error {
	fav := domain.Favorite{UserID: userID, PhotoSetID: photosetID}
	return r.db.Where("user_id = ? AND photoset_id = ?", userID, photosetID).FirstOrCreate(&fav).Error
}

// Remove 取消收藏（不存在也返回 nil）
func (r *FavoriteRepository) Remove(userID, photosetID uint) error {
	return r.db.Where("user_id = ? AND photoset_id = ?", userID, photosetID).Delete(&domain.Favorite{}).Error
}

// IsFavorited 检查是否已收藏
func (r *FavoriteRepository) IsFavorited(userID, photosetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Favorite{}).Where("user_id = ? AND photoset_id = ?", userID, photosetID).Count(&count).Error
	return count > 0, err
}

// List 用户收藏列表（分页，JOIN photosets 返回完整套图信息）
func (r *FavoriteRepository) List(userID uint, page, pageSize int) ([]domain.Favorite, int64, error) {
	var favorites []domain.Favorite
	var total int64

	query := r.db.Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("PhotoSet.User").Preload("PhotoSet.Tags").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&favorites).Error

	return favorites, total, err
}
