package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostShareRepository handles database operations for post shares
type PostShareRepository struct {
	DB *gorm.DB
}

// NewPostShareRepository creates a new PostShareRepository
func NewPostShareRepository(db *gorm.DB) *PostShareRepository {
	return &PostShareRepository{DB: db}
}

// Create creates a new share record
func (r *PostShareRepository) Create(share *domain.PostShare) error {
	return r.DB.Create(share).Error
}

// CountByPostID counts shares for a post
func (r *PostShareRepository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.PostShare{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountByPostIDAndPlatform counts shares for a post grouped by platform
func (r *PostShareRepository) CountByPostIDAndPlatform(postID uint) (map[string]int64, error) {
	type PlatformCount struct {
		Platform string
		Count    int64
	}
	var results []PlatformCount
	err := r.DB.Model(&domain.PostShare{}).
		Select("platform, COUNT(*) as count").
		Where("post_id = ?", postID).
		Group("platform").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Platform] = r.Count
	}
	return counts, nil
}
