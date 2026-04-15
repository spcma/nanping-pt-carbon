package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Factory 语法糖日志接口设计（Sugar Logger）
type Factory interface {
	// 链式调用 API
	With(fields ...zap.Field) Factory
	WithCategory(category string) Factory
	WithTraceID(traceID string) Factory
	WithRequest(c *gin.Context) Factory
	WithUser(userID int64, username string) Factory

	// 标准化的日志方法
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)

	// 业务特定方法
	RequestReceived(method, path string, fields ...zap.Field)
	RequestCompleted(statusCode int, latency time.Duration, fields ...zap.Field)
	DatabaseOperation(operation, table string, duration time.Duration, fields ...zap.Field)
	ExternalCall(service, endpoint string, duration time.Duration, fields ...zap.Field)
	BusinessEvent(eventType string, fields ...zap.Field)
}

// SugarLogger 语法糖日志器实现
type SugarLogger struct {
	baseLogger *zap.Logger
	fields     []zap.Field
	category   string
}

// NewSugarLogger 创建新的语法糖日志器
func NewSugarLogger(category string) Factory {
	baseLogger := Get(category)
	if baseLogger == nil {
		// 如果 globalLogger 未初始化或获取失败，使用降级处理
		baseLogger = zap.NewNop() // 使用 nop logger 避免 panic
	}

	return &SugarLogger{
		baseLogger: baseLogger,
		category:   category,
		fields:     make([]zap.Field, 0),
	}
}

// With 添加通用字段
func (sl *SugarLogger) With(fields ...zap.Field) Factory {
	// 创建新的切片，避免并发修改原切片
	newFields := make([]zap.Field, len(sl.fields), len(sl.fields)+len(fields))
	copy(newFields, sl.fields)
	newFields = append(newFields, fields...)

	return &SugarLogger{
		baseLogger: sl.baseLogger,
		category:   sl.category,
		fields:     newFields,
	}
}

// WithCategory 设置日志分类
func (sl *SugarLogger) WithCategory(category string) Factory {
	return &SugarLogger{
		baseLogger: Get(category),
		category:   category,
		fields:     sl.fields,
	}
}

// WithTraceID 添加追踪 ID
func (sl *SugarLogger) WithTraceID(traceID string) Factory {
	return sl.With(zap.String("trace_id", traceID))
}

// WithRequest 从 Gin 上下文提取请求信息
func (sl *SugarLogger) WithRequest(c *gin.Context) Factory {
	if c == nil {
		return sl
	}

	fields := []zap.Field{
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	}

	// 尝试获取追踪 ID
	if traceID, exists := c.Get("trace_id"); exists {
		if tid, ok := traceID.(string); ok {
			fields = append(fields, zap.String("trace_id", tid))
		} else if tid, ok := traceID.(int64); ok {
			fields = append(fields, zap.Int64("trace_id", tid))
		}
	}

	// 尝试获取用户信息
	if userID, exists := c.Get("id"); exists {
		if uid, ok := userID.(int64); ok {
			fields = append(fields, zap.Int64("user_id", uid))
		}
	}
	if username, exists := c.Get("username"); exists {
		if un, ok := username.(string); ok {
			fields = append(fields, zap.String("username", un))
		}
	}

	return sl.With(fields...)
}

// WithUser 添加用户信息
func (sl *SugarLogger) WithUser(userID int64, username string) Factory {
	fields := []zap.Field{
		zap.Int64("user_id", userID),
		zap.String("username", username),
	}
	return sl.With(fields...)
}

// Debug 记录调试日志
func (sl *SugarLogger) Debug(msg string, fields ...zap.Field) {
	// 合并字段时创建新切片，避免并发问题
	allFields := make([]zap.Field, len(sl.fields)+len(fields))
	copy(allFields, sl.fields)
	copy(allFields[len(sl.fields):], fields)
	sl.baseLogger.Debug(msg, allFields...)
}

// Info 记录信息日志
func (sl *SugarLogger) Info(msg string, fields ...zap.Field) {
	// 合并字段时创建新切片，避免并发问题
	allFields := make([]zap.Field, len(sl.fields)+len(fields))
	copy(allFields, sl.fields)
	copy(allFields[len(sl.fields):], fields)
	sl.baseLogger.Info(msg, allFields...)
}

// Warn 记录警告日志
func (sl *SugarLogger) Warn(msg string, fields ...zap.Field) {
	// 合并字段时创建新切片，避免并发问题
	allFields := make([]zap.Field, len(sl.fields)+len(fields))
	copy(allFields, sl.fields)
	copy(allFields[len(sl.fields):], fields)
	sl.baseLogger.Warn(msg, allFields...)
}

// Error 记录错误日志
func (sl *SugarLogger) Error(msg string, fields ...zap.Field) {
	// 合并字段时创建新切片，避免并发问题
	allFields := make([]zap.Field, len(sl.fields)+len(fields))
	copy(allFields, sl.fields)
	copy(allFields[len(sl.fields):], fields)
	sl.baseLogger.Error(msg, allFields...)
}

// RequestReceived 记录请求接收事件
func (sl *SugarLogger) RequestReceived(method, path string, fields ...zap.Field) {
	eventFields := []zap.Field{
		zap.String("event_type", "request_received"),
		zap.String("http_method", method),
		zap.String("http_path", path),
	}
	eventFields = append(eventFields, fields...)
	sl.Info("HTTP request received", eventFields...)
}

// RequestCompleted 记录请求完成事件
func (sl *SugarLogger) RequestCompleted(statusCode int, latency time.Duration, fields ...zap.Field) {
	eventFields := []zap.Field{
		zap.String("event_type", "request_completed"),
		zap.Int("status_code", statusCode),
		zap.Duration("latency", latency),
		zap.String("latency_ms", formatDuration(latency)),
	}
	eventFields = append(eventFields, fields...)

	if statusCode >= 500 {
		sl.Error("HTTP request completed with server error", eventFields...)
	} else if statusCode >= 400 {
		sl.Warn("HTTP request completed with client error", eventFields...)
	} else {
		sl.Info("HTTP request completed successfully", eventFields...)
	}
}

// DatabaseOperation 记录数据库操作
func (sl *SugarLogger) DatabaseOperation(operation, table string, duration time.Duration, fields ...zap.Field) {
	eventFields := []zap.Field{
		zap.String("event_type", "database_operation"),
		zap.String("db_operation", operation),
		zap.String("db_table", table),
		zap.Duration("duration", duration),
		zap.String("duration_ms", formatDuration(duration)),
	}
	eventFields = append(eventFields, fields...)

	// 慢查询阈值调整为 500ms
	if duration > 500*time.Millisecond {
		sl.Warn("Slow database operation detected", eventFields...)
	} else if duration > 100*time.Millisecond {
		// 新增：记录较慢但不是特别慢的查询
		sl.Info("Database operation completed (noticeable latency)", eventFields...)
	} else {
		sl.Debug("Database operation completed", eventFields...)
	}
}

// ExternalCall 记录外部服务调用
func (sl *SugarLogger) ExternalCall(service, endpoint string, duration time.Duration, fields ...zap.Field) {
	eventFields := []zap.Field{
		zap.String("event_type", "external_call"),
		zap.String("service_name", service),
		zap.String("endpoint", endpoint),
		zap.Duration("duration", duration),
		zap.String("duration_ms", formatDuration(duration)),
	}
	eventFields = append(eventFields, fields...)

	// 外部服务调用慢查询阈值调整
	if duration > 2*time.Second {
		sl.Warn("Slow external service call detected", eventFields...)
	} else if duration > 500*time.Millisecond {
		sl.Info("External service call completed (noticeable latency)", eventFields...)
	} else {
		sl.Debug("External service call completed", eventFields...)
	}
}

// BusinessEvent 记录业务事件
func (sl *SugarLogger) BusinessEvent(eventType string, fields ...zap.Field) {
	eventFields := []zap.Field{
		zap.String("event_type", eventType),
	}
	eventFields = append(eventFields, fields...)
	sl.Info("Business event occurred", eventFields...)
}

// formatDuration 格式化持续时间为毫秒字符串
func formatDuration(d time.Duration) string {
	return d.Truncate(time.Millisecond).String()
}

// 全局便捷函数
var (
	defaultL Factory // 默认分类的日志器

	RuntimeL   Factory
	ErrorL     Factory
	DebugL     Factory
	IamL       Factory
	TrafficL   Factory
	HTTPL      Factory
	InitL      Factory
	IpfsL      Factory
	SchedulerL Factory
)

// InitGlobalLoggers 初始化全局日志器（应在应用启动时调用）
func InitGlobalLoggers() {
	if globalLogger == nil {
		// 如果 logger 还未初始化，直接返回
		return
	}
	initGlobalLoggersInternal()
}

// initGlobalLoggersInternal 内部初始化函数（由 Initialize() 调用）
func initGlobalLoggersInternal() {
	defaultL = NewSugarLogger("runtime")
	RuntimeL = NewSugarLogger("runtime")
	ErrorL = NewSugarLogger("error")
	DebugL = NewSugarLogger("debug")
	IamL = NewSugarLogger("iam")
	HTTPL = NewSugarLogger("http")
	InitL = NewSugarLogger("initializer")
	IpfsL = NewSugarLogger("ipfs")
	SchedulerL = NewSugarLogger("scheduler")
}

// WithCategory 全局 WithCategory 函数
func WithCategory(category string) Factory {
	return NewSugarLogger(category)
}

// GlobalWarn 全局 Warn 函数
func GlobalWarn(msg string, fields ...zap.Field) {
	if defaultL != nil {
		defaultL.Warn(msg, fields...)
	}
}
