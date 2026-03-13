package logger

import "context"

type traceIDKeyType struct{}

var traceIDKey = traceIDKeyType{}

// 写入 trace_id
func WithTrace(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceIDKey, traceID)
}

// 读取 trace_id
func TraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(traceIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
