package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PhotoSetReviewRepository 套图评价数据访问层
type PhotoSetReviewRepository struct {
	db *gorm.DB
}

// NewPhotoSetReviewRepository 创建套图评价数据访问层
func NewPhotoSetReviewRepository(db *gorm.DB) *PhotoSetReviewRepository {
	return &PhotoSetReviewRepository{db: db}
}

// Create 创建评价
func (r *PhotoSetReviewRepository) Create(review *domain.PhotoSetReview) error {
	return r.db.Create(review).Error
}

// Update 更新评价
func (r *PhotoSetReviewRepository) Update(review *domain.PhotoSetReview) error {
	return r.db.Save(review).Error
}

// Delete 删除评价
func (r *PhotoSetReviewRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.PhotoSetReview{}).Error
}

// GetByID 根据ID获取评价
func (r *PhotoSetReviewRepository) GetByID(id uint) (*domain.PhotoSetReview, error) {
	var review domain.PhotoSetReview
	err := r.db.Preload("User").Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// GetByUserAndPhotoset 获取用户对某套图的评价
func (r *PhotoSetReviewRepository) GetByUserAndPhotoset(userID uint, photosetID uint) (*domain.PhotoSetReview, error) {
	var review domain.PhotoSetReview
	err := r.db.Where("user_id = ? AND photoset_id = ?", userID, photosetID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// ListByPhotoSet 获取套图的评价列表
func (r *PhotoSetReviewRepository) ListByPhotoSet(photosetID uint, page, pageSize int, sortBy string) ([]domain.PhotoSetReview, int64, error) {
	var reviews []domain.PhotoSetReview
	var total int64

	query := r.db.Model(&domain.PhotoSetReview{}).Where("photoset_id = ?", photosetID)

	// 统计总数
	query.Count(&total)

	// 排序
	switch sortBy {
	case "rating_desc":
		query = query.Order("rating DESC, created_at DESC")
	case "rating_asc":
		query = query.Order("rating ASC, created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 分页
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&reviews).Error

	return reviews, total, err
}

// GetSummary 获取评价汇总
func (r *PhotoSetReviewRepository) GetSummary(photosetID uint) (*domain.ReviewSummary, error) {
	summary := &domain.ReviewSummary{}

	// 总数和平均分
	r.db.Model(&domain.PhotoSetReview{}).
		Where("photoset_id = ?", photosetID).
		Select("COUNT(*) as total_count, COALESCE(AVG(rating), 0) as average_rating").
		Row().
		Scan(&summary.TotalCount, &summary.AverageRating)

	// 各评分数量
	r.db.Model(&domain.PhotoSetReview{}).
		Where("photoset_id = ? AND rating = 5", photosetID).
		Count(&summary.Rating5Count)
	r.db.Model(&domain.PhotoSetReview{}).
		Where("photoset_id = ? AND rating = 4", photosetID).
		Count(&summary.Rating4Count)
	r.db.Model(&domain.PhotoSetReview{}).
		Where("photoset_id = ? AND rating = 3", photosetID).
		Count(&summary.Rating3Count)
	r.db.Model(&domain.PhotoSetReview{}).
		Where("photoset_id = ? AND rating = 2", photosetID).
		Count(&summary.Rating2Count)
	r.db.Model(&domain.PhotoSetReview{}).
		Where("photoset_id = ? AND rating = 1", photosetID).
		Count(&summary.Rating1Count)

	// 热门标签
	var tagCounts []domain.TagCount
	r.db.Raw(`
		SELECT tag, COUNT(*) as count FROM (
			SELECT TRIM(SUBSTRING_INDEX(SUBSTRING_INDEX(tags, ',', n.n), ',', -1)) as tag
			FROM photoset_reviews
			CROSS JOIN (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5) n
			WHERE photoset_id = ? AND tags != '' AND n.n <= 1 + LENGTH(tags) - LENGTH(REPLACE(tags, ',', ''))
		) t
		WHERE tag != ''
		GROUP BY tag
		ORDER BY count DESC
		LIMIT 5
	`, photosetID).Scan(&tagCounts)

	summary.TopTags = tagCounts

	return summary, nil
}

// HasUserReviewed 检查用户是否已评价
func (r *PhotoSetReviewRepository) HasUserReviewed(userID uint, photosetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.PhotoSetReview{}).
		Where("user_id = ? AND photoset_id = ?", userID, photosetID).
		Count(&count).Error
	return count > 0, err
}
