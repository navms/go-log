package xlog_test

import (
	"context"
	"testing"

	"github.com/navms/go-log"
	"github.com/navms/go-log/logtest"
)

func TestWithContext(t *testing.T) {
	obs, l := logtest.NewObserver(xlog.InfoLevel)
	ctx := xlog.WithTraceID(context.Background(), "trace-1")
	ctx = xlog.WithSpanID(ctx, "span-1")
	l.WithContext(ctx).Info("req")
	logtest.RequireLogContains(t, obs, "req", map[string]any{
		"trace_id": "trace-1",
		"span_id":  "span-1",
	})
}

func TestNamed(t *testing.T) {
	obs, l := logtest.NewObserver(xlog.InfoLevel)
	l.Named("http").Info("listen")
	entries := obs.Entries()
	if len(entries) == 0 {
		t.Fatal("no entries")
	}
}
