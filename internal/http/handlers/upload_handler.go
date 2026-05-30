package handlers

import (
	"bytes"
	"fmt"
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
	fmt.Printf("[UPLOAD-DEBUG] ====== 开始上传 ======\n")

	// Step 1: 获取文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		fmt.Printf("[UPLOAD-DEBUG] Step1 失败: 获取文件出错: %v\n", err)
		response.Error(c, http.StatusBadRequest, "请选择要上传的图片")
		return
	}
	defer file.Close()
	fmt.Printf("[UPLOAD-DEBUG] Step1 OK: 文件名=%s, 大小=%d bytes\n", header.Filename, header.Size)

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
	fmt.Printf("[UPLOAD-DEBUG] Step3 OK: 扩展名=%s\n", ext)

	// Step 4: 读取文件内容并验证 MIME 类型
	buf, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("[UPLOAD-DEBUG] Step4 失败: 读取文件出错: %v\n", err)
		response.Error(c, http.StatusInternalServerError, "读取文件失败")
		return
	}
	mtype := mimetype.Detect(buf)
	fmt.Printf("[UPLOAD-DEBUG] Step4 OK: 原始MIME=%s, 缓冲区大小=%d bytes\n", mtype.String(), len(buf))
	if !allowedMIMETypes[mtype.String()] {
		fmt.Printf("[UPLOAD-DEBUG] Step4 失败: MIME类型不允许: %s\n", mtype.String())
		response.Error(c, http.StatusBadRequest, "文件类型不合法，仅支持真实图片文件")
		return
	}

	// Step 5: 解析 type 参数
	uploadType := storage.UploadTypePhoto
	typeParam := c.PostForm("type")
	if typeParam == "cover" {
		uploadType = storage.UploadTypeCover
	}
	fmt.Printf("[UPLOAD-DEBUG] Step5 OK: uploadType=%s, type参数=%s\n", uploadType, typeParam)

	// Step 6: 添加水印
	if h.watermark != nil && uploadType == storage.UploadTypePhoto {
		fmt.Printf("[UPLOAD-DEBUG] Step6: 水印服务非nil, uploadType=photo, 开始添加水印...\n")
		watermarked, wmErr := h.watermark.AddWatermark(buf, mtype.String())
		fmt.Printf("[UPLOAD-DEBUG] Step6: 水印处理完成, 错误=%v, 原大小=%d, 新大小=%d\n", wmErr, len(buf), len(watermarked))
		if wmErr != nil {
			fmt.Printf("[UPLOAD-DEBUG] Step6: 水印处理出错(将使用原图): %v\n", wmErr)
		}
		if wmErr == nil && len(watermarked) > 0 {
			buf = watermarked
			wmMtype := mimetype.Detect(buf)
			fmt.Printf("[UPLOAD-DEBUG] Step6: 水印后MIME=%s\n", wmMtype.String())
			if allowedMIMETypes[wmMtype.String()] {
				mtype = wmMtype
			}
		}
	} else {
		fmt.Printf("[UPLOAD-DEBUG] Step6: 跳过水印 (watermark=%v, uploadType=%s)\n", h.watermark != nil, uploadType)
	}
	fmt.Printf("[UPLOAD-DEBUG] Step6 OK: 最终MIME=%s, 缓冲区大小=%d bytes\n", mtype.String(), len(buf))

	// Step 7: 解析 photoset_id
	var photosetID uint
	pidStr := c.PostForm("photoset_id")
	if pidStr != "" {
		if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
			photosetID = uint(pid)
		}
	}
	fmt.Printf("[UPLOAD-DEBUG] Step7 OK: photosetID=%d, photoset_id参数=%s\n", photosetID, pidStr)

	// Step 8: 包装文件并设置 header
	wrappedFile := bytesFile{bytes.NewReader(buf)}
	header.Header = make(textproto.MIMEHeader)
	header.Header.Set("Content-Type", mtype.String())
	header.Header.Set("Content-Disposition",
		`form-data; name="image"; filename="`+header.Filename+`"`)
	fmt.Printf("[UPLOAD-DEBUG] Step8 OK: 文件已包装, Content-Type=%s\n", mtype.String())

	// Step 9: 上传到存储
	fmt.Printf("[UPLOAD-DEBUG] Step9: 开始上传到存储...\n")
	url, err := h.storage.UploadWithType(wrappedFile, header, uploadType, photosetID)
	if err != nil {
		fmt.Printf("[UPLOAD-DEBUG] Step9 失败: 存储上传出错: %v\n", err)
		response.Error(c, http.StatusInternalServerError, "上传失败: "+err.Error())
		return
	}

	fmt.Printf("[UPLOAD-DEBUG] Step9 OK: 上传成功! URL=%s\n", url)
	fmt.Printf("[UPLOAD-DEBUG] ====== 上传完成 ======\n")
	response.Success(c, gin.H{"url": url})
}
