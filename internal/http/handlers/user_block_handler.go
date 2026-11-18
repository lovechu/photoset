package handlers

import (
	"photoset/internal/pkg/response"
	"photoset/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserBlockHandler handles user block HTTP requests
type UserBlockHandler struct {
	blockService *service.UserBlockService
}

// NewUserBlockHandler creates a new UserBlockHandler
func NewUserBlockHandler(blockService *service.UserBlockService) *UserBlockHandler {
	return &UserBlockHandler{blockService: blockService}
}

// BlockUser handles blocking a user
func (h *UserBlockHandler) BlockUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	blockedID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.blockService.BlockUser(userID.(uint), uint(blockedID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "拉黑成功", nil)
}

// MuteUser handles muting a user
func (h *UserBlockHandler) MuteUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	mutedID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.blockService.MuteUser(userID.(uint), uint(mutedID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "静音成功", nil)
}

// UnblockUser handles unblocking/unmuting a user
func (h *UserBlockHandler) UnblockUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	blockedID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.blockService.UnblockUser(userID.(uint), uint(blockedID)); err != nil {
		response.ServerError(c, "操作失败")
		return
	}

	response.SuccessWithMessage(c, "取消成功", nil)
}

// GetBlockedUsers handles getting list of blocked users
func (h *UserBlockHandler) GetBlockedUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	blocks, total, err := h.blockService.GetBlockedUsers(userID.(uint), page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	// Format response
	type BlockItem struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		BlockedID uint   `json:"blocked_id"`
		BlockType string `json:"block_type"`
		CreatedAt string `json:"created_at"`
		Username  string `json:"username"`
		Avatar    string `json:"avatar"`
	}

	items := make([]BlockItem, 0, len(blocks))
	for _, block := range blocks {
		item := BlockItem{
			ID:        block.ID,
			UserID:    block.UserID,
			BlockedID: block.BlockedID,
			BlockType: block.BlockType,
			CreatedAt: block.CreatedAt.Format("2006-01-02 15:04:05"),
			Username:  block.Blocked.Nickname,
			Avatar:    block.Blocked.Avatar,
		}
		items = append(items, item)
	}

	response.Success(c, gin.H{
		"list":  items,
		"total": total,
	})
}

// GetBlockStatus handles getting block status with a user
func (h *UserBlockHandler) GetBlockStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	status, err := h.blockService.GetBlockStatus(userID.(uint), uint(targetID))
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"status": status,
	})
}
