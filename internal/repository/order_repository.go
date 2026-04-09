package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create 创建订单
func (r *OrderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

// FindByID 根据ID查询订单
func (r *OrderRepository) FindByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Membership").Preload("PhotoSet").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Update 更新订单
func (r *OrderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

// ListByUserID 根据用户ID查询订单列表（分页）
func (r *OrderRepository) ListByUserID(userID uint, page, pageSize int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	query := r.db.Model(&domain.Order{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Preload("Membership").Preload("PhotoSet").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}

// CountPaidByUserAndPhotoSet 查询用户对某套图是否已支付
func (r *OrderRepository) HasPaidOrder(userID, photosetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Order{}).
		Where("user_id = ? AND photoset_id = ? AND status = ?", userID, photosetID, "paid").
		Count(&count).Error
	return count > 0, err
}

// CountStats 统计订单总数和总收入
func (r *OrderRepository) CountStats() (totalOrders int64, totalRevenue float64, err error) {
	if err = r.db.Model(&domain.Order{}).Where("status = ?", "paid").Count(&totalOrders).Error; err != nil {
		return
	}
	if err = r.db.Model(&domain.Order{}).
		Where("status = ?", "paid").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return
	}
	return
}

// FindPaidOrderByUser 根据用户ID和订单ID查询已支付订单
func (r *OrderRepository) FindPaidOrderByUser(userID, orderID uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Where("id = ? AND user_id = ? AND status = ?", orderID, userID, "paid").
		Preload("Membership").First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindPaidOrder 根据订单ID查询已支付订单（admin用）
func (r *OrderRepository) FindPaidOrder(orderID uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Where("id = ? AND status = ?", orderID, "paid").
		Preload("Membership").First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
