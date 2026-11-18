package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"photoset/internal/database"
)

// 付费状态缓存常量
const (
	CachePrefixPaidStatus = "photoset:paid:"
	PaidStatusTTL         = 30 * time.Minute // 图片归属/付费状态缓存 30 分钟
	MissingStatusTTL      = 10 * time.Second // "图片未命中"缓存只保留极短时间，避免上传预览早于发帖导致长期 403
)

// imagePaidCacheKey 构建单张图片的付费状态缓存 key（供查询与主动失效共用）
func imagePaidCacheKey(imageURL string) string {
	return "photoset:imgpaid:" + imageURL
}

// InvalidateImagePaidCache 在帖子/回复发布后主动清除正文里所有图片的付费状态缓存。
// 背景：图片上传后会先被预览请求，此时正文尚未写入，ImagePaidStatus 会把该图缓存为
// "missing"（fail-closed 返回 403）。发帖/回帖成功后调用本函数失效这些 key，
// 下次访问会重新走库判断（此时 content 已包含图片 URL，识别为社区图片并放行）。
// 缓存键使用 cleanPath（/uploads/... 相对路径），这里同时覆盖绝对 URL 与相对路径两种写法。
func InvalidateImagePaidCache(content string) {
	if content == "" {
		return
	}
	cache := NewCacheService()
	ctx := context.Background()

	seen := make(map[string]bool)
	var keys []string
	// 匹配 /uploads/... 图片地址（Markdown/裸 URL/HTML 里的均会被这段命中）
	re := regexp.MustCompile(`/uploads/[^\s)"'<>\\]+?\.(?:jpg|jpeg|png|gif|webp|avif|bmp)(?:[?#][^\s)"'<>]*)?`)
	for _, m := range re.FindAllString(content, -1) {
		for _, u := range []string{m, "/" + strings.TrimPrefix(m, "/")} {
			if seen[u] {
				continue
			}
			seen[u] = true
			keys = append(keys, imagePaidCacheKey(u))
		}
	}
	if len(keys) > 0 {
		_ = cache.Delete(ctx, keys...)
	}
}

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
	key := imagePaidCacheKey(imageURL)

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
			// 社区帖子/回复正文里引用的图片（App 发帖把图片裸 URL 放进 content）
			// 属于公开内容，不绑定任何付费套图，应放行显示。
			// 若走到这里仍未命中 photos / cover，再查是否被社区内容引用。
			if isCommunityImage(imageURL) {
				_ = cache.Set(ctx, key, "0", PaidStatusTTL)
				return false, true, nil
			}
			// missing 只缓存极短时间：图片常在发帖前被上传预览请求过，
			// 若缓存 30 分钟会导致发帖后 30 分钟内一直 403（缓存早于数据写入）。
			_ = cache.Set(ctx, key, "missing", MissingStatusTTL)
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

// isCommunityImage 判断图片 URL 是否被社区帖子/回复正文引用。
// App 发帖把上传后的图片裸 URL 直接写入 content（如 https://域名/uploads/photos/...），
// 这些图片不绑定任何套图，属于公开社区内容，必须放行，否则会被付费签名策略误杀返回 403。
// 调用方已做 30 分钟缓存，这里的两条 LIKE 查询只在缓存 miss 且非套图图片时执行。
func isCommunityImage(imageURL string) bool {
	if imageURL == "" {
		return false
	}
	// 兼容绝对 URL 与相对路径两种存储形式
	needles := []string{imageURL}
	if idx := strings.Index(imageURL, "/uploads/"); idx >= 0 {
		needles = append(needles, imageURL[idx:])
	}

	db := database.GetMySQL()
	for _, needle := range needles {
		for _, table := range []string{"posts", "post_replies"} {
			var count int64
			if err := db.Table(table).
				Where("content LIKE ?", "%"+needle+"%").
				Limit(1).
				Count(&count).Error; err == nil && count > 0 {
				return true
			}
		}
	}
	return false
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
