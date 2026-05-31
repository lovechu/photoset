package middleware

import (
	"time"

	"photoset/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

// RequestID 为每个请求生成唯一 ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDKey, requestID)
		c.Next()
	}
}

// Logger 结构化请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 获取 Request ID
		requestID, _ := c.Get(RequestIDKey)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		// 根据状态码选择日志级别
		logArgs := []any{
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"ip", clientIP,
			"body_size", bodySize,
			"request_id", requestID,
		}

		if query != "" {
			logArgs = append(logArgs, "query", query)
		}

		// 仅在 Debug 级别记录 User-Agent（减少日志量）
		// Debug 级别时添加更多详细信息
		logArgs = append(logArgs, "user_agent", userAgent)

		switch {
		case statusCode >= 500:
			logger.Error("HTTP request", logArgs...)
		case statusCode >= 400:
			logger.Warn("HTTP request", logArgs...)
		default:
			logger.Info("HTTP request", logArgs...)
		}
	}
}
