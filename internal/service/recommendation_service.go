package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"

	"gorm.io/gorm"
)

// RecommendationService provides recommendation algorithms
type RecommendationService struct {
	db              *gorm.DB
	postRepo        *repository.PostRepository
	postLikeRepo    *repository.PostLikeRepository
	postShareRepo   *repository.PostShareRepository
	postReplyRepo   *repository.PostReplyRepository
	postTagRepo     *repository.PostTagRepository
	postTopicRepo   *repository.PostTopicRepository
	tagRepo         *repository.TagRepository
	topicRepo       *repository.TopicRepository
	hotPostsService *HotPostsService
}

// NewRecommendationService creates a new RecommendationService
func NewRecommendationService(
	db *gorm.DB,
	postRepo *repository.PostRepository,
	postLikeRepo *repository.PostLikeRepository,
	postShareRepo *repository.PostShareRepository,
	postReplyRepo *repository.PostReplyRepository,
	postTagRepo *repository.PostTagRepository,
	postTopicRepo *repository.PostTopicRepository,
	tagRepo *repository.TagRepository,
	topicRepo *repository.TopicRepository,
	hotPostsService *HotPostsService,
) *RecommendationService {
	return &RecommendationService{
		db:              db,
		postRepo:        postRepo,
		postLikeRepo:    postLikeRepo,
		postShareRepo:   postShareRepo,
		postReplyRepo:   postReplyRepo,
		postTagRepo:     postTagRepo,
		postTopicRepo:   postTopicRepo,
		tagRepo:         tagRepo,
		topicRepo:       topicRepo,
		hotPostsService: hotPostsService,
	}
}

// RecommendationType represents the type of recommendation
type RecommendationType string

const (
	RecommendationTypeInterest  RecommendationType = "interest"  // 基于用户兴趣
	RecommendationTypeHot       RecommendationType = "hot"       // 基于热度
	RecommendationTypeSimilar   RecommendationType = "similar"   // 基于相似用户
	RecommendationTypeMixed     RecommendationType = "mixed"     // 混合推荐
)

// RecommendationRequest represents a recommendation request
type RecommendationRequest struct {
	UserID   uint               `json:"user_id"`
	Type     RecommendationType `json:"type"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// RecommendationResponse represents a recommendation response
type RecommendationResponse struct {
	Posts      []domain.Post `json:"posts"`
	Total      int64         `json:"total"`
	Type       RecommendationType `json:"type"`
	Algorithm  string        `json:"algorithm"`
}

// GetRecommendations gets recommended posts based on type
func (s *RecommendationService) GetRecommendations(req *RecommendationRequest) (*RecommendationResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}

	switch req.Type {
	case RecommendationTypeInterest:
		return s.getInterestBasedRecommendations(req.UserID, req.Page, req.PageSize)
	case RecommendationTypeHot:
		return s.getHotBasedRecommendations(req.Page, req.PageSize)
	case RecommendationTypeSimilar:
		return s.getSimilarUserRecommendations(req.UserID, req.Page, req.PageSize)
	case RecommendationTypeMixed:
		return s.getMixedRecommendations(req.UserID, req.Page, req.PageSize)
	default:
		return s.getMixedRecommendations(req.UserID, req.Page, req.PageSize)
	}
}

// getInterestBasedRecommendations gets posts based on user's liked tags/topics
func (s *RecommendationService) getInterestBasedRecommendations(userID uint, page, pageSize int) (*RecommendationResponse, error) {
	// Get user's liked post IDs (last 100 likes)
	var likedPostIDs []uint
	err := s.postLikeRepo.DB.Model(&domain.PostLike{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(100).
		Pluck("post_id", &likedPostIDs).Error
	if err != nil {
		return nil, err
	}

	if len(likedPostIDs) == 0 {
		// Fall back to hot posts if no likes
		return s.getHotBasedRecommendations(page, pageSize)
	}

	// Get tags from liked posts
	var tagIDs []uint
	err = s.postTagRepo.DB.Model(&PostTag{}).
		Where("post_id IN ?", likedPostIDs).
		Pluck("tag_id", &tagIDs).Error
	if err != nil {
		return nil, err
	}

	// Get topics from liked posts
	var topicIDs []uint
	err = s.postTopicRepo.DB.Model(&PostTopic{}).
		Where("post_id IN ?", likedPostIDs).
		Pluck("topic_id", &topicIDs).Error
	if err != nil {
		return nil, err
	}

	// Find posts with similar tags/topics, excluding already liked posts
	query := s.postRepo.DB.Model(&domain.Post{}).
		Where("status = ?", "approved").
		Where("id NOT IN ?", likedPostIDs)

	// If we have tags or topics, use them for filtering
	if len(tagIDs) > 0 || len(topicIDs) > 0 {
		if len(tagIDs) > 0 {
			// Get post IDs with these tags
			var postIDsWithTags []uint
			s.postTagRepo.DB.Model(&PostTag{}).
				Where("tag_id IN ?", tagIDs).
				Distinct("post_id").
				Pluck("post_id", &postIDsWithTags)
			
			if len(postIDsWithTags) > 0 {
				query = query.Where("id IN ?", postIDsWithTags)
			}
		}
		if len(topicIDs) > 0 {
			// Get post IDs with these topics
			var postIDsWithTopics []uint
			s.postTopicRepo.DB.Model(&PostTopic{}).
				Where("topic_id IN ?", topicIDs).
				Distinct("post_id").
				Pluck("post_id", &postIDsWithTopics)
			
			if len(postIDsWithTopics) > 0 {
				query = query.Where("id IN ?", postIDsWithTopics)
			}
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Order by recent activity and popularity
	offset := (page - 1) * pageSize
	var posts []domain.Post
	err = query.
		Preload("User").
		Preload("Tags").
		Preload("Topics").
		Order("created_at DESC, like_count DESC, reply_count DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return &RecommendationResponse{
		Posts:     posts,
		Total:     total,
		Type:      RecommendationTypeInterest,
		Algorithm: "tag_topic_based",
	}, nil
}

// getHotBasedRecommendations gets hot posts
func (s *RecommendationService) getHotBasedRecommendations(page, pageSize int) (*RecommendationResponse, error) {
	posts, total, err := s.hotPostsService.GetHotPosts(page, pageSize)
	if err != nil {
		return nil, err
	}

	return &RecommendationResponse{
		Posts:     posts,
		Total:     total,
		Type:      RecommendationTypeHot,
		Algorithm: "hot_score",
	}, nil
}

// getSimilarUserRecommendations gets posts from users with similar interests
func (s *RecommendationService) getSimilarUserRecommendations(userID uint, page, pageSize int) (*RecommendationResponse, error) {
	// Get user's liked post IDs
	var likedPostIDs []uint
	err := s.postLikeRepo.DB.Model(&domain.PostLike{}).
		Where("user_id = ?", userID).
		Pluck("post_id", &likedPostIDs).Error
	if err != nil {
		return nil, err
	}

	if len(likedPostIDs) == 0 {
		return s.getHotBasedRecommendations(page, pageSize)
	}

	// Find users who liked the same posts (similar users)
	var similarUserIDs []uint
	err = s.postLikeRepo.DB.Model(&domain.PostLike{}).
		Where("post_id IN ? AND user_id != ?", likedPostIDs, userID).
		Select("DISTINCT user_id").
		Pluck("user_id", &similarUserIDs).Error
	if err != nil {
		return nil, err
	}

	if len(similarUserIDs) == 0 {
		return s.getHotBasedRecommendations(page, pageSize)
	}

	// Limit to top 50 similar users
	if len(similarUserIDs) > 50 {
		similarUserIDs = similarUserIDs[:50]
	}

	// Get posts liked by similar users but not by current user
	var recommendedPostIDs []uint
	err = s.postLikeRepo.DB.Model(&domain.PostLike{}).
		Where("user_id IN ? AND post_id NOT IN ?", similarUserIDs, likedPostIDs).
		Select("post_id").
		Group("post_id").
		Order("COUNT(*) DESC").
		Limit(200).
		Pluck("post_id", &recommendedPostIDs).Error
	if err != nil {
		return nil, err
	}

	if len(recommendedPostIDs) == 0 {
		return s.getHotBasedRecommendations(page, pageSize)
	}

	// Get posts with pagination
	query := s.postRepo.DB.Model(&domain.Post{}).
		Where("id IN ? AND status = ?", recommendedPostIDs, "approved")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	var posts []domain.Post
	err = query.
		Preload("User").
		Preload("Tags").
		Preload("Topics").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return &RecommendationResponse{
		Posts:     posts,
		Total:     total,
		Type:      RecommendationTypeSimilar,
		Algorithm: "collaborative_filtering",
	}, nil
}

// getMixedRecommendations gets mixed recommendations
func (s *RecommendationService) getMixedRecommendations(userID uint, page, pageSize int) (*RecommendationResponse, error) {
	// Calculate distribution: 40% interest, 30% hot, 30% similar
	interestSize := pageSize * 40 / 100
	hotSize := pageSize * 30 / 100
	similarSize := pageSize - interestSize - hotSize

	if interestSize == 0 {
		interestSize = 1
	}
	if hotSize == 0 {
		hotSize = 1
	}
	if similarSize == 0 {
		similarSize = 1
	}

	// Get posts from each algorithm
	interestResp, err := s.getInterestBasedRecommendations(userID, 1, interestSize)
	if err != nil {
		interestResp = &RecommendationResponse{Posts: []domain.Post{}}
	}

	hotResp, err := s.getHotBasedRecommendations(1, hotSize)
	if err != nil {
		hotResp = &RecommendationResponse{Posts: []domain.Post{}}
	}

	similarResp, err := s.getSimilarUserRecommendations(userID, 1, similarSize)
	if err != nil {
		similarResp = &RecommendationResponse{Posts: []domain.Post{}}
	}

	// Merge and deduplicate posts
	postMap := make(map[uint]domain.Post)
	for _, post := range interestResp.Posts {
		postMap[post.ID] = post
	}
	for _, post := range hotResp.Posts {
		if _, exists := postMap[post.ID]; !exists {
			postMap[post.ID] = post
		}
	}
	for _, post := range similarResp.Posts {
		if _, exists := postMap[post.ID]; !exists {
			postMap[post.ID] = post
		}
	}

	// Convert to slice
	posts := make([]domain.Post, 0, len(postMap))
	for _, post := range postMap {
		posts = append(posts, post)
	}

	// Sort by created_at descending
	sortPostsByCreatedAt(posts)

	// Apply pagination
	total := int64(len(posts))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(posts) {
		start = len(posts)
	}
	if end > len(posts) {
		end = len(posts)
	}
	pagedPosts := posts[start:end]

	return &RecommendationResponse{
		Posts:     pagedPosts,
		Total:     total,
		Type:      RecommendationTypeMixed,
		Algorithm: "mixed_hybrid",
	}, nil
}

// sortPostsByCreatedAt sorts posts by created_at descending
func sortPostsByCreatedAt(posts []domain.Post) {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].CreatedAt.Before(posts[j].CreatedAt) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
}

// GetUserInterestTags gets user's interest tags based on liked posts
func (s *RecommendationService) GetUserInterestTags(userID uint, limit int) ([]domain.Tag, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get user's liked post IDs
	var likedPostIDs []uint
	err := s.postLikeRepo.DB.Model(&domain.PostLike{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(200).
		Pluck("post_id", &likedPostIDs).Error
	if err != nil {
		return nil, err
	}

	if len(likedPostIDs) == 0 {
		return []domain.Tag{}, nil
	}

	// Get tags from liked posts
	var tagIDs []uint
	err = s.postTagRepo.DB.Model(&PostTag{}).
		Where("post_id IN ?", likedPostIDs).
		Select("tag_id").
		Group("tag_id").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("tag_id", &tagIDs).Error
	if err != nil {
		return nil, err
	}

	if len(tagIDs) == 0 {
		return []domain.Tag{}, nil
	}

	// Get tag details
	var tags []domain.Tag
	err = s.tagRepo.DB.Where("id IN ?", tagIDs).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// GetUserInterestTopics gets user's interest topics based on liked posts
func (s *RecommendationService) GetUserInterestTopics(userID uint, limit int) ([]domain.Topic, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get user's liked post IDs
	var likedPostIDs []uint
	err := s.postLikeRepo.DB.Model(&domain.PostLike{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(200).
		Pluck("post_id", &likedPostIDs).Error
	if err != nil {
		return nil, err
	}

	if len(likedPostIDs) == 0 {
		return []domain.Topic{}, nil
	}

	// Get topics from liked posts
	var topicIDs []uint
	err = s.postTopicRepo.DB.Model(&PostTopic{}).
		Where("post_id IN ?", likedPostIDs).
		Select("topic_id").
		Group("topic_id").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("topic_id", &topicIDs).Error
	if err != nil {
		return nil, err
	}

	if len(topicIDs) == 0 {
		return []domain.Topic{}, nil
	}

	// Get topic details
	var topics []domain.Topic
	err = s.topicRepo.DB.Where("id IN ?", topicIDs).Find(&topics).Error
	if err != nil {
		return nil, err
	}

	return topics, nil
}

// PostTag represents the post_tags join table
type PostTag struct {
	PostID uint `gorm:"column:post_id"`
	TagID  uint `gorm:"column:tag_id"`
}

// TableName specifies the table name
func (PostTag) TableName() string {
	return "post_tags"
}

// PostTopic represents the post_topics join table
type PostTopic struct {
	PostID  uint `gorm:"column:post_id"`
	TopicID uint `gorm:"column:topic_id"`
}

// TableName specifies the table name
func (PostTopic) TableName() string {
	return "post_topics"
}