package handlers

import (
	"net/http"
	"runtime"
	"runtime/pprof"
	"time"

	"photoset/internal/http/middleware"
	"photoset/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type PerformanceHandler struct{}

func NewPerformanceHandler() *PerformanceHandler {
	return &PerformanceHandler{}
}

// GetMetrics 获取性能指标
// @Summary      性能指标
// @Description  获取系统性能指标（请求量、延迟、内存、goroutine 等）
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "性能指标"
// @Router       /api/system/metrics [get]
func (h *PerformanceHandler) GetMetrics(c *gin.Context) {
	gm := middleware.GetGlobalMetrics()
	snapshot := gm.Snapshot()
	response.Success(c, snapshot)
}

// GetGoroutines 获取 Goroutine 详细信息（用于排查泄露）
// @Summary      Goroutine 信息
// @Description  获取当前所有 Goroutine 的堆栈信息（调试用）
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "Goroutine 堆栈"
// @Security     BearerAuth
// @Router       /api/admin/system/goroutines [get]
func (h *PerformanceHandler) GetGoroutines(c *gin.Context) {
	count := runtime.NumGoroutine()
	var buf []byte
	n := 1024 * 1024 // 1MB buffer
	for {
		buf = make([]byte, n)
		n = runtime.Stack(buf, true)
		if n < len(buf) {
			break
		}
		n *= 2
	}

	response.Success(c, gin.H{
		"goroutine_count": count,
		"stack_trace":     string(buf[:n]),
	})
}

// StartCPUProfile 开始 CPU 分析
// @Summary      开始 CPU Profile
// @Description  开始 CPU 性能分析采样
// @Tags         System
// @Accept       json
// @Produce      json
// @Param        duration  query  int  false  "采样时长（秒）"  default(30)
// @Success      200  {object}  response.Response  "Profile 文件"
// @Security     BearerAuth
// @Router       /api/admin/system/profile/cpu [get]
func (h *PerformanceHandler) StartCPUProfile(c *gin.Context) {
	duration := 30 // 默认30秒
	var dur time.Duration

	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=cpu.prof")
	c.Writer.WriteHeader(http.StatusOK)

	if d := c.Query("duration"); d != "" {
		f, _ := time.ParseDuration(d + "s")
		if f > 0 && f <= 60*time.Second {
			dur = f
		}
	}
	if dur == 0 {
		dur = time.Duration(duration) * time.Second
	}

	pprof.StartCPUProfile(c.Writer)
	time.Sleep(dur)
	pprof.StopCPUProfile()
}

// GetHeapProfile 获取内存快照
// @Summary      内存分析
// @Description  获取当前堆内存分析快照（pprof 格式）
// @Tags         System
// @Accept       json
// @Produce      octet-stream
// @Success      200  {file}  application/octet-stream  "Heap Profile"
// @Security     BearerAuth
// @Router       /api/admin/system/profile/heap [get]
func (h *PerformanceHandler) GetHeapProfile(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=heap.prof")
	c.Writer.WriteHeader(http.StatusOK)

	runtime.GC() // 先触发 GC 得到更准确的内存使用情况
	pprof.WriteHeapProfile(c.Writer)
}
