package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"photoset/internal/http/middleware"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// ReviewHandler 评价处理器
type ReviewHandler struct {
	reviewService *service.ReviewService
}

// NewReviewHandler 创建评价处理器
func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReviewRequest 创建评价请求
type CreateReviewRequest struct {
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Tags        string `json:"tags"`
	Content     string `json:"content"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// @Summary      创建评价
// @Description  对套图进行评分（1-5星）和评价
// @Tags         Review
// @Accept       json
// @Produce      json
// @Param        id   path  int   true  "套图ID"
// @Param        body body  object  true  "评价请求 {rating, tags, content, is_anonymous}"
// @Success      200  {object}  object  "评价结果"
// @Failure      400  {object}  object  "参数错误"
// @Security     BearerAuth
// @Router       /api/photosets/{id}/reviews [post]
// Create 创建评价
func (h *ReviewHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请先登录"})
		return
	}

	photosetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套图ID"})
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	review, err := h.reviewService.CreateReview(userID, uint(photosetID), req.Rating, req.Tags, req.Content, req.IsAnonymous)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": review})
}

// @Summary      更新评价
// @Description  更新自己对套图的评价
// @Tags         Review
// @Accept       json
// @Produce      json
// @Param        id       path  int   true  "套图ID"
// @Param        reviewId path  int   true  "评价ID"
// @Param        body     body  object  true  "更新请求 {rating, tags, content, is_anonymous}"
// @Success      200  {object}  object  "更新结果"
// @Failure      400  {object}  object  "参数错误"
// @Security     BearerAuth
// @Router       /api/photosets/{id}/reviews/{reviewId} [put]
// Update 更新评价
func (h *ReviewHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请先登录"})
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("reviewId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评价ID"})
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	review, err := h.reviewService.UpdateReview(uint(reviewID), userID, req.Rating, req.Tags, req.Content, req.IsAnonymous)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": review})
}

// @Summary      删除评价
// @Description  删除自己对套图的评价
// @Tags         Review
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "套图ID"
// @Param        reviewId path  int  true  "评价ID"
// @Success      200  {object}  object  "删除成功"
// @Failure      400  {object}  object  "参数错误"
// @Security     BearerAuth
// @Router       /api/photosets/{id}/reviews/{reviewId} [delete]
// Delete 删除评价
func (h *ReviewHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请先登录"})
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("reviewId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评价ID"})
		return
	}

	if err := h.reviewService.DeleteReview(uint(reviewID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// @Summary      获取评价列表
// @Description  分页获取套图评价列表，支持排序
// @Tags         Review
// @Accept       json
// @Produce      json
// @Param        id        path  int     true   "套图ID"
// @Param        page      query int     false  "页码"                      default(1)
// @Param        page_size query int     false  "每页数量"                  default(20)
// @Param        sort_by   query string  false  "排序方式(newest/highest/lowest)" default(newest)
// @Success      200  {object}  object  "评价列表"
// @Failure      400  {object}  object  "参数错误"
// @Router       /api/photosets/{id}/reviews [get]
// List 获取评价列表
func (h *ReviewHandler) List(c *gin.Context) {
	photosetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套图ID"})
		return
	}

	page := 1
	pageSize := 20
	sortBy := "newest"

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	if s := c.Query("sort_by"); s != "" {
		sortBy = s
	}

	reviews, total, err := h.reviewService.ListReviews(uint(photosetID), page, pageSize, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      reviews,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// @Summary      获取评价汇总
// @Description  获取套图的评价统计汇总（评分分布、标签统计等）
// @Tags         Review
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "套图ID"
// @Success      200  {object}  object  "评价汇总"
// @Failure      400  {object}  object  "参数错误"
// @Router       /api/photosets/{id}/reviews/summary [get]
// GetSummary 获取评价汇总
func (h *ReviewHandler) GetSummary(c *gin.Context) {
	photosetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套图ID"})
		return
	}

	summary, err := h.reviewService.GetReviewSummary(uint(photosetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评价汇总失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": summary})
}

// @Summary      获取我的评价
// @Description  获取当前用户对指定套图的评价
// @Tags         Review
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "套图ID"
// @Success      200  {object}  object  "我的评价"
// @Failure      401  {object}  object  "未登录"
// @Security     BearerAuth
// @Router       /api/photosets/{id}/reviews/mine [get]
// GetMyReview 获取我的评价
func (h *ReviewHandler) GetMyReview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请先登录"})
		return
	}

	photosetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套图ID"})
		return
	}

	review, err := h.reviewService.GetUserReview(userID, uint(photosetID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": review})
}
