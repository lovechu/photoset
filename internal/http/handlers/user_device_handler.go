package handlers

import (
	"net/http"
	"photoset/internal/http/middleware"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// UserDeviceHandler handles user device HTTP requests
type UserDeviceHandler struct {
	deviceService *service.UserDeviceService
}

// NewUserDeviceHandler creates a new UserDeviceHandler
func NewUserDeviceHandler(deviceService *service.UserDeviceService) *UserDeviceHandler {
	return &UserDeviceHandler{deviceService: deviceService}
}

// @Summary      设备列表
// @Description  获取当前用户的活跃设备列表
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  object  "设备列表"
// @Failure      401  {object}  object  "未登录"
// @Security     BearerAuth
// @Router       /api/user/devices [get]
// GetUserDevices handles getting user's active devices
func (h *UserDeviceHandler) GetUserDevices(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	devices, err := h.deviceService.GetUserDevices(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取设备列表失败"})
		return
	}

	// Format response
	type DeviceItem struct {
		ID           uint   `json:"id"`
		DeviceID     string `json:"device_id"`
		DeviceName   string `json:"device_name"`
		DeviceType   string `json:"device_type"`
		OS           string `json:"os"`
		Browser      string `json:"browser"`
		IP           string `json:"ip"`
		IPLocation   string `json:"ip_location"`
		LastActiveAt string `json:"last_active_at"`
		IsActive     bool   `json:"is_active"`
	}

	items := make([]DeviceItem, 0, len(devices))
	for _, device := range devices {
		item := DeviceItem{
			ID:           device.ID,
			DeviceID:     device.DeviceID,
			DeviceName:   device.DeviceName,
			DeviceType:   device.DeviceType,
			OS:           device.OS,
			Browser:      device.Browser,
			IP:           device.IP,
			IPLocation:   device.IPLocation,
			LastActiveAt: device.LastActiveAt.Format("2006-01-02 15:04:05"),
			IsActive:     device.IsActive,
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": items,
	})
}

// @Summary      停用设备
// @Description  停用指定设备
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        deviceId  path  string  true  "设备ID"
// @Success      200  {object}  object  "操作成功"
// @Failure      400  {object}  object  "设备ID为空"
// @Security     BearerAuth
// @Router       /api/user/devices/{deviceId} [delete]
// DeactivateDevice handles deactivating a specific device
func (h *UserDeviceHandler) DeactivateDevice(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	if err := h.deviceService.DeactivateDevice(userID, deviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "停用设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设备已停用",
	})
}

// @Summary      停用其他设备
// @Description  停用除当前设备外的所有其他设备
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        X-Device-ID  header  string  false  "当前设备ID"
// @Success      200  {object}  object  "操作成功"
// @Failure      400  {object}  object  "当前设备ID为空"
// @Security     BearerAuth
// @Router       /api/user/devices [delete]
// DeactivateOtherDevices handles deactivating all devices except current one
func (h *UserDeviceHandler) DeactivateOtherDevices(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	// Get current device ID from request header or query parameter
	currentDeviceID := c.GetHeader("X-Device-ID")
	if currentDeviceID == "" {
		currentDeviceID = c.Query("current_device_id")
	}
	
	if currentDeviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前设备ID不能为空"})
		return
	}

	if err := h.deviceService.DeactivateOtherDevices(userID, currentDeviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已停用其他设备",
	})
}