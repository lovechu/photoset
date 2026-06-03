package handlers

import (
	"context"
	"time"

	"photoset/internal/database"
	"photoset/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	// 检查数据库连接
	dbStatus := "ok"
	db := database.GetMySQL()
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "error"
		}
	} else {
		dbStatus = "error"
	}

	// 检查 Redis 连接
	redisStatus := "ok"
	if database.IsRedisAvailable() && database.RedisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := database.RedisClient.Ping(ctx).Err(); err != nil {
			redisStatus = "error"
		}
	} else {
		redisStatus = "error"
	}

	// 判断整体健康状态
	overall := "ok"
	if dbStatus == "error" || redisStatus == "error" {
		overall = "degraded"
	}

	response.Success(c, gin.H{
		"status":   overall,
		"database": dbStatus,
		"redis":    redisStatus,
		"time":     time.Now().Format("2006-01-02 15:04:05"),
	})
}

