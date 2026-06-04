package middleware

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// RouteMetrics 路由级别性能指标
type RouteMetrics struct {
	mu         sync.RWMutex
	path       string
	count      int64         // 总请求数
	totalTime  time.Duration // 总耗时
	minTime    time.Duration // 最小耗时
	maxTime    time.Duration // 最大耗时
	status2xx  int64         // 2xx 响应数
	status3xx  int64         // 3xx 响应数
	status4xx  int64         // 4xx 响应数
	status5xx  int64         // 5xx 响应数
}

// GlobalMetrics 全局性能指标
type GlobalMetrics struct {
	mu             sync.RWMutex
	routes         map[string]*RouteMetrics
	startTime      time.Time
	totalRequests  int64
	totalErrors5xx int64
	activeRequests int64 // 当前活跃请求数
}

var globalMetrics = &GlobalMetrics{
	routes:    make(map[string]*RouteMetrics),
	startTime: time.Now(),
}

// RecordRequest 记录请求指标
func (gm *GlobalMetrics) RecordRequest(path, method string, statusCode int, latency time.Duration) {
	atomic.AddInt64(&gm.totalRequests, 1)

	key := method + " " + path
	gm.mu.Lock()
	metrics, exists := gm.routes[key]
	if !exists {
		metrics = &RouteMetrics{path: key}
		gm.routes[key] = metrics
	}
	gm.mu.Unlock()

	atomic.AddInt64(&metrics.count, 1)

	// 更新延迟指标
	metrics.mu.Lock()
	metrics.totalTime += latency
	if metrics.minTime == 0 || latency < metrics.minTime {
		metrics.minTime = latency
	}
	if latency > metrics.maxTime {
		metrics.maxTime = latency
	}
	metrics.mu.Unlock()

	// 更新状态码计数
	switch {
	case statusCode >= 500:
		atomic.AddInt64(&metrics.status5xx, 1)
		atomic.AddInt64(&gm.totalErrors5xx, 1)
	case statusCode >= 400:
		atomic.AddInt64(&metrics.status4xx, 1)
	case statusCode >= 300:
		atomic.AddInt64(&metrics.status3xx, 1)
	default:
		atomic.AddInt64(&metrics.status2xx, 1)
	}
}

// IncrementActive 增加活跃连接数
func (gm *GlobalMetrics) IncrementActive() {
	atomic.AddInt64(&gm.activeRequests, 1)
}

// DecrementActive 减少活跃连接数
func (gm *GlobalMetrics) DecrementActive() {
	atomic.AddInt64(&gm.activeRequests, -1)
}

// RouteMetricSnapshot 路由指标快照
type RouteMetricSnapshot struct {
	Path      string  `json:"path"`
	Count     int64   `json:"count"`
	AvgTime   float64 `json:"avg_time_ms"` // 平均响应时间（毫秒）
	MinTime   float64 `json:"min_time_ms"`
	MaxTime   float64 `json:"max_time_ms"`
	Status2xx int64   `json:"status_2xx"`
	Status3xx int64   `json:"status_3xx"`
	Status4xx int64   `json:"status_4xx"`
	Status5xx int64   `json:"status_5xx"`
}

// Snapshot 获取全局指标快照
func (gm *GlobalMetrics) Snapshot() map[string]interface{} {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	routes := make([]RouteMetricSnapshot, 0, len(gm.routes))
	for _, m := range gm.routes {
		m.mu.RLock()
		avgTime := float64(0)
		count := atomic.LoadInt64(&m.count)
		if count > 0 {
			avgTime = float64(m.totalTime.Microseconds()) / float64(count) / 1000.0
		}
		routes = append(routes, RouteMetricSnapshot{
			Path:      m.path,
			Count:     count,
			AvgTime:   avgTime,
			MinTime:   float64(m.minTime.Microseconds()) / 1000.0,
			MaxTime:   float64(m.maxTime.Microseconds()) / 1000.0,
			Status2xx: atomic.LoadInt64(&m.status2xx),
			Status4xx: atomic.LoadInt64(&m.status4xx),
			Status5xx: atomic.LoadInt64(&m.status5xx),
		})
		m.mu.RUnlock()
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"uptime_seconds":   int64(time.Since(gm.startTime).Seconds()),
		"total_requests":   atomic.LoadInt64(&gm.totalRequests),
		"total_errors_5xx": atomic.LoadInt64(&gm.totalErrors5xx),
		"active_requests":  atomic.LoadInt64(&gm.activeRequests),
		"num_goroutine":    runtime.NumGoroutine(),
		"memory": map[string]interface{}{
			"alloc_mb":            float64(memStats.Alloc) / 1024 / 1024,
			"total_alloc_mb":      float64(memStats.TotalAlloc) / 1024 / 1024,
			"sys_mb":              float64(memStats.Sys) / 1024 / 1024,
			"gc_cycles":           memStats.NumGC,
			"heap_objects":        memStats.HeapObjects,
			"stack_inuse_mb":      float64(memStats.StackInuse) / 1024 / 1024,
			"gc_pause_total_ms":   float64(memStats.PauseTotalNs) / 1e6,
			"last_gc_pause_ms":    float64(memStats.PauseNs[(memStats.NumGC+255)%256]) / 1e6,
			"last_gc_seconds_ago": uint64(time.Since(time.Unix(0, int64(memStats.LastGC))).Seconds()),
			"num_forced_gc":       memStats.NumForcedGC,
			"next_gc_mb":          float64(memStats.NextGC) / 1024 / 1024,
		},
		"routes": routes,
	}
}

// GetGlobalMetrics 获取全局指标实例
func GetGlobalMetrics() *GlobalMetrics {
	return globalMetrics
}

// Metrics 性能指标中间件 — 记录每个请求的耗时和状态码
func Metrics() func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		method := c.Request.Method

		if path == "" {
			path = c.Request.URL.Path
		}

		globalMetrics.IncrementActive()

		c.Next()

		globalMetrics.DecrementActive()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		globalMetrics.RecordRequest(path, method, statusCode, latency)
	}
}
