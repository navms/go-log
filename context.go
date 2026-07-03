package xlog

import "context"

type traceKey struct{}
type spanKey struct{}

var (
	traceIDKey = traceKey{}
	spanIDKey  = spanKey{}
)

// WithTraceID stores trace_id in ctx.
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

// WithSpanID stores span_id in ctx.
func WithSpanID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, spanIDKey, id)
}

// TraceIDFromContext returns the trace id if present.
func TraceIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(traceIDKey).(string)
	return v, ok && v != ""
}

// SpanIDFromContext returns the span id if present.
func SpanIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(spanIDKey).(string)
	return v, ok && v != ""
}

func contextFields(ctx context.Context, keys ContextKeys) []any {
	if keys.TraceID == "" {
		keys = DefaultContextKeys()
	}
	var kvs []any
	if id, ok := TraceIDFromContext(ctx); ok {
		kvs = append(kvs, keys.TraceID, id)
	}
	if id, ok := SpanIDFromContext(ctx); ok {
		kvs = append(kvs, keys.SpanID, id)
	}
	// OpenTelemetry trace context via unexported keys is not accessible;
	// callers should use WithTraceID/WithSpanID or set values with matching keys.
	if id, ok := ctx.Value(keys.TraceID).(string); ok && id != "" {
		if len(kvs) == 0 || !containsKey(kvs, keys.TraceID) {
			kvs = append(kvs, keys.TraceID, id)
		}
	}
	if id, ok := ctx.Value(keys.SpanID).(string); ok && id != "" {
		if !containsKey(kvs, keys.SpanID) {
			kvs = append(kvs, keys.SpanID, id)
		}
	}
	return kvs
}

func containsKey(kvs []any, key string) bool {
	for i := 0; i+1 < len(kvs); i += 2 {
		if k, ok := kvs[i].(string); ok && k == key {
			return true
		}
	}
	return false
}

// WithContext returns a logger that includes fields extracted from ctx.
func WithContext(ctx context.Context) Logger {
	return defaultLogger.WithContext(ctx)
}
