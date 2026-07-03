package zap

import "go.uber.org/zap/zapcore"

// TraceLevel is one step below Debug in zap's level ordering.
const TraceLevel = zapcore.Level(-2)

// ParseLevel parses a level string including "trace".
func ParseLevel(s string) (zapcore.Level, error) {
	switch s {
	case "trace":
		return TraceLevel, nil
	default:
		var l zapcore.Level
		err := l.UnmarshalText([]byte(s))
		return l, err
	}
}

// LevelString returns the string representation including trace.
func LevelString(l zapcore.Level) string {
	if l == TraceLevel {
		return "trace"
	}
	return l.String()
}
