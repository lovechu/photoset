package admin

import (
	"fmt"
	"regexp"
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
	categoryRepo    *repository.PostCategoryRepository
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
		categoryRepo:  repository.NewPostCategoryRepository(db),
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

	var postID uint
	if postIDStr != "" {
		id, err := strconv.ParseUint(postIDStr, 10, 64)
		if err == nil {
			postID = uint(id)
		}
	}

	replies, total, err := h.replyRepo.ListForAdmin(page, pageSize, postID)
	if err != nil {
		response.ServerError(c, "failed to get replies")
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
		Word        string `json:"word" binding:"required"`
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
		Word:        req.Word,
		Replacement: req.Replacement,
		IsActive:    true,
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
		Word        string `json:"word"`
		Replacement string `json:"replacement"`
		IsActive    *bool  `json:"is_active"`
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
	levelStr := c.DefaultQuery("level", "0")
	level, _ := strconv.Atoi(levelStr)

	users, total, err := h.pointRepo.ListForAdmin(page, pageSize, level)
	if err != nil {
		response.ServerError(c, "failed to get users")
		return
	}

	response.Success(c, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
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
	// Count total posts
	var totalPosts int64
	if err := h.postRepo.DB.Model(&domain.Post{}).Count(&totalPosts).Error; err != nil {
		response.ServerError(c, "failed to count posts")
		return
	}

	// Count total replies
	var totalReplies int64
	if err := h.replyRepo.DB.Model(&domain.PostReply{}).Count(&totalReplies).Error; err != nil {
		response.ServerError(c, "failed to count replies")
		return
	}

	// Count total users with points (unique user_ids in user_points table)
	var totalUsers int64
	if err := h.pointRepo.DB.Model(&domain.UserPoint{}).Count(&totalUsers).Error; err != nil {
		response.ServerError(c, "failed to count users")
		return
	}

	// Count total reports
	var totalReports int64
	if err := h.reportRepo.DB.Model(&domain.PostReport{}).Count(&totalReports).Error; err != nil {
		response.ServerError(c, "failed to count reports")
		return
	}

	// Count pending reports (status = 'pending')
	var pendingReports int64
	if err := h.reportRepo.DB.Model(&domain.PostReport{}).Where("status = ?", "pending").Count(&pendingReports).Error; err != nil {
		response.ServerError(c, "failed to count pending reports")
		return
	}

	stats := gin.H{
		"total_posts":     totalPosts,
		"total_replies":   totalReplies,
		"total_users":     totalUsers,
		"total_reports":   totalReports,
		"pending_reports": pendingReports,
	}

	response.Success(c, gin.H{"stats": stats})
}

// ============ Category Management ============

// categoryKeyPattern validates category key: lowercase letters, digits, and underscores only
var categoryKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// ListCategories returns all post categories with post counts
func (h *AdminCommunityHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryRepo.ListCategories()
	if err != nil {
		response.ServerError(c, "failed to get categories")
		return
	}

	type CategoryWithCount struct {
		domain.CommunityCategory
		PostCount int64 `json:"post_count"`
	}

	results := make([]CategoryWithCount, 0, len(categories))
	for _, cat := range categories {
		count, _ := h.categoryRepo.CountPostsByCategory(cat.Key)
		results = append(results, CategoryWithCount{
			CommunityCategory: cat,
			PostCount:    count,
		})
	}

	response.Success(c, gin.H{"categories": results})
}

// CreateCategory creates a new post category
func (h *AdminCommunityHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Key         string `json:"key" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "key and name are required")
		return
	}

	// Validate key format
	if !categoryKeyPattern.MatchString(req.Key) {
		response.Error(c, 400, "category key must start with a lowercase letter and contain only lowercase letters, digits, and underscores")
		return
	}

	// Check key uniqueness
	if _, err := h.categoryRepo.GetCategoryByKey(req.Key); err == nil {
		response.Error(c, 400, "category key already exists")
		return
	}

	cat := &domain.CommunityCategory{
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
	}

	if cat.Color == "" {
		cat.Color = "#409EFF"
	}

	if err := h.categoryRepo.CreateCategory(cat); err != nil {
		response.ServerError(c, "failed to create category")
		return
	}

	response.Success(c, gin.H{"category": cat})
}

// UpdateCategory updates an existing post category
func (h *AdminCommunityHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
		SortOrder   *int   `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	// Ensure category exists
	if _, err := h.categoryRepo.GetCategoryByID(uint(id)); err != nil {
		response.NotFound(c, "category not found")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	// Do not allow modifying the key
	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}

	if err := h.categoryRepo.UpdateCategory(uint(id), updates); err != nil {
		response.ServerError(c, "failed to update category")
		return
	}

	response.Success(c, gin.H{"message": "category updated successfully"})
}

// DeleteCategory deletes a post category (only if no posts use it)
func (h *AdminCommunityHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	cat, err := h.categoryRepo.GetCategoryByID(uint(id))
	if err != nil {
		response.NotFound(c, "category not found")
		return
	}

	// Check if any posts use this category
	count, err := h.categoryRepo.CountPostsByCategory(cat.Key)
	if err != nil {
		response.ServerError(c, "failed to check category usage")
		return
	}
	if count > 0 {
		response.Error(c, 400, fmt.Sprintf("该分类下有 %d 篇帖子，无法删除", count))
		return
	}

	if err := h.categoryRepo.DeleteCategory(uint(id)); err != nil {
		response.ServerError(c, "failed to delete category")
		return
	}

	response.Success(c, gin.H{"message": "category deleted successfully"})
}

// SortCategories batch-updates sort_order for categories
func (h *AdminCommunityHandler) SortCategories(c *gin.Context) {
	var req []struct {
		ID        uint `json:"id" binding:"required"`
		SortOrder int  `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: expected array of {id, sort_order}")
		return
	}

	for _, item := range req {
		if err := h.categoryRepo.UpdateCategory(item.ID, map[string]interface{}{
			"sort_order": item.SortOrder,
		}); err != nil {
			response.ServerError(c, fmt.Sprintf("failed to update category %d: %v", item.ID, err))
			return
		}
	}

	response.Success(c, gin.H{"message": "categories sorted successfully"})
}
