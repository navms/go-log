package xlog

// Option configures a logger.
type Option func(*Config)

func WithLevel(level Level) Option {
	return func(c *Config) { c.Level = level }
}

func WithJSON() Option {
	return func(c *Config) { c.Format = FormatJSON }
}

func WithConsole() Option {
	return func(c *Config) { c.Format = FormatConsole }
}

func WithDevelopment() Option {
	return func(c *Config) { c.Development = true }
}

func WithOutputs(outputs ...Output) Option {
	return func(c *Config) {
		c.Outputs = append([]Output(nil), outputs...)
	}
}

func WithFile(path string, rot Rotation) Option {
	return func(c *Config) {
		c.File.Path = path
		c.File.Rotation = rot
		hasFile := false
		for _, o := range c.Outputs {
			if o == OutputFile {
				hasFile = true
				break
			}
		}
		if !hasFile {
			c.Outputs = append(c.Outputs, OutputFile)
		}
	}
}

func WithRotation(rot Rotation) Option {
	return func(c *Config) { c.File.Rotation = rot }
}

func WithSyslog(tag string) Option {
	return func(c *Config) {
		c.SyslogTag = tag
		hasSyslog := false
		for _, o := range c.Outputs {
			if o == OutputSyslog {
				hasSyslog = true
				break
			}
		}
		if !hasSyslog {
			c.Outputs = append(c.Outputs, OutputSyslog)
		}
	}
}

func WithSyslogNetwork(network string) Option {
	return func(c *Config) { c.SyslogNetwork = network }
}

func WithHook(h Hook) Option {
	return func(c *Config) { c.Hooks = append(c.Hooks, h) }
}

func WithFields(keysAndValues ...any) Option {
	return func(c *Config) {
		if c.InitialFields == nil {
			c.InitialFields = make(map[string]any)
		}
		for i := 0; i+1 < len(keysAndValues); i += 2 {
			key, ok := keysAndValues[i].(string)
			if !ok {
				continue
			}
			c.InitialFields[key] = keysAndValues[i+1]
		}
	}
}

func WithContextKeys(keys ContextKeys) Option {
	return func(c *Config) { c.ContextKeys = keys }
}

// WithSampling enables sampling for trace/debug/info levels.
// Warn and above are always logged. Disabled when initial or thereafter is zero.
func WithSampling(initial, thereafter int) Option {
	return func(c *Config) {
		c.Sampling = Sampling{Initial: initial, Thereafter: thereafter}
	}
}

// WithSyncOnFatal controls whether buffers are synced before fatal exit.
func WithSyncOnFatal(enabled bool) Option {
	return func(c *Config) { c.SyncOnFatal = enabled }
}

// WithProductionPreset applies production defaults: JSON file output with platform log path.
func WithProductionPreset(appName string) Option {
	return func(c *Config) {
		preset := NewProductionConfig(appName)
		c.Level = preset.Level
		c.Format = preset.Format
		c.Outputs = append([]Output(nil), preset.Outputs...)
		c.File = preset.File
		c.Development = preset.Development
		c.SyncOnFatal = preset.SyncOnFatal
	}
}

// WithDevelopmentPreset applies development defaults: colored console output.
func WithDevelopmentPreset() Option {
	return func(c *Config) {
		preset := NewDevelopmentConfig()
		c.Level = preset.Level
		c.Format = preset.Format
		c.Outputs = append([]Output(nil), preset.Outputs...)
		c.Development = preset.Development
		c.SyncOnFatal = preset.SyncOnFatal
	}
}

// WithEncoderConfig sets custom encoder field keys and formatting.
// Unset fields keep defaults; use Disable* to omit fields from output.
func WithEncoderConfig(ec EncoderConfig) Option {
	return func(c *Config) {
		ecCopy := ec
		c.Encoder = &ecCopy
	}
}
