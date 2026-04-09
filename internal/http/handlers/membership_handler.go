package handlers

import (
	"photoset/internal/pkg/response"
	"photoset/internal/repository"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	repo *repository.MembershipRepository
}

func NewMembershipHandler(repo *repository.MembershipRepository) *MembershipHandler {
	return &MembershipHandler{repo: repo}
}

// List 获取上架的会员套餐列表（公开接口）
func (h *MembershipHandler) List(c *gin.Context) {
	plans, err := h.repo.ListActive()
	if err != nil {
		response.Error(c, 500, "获取套餐列表失败")
		return
	}
	response.Success(c, plans)
}
