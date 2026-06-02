package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/logger"
	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service      *service.OrderService
	alipayService *service.AlipayService
}

func NewOrderHandler(service *service.OrderService, alipayService *service.AlipayService) *OrderHandler {
	return &OrderHandler{
		service:      service,
		alipayService: alipayService,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Type         string `json:"type" binding:"required,oneof=membership single"`
	MembershipID *uint  `json:"membership_id"`
	PhotoSetID   *uint  `json:"photoset_id"`
}

// PayRequest 支付请求
type PayRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay mock"`
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
		logger.Warn("Order creation failed", "user_id", userID, "type", req.Type, "error", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Info("Order created", "order_id", order.ID, "user_id", userID, "type", req.Type)
	response.Success(c, order)
}

// Pay 支付订单
func (h *OrderHandler) Pay(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订单ID")
		return
	}

	var req PayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	switch req.PaymentMethod {
	case "alipay":
		// 支付宝支付
		if h.alipayService == nil {
			response.Error(c, http.StatusServiceUnavailable, "支付宝支付未配置")
			return
		}

		payURL, err := h.service.CreateAlipayOrder(userID.(uint), uint(orderID), h.alipayService)
		if err != nil {
			logger.Warn("支付宝支付创建失败", "order_id", orderID, "user_id", userID, "error", err)
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		logger.Info("支付宝支付订单创建成功", "order_id", orderID, "user_id", userID)
		response.Success(c, gin.H{
			"payment_method": "alipay",
			"pay_url":        payURL,
		})

	case "mock":
		// 模拟支付（开发测试用）
		token, err := h.service.MockPay(userID.(uint), uint(orderID))
		if err != nil {
			logger.Warn("模拟支付失败", "order_id", orderID, "user_id", userID, "error", err)
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		logger.Info("模拟支付成功", "order_id", orderID, "user_id", userID)
		response.Success(c, gin.H{
			"message":        "支付成功",
			"payment_method": "mock",
			"token":          token,
		})

	default:
		response.Error(c, http.StatusBadRequest, "不支持的支付方式")
	}
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

// AlipayNotify 支付宝异步回调
func (h *OrderHandler) AlipayNotify(c *gin.Context) {
	if h.alipayService == nil {
		logger.Error("支付宝回调: 服务未配置")
		c.String(http.StatusOK, "fail")
		return
	}

	notification, err := h.alipayService.VerifyNotification(c.Request)
	if err != nil {
		logger.Error("支付宝回调验证失败", "error", err)
		c.String(http.StatusOK, "fail")
		return
	}

	if err := h.alipayService.HandlePaymentSuccess(notification); err != nil {
		logger.Error("支付宝回调处理失败", "error", err)
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}
