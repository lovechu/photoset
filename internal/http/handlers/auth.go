package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/jwt"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
)

type AuthHandler struct {
	userService    service.UserService
	captchaService service.CaptchaService
}

func NewAuthHandler(userService service.UserService, captchaService service.CaptchaService) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		captchaService: captchaService,
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
			"id":       user.ID,
			"nickname": user.Nickname,
			"email":    user.Email,
			"role":     user.Role,
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
