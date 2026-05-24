package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
	"time"
)

// HotPostsService provides hot posts calculation
type HotPostsService struct {
	postRepo *repository.PostRepository
}

// NewHotPostsService creates a new HotPostsService
func NewHotPostsService(postRepo *repository.PostRepository) *HotPostsService {
	return &HotPostsService{postRepo: postRepo}
}

// GetHotPosts gets hot posts within 7 days
// Algorithm: (reply_count * 2 + like_count + view_count / 10)
func (s *HotPostsService) GetHotPosts(page, pageSize int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// Only consider posts within 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	query := s.postRepo.DB.Model(&domain.Post{}).
		Where("status = ?", "approved").
		Where("created_at >= ?", sevenDaysAgo)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Order by hotness formula
	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("reply_count * 2 + like_count + view_count / 10 DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetHotPostsRaw gets hot posts with raw SQL for more complex calculations
func (s *HotPostsService) GetHotPostsRaw(page, pageSize int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// Count total
	err := s.postRepo.DB.Model(&domain.Post{}).
		Where("status = ? AND created_at >= ?", "approved", sevenDaysAgo).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Use raw SQL for precise hotness calculation
	offset := (page - 1) * pageSize
	err = s.postRepo.DB.Raw(`
		SELECT posts.*, 
		       (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) as like_count,
		       (SELECT COUNT(*) FROM post_replies WHERE post_replies.post_id = posts.id) as reply_count
		FROM posts 
		WHERE posts.status = 'approved' 
			  AND posts.created_at >= ?
		ORDER BY (reply_count * 2 + like_count + view_count / 10) DESC, created_at DESC
		LIMIT ? OFFSET ?
	`, sevenDaysAgo, pageSize, offset).Scan(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	// Load associations
	for i := range posts {
		s.postRepo.DB.Preload("User").First(&posts[i], posts[i].ID)
	}

	return posts, total, nil
}
