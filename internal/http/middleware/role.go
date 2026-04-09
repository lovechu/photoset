package middleware

import (
	"net/http"
	"strings"

	"photoset/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireRoles 角色校验中间件
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role, exists := c.Get(RoleKey)
		if !exists {
			response.Error(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		// 检查角色是否在允许列表中
		userRole, ok := role.(string)
		if !ok {
			response.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		// 检查角色权限
		hasPermission := false
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(userRole, allowedRole) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Error(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
