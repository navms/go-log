package xlog_test

import (
	"strings"
	"testing"

	"github.com/navms/go-log"
	"github.com/navms/go-log/logtest"
)

func TestInit(t *testing.T) {
	if err := xlog.Init(xlog.WithLevel(xlog.DebugLevel), xlog.WithConsole()); err != nil {
		t.Fatal(err)
	}
	xlog.Info("init ok")
}

func logFromGlobalHelper(t *testing.T) {
	t.Helper()
	xlog.Info("caller site")
}

func TestGlobalInfoCaller(t *testing.T) {
	var gotCaller string
	rec := &callerRecorder{}
	if err := xlog.Init(
		xlog.WithLevel(xlog.InfoLevel),
		xlog.WithConsole(),
		xlog.WithHook(rec),
	); err != nil {
		t.Fatal(err)
	}
	logFromGlobalHelper(t)
	if rec.caller == "" {
		t.Fatal("expected caller in hook entry")
	}
	gotCaller = rec.caller
	if strings.Contains(gotCaller, "global.go") {
		t.Fatalf("caller = %q, should not point at global wrapper", gotCaller)
	}
	if !strings.Contains(gotCaller, "global_test.go") {
		t.Fatalf("caller = %q, want global_test.go call site", gotCaller)
	}
}

func TestLoggerInfoCaller(t *testing.T) {
	rec := &callerRecorder{}
	l, err := xlog.New(
		xlog.WithLevel(xlog.InfoLevel),
		xlog.WithConsole(),
		xlog.WithHook(rec),
	)
	if err != nil {
		t.Fatal(err)
	}
	logFromLoggerHelper(t, l)
	if rec.caller == "" {
		t.Fatal("expected caller in hook entry")
	}
	if strings.Contains(rec.caller, "logger.go") {
		t.Fatalf("caller = %q, should not point at logger wrapper", rec.caller)
	}
	if !strings.Contains(rec.caller, "global_test.go") {
		t.Fatalf("caller = %q, want global_test.go call site", rec.caller)
	}
}

func logFromLoggerHelper(t *testing.T, l xlog.Logger) {
	t.Helper()
	l.Info("caller site")
}

type callerRecorder struct {
	caller string
}

func (r *callerRecorder) Handle(e xlog.Entry) error {
	r.caller = e.Caller
	return nil
}

func TestWithFields(t *testing.T) {
	obs, l := logtest.NewObserver(xlog.InfoLevel)
	child := l.With("service", "api")
	child.Info("started")
	logtest.RequireLogContains(t, obs, "started", map[string]any{"service": "api"})
}

func TestInitFromEnv(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_FORMAT", "json")
	cfg := xlog.ConfigFromEnv()
	if cfg.Level != xlog.DebugLevel {
		t.Fatalf("level = %v", cfg.Level)
	}
	if cfg.Format != xlog.FormatJSON {
		t.Fatalf("format = %v", cfg.Format)
	}
}
