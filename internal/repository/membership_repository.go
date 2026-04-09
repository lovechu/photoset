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
