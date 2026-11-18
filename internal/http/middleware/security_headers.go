package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 为所有响应补充基础安全响应头。
// CSP 使用 Report-Only 便于灰度，避免直接破坏前端内联脚本/样式；
// 生产可根据实际情况改为正式 CSP 策略。
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
