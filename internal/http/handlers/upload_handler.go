package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"

	"photoset/internal/logger"
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

// allowedVideoMIMETypes 允许的视频 MIME 类型
var allowedVideoMIMETypes = map[string]bool{
	"video/mp4":  true,
	"video/quicktime": true, // .mov
	"video/x-msvideo": true, // .avi
	"video/webm": true,
	"video/x-matroska": true, // .mkv
}

// bytesFile 将 []byte 包装为 multipart.File 接口
type bytesFile struct {
	*bytes.Reader
}

func (bytesFile) Close() error { return nil }

func (h *UploadHandler) UploadImage(c *gin.Context) {
	logger.Debug("Upload: started")

	// Step 1: 获取文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		logger.Debug("Upload Step1: get file failed", "error", err)
		response.Error(c, http.StatusBadRequest, "请选择要上传的图片")
		return
	}
	defer file.Close()
	logger.Debug("Upload Step1: OK", "filename", header.Filename, "size", header.Size)

	// Step 2: 大小校验
	if header.Size > 10*1024*1024 {
		response.Error(c, http.StatusBadRequest, "图片大小不能超过10MB")
		return
	}

	// Step 3: 扩展名校验
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true, ".avif": true}
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "仅支持 JPG、PNG、WebP、GIF、AVIF 格式")
		return
	}
	logger.Debug("Upload Step3: OK", "ext", ext)

	// Step 4: 读取文件内容并验证 MIME 类型
	buf, err := io.ReadAll(file)
	if err != nil {
		logger.Debug("Upload Step4: read file failed", "error", err)
		response.Error(c, http.StatusInternalServerError, "读取文件失败")
		return
	}
	mtype := mimetype.Detect(buf)
	logger.Debug("Upload Step4: OK", "mime", mtype.String(), "size", len(buf))
	if !allowedMIMETypes[mtype.String()] {
		logger.Warn("Upload Step4: MIME type not allowed", "mime", mtype.String())
		response.Error(c, http.StatusBadRequest, "文件类型不合法，仅支持真实图片文件")
		return
	}

	// Step 5: 解析 type 参数
	uploadType := storage.UploadTypePhoto
	typeParam := c.PostForm("type")
	if typeParam == "cover" {
		uploadType = storage.UploadTypeCover
	} else if typeParam == "avatar" {
		uploadType = storage.UploadTypeAvatar
	} else if typeParam == "user_cover" {
		uploadType = storage.UploadTypeUserCover
	}
	logger.Debug("Upload Step5: OK", "uploadType", uploadType, "typeParam", typeParam)

	// Step 6: 添加水印
	if h.watermark != nil && uploadType == storage.UploadTypePhoto {
		logger.Debug("Upload Step6: adding watermark")
		watermarked, wmErr := h.watermark.AddWatermark(buf, mtype.String())
		if wmErr != nil {
			logger.Warn("Upload Step6: watermark failed, using original", "error", wmErr)
		}
		if wmErr == nil && len(watermarked) > 0 {
			buf = watermarked
			wmMtype := mimetype.Detect(buf)
			if allowedMIMETypes[wmMtype.String()] {
				mtype = wmMtype
			}
		}
	}
	logger.Debug("Upload Step6: OK", "mime", mtype.String(), "size", len(buf))

	// Step 7: 解析 photoset_id
	var photosetID uint
	pidStr := c.PostForm("photoset_id")
	if pidStr != "" {
		if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
			photosetID = uint(pid)
		}
	}
	logger.Debug("Upload Step7: OK", "photosetID", photosetID)

	// Step 8: 包装文件并设置 header
	wrappedFile := bytesFile{bytes.NewReader(buf)}
	header.Header = make(textproto.MIMEHeader)
	header.Header.Set("Content-Type", mtype.String())
	header.Header.Set("Content-Disposition",
		`form-data; name="image"; filename="`+header.Filename+`"`)
	logger.Debug("Upload Step8: OK", "contentType", mtype.String())

	// Step 9: 上传到存储
	logger.Debug("Upload Step9: uploading to storage")
	url, err := h.storage.UploadWithType(wrappedFile, header, uploadType, photosetID)
	if err != nil {
		logger.Error("Upload Step9: storage upload failed", "error", err)
		response.Error(c, http.StatusInternalServerError, "上传失败")
		return
	}

	logger.Debug("Upload Step9: upload success")
	response.Success(c, gin.H{"url": url})
}

// UploadVideo handles video file uploads
func (h *UploadHandler) UploadVideo(c *gin.Context) {
	logger.Debug("UploadVideo: started")

	// Step 1: 获取文件
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		logger.Debug("UploadVideo Step1: get file failed", "error", err)
		response.Error(c, http.StatusBadRequest, "请选择要上传的视频")
		return
	}
	defer file.Close()
	logger.Debug("UploadVideo Step1: OK", "filename", header.Filename, "size", header.Size)

	// Step 2: 大小校验（最大500MB）
	if header.Size > 500*1024*1024 {
		response.Error(c, http.StatusBadRequest, "视频大小不能超过500MB")
		return
	}

	// Step 3: 扩展名校验
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".mp4": true, ".mov": true, ".avi": true, ".webm": true, ".mkv": true}
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "仅支持 MP4、MOV、AVI、WebM、MKV 格式")
		return
	}
	logger.Debug("UploadVideo Step3: OK", "ext", ext)

	// Step 4: 读取文件内容并验证 MIME 类型
	buf, err := io.ReadAll(file)
	if err != nil {
		logger.Debug("UploadVideo Step4: read file failed", "error", err)
		response.Error(c, http.StatusInternalServerError, "读取文件失败")
		return
	}
	mtype := mimetype.Detect(buf)
	logger.Debug("UploadVideo Step4: OK", "mime", mtype.String(), "size", len(buf))
	if !allowedVideoMIMETypes[mtype.String()] {
		logger.Warn("UploadVideo Step4: MIME type not allowed", "mime", mtype.String())
		response.Error(c, http.StatusBadRequest, "文件类型不合法，仅支持真实视频文件")
		return
	}

	// Step 5: 包装文件并设置 header
	wrappedFile := bytesFile{bytes.NewReader(buf)}
	header.Header = make(textproto.MIMEHeader)
	header.Header.Set("Content-Type", mtype.String())
	header.Header.Set("Content-Disposition",
		`form-data; name="video"; filename="`+header.Filename+`"`)
	logger.Debug("UploadVideo Step5: OK", "contentType", mtype.String())

	// Step 6: 上传到存储
	logger.Debug("UploadVideo Step6: uploading to storage")
	url, err := h.storage.UploadWithType(wrappedFile, header, storage.UploadTypeVideo, 0)
	if err != nil {
		logger.Error("UploadVideo Step6: storage upload failed", "error", err)
		response.Error(c, http.StatusInternalServerError, "上传失败")
		return
	}

	logger.Debug("UploadVideo Step6: upload success")
	response.Success(c, gin.H{"url": url})
}
