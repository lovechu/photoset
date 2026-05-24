package handlers

import (
	"strconv"

	"photoset/internal/domain"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommunityHandler handles community-related requests
type CommunityHandler struct {
	communityService  *service.CommunityService
	pointService     *service.PointService
	hotPostsService  *service.HotPostsService
	postRepo         *repository.PostRepository
	replyRepo        *repository.PostReplyRepository
	likeRepo         *repository.PostLikeRepository
	replyLikeRepo    *repository.PostReplyLikeRepository
	pointRepo        *repository.UserPointRepository
	reportRepo       *repository.PostReportRepository
}

// NewCommunityHandler creates a new CommunityHandler
func NewCommunityHandler(
	db *gorm.DB,
	communityService *service.CommunityService,
	pointService *service.PointService,
	hotPostsService *service.HotPostsService,
) *CommunityHandler {
	return &CommunityHandler{
		communityService: communityService,
		pointService:    pointService,
		hotPostsService: hotPostsService,
		postRepo:        repository.NewPostRepository(db),
		replyRepo:       repository.NewPostReplyRepository(db),
		likeRepo:        repository.NewPostLikeRepository(db),
		replyLikeRepo:   repository.NewPostReplyLikeRepository(db),
		pointRepo:       repository.NewUserPointRepository(db),
		reportRepo:      repository.NewPostReportRepository(db),
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
		response.ServerError(c, "failed to get posts")
		return
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

	response.Success(c, gin.H{"post": post})
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
		if err == domain.ErrTitleRequired {
			response.BadRequest(c, "title is required")
		} else if err == domain.ErrContentRequired {
			response.BadRequest(c, "content is required")
		} else if err == domain.ErrDailyLimitReached {
			response.Error(c, 400, "daily post limit reached (max 5 posts per day)")
		} else {
			response.ServerError(c, "failed to create post: "+err.Error())
		}
		return
	}

	response.Success(c, gin.H{"post": post})
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

	response.SuccessWithMessage(c, action, gin.H{"like_count": likeCount})
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

	response.SuccessWithMessage(c, action, gin.H{"like_count": likeCount})
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

// GetCategories gets available categories
func (h *CommunityHandler) GetCategories(c *gin.Context) {
	categories := []gin.H{
		{"slug": "discussion", "name": "Discussion"},
		{"slug": "qa", "name": "Q&A"},
		{"slug": "showcase", "name": "Showcase"},
		{"slug": "suggestion", "name": "Suggestion"},
	}
	response.Success(c, gin.H{"categories": categories})
}

// GetHotPosts gets hot posts
func (h *CommunityHandler) GetHotPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := h.hotPostsService.GetHotPosts(page, pageSize)
	if err != nil {
		response.ServerError(c, "failed to get hot posts")
		return
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
		"posts": posts,
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

	response.Success(c, gin.H{
		"replies": replies,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
		},
	})
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
