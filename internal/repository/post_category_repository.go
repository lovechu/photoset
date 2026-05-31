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

// CategoryPostCount holds a category key and its post count
type CategoryPostCount struct {
	Category string `gorm:"column:category"`
	Count    int64  `gorm:"column:cnt"`
}

// BatchCountPosts counts posts for all categories in one query using GROUP BY.
// Returns a map: category_key -> post_count.
// This avoids the N+1 query problem in ListCategories.
func (r *PostCategoryRepository) BatchCountPosts() (map[string]int64, error) {
	var results []CategoryPostCount
	err := r.DB.Raw(`
		SELECT category, COUNT(*) AS cnt
		FROM posts
		WHERE deleted_at IS NULL
		GROUP BY category
	`).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int64, len(results))
	for _, r := range results {
		countMap[r.Category] = r.Count
	}
	return countMap, nil
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
