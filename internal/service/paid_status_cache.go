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
	PaidStatusTTL         = 30 * time.Minute // 缓存 30 分钟
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

type freeStatus struct {
	IsFree int8 `gorm:"column:is_free"`
}

// ImagePaidStatus 根据图片 URL 反查其归属套图的付费状态。
// 历史图片路径第一段是时间戳（非套图ID），新路径第一段才是套图ID，
// 因此不能靠解析路径数字判断，必须用完整 URL 查 photos 表归属。
// 返回值：isPaid（是否付费）、exists（图片记录是否存在）、error。
// 图片不存在或查询出错时调用方必须 fail-closed（要求签名或拒绝）。
func ImagePaidStatus(imageURL string) (isPaid bool, exists bool, err error) {
	ctx := context.Background()
	cache := NewCacheService()
	key := "photoset:imgpaid:" + imageURL

	// 缓存存 "1"=付费 / "0"=免费 / "missing"=不存在
	var cached string
	if e := cache.Get(ctx, key, &cached); e == nil {
		switch cached {
		case "1":
			return true, true, nil
		case "0":
			return false, true, nil
		case "missing":
			return false, false, nil
		}
	}

	var result freeStatus
	err = database.GetMySQL().
		Table("photos").
		Select("photosets.is_free").
		Joins("JOIN photosets ON photosets.id = photos.photoset_id").
		Where("photos.url = ? AND photos.deleted_at IS NULL", imageURL).
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return false, false, fmt.Errorf("查询图片归属失败: %w", err)
	}
	// Scan 未找到记录时 result.IsFree 为零值 0，无法区分"付费"与"不存在"，
	// 因此单独 Count 确认存在性。
	var count int64
	if err = database.GetMySQL().
		Table("photos").
		Where("url = ? AND deleted_at IS NULL", imageURL).
		Count(&count).Error; err != nil {
		return false, false, fmt.Errorf("查询图片存在性失败: %w", err)
	}
	if count == 0 {
		// 历史数据中封面图存在 photosets.cover 而非 photos 表（如 /uploads/images/...）。
		// 封面按设计公开，查到则按封面归属套图的付费状态返回，避免误伤正常访问。
		var cover freeStatus
		err = database.GetMySQL().
			Table("photosets").
			Select("is_free").
			Where("cover = ? AND deleted_at IS NULL", imageURL).
			Limit(1).
			Scan(&cover).Error
		if err != nil {
			return false, false, fmt.Errorf("查询封面归属失败: %w", err)
		}
		// 用 cover 字段是否为空判断是否命中（is_free=0 也可能是命中）
		var coverCount int64
		if err = database.GetMySQL().
			Table("photosets").
			Where("cover = ? AND deleted_at IS NULL", imageURL).
			Count(&coverCount).Error; err != nil {
			return false, false, fmt.Errorf("查询封面存在性失败: %w", err)
		}
		if coverCount == 0 {
			_ = cache.Set(ctx, key, "missing", PaidStatusTTL)
			return false, false, nil
		}
		isPaid = cover.IsFree == 0
		val := "0"
		if isPaid {
			val = "1"
		}
		_ = cache.Set(ctx, key, val, PaidStatusTTL)
		return isPaid, true, nil
	}

	isPaid = result.IsFree == 0
	val := "0"
	if isPaid {
		val = "1"
	}
	_ = cache.Set(ctx, key, val, PaidStatusTTL)
	return isPaid, true, nil
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
