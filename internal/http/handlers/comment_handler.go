package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	repo *repository.CommentRepository
}

func NewCommentHandler(repo *repository.CommentRepository) *CommentHandler {
	return &CommentHandler{repo: repo}
}

// sanitizeContent 安全过滤评论内容
// - 移除 HTML 标签（防 XSS）
// - 保留表情标记 [表情名] 格式
func sanitizeContent(content string) string {
	// 移除所有 HTML 标签（script, iframe, 等）
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	cleaned := htmlTagRegex.ReplaceAllString(content, "")

	// 移除危险的 JS 事件属性（如 onclick=, onerror= 等）
	onEventRegex := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	cleaned = onEventRegex.ReplaceAllString(cleaned, "")

	// 移除 javascript: 协议
	jsProtocolRegex := regexp.MustCompile(`(?i)javascript\s*:`)
	cleaned = jsProtocolRegex.ReplaceAllString(cleaned, "")

	// 移除 data: 协议（防止 base64 图片 XSS）
	dataProtocolRegex := regexp.MustCompile(`(?i)data\s*:`)
	cleaned = dataProtocolRegex.ReplaceAllString(cleaned, "")

	// 防止 SQL 注入尝试（移除常见关键字，但不影响正常文本）
	sqlKeywords := []string{"DROP TABLE", "DELETE FROM", "INSERT INTO", "UPDATE SET", "UNION SELECT"}
	for _, keyword := range sqlKeywords {
		if strings.Contains(strings.ToUpper(cleaned), keyword) {
			// 如果包含 SQL 注入尝试，清空内容
			cleaned = "[内容被过滤]"
			break
		}
	}

	// 截断超长内容（防止资源消耗）
	if len(cleaned) > 1000 {
		cleaned = cleaned[:1000] + "..."
	}

	return cleaned
}

// List 获取套图评论列表
func (h *CommentHandler) List(c *gin.Context) {
	photosetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	var req struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 获取当前用户ID（可选）
	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	comments, total, err := h.repo.GetByPhotosetID(uint(photosetID), userID, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取评论失败")
		return
	}

	// 构建评论列表，包含回复数量和当前用户点赞状态
	var list []gin.H
	for _, comment := range comments {
		// 获取该评论的前3条回复
		replies, _, _ := h.repo.GetReplies(comment.ID, userID, 1, 3)

		// 获取当前用户对该评论的点赞状态
		isLiked := false
		if userID > 0 {
			isLiked, _ = h.repo.IsLiked(comment.ID, userID)
		}

		// 获取该评论的总回复数
		_, replyTotal, _ := h.repo.GetReplies(comment.ID, userID, 1, 1000)

		// 构建回复列表
		var replyList []gin.H
		for _, reply := range replies {
			replyIsLiked := false
			if userID > 0 {
				replyIsLiked, _ = h.repo.IsLiked(reply.ID, userID)
			}
			replyList = append(replyList, gin.H{
				"id":         reply.ID,
				"content":    reply.Content,
				"image_url":  reply.ImageURL,
				"like_count": reply.LikeCount,
				"is_liked":   replyIsLiked,
				"created_at": reply.CreatedAt.Format("2006-01-02T15:04:05Z"),
				"user": gin.H{
					"id":       reply.User.ID,
					"nickname": reply.User.Nickname,
				},
				"parent_id": comment.ID,
			})
		}

		list = append(list, gin.H{
			"id":          comment.ID,
			"content":     comment.Content,
			"image_url":   comment.ImageURL,
			"like_count":  comment.LikeCount,
			"is_liked":    isLiked,
			"reply_count": int(replyTotal),
			"created_at":  comment.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"user": gin.H{
				"id":       comment.User.ID,
				"nickname": comment.User.Nickname,
			},
			"replies": replyList,
		})
	}

	// 获取评论总数
	commentCount, _ := h.repo.GetCommentCount(uint(photosetID))

	response.Success(c, gin.H{
		"list":         list,
		"total":        total,
		"page":         req.Page,
		"page_size":    req.PageSize,
		"comment_count": commentCount,
	})
}

// GetReplies 获取评论的回复列表
func (h *CommentHandler) GetReplies(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	var req struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	}

	replies, total, err := h.repo.GetReplies(uint(commentID), userID, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取回复失败")
		return
	}

	var list []gin.H
	for _, reply := range replies {
		isLiked := false
		if userID > 0 {
			isLiked, _ = h.repo.IsLiked(reply.ID, userID)
		}
		list = append(list, gin.H{
			"id":         reply.ID,
			"content":    reply.Content,
			"image_url":  reply.ImageURL,
			"like_count": reply.LikeCount,
			"is_liked":   isLiked,
			"created_at": reply.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"user": gin.H{
				"id":       reply.User.ID,
				"nickname": reply.User.Nickname,
			},
			"parent_id": commentID,
		})
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// Create 发表评论
func (h *CommentHandler) Create(c *gin.Context) {
	photosetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	var req struct {
		Content  string `json:"content" binding:"required,max=500"`
		ImageURL string `json:"image_url"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "评论内容不能为空且不超过500字")
		return
	}

	// 安全过滤评论内容（移除危险标签，保留表情格式）
	sanitizedContent := sanitizeContent(req.Content)

	comment := &domain.Comment{
		PhotoSetID: uint(photosetID),
		UserID:     userID.(uint),
		Content:    sanitizedContent,
		ImageURL:   req.ImageURL,
		ParentID:   req.ParentID,
	}

	if err := h.repo.Create(comment); err != nil {
		log.Printf("[CommentHandler] Create ERROR: %v, comment=%+v", err, comment)
		response.Error(c, http.StatusInternalServerError, "发表评论失败: "+err.Error())
		return
	}

	// 重新查询以获取关联的用户信息
	created, err := h.repo.GetByID(comment.ID)
	if err != nil {
		response.Success(c, gin.H{"id": comment.ID})
		return
	}

	response.Success(c, gin.H{
		"id":         created.ID,
		"content":    created.Content,
		"image_url":  created.ImageURL,
		"like_count": created.LikeCount,
		"is_liked":   false,
		"created_at": created.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"user": gin.H{
			"id":       created.User.ID,
			"nickname": created.User.Nickname,
		},
		"parent_id": req.ParentID,
	})
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	// 检查评论是否存在及权限
	comment, err := h.repo.GetByID(uint(commentID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "评论不存在")
		return
	}

	// 只有作者本人或管理员可以删除
	if comment.UserID != userID.(uint) && userRole.(string) != "admin" {
		response.Error(c, http.StatusForbidden, "无权删除此评论")
		return
	}

	if err := h.repo.Delete(uint(commentID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除评论失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ToggleLike 点赞/取消点赞评论
func (h *CommentHandler) ToggleLike(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	isLiked, err := h.repo.ToggleLike(uint(commentID), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "操作失败")
		return
	}

	// 获取更新后的评论
	comment, _ := h.repo.GetByID(uint(commentID))
	likeCount := 0
	if comment != nil {
		likeCount = comment.LikeCount
	}

	response.Success(c, gin.H{
		"is_liked":   isLiked,
		"like_count": likeCount,
	})
}
