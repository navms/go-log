package xlog

import (
	"net/http"
	"sync"

	"go.uber.org/zap"
)

var (
	mu            sync.RWMutex
	defaultLogger = NewNoop()
	atomicLevel   zap.AtomicLevel
	hasAtomic     bool
)

// Init builds the default logger from options and sets it as the global logger.
func Init(opts ...Option) error {
	l, err := New(opts...)
	if err != nil {
		return err
	}
	SetDefault(l)
	return nil
}

// InitFromEnv loads config from environment variables then Init.
func InitFromEnv(opts ...Option) error {
	cfg := ConfigFromEnv()
	for _, o := range opts {
		o(&cfg)
	}
	l, err := newFromConfig(cfg)
	if err != nil {
		return err
	}
	SetDefault(l)
	return nil
}

// SetDefault replaces the global logger.
func SetDefault(l Logger) {
	mu.Lock()
	defaultLogger = l
	hasAtomic = false
	if zl, ok := l.(*zapLogger); ok {
		atomicLevel = zl.backend.AtomicLevel()
		hasAtomic = true
	}
	mu.Unlock()
}

// Default returns the global logger.
func Default() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return defaultLogger
}

func withZapLogger(zapFn func(*zapLogger), fallback func(Logger)) {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	if zl, ok := l.(*zapLogger); ok {
		zapFn(zl)
		return
	}
	fallback(l)
}

func Trace(msg string, kv ...any) {
	withZapLogger(
		func(zl *zapLogger) { zl.traceForGlobal(msg, kv...) },
		func(l Logger) { l.Trace(msg, kv...) },
	)
}

func Debug(msg string, kv ...any) {
	withZapLogger(
		func(zl *zapLogger) { zl.sugarForGlobal().Debugw(msg, kv...) },
		func(l Logger) { l.Debug(msg, kv...) },
	)
}

func Info(msg string, kv ...any) {
	withZapLogger(
		func(zl *zapLogger) { zl.sugarForGlobal().Infow(msg, kv...) },
		func(l Logger) { l.Info(msg, kv...) },
	)
}

func Warn(msg string, kv ...any) {
	withZapLogger(
		func(zl *zapLogger) { zl.sugarForGlobal().Warnw(msg, kv...) },
		func(l Logger) { l.Warn(msg, kv...) },
	)
}

func Error(msg string, kv ...any) {
	withZapLogger(
		func(zl *zapLogger) { zl.sugarForGlobal().Errorw(msg, kv...) },
		func(l Logger) { l.Error(msg, kv...) },
	)
}

func Fatal(msg string, kv ...any) {
	withZapLogger(
		func(zl *zapLogger) { zl.sugarForGlobal().Fatalw(msg, kv...) },
		func(l Logger) { l.Fatal(msg, kv...) },
	)
}

func With(kv ...any) Logger { return Default().With(kv...) }

// Typed returns the typed logger for the global logger.
func Typed() TypedLogger { return Default().Typed() }

// Sync flushes the global logger.
func Sync() error { return Default().Sync() }

// SetLevel changes the global logger level at runtime.
func SetLevel(level Level) { Default().SetLevel(level) }

// LevelHandler returns an HTTP handler to adjust log level dynamically (JSON body).
func LevelHandler() http.Handler {
	mu.RLock()
	al := atomicLevel
	ok := hasAtomic
	mu.RUnlock()
	if !ok {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"logger not initialized with dynamic level"}`))
		})
	}
	return al
}
