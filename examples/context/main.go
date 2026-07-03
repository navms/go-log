// Context 链路：自动注入 trace_id / span_id
package main

import (
	"context"

	"github.com/navms/go-log"
)

func main() {
	if err := xlog.Init(xlog.WithJSON()); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	ctx := xlog.WithTraceID(context.Background(), "trace-abc-123")
	ctx = xlog.WithSpanID(ctx, "span-456")

	handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
	// WithContext 自动把 trace_id / span_id 写入日志字段
	xlog.WithContext(ctx).Info("request handled", "latency_ms", 42)
}
