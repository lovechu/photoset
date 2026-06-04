package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"photoset/internal/config"
	"photoset/internal/http/middleware"
	"photoset/internal/repository"
)

type ShareHandler struct {
	photosetRepo *repository.PhotoSetRepository
	cfg          *config.Config
}

func NewShareHandler(photosetRepo *repository.PhotoSetRepository, cfg *config.Config) *ShareHandler {
	return &ShareHandler{
		photosetRepo: photosetRepo,
		cfg:          cfg,
	}
}

// @Summary      生成套图分享链接
// @Description  为指定套图生成带签名的分享链接
// @Tags         Share
// @Accept       json
// @Produce      json
// @Param        id path string true "套图ID"
// @Param        source query string false "分享来源(app/web/miniapp)，默认app"
// @Success      200 {object} response.Response{data=object} "分享链接"
// @Failure      400 {object} response.Response "缺少套图ID"
// @Failure      404 {object} response.Response "套图不存在"
// @Security     BearerAuth
// @Router       /api/photosets/{id}/share-link [get]
// GenerateShareLink 生成套图分享链接
func (h *ShareHandler) GenerateShareLink(c *gin.Context) {
	photosetID := c.Param("id")
	if photosetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少套图ID"})
		return
	}

	// 验证套图存在
	photosetIDUint, err := strconv.ParseUint(photosetID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套图ID"})
		return
	}
	_, err = h.photosetRepo.FindByID(uint(photosetIDUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "套图不存在"})
		return
	}

	// 获取分享来源
	source := c.Query("source") // app, web, miniapp
	if source == "" {
		source = "app"
	}

	// 生成带签名的分享链接
	shareLink := h.generateSignedLink(photosetID, source)

	// 记录分享行为（可选，用于统计）
	userID, _ := middleware.GetUserID(c)
	if userID > 0 {
		// 可以在这里记录分享统计
		_ = userID
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"share_link": shareLink,
			"photoset_id": photosetID,
			"source": source,
		},
	})
}

// generateSignedLink 生成带签名的分享链接
func (h *ShareHandler) generateSignedLink(photosetID, source string) string {
	// 获取前端域名
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://photoset.app"
	}

	// 生成签名
	secret := os.Getenv("SHARE_SECRET")
	if secret == "" {
		secret = "photoset-share-secret-2024"
	}

	timestamp := time.Now().Unix()
	message := fmt.Sprintf("%s:%s:%s:%d", photosetID, source, secret, timestamp)
	signature := generateHMAC(message, secret)

	// 构建分享链接
	shareLink := fmt.Sprintf("%s/share/%s?s=%s&t=%d&sig=%s",
		frontendURL, photosetID, source, timestamp, signature)

	return shareLink
}

// generateHMAC 生成 HMAC-SHA256 签名
func generateHMAC(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// @Summary      验证分享链接
// @Description  验证分享链接的签名是否有效
// @Tags         Share
// @Accept       json
// @Produce      json
// @Param        id path string true "套图ID"
// @Param        s query string true "分享来源"
// @Param        t query string true "时间戳"
// @Param        sig query string true "签名"
// @Success      200 {object} response.Response{data=object} "验证通过"
// @Failure      400 {object} response.Response "无效的分享链接或签名无效"
// @Router       /api/share/{id} [get]
// VerifyShareLink 验证分享链接签名
func (h *ShareHandler) VerifyShareLink(c *gin.Context) {
	photosetID := c.Param("id")
	source := c.Query("s")
	timestamp := c.Query("t")
	signature := c.Query("sig")

	if photosetID == "" || timestamp == "" || signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的分享链接"})
		return
	}

	// 验证签名
	secret := os.Getenv("SHARE_SECRET")
	if secret == "" {
		secret = "photoset-share-secret-2024"
	}

	message := fmt.Sprintf("%s:%s:%s:%s", photosetID, source, secret, timestamp)
	expectedSignature := generateHMAC(message, secret)

	if signature != expectedSignature {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "分享链接签名无效"})
		return
	}

	// 检查链接是否过期（可选，例如30天过期）
	// 这里可以添加时间戳过期检查

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"photoset_id": photosetID,
			"source":      source,
			"valid":       true,
		},
	})
}