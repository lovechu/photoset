package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"photoset/internal/config"
	"photoset/internal/domain"
	"photoset/internal/logger"
	"photoset/internal/repository"

	"github.com/smartwalle/alipay/v3"
)

type AlipayService struct {
	mu        sync.RWMutex
	client    *alipay.Client
	orderRepo *repository.OrderRepository
	config    *config.AlipayConfig
}

func NewAlipayService(cfg *config.AlipayConfig, orderRepo *repository.OrderRepository) (*AlipayService, error) {
	if cfg.AppID == "" || cfg.PrivateKey == "" || cfg.PublicKey == "" {
		return nil, errors.New("支付宝配置不完整，请检查 ALIPAY_APP_ID、ALIPAY_PRIVATE_KEY、ALIPAY_PUBLIC_KEY")
	}

	client, err := alipay.New(cfg.AppID, cfg.PrivateKey, cfg.IsSandbox)
	if err != nil {
		return nil, fmt.Errorf("初始化支付宝客户端失败: %w", err)
	}

	// 加载支付宝公钥
	if err := client.LoadAliPayPublicKey(cfg.PublicKey); err != nil {
		return nil, fmt.Errorf("加载支付宝公钥失败: %w", err)
	}

	return &AlipayService{
		client:    client,
		orderRepo: orderRepo,
		config:    cfg,
	}, nil
}

// Reload 使用新配置重新初始化支付宝客户端
func (s *AlipayService) Reload(cfg *config.AlipayConfig) error {
	if cfg.AppID == "" || cfg.PrivateKey == "" || cfg.PublicKey == "" {
		return errors.New("支付宝配置不完整")
	}

	client, err := alipay.New(cfg.AppID, cfg.PrivateKey, cfg.IsSandbox)
	if err != nil {
		return fmt.Errorf("重新初始化支付宝客户端失败: %w", err)
	}

	if err := client.LoadAliPayPublicKey(cfg.PublicKey); err != nil {
		return fmt.Errorf("重新加载支付宝公钥失败: %w", err)
	}

	s.mu.Lock()
	s.client = client
	s.config = cfg
	s.mu.Unlock()

	logger.Info("支付宝服务配置已更新")
	return nil
}

// ReloadFromSettings 从数据库配置重载支付宝服务
func (s *AlipayService) ReloadFromSettings(settings map[string]string) {
	cfg := &config.AlipayConfig{
		AppID:      settings["alipay_app_id"],
		PrivateKey: settings["alipay_private_key"],
		PublicKey:  settings["alipay_public_key"],
		NotifyURL:  settings["alipay_notify_url"],
		ReturnURL:  settings["alipay_return_url"],
		IsSandbox:  settings["alipay_sandbox"] == "true",
	}

	if cfg.AppID == "" || cfg.PrivateKey == "" || cfg.PublicKey == "" {
		return // 配置不完整，不重载
	}

	if err := s.Reload(cfg); err != nil {
		logger.Warn("从数据库重载支付宝配置失败", "error", err)
	}
}

// GetConfig 获取当前配置（脱敏）
func (s *AlipayService) GetConfig() map[string]interface{} {
	s.mu.RLock()
	cfg := s.config
	s.mu.RUnlock()

	if cfg == nil {
		return nil
	}

	return map[string]interface{}{
		"alipay_app_id":        cfg.AppID,
		"alipay_private_key_set": cfg.PrivateKey != "",
		"alipay_public_key_set":  cfg.PublicKey != "",
		"alipay_notify_url":    cfg.NotifyURL,
		"alipay_return_url":    cfg.ReturnURL,
		"alipay_sandbox":       cfg.IsSandbox,
	}
}

// CreatePayment 创建支付宝支付订单
func (s *AlipayService) CreatePayment(order *domain.Order) (string, error) {
	s.mu.RLock()
	client := s.client
	cfg := s.config
	s.mu.RUnlock()

	if client == nil {
		return "", errors.New("支付宝服务未初始化")
	}

	// 构建支付参数
	var p = alipay.TradePagePay{}
	p.NotifyURL = cfg.NotifyURL
	p.ReturnURL = cfg.ReturnURL
	p.Subject = fmt.Sprintf("PhotoSet订单-%s", order.OrderNo)
	p.OutTradeNo = order.OrderNo
	p.TotalAmount = fmt.Sprintf("%.2f", order.Amount)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	// 发起支付请求
	url, err := client.TradePagePay(p)
	if err != nil {
		logger.Error("创建支付宝支付订单失败", "error", err, "order_no", order.OrderNo)
		return "", fmt.Errorf("创建支付订单失败: %w", err)
	}

	return url.String(), nil
}

// VerifyNotification 验证支付宝异步通知
func (s *AlipayService) VerifyNotification(req *http.Request) (*alipay.Notification, error) {
	s.mu.RLock()
	client := s.client
	s.mu.RUnlock()

	if client == nil {
		return nil, errors.New("支付宝服务未初始化")
	}

	notification, err := client.GetTradeNotification(req)
	if err != nil {
		logger.Error("解析支付宝通知失败", "error", err)
		return nil, fmt.Errorf("解析通知失败: %w", err)
	}

	return notification, nil
}

// HandlePaymentSuccess 处理支付成功
func (s *AlipayService) HandlePaymentSuccess(notification *alipay.Notification) error {
	// 查找订单
	order, err := s.orderRepo.FindByOrderNo(notification.OutTradeNo)
	if err != nil {
		logger.Error("支付宝回调: 订单不存在", "order_no", notification.OutTradeNo)
		return errors.New("订单不存在")
	}

	// 检查订单状态
	if order.Status != "pending" {
		logger.Info("支付宝回调: 订单已处理", "order_no", notification.OutTradeNo, "status", order.Status)
		return nil // 已处理过的订单，直接返回成功
	}

	// 验证金额
	expectedAmount := fmt.Sprintf("%.2f", order.Amount)
	if notification.TotalAmount != expectedAmount {
		logger.Error("支付宝回调: 金额不匹配", "order_no", notification.OutTradeNo, 
			"expected", expectedAmount, "actual", notification.TotalAmount)
		return errors.New("金额不匹配")
	}

	// 更新订单状态
	order.Status = "paid"
	order.PaymentNo = notification.TradeNo
	order.PaymentMethod = "alipay"

	if err := s.orderRepo.Update(order); err != nil {
		logger.Error("支付宝回调: 更新订单失败", "error", err, "order_no", notification.OutTradeNo)
		return errors.New("更新订单失败")
	}

	logger.Info("支付宝回调: 支付成功", "order_no", notification.OutTradeNo, "trade_no", notification.TradeNo)
	return nil
}

// QueryOrder 查询支付宝订单状态
func (s *AlipayService) QueryOrder(orderNo string) (string, error) {
	s.mu.RLock()
	client := s.client
	s.mu.RUnlock()

	if client == nil {
		return "", errors.New("支付宝服务未初始化")
	}

	var p = alipay.TradeQuery{}
	p.OutTradeNo = orderNo

	result, err := client.TradeQuery(context.Background(), p)
	if err != nil {
		logger.Error("查询支付宝订单失败", "error", err, "order_no", orderNo)
		return "", fmt.Errorf("查询订单失败: %w", err)
	}

	if result.Code != "10000" {
		return "", fmt.Errorf("查询失败: %s - %s", result.Code, result.Msg)
	}

	return string(result.TradeStatus), nil
}

// CloseOrder 关闭支付宝订单
func (s *AlipayService) CloseOrder(orderNo string) error {
	s.mu.RLock()
	client := s.client
	s.mu.RUnlock()

	if client == nil {
		return errors.New("支付宝服务未初始化")
	}

	var p = alipay.TradeClose{}
	p.OutTradeNo = orderNo

	result, err := client.TradeClose(context.Background(), p)
	if err != nil {
		logger.Error("关闭支付宝订单失败", "error", err, "order_no", orderNo)
		return fmt.Errorf("关闭订单失败: %w", err)
	}

	if result.Code != "10000" {
		return fmt.Errorf("关闭失败: %s - %s", result.Code, result.Msg)
	}

	return nil
}

// RefundOrder 支付宝退款
func (s *AlipayService) RefundOrder(orderNo string, refundAmount float64, reason string) error {
	s.mu.RLock()
	client := s.client
	s.mu.RUnlock()

	if client == nil {
		return errors.New("支付宝服务未初始化")
	}

	order, err := s.orderRepo.FindByOrderNo(orderNo)
	if err != nil {
		return errors.New("订单不存在")
	}

	var p = alipay.TradeRefund{}
	p.OutTradeNo = orderNo
	p.RefundAmount = fmt.Sprintf("%.2f", refundAmount)
	p.RefundReason = reason

	result, err := client.TradeRefund(context.Background(), p)
	if err != nil {
		logger.Error("支付宝退款失败", "error", err, "order_no", orderNo)
		return fmt.Errorf("退款请求失败: %w", err)
	}

	if result.Code != "10000" {
		return fmt.Errorf("退款失败: %s - %s", result.Code, result.Msg)
	}

	// 更新订单状态
	order.Status = "refunded"
	order.PaymentNo = result.TradeNo
	if err := s.orderRepo.Update(order); err != nil {
		logger.Error("更新退款订单失败", "error", err, "order_no", orderNo)
		return errors.New("更新订单状态失败")
	}

	logger.Info("支付宝退款成功", "order_no", orderNo, "refund_fee", result.RefundFee)
	return nil
}