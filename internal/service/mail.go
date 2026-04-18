package service

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"photoset/internal/repository"
)

// SMTP 连接超时时间
const smtpTimeout = 10 * time.Second

// MailConfig SMTP 配置
type MailConfig struct {
	Host     string
	Port     int
	From     string
	Username string
	Password string
	To       string
	UseTLS   bool
}

// MailService 邮件服务
type MailService struct {
	settingRepo *repository.SiteSettingRepository
	config      *MailConfig
	configMu    sync.RWMutex
}

// NewMailService 创建邮件服务
func NewMailService() *MailService {
	return &MailService{
		settingRepo: repository.NewSiteSettingRepository(),
	}
}

// getConfig 获取 SMTP 配置
func (s *MailService) getConfig() (*MailConfig, error) {
	s.configMu.RLock()
	if s.config != nil {
		defer s.configMu.RUnlock()
		return s.config, nil
	}
	s.configMu.RUnlock()

	settings, err := s.settingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	cfg := &MailConfig{
		Host:     settings["smtp_host"],
		Username: settings["smtp_user"],  // 前端保存的字段是 smtp_user
		Password: settings["smtp_password"],
		From:     settings["smtp_from"],
		UseTLS:   settings["smtp_use_tls"] == "true",
	}

	// 解析端口
	if portStr := settings["smtp_port"]; portStr != "" {
		var port int
		if _, err := fmt.Sscanf(portStr, "%d", &port); err == nil {
			cfg.Port = port
		}
	}

	// 默认端口
	if cfg.Port == 0 {
		if cfg.UseTLS {
			cfg.Port = 465
		} else {
			cfg.Port = 587
		}
	}

	s.configMu.Lock()
	s.config = cfg
	s.configMu.Unlock()

	return cfg, nil
}

// InvalidateConfig 使配置缓存失效（配置更新后调用）
func (s *MailService) InvalidateConfig() {
	s.configMu.Lock()
	s.config = nil
	s.configMu.Unlock()
}

// Send 发送邮件
func (s *MailService) Send(to, subject, body string) error {
	cfg, err := s.getConfig()
	if err != nil {
		return fmt.Errorf("获取SMTP配置失败: %w", err)
	}

	if cfg.Host == "" || cfg.Username == "" || cfg.Password == "" {
		return fmt.Errorf("SMTP配置不完整")
	}

	return s.sendMail(cfg, to, subject, body)
}

// sendMail 使用配置发送邮件
func (s *MailService) sendMail(cfg *MailConfig, to, subject, body string) error {
	// 构建邮件内容
	headers := make(map[string]string)
	headers["From"] = cfg.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// SMTP 认证
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	if cfg.UseTLS && cfg.Port == 465 {
		// 使用 TLS 连接（端口 465）
		return s.sendWithTLS(cfg, to, auth, msg.String())
	}

	// STARTTLS（端口 587）
	return s.sendWithSTARTTLS(cfg, to, auth, msg.String())
}

// sendWithTLS 使用 TLS 发送邮件
func (s *MailService) sendWithTLS(cfg *MailConfig, to string, auth smtp.Auth, msg string) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// 带超时的 TCP 连接
	conn, err := net.DialTimeout("tcp", addr, smtpTimeout)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败（超时%d秒）: %w", int(smtpTimeout.Seconds()), err)
	}

	// 升级到 TLS
	tlsConfig := &tls.Config{
		ServerName: cfg.Host,
	}
	tlsConn := tls.Client(conn, tlsConfig)
	defer tlsConn.Close()

	client, err := smtp.NewClient(tlsConn, cfg.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	// 认证
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	// 设置发送者
	if err := client.Mail(cfg.Username); err != nil {
		return fmt.Errorf("设置发送者失败: %w", err)
	}

	// 设置接收者
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("设置接收者失败: %w", err)
	}

	// 发送内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送数据失败: %w", err)
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("关闭数据流失败: %w", err)
	}

	return client.Quit()
}

// sendWithSTARTTLS 使用 STARTTLS 发送邮件
func (s *MailService) sendWithSTARTTLS(cfg *MailConfig, to string, auth smtp.Auth, msg string) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// 带超时连接服务器
	conn, err := net.DialTimeout("tcp", addr, smtpTimeout)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败（超时%d秒）: %w", int(smtpTimeout.Seconds()), err)
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	// STARTTLS
	tlsConfig := &tls.Config{
		ServerName: cfg.Host,
	}

	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS失败: %w", err)
	}

	// 认证
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	// 设置发送者
	if err := client.Mail(cfg.Username); err != nil {
		return fmt.Errorf("设置发送者失败: %w", err)
	}

	// 设置接收者
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("设置接收者失败: %w", err)
	}

	// 发送内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送数据失败: %w", err)
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("关闭数据流失败: %w", err)
	}

	return client.Quit()
}

// SendTemplate 发送模板邮件
func (s *MailService) SendTemplate(to, subject, template string, data map[string]interface{}) error {
	// 简单的模板替换
	body := template
	for k, v := range data {
		placeholder := fmt.Sprintf("{{.%s}}", k)
		body = strings.ReplaceAll(body, placeholder, fmt.Sprintf("%v", v))
	}

	return s.Send(to, subject, body)
}

// TestConnection 测试 SMTP 连接
func (s *MailService) TestConnection() (bool, string) {
	cfg, err := s.getConfig()
	if err != nil {
		return false, fmt.Sprintf("获取配置失败: %v", err)
	}

	if cfg.Host == "" {
		return false, "SMTP服务器地址未配置"
	}

	if cfg.Username == "" || cfg.Password == "" {
		return false, "SMTP用户名或密码未配置"
	}

	// 尝试发送测试邮件到自己
	testSubject := "PhotoSet 邮件服务测试"
	testBody := "<h1>测试成功！</h1><p>这是一封来自 PhotoSet 的测试邮件。</p>"

	if err := s.sendMail(cfg, cfg.Username, testSubject, testBody); err != nil {
		// 有些服务器会拒绝发给自己的测试邮件，这不代表配置错误
		if strings.Contains(err.Error(), "535") || strings.Contains(err.Error(), "authentication") {
			return false, fmt.Sprintf("认证失败，请检查用户名和密码: %v", err)
		}
		if strings.Contains(err.Error(), "connection") || strings.Contains(err.Error(), "refused") {
			return false, fmt.Sprintf("无法连接到SMTP服务器: %v", err)
		}
		// 其他错误可能是发送限制，但配置本身可能是正确的
		return false, fmt.Sprintf("连接失败: %v", err)
	}

	return true, "SMTP连接测试成功！"
}

// GetConfigInfo 获取配置信息（不包含密码）
func (s *MailService) GetConfigInfo() map[string]interface{} {
	cfg, err := s.getConfig()
	if err != nil {
		return map[string]interface{}{
			"configured": false,
			"error":      err.Error(),
		}
	}

	info := map[string]interface{}{
		"configured": cfg.Host != "",
		"host":       cfg.Host,
		"port":       cfg.Port,
		"username":   cfg.Username,
		"from":       cfg.From,
		"use_tls":    cfg.UseTLS,
		"password_set": cfg.Password != "",
	}

	return info
}
