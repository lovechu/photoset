package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"

	"photoset/internal/pkg/response"
	"photoset/internal/service"
	"photoset/internal/storage"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storage   storage.Storage
	watermark *service.WatermarkService
}

func NewUploadHandler(s storage.Storage) *UploadHandler {
	return &UploadHandler{
		storage:   s,
		watermark: service.InitWatermarkService(),
	}
}

// allowedMIMETypes 允许的图片 MIME 类型（基于 magic bytes）
var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
	"image/avif": true,
}

// bytesFile 将 []byte 包装为 multipart.File 接口
type bytesFile struct {
	*bytes.Reader
}

func (bytesFile) Close() error { return nil }

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

	// 扩展名白名单（第一道防线）
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true, ".avif": true}
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "仅支持 JPG、PNG、WebP、GIF、AVIF 格式")
		return
	}

	// 读取文件内容以验证 magic bytes（第二道防线，防止伪装扩展名）
	buf, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "读取文件失败")
		return
	}
	mtype := mimetype.Detect(buf)
	if !allowedMIMETypes[mtype.String()] {
		response.Error(c, http.StatusBadRequest, "文件类型不合法，仅支持真实图片文件")
		return
	}

	// 解析 type 参数：cover（封面图）/ photo（付费照片，默认）
	uploadType := storage.UploadTypePhoto
	if c.PostForm("type") == "cover" {
		uploadType = storage.UploadTypeCover
	}

	// 添加水印（仅对付费照片添加水印，封面图不加水印）
	if h.watermark != nil && uploadType == storage.UploadTypePhoto {
		watermarked, err := h.watermark.AddWatermark(buf, mtype.String())
		if err == nil && len(watermarked) > 0 {
			buf = watermarked
			// 重新检测水印后的 MIME 类型（通常是 PNG）
			wmMtype := mimetype.Detect(buf)
			if allowedMIMETypes[wmMtype.String()] {
				mtype = wmMtype
			}
		}
	}

	// 解析 photoset_id：嵌入路径用于签名校验（photo 类型必须，cover 类型可选）
	var photosetID uint
	if pidStr := c.PostForm("photoset_id"); pidStr != "" {
		if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
			photosetID = uint(pid)
		}
	}

	// 将 buf 重新包装为 multipart.File，并更新 Content-Type
	wrappedFile := bytesFile{bytes.NewReader(buf)}
	header.Header = make(textproto.MIMEHeader)
	header.Header.Set("Content-Type", mtype.String())
	header.Header.Set("Content-Disposition",
		`form-data; name="image"; filename="`+header.Filename+`"`)

	url, err := h.storage.UploadWithType(wrappedFile, header, uploadType, photosetID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "上传失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"url": url})
}
