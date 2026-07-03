package xlog

import (
	"context"
	stdlog "log"

	"go.uber.org/zap"
)

type noopLogger struct{}

// NewNoop returns a logger that discards all output.
func NewNoop() Logger { return noopLogger{} }

func (noopLogger) Trace(string, ...any) {}
func (noopLogger) Debug(string, ...any) {}
func (noopLogger) Info(string, ...any)  {}
func (noopLogger) Warn(string, ...any)  {}
func (noopLogger) Error(string, ...any) {}
func (noopLogger) Fatal(string, ...any) {}

func (n noopLogger) With(...any) Logger                 { return n }
func (n noopLogger) WithContext(context.Context) Logger { return n }
func (n noopLogger) Named(string) Logger                { return n }
func (noopLogger) Level() Level                         { return InfoLevel }
func (noopLogger) SetLevel(Level)                       {}
func (noopLogger) Sync() error                          { return nil }
func (n noopLogger) Typed() TypedLogger                 { return noopTyped{} }
func (noopLogger) StdLogger() *stdlog.Logger {
	return stdlog.New(noopWriter{}, "", 0)
}

type noopTyped struct{}

func (noopTyped) Trace(string, ...Field) {}
func (noopTyped) Debug(string, ...Field) {}
func (noopTyped) Info(string, ...Field)  {}
func (noopTyped) Warn(string, ...Field)  {}
func (noopTyped) Error(string, ...Field) {}
func (noopTyped) Fatal(string, ...Field) {}
func (n noopTyped) With(...Field) TypedLogger { return n }
func (noopTyped) Zap() *zap.Logger            { return zap.NewNop() }

type noopWriter struct{}

func (noopWriter) Write(p []byte) (int, error) { return len(p), nil }
