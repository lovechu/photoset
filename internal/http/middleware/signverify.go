package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/config"
	"photoset/internal/pkg/signurl"
)

// SignVerify 图片签名验证中间件
// 无 sign 参数的请求直接放行（免费图片）
// 有 sign 参数的验证签名和有效期
func SignVerify() gin.HandlerFunc {
	cfg := config.Load()

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 只拦截 /uploads/ 路径
		if !strings.HasPrefix(path, "/uploads/") {
			c.Next()
			return
		}

		query := c.Request.URL.Query()
		// 没有 sign 参数的请求直接放行（免费图片）
		if query.Get("sign") == "" {
			c.Next()
			return
		}

		// 有 sign 参数则验证
		if !signurl.VerifyURL(c.Request.URL.String(), cfg.Storage.SignSecret) {
			c.AbortWithStatus(http.StatusForbidden)
			c.Writer.Write([]byte("403 Forbidden - Link expired or invalid"))
			return
		}

		c.Next()
	}
}
