package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"photoset/internal/database"
	"photoset/internal/http/middleware"
	"photoset/internal/logger"
	"photoset/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Feedback 用户反馈数据模型
type Feedback struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"index"`
	Category   string    `json:"category" gorm:"size:50"`
	Content    string    `json:"content" gorm:"type:text"`
	Images     string    `json:"images" gorm:"type:text"` // JSON 数组字符串
	DeviceInfo string    `json:"device_info" gorm:"size:500"`
	Contact    string    `json:"contact" gorm:"size:200"`
	Status     string    `json:"status" gorm:"size:20;default:'pending'"` // pending, processing, resolved
	Reply      string    `json:"reply" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}

// FeedbackHandler 反馈处理器
type FeedbackHandler struct {
	db *gorm.DB
}

// NewFeedbackHandler 创建反馈处理器
func NewFeedbackHandler() *FeedbackHandler {
	db := database.GetMySQL()
	// 自动迁移
	db.AutoMigrate(&Feedback{})
	return &FeedbackHandler{db: db}
}

// CreateFeedbackRequest 创建反馈请求
type CreateFeedbackRequest struct {
	Category   string   `json:"category" binding:"required"`
	Content    string   `json:"content" binding:"required,min=10,max=2000"`
	Images     []string `json:"images"`
	DeviceInfo string   `json:"device_info"`
	Contact    string   `json:"contact"`
}

// @Summary      提交反馈
// @Description  用户提交反馈意见，支持上传图片、设备信息等
// @Tags         Feedback
// @Accept       json
// @Produce      json
// @Param        body body CreateFeedbackRequest true "反馈请求参数"
// @Success      200 {object} response.Response{data=object} "反馈提交成功"
// @Failure      400 {object} response.Response "参数错误"
// @Failure      500 {object} response.Response "提交失败"
// @Router       /api/feedback [post]
// Create 创建反馈
func (h *FeedbackHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 限制图片数量
	if len(req.Images) > 5 {
		response.Error(c, http.StatusBadRequest, "最多上传5张图片")
		return
	}

	// 将图片列表转为 JSON 字符串
	imagesJSON := "[]"
	if len(req.Images) > 0 {
		var quoted []string
		for _, img := range req.Images {
			quoted = append(quoted, fmt.Sprintf(`"%s"`, img))
		}
		imagesJSON = "[" + strings.Join(quoted, ",") + "]"
	}

	feedback := Feedback{
		UserID:     userID,
		Category:   req.Category,
		Content:    req.Content,
		Images:     imagesJSON,
		DeviceInfo: req.DeviceInfo,
		Contact:    req.Contact,
		Status:     "pending",
	}

	if err := h.db.Create(&feedback).Error; err != nil {
		logger.Error("创建反馈失败", "error", err)
		response.Error(c, http.StatusInternalServerError, "提交反馈失败")
		return
	}

	logger.Info("用户反馈已提交", "user_id", userID, "category", req.Category)
	response.Success(c, gin.H{"id": feedback.ID, "message": "感谢您的反馈，我们会尽快处理"})
}

// @Summary      上传反馈图片
// @Description  上传反馈附带的图片，支持JPG、PNG、WebP格式，最大5MB
// @Tags         Feedback
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "反馈图片文件"
// @Success      200 {object} response.Response{data=object} "上传成功，返回图片URL"
// @Failure      400 {object} response.Response "文件错误或格式不支持"
// @Failure      500 {object} response.Response "上传失败"
// @Router       /api/feedback/image [post]
// UploadFeedbackImage 上传反馈图片
func (h *FeedbackHandler) UploadFeedbackImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的图片")
		return
	}
	defer file.Close()

	// 大小校验（5MB）
	if header.Size > 5*1024*1024 {
		response.Error(c, http.StatusBadRequest, "图片大小不能超过5MB")
		return
	}

	// 扩展名校验
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "仅支持 JPG、PNG、WebP 格式")
		return
	}

	// 读取文件内容
	buf := make([]byte, header.Size)
	if _, err := file.Read(buf); err != nil {
		response.Error(c, http.StatusInternalServerError, "读取文件失败")
		return
	}

	// 存储到本地 uploads/feedback/ 目录
	filename := fmt.Sprintf("feedback_%d_%s", time.Now().UnixNano(), header.Filename)
	savePath := fmt.Sprintf("./uploads/feedback/%s", filename)

	// 确保目录存在
	if err := ensureDir("./uploads/feedback/"); err != nil {
		logger.Error("创建反馈图片目录失败", "error", err)
		response.Error(c, http.StatusInternalServerError, "上传失败")
		return
	}

	if err := writeFile(savePath, buf); err != nil {
		logger.Error("保存反馈图片失败", "error", err)
		response.Error(c, http.StatusInternalServerError, "上传失败")
		return
	}

	url := fmt.Sprintf("/uploads/feedback/%s", filename)
	response.Success(c, gin.H{"url": url})
}

// ensureDir 确保目录存在
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// writeFile 写入文件
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
