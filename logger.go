package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var std = New(nil)

// Env 运行环境
type Env string

// Level 日志级别
type Level = zapcore.Level

// Output 输出方式
type Output string

const (
	Dev  Env = "dev"
	Prod Env = "prod"

	Console        Output = "console"
	File           Output = "file"
	FileAndConsole Output = "file-and-console"
)

const (
	DebugL      = zapcore.DebugLevel
	InfoL       = zapcore.InfoLevel
	WarnL       = zapcore.WarnLevel
	ErrorL      = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicL      = zapcore.PanicLevel
	FatalL      = zapcore.FatalLevel
)

type Logger struct {
	l  *zap.Logger
	s  *zap.SugaredLogger
	al *zap.AtomicLevel

	config     *Config
	writers    []io.Writer
	fileLogger *lumberjack.Logger
}

func New(cfg *Config) *Logger {
	if cfg == nil {
		return New(NewDevelopConfig())
	}

	log := &Logger{
		config: cfg,
	}
	if err := log.init(); err != nil {
		panic(err)
	}

	return log
}

func (l *Logger) init() error {
	if err := l.validateConfig(); err != nil {
		return err
	}

	levelAt := zap.NewAtomicLevelAt(l.config.Level)
	l.al = &levelAt

	var ws []zapcore.WriteSyncer
	var fl *lumberjack.Logger

	// 配置文件输出
	if l.config.Output == File || l.config.Output == FileAndConsole {
		fl = &lumberjack.Logger{
			Filename:   l.config.Filename,
			MaxSize:    l.config.MaxSize,
			MaxBackups: l.config.MaxBackups,
			MaxAge:     l.config.MaxAge,
			Compress:   l.config.Compress,
			LocalTime:  true,
		}
		ws = append(ws, zapcore.AddSync(fl))
		l.fileLogger = fl
	}

	// 配置标准输出
	if l.config.Output == FileAndConsole || l.config.Output == Console {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}

	// 添加动态 writers
	for _, w := range l.writers {
		ws = append(ws, zapcore.AddSync(w))
	}

	if len(ws) == 0 {
		return errors.New("logger has no output write syncer configured")
	}

	// 合并所有输出
	writeSyncer := zapcore.NewMultiWriteSyncer(ws...)

	// 构建 encoder
	var encoder zapcore.Encoder
	ec := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	switch l.config.Env {
	case Dev:
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(ec)
	case Prod:
		ec.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(ec)
	default:
		ec.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewConsoleEncoder(ec)
	}

	// 采样级别
	sampledLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return l.al.Enabled(level) && (level == zapcore.DebugLevel || level == zapcore.InfoLevel)
	})
	unsampledLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return l.al.Enabled(level) && level >= zapcore.WarnLevel
	})

	sampledCore := zapcore.NewCore(encoder, writeSyncer, sampledLevel)
	unsampledCore := zapcore.NewCore(encoder, writeSyncer, unsampledLevel)

	// 仅对 Debug/Info 做采样，Warn 及以上不采样
	if l.config.SamplingInitial > 0 && l.config.SamplingThereafter > 0 {
		sampledCore = zapcore.NewSamplerWithOptions(
			sampledCore,
			time.Second,
			l.config.SamplingInitial,
			l.config.SamplingThereafter,
		)
	}
	core := zapcore.NewTee(sampledCore, unsampledCore)

	opts := []zap.Option{
		zap.ErrorOutput(zapcore.AddSync(os.Stderr)),
		zap.AddStacktrace(zap.PanicLevel),
		zap.WithFatalHook(SyncFatalHook{logger: l}),
	}
	if l.config.ShowCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(2))
	}

	zapLogger := zap.New(core, opts...)
	l.l = zapLogger
	l.s = zapLogger.Sugar()

	if l.config.RedirectStdLog {
		zap.RedirectStdLog(l.l)
	}
	return nil
}

func (l *Logger) validateConfig() error {
	if l.config == nil {
		return errors.New("logger config is nil")
	}

	switch l.config.Output {
	case Console, File, FileAndConsole:
	default:
		return fmt.Errorf("invalid logger output: %q", l.config.Output)
	}

	if (l.config.Output == File || l.config.Output == FileAndConsole) && strings.TrimSpace(l.config.Filename) == "" {
		return errors.New("filename is required when output includes file")
	}

	return nil
}

func (l *Logger) SetLevel(level Level) {
	if l.al != nil {
		l.al.SetLevel(level)
	}
}

func (l *Logger) Sync() error { return l.l.Sync() }

func (l *Logger) Close() error {
	if l == nil {
		return nil
	}
	err := l.Sync()
	if err != nil {
		return err
	}
	if l.fileLogger != nil {
		err = l.fileLogger.Close()
		if err != nil {
			return err
		}
		l.fileLogger = nil
	}
	return nil
}

func (l *Logger) Debug(msg string, fields ...zap.Field) { l.l.Debug(msg, fields...) }
func (l *Logger) Info(msg string, fields ...zap.Field)  { l.l.Info(msg, fields...) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.l.Warn(msg, fields...) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.l.Error(msg, fields...) }
func (l *Logger) Panic(msg string, fields ...zap.Field) { l.l.Panic(msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...zap.Field) { l.l.Fatal(msg, fields...) }
func (l *Logger) Debugf(template string, args ...any)   { l.s.Debugf(template, args...) }
func (l *Logger) Infof(template string, args ...any)    { l.s.Infof(template, args...) }
func (l *Logger) Warnf(template string, args ...any)    { l.s.Warnf(template, args...) }
func (l *Logger) Errorf(template string, args ...any)   { l.s.Errorf(template, args...) }
func (l *Logger) Panicf(template string, args ...any)   { l.s.Panicf(template, args...) }
func (l *Logger) Fatalf(template string, args ...any)   { l.s.Fatalf(template, args...) }

func Default() *Logger                      { return std }
func ReplaceDefault(l *Logger)              { std = l }
func Debug(msg string, fields ...zap.Field) { std.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { std.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { std.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { std.Error(msg, fields...) }
func Panic(msg string, fields ...zap.Field) { std.Panic(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { std.Fatal(msg, fields...) }
func Debugf(template string, args ...any)   { std.Debugf(template, args...) }
func Infof(template string, args ...any)    { std.Infof(template, args...) }
func Warnf(template string, args ...any)    { std.Warnf(template, args...) }
func Errorf(template string, args ...any)   { std.Errorf(template, args...) }
func Panicf(template string, args ...any)   { std.Panicf(template, args...) }
func Fatalf(template string, args ...any)   { std.Fatalf(template, args...) }
func Sync() error                           { return std.Sync() }
func Close() error                          { return std.Close() }
func SetLevel(level Level)                  { std.SetLevel(level) }
