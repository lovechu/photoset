package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"photoset/internal/database"
)

// CacheService Redis 缓存服务
type CacheService struct{}

func NewCacheService() *CacheService {
	return &CacheService{}
}

// 缓存 key 前缀
const (
	CachePrefixPhotosetList   = "photoset:list:"    // photoset:list:{page}:{pageSize}:{tag}:{keyword}
	CachePrefixPhotosetDetail = "photoset:detail:"   // photoset:detail:{id}
	CachePrefixTags           = "tags:all"            // tags:all
)

// Set 设置缓存
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if database.RedisClient == nil {
		return nil // 静默失败
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return database.RedisClient.Set(ctx, key, data, ttl).Err()
}

// Get 获取缓存，返回 error 表示缓存未命中
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	if database.RedisClient == nil {
		return fmt.Errorf("redis not initialized")
	}
	data, err := database.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err // redis.Nil 也返回，调用方判断
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存（用于主动失效）
func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
	if database.RedisClient == nil {
		return nil
	}
	if len(keys) == 0 {
		return nil
	}
	return database.RedisClient.Del(ctx, keys...).Err()
}

// DeleteByPattern 按通配符删除缓存（用于列表缓存失效）
func (s *CacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	if database.RedisClient == nil {
		return nil
	}
	iter := database.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) == 0 {
		return nil
	}
	return database.RedisClient.Del(ctx, keys...).Err()
}

// PhotosetListCacheKey 构建套图列表缓存 key
func PhotosetListCacheKey(page, pageSize int, tag, keyword string) string {
	return fmt.Sprintf("%s%d:%d:%s:%s", CachePrefixPhotosetList, page, pageSize, tag, keyword)
}

// PhotosetDetailCacheKey 构建套图详情缓存 key
func PhotosetDetailCacheKey(id uint) string {
	return fmt.Sprintf("%s%d", CachePrefixPhotosetDetail, id)
}
