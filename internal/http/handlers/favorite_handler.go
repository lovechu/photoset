package handlers

import (
	"net/http"
	"strconv"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteHandler(repo *repository.FavoriteRepository) *FavoriteHandler {
	return &FavoriteHandler{repo: repo}
}

// Add 收藏套图
func (h *FavoriteHandler) Add(c *gin.Context) {
	photosetID, err := strconv.ParseUint(c.Param("photosetId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.repo.Add(userID.(uint), uint(photosetID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "收藏失败")
		return
	}
	response.Success(c, nil)
}

// Remove 取消收藏
func (h *FavoriteHandler) Remove(c *gin.Context) {
	photosetID, err := strconv.ParseUint(c.Param("photosetId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套图ID")
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.repo.Remove(userID.(uint), uint(photosetID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "取消收藏失败")
		return
	}
	response.Success(c, nil)
}

// List 我的收藏列表
func (h *FavoriteHandler) List(c *gin.Context) {
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
		req.PageSize = 12
	}

	userID, _ := c.Get("user_id")

	favorites, total, err := h.repo.List(userID.(uint), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取收藏列表失败")
		return
	}

	// 返回完整的 FavoriteModel 结构，前端需要 id, user_id, photoset_id, created_at, photoset
	var list []interface{}
	for _, fav := range favorites {
		list = append(list, gin.H{
			"id":          fav.ID,
			"user_id":     fav.UserID,
			"photoset_id": fav.PhotoSetID,
			"created_at":  fav.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"photoset":    fav.PhotoSet,
		})
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}
