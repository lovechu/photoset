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
func (r *PhotoSetRepository) List(page, pageSize int, tag string, keyword string, userID uint, onlyMine bool) ([]domain.PhotoSet, int64, error) {
	var photosets []domain.PhotoSet
	var total int64

	query := r.db.Model(&domain.PhotoSet{})

	// 按标签筛选
	if tag != "" {
		query = query.Joins("INNER JOIN photoset_tags ON photosets.id = photoset_tags.photoset_id").
			Joins("INNER JOIN tags ON photoset_tags.tag_id = tags.id").
			Where("tags.name = ?", tag)
	}

	// 关键词搜索（FULLTEXT 全文检索）
	if keyword != "" {
		query = query.Where("MATCH(photosets.title, photosets.description) AGAINST(? IN BOOLEAN MODE)", keyword)
	}

	// 只看自己的套图
	if onlyMine && userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，使用子查询获取 photo_count
	offset := (page - 1) * pageSize
	err := r.db.Table("photosets").
		Select("photosets.*, (SELECT COUNT(*) FROM photos WHERE photos.photo_set_id = photosets.id) AS photo_count").
		Preload("User").Preload("Tags").
		Where(query).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Scan(&photosets).Error

	if err != nil {
		return nil, 0, err
	}

	return photosets, total, nil
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
			"tag_id":     tagID,
		}
		if err := r.db.Table("photoset_tags").Create(&photosetTag).Error; err != nil {
			return err
		}
	}
	return nil
}

// Update 更新套图基础信息（不包含图片列表）
func (r *PhotoSetRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&domain.PhotoSet{}).Where("id = ?", id).Updates(updates).Error
}

// ReplaceTags 替换套图标签（先删后插）
func (r *PhotoSetRepository) ReplaceTags(photosetID uint, tagNames []string) error {
	// 删除旧关联
	if err := r.db.Exec("DELETE FROM photoset_tags WHERE photoset_id = ?", photosetID).Error; err != nil {
		return err
	}
	if len(tagNames) == 0 {
		return nil
	}
	// 写入新关联
	for _, name := range tagNames {
		tag, err := r.FindTagByName(name)
		if err != nil {
			tag = &domain.Tag{Name: name}
			if err := r.CreateTag(tag); err != nil {
				return err
			}
		}
		photosetTag := map[string]interface{}{
			"photoset_id": photosetID,
			"tag_id":      tag.ID,
		}
		r.db.Table("photoset_tags").Create(&photosetTag)
	}
	return nil
}

// Delete 软删除套图（GORM 软删除，同时级联删除关联 photos 和 tags）
func (r *PhotoSetRepository) Delete(id uint) error {
	// 先删除关联的 photos
	if err := r.db.Where("photo_set_id = ?", id).Delete(&domain.Photo{}).Error; err != nil {
		return err
	}
	// 删除关联的 photoset_tags
	if err := r.db.Exec("DELETE FROM photoset_tags WHERE photoset_id = ?", id).Error; err != nil {
		return err
	}
	// 软删除套图本身
	return r.db.Delete(&domain.PhotoSet{}, id).Error
}

// ReplacePhotos 替换套图图片（先删后插）
func (r *PhotoSetRepository) ReplacePhotos(photosetID uint, photos []domain.Photo) error {
	if err := r.db.Where("photo_set_id = ?", photosetID).Delete(&domain.Photo{}).Error; err != nil {
		return err
	}
	if len(photos) == 0 {
		return nil
	}
	for i := range photos {
		photos[i].PhotoSetID = photosetID
	}
	return r.db.Create(&photos).Error
}
