package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
)

type ExploreHandler struct {
	exploreService *service.ExploreService
}

func NewExploreHandler(exploreService *service.ExploreService) *ExploreHandler {
	return &ExploreHandler{
		exploreService: exploreService,
	}
}

// @Summary      探索发现
// @Description  获取个性化的探索发现内容流，支持分页
// @Tags         Explore
// @Accept       json
// @Produce      json
// @Param        page      query int  false  "页码"       default(1)
// @Param        page_size query int  false  "每页数量"   default(20)
// @Success      200  {object}  object  "发现内容"
// @Failure      500  {object}  object  "服务器错误"
// @Router       /api/explore/feed [get]
// GetExploreFeed 获取探索/发现页 Feed
func (h *ExploreHandler) GetExploreFeed(c *gin.Context) {
	// 获取分页参数
	page := 1
	pageSize := 20
	
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	// 限制每页数量
	if pageSize > 50 {
		pageSize = 50
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}

	// 获取用户ID（可选，用于个性化推荐）
	userID, _ := middleware.GetUserID(c)

	// 获取探索 Feed
	feed, err := h.exploreService.GetExploreFeed(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取探索内容失败")
		return
	}

	response.Success(c, feed)
}

// @Summary      热门套图
// @Description  获取热门套图排行榜列表
// @Tags         Explore
// @Accept       json
// @Produce      json
// @Param        page      query int  false  "页码"       default(1)
// @Param        page_size query int  false  "每页数量"   default(20)
// @Success      200  {object}  object  "热门套图列表"
// @Router       /api/explore/hot [get]
// GetHotPhotosets 获取热门套图
func (h *ExploreHandler) GetHotPhotosets(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	if page < 1 {
		page = 1
	}
	if pageSize > 50 {
		pageSize = 50
	}

	feed, err := h.exploreService.GetHotPhotosets(page, pageSize)
	if err != nil {
		response.ServerError(c, "获取热门套图失败")
		return
	}

	response.Success(c, feed)
}