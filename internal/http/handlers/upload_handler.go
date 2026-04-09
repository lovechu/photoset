package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"photoset/internal/pkg/response"
	"photoset/internal/storage"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storage storage.Storage
}

func NewUploadHandler(s storage.Storage) *UploadHandler {
	return &UploadHandler{storage: s}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的图片")
		return
	}
	defer file.Close()

	if header.Size > 10*1024*1024 {
		response.Error(c, http.StatusBadRequest, "图片大小不能超过10MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "仅支持 JPG、PNG、WebP 格式")
		return
	}

	url, err := h.storage.Upload(file, header)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "上传失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"url": url})
}
