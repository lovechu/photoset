package service

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"photoset/internal/domain"
	"photoset/internal/logger"
	"photoset/internal/repository"
	"photoset/internal/storage"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CommunityService provides community business logic
type CommunityService struct {
	postRepo      *repository.PostRepository
	replyRepo     *repository.PostReplyRepository
	likeRepo      *repository.PostLikeRepository
	replyLikeRepo *repository.PostReplyLikeRepository
	shareRepo     *repository.PostShareRepository
	pointRepo     *repository.UserPointRepository
	reportRepo    *repository.PostReportRepository
	categoryRepo  *repository.PostCategoryRepository
	draftRepo     *repository.DraftRepository
	tagRepo       *repository.TagRepository
	postTagRepo   *repository.PostTagRepository
	topicRepo     *repository.TopicRepository
	postTopicRepo *repository.PostTopicRepository
	pointService  *PointService
	filterService *SensitiveFilterService
	mentionService *MentionService
	storage       storage.Storage
}

// NewCommunityService creates a new CommunityService
func NewCommunityService(
	postRepo *repository.PostRepository,
	replyRepo *repository.PostReplyRepository,
	likeRepo *repository.PostLikeRepository,
	replyLikeRepo *repository.PostReplyLikeRepository,
	shareRepo *repository.PostShareRepository,
	pointRepo *repository.UserPointRepository,
	reportRepo *repository.PostReportRepository,
	categoryRepo *repository.PostCategoryRepository,
	draftRepo *repository.DraftRepository,
	tagRepo *repository.TagRepository,
	postTagRepo *repository.PostTagRepository,
	topicRepo *repository.TopicRepository,
	postTopicRepo *repository.PostTopicRepository,
	pointService *PointService,
	filterService *SensitiveFilterService,
	mentionService *MentionService,
	stor storage.Storage,
) *CommunityService {
	return &CommunityService{
		postRepo:      postRepo,
		replyRepo:     replyRepo,
		likeRepo:      likeRepo,
		replyLikeRepo: replyLikeRepo,
		shareRepo:     shareRepo,
		pointRepo:     pointRepo,
		reportRepo:    reportRepo,
		categoryRepo:  categoryRepo,
		draftRepo:     draftRepo,
		tagRepo:       tagRepo,
		postTagRepo:   postTagRepo,
		topicRepo:     topicRepo,
		postTopicRepo: postTopicRepo,
		pointService:  pointService,
		filterService: filterService,
		mentionService: mentionService,
		storage:       stor,
	}
}

// CreatePost creates a new post with sensitive word filtering
func (s *CommunityService) CreatePost(userID uint, req *CreatePostRequest, authorIPLocation string) (*domain.Post, error) {
	// Validate request
	post := &domain.Post{
		UserID:            userID,
		Title:             req.Title,
		Content:           req.Content,
		PhotosetID:        req.PhotosetID,
		Category:          req.Category,
		PostType:          req.PostType,
		Visibility:        req.Visibility,
		IsOriginal:        req.IsOriginal,
		AuthorIPLocation:  authorIPLocation,
		Status:            string(domain.PostStatusApproved), // Auto-approve on creation
	}

	// 解析定时发布时间
	if req.ScheduledAt != nil && *req.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			return nil, fmt.Errorf("invalid scheduled_at format, expected RFC3339: %w", err)
		}
		post.ScheduledAt = &t
		// 定时发布的帖子状态设为 pending
		post.Status = string(domain.PostStatusPending)
	}

	if err := post.Validate(); err != nil {
		return nil, err
	}

	// Validate category exists in DB (replaces hardcoded constant check)
	if post.Category != "" {
		activeKeys, err := s.categoryRepo.GetActiveKeys()
		if err != nil {
			return nil, err
		}
		if !slices.Contains(activeKeys, post.Category) {
			return nil, domain.ErrInvalidCategory
		}
	}

	// Filter sensitive words
	filteredTitle, _ := s.filterService.FilterTextAdvanced(req.Title)
	filteredContent, _ := s.filterService.FilterTextAdvanced(req.Content)
	post.Title = filteredTitle
	post.Content = filteredContent

	// Create post in transaction
	err := s.postRepo.DB.Transaction(func(tx *gorm.DB) error {
		// Create post using tx (transaction connection)
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Handle tags
	if len(req.Tags) > 0 {
		tags, err := s.tagRepo.FindOrCreateBatch(req.Tags)
		if err != nil {
			// Log error but don't fail post creation
			_ = err
		} else {
			if err := s.postTagRepo.AddTagsToPost(post.ID, tags); err != nil {
				// Log error but don't fail post creation
				_ = err
			}
		}
	}

	// Handle topics with validation and filtering
	if len(req.Topics) > 0 {
		// 限制话题数量（最多3个）
		topicsToProcess := req.Topics
		if len(topicsToProcess) > 3 {
			topicsToProcess = topicsToProcess[:3]
		}
		
		// 过滤和验证话题名称
		var validTopics []string
		seenTopics := make(map[string]bool)
		
		for _, topicName := range topicsToProcess {
			// 去除首尾空格和#
			cleanName := topicName
			for len(cleanName) > 0 && cleanName[0] == '#' {
				cleanName = cleanName[1:]
			}
			for len(cleanName) > 0 && cleanName[len(cleanName)-1] == '#' {
				cleanName = cleanName[:len(cleanName)-1]
			}
			cleanName = strings.TrimSpace(cleanName)
			
			// 跳过空话题
			if cleanName == "" {
				continue
			}
			
			// 长度验证：2-100个字符
			if len([]rune(cleanName)) < 2 || len([]rune(cleanName)) > 100 {
				continue
			}
			
			// 格式验证：只允许中文、字母、数字、下划线
			validPattern := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_]+$`)
			if !validPattern.MatchString(cleanName) {
				continue
			}
			
			// 不能以数字开头
			if len(cleanName) > 0 && cleanName[0] >= '0' && cleanName[0] <= '9' {
				continue
			}
			
			// 敏感词过滤
			filteredName, _ := s.filterService.FilterTextAdvanced(cleanName)
			if filteredName != cleanName {
				// 包含敏感词，跳过
				continue
			}
			
			// 去重
			if !seenTopics[cleanName] {
				seenTopics[cleanName] = true
				validTopics = append(validTopics, cleanName)
			}
		}
		
		// 处理有效的话题
		if len(validTopics) > 0 {
			topics, err := s.topicRepo.FindOrCreateBatch(validTopics)
			if err != nil {
				// Log error but don't fail post creation
				_ = err
			} else {
				if err := s.postTopicRepo.AddTopicsToPost(post.ID, topics); err != nil {
					// Log error but don't fail post creation
					_ = err
				}
				// Increment post count for each topic
				for _, topic := range topics {
					go s.topicRepo.IncrementPostCount(topic.ID)
				}
			}
		}
	}

	// Add points after successful post creation (non-critical, can fail independently)
	go s.pointService.AddPointsForPost(userID)

	// Send mention notifications (non-critical, async)
	go s.mentionService.SendMentionNotifications(userID, post.ID, req.Content)

	// Load associations
	post, err = s.postRepo.FindByID(post.ID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// CreateReply creates a new reply with sensitive word filtering
func (s *CommunityService) CreateReply(userID, postID uint, req *CreateReplyRequest, authorIPLocation string) (*domain.PostReply, error) {
	// Check if post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, domain.ErrPostNotFound
	}

	// Validate request
	reply := &domain.PostReply{
		PostID:            postID,
		UserID:            userID,
		Content:           req.Content,
		ParentReplyID:     req.ParentReplyID,
		AuthorIPLocation:  authorIPLocation,
	}

	if err := reply.Validate(); err != nil {
		return nil, err
	}

	// Filter sensitive words
	filteredContent, _ := s.filterService.FilterTextAdvanced(req.Content)
	reply.Content = filteredContent

	// Create reply in transaction
	err = s.replyRepo.DB.Transaction(func(tx *gorm.DB) error {
		// Create reply using tx (transaction connection)
		if err := tx.Create(reply).Error; err != nil {
			return err
		}

		// Increment post reply count using tx
		if err := tx.Model(&domain.Post{}).Where("id = ?", postID).Update("reply_count", gorm.Expr("reply_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Add points after successful reply creation (non-critical, can fail independently)
	go s.pointService.AddPointsForReply(userID)

	// Send mention notifications (non-critical, async)
	go s.mentionService.SendMentionNotifications(userID, postID, req.Content)

	// Load associations
	reply, err = s.replyRepo.FindByID(reply.ID)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// TogglePostLike toggles like status for a post (with row lock to prevent race conditions)
func (s *CommunityService) TogglePostLike(userID, postID uint) (string, int, error) {
	// Check if post exists
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return "", 0, domain.ErrPostNotFound
	}

	var action string
	var likeCount int

	err = s.likeRepo.DB.Transaction(func(tx *gorm.DB) error {
		// Lock the post row to prevent race conditions on like toggle
		var lockedPost domain.Post
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedPost, postID).Error; err != nil {
			return err
		}

		// Check if already liked (inside transaction, with row lock)
		var existingLike domain.PostLike
		likeErr := tx.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike).Error

		if likeErr == nil {
			// Already liked → unlike
			if err := tx.Delete(&existingLike).Error; err != nil {
				return err
			}
			if err := tx.Model(&domain.Post{}).Where("id = ?", postID).Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error; err != nil {
				return err
			}
			action = "unliked"
		} else if errors.Is(likeErr, gorm.ErrRecordNotFound) {
			// Not liked → like
			like := &domain.PostLike{
				UserID: userID,
				PostID: postID,
			}
			if err := tx.Create(like).Error; err != nil {
				return err
			}
			if err := tx.Model(&domain.Post{}).Where("id = ?", postID).Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
				return err
			}
			action = "liked"
		} else {
			return likeErr
		}

		return nil
	})

	if err != nil {
		return "", 0, err
	}

	// Add points to post author after successful like (non-critical, can fail independently)
	if action == "liked" {
		go s.pointService.AddPointsForLiked(post.UserID, 2)
	}

	// Get updated like count
	var count int64
	s.likeRepo.DB.Model(&domain.PostLike{}).Where("post_id = ?", postID).Count(&count)
	likeCount = int(count)

	return action, likeCount, nil
}

// ToggleReplyLike toggles like status for a reply (with row lock)
func (s *CommunityService) ToggleReplyLike(userID, replyID uint) (string, int, error) {
	// Check if reply exists
	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return "", 0, domain.ErrReplyNotFound
	}

	var action string
	var likeCount int

	err = s.replyLikeRepo.DB.Transaction(func(tx *gorm.DB) error {
		// Lock the reply row to prevent race conditions on like toggle
		var lockedReply domain.PostReply
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedReply, replyID).Error; err != nil {
			return err
		}

		// Check if already liked (inside transaction, with lock)
		var existingLike domain.PostReplyLike
		likeErr := tx.Where("user_id = ? AND reply_id = ?", userID, replyID).First(&existingLike).Error

		if likeErr == nil {
			// Already liked → unlike
			if err := tx.Delete(&existingLike).Error; err != nil {
				return err
			}
			if err := tx.Model(&domain.PostReply{}).Where("id = ?", replyID).Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error; err != nil {
				return err
			}
			action = "unliked"
		} else if errors.Is(likeErr, gorm.ErrRecordNotFound) {
			// Not liked → like
			like := &domain.PostReplyLike{
				UserID: userID,
				ReplyID: replyID,
			}
			if err := tx.Create(like).Error; err != nil {
				return err
			}
			if err := tx.Model(&domain.PostReply{}).Where("id = ?", replyID).Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
				return err
			}
			action = "liked"
		} else {
			return likeErr
		}

		return nil
	})

	if err != nil {
		return "", 0, err
	}

	// Add points to reply author after successful like (non-critical, can fail independently)
	if action == "liked" {
		go s.pointService.AddPointsForLiked(reply.UserID, 1)
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
		Reason:     reason,
		Status:     string(domain.ReportStatusPending),
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
		Reason:     reason,
		Status:     string(domain.ReportStatusPending),
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

// UpdatePost updates a post (only by the owner)
func (s *CommunityService) UpdatePost(userID, postID uint, req *UpdatePostRequest) (*domain.Post, error) {
	// Find the post
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, domain.ErrPostNotFound
	}

	// Check ownership
	if post.UserID != userID {
		return nil, domain.ErrPermissionDenied
	}

	// Build updates map
	updates := map[string]interface{}{}

	if req.Title != nil {
		filteredTitle, _ := s.filterService.FilterTextAdvanced(*req.Title)
		updates["title"] = filteredTitle
	}
	if req.Content != nil {
		filteredContent, _ := s.filterService.FilterTextAdvanced(*req.Content)
		updates["content"] = filteredContent
	}
	if req.Category != nil {
		// Validate category exists
		activeKeys, err := s.categoryRepo.GetActiveKeys()
		if err != nil {
			return nil, err
		}
		if !slices.Contains(activeKeys, *req.Category) {
			return nil, domain.ErrInvalidCategory
		}
		updates["category"] = *req.Category
	}
	if req.PostType != nil {
		validPostTypes := []domain.PostType{domain.PostTypeDynamic, domain.PostTypeArticle, domain.PostTypeQuestion, domain.PostTypeSuggest, domain.PostTypeQuick}
		valid := false
		for _, t := range validPostTypes {
			if domain.PostType(*req.PostType) == t {
				valid = true
				break
			}
		}
		if !valid {
			return nil, domain.ErrInvalidCategory
		}
		updates["post_type"] = *req.PostType
	}
	if req.Visibility != nil {
		validVisibilities := []domain.PostVisibility{domain.VisibilityPublic, domain.VisibilityMember, domain.VisibilityVIP, domain.VisibilityAdmin}
		valid := false
		for _, v := range validVisibilities {
			if domain.PostVisibility(*req.Visibility) == v {
				valid = true
				break
			}
		}
		if !valid {
			return nil, domain.ErrInvalidVisibility
		}
		updates["visibility"] = *req.Visibility
	}

	if len(updates) == 0 && req.Tags == nil && req.Topics == nil {
		return post, nil // Nothing to update
	}

	if len(updates) > 0 {
		if err := s.postRepo.Update(postID, updates); err != nil {
			return nil, err
		}
	}

	// Handle tags update
	if req.Tags != nil {
		if len(req.Tags) == 0 {
			// Remove all tags
			if err := s.postTagRepo.RemoveAllTagsFromPost(postID); err != nil {
				// Log error but don't fail update
				_ = err
			}
		} else {
			tags, err := s.tagRepo.FindOrCreateBatch(req.Tags)
			if err != nil {
				// Log error but don't fail update
				_ = err
			} else {
				if err := s.postTagRepo.AddTagsToPost(postID, tags); err != nil {
					// Log error but don't fail update
					_ = err
				}
			}
		}
	}

	// Handle topics update
	if req.Topics != nil {
		// Get old topics to decrement count
		oldTopics, _ := s.postTopicRepo.GetTopicsByPostID(postID)

		if len(req.Topics) == 0 {
			// Remove all topics
			if err := s.postTopicRepo.RemoveAllTopicsFromPost(postID); err != nil {
				// Log error but don't fail update
				_ = err
			}
			// Decrement post count for old topics
			for _, topic := range oldTopics {
				go s.topicRepo.DecrementPostCount(topic.ID)
			}
		} else {
			topics, err := s.topicRepo.FindOrCreateBatch(req.Topics)
			if err != nil {
				// Log error but don't fail update
				_ = err
			} else {
				if err := s.postTopicRepo.AddTopicsToPost(postID, topics); err != nil {
					// Log error but don't fail update
					_ = err
				}
				// Decrement post count for old topics
				for _, topic := range oldTopics {
					go s.topicRepo.DecrementPostCount(topic.ID)
				}
				// Increment post count for new topics
				for _, topic := range topics {
					go s.topicRepo.IncrementPostCount(topic.ID)
				}
			}
		}
	}

	// Reload post
	return s.postRepo.FindByID(postID)
}

// DeletePost deletes a post (only by the owner)
func (s *CommunityService) DeletePost(userID, postID uint) error {
	// Find the post
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return domain.ErrPostNotFound
	}

	// Check ownership
	if post.UserID != userID {
		return domain.ErrPermissionDenied
	}

	// 异步清理帖子内容中引用的上传文件
	go s.cleanupPostFiles(post.Content)

	return s.postRepo.Delete(postID)
}

// AdminDeletePost 管理员删除帖子（无需验证所有权）
func (s *CommunityService) AdminDeletePost(postID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return domain.ErrPostNotFound
	}

	// 异步清理帖子内容中引用的上传文件
	go s.cleanupPostFiles(post.Content)

	return s.postRepo.Delete(postID)
}

// cleanupPostFiles 从帖子内容中提取上传文件 URL 并删除
// 仅删除匹配本系统上传路径的文件，忽略外部 URL
func (s *CommunityService) cleanupPostFiles(content string) {
	if s.storage == nil || content == "" {
		return
	}

	// 匹配 /uploads/ 开头的本地上传路径
	// 也匹配包含 uploads/ 的完整 CDN URL
	re := regexp.MustCompile(`(?:https?://[^\s"'()\[\]]+)?/uploads/(?:photos|covers|videos)/[^\s"'()\[\]]+\.\w+`)
	matches := re.FindAllString(content, -1)

	// 去重
	seen := make(map[string]bool)
	for _, url := range matches {
		if seen[url] {
			continue
		}
		seen[url] = true
		if err := s.storage.Delete(url); err != nil {
			logger.Warn("清理帖子文件失败", "url", url, "error", err)
		}
	}
}

// UpdateReply updates a reply (only by the owner)
func (s *CommunityService) UpdateReply(userID, replyID uint, req *UpdateReplyRequest) (*domain.PostReply, error) {
	// Find the reply
	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return nil, domain.ErrReplyNotFound
	}

	// Check ownership
	if reply.UserID != userID {
		return nil, domain.ErrPermissionDenied
	}

	// Filter sensitive words
	filteredContent, _ := s.filterService.FilterTextAdvanced(req.Content)

	updates := map[string]interface{}{
		"content": filteredContent,
	}

	if err := s.replyRepo.Update(replyID, updates); err != nil {
		return nil, err
	}

	// Reload reply
	return s.replyRepo.FindByID(replyID)
}

// DeleteReply deletes a reply (only by the owner)
func (s *CommunityService) DeleteReply(userID, replyID uint) error {
	// Find the reply
	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return domain.ErrReplyNotFound
	}

	// Check ownership
	if reply.UserID != userID {
		return domain.ErrPermissionDenied
	}

	// Delete reply and decrement post reply count
	if err := s.replyRepo.Delete(replyID); err != nil {
		return err
	}

	// Decrement post reply count (non-critical)
	s.postRepo.DecrementReplyCount(reply.PostID)

	return nil
}

// SaveDraft saves a draft (create or update)
func (s *CommunityService) SaveDraft(userID uint, req *SaveDraftRequest) (*domain.Draft, error) {
	// Limit to 20 drafts per user
	count, err := s.draftRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	if req.ID != nil {
		// Update existing draft
		draft, err := s.draftRepo.FindByID(*req.ID)
		if err != nil {
			return nil, err
		}
		// Check ownership
		if draft.UserID != userID {
			return nil, domain.ErrPermissionDenied
		}

		updates := map[string]interface{}{}
		if req.Title != "" {
			updates["title"] = req.Title
		}
		if req.Content != "" {
			updates["content"] = req.Content
		}
		if req.Category != "" {
			updates["category"] = req.Category
		}
		if req.PostType != "" {
			updates["post_type"] = req.PostType
		}
		if req.Visibility != "" {
			updates["visibility"] = req.Visibility
		}

		if len(updates) > 0 {
			if err := s.draftRepo.Update(*req.ID, updates); err != nil {
				return nil, err
			}
		}

		return s.draftRepo.FindByID(*req.ID)
	}

	// Create new draft
	if count >= 20 {
		return nil, domain.ErrDraftLimitReached
	}

	draft := &domain.Draft{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Category:   req.Category,
		PostType:   req.PostType,
		Visibility: req.Visibility,
	}

	if draft.Category == "" {
		draft.Category = "discussion"
	}
	if draft.PostType == "" {
		draft.PostType = "dynamic"
	}
	if draft.Visibility == "" {
		draft.Visibility = "public"
	}

	if err := s.draftRepo.Create(draft); err != nil {
		return nil, err
	}

	return draft, nil
}

// GetDrafts gets user's drafts with pagination
func (s *CommunityService) GetDrafts(userID uint, page, pageSize int) ([]domain.Draft, int64, error) {
	return s.draftRepo.ListByUserID(userID, page, pageSize)
}

// DeleteDraft deletes a draft (only by the owner)
func (s *CommunityService) DeleteDraft(userID, draftID uint) error {
	draft, err := s.draftRepo.FindByID(draftID)
	if err != nil {
		return domain.ErrPostNotFound // reuse error
	}

	if draft.UserID != userID {
		return domain.ErrPermissionDenied
	}

	return s.draftRepo.Delete(draftID)
}

// GetPopularTags returns popular tags
func (s *CommunityService) GetPopularTags(limit int) ([]domain.Tag, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.tagRepo.GetPopularTags(limit)
}

// SearchTags searches tags by keyword
func (s *CommunityService) SearchTags(keyword string, limit int) ([]domain.Tag, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.tagRepo.SearchByName(keyword, limit)
}

// GetPostsByTagName gets posts by tag name
func (s *CommunityService) GetPostsByTagName(tagName string, page, pageSize int) ([]domain.Post, int64, error) {
	tag, err := s.tagRepo.FindByName(tagName)
	if err != nil {
		return nil, 0, domain.ErrTagNotFound
	}

	postIDs, err := s.postTagRepo.GetPostIDsByTagID(tag.ID)
	if err != nil {
		return nil, 0, err
	}

	if len(postIDs) == 0 {
		return []domain.Post{}, 0, nil
	}

	// Get posts with pagination
	var posts []domain.Post
	var total int64

	offset := (page - 1) * pageSize
	err = s.postRepo.DB.
		Where("id IN ?", postIDs).
		Preload("User").
		Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.postRepo.DB.
		Model(&domain.Post{}).
		Where("id IN ?", postIDs).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetHotTopics returns hot topics
func (s *CommunityService) GetHotTopics(limit int) ([]domain.Topic, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.topicRepo.GetHotTopics(limit)
}

// SearchTopics searches topics by keyword
func (s *CommunityService) SearchTopics(keyword string, limit int) ([]domain.Topic, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.topicRepo.SearchByName(keyword, limit)
}

// GetPostsByTopicName gets posts by topic name
func (s *CommunityService) GetPostsByTopicName(topicName string, page, pageSize int) ([]domain.Post, int64, error) {
	topic, err := s.topicRepo.FindByName(topicName)
	if err != nil {
		return nil, 0, domain.ErrTopicNotFound
	}

	postIDs, err := s.postTopicRepo.GetPostIDsByTopicID(topic.ID)
	if err != nil {
		return nil, 0, err
	}

	if len(postIDs) == 0 {
		return []domain.Post{}, 0, nil
	}

	// Get posts with pagination
	var posts []domain.Post
	var total int64

	offset := (page - 1) * pageSize
	err = s.postRepo.DB.
		Where("id IN ?", postIDs).
		Preload("User").
		Preload("Tags").
		Preload("Topics").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.postRepo.DB.
		Model(&domain.Post{}).
		Where("id IN ?", postIDs).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// Request types
type CreatePostRequest struct {
	Title       string   `json:"title" binding:""`
	Content     string   `json:"content" binding:"required"`
	PhotosetID  *uint    `json:"photoset_id"`
	Category    string   `json:"category"`
	PostType    string   `json:"post_type"`
	Visibility  string   `json:"visibility"`
	IsOriginal  bool     `json:"is_original"`
	ScheduledAt *string  `json:"scheduled_at"`
	Tags        []string `json:"tags"`
	Topics      []string `json:"topics"`
}

type UpdatePostRequest struct {
	Title      *string  `json:"title"`
	Content    *string  `json:"content"`
	Category   *string  `json:"category"`
	PostType   *string  `json:"post_type"`
	Visibility *string  `json:"visibility"`
	Tags       []string `json:"tags"`
	Topics     []string `json:"topics"`
}

type CreateReplyRequest struct {
	Content       string `json:"content" binding:"required"`
	ParentReplyID *uint  `json:"parent_reply_id"`
}

type UpdateReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

type SaveDraftRequest struct {
	ID         *uint  `json:"id"`         // nil for create, non-nil for update
	Title      string `json:"title"`
	Content    string `json:"content" binding:"required"`
	Category   string `json:"category"`
	PostType   string `json:"post_type"`
	Visibility string `json:"visibility"`
}

// RecordShare records a share action and increments post share count atomically
func (s *CommunityService) RecordShare(userID, postID uint, platform string) error {
	// Validate platform
	validPlatforms := map[string]bool{
		"wechat": true,
		"weibo":  true,
		"link":   true,
		"other":  true,
	}
	if !validPlatforms[platform] {
		platform = "other"
	}

	// Check post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		return domain.ErrPostNotFound
	}

	// Record share and increment count in transaction
	return s.shareRepo.DB.Transaction(func(tx *gorm.DB) error {
		share := &domain.PostShare{
			UserID:   userID,
			PostID:   postID,
			Platform: platform,
		}
		if err := tx.Create(share).Error; err != nil {
			return err
		}
		// Atomic increment
		if err := tx.Model(&domain.Post{}).Where("id = ?", postID).Update("share_count", gorm.Expr("share_count + 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetShareStats returns total and per-platform share counts for a post
func (s *CommunityService) GetShareStats(postID uint) (int64, map[string]int64, error) {
	// Check post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		return 0, nil, domain.ErrPostNotFound
	}

	total, err := s.shareRepo.CountByPostID(postID)
	if err != nil {
		return 0, nil, err
	}

	platformCounts, err := s.shareRepo.CountByPostIDAndPlatform(postID)
	if err != nil {
		return 0, nil, err
	}

	return total, platformCounts, nil
}
