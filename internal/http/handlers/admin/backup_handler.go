package admin

import (
	"net/http"
	"strconv"

	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	backupService *service.BackupService
}

func NewBackupHandler(backupService *service.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// CreateBackup 创建备份
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	backup, err := h.backupService.CreateBackup()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建备份失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":  "备份创建成功",
		"backup":   backup,
	})
}

// ListBackups 获取备份列表
func (h *BackupHandler) ListBackups(c *gin.Context) {
	backups, err := h.backupService.ListBackups()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取备份列表失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  backups,
		"total": len(backups),
	})
}

// DownloadBackup 下载备份文件
func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.Error(c, http.StatusBadRequest, "文件名不能为空")
		return
	}

	path, err := h.backupService.GetBackupPath(filename)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(path)
}

// DeleteBackup 删除备份文件
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.Error(c, http.StatusBadRequest, "文件名不能为空")
		return
	}

	if err := h.backupService.DeleteBackup(filename); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除备份失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "备份文件已删除",
	})
}

// GetBackupStats 获取备份统计信息
func (h *BackupHandler) GetBackupStats(c *gin.Context) {
	backups, err := h.backupService.ListBackups()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取备份统计失败: "+err.Error())
		return
	}

	var totalSize int64
	for _, b := range backups {
		totalSize += b.Size
	}

	response.Success(c, gin.H{
		"count":      len(backups),
		"total_size": totalSize,
		"total_size_str": formatFileSize(totalSize),
		"latest":     nil,
	})

	if len(backups) > 0 {
		response.Success(c, gin.H{
			"count":          len(backups),
			"total_size":     totalSize,
			"total_size_str": formatFileSize(totalSize),
			"latest":         backups[0],
		})
	}
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return strconv.FormatFloat(float64(bytes)/float64(div), 'f', 1, 64) + " " + string("KMGTPE"[exp]) + "B"
}
