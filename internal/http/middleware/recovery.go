package middleware

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"photoset/internal/pkg/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC] 路径: %s %s\n", c.Request.Method, c.Request.URL.Path)
				fmt.Printf("[PANIC] 错误: %v\n", err)
				fmt.Printf("[PANIC] 堆栈:\n%s\n", string(debug.Stack()))
				log.Printf("Panic recovered: %v", err)
				response.ServerError(c, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
