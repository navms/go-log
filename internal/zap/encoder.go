package zap

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// EncoderSpec holds encoder customization (internal, no xlog import).
type EncoderSpec struct {
	TimeKey       string
	LevelKey      string
	NameKey       string
	CallerKey     string
	MessageKey    string
	StacktraceKey string
	TimeFormat    string
	LevelFormat   string
	CallerFormat  string
	DurationFormat string

	DisableTime       bool
	DisableLevel      bool
	DisableCaller     bool
	DisableStacktrace bool
	DisableLoggerName bool
}

func newEncoder(format string, development bool, spec *EncoderSpec) zapcore.Encoder {
	cfg := buildEncoderConfig(development, spec)
	switch format {
	case "json":
		return zapcore.NewJSONEncoder(cfg)
	default:
		return zapcore.NewConsoleEncoder(cfg)
	}
}

func buildEncoderConfig(development bool, spec *EncoderSpec) zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    pickLevelEncoder(development, spec),
	}

	if spec != nil {
		if spec.DisableTime {
			cfg.TimeKey = zapcore.OmitKey
		} else if spec.TimeKey != "" {
			cfg.TimeKey = spec.TimeKey
		}
		if spec.DisableLevel {
			cfg.LevelKey = zapcore.OmitKey
		} else if spec.LevelKey != "" {
			cfg.LevelKey = spec.LevelKey
		}
		if spec.DisableLoggerName {
			cfg.NameKey = zapcore.OmitKey
		} else if spec.NameKey != "" {
			cfg.NameKey = spec.NameKey
		}
		if spec.DisableCaller {
			cfg.CallerKey = zapcore.OmitKey
		} else if spec.CallerKey != "" {
			cfg.CallerKey = spec.CallerKey
		}
		if spec.DisableStacktrace {
			cfg.StacktraceKey = zapcore.OmitKey
		} else if spec.StacktraceKey != "" {
			cfg.StacktraceKey = spec.StacktraceKey
		}
		if spec.MessageKey != "" {
			cfg.MessageKey = spec.MessageKey
		}
		if spec.TimeFormat != "" {
			layout := spec.TimeFormat
			cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(layout))
			}
		}
		if spec.CallerFormat != "" {
			cfg.EncodeCaller = pickCallerEncoder(spec.CallerFormat)
		}
		if spec.DurationFormat != "" {
			cfg.EncodeDuration = pickDurationEncoder(spec.DurationFormat)
		}
		if spec.LevelFormat != "" {
			cfg.EncodeLevel = pickLevelEncoderFromFormat(spec.LevelFormat)
		}
	}

	return cfg
}

func pickLevelEncoder(development bool, spec *EncoderSpec) zapcore.LevelEncoder {
	if spec != nil && spec.LevelFormat != "" {
		return pickLevelEncoderFromFormat(spec.LevelFormat)
	}
	if development {
		return traceAwareCapitalColorLevelEncoder
	}
	return traceAwareLowercaseLevelEncoder
}

func pickLevelEncoderFromFormat(format string) zapcore.LevelEncoder {
	switch format {
	case "capitalColor":
		return traceAwareCapitalColorLevelEncoder
	case "capital":
		return traceAwareCapitalLevelEncoder
	default:
		return traceAwareLowercaseLevelEncoder
	}
}

func traceAwareCapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == TraceLevel {
		enc.AppendString("TRACE")
		return
	}
	zapcore.CapitalLevelEncoder(l, enc)
}

func pickCallerEncoder(format string) zapcore.CallerEncoder {
	switch format {
	case "full":
		return zapcore.FullCallerEncoder
	default:
		return zapcore.ShortCallerEncoder
	}
}

func pickDurationEncoder(format string) zapcore.DurationEncoder {
	switch format {
	case "seconds":
		return zapcore.SecondsDurationEncoder
	case "string":
		return zapcore.StringDurationEncoder
	default:
		return zapcore.MillisDurationEncoder
	}
}

func traceAwareLowercaseLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == TraceLevel {
		enc.AppendString("trace")
		return
	}
	zapcore.LowercaseLevelEncoder(l, enc)
}

func traceAwareCapitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == TraceLevel {
		enc.AppendString("\x1b[36mTRACE\x1b[0m")
		return
	}
	zapcore.CapitalColorLevelEncoder(l, enc)
}

// DisableCaller reports whether caller capture should be skipped.
func (s *EncoderSpec) DisableCallerCapture() bool {
	return s != nil && s.DisableCaller
}
