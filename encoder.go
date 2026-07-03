package xlog

import zapbackend "github.com/navms/go-log/internal/zap"

func toEncoderSpec(ec *EncoderConfig) *zapbackend.EncoderSpec {
	if ec == nil {
		return nil
	}
	return &zapbackend.EncoderSpec{
		TimeKey:           ec.TimeKey,
		LevelKey:          ec.LevelKey,
		NameKey:           ec.NameKey,
		CallerKey:         ec.CallerKey,
		MessageKey:        ec.MessageKey,
		StacktraceKey:     ec.StacktraceKey,
		TimeFormat:        ec.TimeFormat,
		LevelFormat:       string(ec.LevelFormat),
		CallerFormat:      string(ec.CallerFormat),
		DurationFormat:    string(ec.DurationFormat),
		DisableTime:       ec.DisableTime,
		DisableLevel:      ec.DisableLevel,
		DisableCaller:     ec.DisableCaller,
		DisableStacktrace: ec.DisableStacktrace,
		DisableLoggerName: ec.DisableLoggerName,
	}
}
