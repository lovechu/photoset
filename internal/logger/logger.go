package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Log levels
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// Logger 全局日志实例
var Logger *slog.Logger

// Init 初始化日志系统
// level: debug/info/warn/error
// format: json/text
func Init(level, format string) {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case LevelDebug:
		logLevel = slog.LevelDebug
	case LevelInfo:
		logLevel = slog.LevelInfo
	case LevelWarn:
		logLevel = slog.LevelWarn
	case LevelError:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
		AddSource: logLevel == slog.LevelDebug, // DEBUG 级别显示源码位置
	}

	var handler slog.Handler
	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// Debug 调试日志
func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

// Info 信息日志
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

// Warn 警告日志
func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

// Error 错误日志
func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

// With 创建带字段的子日志器
func With(args ...any) *slog.Logger {
	return Logger.With(args...)
}

// Module 为日志添加模块标签
func Module(name string) *slog.Logger {
	return Logger.With("module", name)
}
