package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostCategoryRepository handles database operations for post categories
type PostCategoryRepository struct {
	DB *gorm.DB
}

// NewPostCategoryRepository creates a new PostCategoryRepository
func NewPostCategoryRepository(db *gorm.DB) *PostCategoryRepository {
	return &PostCategoryRepository{DB: db}
}

// ListCategories returns all categories sorted by sort_order DESC, then id ASC
func (r *PostCategoryRepository) ListCategories() ([]domain.CommunityCategory, error) {
	var categories []domain.CommunityCategory
	err := r.DB.Order("sort_order DESC, id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByKey finds a category by its unique key
func (r *PostCategoryRepository) GetCategoryByKey(key string) (*domain.CommunityCategory, error) {
	var category domain.CommunityCategory
	err := r.DB.Where("`key` = ?", key).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoryByID finds a category by its primary key
func (r *PostCategoryRepository) GetCategoryByID(id uint) (*domain.CommunityCategory, error) {
	var category domain.CommunityCategory
	err := r.DB.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CreateCategory inserts a new category record
func (r *PostCategoryRepository) CreateCategory(cat *domain.CommunityCategory) error {
	return r.DB.Create(cat).Error
}

// UpdateCategory updates a category with the given map of fields
func (r *PostCategoryRepository) UpdateCategory(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.CommunityCategory{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategory removes a category record by ID
func (r *PostCategoryRepository) DeleteCategory(id uint) error {
	return r.DB.Delete(&domain.CommunityCategory{}, id).Error
}

// CountPostsByCategory counts the number of posts in a given category by key
func (r *PostCategoryRepository) CountPostsByCategory(key string) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Post{}).Where("category = ?", key).Count(&count).Error
	return count, err
}

// GetActiveKeys returns all category keys currently in the database
// Used for validating post category input
func (r *PostCategoryRepository) GetActiveKeys() ([]string, error) {
	var keys []string
	err := r.DB.Model(&domain.CommunityCategory{}).Pluck("`key`", &keys).Error
	if err != nil {
		return nil, err
	}
	return keys, nil
}
