package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type UserCollectionRepository struct {
	db *gorm.DB
}

func NewUserCollectionRepository(db *gorm.DB) *UserCollectionRepository {
	return &UserCollectionRepository{db: db}
}

// Create 创建合集
func (r *UserCollectionRepository) Create(collection *domain.UserCollection) error {
	return r.db.Create(collection).Error
}

// Update 更新合集
func (r *UserCollectionRepository) Update(collection *domain.UserCollection) error {
	return r.db.Save(collection).Error
}

// Delete 删除合集（级联删除项目）
func (r *UserCollectionRepository) Delete(userID, collectionID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除合集项目
		if err := tx.Where("collection_id = ?", collectionID).Delete(&domain.CollectionItem{}).Error; err != nil {
			return err
		}
		// 删除合集
		return tx.Where("id = ? AND user_id = ?", collectionID, userID).Delete(&domain.UserCollection{}).Error
	})
}

// GetByID 获取合集详情
func (r *UserCollectionRepository) GetByID(collectionID uint) (*domain.UserCollection, error) {
	var collection domain.UserCollection
	err := r.db.Preload("Items.PhotoSet.User").Preload("Items.PhotoSet.Tags").
		Where("id = ?", collectionID).
		First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// ListByUser 获取用户合集列表
func (r *UserCollectionRepository) ListByUser(userID uint, page, pageSize int) ([]domain.UserCollection, int64, error) {
	var collections []domain.UserCollection
	var total int64

	if err := r.db.Model(&domain.UserCollection{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&collections).Error

	if err != nil {
		return nil, 0, err
	}

	return collections, total, nil
}

// AddItem 添加套图到合集
func (r *UserCollectionRepository) AddItem(collectionID, photosetID uint) error {
	item := domain.CollectionItem{
		CollectionID: collectionID,
		PhotoSetID:   photosetID,
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查是否已存在
		var count int64
		if err := tx.Model(&domain.CollectionItem{}).
			Where("collection_id = ? AND photoset_id = ?", collectionID, photosetID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil // 已存在，幂等处理
		}

		// 添加项目
		if err := tx.Create(&item).Error; err != nil {
			return err
		}

		// 更新计数
		return tx.Model(&domain.UserCollection{}).
			Where("id = ?", collectionID).
			UpdateColumn("item_count", gorm.Expr("item_count + 1")).Error
	})
}

// RemoveItem 从合集移除套图
func (r *UserCollectionRepository) RemoveItem(collectionID, photosetID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("collection_id = ? AND photoset_id = ?", collectionID, photosetID).
			Delete(&domain.CollectionItem{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return tx.Model(&domain.UserCollection{}).
				Where("id = ?", collectionID).
				UpdateColumn("item_count", gorm.Expr("GREATEST(item_count - 1, 0)")).Error
		}
		return nil
	})
}

// BatchAddItems 批量添加套图到合集
func (r *UserCollectionRepository) BatchAddItems(collectionID uint, photosetIDs []uint) error {
	if len(photosetIDs) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		added := 0
		for _, photosetID := range photosetIDs {
			var count int64
			if err := tx.Model(&domain.CollectionItem{}).
				Where("collection_id = ? AND photoset_id = ?", collectionID, photosetID).
				Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				item := domain.CollectionItem{
					CollectionID: collectionID,
					PhotoSetID:   photosetID,
				}
				if err := tx.Create(&item).Error; err != nil {
					return err
				}
				added++
			}
		}
		if added > 0 {
			return tx.Model(&domain.UserCollection{}).
				Where("id = ?", collectionID).
				UpdateColumn("item_count", gorm.Expr("item_count + ?", added)).Error
		}
		return nil
	})
}

// IsInCollection 检查套图是否在合集中
func (r *UserCollectionRepository) IsInCollection(collectionID, photosetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.CollectionItem{}).
		Where("collection_id = ? AND photoset_id = ?", collectionID, photosetID).
		Count(&count).Error
	return count > 0, err
}

// GetUserCollectionsContaining 获取包含指定套图的用户合集
func (r *UserCollectionRepository) GetUserCollectionsContaining(userID, photosetID uint) ([]domain.UserCollection, error) {
	var collections []domain.UserCollection
	err := r.db.Joins("JOIN collection_items ON collection_items.collection_id = user_collections.id").
		Where("user_collections.user_id = ? AND collection_items.photoset_id = ?", userID, photosetID).
		Find(&collections).Error
	return collections, err
}
