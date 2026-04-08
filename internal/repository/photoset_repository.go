package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

type PhotoSetRepository struct {
	db *gorm.DB
}

func NewPhotoSetRepository(db *gorm.DB) *PhotoSetRepository {
	return &PhotoSetRepository{db: db}
}

// Create 创建套图
func (r *PhotoSetRepository) Create(photoset *domain.PhotoSet) error {
	return r.db.Create(photoset).Error
}

// FindByID 根据ID查询套图
func (r *PhotoSetRepository) FindByID(id uint) (*domain.PhotoSet, error) {
	var photoset domain.PhotoSet
	err := r.db.Preload("User").Preload("Photos", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Preload("Tags").First(&photoset, id).Error
	if err != nil {
		return nil, err
	}
	return &photoset, nil
}

// FindByIDWithoutPhotos 根据ID查询套图（不预加载图片）
func (r *PhotoSetRepository) FindByIDWithoutPhotos(id uint) (*domain.PhotoSet, error) {
	var photoset domain.PhotoSet
	err := r.db.Preload("User").Preload("Tags").First(&photoset, id).Error
	if err != nil {
		return nil, err
	}
	return &photoset, nil
}

// List 查询套图列表
func (r *PhotoSetRepository) List(page, pageSize int, tag string) ([]domain.PhotoSet, int64, error) {
	var photosets []domain.PhotoSet
	var total int64

	query := r.db.Model(&domain.PhotoSet{})

	// 按标签筛选
	if tag != "" {
		query = query.Joins("INNER JOIN photoset_tags ON photosets.id = photoset_tags.photoset_id").
			Joins("INNER JOIN tags ON photoset_tags.tag_id = tags.id").
			Where("tags.name = ?", tag)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&photosets).Error

	return photosets, total, err
}

// FindByName 根据标签名查询标签
func (r *PhotoSetRepository) FindTagByName(name string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// CreateTag 创建标签
func (r *PhotoSetRepository) CreateTag(tag *domain.Tag) error {
	return r.db.Create(tag).Error
}

// ListTags 查询所有标签
func (r *PhotoSetRepository) ListTags() ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.Order("name ASC").Find(&tags).Error
	return tags, err
}

// CreatePhotos 批量创建图片
func (r *PhotoSetRepository) CreatePhotos(photos []domain.Photo) error {
	return r.db.Create(&photos).Error
}

// CreatePhotoSetTags 创建套图标签关联
func (r *PhotoSetRepository) CreatePhotoSetTags(photosetID uint, tagIDs []uint) error {
	for _, tagID := range tagIDs {
		photosetTag := map[string]interface{}{
			"photoset_id": photosetID,
			"tag_id":      tagID,
		}
		if err := r.db.Table("photoset_tags").Create(&photosetTag).Error; err != nil {
			return err
		}
	}
	return nil
}
