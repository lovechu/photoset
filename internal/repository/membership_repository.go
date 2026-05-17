package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type MembershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

// FindByID 根据ID查询会员套餐
func (r *MembershipRepository) FindByID(id uint) (*domain.MembershipPlan, error) {
	var plan domain.MembershipPlan
	err := r.db.First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// ListActive 获取所有上架的套餐
func (r *MembershipRepository) ListActive() ([]domain.MembershipPlan, error) {
	var plans []domain.MembershipPlan
	err := r.db.Where("status = ?", 1).Find(&plans).Error
	return plans, err
}

// ============ Admin APIs ============

// ListAll 分页查询所有会员套餐（管理员用）
func (r *MembershipRepository) ListAll(page, pageSize int) ([]domain.MembershipPlan, int64, error) {
	var plans []domain.MembershipPlan
	var total int64

	query := r.db.Model(&domain.MembershipPlan{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&plans).Error
	return plans, total, err
}

// Create 创建会员套餐
func (r *MembershipRepository) Create(plan *domain.MembershipPlan) error {
	return r.db.Create(plan).Error
}

// Update 更新会员套餐
func (r *MembershipRepository) Update(plan *domain.MembershipPlan) error {
	return r.db.Save(plan).Error
}

// Delete 删除会员套餐
func (r *MembershipRepository) Delete(id uint) error {
	return r.db.Delete(&domain.MembershipPlan{}, id).Error
}
