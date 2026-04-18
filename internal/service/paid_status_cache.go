package service

import (
	"context"
	"fmt"
	"time"

	"photoset/internal/database"
)

// 付费状态缓存常量
const (
	CachePrefixPaidStatus = "photoset:paid:"
	PaidStatusTTL          = 30 * time.Minute // 缓存 30 分钟
)

// PaidStatusCacheKey 构建付费状态缓存 key
func PaidStatusCacheKey(photosetID uint) string {
	return fmt.Sprintf("%s%d", CachePrefixPaidStatus, photosetID)
}

// IsPaid 判断套图是否为付费套图
// 查询顺序：Redis → 数据库 → 回源写缓存 → 返回
// Redis 不可用时静默回退到数据库
func IsPaid(photosetID uint) (bool, error) {
	ctx := context.Background()
	cache := NewCacheService()

	// 1. 查 Redis
	if val, err := cache.GetBool(ctx, PaidStatusCacheKey(photosetID)); err == nil {
		return val, nil
	}

	// 2. Redis 未命中，查数据库
	var result struct {
		IsFree int8 `gorm:"column:is_free"`
	}
	if err := database.GetMySQL().
		Table("photosets").
		Select("is_free").
		Where("id = ?", photosetID).
		Scan(&result).Error; err != nil {
		return false, fmt.Errorf("查询付费状态失败: %w", err)
	}

	// is_free=0 表示付费（收费套图），is_free=1 表示免费
	isPaid := result.IsFree == 0

	// 3. 回源写 Redis（静默失败）
	cache.SetBool(ctx, PaidStatusCacheKey(photosetID), isPaid, PaidStatusTTL)

	return isPaid, nil
}

// SetPaidStatus 主动写入付费状态缓存（套图改价/创建时调用）
func SetPaidStatus(photosetID uint, isPaid bool) {
	ctx := context.Background()
	cache := NewCacheService()
	cache.SetBool(ctx, PaidStatusCacheKey(photosetID), isPaid, PaidStatusTTL)
}

// InvalidatePaidStatus 删除付费状态缓存（套图删除时调用）
func InvalidatePaidStatus(photosetID uint) {
	ctx := context.Background()
	cache := NewCacheService()
	cache.Delete(ctx, PaidStatusCacheKey(photosetID))
}
