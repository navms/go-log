package xlog

// LevelFormat controls how log levels are rendered.
type LevelFormat string

const (
	LevelFormatLowercase     LevelFormat = "lowercase"
	LevelFormatCapital       LevelFormat = "capital"
	LevelFormatCapitalColor  LevelFormat = "capitalColor"
)

// CallerFormat controls caller encoding.
type CallerFormat string

const (
	CallerFormatShort CallerFormat = "short"
	CallerFormatFull  CallerFormat = "full"
)

// DurationFormat controls duration field encoding.
type DurationFormat string

const (
	DurationFormatMillis  DurationFormat = "millis"
	DurationFormatSeconds DurationFormat = "seconds"
	DurationFormatString  DurationFormat = "string"
)

// EncoderConfig customizes log field keys and encodings.
// Unset string fields keep xlog defaults; use Disable* flags to omit a field.
type EncoderConfig struct {
	TimeKey       string
	LevelKey      string
	NameKey       string
	CallerKey     string
	MessageKey    string
	StacktraceKey string

	// TimeFormat is a Go time layout, e.g. "2006-01-02 15:04:05.000".
	// Empty uses ISO8601.
	TimeFormat string

	LevelFormat    LevelFormat
	CallerFormat   CallerFormat
	DurationFormat DurationFormat

	DisableTime       bool
	DisableLevel      bool
	DisableCaller     bool
	DisableStacktrace bool
	DisableLoggerName bool
}

// DefaultEncoderConfig returns the default encoder settings used by xlog.
func DefaultEncoderConfig() EncoderConfig {
	return EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LevelFormat:    LevelFormatLowercase,
		CallerFormat:   CallerFormatShort,
		DurationFormat: DurationFormatMillis,
	}
}
