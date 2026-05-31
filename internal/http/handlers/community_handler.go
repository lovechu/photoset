package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/domain"
	"photoset/internal/http/middleware"
	"photoset/internal/logger"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommunityHandler handles community-related requests
type CommunityHandler struct {
	communityService      *service.CommunityService
	pointService          *service.PointService
	hotPostsService       *service.HotPostsService
	recommendationService *service.RecommendationService
	postRepo              *repository.PostRepository
	replyRepo             *repository.PostReplyRepository
	likeRepo              *repository.PostLikeRepository
	replyLikeRepo         *repository.PostReplyLikeRepository
	pointRepo             *repository.UserPointRepository
	reportRepo            *repository.PostReportRepository
	categoryRepo          *repository.PostCategoryRepository
	followRepo            *repository.FollowRepository
	postFavoriteRepo      *repository.PostFavoriteRepository
	draftRepo             *repository.DraftRepository
	userRepo              repository.UserRepository
}

// NewCommunityHandler creates a new CommunityHandler
func NewCommunityHandler(
	db *gorm.DB,
	communityService *service.CommunityService,
	pointService *service.PointService,
	hotPostsService *service.HotPostsService,
	recommendationService *service.RecommendationService,
) *CommunityHandler {
	return &CommunityHandler{
		communityService:      communityService,
		pointService:          pointService,
		hotPostsService:       hotPostsService,
		recommendationService: recommendationService,
		postRepo:              repository.NewPostRepository(db),
		replyRepo:             repository.NewPostReplyRepository(db),
		likeRepo:              repository.NewPostLikeRepository(db),
		replyLikeRepo:         repository.NewPostReplyLikeRepository(db),
		pointRepo:             repository.NewUserPointRepository(db),
		reportRepo:            repository.NewPostReportRepository(db),
		categoryRepo:          repository.NewPostCategoryRepository(db),
		followRepo:            repository.NewFollowRepository(db),
		postFavoriteRepo:      repository.NewPostFavoriteRepository(db),
		draftRepo:             repository.NewDraftRepository(db),
		userRepo:              repository.NewUserRepository(),
	}
}

// GetPosts gets post list with pagination and filtering
func (h *CommunityHandler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	sortBy := c.DefaultQuery("sort_by", "latest")

	// Get user info (optional)
	userID, _ := middleware.GetUserID(c)
	userRole, _ := middleware.GetUserRole(c)

	// Determine visibility based on user role
	visibility := ""
	if userID == 0 {
		visibility = "public"
	}

	var posts []domain.Post
	var total int64
	var err error

	if sortBy == "hot" {
		posts, total, err = h.hotPostsService.GetHotPosts(page, pageSize)
	} else {
		posts, total, err = h.postRepo.List(page, pageSize, category, visibility, userRole, userID, sortBy)
	}

	if err != nil {
		logger.Error("Failed to get posts", "error", err, "page", page, "category", category, "sort", sortBy)
		response.ServerError(c, "failed to get posts")
		return
	}

	response.Success(c, gin.H{
		"posts": h.postsToResponseList(posts, userID),
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetPostDetail gets post detail and increments view count
func (h *CommunityHandler) GetPostDetail(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	post, err := h.communityService.GetPostDetail(uint(postID))
	if err != nil {
		response.NotFound(c, "post not found")
		return
	}

	// Check visibility permission
	userRole, _ := middleware.GetUserRole(c)
	if !h.canViewPost(post, userRole) {
		response.Forbidden(c, "permission denied to view this post")
		return
	}

	// Get user info for is_liked check
	userID, _ := middleware.GetUserID(c)

	response.Success(c, gin.H{"post": h.postToResponse(*post, userID)})
}

// CreatePost creates a new post
func (h *CommunityHandler) CreatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	var req service.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	post, err := h.communityService.CreatePost(userID, &req)
	if err != nil {
		logger.Warn("Post creation failed", "user_id", userID, "error", err)
		if err == domain.ErrTitleRequired {
			response.BadRequest(c, "title is required")
		} else if err == domain.ErrContentRequired {
			response.BadRequest(c, "content is required")
		} else if err == domain.ErrInvalidCategory {
			response.BadRequest(c, "invalid category")
		} else if err == domain.ErrDailyLimitReached {
			response.Error(c, 400, "daily post limit reached (max 5 posts per day)")
		} else {
			response.ServerError(c, "failed to create post: "+err.Error())
		}
		return
	}
	logger.Info("Post created", "post_id", post.ID, "user_id", userID, "category", req.Category)
	response.Success(c, gin.H{"id": post.ID})
}

// CreateReply creates a reply to a post
func (h *CommunityHandler) CreateReply(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req service.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	reply, err := h.communityService.CreateReply(userID, uint(postID), &req)
	if err != nil {
		logger.Warn("Reply creation failed", "user_id", userID, "post_id", postID, "error", err)
		if err == domain.ErrPostNotFound {
			response.NotFound(c, "post not found")
		} else if err == domain.ErrReplyContentRequired {
			response.BadRequest(c, "reply content is required")
		} else if err == domain.ErrDailyLimitReached {
			response.Error(c, 400, "daily reply limit reached (max 6 replies per day)")
		} else {
			response.ServerError(c, "failed to create reply: "+err.Error())
		}
		return
	}
	logger.Info("Reply created", "reply_id", reply.ID, "post_id", postID, "user_id", userID)

	response.Success(c, gin.H{"reply": reply})
}

// TogglePostLike toggles like status for a post
func (h *CommunityHandler) TogglePostLike(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	action, likeCount, err := h.communityService.TogglePostLike(userID, uint(postID))
	if err != nil {
		if err == domain.ErrPostNotFound {
			response.NotFound(c, "post not found")
		} else {
			response.ServerError(c, "failed to toggle like: "+err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, action, gin.H{"is_liked": action == "liked", "like_count": likeCount})
}

// ToggleReplyLike toggles like status for a reply
func (h *CommunityHandler) ToggleReplyLike(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	replyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	action, likeCount, err := h.communityService.ToggleReplyLike(userID, uint(replyID))
	if err != nil {
		if err == domain.ErrReplyNotFound {
			response.NotFound(c, "reply not found")
		} else {
			response.ServerError(c, "failed to toggle reply like: "+err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, action, gin.H{"is_liked": action == "liked", "like_count": likeCount})
}

// ReportPost reports a post
func (h *CommunityHandler) ReportPost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required")
		return
	}

	if err := h.communityService.ReportPost(userID, uint(postID), req.Reason); err != nil {
		response.ServerError(c, "failed to report post: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "report submitted successfully"})
}

// ReportReply reports a reply
func (h *CommunityHandler) ReportReply(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	replyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required")
		return
	}

	if err := h.communityService.ReportReply(userID, uint(replyID), req.Reason); err != nil {
		response.ServerError(c, "failed to report reply: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "report submitted successfully"})
}

// GetCategories gets available categories from the database
func (h *CommunityHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryRepo.ListCategories()
	if err != nil {
		response.ServerError(c, "failed to get categories")
		return
	}

	// Return simplified format for the public API
	type CategoryItem struct {
		Slug  string `json:"slug"`
		Name  string `json:"name"`
		Color string `json:"color,omitempty"`
	}
	items := make([]CategoryItem, 0, len(categories))
	for _, cat := range categories {
		items = append(items, CategoryItem{
			Slug:  cat.Key,
			Name:  cat.Name,
			Color: cat.Color,
		})
	}

	response.Success(c, gin.H{"categories": items})
}

// GetHotPosts gets hot posts
func (h *CommunityHandler) GetHotPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Get user info (optional) for is_liked check
	userID, _ := middleware.GetUserID(c)

	posts, total, err := h.hotPostsService.GetHotPosts(page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get hot posts")
		return
	}

	response.Success(c, gin.H{
		"posts": h.postsToResponseList(posts, userID),
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetMyPosts gets current user's posts
func (h *CommunityHandler) GetMyPosts(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := h.postRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get your posts")
		return
	}

	response.Success(c, gin.H{
		"posts": h.postsToResponseList(posts, userID),
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetUserPosts gets a specific user's posts
func (h *CommunityHandler) GetUserPosts(c *gin.Context) {
	targetUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Get current user (optional) for is_liked check
	currentUserID, _ := middleware.GetUserID(c)

	posts, total, err := h.postRepo.FindByUserID(uint(targetUserID), page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get user posts")
		return
	}

	response.Success(c, gin.H{
		"posts": h.postsToResponseList(posts, currentUserID),
		"total": total,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetMyReplies gets current user's replies
func (h *CommunityHandler) GetMyReplies(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	replies, total, err := h.replyRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get your replies")
		return
	}

	response.Success(c, gin.H{
		"replies": replies,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetMyPoints gets current user's points and level info
func (h *CommunityHandler) GetMyPoints(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	level, levelName, currentPoints, nextLevelPoints, err := h.pointService.GetLevelInfo(userID)
	if err != nil {
		response.ServerError(c, "failed to get points info")
		return
	}

	response.Success(c, gin.H{
		"points":            currentPoints,
		"level":             level,
		"level_name":        levelName,
		"next_level_points": nextLevelPoints,
	})
}

// GetReplies gets replies for a post
func (h *CommunityHandler) GetReplies(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	replies, err := h.replyRepo.FindByPostID(uint(postID), page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get replies")
		return
	}

	total, _ := h.replyRepo.CountByPostID(uint(postID))

	response.Success(c, gin.H{
		"replies": replies,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// postToResponse converts a domain.Post to a flat gin.H response (with author_id, author_name, is_liked)
func (h *CommunityHandler) postToResponse(post domain.Post, userID uint) gin.H {
	authorName := ""
	authorAvatar := ""
	if post.User.ID != 0 {
		authorName = post.User.Nickname
		authorAvatar = post.User.Avatar
	}

	// Check if current user liked this post
	isLiked := false
	if userID > 0 {
		liked, _ := h.likeRepo.Exists(userID, post.ID)
		isLiked = liked
	}

	// Check if current user is following the author
	isFollowing := false
	if userID > 0 && post.UserID > 0 && userID != post.UserID {
		following, _ := h.followRepo.Exists(userID, post.UserID)
		isFollowing = following
	}

	// Convert tags to response format
	tags := make([]gin.H, len(post.Tags))
	for i, tag := range post.Tags {
		tags[i] = gin.H{
			"id":   tag.ID,
			"name": tag.Name,
		}
	}

	// Convert topics to response format
	topics := make([]gin.H, len(post.Topics))
	for i, topic := range post.Topics {
		topics[i] = gin.H{
			"id":   topic.ID,
			"name": topic.Name,
		}
	}

	return gin.H{
		"id":           post.ID,
		"title":        post.Title,
		"content":      post.Content,
		"category":     post.Category,
		"post_type":    post.PostType,
		"author_id":    post.UserID,
		"author_name":  authorName,
		"author_avatar": authorAvatar,
		"reply_count":  post.ReplyCount,
		"like_count":   post.LikeCount,
		"share_count":  post.ShareCount,
		"view_count":   post.ViewCount,
		"is_pinned":    post.IsPinned,
		"is_essence":   post.IsEssence,
		"is_liked":     isLiked,
		"is_following": isFollowing,
		"status":       post.Status,
		"created_at":   post.CreatedAt,
		"tags":         tags,
		"topics":       topics,
	}
}

// postsToResponseList converts a slice of domain.Post to flat gin.H slices
func (h *CommunityHandler) postsToResponseList(posts []domain.Post, userID uint) []gin.H {
	result := make([]gin.H, len(posts))
	for i, post := range posts {
		result[i] = h.postToResponse(post, userID)
	}
	return result
}

// canViewPost checks if user can view the post based on visibility
func (h *CommunityHandler) canViewPost(post *domain.Post, userRole string) bool {
	if post.Visibility == "public" {
		return true
	}

	if userRole == "" {
		return false
	}

	if post.Visibility == "member" {
		return userRole == "member" || userRole == "creator" || userRole == "vip" || userRole == "admin"
	}

	if post.Visibility == "vip" {
		return userRole == "creator" || userRole == "vip" || userRole == "admin"
	}

	if post.Visibility == "admin" {
		return userRole == "admin"
	}

	return false
}

// GetUserProfile gets a user's public profile info
func (h *CommunityHandler) GetUserProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(uint(userID))
	if err != nil {
		response.ServerError(c, "获取用户信息失败")
		return
	}
	if user == nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// Count user's posts
	var postCount int64
	h.postRepo.DB.Model(&domain.Post{}).Where("user_id = ?", userID).Count(&postCount)

	// Count user's replies
	var replyCount int64
	h.replyRepo.DB.Model(&domain.PostReply{}).Where("user_id = ?", userID).Count(&replyCount)

	// Count user's likes received
	var likeCount int64
	h.likeRepo.DB.Model(&domain.PostLike{}).Where("post_id IN (SELECT id FROM posts WHERE user_id = ?)", userID).Count(&likeCount)

	response.Success(c, gin.H{
		"id":              user.ID,
		"nickname":        user.Nickname,
		"avatar":          user.Avatar,
		"bio":             user.Bio,
		"ip_location":     user.IPLocation,
		"level":           user.Level,
		"following_count": user.FollowingCount,
		"follower_count":  user.FollowerCount,
		"like_count":      likeCount,
		"post_count":      postCount,
		"reply_count":     replyCount,
		"created_at":      user.CreatedAt,
	})
}

// TogglePostFavorite toggles favorite status for a post
func (h *CommunityHandler) TogglePostFavorite(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	// Check if post exists
	_, err = h.postRepo.FindByID(uint(postID))
	if err != nil {
		response.NotFound(c, "post not found")
		return
	}

	isFavorited, err := h.postFavoriteRepo.Toggle(userID, uint(postID))
	if err != nil {
		response.ServerError(c, "failed to toggle favorite")
		return
	}

	message := "收藏成功"
	if !isFavorited {
		message = "已取消收藏"
	}

	response.Success(c, gin.H{
		"is_favorited": isFavorited,
		"message":      message,
	})
}

// CheckPostFavorite checks if user has favorited a post
func (h *CommunityHandler) CheckPostFavorite(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Success(c, gin.H{"is_favorited": false})
		return
	}

	isFavorited, err := h.postFavoriteRepo.IsFavorited(userID, uint(postID))
	if err != nil {
		response.ServerError(c, "failed to check favorite status")
		return
	}

	response.Success(c, gin.H{"is_favorited": isFavorited})
}

// GetMyFavorites gets current user's favorite posts
func (h *CommunityHandler) GetMyFavorites(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	favorites, total, err := h.postFavoriteRepo.List(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get favorites")
		return
	}

	// Convert to response format
	posts := make([]gin.H, len(favorites))
	for i, fav := range favorites {
		posts[i] = h.postToResponse(fav.Post, userID)
	}

	response.Success(c, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// UpdatePost updates a post (only by the owner)
func (h *CommunityHandler) UpdatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req service.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	post, err := h.communityService.UpdatePost(userID, uint(postID), &req)
	if err != nil {
		if err == domain.ErrPostNotFound {
			response.NotFound(c, "post not found")
		} else if err == domain.ErrPermissionDenied {
			response.Forbidden(c, "you can only edit your own posts")
		} else if err == domain.ErrInvalidCategory {
			response.BadRequest(c, "invalid category or post type")
		} else if err == domain.ErrInvalidVisibility {
			response.BadRequest(c, "invalid visibility")
		} else {
			response.ServerError(c, "failed to update post: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"post": h.postToResponse(*post, userID)})
}

// DeletePost deletes a post (only by the owner)
func (h *CommunityHandler) DeletePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	if err := h.communityService.DeletePost(userID, uint(postID)); err != nil {
		if err == domain.ErrPostNotFound {
			response.NotFound(c, "post not found")
		} else if err == domain.ErrPermissionDenied {
			response.Forbidden(c, "you can only delete your own posts")
		} else {
			response.ServerError(c, "failed to delete post: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "post deleted successfully"})
}

// UpdateReply updates a reply (only by the owner)
func (h *CommunityHandler) UpdateReply(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	replyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	var req service.UpdateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	reply, err := h.communityService.UpdateReply(userID, uint(replyID), &req)
	if err != nil {
		if err == domain.ErrReplyNotFound {
			response.NotFound(c, "reply not found")
		} else if err == domain.ErrPermissionDenied {
			response.Forbidden(c, "you can only edit your own replies")
		} else if err == domain.ErrReplyContentRequired {
			response.BadRequest(c, "reply content is required")
		} else {
			response.ServerError(c, "failed to update reply: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"reply": reply})
}

// DeleteReply deletes a reply (only by the owner)
func (h *CommunityHandler) DeleteReply(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	replyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	if err := h.communityService.DeleteReply(userID, uint(replyID)); err != nil {
		if err == domain.ErrReplyNotFound {
			response.NotFound(c, "reply not found")
		} else if err == domain.ErrPermissionDenied {
			response.Forbidden(c, "you can only delete your own replies")
		} else {
			response.ServerError(c, "failed to delete reply: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "reply deleted successfully"})
}

// SaveDraft saves a draft (create or update)
func (h *CommunityHandler) SaveDraft(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	var req service.SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	draft, err := h.communityService.SaveDraft(userID, &req)
	if err != nil {
		if err == domain.ErrPermissionDenied {
			response.Forbidden(c, "you can only edit your own drafts")
		} else if err == domain.ErrDraftLimitReached {
			response.Error(c, 400, "draft limit reached (max 20)")
		} else {
			response.ServerError(c, "failed to save draft: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"draft": draft})
}

// GetMyDrafts gets current user's drafts
func (h *CommunityHandler) GetMyDrafts(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	drafts, total, err := h.communityService.GetDrafts(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get drafts")
		return
	}

	response.Success(c, gin.H{
		"drafts": drafts,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// DeleteDraft deletes a draft
func (h *CommunityHandler) DeleteDraft(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	draftID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid draft id")
		return
	}

	if err := h.communityService.DeleteDraft(userID, uint(draftID)); err != nil {
		if err == domain.ErrPostNotFound {
			response.NotFound(c, "draft not found")
		} else if err == domain.ErrPermissionDenied {
			response.Forbidden(c, "you can only delete your own drafts")
		} else {
			response.ServerError(c, "failed to delete draft: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "draft deleted successfully"})
}

// SearchPosts searches posts by keyword
func (h *CommunityHandler) SearchPosts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "keyword is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Get user info (optional)
	userID, _ := middleware.GetUserID(c)
	userRole, _ := middleware.GetUserRole(c)

	posts, total, err := h.postRepo.Search(keyword, page, pageSize, userID, userRole)
	if err != nil {
		response.ServerError(c, "failed to search posts")
		return
	}

	response.Success(c, gin.H{
		"posts": h.postsToResponseList(posts, userID),
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// SearchUsersForMention searches users for @mention autocomplete
func (h *CommunityHandler) SearchUsersForMention(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "keyword is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	users, err := h.userRepo.SearchByNickname(keyword, limit)
	if err != nil {
		response.ServerError(c, "failed to search users")
		return
	}

	// Format response
	result := make([]gin.H, len(users))
	for i, user := range users {
		result[i] = gin.H{
			"id":       user.ID,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		}
	}

	response.Success(c, gin.H{"users": result})
}

// RecordShare records a share action for a post
func (h *CommunityHandler) RecordShare(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req struct {
		Platform string `json:"platform"`
	}
	// Platform is optional, defaults to "other"
	_ = c.ShouldBindJSON(&req)

	if req.Platform == "" {
		req.Platform = "other"
	}

	if err := h.communityService.RecordShare(userID, uint(postID), req.Platform); err != nil {
		if err == domain.ErrPostNotFound {
			response.NotFound(c, "post not found")
		} else {
			response.ServerError(c, "failed to record share: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "share recorded"})
}

// GetShareStats gets share statistics for a post
func (h *CommunityHandler) GetShareStats(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	total, platformCounts, err := h.communityService.GetShareStats(uint(postID))
	if err != nil {
		if err == domain.ErrPostNotFound {
			response.NotFound(c, "post not found")
		} else {
			response.ServerError(c, "failed to get share stats")
		}
		return
	}

	response.Success(c, gin.H{
		"total":           total,
		"platform_counts": platformCounts,
	})
}

// GetPopularTags gets popular tags
func (h *CommunityHandler) GetPopularTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	tags, err := h.communityService.GetPopularTags(limit)
	if err != nil {
		response.ServerError(c, "failed to get popular tags")
		return
	}

	// Convert to response format
	result := make([]gin.H, len(tags))
	for i, tag := range tags {
		result[i] = gin.H{
			"id":   tag.ID,
			"name": tag.Name,
		}
	}

	response.Success(c, gin.H{"tags": result})
}

// SearchTags searches tags by keyword
func (h *CommunityHandler) SearchTags(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "keyword is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tags, err := h.communityService.SearchTags(keyword, limit)
	if err != nil {
		response.ServerError(c, "failed to search tags")
		return
	}

	// Convert to response format
	result := make([]gin.H, len(tags))
	for i, tag := range tags {
		result[i] = gin.H{
			"id":   tag.ID,
			"name": tag.Name,
		}
	}

	response.Success(c, gin.H{"tags": result})
}

// GetPostsByTagName gets posts by tag name
func (h *CommunityHandler) GetPostsByTagName(c *gin.Context) {
	tagName := c.Param("name")
	if tagName == "" {
		response.BadRequest(c, "tag name is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Get user info (optional)
	userID, _ := middleware.GetUserID(c)

	posts, total, err := h.communityService.GetPostsByTagName(tagName, page, pageSize)
	if err != nil {
		if err == domain.ErrTagNotFound {
			response.NotFound(c, "tag not found")
		} else {
			response.ServerError(c, "failed to get posts by tag")
		}
		return
	}

	response.Success(c, gin.H{
		"posts":     h.postsToResponseList(posts, userID),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetHotTopics gets hot topics
func (h *CommunityHandler) GetHotTopics(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	topics, err := h.communityService.GetHotTopics(limit)
	if err != nil {
		response.ServerError(c, "failed to get hot topics")
		return
	}

	// Convert to response format
	result := make([]gin.H, len(topics))
	for i, topic := range topics {
		result[i] = gin.H{
			"id":         topic.ID,
			"name":       topic.Name,
			"cover":      topic.Cover,
			"description": topic.Description,
			"post_count": topic.PostCount,
			"is_hot":     topic.IsHot,
		}
	}

	response.Success(c, gin.H{"topics": result})
}

// SearchTopics searches topics by keyword
func (h *CommunityHandler) SearchTopics(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "keyword is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	topics, err := h.communityService.SearchTopics(keyword, limit)
	if err != nil {
		response.ServerError(c, "failed to search topics")
		return
	}

	// Convert to response format
	result := make([]gin.H, len(topics))
	for i, topic := range topics {
		result[i] = gin.H{
			"id":         topic.ID,
			"name":       topic.Name,
			"cover":      topic.Cover,
			"description": topic.Description,
			"post_count": topic.PostCount,
			"is_hot":     topic.IsHot,
		}
	}

	response.Success(c, gin.H{"topics": result})
}

// GetPostsByTopicName gets posts by topic name
func (h *CommunityHandler) GetPostsByTopicName(c *gin.Context) {
	topicName := c.Param("name")
	if topicName == "" {
		response.BadRequest(c, "topic name is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Get user info (optional)
	userID, _ := middleware.GetUserID(c)

	posts, total, err := h.communityService.GetPostsByTopicName(topicName, page, pageSize)
	if err != nil {
		if err == domain.ErrTopicNotFound {
			response.NotFound(c, "topic not found")
		} else {
			response.ServerError(c, "failed to get posts by topic")
		}
		return
	}

	response.Success(c, gin.H{
		"posts":     h.postsToResponseList(posts, userID),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetRecommendations gets recommended posts based on algorithm type
func (h *CommunityHandler) GetRecommendations(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	recType := service.RecommendationType(c.DefaultQuery("type", "mixed"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req := &service.RecommendationRequest{
		UserID:   userID,
		Type:     recType,
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.recommendationService.GetRecommendations(req)
	if err != nil {
		response.ServerError(c, "failed to get recommendations: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"posts":     h.postsToResponseList(result.Posts, userID),
		"total":     result.Total,
		"type":      result.Type,
		"algorithm": result.Algorithm,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     result.Total,
		},
	})
}

// GetUserInterestProfile gets user's interest profile (tags and topics)
func (h *CommunityHandler) GetUserInterestProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "please login first")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tags, err := h.recommendationService.GetUserInterestTags(userID, limit)
	if err != nil {
		response.ServerError(c, "failed to get interest tags")
		return
	}

	topics, err := h.recommendationService.GetUserInterestTopics(userID, limit)
	if err != nil {
		response.ServerError(c, "failed to get interest topics")
		return
	}

	tagResult := make([]gin.H, len(tags))
	for i, tag := range tags {
		tagResult[i] = gin.H{
			"id":   tag.ID,
			"name": tag.Name,
		}
	}

	topicResult := make([]gin.H, len(topics))
	for i, topic := range topics {
		topicResult[i] = gin.H{
			"id":         topic.ID,
			"name":       topic.Name,
			"cover":      topic.Cover,
			"description": topic.Description,
			"post_count": topic.PostCount,
		}
	}

	response.Success(c, gin.H{
		"interest_tags":   tagResult,
		"interest_topics": topicResult,
	})
}
