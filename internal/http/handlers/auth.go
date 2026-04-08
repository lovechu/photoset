package handlers

import (
	"github.com/gin-gonic/gin"
	"photoset/internal/http/middleware"
	"photoset/internal/pkg/jwt"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
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
		response.Unauthorized(c, "user not found in context")
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
