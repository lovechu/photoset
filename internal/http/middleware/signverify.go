package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/config"
	"photoset/internal/pkg/signurl"
	"photoset/internal/service"
)

// SignVerify 图片签名验证中间件
//
// 路径规则：
//   - /uploads/covers/              → 封面图，免费公开访问
//   - /uploads/photos/{id}/...      → 根据套图 is_free 决定是否验签
//   - /uploads/images/...            → 旧路径兼容，无 sign 则放行
func SignVerify(cfg *config.Config) gin.HandlerFunc {
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
					isPaid, err := service.IsPaid(uint(photosetID))
					if err != nil || isPaid { // 出错时保守处理：视为付费
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


