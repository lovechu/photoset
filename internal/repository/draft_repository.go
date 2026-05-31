package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// DraftRepository handles database operations for drafts
type DraftRepository struct {
	DB *gorm.DB
}

// NewDraftRepository creates a new DraftRepository
func NewDraftRepository(db *gorm.DB) *DraftRepository {
	return &DraftRepository{DB: db}
}

// Create creates a new draft
func (r *DraftRepository) Create(draft *domain.Draft) error {
	return r.DB.Create(draft).Error
}

// FindByID finds a draft by ID (only for the owner)
func (r *DraftRepository) FindByID(id uint) (*domain.Draft, error) {
	var draft domain.Draft
	err := r.DB.First(&draft, id).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

// ListByUserID returns all drafts for a user
func (r *DraftRepository) ListByUserID(userID uint, page, pageSize int) ([]domain.Draft, int64, error) {
	var drafts []domain.Draft
	var total int64

	query := r.DB.Model(&domain.Draft{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&drafts).Error
	if err != nil {
		return nil, 0, err
	}

	return drafts, total, nil
}

// Update updates a draft
func (r *DraftRepository) Update(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.Draft{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a draft (soft delete)
func (r *DraftRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.Draft{}, id).Error
}

// CountByUserID counts drafts for a user
func (r *DraftRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Draft{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
