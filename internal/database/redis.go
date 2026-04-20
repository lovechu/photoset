package database

import (
	"context"
	"fmt"
	"photoset/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisAvailable bool // Redis 是否初始化成功

func InitRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		RedisAvailable = false
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	RedisAvailable = true
	return nil
}

// IsRedisAvailable 返回 Redis 是否可用
func IsRedisAvailable() bool {
	return RedisAvailable
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
