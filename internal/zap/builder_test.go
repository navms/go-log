package zap_test

import (
	"os"
	"path/filepath"
	"testing"

	zapbackend "github.com/navms/go-log/internal/zap"
)

func TestBuildConsoleJSON(t *testing.T) {
	for _, format := range []string{"console", "json"} {
		b, err := zapbackend.New(zapbackend.BuildConfig{
			Level:  0,
			Format: format,
			Output: zapbackend.OutputSpec{Stdout: true},
		})
		if err != nil {
			t.Fatalf("format %s: %v", format, err)
		}
		b.Sugar().Info("snapshot")
		_ = b.Sync()
	}
}

func TestFileRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	b, err := zapbackend.New(zapbackend.BuildConfig{
		Level:  0,
		Format: "json",
		Output: zapbackend.OutputSpec{
			File: &zapbackend.FileSpec{
				Path:     path,
				MaxSize:  1,
				MaxBackups: 1,
				MaxAge:   1,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	b.Sugar().Info("file log")
	_ = b.Sync()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("log file: %v", err)
	}
}

func TestHookCore(t *testing.T) {
	var n int
	b, err := zapbackend.New(zapbackend.BuildConfig{
		Level:  0,
		Format: "json",
		Output: zapbackend.OutputSpec{Stdout: true},
		Hooks: []zapbackend.HookFunc{
			func(e zapbackend.HookEntry) error {
				n++
				return nil
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	b.Sugar().Info("hook")
	if n != 1 {
		t.Fatalf("hook calls = %d", n)
	}
}

func TestTraceLevel(t *testing.T) {
	b, err := zapbackend.New(zapbackend.BuildConfig{
		Level:  zapbackend.TraceLevel,
		Format: "console",
		Output: zapbackend.OutputSpec{Stdout: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if ce := b.Zap().Check(zapbackend.TraceLevel, "trace msg"); ce != nil {
		ce.Write()
	}
}
