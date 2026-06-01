package admin

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"photoset/internal/database"
	"photoset/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

type SystemStatus struct {
	Server   ServerInfo   `json:"server"`
	Memory   MemoryInfo   `json:"memory"`
	Database DatabaseInfo `json:"database"`
	Redis    RedisInfo    `json:"redis"`
	Disk     DiskInfo     `json:"disk"`
}

type ServerInfo struct {
	Uptime     string `json:"uptime"`
	GoVersion  string `json:"go_version"`
	Goroutines int    `json:"goroutines"`
	OS         string `json:"os"`
	Hostname   string `json:"hostname"`
}

type MemoryInfo struct {
	AllocMB      uint64 `json:"alloc_mb"`
	TotalAllocMB uint64 `json:"total_alloc_mb"`
	SysMB        uint64 `json:"sys_mb"`
	NumGC        uint32 `json:"gc_count"`
}

type DatabaseInfo struct {
	Status          string `json:"status"`
	OpenConnections int    `json:"open_connections"`
	InUse           int    `json:"in_use"`
	Idle            int    `json:"idle"`
	WaitCount       int64  `json:"wait_count"`
	WaitDuration    string `json:"wait_duration"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
}

type RedisInfo struct {
	Status           string  `json:"status"`
	UsedMemoryMB     float64 `json:"used_memory_mb"`
	ConnectedClients int64   `json:"connected_clients"`
	UptimeInSeconds   int64   `json:"uptime_seconds"`
}

type DiskInfo struct {
	UploadsSizeMB float64 `json:"uploads_size_mb"`
	UploadsCount  int     `json:"uploads_count"`
}

var startTime = time.Now()

// GetSystemStatus 获取系统状态
func (h *SystemHandler) GetSystemStatus(c *gin.Context) {
	status := SystemStatus{
		Server:   h.getServerInfo(),
		Memory:   h.getMemoryInfo(),
		Database: h.getDatabaseInfo(),
		Redis:    h.getRedisInfo(),
		Disk:     h.getDiskInfo(),
	}

	response.Success(c, status)
}

// HealthCheck 健康检查
func (h *SystemHandler) HealthCheck(c *gin.Context) {
	dbStatus := "connected"
	db := database.GetMySQL()
	if db == nil {
		dbStatus = "disconnected"
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			dbStatus = "error"
		} else if err := sqlDB.Ping(); err != nil {
			dbStatus = "error"
		}
	}

	redisStatus := "disconnected"
	if database.IsRedisAvailable() && database.RedisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := database.RedisClient.Ping(ctx).Err(); err == nil {
			redisStatus = "connected"
		} else {
			redisStatus = "error"
		}
	}

	status := "healthy"
	httpStatus := http.StatusOK
	if dbStatus != "connected" {
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status":    status,
		"database":  dbStatus,
		"redis":     redisStatus,
		"timestamp": time.Now().Unix(),
	})
}

func (h *SystemHandler) getServerInfo() ServerInfo {
	hostname, _ := os.Hostname()
	return ServerInfo{
		Uptime:     formatDuration(time.Since(startTime)),
		GoVersion:  runtime.Version(),
		Goroutines: runtime.NumGoroutine(),
		OS:         runtime.GOOS + "/" + runtime.GOARCH,
		Hostname:   hostname,
	}
}

func (h *SystemHandler) getMemoryInfo() MemoryInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return MemoryInfo{
		AllocMB:      bToMB(m.Alloc),
		TotalAllocMB: bToMB(m.TotalAlloc),
		SysMB:        bToMB(m.Sys),
		NumGC:        m.NumGC,
	}
}

func (h *SystemHandler) getDatabaseInfo() DatabaseInfo {
	info := DatabaseInfo{
		Status: "disconnected",
	}

	db := database.GetMySQL()
	if db == nil {
		return info
	}

	sqlDB, err := db.DB()
	if err != nil {
		info.Status = "error"
		return info
	}

	if err := sqlDB.Ping(); err != nil {
		info.Status = "error"
		return info
	}

	stats := sqlDB.Stats()
	info.Status = "connected"
	info.OpenConnections = stats.OpenConnections
	info.InUse = stats.InUse
	info.Idle = stats.Idle
	info.WaitCount = stats.WaitCount
	info.WaitDuration = stats.WaitDuration.String()
	info.MaxOpenConns = stats.MaxOpenConnections
	info.MaxIdleConns = 0 // sql.DBStats 不暴露 MaxIdleConns，暂设为0

	return info
}

func (h *SystemHandler) getRedisInfo() RedisInfo {
	info := RedisInfo{
		Status: "disconnected",
	}

	if !database.IsRedisAvailable() || database.RedisClient == nil {
		return info
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := database.RedisClient.Ping(ctx).Err(); err != nil {
		info.Status = "error"
		return info
	}

	info.Status = "connected"

	// 获取 Redis 信息
	redisInfo, err := database.RedisClient.Info(ctx, "memory", "clients", "server").Result()
	if err == nil {
		info.parseRedisInfo(redisInfo)
	}

	return info
}

func (info *RedisInfo) parseRedisInfo(redisInfo string) {
	lines := strings.Split(redisInfo, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "used_memory:") {
			memStr := strings.TrimPrefix(line, "used_memory:")
			if bytes := parseUint64(memStr); bytes > 0 {
				info.UsedMemoryMB = float64(bytes) / 1024 / 1024
			}
		} else if strings.HasPrefix(line, "connected_clients:") {
			clientsStr := strings.TrimPrefix(line, "connected_clients:")
			info.ConnectedClients = parseInt64(clientsStr)
		} else if strings.HasPrefix(line, "uptime_in_seconds:") {
			uptimeStr := strings.TrimPrefix(line, "uptime_in_seconds:")
			info.UptimeInSeconds = parseInt64(uptimeStr)
		}
	}
}

func (h *SystemHandler) getDiskInfo() DiskInfo {
	info := DiskInfo{}

	// 获取 uploads 目录大小
	uploadsPath := "./uploads"
	if _, err := os.Stat(uploadsPath); err == nil {
		size, count := dirSize(uploadsPath)
		info.UploadsSizeMB = float64(size) / 1024 / 1024
		info.UploadsCount = count
	}

	return info
}

func dirSize(path string) (int64, int) {
	var size int64
	var count int
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		return nil
	})
	return size, count
}

func bToMB(b uint64) uint64 {
	return b / 1024 / 1024
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	}
	return fmt.Sprintf("%d分钟", minutes)
}

func parseUint64(s string) uint64 {
	var n uint64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + uint64(c-'0')
		} else {
			break
		}
	}
	return n
}

func parseInt64(s string) int64 {
	var n int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		} else {
			break
		}
	}
	return n
}
