package xlog

import (
	"context"
	stdlog "log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	zapbackend "github.com/navms/go-log/internal/zap"
)

// Logger is the sugared logging interface for application code.
type Logger interface {
	Trace(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)

	With(keysAndValues ...any) Logger
	WithContext(ctx context.Context) Logger
	Named(name string) Logger

	Level() Level
	SetLevel(level Level)
	Sync() error

	Typed() TypedLogger
	StdLogger() *stdlog.Logger
}

// TypedLogger is the zero-allocation logging interface for hot paths.
type TypedLogger interface {
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	With(fields ...Field) TypedLogger
	Zap() *zap.Logger
}

type zapLogger struct {
	backend *zapbackend.Backend
	ctxKeys ContextKeys
}

// New creates a Logger from options.
func New(opts ...Option) (Logger, error) {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return newFromConfig(cfg)
}

// NewWithConfig creates a Logger from an explicit Config, optionally overlaid with options.
func NewWithConfig(cfg Config, opts ...Option) (Logger, error) {
	for _, o := range opts {
		o(&cfg)
	}
	return newFromConfig(cfg)
}

func newFromConfig(cfg Config) (Logger, error) {
	bcfg, err := toBuildConfig(cfg)
	if err != nil {
		return nil, err
	}
	backend, err := zapbackend.New(bcfg)
	if err != nil {
		return nil, err
	}
	if cfg.ContextKeys.TraceID == "" && cfg.ContextKeys.SpanID == "" {
		cfg.ContextKeys = DefaultContextKeys()
	}
	return &zapLogger{backend: backend, ctxKeys: cfg.ContextKeys}, nil
}

func toBuildConfig(cfg Config) (zapbackend.BuildConfig, error) {
	zlvl, err := zapbackend.ParseLevel(cfg.Level.String())
	if err != nil {
		return zapbackend.BuildConfig{}, err
	}

	out := zapbackend.OutputSpec{}
	for _, o := range cfg.Outputs {
		switch o {
		case OutputStdout:
			out.Stdout = true
		case OutputStderr:
			out.Stderr = true
		case OutputFile:
			out.File = &zapbackend.FileSpec{
				Path:       cfg.File.Path,
				MaxSize:    cfg.File.Rotation.MaxSize,
				MaxBackups: cfg.File.Rotation.MaxBackups,
				MaxAge:     cfg.File.Rotation.MaxAge,
				Compress:   cfg.File.Rotation.Compress,
			}
		case OutputSyslog:
			out.Syslog = &zapbackend.SyslogSpec{
				Tag:     cfg.SyslogTag,
				Network: cfg.SyslogNetwork,
			}
		}
	}

	hooks := make([]zapbackend.HookFunc, len(cfg.Hooks))
	for i, h := range cfg.Hooks {
		hook := h
		hooks[i] = func(e zapbackend.HookEntry) error {
			return hook.Handle(fromHookEntry(e))
		}
	}

	return zapbackend.BuildConfig{
		Level:       zlvl,
		Format:      string(cfg.Format),
		Development: cfg.Development,
		Output:      out,
		Initial:     cfg.InitialFields,
		Hooks:       hooks,
		Sampling: zapbackend.SamplingSpec{
			Initial:    cfg.Sampling.Initial,
			Thereafter: cfg.Sampling.Thereafter,
		},
		SyncOnFatal: cfg.SyncOnFatal,
		Encoder:     toEncoderSpec(cfg.Encoder),
	}, nil
}

func fromHookEntry(e zapbackend.HookEntry) Entry {
	return Entry{
		Level:   Level(e.Level),
		Time:    time.Unix(0, e.Time),
		Message: e.Message,
		Fields:  e.Fields,
		Caller:  e.Caller,
	}
}

func (l *zapLogger) Trace(msg string, kvs ...any) {
	if ce := l.backend.Zap().Check(zapbackend.TraceLevel, msg); ce != nil {
		ce.Write(sugarKVToZapFields(kvs)...)
	}
}

func (l *zapLogger) traceForGlobal(msg string, kvs ...any) {
	zl := l.backend.Zap().WithOptions(zap.AddCallerSkip(1))
	if ce := zl.Check(zapbackend.TraceLevel, msg); ce != nil {
		ce.Write(sugarKVToZapFields(kvs)...)
	}
}

func (l *zapLogger) Debug(msg string, kvs ...any) { l.backend.Sugar().Debugw(msg, kvs...) }
func (l *zapLogger) Info(msg string, kvs ...any)  { l.backend.Sugar().Infow(msg, kvs...) }
func (l *zapLogger) Warn(msg string, kvs ...any)  { l.backend.Sugar().Warnw(msg, kvs...) }
func (l *zapLogger) Error(msg string, kvs ...any) { l.backend.Sugar().Errorw(msg, kvs...) }
func (l *zapLogger) Fatal(msg string, kvs ...any) { l.backend.Sugar().Fatalw(msg, kvs...) }

func (l *zapLogger) With(kvs ...any) Logger {
	return &zapLogger{backend: l.backend.WithSugar(kvs...), ctxKeys: l.ctxKeys}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	kvs := contextFields(ctx, l.ctxKeys)
	if len(kvs) == 0 {
		return l
	}
	return l.With(kvs...)
}

func (l *zapLogger) Named(name string) Logger {
	return &zapLogger{backend: l.backend.Named(name), ctxKeys: l.ctxKeys}
}

func (l *zapLogger) Level() Level {
	return fromZapLevel(l.backend.Level())
}

func (l *zapLogger) SetLevel(level Level) {
	l.backend.SetLevel(toZapLevel(level))
}

func (l *zapLogger) Sync() error { return l.backend.Sync() }

func (l *zapLogger) sugarForGlobal() *zap.SugaredLogger {
	return l.backend.Zap().WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func (l *zapLogger) Typed() TypedLogger { return &typedLogger{backend: l.backend} }

func (l *zapLogger) StdLogger() *stdlog.Logger {
	return stdlog.New(&stdWriter{l: l}, "", 0)
}

func fromZapLevel(l zapcore.Level) Level {
	if l == zapbackend.TraceLevel {
		return TraceLevel
	}
	return Level(l)
}

func toZapLevel(l Level) zapcore.Level {
	if l == TraceLevel {
		return zapbackend.TraceLevel
	}
	return zapcore.Level(l)
}

func sugarKVToZapFields(kvs []any) []zap.Field {
	if len(kvs) == 0 {
		return nil
	}
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "(MISSING)")
	}
	fields := make([]zap.Field, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			key = "invalid_key"
		}
		fields = append(fields, zap.Any(key, kvs[i+1]))
	}
	return fields
}

type typedLogger struct {
	backend *zapbackend.Backend
}

func (t *typedLogger) Trace(msg string, fields ...Field) {
	if ce := t.backend.Zap().Check(zapbackend.TraceLevel, msg); ce != nil {
		ce.Write(fieldsToZap(fields)...)
	}
}

func (t *typedLogger) Debug(msg string, fields ...Field) {
	t.backend.Zap().Debug(msg, fieldsToZap(fields)...)
}

func (t *typedLogger) Info(msg string, fields ...Field) {
	t.backend.Zap().Info(msg, fieldsToZap(fields)...)
}

func (t *typedLogger) Warn(msg string, fields ...Field) {
	t.backend.Zap().Warn(msg, fieldsToZap(fields)...)
}

func (t *typedLogger) Error(msg string, fields ...Field) {
	t.backend.Zap().Error(msg, fieldsToZap(fields)...)
}

func (t *typedLogger) Fatal(msg string, fields ...Field) {
	t.backend.Zap().Fatal(msg, fieldsToZap(fields)...)
}

func (t *typedLogger) With(fields ...Field) TypedLogger {
	zf := fieldsToZap(fields)
	return &typedLogger{backend: t.backend.WithFields(zf...)}
}

func (t *typedLogger) Zap() *zap.Logger { return t.backend.Zap() }

func fieldsToZap(fields []Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}
	out := make([]zap.Field, len(fields))
	for i, f := range fields {
		out[i] = f.toZap()
	}
	return out
}
