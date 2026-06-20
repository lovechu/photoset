package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/pkg/email"
	"time"
)

// EmailVerificationService 邮箱验证码服务
type EmailVerificationService struct {
	siteSettingRepo interface {
		GetAll() (map[string]string, error)
	}
}

func NewEmailVerificationService(siteSettingRepo interface {
	GetAll() (map[string]string, error)
}) *EmailVerificationService {
	return &EmailVerificationService{
		siteSettingRepo: siteSettingRepo,
	}
}

// GenerateCode 生成6位随机数字验证码
// 使用 crypto/rand 确保密码学安全，防止可预测攻击
func GenerateCode() string {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// crypto/rand 极少失败；若失败则用时间戳兜底（仍比 math/rand 安全）
		n = big.NewInt(time.Now().UnixNano() % 900000)
	}
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

// SendVerificationCode 发送邮箱验证码
func (s *EmailVerificationService) SendVerificationCode(toEmail string, purpose string) error {
	db := database.GetMySQL()

	// 检查发送频率：同一邮箱同一用途 60 秒内只能发送一次
	var recentCode domain.EmailVerificationCode
	err := db.Where("email = ? AND purpose = ? AND created_at > ?", toEmail, purpose, time.Now().Add(-60*time.Second)).First(&recentCode).Error
	if err == nil {
		return errors.New("验证码已发送，请60秒后重试")
	}

	// 检查每天发送上限（同一邮箱同一用途每天最多 10 次）
	var todayCount int64
	today := time.Now().Truncate(24 * time.Hour)
	db.Model(&domain.EmailVerificationCode{}).
		Where("email = ? AND purpose = ? AND created_at >= ?", toEmail, purpose, today).
		Count(&todayCount)
	if todayCount >= 10 {
		return errors.New("今日发送次数已达上限，请明天再试")
	}

	// 将该邮箱之前未使用的验证码全部标记为已使用（确保同时只有一个有效验证码）
	db.Model(&domain.EmailVerificationCode{}).
		Where("email = ? AND purpose = ? AND used = ?", toEmail, purpose, false).
		Update("used", true)

	// 生成验证码
	code := GenerateCode()

	// 保存到数据库（10 分钟有效）
	verificationCode := &domain.EmailVerificationCode{
		Email:   toEmail,
		Code:    code,
		Used:    false,
		Expire:  time.Now().Add(10 * time.Minute),
		Purpose: purpose,
	}
	if err := db.Create(verificationCode).Error; err != nil {
		return errors.New("保存验证码失败")
	}

	// 获取邮件配置
	settings, _ := s.siteSettingRepo.GetAll()
	mailCfg := email.GetEmailConfigFromSettings(settings)
	mailCfg.Host = email.NormalizeHost(mailCfg.Host)
	if !mailCfg.IsConfigured() {
		return errors.New("邮件服务未配置")
	}

	// 获取站点名称
	siteName := settings["site_title"]
	if siteName == "" {
		siteName = "PhotoSet"
	}

	// 构建邮件内容
	subject := fmt.Sprintf("[%s] 邮箱验证码", siteName)
	body := BuildVerificationCodeBody(siteName, code, purpose)

	// 发送邮件
	if err := email.SendMail(mailCfg, toEmail, subject, body); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// VerifyCode 验证邮箱验证码
// 包含暴力破解保护：同一验证码最多尝试 5 次，超过后自动失效
func (s *EmailVerificationService) VerifyCode(email, code, purpose string) error {
	db := database.GetMySQL()

	// 查询该邮箱+用途最近一条未使用、未过期的验证码
	var verificationCode domain.EmailVerificationCode
	err := db.Where("email = ? AND purpose = ? AND used = ? AND expire > ?",
		email, purpose, false, time.Now()).
		Order("created_at DESC").
		First(&verificationCode).Error
	if err != nil {
		return errors.New("验证码错误或已过期")
	}

	// 暴力破解保护：超过最大尝试次数则失效
	if verificationCode.Attempts >= 5 {
		db.Model(&verificationCode).Update("used", true)
		return errors.New("尝试次数过多，请重新获取验证码")
	}

	// 验证码不匹配，增加尝试次数
	if verificationCode.Code != code {
		db.Model(&verificationCode).Update("attempts", verificationCode.Attempts+1)
		return errors.New("验证码错误或已过期")
	}

	// 验证成功，标记为已使用
	db.Model(&verificationCode).Update("used", true)

	return nil
}

// BuildVerificationCodeBody 构建验证码邮件内容
func BuildVerificationCodeBody(siteName, code, purpose string) string {
	purposeText := "绑定邮箱"
	if purpose == "verify" {
		purposeText = "邮箱验证"
	}

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
    <h2 style="margin:0 0 16px;color:#333;font-size:20px;">%s验证码</h2>
    <p style="color:#666;line-height:1.6;margin:0 0 24px;">您好，您正在进行%s操作。请使用以下验证码完成验证：</p>
    <div style="text-align:center;margin:32px 0;">
      <div style="display:inline-block;padding:16px 40px;background:#f0f0ff;border-radius:12px;font-size:32px;font-weight:bold;color:#667eea;letter-spacing:8px;">%s</div>
    </div>
    <p style="color:#999;font-size:13px;line-height:1.6;margin:24px 0 0;">此验证码有效期为 10 分钟。如非本人操作，请忽略此邮件。</p>
  </div>
  <div style="border-top:1px solid #eee;padding:16px 32px;text-align:center;">
    <p style="margin:0;color:#bbb;font-size:12px;">此邮件由系统自动发送，请勿回复。</p>
  </div>
</div>
</body>
</html>`, siteName, purposeText, purposeText, code)
}
