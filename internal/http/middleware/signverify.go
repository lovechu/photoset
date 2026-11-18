package middleware

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"photoset/internal/config"
	"photoset/internal/pkg/signurl"
	"photoset/internal/service"
)

// SignVerify 图片签名验证中间件（fail-closed）
//
// 路径规则：
//   - /uploads/covers/   → 封面图（含付费套图封面），免费公开访问
//   - /uploads/videos/   → 视频，免费公开访问
//   - /uploads/avatars/  → 头像，免费公开访问
//   - /uploads/photos/   → 正文图片：按完整 URL 反查归属套图付费状态，付费必须验签
//   - 其他路径            → 默认要求签名（历史 /uploads/images/ 亦不再无脑放行）
//
// 安全要点：
//  1. 路径先规范化并拒绝含 ".." 的穿越，避免 FileServer path.Clean 后绕过前缀判断；
//  2. 付费状态用完整 URL 查库（历史路径第一段是时间戳，不能当套图ID解析）；
//  3. 查不到记录 / 查询出错 / 路径格式异常一律 fail-closed（要求签名），不再放行。
func SignVerify(cfg *config.Config) gin.HandlerFunc {
	forbidden := func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
		_, _ = c.Writer.Write([]byte("403 Forbidden - Link expired or invalid"))
	}

	return func(c *gin.Context) {
		rawPath := c.Request.URL.Path
		if !strings.HasPrefix(rawPath, "/uploads/") {
			c.Next()
			return
		}

		// 路径穿越防护：规范化后必须仍以 /uploads/ 开头，且原路径不得含 .. 段
		cleanPath := path.Clean(rawPath)
		if !strings.HasPrefix(cleanPath, "/uploads/") || hasDotDotSegment(rawPath) {
			forbidden(c)
			return
		}

		// 免费公开目录
		if strings.HasPrefix(cleanPath, "/uploads/covers/") ||
			strings.HasPrefix(cleanPath, "/uploads/videos/") ||
			strings.HasPrefix(cleanPath, "/uploads/avatars/") {
			c.Next()
			return
		}

		// 正文图片路径（含历史 /uploads/images/）：按完整 URL 反查归属套图的付费状态。
		// 历史路径第一段是时间戳，新路径第一段是套图ID，二者都不能靠路径数字判断付费，
		// 统一用 URL 查库，避免付费图片藏在 images 旧路径下绕过验签。
		if strings.HasPrefix(cleanPath, "/uploads/photos/") ||
			strings.HasPrefix(cleanPath, "/uploads/images/") {
			isPaid, exists, err := service.ImagePaidStatus(cleanPath)
			if err != nil || !exists {
				// fail-closed：未知图片或查询失败按付费处理，必须验签
				if !signurl.VerifyURL(c.Request.URL.String(), cfg.Storage.SignSecret) {
					forbidden(c)
					return
				}
				c.Next()
				return
			}
			if isPaid && !signurl.VerifyURL(c.Request.URL.String(), cfg.Storage.SignSecret) {
				forbidden(c)
				return
			}
			// 已确认免费套图的图片，放行
			c.Next()
			return
		}

		// 其他路径：默认要求签名
		if !signurl.VerifyURL(c.Request.URL.String(), cfg.Storage.SignSecret) {
			forbidden(c)
			return
		}
		c.Next()
	}
}

// hasDotDotSegment 检查路径是否包含 ".." 路径段（含被编码后已解码为 .. 的情况）
func hasDotDotSegment(p string) bool {
	for _, seg := range strings.Split(p, "/") {
		if seg == ".." {
			return true
		}
	}
	return false
}
