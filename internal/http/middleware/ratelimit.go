package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"photoset/internal/database"
	"photoset/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// RateLimit 限流中间件
// limit: 时间窗口内最大请求数
// window: 时间窗口（秒）
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端 IP（优先 X-Forwarded-For，再取 RemoteAddr）
		clientIP := c.ClientIP()

		key := fmt.Sprintf("ratelimit:%s:%s", strings.Replace(c.FullPath(), "/", ":", -1), clientIP)

		count, err := database.RedisClient.Incr(c.Request.Context(), key).Result()
		if err != nil {
			// Redis 出错时放行，不阻塞正常请求
			c.Next()
			return
		}

		if count == 1 {
			database.RedisClient.Expire(c.Request.Context(), key, window)
		}

		if count > int64(limit) {
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimit 登录接口专用限流：同一 IP 每分钟最多 10 次
func LoginRateLimit() gin.HandlerFunc {
	return RateLimit(10, 1*time.Minute)
}

// RegisterRateLimit 注册接口专用限流：同一 IP 每分钟最多 5 次
func RegisterRateLimit() gin.HandlerFunc {
	return RateLimit(5, 1*time.Minute)
}

// CaptchaRateLimit 验证码获取限流：同一 IP 每分钟最多 30 次
func CaptchaRateLimit() gin.HandlerFunc {
	return RateLimit(30, 1*time.Minute)
}
