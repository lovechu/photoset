package handlers

import (
	"net/http"
	"photoset/internal/domain"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"strconv"

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

// ============ Admin APIs ============

// AdminList 管理员分页查询所有会员套餐
func (h *MembershipHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	plans, total, err := h.repo.ListAll(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取套餐列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  plans,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// AdminCreate 管理员创建会员套餐
func (h *MembershipHandler) AdminCreate(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Duration    int     `json:"duration" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Description string  `json:"description"`
		Status      *int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	plan := &domain.MembershipPlan{
		Name:        req.Name,
		Duration:    req.Duration,
		Price:       req.Price,
		Description: req.Description,
		Status:      1, // 默认上架
	}
	if req.Status != nil {
		plan.Status = *req.Status
	}

	if err := h.repo.Create(plan); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建套餐失败")
		return
	}

	response.Success(c, plan)
}

// AdminUpdate 管理员更新会员套餐
func (h *MembershipHandler) AdminUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套餐ID")
		return
	}

	plan, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "套餐不存在")
		return
	}

	var req struct {
		Name        *string  `json:"name"`
		Duration    *int     `json:"duration"`
		Price       *float64 `json:"price"`
		Description *string  `json:"description"`
		Status      *int     `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.Duration != nil {
		plan.Duration = *req.Duration
	}
	if req.Price != nil {
		plan.Price = *req.Price
	}
	if req.Description != nil {
		plan.Description = *req.Description
	}
	if req.Status != nil {
		plan.Status = *req.Status
	}

	if err := h.repo.Update(plan); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新套餐失败")
		return
	}

	response.Success(c, plan)
}

// AdminDelete 管理员删除会员套餐
func (h *MembershipHandler) AdminDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的套餐ID")
		return
	}

	if _, err := h.repo.FindByID(uint(id)); err != nil {
		response.Error(c, http.StatusNotFound, "套餐不存在")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除套餐失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
