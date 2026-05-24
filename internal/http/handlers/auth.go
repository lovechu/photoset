package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/email"
	"photoset/internal/pkg/jwt"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"photoset/internal/service"
)

type AuthHandler struct {
	userService         service.UserService
	captchaService      service.CaptchaService
	siteSettingRepo     *repository.SiteSettingRepository
}

func NewAuthHandler(userService service.UserService, captchaService service.CaptchaService, siteSettingRepo *repository.SiteSettingRepository) *AuthHandler {
	return &AuthHandler{
		userService:     userService,
		captchaService:  captchaService,
		siteSettingRepo: siteSettingRepo,
	}
}

type RegisterRequest struct {
	Nickname    string `json:"nickname" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	// 验证图形验证码
	if !h.captchaService.Verify(req.CaptchaID, req.CaptchaCode, "register") {
		response.Error(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	user, err := h.userService.Register(req.Nickname, req.Email, req.Password)
	if err != nil {
		response.Error(c, -1, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user": user,
	})
}

type LoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	// 验证图形验证码
	if !h.captchaService.Verify(req.CaptchaID, req.CaptchaCode, "login") {
		response.Error(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, -1, err.Error())
		return
	}

	token, err := jwt.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		response.ServerError(c, "failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"role":       user.Role,
			"status":     user.Status,
			"created_at": user.CreatedAt,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		// 没有用户登录，返回空的用户信息（不是401错误）
		response.Success(c, gin.H{
			"user": nil,
		})
		return
	}

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.ServerError(c, "failed to get user profile")
		return
	}

	response.Success(c, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"role":       user.Role,
			"status":     user.Status,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// ChangePassword 用户修改自己的密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误，新密码长度不能少于6位")
		return
	}

	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "密码修改成功"})
}

// ForgotPassword 请求密码重置（发送重置邮件）
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		CaptchaID   string `json:"captcha_id" binding:"required"`
		CaptchaCode string `json:"captcha_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 验证图形验证码
	if !h.captchaService.Verify(req.CaptchaID, req.CaptchaCode, "forgot") {
		response.Error(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	// 检查邮件配置
	settings, _ := h.siteSettingRepo.GetAll()
	mailCfg := email.GetEmailConfigFromSettings(settings)
	mailCfg.Host = email.NormalizeHost(mailCfg.Host)
	if !mailCfg.IsConfigured() {
		response.Error(c, http.StatusBadRequest, "邮件服务未配置，请联系管理员配置 SMTP")
		return
	}

	// 生成重置 token
	token, err := h.userService.RequestPasswordReset(req.Email)
	if err != nil {
		// 不暴露具体原因（防止用户枚举）
		// 但如果是"未注册"的错误还是要告知
		if strings.Contains(err.Error(), "未注册") || strings.Contains(err.Error(), "禁用") {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "发送重置邮件失败，请稍后重试")
		return
	}

	// 获取站点名称
	siteName := settings["site_title"]
	if siteName == "" {
		siteName = "PhotoSet"
	}

	// 构建重置 URL
	siteURL := settings["site_url"]
	if siteURL == "" {
		siteURL = c.Request.Header.Get("Origin")
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}
	// 去掉末尾的斜杠
	siteURL = strings.TrimRight(siteURL, "/")

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", siteURL, token)

	// 构建邮件内容
	body := email.BuildResetPasswordBody(siteName, resetURL)
	subject := fmt.Sprintf("[%s] 密码重置请求", siteName)

	// 发送邮件
	if err := email.SendMail(mailCfg, req.Email, subject, body); err != nil {
		response.Error(c, http.StatusInternalServerError, "发送重置邮件失败："+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "重置邮件已发送，请查收邮箱"})
}

// ResetPasswordByToken 通过 token 重置密码
func (h *AuthHandler) ResetPasswordByToken(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.userService.ResetPasswordByToken(req.Token, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "密码重置成功，请使用新密码登录"})
}

// CheckEmailConfig 检查邮件配置是否可用（公开接口，前端在忘记密码页面判断是否显示）
func (h *AuthHandler) CheckEmailConfig(c *gin.Context) {
	settings, _ := h.siteSettingRepo.GetAll()
	mailCfg := email.GetEmailConfigFromSettings(settings)
	mailCfg.Host = email.NormalizeHost(mailCfg.Host)

	response.Success(c, gin.H{
		"configured": mailCfg.IsConfigured(),
	})
}
