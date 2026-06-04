package service

import (
	"errors"
	"photoset/internal/domain"
	"photoset/internal/repository"
)

// ReviewService 评价服务
type ReviewService struct {
	reviewRepo   *repository.PhotoSetReviewRepository
	photosetRepo *repository.PhotoSetRepository
}

// NewReviewService 创建评价服务
func NewReviewService(reviewRepo *repository.PhotoSetReviewRepository, photosetRepo *repository.PhotoSetRepository) *ReviewService {
	return &ReviewService{
		reviewRepo:   reviewRepo,
		photosetRepo: photosetRepo,
	}
}

// CreateReview 创建评价
func (s *ReviewService) CreateReview(userID uint, photosetID uint, rating int, tags string, content string, isAnonymous bool) (*domain.PhotoSetReview, error) {
	// 验证评分范围
	if rating < 1 || rating > 5 {
		return nil, errors.New("评分必须在 1-5 之间")
	}

	// 检查套图是否存在
	_, err := s.photosetRepo.FindByID(photosetID)
	if err != nil {
		return nil, errors.New("套图不存在")
	}

	// 检查是否已评价
	hasReviewed, err := s.reviewRepo.HasUserReviewed(userID, photosetID)
	if err != nil {
		return nil, err
	}
	if hasReviewed {
		return nil, errors.New("您已评价过该套图")
	}

	review := &domain.PhotoSetReview{
		UserID:      userID,
		PhotoSetID:  photosetID,
		Rating:      rating,
		Tags:        tags,
		Content:     content,
		IsAnonymous: isAnonymous,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	return s.reviewRepo.GetByID(review.ID)
}

// UpdateReview 更新评价
func (s *ReviewService) UpdateReview(reviewID uint, userID uint, rating int, tags string, content string, isAnonymous bool) (*domain.PhotoSetReview, error) {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, errors.New("评价不存在")
	}

	if review.UserID != userID {
		return nil, errors.New("无权修改此评价")
	}

	if rating < 1 || rating > 5 {
		return nil, errors.New("评分必须在 1-5 之间")
	}

	review.Rating = rating
	review.Tags = tags
	review.Content = content
	review.IsAnonymous = isAnonymous

	if err := s.reviewRepo.Update(review); err != nil {
		return nil, err
	}

	return s.reviewRepo.GetByID(review.ID)
}

// DeleteReview 删除评价
func (s *ReviewService) DeleteReview(reviewID uint, userID uint) error {
	return s.reviewRepo.Delete(reviewID, userID)
}

// GetReview 获取评价详情
func (s *ReviewService) GetReview(reviewID uint) (*domain.PhotoSetReview, error) {
	return s.reviewRepo.GetByID(reviewID)
}

// GetUserReview 获取用户对某套图的评价
func (s *ReviewService) GetUserReview(userID uint, photosetID uint) (*domain.PhotoSetReview, error) {
	return s.reviewRepo.GetByUserAndPhotoset(userID, photosetID)
}

// ListReviews 获取套图评价列表
func (s *ReviewService) ListReviews(photosetID uint, page, pageSize int, sortBy string) ([]domain.PhotoSetReview, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.reviewRepo.ListByPhotoSet(photosetID, page, pageSize, sortBy)
}

// GetReviewSummary 获取评价汇总
func (s *ReviewService) GetReviewSummary(photosetID uint) (*domain.ReviewSummary, error) {
	return s.reviewRepo.GetSummary(photosetID)
}
