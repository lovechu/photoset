package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/pkg/response"
	"photoset/internal/service"
)

const OAuthTokenKey = "oauth_access_token"
const OAuthUserIDKey = "oauth_user_id"
const OAuthScopesKey = "oauth_scopes"

// OAuthAuth OAuth2 鉴权中间件 - 验证 OAuth2 访问令牌
func OAuthAuth(oauthService service.OAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractOAuthToken(c)
		if token == "" {
			response.Unauthorized(c, "missing access_token")
			c.Abort()
			return
		}

		// 验证访问令牌并获取用户信息
		userInfo, err := oauthService.GetUserInfo(token)
		if err != nil {
			response.Unauthorized(c, "invalid or expired access_token")
			c.Abort()
			return
		}

		// 获取用户授权的权限范围
		scopes, err := oauthService.GetUserScopes(token)
		if err != nil {
			response.Unauthorized(c, "invalid access_token")
			c.Abort()
			return
		}

		// 将用户信息和权限范围写入上下文
		c.Set(OAuthTokenKey, token)
		if userID, ok := userInfo["id"].(uint); ok {
			c.Set(OAuthUserIDKey, userID)
		}
		c.Set(OAuthScopesKey, scopes)

		c.Next()
	}
}

// extractOAuthToken 从请求中提取 OAuth2 访问令牌
func extractOAuthToken(c *gin.Context) string {
	// 优先从 Authorization header 提取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// 从 query 参数提取
	return c.Query("access_token")
}

// GetOAuthUserID 从上下文中获取 OAuth 用户 ID
func GetOAuthUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(OAuthUserIDKey)
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetOAuthScopes 从上下文中获取 OAuth 权限范围
func GetOAuthScopes(c *gin.Context) ([]string, bool) {
	scopes, exists := c.Get(OAuthScopesKey)
	if !exists {
		return nil, false
	}
	return scopes.([]string), true
}

// RequireScope 检查是否具有指定权限范围
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopes, exists := GetOAuthScopes(c)
		if !exists {
			response.Forbidden(c, "no scopes found")
			c.Abort()
			return
		}

		for _, s := range scopes {
			if s == scope {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient scope: "+scope)
		c.Abort()
	}
}