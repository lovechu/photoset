package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PointsMallRepository handles database operations for points mall
type PointsMallRepository struct {
	DB *gorm.DB
}

// NewPointsMallRepository creates a new PointsMallRepository
func NewPointsMallRepository(db *gorm.DB) *PointsMallRepository {
	return &PointsMallRepository{DB: db}
}

// GetAllItems returns all active points mall items
func (r *PointsMallRepository) GetAllItems() ([]domain.PointsMallItem, error) {
	var items []domain.PointsMallItem
	err := r.DB.Where("is_active = ?", true).Order("sort_order ASC").Find(&items).Error
	return items, err
}

// GetItemsByCategory returns items by category
func (r *PointsMallRepository) GetItemsByCategory(category string) ([]domain.PointsMallItem, error) {
	var items []domain.PointsMallItem
	err := r.DB.Where("is_active = ? AND category = ?", true, category).
		Order("sort_order ASC").
		Find(&items).Error
	return items, err
}

// GetItemByID returns an item by ID
func (r *PointsMallRepository) GetItemByID(id uint) (*domain.PointsMallItem, error) {
	var item domain.PointsMallItem
	err := r.DB.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateItem creates a new item
func (r *PointsMallRepository) CreateItem(item *domain.PointsMallItem) error {
	return r.DB.Create(item).Error
}

// UpdateItem updates an item
func (r *PointsMallRepository) UpdateItem(item *domain.PointsMallItem) error {
	return r.DB.Save(item).Error
}

// ExchangeItem exchanges points for an item (with transaction)
func (r *PointsMallRepository) ExchangeItem(userID uint, item *domain.PointsMallItem, userPoints *domain.UserPoint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Check stock again in transaction
		if !item.IsUnlimited {
			var freshItem domain.PointsMallItem
			if err := tx.First(&freshItem, item.ID).Error; err != nil {
				return err
			}
			if freshItem.GetRemainingStock() <= 0 {
				return domain.ErrOutOfStock
			}
			// Update stock
			if err := tx.Model(&freshItem).Update("used_stock", gorm.Expr("used_stock + 1")).Error; err != nil {
				return err
			}
		}

		// Deduct points
		if err := tx.Model(userPoints).Update("points", gorm.Expr("points - ?", item.PointsCost)).Error; err != nil {
			return err
		}

		// Create exchange record
		exchange := domain.UserPointsExchange{
			UserID:    userID,
			ItemID:    item.ID,
			Points:    item.PointsCost,
			ItemName:  item.Name,
			ItemType:  item.ItemType,
			ItemValue: item.ItemValue,
			Status:    "completed",
		}
		if err := tx.Create(&exchange).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserExchanges returns user's exchange history
func (r *PointsMallRepository) GetUserExchanges(userID uint, page, pageSize int) ([]domain.UserPointsExchange, int64, error) {
	var exchanges []domain.UserPointsExchange
	var total int64

	query := r.DB.Model(&domain.UserPointsExchange{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&exchanges).Error
	return exchanges, total, err
}

// HasExchanged checks if user has already exchanged a specific item
func (r *PointsMallRepository) HasExchanged(userID uint, itemID uint) bool {
	var count int64
	r.DB.Model(&domain.UserPointsExchange{}).
		Where("user_id = ? AND item_id = ? AND status = ?", userID, itemID, "completed").
		Count(&count)
	return count > 0
}

// InitDefaultItems initializes default items if not exists
func (r *PointsMallRepository) InitDefaultItems() error {
	var count int64
	r.DB.Model(&domain.PointsMallItem{}).Count(&count)
	if count > 0 {
		return nil // Already initialized
	}

	items := domain.DefaultPointsMallItems()
	for i := range items {
		if err := r.DB.Create(&items[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
