package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/pkg/jwt"
	"photoset/internal/pkg/response"
)

const UserKey = "user_id"
const RoleKey = "user_role"

const (
	RoleGuest   = "guest"
	RoleUser    = "user"
	RoleMember  = "member"
	RoleCreator = "creator"
	RoleVIP     = "vip"
	RoleAdmin   = "admin"
)

// Auth 强制鉴权中间件 - 必须提供有效 token
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Unauthorized(c, "missing authorization token")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserKey, claims.UserID)
		c.Set(RoleKey, claims.Role)

		c.Next()
	}
}

// extractToken extracts JWT token from Authorization header or query parameter.
// 查询参数仅对 WebSocket 升级请求放行（WebSocket 握手时无法携带自定义 header，
// 只能用 URL 传 token），普通 HTTP 请求一律拒绝查询参数中的 token，
// 避免 token 泄露到服务器日志、代理日志、CDN 日志、Referer header 中。
func extractToken(c *gin.Context) string {
	// First try Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	// 查询参数仅对 WebSocket 升级请求放行（修复 AUTH-003）
	if strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
		return c.Query("token")
	}
	return ""
}

// OptionalAuth 可选鉴权中间件 - 没有按游客处理,有则解析
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			// 没有提供 token,按游客处理,直接放行
			c.Next()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			// token 无效或过期,按游客处理,直接放行
			c.Next()
			return
		}

		// token 有效,写入上下文
		c.Set(UserKey, claims.UserID)
		c.Set(RoleKey, claims.Role)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserKey)
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(RoleKey)
	if !exists {
		return "", false
	}
	return role.(string), true
}

// AdminOnly middleware - only allow admin users
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetUserRole(c)
		if !exists {
			response.Unauthorized(c, "please login first")
			c.Abort()
			return
		}

		if role != RoleAdmin {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

