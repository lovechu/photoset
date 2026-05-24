package admin

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

// AdminCommunityHandler handles admin community requests
type AdminCommunityHandler struct {
	postRepo        *repository.PostRepository
	replyRepo       *repository.PostReplyRepository
	likeRepo        *repository.PostLikeRepository
	replyLikeRepo   *repository.PostReplyLikeRepository
	pointRepo       *repository.UserPointRepository
	wordRepo        *repository.SensitiveWordRepository
	reportRepo      *repository.PostReportRepository
	pointService    *service.PointService
	filterService   *service.SensitiveFilterService
}

// NewAdminCommunityHandler creates a new AdminCommunityHandler
func NewAdminCommunityHandler(db *gorm.DB) *AdminCommunityHandler {
	return &AdminCommunityHandler{
		postRepo:      repository.NewPostRepository(db),
		replyRepo:     repository.NewPostReplyRepository(db),
		likeRepo:      repository.NewPostLikeRepository(db),
		replyLikeRepo: repository.NewPostReplyLikeRepository(db),
		pointRepo:     repository.NewUserPointRepository(db),
		wordRepo:      repository.NewSensitiveWordRepository(db),
		reportRepo:    repository.NewPostReportRepository(db),
		pointService:  service.NewPointService(repository.NewUserPointRepository(db)),
		filterService: service.NewSensitiveFilterService(repository.NewSensitiveWordRepository(db)),
	}
}

// ============ Post Management ============

// GetPosts gets all posts (admin)
func (h *AdminCommunityHandler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	posts, total, err := h.postRepo.ListForAdmin(page, pageSize, status)
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

// TogglePin toggles pin status of a post
func (h *AdminCommunityHandler) TogglePin(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	if err := h.postRepo.TogglePin(uint(postID)); err != nil {
		response.ServerError(c, "failed to toggle pin")
		return
	}

	response.Success(c, gin.H{"message": "pin status toggled successfully"})
}

// UpdateStatus updates post status
func (h *AdminCommunityHandler) UpdateStatus(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "status is required")
		return
	}

	if err := h.postRepo.UpdateStatus(uint(postID), req.Status); err != nil {
		response.ServerError(c, "failed to update status")
		return
	}

	response.Success(c, gin.H{"message": "status updated successfully"})
}

// DeletePost deletes a post (hard delete)
func (h *AdminCommunityHandler) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	// Get post author to deduct points
	post, err := h.postRepo.FindByID(uint(postID))
	if err == nil {
		// Deduct points from author
		h.pointService.DeductPointsForDelete(post.UserID)
	}

	// Delete post (cascade will delete related records)
	if err := h.postRepo.Delete(uint(postID)); err != nil {
		response.ServerError(c, "failed to delete post")
		return
	}

	response.Success(c, gin.H{"message": "post deleted successfully"})
}

// ============ Reply Management ============

// GetReplies gets all replies (admin)
func (h *AdminCommunityHandler) GetReplies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	postIDStr := c.Query("post_id")

	// TODO: Implement admin reply list with filtering
	_ = postIDStr

	response.Success(c, gin.H{
		"replies": []interface{}{},
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     0,
		},
	})
}

// DeleteReply deletes a reply (hard delete)
func (h *AdminCommunityHandler) DeleteReply(c *gin.Context) {
	replyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	// Get reply to decrement post reply count
	reply, err := h.replyRepo.FindByID(uint(replyID))
	if err == nil {
		h.postRepo.DecrementReplyCount(reply.PostID)
	}

	// Delete reply
	if err := h.replyRepo.Delete(uint(replyID)); err != nil {
		response.ServerError(c, "failed to delete reply")
		return
	}

	response.Success(c, gin.H{"message": "reply deleted successfully"})
}

// ============ Sensitive Word Management ============

// GetKeywords gets all sensitive words
func (h *AdminCommunityHandler) GetKeywords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	isActiveStr := c.Query("is_active")

	var isActive *bool
	if isActiveStr != "" {
		v := isActiveStr == "true"
		isActive = &v
	}

	words, total, err := h.wordRepo.List(page, pageSize, isActive)
	if err != nil {
		response.ServerError(c, "failed to get keywords")
		return
	}

	response.Success(c, gin.H{
		"keywords": words,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// AddKeyword adds a new sensitive word
func (h *AdminCommunityHandler) AddKeyword(c *gin.Context) {
	var req struct {
		Word       string `json:"word" binding:"required"`
		Replacement string `json:"replacement"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "word is required")
		return
	}

	if req.Replacement == "" {
		req.Replacement = "***"
	}

	word := &domain.SensitiveWord{
		Word:       req.Word,
		Replacement: req.Replacement,
		IsActive:   true,
	}

	if err := h.wordRepo.Create(word); err != nil {
		response.ServerError(c, "failed to add keyword")
		return
	}

	// Reload sensitive words
	service.InitSensitiveWords(h.wordRepo)

	response.Success(c, gin.H{"keyword": word})
}

// UpdateKeyword updates a sensitive word
func (h *AdminCommunityHandler) UpdateKeyword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid keyword id")
		return
	}

	var req struct {
		Word       string `json:"word"`
		Replacement string `json:"replacement"`
		IsActive   *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	updates := map[string]interface{}{}
	if req.Word != "" {
		updates["word"] = req.Word
	}
	if req.Replacement != "" {
		updates["replacement"] = req.Replacement
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.wordRepo.Update(uint(id), updates); err != nil {
		response.ServerError(c, "failed to update keyword")
		return
	}

	// Reload sensitive words
	service.InitSensitiveWords(h.wordRepo)

	response.Success(c, gin.H{"message": "keyword updated successfully"})
}

// DeleteKeyword deletes a sensitive word
func (h *AdminCommunityHandler) DeleteKeyword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid keyword id")
		return
	}

	if err := h.wordRepo.Delete(uint(id)); err != nil {
		response.ServerError(c, "failed to delete keyword")
		return
	}

	// Reload sensitive words
	service.InitSensitiveWords(h.wordRepo)

	response.Success(c, gin.H{"message": "keyword deleted successfully"})
}

// ReloadKeywords reloads sensitive words from database
func (h *AdminCommunityHandler) ReloadKeywords(c *gin.Context) {
	if err := service.InitSensitiveWords(h.wordRepo); err != nil {
		response.ServerError(c, "failed to reload keywords")
		return
	}

	response.Success(c, gin.H{"message": "keywords reloaded successfully"})
}

// ============ Report Management ============

// GetReports gets all reports
func (h *AdminCommunityHandler) GetReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	reports, total, err := h.reportRepo.List(page, pageSize, status)
	if err != nil {
		response.ServerError(c, "failed to get reports")
		return
	}

	response.Success(c, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// ResolveReport resolves or rejects a report
func (h *AdminCommunityHandler) ResolveReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid report id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "status is required")
		return
	}

	// Get handler ID from context
	handlerID, _ := middleware.GetUserID(c)

	if err := h.reportRepo.Process(uint(id), handlerID, req.Status, req.Note); err != nil {
		response.ServerError(c, "failed to process report")
		return
	}

	response.Success(c, gin.H{"message": "report processed successfully"})
}

// ============ User Points Management ============

// GetUsers gets user points list
func (h *AdminCommunityHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// TODO: Implement user points list with pagination
	response.Success(c, gin.H{
		"users": []interface{}{},
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     0,
		},
	})
}

// AdjustPoints adjusts user points
func (h *AdminCommunityHandler) AdjustPoints(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Points int    `json:"points" binding:"required"`
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "points and reason are required")
		return
	}

	if err := h.pointService.AdjustPoints(uint(id), req.Points, req.Reason); err != nil {
		response.ServerError(c, "failed to adjust points")
		return
	}

	response.Success(c, gin.H{"message": "points adjusted successfully"})
}

// ============ Statistics ============

// GetStats gets community statistics
func (h *AdminCommunityHandler) GetStats(c *gin.Context) {
	// TODO: Implement statistics
	stats := gin.H{
		"total_posts":    0,
		"total_replies":  0,
		"total_users":    0,
		"total_reports":  0,
		"pending_reports": 0,
	}

	response.Success(c, gin.H{"stats": stats})
}
