package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// FollowHandler handles follow-related requests
type FollowHandler struct {
	followService service.FollowService
}

// NewFollowHandler creates a new FollowHandler
func NewFollowHandler(followService service.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

// Follow 关注用户
func (h *FollowHandler) Follow(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := h.followService.Follow(userID, uint(followingID)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "关注成功"})
}

// Unfollow 取消关注
func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := h.followService.Unfollow(userID, uint(followingID)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "取消关注成功"})
}

// CheckFollowing 检查是否关注了某用户
func (h *FollowHandler) CheckFollowing(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	isFollowing, err := h.followService.IsFollowing(userID, uint(followingID))
	if err != nil {
		response.ServerError(c, "检查关注状态失败")
		return
	}

	response.Success(c, gin.H{"is_following": isFollowing})
}

// BatchCheckFollowing 批量检查关注状态
func (h *FollowHandler) BatchCheckFollowing(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	var req struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	followMap, err := h.followService.BatchCheckFollowing(userID, req.UserIDs)
	if err != nil {
		response.ServerError(c, "批量检查关注状态失败")
		return
	}

	response.Success(c, gin.H{"following_map": followMap})
}

// GetFollowingList 获取关注列表
func (h *FollowHandler) GetFollowingList(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.followService.GetFollowingList(uint(userID), page, pageSize)
	if err != nil {
		response.ServerError(c, "获取关注列表失败")
		return
	}

	// 转换为简化的用户信息列表
	userList := make([]gin.H, 0, len(users))
	for _, u := range users {
		userList = append(userList, gin.H{
			"id":        u.ID,
			"nickname":  u.Nickname,
			"avatar":    u.Avatar,
			"bio":       u.Bio,
			"level":     u.Level,
		})
	}

	response.Success(c, gin.H{
		"users": userList,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// GetFollowerList 获取粉丝列表
func (h *FollowHandler) GetFollowerList(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.followService.GetFollowerList(uint(userID), page, pageSize)
	if err != nil {
		response.ServerError(c, "获取粉丝列表失败")
		return
	}

	// 转换为简化的用户信息列表
	userList := make([]gin.H, 0, len(users))
	for _, u := range users {
		userList = append(userList, gin.H{
			"id":        u.ID,
			"nickname":  u.Nickname,
			"avatar":    u.Avatar,
			"bio":       u.Bio,
			"level":     u.Level,
		})
	}

	response.Success(c, gin.H{
		"users": userList,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}
