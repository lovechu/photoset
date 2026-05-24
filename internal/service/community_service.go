package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"

	"gorm.io/gorm"
)

// CommunityService provides community business logic
type CommunityService struct {
	postRepo      *repository.PostRepository
	replyRepo     *repository.PostReplyRepository
	likeRepo      *repository.PostLikeRepository
	replyLikeRepo *repository.PostReplyLikeRepository
	pointRepo     *repository.UserPointRepository
	reportRepo    *repository.PostReportRepository
	pointService  *PointService
	filterService *SensitiveFilterService
}

// NewCommunityService creates a new CommunityService
func NewCommunityService(
	postRepo *repository.PostRepository,
	replyRepo *repository.PostReplyRepository,
	likeRepo *repository.PostLikeRepository,
	replyLikeRepo *repository.PostReplyLikeRepository,
	pointRepo *repository.UserPointRepository,
	reportRepo *repository.PostReportRepository,
	pointService *PointService,
	filterService *SensitiveFilterService,
) *CommunityService {
	return &CommunityService{
		postRepo:      postRepo,
		replyRepo:     replyRepo,
		likeRepo:      likeRepo,
		replyLikeRepo: replyLikeRepo,
		pointRepo:     pointRepo,
		reportRepo:    reportRepo,
		pointService:  pointService,
		filterService: filterService,
	}
}

// CreatePost creates a new post with sensitive word filtering
func (s *CommunityService) CreatePost(userID uint, req *CreatePostRequest) (*domain.Post, error) {
	// Validate request
	post := &domain.Post{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		PhotosetID: req.PhotosetID,
		Category:    req.Category,
		Visibility:  req.Visibility,
		Status:     string(domain.PostStatusApproved), // First post, then review
	}

	if err := post.Validate(); err != nil {
		return nil, err
	}

	// Filter sensitive words
	filteredTitle, _ := s.filterService.FilterText(req.Title)
	filteredContent, _ := s.filterService.FilterText(req.Content)
	post.Title = filteredTitle
	post.Content = filteredContent

	// Create post in transaction
	err := s.postRepo.DB.Transaction(func(tx *gorm.DB) error {
		// Create post
		if err := s.postRepo.Create(post); err != nil {
			return err
		}

		// Add points
		if err := s.pointService.AddPointsForPost(userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load associations
	post, _ = s.postRepo.FindByID(post.ID)
	return post, nil
}

// CreateReply creates a new reply with sensitive word filtering
func (s *CommunityService) CreateReply(userID, postID uint, req *CreateReplyRequest) (*domain.PostReply, error) {
	// Check if post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, domain.ErrPostNotFound
	}

	// Validate request
	reply := &domain.PostReply{
		PostID:        postID,
		UserID:        userID,
		Content:       req.Content,
		ParentReplyID: req.ParentReplyID,
	}

	if err := reply.Validate(); err != nil {
		return nil, err
	}

	// Filter sensitive words
	filteredContent, _ := s.filterService.FilterText(req.Content)
	reply.Content = filteredContent

	// Create reply in transaction
	err = s.replyRepo.DB.Transaction(func(tx *gorm.DB) error {
		// Create reply
		if err := s.replyRepo.Create(reply); err != nil {
			return err
		}

		// Increment post reply count
		if err := s.postRepo.IncrementReplyCount(postID); err != nil {
			return err
		}

		// Add points
		if err := s.pointService.AddPointsForReply(userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load associations
	reply, _ = s.replyRepo.FindByID(reply.ID)
	return reply, nil
}

// TogglePostLike toggles like status for a post
func (s *CommunityService) TogglePostLike(userID, postID uint) (string, int, error) {
	// Check if post exists
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return "", 0, domain.ErrPostNotFound
	}

	// Check if already liked
	liked, err := s.likeRepo.Exists(userID, postID)
	if err != nil {
		return "", 0, err
	}

	var action string
	var likeCount int

	if liked {
		// Unlike: delete like record and decrement count
		err = s.likeRepo.DB.Transaction(func(tx *gorm.DB) error {
			if err := s.likeRepo.Delete(userID, postID); err != nil {
				return err
			}
			return s.postRepo.DecrementLikeCount(postID)
		})
		action = "unliked"
	} else {
		// Like: create like record and increment count
		like := &domain.PostLike{
			UserID: userID,
			PostID: postID,
		}
		err = s.likeRepo.DB.Transaction(func(tx *gorm.DB) error {
			if err := s.likeRepo.Create(like); err != nil {
				return err
			}
			if err := s.postRepo.IncrementLikeCount(postID); err != nil {
				return err
			}
			// Add points to post author
			return s.pointService.AddPointsForLiked(post.UserID, 2)
		})
		action = "liked"
	}

	if err != nil {
		return "", 0, err
	}

	// Get updated like count
	var count int64
	s.likeRepo.DB.Model(&domain.PostLike{}).Where("post_id = ?", postID).Count(&count)
	likeCount = int(count)

	return action, likeCount, nil
}

// ToggleReplyLike toggles like status for a reply
func (s *CommunityService) ToggleReplyLike(userID, replyID uint) (string, int, error) {
	// Check if reply exists
	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return "", 0, domain.ErrReplyNotFound
	}

	// Check if already liked
	liked, err := s.replyLikeRepo.Exists(userID, replyID)
	if err != nil {
		return "", 0, err
	}

	var action string
	var likeCount int

	if liked {
		// Unlike
		err = s.replyLikeRepo.DB.Transaction(func(tx *gorm.DB) error {
			if err := s.replyLikeRepo.Delete(userID, replyID); err != nil {
				return err
			}
			return s.replyRepo.DecrementLikeCount(replyID)
		})
		action = "unliked"
	} else {
		// Like
		like := &domain.PostReplyLike{
			UserID: userID,
			ReplyID: replyID,
		}
		err = s.replyLikeRepo.DB.Transaction(func(tx *gorm.DB) error {
			if err := s.replyLikeRepo.Create(like); err != nil {
				return err
			}
			if err := s.replyRepo.IncrementLikeCount(replyID); err != nil {
				return err
			}
			// Add points to reply author
			return s.pointService.AddPointsForLiked(reply.UserID, 1)
		})
		action = "liked"
	}

	if err != nil {
		return "", 0, err
	}

	// Get updated like count
	var count int64
	s.replyLikeRepo.DB.Model(&domain.PostReplyLike{}).Where("reply_id = ?", replyID).Count(&count)
	likeCount = int(count)

	return action, likeCount, nil
}

// ReportPost reports a post
func (s *CommunityService) ReportPost(reporterID, postID uint, reason string) error {
	if reason == "" {
		return domain.ErrReportReasonRequired
	}

	// Check if post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		return domain.ErrPostNotFound
	}

	report := &domain.PostReport{
		PostID:     &postID,
		ReporterID: reporterID,
		Reason:      reason,
		Status:      string(domain.ReportStatusPending),
	}

	return s.reportRepo.Create(report)
}

// ReportReply reports a reply
func (s *CommunityService) ReportReply(reporterID, replyID uint, reason string) error {
	if reason == "" {
		return domain.ErrReportReasonRequired
	}

	// Check if reply exists
	_, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return domain.ErrReplyNotFound
	}

	report := &domain.PostReport{
		ReplyID:    &replyID,
		ReporterID: reporterID,
		Reason:      reason,
		Status:      string(domain.ReportStatusPending),
	}

	return s.reportRepo.Create(report)
}

// GetPostDetail gets post detail and increments view count
func (s *CommunityService) GetPostDetail(postID uint) (*domain.Post, error) {
	post, err := s.postRepo.FindByIDWithCounts(postID)
	if err != nil {
		return nil, domain.ErrPostNotFound
	}

	// Increment view count
	s.postRepo.IncrementViewCount(postID)

	return post, nil
}

// Request types
type CreatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	PhotosetID *uint `json:"photoset_id"`
	Category   string `json:"category"`
	Visibility string `json:"visibility"`
}

type CreateReplyRequest struct {
	Content       string `json:"content" binding:"required"`
	ParentReplyID *uint `json:"parent_reply_id"`
}
