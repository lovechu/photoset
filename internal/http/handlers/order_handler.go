package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Type         string `json:"type" binding:"required,oneof=membership single"`
	MembershipID *uint  `json:"membership_id"`
	PhotoSetID   *uint  `json:"photoset_id"`
}

// Create 创建订单
func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	order, err := h.service.CreateOrder(userID.(uint), req.Type, req.MembershipID, req.PhotoSetID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, order)
}

// Pay 模拟支付
func (h *OrderHandler) Pay(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}

	userID, _ := c.Get("user_id")

	token, err := h.service.MockPay(userID.(uint), uint(orderID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "支付成功",
		"token":   token,
	})
}

// List 我的订单列表
func (h *OrderHandler) List(c *gin.Context) {
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

	userID, _ := c.Get("user_id")

	orders, total, err := h.service.GetOrderList(userID.(uint), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// Refund 用户自助退款
func (h *OrderHandler) Refund(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.service.RefundOrder(userID.(uint), uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "退款成功"})
}
