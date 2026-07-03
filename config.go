package xlog

// Format selects log output encoding.
type Format string

const (
	FormatConsole Format = "console"
	FormatJSON    Format = "json"
)

// Output selects log sink targets.
type Output string

const (
	OutputStdout Output = "stdout"
	OutputStderr Output = "stderr"
	OutputFile   Output = "file"
	OutputSyslog Output = "syslog"
)

// Rotation configures file log rotation (lumberjack).
type Rotation struct {
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

// FileConfig holds file output settings.
type FileConfig struct {
	Path     string
	Rotation Rotation
}

// ContextKeys configures field names extracted from context.Context.
type ContextKeys struct {
	TraceID string
	SpanID  string
}

// DefaultContextKeys returns the default context field names.
func DefaultContextKeys() ContextKeys {
	return ContextKeys{TraceID: "trace_id", SpanID: "span_id"}
}

// Sampling configures log sampling for verbose levels (trace/debug/info).
// Warn and above are always logged without sampling.
// Sampling is disabled when Initial or Thereafter is zero.
type Sampling struct {
	Initial    int // log first N entries per tick at full rate
	Thereafter int // then log 1 of every M entries
}

// Config holds logger configuration.
type Config struct {
	Level         Level
	Format        Format
	Outputs       []Output
	File          FileConfig
	SyslogNetwork string // e.g. "udp", "tcp"; empty uses default
	SyslogTag     string
	Development   bool
	Sampling      Sampling
	SyncOnFatal   bool // sync buffers before fatal exit; default true
	Encoder       *EncoderConfig
	InitialFields map[string]any
	ContextKeys   ContextKeys
	Hooks         []Hook
}

// DefaultConfig returns sensible defaults (info, console, stdout).
func DefaultConfig() Config {
	return Config{
		Level:       InfoLevel,
		Format:      FormatConsole,
		Outputs:     []Output{OutputStdout},
		SyncOnFatal: true,
		ContextKeys: DefaultContextKeys(),
		File: FileConfig{
			Rotation: Rotation{MaxSize: 100, MaxBackups: 7, MaxAge: 30, Compress: true},
		},
		Encoder: &EncoderConfig{
			TimeFormat: "2006-01-02 15:04:05",
		},
		SyslogTag: "xlog",
	}
}
