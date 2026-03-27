package logger

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TraceIDHeader 请求/响应头中 traceid 的 key
const TraceIDHeader = "X-Trace-ID"

// LoggingMiddleware 日志中间件，自动处理请求上下文和追踪
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从请求头获取 traceid（支持分布式追踪）
		traceID := c.GetHeader(TraceIDHeader)
		if traceID == "" {
			// 没有则生成新的 traceid
			traceID = generateTraceID()
		}

		// 存储到 Gin 上下文
		c.Set("trace_id", traceID)

		// 同时存储到 Go context（便于跨层传递）
		ctx := WithTrace(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// 在响应头中返回 traceid（便于客户端追踪）
		c.Header(TraceIDHeader, traceID)

		// 记录请求开始时间
		startTime := time.Now()

		// 创建带上下文的日志器
		ctxLogger := NewSugarLogger("http").
			WithTraceID(traceID).
			WithRequest(c)

		// 将日志器存储到上下文中
		c.Set("logger", ctxLogger)

		// 记录请求接收
		ctxLogger.RequestReceived(c.Request.Method, c.Request.URL.Path)

		// 继续处理请求
		c.Next()

		// 计算请求耗时
		latency := time.Since(startTime)

		// 记录请求完成
		statusCode := c.Writer.Status()
		ctxLogger.RequestCompleted(statusCode, latency)

		// 如果有错误，额外记录
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				ErrorLogger.Error("Request error occurred",
					zap.String("trace_id", traceID),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Int("status_code", statusCode),
					zap.Duration("latency", latency),
					zap.String("error_message", err.Error()),
				)
			}
		}
	}
}

// GetLoggerFromContext 从 Gin 上下文获取日志器
func GetLoggerFromContext(c *gin.Context) Factory {
	if logger, exists := c.Get("logger"); exists {
		if l, ok := logger.(Factory); ok {
			return l
		}
	}
	// 返回默认日志器
	return defaultLogger
}

// generateTraceID 生成追踪 ID
func generateTraceID() string {
	return uuid.New().String()
}

// RecoveryMiddleware 带日志的恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取日志器
				logger := GetLoggerFromContext(c)

				// 记录 panic 信息
				logger.Error("Panic recovered",
					zap.Any("panic_details", err),
					zap.String("stack", zap.Stack("").String),
				)

				// 返回 500 错误
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// TimeoutMiddleware 超时中间件（带日志）
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLoggerFromContext(c)

		// 创建超时上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换请求上下文
		c.Request = c.Request.WithContext(ctx)

		// 启动定时器监控
		done := make(chan bool, 1)
		go func() {
			c.Next()
			done <- true
		}()

		select {
		case <-done:
			// 正常完成
			return
		case <-ctx.Done():
			// 超时
			logger.Warn("Request timeout",
				zap.String("trace_id", getTraceID(c)),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Duration("timeout_duration", timeout),
			)
			c.AbortWithStatus(408) // Request Timeout
		}
	}
}

// GetTraceID 从上下文获取追踪 ID（公共方法）
func GetTraceID(c *gin.Context) string {
	// 优先从 Go context 获取
	if traceID := TraceID(c.Request.Context()); traceID != "" {
		return traceID
	}
	// 其次从 Gin context 获取
	if traceID, exists := c.Get("trace_id"); exists {
		if tid, ok := traceID.(string); ok {
			return tid
		}
	}
	return ""
}

// getTraceID 从上下文获取追踪 ID（内部兼容方法）
func getTraceID(c *gin.Context) string {
	return GetTraceID(c)
}

// CORSMiddleware CORS 中间件（带日志）
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLoggerFromContext(c)

		origin := c.Request.Header.Get("Origin")

		// 记录 CORS 预检请求
		if c.Request.Method == "OPTIONS" {
			logger.Debug("CORS preflight request",
				zap.String("trace_id", getTraceID(c)),
				zap.String("path", c.Request.URL.Path),
				zap.String("origin", origin),
				zap.String("access_control_request_method", c.Request.Header.Get("Access-Control-Request-Method")),
			)
		}

		// 设置 CORS 头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
