package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"photoset/internal/database"
)

// CacheService Redis 缓存服务
type CacheService struct{}

// 缓存常量
const (
	CachePrefixPhotosetList   = "photoset:list:"
	CachePrefixPhotosetDetail = "photoset:detail:"
	CachePrefixTags           = "tags:all"
)

func NewCacheService() *CacheService {
	return &CacheService{}
}

// Set 设置缓存
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !database.IsRedisAvailable() {
		log.Printf("[Cache] WARNING: Redis unavailable, cache SET skipped (key=%s)", key)
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return database.RedisClient.Set(ctx, key, data, ttl).Err()
}

// Get 获取缓存，返回 error 表示缓存未命中
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	if !database.IsRedisAvailable() {
		return fmt.Errorf("redis unavailable")
	}
	data, err := database.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err // redis.Nil 也返回，调用方判断
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存（用于主动失效）
func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
	if !database.IsRedisAvailable() {
		return nil
	}
	if len(keys) == 0 {
		return nil
	}
	return database.RedisClient.Del(ctx, keys...).Err()
}

// DeleteByPattern 按通配符删除缓存（用于列表缓存失效）
func (s *CacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	if !database.IsRedisAvailable() {
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

// SetBool 设置布尔值缓存
func (s *CacheService) SetBool(ctx context.Context, key string, val bool, ttl time.Duration) error {
	if !database.IsRedisAvailable() {
		log.Printf("[Cache] WARNING: Redis unavailable, cache SETBool skipped (key=%s)", key)
		return nil
	}
	return database.RedisClient.Set(ctx, key, val, ttl).Err()
}

// GetBool 获取布尔值缓存，error 表示未命中
func (s *CacheService) GetBool(ctx context.Context, key string) (bool, error) {
	if !database.IsRedisAvailable() {
		return false, fmt.Errorf("redis unavailable")
	}
	val, err := database.RedisClient.Get(ctx, key).Int()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

// PhotosetListCacheKey 构建套图列表缓存 key
func PhotosetListCacheKey(page, pageSize int, tag, keyword string) string {
	return fmt.Sprintf("%s%d:%d:%s:%s", CachePrefixPhotosetList, page, pageSize, tag, keyword)
}

// PhotosetDetailCacheKey 构建套图详情缓存 key
func PhotosetDetailCacheKey(id uint) string {
	return fmt.Sprintf("%s%d", CachePrefixPhotosetDetail, id)
}

// PhotosetAdvancedListCacheKey 构建高级筛选套图列表缓存 key
func PhotosetAdvancedListCacheKey(
	page, pageSize int,
	tag, keyword string,
	userID uint,
	onlyMine bool,
	category string,
	priceMin, priceMax float64,
	isFree *bool,
	sortBy, timeRange string,
	filterUserID uint,
) string {
	// 基础参数
	key := fmt.Sprintf("%s%d:%d:%s:%s", CachePrefixPhotosetList, page, pageSize, tag, keyword)
	
	// 用户相关参数
	if onlyMine && userID > 0 {
		key = fmt.Sprintf("%s:mine:%d", key, userID)
	}
	
	// 分类参数
	if category != "" {
		key = fmt.Sprintf("%s:cat:%s", key, category)
	}
	
	// 价格参数
	if priceMin > 0 {
		key = fmt.Sprintf("%s:min%.2f", key, priceMin)
	}
	if priceMax > 0 {
		key = fmt.Sprintf("%s:max%.2f", key, priceMax)
	}
	
	// 是否免费参数
	if isFree != nil {
		if *isFree {
			key = fmt.Sprintf("%s:free", key)
		} else {
			key = fmt.Sprintf("%s:paid", key)
		}
	}
	
	// 排序参数
	if sortBy != "" {
		key = fmt.Sprintf("%s:sort:%s", key, sortBy)
	}
	
	// 时间范围参数
	if timeRange != "" {
		key = fmt.Sprintf("%s:time:%s", key, timeRange)
	}
	
	// 筛选特定用户参数
	if filterUserID > 0 {
		key = fmt.Sprintf("%s:user:%d", key, filterUserID)
	}
	
	return key
}
