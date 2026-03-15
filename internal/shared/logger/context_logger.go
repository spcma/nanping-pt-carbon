package logger

import (
	"context"
	"go.uber.org/zap"
)

func fieldsWithTrace(ctx context.Context, fields []zap.Field) []zap.Field {
	traceID := TraceID(ctx)
	if traceID == "" {
		return fields
	}

	// 避免修改原 slice
	newFields := make([]zap.Field, 0, len(fields)+1)
	newFields = append(newFields, zap.String("trace_id", traceID))
	newFields = append(newFields, fields...)

	return newFields
}

func DebugCtx(ctx context.Context, category, msg string, fields ...zap.Field) {
	Debug(category, msg, fieldsWithTrace(ctx, fields)...)
}

func InfoCtx(ctx context.Context, category, msg string, fields ...zap.Field) {
	Info(category, msg, fieldsWithTrace(ctx, fields)...)
}

func WarnCtx(ctx context.Context, category, msg string, fields ...zap.Field) {
	Warn(category, msg, fieldsWithTrace(ctx, fields)...)
}

func ErrorCtx(ctx context.Context, category, msg string, fields ...zap.Field) {
	Error(category, msg, fieldsWithTrace(ctx, fields)...)
}
