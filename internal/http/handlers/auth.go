package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/logger"
	"photoset/internal/pkg/email"
	"photoset/internal/pkg/jwt"
	"photoset/internal/pkg/response"
	"photoset/internal/repository"
	"photoset/internal/service"
)

type AuthHandler struct {
	userService          service.UserService
	captchaService       service.CaptchaService
	siteSettingRepo      *repository.SiteSettingRepository
	ipGeoService         *service.IPGeoService
	loginHistoryService  *service.LoginHistoryService
	userDeviceService    *service.UserDeviceService
	emailVerificationSvc *service.EmailVerificationService
}

func NewAuthHandler(
	userService service.UserService,
	captchaService service.CaptchaService,
	siteSettingRepo *repository.SiteSettingRepository,
	loginHistoryService *service.LoginHistoryService,
	userDeviceService *service.UserDeviceService,
	emailVerificationSvc *service.EmailVerificationService,
) *AuthHandler {
	return &AuthHandler{
		userService:          userService,
		captchaService:       captchaService,
		siteSettingRepo:      siteSettingRepo,
		ipGeoService:         service.NewIPGeoService(),
		loginHistoryService:  loginHistoryService,
		userDeviceService:    userDeviceService,
		emailVerificationSvc: emailVerificationSvc,
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
		logger.Warn("User registration failed", "email", req.Email, "error", err)
		response.Error(c, -1, err.Error())
		return
	}

	// 根据IP自动获取用户地理位置
	if clientIP := c.ClientIP(); clientIP != "" {
		if ipLocation := h.ipGeoService.GetLocation(clientIP); ipLocation != "" {
			h.userService.UpdateProfile(user.ID, "", "", "", ipLocation)
			user.IPLocation = ipLocation
		}
	}

	logger.Info("User registered successfully", "user_id", user.ID, "email", req.Email)

	response.Success(c, gin.H{
		"user": user,
	})
}

type LoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
	DeviceID    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
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
		logger.Warn("User login failed", "email", req.Email, "error", err)
		response.Error(c, -1, err.Error())
		return
	}

	// 根据IP更新用户地理位置（每次登录时更新）
	clientIP := c.ClientIP()
	var ipLocation string
	if clientIP != "" {
		ipLocation = h.ipGeoService.GetLocation(clientIP)
		if ipLocation != "" {
			h.userService.UpdateProfile(user.ID, "", "", "", ipLocation)
			user.IPLocation = ipLocation
		}
	}

	// 记录登录历史
	userAgent := c.Request.UserAgent()
	if h.loginHistoryService != nil {
		if err := h.loginHistoryService.CreateLoginHistory(
			user.ID, clientIP, ipLocation, userAgent, "password", true, "",
		); err != nil {
			logger.Warn("Failed to record login history", "user_id", user.ID, "error", err)
		}
	}

	// 注册或更新设备信息
	if h.userDeviceService != nil && req.DeviceID != "" {
		deviceID := req.DeviceID
		deviceName := req.DeviceName
		if err := h.userDeviceService.RegisterOrUpdateDevice(
			user.ID, deviceID, deviceName, userAgent, clientIP, ipLocation,
		); err != nil {
			logger.Warn("Failed to register device", "user_id", user.ID, "error", err)
		}
	}

	token, err := jwt.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		logger.Error("Failed to generate JWT token", "user_id", user.ID, "error", err)
		response.ServerError(c, "failed to generate token")
		return
	}
	logger.Info("User logged in successfully", "user_id", user.ID, "email", req.Email)

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"nickname":        user.Nickname,
			"email":           user.Email,
			"avatar":          user.Avatar,
			"bio":             user.Bio,
			"ip_location":     user.IPLocation,
			"level":           user.Level,
			"exp":             user.Exp,
			"circle_count":    user.CircleCount,
			"following_count": user.FollowingCount,
			"follower_count":  user.FollowerCount,
			"like_count":      user.LikeCount,
			"role":            user.Role,
			"status":          user.Status,
			"created_at":      user.CreatedAt,
			"updated_at":      user.UpdatedAt,
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
		logger.Error("Failed to get user profile", "user_id", userID, "error", err)
		response.ServerError(c, "failed to get user profile")
		return
	}

	// 每次访问个人页时，实时检查当前IP是否变化，若有变化则更新
	if clientIP := c.ClientIP(); clientIP != "" {
		if newLocation := h.ipGeoService.GetLocation(clientIP); newLocation != "" && newLocation != user.IPLocation {
			h.userService.UpdateProfile(user.ID, "", "", "", newLocation)
			user.IPLocation = newLocation
		}
	}

	response.Success(c, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"bio":        user.Bio,
			"ip_location": user.IPLocation,
			"level":      user.Level,
			"exp":        user.Exp,
			"circle_count": user.CircleCount,
			"following_count": user.FollowingCount,
			"follower_count": user.FollowerCount,
			"like_count": user.LikeCount,
			"role":       user.Role,
			"status":     user.Status,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	var req struct {
		DeviceID string `json:"device_id"`
	}
	// 忽略解析错误，device_id 可选
	_ = c.ShouldBindJSON(&req)

	// 如果有设备 ID，下线该设备
	if req.DeviceID != "" && h.userDeviceService != nil {
		if err := h.userDeviceService.DeactivateDevice(userID, req.DeviceID); err != nil {
			logger.Warn("Failed to deactivate device on logout", "user_id", userID, "error", err)
		}
	}

	logger.Info("User logged out", "user_id", userID)
	response.Success(c, gin.H{"message": "已登出"})
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

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	var req struct {
		Nickname   string `json:"nickname"`
		Bio        string `json:"bio"`
		Avatar     string `json:"avatar"`
		IPLocation  string `json:"ip_location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// DEBUG: 打印收到的参数，确认 avatar 是否传入
	log.Printf("[UpdateProfile] userID=%d nickname=%q bio=%q avatar=%q ip_location=%q",
		userID, req.Nickname, req.Bio, req.Avatar, req.IPLocation)

	user, err := h.userService.UpdateProfile(userID, req.Nickname, req.Bio, req.Avatar, req.IPLocation)
	if err != nil {
		log.Printf("[UpdateProfile] UpdateProfile error: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("[UpdateProfile] SUCCESS user.ID=%d avatar=%q", user.ID, user.Avatar)

	response.Success(c, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"bio":        user.Bio,
			"ip_location": user.IPLocation,
			"level":      user.Level,
			"exp":        user.Exp,
			"circle_count": user.CircleCount,
			"following_count": user.FollowingCount,
			"follower_count": user.FollowerCount,
			"like_count": user.LikeCount,
			"role":       user.Role,
			"status":     user.Status,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
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

// SendEmailCode 发送邮箱验证码（绑定邮箱或注册验证）
func (h *AuthHandler) SendEmailCode(c *gin.Context) {
	if h.emailVerificationSvc == nil {
		response.Error(c, http.StatusInternalServerError, "邮箱验证服务未初始化")
		return
	}

	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Purpose string `json:"purpose" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	if req.Purpose != "bind" && req.Purpose != "verify" {
		response.BadRequest(c, "无效的验证码用途")
		return
	}

	if err := h.emailVerificationSvc.SendVerificationCode(req.Email, req.Purpose); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "验证码已发送，请查收邮箱"})
}

// VerifyEmailCode 验证邮箱验证码并确认邮箱所有权
func (h *AuthHandler) VerifyEmailCode(c *gin.Context) {
	if h.emailVerificationSvc == nil {
		response.Error(c, http.StatusInternalServerError, "邮箱验证服务未初始化")
		return
	}

	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Code    string `json:"code" binding:"required"`
		Purpose string `json:"purpose" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	if err := h.emailVerificationSvc.VerifyCode(req.Email, req.Code, req.Purpose); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "验证成功"})
}

// BindEmail 绑定邮箱（需登录，验证码已通过 VerifyEmailCode 验证）
func (h *AuthHandler) BindEmail(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "请先登录")
		return
	}

	if h.emailVerificationSvc == nil {
		response.Error(c, http.StatusInternalServerError, "邮箱验证服务未初始化")
		return
	}

	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Code    string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	// 验证验证码
	if err := h.emailVerificationSvc.VerifyCode(req.Email, req.Code, "bind"); err != nil {
		response.Error(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	// 更新用户邮箱
	user, err := h.userService.UpdateEmail(userID, req.Email)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "邮箱绑定成功",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
