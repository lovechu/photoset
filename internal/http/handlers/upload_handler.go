package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"photoset/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	UploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{UploadDir: uploadDir}
}

// UploadImage 上传单张图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的图片")
		return
	}
	defer file.Close()

	// 校验文件大小（10MB）
	if header.Size > 10*1024*1024 {
		response.Error(c, http.StatusBadRequest, "图片大小不能超过10MB")
		return
	}

	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "仅支持 JPG、PNG、WebP 格式")
		return
	}

	// 生成存储路径：/uploads/images/2026/04/08/uuid.jpg
	now := time.Now()
	dirPath := filepath.Join(h.UploadDir, "images", now.Format("2006"), now.Format("01"), now.Format("02"))
	os.MkdirAll(dirPath, 0755)

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(dirPath, filename)

	// 保存文件
	dst, err := os.Create(fullPath)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "保存图片失败")
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存图片失败")
		return
	}

	// 返回可访问的 URL
	relPath := fmt.Sprintf("/uploads/images/%s/%s/%s/%s",
		now.Format("2006"), now.Format("01"), now.Format("02"), filename)

	response.Success(c, gin.H{"url": relPath})
}
