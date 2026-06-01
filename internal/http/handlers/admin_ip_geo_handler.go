package handlers

import (
	"net/http"
	"strconv"

	"photoset/internal/pkg/response"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminIPGeoHandler IP地理位置后台管理处理器
type AdminIPGeoHandler struct {
	ipGeoService *service.IPGeoService
}

// NewAdminIPGeoHandler 创建IP地理位置后台管理处理器
func NewAdminIPGeoHandler() *AdminIPGeoHandler {
	return &AdminIPGeoHandler{
		ipGeoService: service.NewIPGeoService(),
	}
}

// GetConfig 获取IP地理位置配置
func (h *AdminIPGeoHandler) GetConfig(c *gin.Context) {
	config := h.ipGeoService.GetConfig()
	dbInfo := h.ipGeoService.GetDatabaseInfo()

	response.Success(c, gin.H{
		"config":   config,
		"database": dbInfo,
	})
}

// UpdateConfig 更新IP地理位置配置
func (h *AdminIPGeoHandler) UpdateConfig(c *gin.Context) {
	var req struct {
		Enabled        *bool   `json:"enabled"`
		Mirror         *string `json:"mirror"`
		DownloadURLV4  *string `json:"download_url_v4"`
		DownloadURLV6  *string `json:"download_url_v6"`
		UpdateDays     *int    `json:"update_days"`
		DatabasePathV4 *string `json:"database_path_v4"`
		DatabasePathV6 *string `json:"database_path_v6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取当前配置
	config := h.ipGeoService.GetConfig()

	// 更新配置
	if req.Enabled != nil {
		config.Enabled = *req.Enabled
	}
	if req.Mirror != nil {
		config.Mirror = *req.Mirror
	}
	if req.DownloadURLV4 != nil {
		config.DownloadURLV4 = *req.DownloadURLV4
	}
	if req.DownloadURLV6 != nil {
		config.DownloadURLV6 = *req.DownloadURLV6
	}
	if req.UpdateDays != nil {
		config.UpdateDays = *req.UpdateDays
	}
	if req.DatabasePathV4 != nil {
		config.DatabasePathV4 = *req.DatabasePathV4
	}
	if req.DatabasePathV6 != nil {
		config.DatabasePathV6 = *req.DatabasePathV6
	}

	// 保存配置
	if err := h.ipGeoService.UpdateConfig(config); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存配置失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "配置更新成功"})
}

// UpdateDatabase 手动更新数据库
func (h *AdminIPGeoHandler) UpdateDatabase(c *gin.Context) {
	if err := h.ipGeoService.UpdateDatabase(); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新数据库失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "数据库更新成功"})
}

// GetDatabaseInfo 获取数据库信息
func (h *AdminIPGeoHandler) GetDatabaseInfo(c *gin.Context) {
	info := h.ipGeoService.GetDatabaseInfo()
	response.Success(c, info)
}

// TestIP 查询测试IP
func (h *AdminIPGeoHandler) TestIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		response.Error(c, http.StatusBadRequest, "请提供IP地址")
		return
	}

	// 查询IP地理位置
	location := h.ipGeoService.GetLocation(ip)
	fullLocation := h.ipGeoService.GetFullLocation(ip)

	response.Success(c, gin.H{
		"ip":            ip,
		"location":      location,
		"full_location": fullLocation,
	})
}

// GetUpdateLogs 获取更新日志（从站点设置中获取）
func (h *AdminIPGeoHandler) GetUpdateLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 这里可以扩展为从数据库获取更新日志
	// 目前返回简单的配置信息
	config := h.ipGeoService.GetConfig()
	dbInfo := h.ipGeoService.GetDatabaseInfo()

	response.Success(c, gin.H{
		"list": []gin.H{
			{
				"config":   config,
				"database": dbInfo,
			},
		},
		"total":     1,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetStatus 获取IP地理位置服务状态
func (h *AdminIPGeoHandler) GetStatus(c *gin.Context) {
	config := h.ipGeoService.GetConfig()
	dbInfo := h.ipGeoService.GetDatabaseInfo()

	status := "disabled"
	if config.Enabled {
		status = "enabled"
	}

	// 检查数据库文件是否存在
	dbStatus := "missing"
	if exists, ok := dbInfo["file_exists_v4"].(bool); ok && exists {
		dbStatus = "exists"
	}

	response.Success(c, gin.H{
		"status":           status,
		"database":         dbStatus,
		"last_update":      config.LastUpdate,
		"update_days":      config.UpdateDays,
		"download_url_v4":  config.DownloadURLV4,
		"download_url_v6":  config.DownloadURLV6,
		"file_size_v4":     dbInfo["file_size_v4"],
		"file_size_v6":     dbInfo["file_size_v6"],
		"ipv4_loaded":      dbInfo["ipv4_loaded"],
		"ipv6_loaded":      dbInfo["ipv6_loaded"],
	})
}
