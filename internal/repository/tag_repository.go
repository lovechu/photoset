package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// TagRepository handles database operations for tags
type TagRepository struct {
	DB *gorm.DB
}

// NewTagRepository creates a new TagRepository
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{DB: db}
}

// Create creates a new tag
func (r *TagRepository) Create(tag *domain.Tag) error {
	return r.DB.Create(tag).Error
}

// FindByID finds a tag by ID
func (r *TagRepository) FindByID(id uint) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.DB.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindByName finds a tag by name
func (r *TagRepository) FindByName(name string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindOrCreate finds a tag by name or creates it if it doesn't exist
func (r *TagRepository) FindOrCreate(name string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		tag = domain.Tag{Name: name}
		if err := r.DB.Create(&tag).Error; err != nil {
			return nil, err
		}
		return &tag, nil
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindOrCreateBatch finds or creates multiple tags by name
func (r *TagRepository) FindOrCreateBatch(names []string) ([]domain.Tag, error) {
	var tags []domain.Tag
	for _, name := range names {
		tag, err := r.FindOrCreate(name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

// GetAll returns all tags
func (r *TagRepository) GetAll() ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.DB.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetPopularTags returns popular tags sorted by usage count
func (r *TagRepository) GetPopularTags(limit int) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.DB.
		Select("tags.*, COUNT(post_tags.post_id) as post_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Group("tags.id").
		Order("post_count DESC").
		Limit(limit).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// Delete deletes a tag by ID
func (r *TagRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.Tag{}, id).Error
}

// SearchByName searches tags by name (partial match)
func (r *TagRepository) SearchByName(keyword string, limit int) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.DB.Where("name LIKE ?", "%"+keyword+"%").
		Limit(limit).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
