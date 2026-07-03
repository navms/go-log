package zap

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SamplingSpec configures log sampling for verbose levels.
type SamplingSpec struct {
	Initial    int
	Thereafter int
}

// BuildConfig is the internal configuration for constructing a zap backend.
type BuildConfig struct {
	Level       zapcore.Level
	Format      string // "console" | "json"
	Development bool
	Output      OutputSpec
	Initial     map[string]any
	Hooks       []HookFunc
	Sampling    SamplingSpec
	SyncOnFatal bool
	Encoder     *EncoderSpec
}

// Backend wraps zap loggers and dynamic level control.
type Backend struct {
	zap         *zap.Logger
	sugar       *zap.SugaredLogger
	atomicLevel zap.AtomicLevel
}

// New builds a Backend from BuildConfig.
func New(cfg BuildConfig) (*Backend, error) {
	ws, err := buildWriteSyncer(cfg.Output)
	if err != nil {
		return nil, err
	}

	atomicLevel := zap.NewAtomicLevelAt(cfg.Level)
	enc := newEncoder(cfg.Format, cfg.Development, cfg.Encoder)
	core := buildCore(enc, ws, atomicLevel, cfg.Sampling)
	core = newHookCore(core, cfg.Hooks)

	opts := []zap.Option{zap.AddCallerSkip(1)}
	if cfg.Encoder == nil || !cfg.Encoder.DisableCallerCapture() {
		opts = append(opts, zap.AddCaller())
	}
	if cfg.Development {
		opts = append(opts, zap.Development())
	}
	if cfg.SyncOnFatal {
		opts = append(opts, zap.WithFatalHook(syncFatalHook{syncFn: ws.Sync}))
	}

	fields := mapToFields(cfg.Initial)
	zl := zap.New(core, opts...).With(fields...)
	return &Backend{
		zap:         zl,
		sugar:       zl.Sugar(),
		atomicLevel: atomicLevel,
	}, nil
}

func buildCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, level zap.AtomicLevel, sampling SamplingSpec) zapcore.Core {
	if sampling.Initial <= 0 || sampling.Thereafter <= 0 {
		return zapcore.NewCore(enc, ws, level)
	}

	low := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return level.Enabled(l) && l < zapcore.WarnLevel
	})
	high := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return level.Enabled(l) && l >= zapcore.WarnLevel
	})

	lowCore := zapcore.NewCore(enc, ws, low)
	lowCore = zapcore.NewSamplerWithOptions(lowCore, time.Second, sampling.Initial, sampling.Thereafter)
	highCore := zapcore.NewCore(enc, ws, high)
	return zapcore.NewTee(lowCore, highCore)
}

func mapToFields(m map[string]any) []zap.Field {
	if len(m) == 0 {
		return nil
	}
	fields := make([]zap.Field, 0, len(m))
	for k, v := range m {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// Zap returns the underlying zap.Logger.
func (b *Backend) Zap() *zap.Logger { return b.zap }

// Sugar returns the sugared logger.
func (b *Backend) Sugar() *zap.SugaredLogger { return b.sugar }

// AtomicLevel returns the dynamic level controller.
func (b *Backend) AtomicLevel() zap.AtomicLevel { return b.atomicLevel }

// Sync flushes buffered logs.
func (b *Backend) Sync() error { return b.zap.Sync() }

// SetLevel changes the minimum enabled level at runtime.
func (b *Backend) SetLevel(l zapcore.Level) { b.atomicLevel.SetLevel(l) }

// Level returns the current minimum level.
func (b *Backend) Level() zapcore.Level { return b.atomicLevel.Level() }

// WithSugar returns a new Backend with additional sugared fields.
func (b *Backend) WithSugar(keysAndValues ...any) *Backend {
	return &Backend{
		zap:         b.zap.With(sugarKVToFields(keysAndValues)...),
		sugar:       b.sugar.With(keysAndValues...),
		atomicLevel: b.atomicLevel,
	}
}

// WithFields returns a new Backend with typed fields.
func (b *Backend) WithFields(fields ...zap.Field) *Backend {
	zl := b.zap.With(fields...)
	return &Backend{
		zap:         zl,
		sugar:       zl.Sugar(),
		atomicLevel: b.atomicLevel,
	}
}

// Named returns a new Backend with a name segment.
func (b *Backend) Named(name string) *Backend {
	zl := b.zap.Named(name)
	return &Backend{
		zap:         zl,
		sugar:       zl.Sugar(),
		atomicLevel: b.atomicLevel,
	}
}

func sugarKVToFields(keysAndValues []any) []zap.Field {
	if len(keysAndValues) == 0 {
		return nil
	}
	if len(keysAndValues)%2 != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")
	}
	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keysAndValues[i])
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}
	return fields
}
