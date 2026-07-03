package xlog_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/navms/go-log"
)

func TestDefaultLogFilePath(t *testing.T) {
	path := xlog.DefaultLogFilePath("myapp")
	if !strings.HasSuffix(path, filepath.Join("myapp", "logs", "app.log")) &&
		!strings.HasSuffix(path, filepath.Join("myapp", "app.log")) &&
		!strings.Contains(path, "myapp") {
		t.Fatalf("unexpected path: %s", path)
	}
	switch runtime.GOOS {
	case "darwin":
		if !strings.Contains(path, filepath.Join("Library", "Logs")) {
			t.Fatalf("darwin path = %s", path)
		}
	}
}

func TestProductionConfigPreset(t *testing.T) {
	cfg := xlog.NewProductionConfig("demo")
	if cfg.Format != xlog.FormatJSON {
		t.Fatalf("format = %s", cfg.Format)
	}
	if len(cfg.Outputs) != 1 || cfg.Outputs[0] != xlog.OutputFile {
		t.Fatalf("outputs = %v", cfg.Outputs)
	}
	if cfg.File.Path == "" {
		t.Fatal("empty file path")
	}
	if !cfg.SyncOnFatal {
		t.Fatal("sync on fatal expected")
	}
}

func TestDevelopmentConfigPreset(t *testing.T) {
	cfg := xlog.NewDevelopmentConfig()
	if cfg.Format != xlog.FormatConsole || !cfg.Development {
		t.Fatalf("dev config = %+v", cfg)
	}
}

func TestNewWithConfig(t *testing.T) {
	l, err := xlog.NewWithConfig(xlog.NewDevelopmentConfig())
	if err != nil {
		t.Fatal(err)
	}
	l.Debug("preset ok")
}

func TestSamplingOption(t *testing.T) {
	l, err := xlog.New(
		xlog.WithConsole(),
		xlog.WithLevel(xlog.DebugLevel),
		xlog.WithSampling(10, 10),
	)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		l.Debug("sample", "i", i)
	}
	l.Warn("always logged")
}
