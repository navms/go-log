package xlog_test

import (
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
