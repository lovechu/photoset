package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/pkg/jwt"
	"photoset/internal/pkg/response"
)

const UserKey = "user_id"
const RoleKey = "user_role"

// Auth 强制鉴权中间件 - 必须提供有效 token
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
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

// OptionalAuth 可选鉴权中间件 - 没有按游客处理,有则解析
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有提供 token,按游客处理,直接放行
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// token 格式错误,按游客处理,直接放行
			c.Next()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
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

