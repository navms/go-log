package zap

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestHookCoreWrite(t *testing.T) {
	var called bool
	enc := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{MessageKey: "msg"})
	core := zapcore.NewCore(enc, zapcore.AddSync(discardWriter{}), zapcore.DebugLevel)
	wrapped := newHookCore(core, []HookFunc{
		func(e HookEntry) error {
			called = true
			if e.Message != "hi" {
				t.Fatalf("msg = %q", e.Message)
			}
			return nil
		},
	})
	ent := zapcore.Entry{Level: zapcore.InfoLevel, Message: "hi"}
	if err := wrapped.Write(ent, nil); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("hook not called")
	}
}

func TestZapLoggerWithHook(t *testing.T) {
	var called bool
	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "message",
		LevelKey:   "level",
	})
	base := zapcore.NewCore(enc, zapcore.AddSync(discardWriter{}), zapcore.DebugLevel)
	hc := newHookCore(base, []HookFunc{
		func(e HookEntry) error {
			called = true
			return nil
		},
	})
	l := zap.New(hc, zap.AddCaller())
	l.Info("hook")
	if !called {
		t.Fatal("hook not called via zap.Logger")
	}
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }
