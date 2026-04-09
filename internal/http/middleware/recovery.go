package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"photoset/internal/pkg/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				response.ServerError(c, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
