package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"gopkg.in/gomail.v2"
)

// Config 邮件配置（从站点设置读取）
type Config struct {
	Host     string // SMTP 服务器地址
	Port     int    // SMTP 端口
	Username string // SMTP 用户名（发件邮箱）
	Password string // SMTP 密码或授权码
	FromName string // 发件人名称
	FromAddr string // 发件人地址（如不填则用 Username）
}

// SendMail 发送邮件
func SendMail(cfg *Config, to, subject, body string) error {
	if cfg.Host == "" || cfg.Username == "" {
		return fmt.Errorf("邮件服务未配置")
	}

	fromAddr := cfg.FromAddr
	if fromAddr == "" {
		fromAddr = cfg.Username
	}

	fromName := cfg.FromName
	if fromName == "" {
		fromName = "PhotoSet"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromAddr, fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// BuildResetPasswordBody 构建密码重置邮件内容
func BuildResetPasswordBody(siteName, resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="margin:0;padding:0;background:#f5f5f5;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<div style="max-width:600px;margin:40px auto;background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);">
  <div style="background:linear-gradient(135deg,#667eea 0%%,#764ba2 100%%);padding:32px;text-align:center;">
    <h1 style="margin:0;color:#fff;font-size:24px;">%s</h1>
  </div>
  <div style="padding:32px;">
    <h2 style="margin:0 0 16px;color:#333;font-size:20px;">密码重置请求</h2>
    <p style="color:#666;line-height:1.6;margin:0 0 24px;">您好，我们收到了您的密码重置请求。请点击下方按钮重置密码：</p>
    <div style="text-align:center;margin:32px 0;">
      <a href="%s" style="display:inline-block;padding:12px 32px;background:linear-gradient(135deg,#667eea 0%%,#764ba2 100%%);color:#fff;text-decoration:none;border-radius:8px;font-size:16px;font-weight:500;">重置密码</a>
    </div>
    <p style="color:#999;font-size:13px;line-height:1.6;margin:24px 0 0;">此链接有效期为 30 分钟。如果您没有发送此请求，请忽略此邮件。</p>
  </div>
  <div style="border-top:1px solid #eee;padding:16px 32px;text-align:center;">
    <p style="margin:0;color:#bbb;font-size:12px;">此邮件由系统自动发送，请勿回复。</p>
  </div>
</div>
</body>
</html>`, siteName, resetURL)
}

// CheckSMTPConfig 检查 SMTP 配置是否完整
func CheckSMTPConfig(cfg *Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("SMTP 服务器地址未配置")
	}
	if cfg.Username == "" {
		return fmt.Errorf("SMTP 用户名未配置")
	}
	if cfg.Password == "" {
		return fmt.Errorf("SMTP 密码未配置")
	}
	return nil
}

// TestConnection 测试 SMTP 连接
func TestConnection(cfg *Config) error {
	if err := CheckSMTPConfig(cfg); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接 SMTP 服务器失败: %w", err)
	}
	defer client.Close()

	// 验证认证信息
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 认证失败: %w", err)
	}

	return nil
}

// GetEmailConfigFromSettings 从站点设置 map 中提取邮件配置
func GetEmailConfigFromSettings(settings map[string]string) *Config {
	port := 587
	if p := settings["smtp_port"]; p != "" {
		fmt.Sscanf(p, "%d", &port)
	}
	return &Config{
		Host:     settings["smtp_host"],
		Port:     port,
		Username: settings["smtp_user"],
		Password: settings["smtp_password"],
		FromName: settings["site_title"],
		FromAddr: settings["smtp_from_addr"],
	}
}

// IsConfigured 检查邮件配置是否可用
func (cfg *Config) IsConfigured() bool {
	return cfg.Host != "" && cfg.Username != "" && cfg.Password != ""
}

// NormalizeHost 去除 host 中的协议前缀（用户可能填 smtp://xxx）
func NormalizeHost(host string) string {
	host = strings.TrimSpace(host)
	host = strings.TrimPrefix(host, "smtp://")
	host = strings.TrimPrefix(host, "smtps://")
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	return host
}
