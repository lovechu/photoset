package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/pkg/signurl"
)

const (
	paidStatusCachePrefix = "photoset:paid:"
	paidStatusCacheTTL    = 30 * time.Minute
)

// paidStatus 从 Redis 缓存 + 数据库获取套图付费状态
// 返回 true=付费，false=免费；缓存未命中时自动回源写缓存
func paidStatus(photosetID uint) bool {
	ctx := context.Background()

	// 1. 查 Redis
	key := fmt.Sprintf("%s%d", paidStatusCachePrefix, photosetID)
	if database.RedisClient != nil {
		if val, err := database.RedisClient.Get(ctx, key).Int(); err == nil {
			return val == 1
		}
	}

	// 2. Redis 未命中，查数据库
	var isFree int8
	database.GetMySQL().
		Table("photosets").
		Select("is_free").
		Where("id = ?", photosetID).
		Scan(&isFree)

	// is_free=0 → 付费(paid=true)，is_free=1 → 免费(paid=false)
	isPaid := isFree == 0

	// 3. 回源写 Redis（静默失败）
	if database.RedisClient != nil {
		if isPaid {
			database.RedisClient.Set(ctx, key, 1, paidStatusCacheTTL)
		} else {
			database.RedisClient.Set(ctx, key, 0, paidStatusCacheTTL)
		}
	}

	return isPaid
}

// SignVerify 图片签名验证中间件
//
// 路径规则：
//   - /uploads/covers/              → 封面图，免费公开访问
//   - /uploads/photos/{id}/...      → 根据套图 is_free 决定是否验签
//   - /uploads/images/...            → 旧路径兼容，无 sign 则放行
func SignVerify() gin.HandlerFunc {
	cfg := config.Load()

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 只拦截 /uploads/ 路径
		if !strings.HasPrefix(path, "/uploads/") {
			c.Next()
			return
		}

		// 封面图路径：永远免费
		if strings.HasPrefix(path, "/uploads/covers/") {
			c.Next()
			return
		}

		// 新版付费照片路径：/uploads/photos/{photosetID}/{MM}/{uuid}.ext
		// 根据套图是否付费决定是否验签
		if strings.HasPrefix(path, "/uploads/photos/") {
			// 从路径解析 photosetID：/uploads/photos/123/01/uuid.jpg → "123"
			parts := strings.Split(path, "/")
			if len(parts) >= 4 {
				if photosetID, err := strconv.ParseUint(parts[3], 10, 32); err == nil && photosetID > 0 {
					if paidStatus(uint(photosetID)) {
						// 付费套图：必须有签名且签名有效
						if !signurl.VerifyURL(c.Request.URL.String(), cfg.Storage.SignSecret) {
							c.AbortWithStatus(http.StatusForbidden)
							c.Writer.Write([]byte("403 Forbidden - Link expired or invalid"))
							return
						}
					}
					// 免费套图：无签名也放行
					c.Next()
					return
				}
			}
			// 路径格式不符，默认按有签名处理
			c.Next()
			return
		}

		// 其他旧路径（/uploads/images/ 等历史数据）兼容：无 sign 参数直接放行
		if c.Request.URL.Query().Get("sign") == "" {
			c.Next()
			return
		}

		if !signurl.VerifyURL(c.Request.URL.String(), cfg.Storage.SignSecret) {
			c.AbortWithStatus(http.StatusForbidden)
			c.Writer.Write([]byte("403 Forbidden - Link expired or invalid"))
			return
		}

		c.Next()
	}
}


