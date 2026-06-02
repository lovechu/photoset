package service

import (
	"errors"
	"sync"
	"photoset/internal/config"
	"photoset/internal/logger"
	"photoset/internal/repository"
)

type WechatPayService struct {
	mu        sync.RWMutex
	orderRepo *repository.OrderRepository
	config    *config.WechatPayConfig
}

func NewWechatPayService(cfg *config.WechatPayConfig, orderRepo *repository.OrderRepository) *WechatPayService {
	return &WechatPayService{
		orderRepo: orderRepo,
		config:    cfg,
	}
}

// Reload 使用新配置重载微信支付服务
func (s *WechatPayService) Reload(cfg *config.WechatPayConfig) error {
	if cfg.MchID == "" || cfg.APIKey == "" {
		return errors.New("微信支付配置不完整")
	}

	s.mu.Lock()
	s.config = cfg
	s.mu.Unlock()

	logger.Info("微信支付服务配置已更新")
	return nil
}

// ReloadFromSettings 从数据库配置重载微信支付服务
func (s *WechatPayService) ReloadFromSettings(settings map[string]string) {
	cfg := &config.WechatPayConfig{
		AppID:     settings["wechat_app_id"],
		MchID:     settings["wechat_mch_id"],
		APIKey:    settings["wechat_api_key"],
		CertPath:  settings["wechat_cert_path"],
		NotifyURL: settings["wechat_notify_url"],
	}

	if cfg.MchID == "" || cfg.APIKey == "" {
		return // 配置不完整，不重载
	}

	if err := s.Reload(cfg); err != nil {
		logger.Warn("从数据库重载微信支付配置失败", "error", err)
	}
}

// GetConfig 获取当前配置（脱敏）
func (s *WechatPayService) GetConfig() map[string]interface{} {
	s.mu.RLock()
	cfg := s.config
	s.mu.RUnlock()

	if cfg == nil {
		return nil
	}

	return map[string]interface{}{
		"wechat_app_id":     cfg.AppID,
		"wechat_mch_id":     cfg.MchID,
		"wechat_api_key_set": cfg.APIKey != "",
		"wechat_cert_path":  cfg.CertPath,
		"wechat_notify_url": cfg.NotifyURL,
	}
}

// ValidateConfig 验证配置是否可用
func (s *WechatPayService) ValidateConfig() error {
	s.mu.RLock()
	cfg := s.config
	s.mu.RUnlock()

	if cfg == nil {
		return errors.New("微信支付服务未初始化")
	}
	if cfg.MchID == "" {
		return errors.New("商户号未配置")
	}
	if cfg.APIKey == "" {
		return errors.New("API密钥未配置")
	}
	return nil
}
