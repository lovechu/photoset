package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// SensitiveWordRepository handles database operations for sensitive words
type SensitiveWordRepository struct {
	DB *gorm.DB
}

// NewSensitiveWordRepository creates a new SensitiveWordRepository
func NewSensitiveWordRepository(db *gorm.DB) *SensitiveWordRepository {
	return &SensitiveWordRepository{DB: db}
}

// Create creates a new sensitive word
func (r *SensitiveWordRepository) Create(word *domain.SensitiveWord) error {
	return r.DB.Create(word).Error
}

// FindByID finds a sensitive word by ID
func (r *SensitiveWordRepository) FindByID(id uint) (*domain.SensitiveWord, error) {
	var word domain.SensitiveWord
	err := r.DB.First(&word, id).Error
	if err != nil {
		return nil, err
	}
	return &word, nil
}

// FindByWord finds a sensitive word by word text
func (r *SensitiveWordRepository) FindByWord(word string) (*domain.SensitiveWord, error) {
	var sensitiveWord domain.SensitiveWord
	err := r.DB.Where("word = ?", word).First(&sensitiveWord).Error
	if err != nil {
		return nil, err
	}
	return &sensitiveWord, nil
}

// List returns all sensitive words with pagination
func (r *SensitiveWordRepository) List(page, pageSize int, isActive *bool) ([]domain.SensitiveWord, int64, error) {
	var words []domain.SensitiveWord
	var total int64

	query := r.DB.Model(&domain.SensitiveWord{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&words).Error
	if err != nil {
		return nil, 0, err
	}

	return words, total, nil
}

// Update updates a sensitive word
func (r *SensitiveWordRepository) Update(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.SensitiveWord{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a sensitive word
func (r *SensitiveWordRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.SensitiveWord{}, id).Error
}

// LoadAllActive loads all active sensitive words
func (r *SensitiveWordRepository) LoadAllActive() ([]domain.SensitiveWord, error) {
	var words []domain.SensitiveWord
	err := r.DB.Where("is_active = ?", true).Find(&words).Error
	if err != nil {
		return nil, err
	}
	return words, nil
}
