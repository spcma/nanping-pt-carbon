package logger

import (
	"go.uber.org/zap"
)

// 标准化字段便捷函数

// WithTraceIDField 添加追踪 ID 字段
func WithTraceIDField(traceID interface{}) zap.Field {
	switch v := traceID.(type) {
	case string:
		return zap.String("trace_id", v)
	case int64:
		return zap.Int64("trace_id", v)
	default:
		return zap.Any("trace_id", traceID)
	}
}

// WithRequestID 添加请求 ID
func WithRequestID(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}

// WithUserID 添加用户 ID
func WithUserID(userID int64) zap.Field {
	return zap.Int64("user_id", userID)
}

// WithUsername 添加用户名
func WithUsername(username string) zap.Field {
	return zap.String("username", username)
}

// WithPath 添加请求路径
func WithPath(path string) zap.Field {
	return zap.String("path", path)
}

// WithMethod 添加请求方法
func WithMethod(method string) zap.Field {
	return zap.String("method", method)
}

// WithStatusCode 添加状态码
func WithStatusCode(statusCode int) zap.Field {
	return zap.Int("status_code", statusCode)
}

// WithLatency 添加延迟
func WithLatency(latency interface{}) zap.Field {
	switch v := latency.(type) {
	case float64:
		return zap.Float64("latency", v)
	case int64:
		return zap.Int64("latency", v)
	default:
		return zap.Any("latency", latency)
	}
}

// WithClientIP 添加客户端 IP
func WithClientIP(clientIP string) zap.Field {
	return zap.String("client_ip", clientIP)
}

// WithUserAgent 添加 User-Agent
func WithUserAgent(userAgent string) zap.Field {
	return zap.String("user_agent", userAgent)
}

// WithEventType 添加事件类型
func WithEventType(eventType string) zap.Field {
	return zap.String("event_type", eventType)
}

// WithError 添加错误字段
func WithError(err error) zap.Field {
	if err != nil {
		return zap.Error(err)
	}
	return zap.Skip()
}

// WithDuration 添加持续时间
func WithDuration(duration interface{}) zap.Field {
	switch v := duration.(type) {
	case float64:
		return zap.Float64("duration", v)
	case int64:
		return zap.Int64("duration", v)
	default:
		return zap.Any("duration", duration)
	}
}

// WithTableName 添加表名
func WithTableName(tableName string) zap.Field {
	return zap.String("table_name", tableName)
}

// WithOperation 添加操作名
func WithOperation(operation string) zap.Field {
	return zap.String("operation", operation)
}

// WithServiceName 添加服务名
func WithServiceName(serviceName string) zap.Field {
	return zap.String("service_name", serviceName)
}

// WithEndpoint 添加端点
func WithEndpoint(endpoint string) zap.Field {
	return zap.String("endpoint", endpoint)
}

// WithMessage 添加消息
func WithMessage(message string) zap.Field {
	return zap.String("message", message)
}

// WithData 添加数据
func WithData(data interface{}) zap.Field {
	return zap.Any("data", data)
}

// 业务场景组合字段

// WithDeviceInfo 添加设备信息
func WithDeviceInfo(deviceCode, deviceName string) []zap.Field {
	return []zap.Field{
		WithDeviceCode(deviceCode),
		zap.String("device_name", deviceName),
	}
}

// WithBusInfo 添加车辆信息
func WithBusInfo(busID int64, license string) []zap.Field {
	return []zap.Field{
		WithBusID(busID),
		zap.String("bus_license", license),
	}
}

// WithProcessingResult 添加处理结果
func WithProcessingResult(success bool, message string) []zap.Field {
	return []zap.Field{
		zap.Bool("processing_success", success),
		WithMessage(message),
	}
}

// WithTiming 添加时间信息
func WithTiming(start, end int64) []zap.Field {
	duration := end - start
	return []zap.Field{
		zap.Int64("start_time", start),
		zap.Int64("end_time", end),
		zap.Int64("duration_ms", duration),
	}
}

// WithErrorDetails 添加错误详情
func WithErrorDetails(err error, context string) []zap.Field {
	if err == nil {
		return []zap.Field{}
	}

	return []zap.Field{
		WithError(err),
		zap.String("error_context", context),
	}
}

// WithPerformanceMetrics 添加性能指标
func WithPerformanceMetrics(latencyMs int64, memoryMB float64, cpuPercent float64) []zap.Field {
	return []zap.Field{
		zap.Int64("latency_ms", latencyMs),
		zap.Float64("memory_mb", memoryMB),
		zap.Float64("cpu_percent", cpuPercent),
	}
}

// WithPagination 添加分页信息
func WithPagination(pageNum, pageSize, totalCount int64) []zap.Field {
	return []zap.Field{
		zap.Int64("page_num", pageNum),
		zap.Int64("page_size", pageSize),
		zap.Int64("total_count", totalCount),
		zap.Int64("total_pages", (totalCount+pageSize-1)/pageSize),
	}
}

// WithCacheInfo 添加缓存信息
func WithCacheInfo(cacheKey string, hit bool, ttlSeconds int64) []zap.Field {
	return []zap.Field{
		zap.String("cache_key", cacheKey),
		zap.Bool("cache_hit", hit),
		zap.Int64("ttl_seconds", ttlSeconds),
	}
}

// WithFileInfo 添加文件信息
func WithFileInfo(fileName, filePath string, fileSize int64) []zap.Field {
	return []zap.Field{
		zap.String("file_name", fileName),
		zap.String("file_path", filePath),
		zap.Int64("file_size", fileSize),
	}
}

// WithConnectionInfo 添加连接信息
func WithConnectionInfo(host, port string, isConnected bool) []zap.Field {
	return []zap.Field{
		zap.String("host", host),
		zap.String("port", port),
		zap.Bool("is_connected", isConnected),
	}
}

// 特定业务字段

// WithDeviceCode 添加设备编码
func WithDeviceCode(deviceCode string) zap.Field {
	return zap.String("device_code", deviceCode)
}

// WithBusID 添加车辆 ID
func WithBusID(busID int64) zap.Field {
	return zap.Int64("bus_id", busID)
}

// WithBusRouteID 添加线路 ID
func WithBusRouteID(busRouteID int64) zap.Field {
	return zap.Int64("bus_route_id", busRouteID)
}
