package xlog

import (
	"os"
	"strconv"
	"strings"
)

// ConfigFromEnv builds a Config from environment variables, overlaid on defaults.
//
//	LOG_LEVEL=debug
//	LOG_FORMAT=json|console
//	LOG_OUTPUT=stdout,file,syslog  (comma-separated)
//	LOG_FILE_PATH=/var/log/app.log
//	LOG_FILE_MAX_SIZE=100
//	LOG_FILE_MAX_BACKUPS=7
//	LOG_FILE_MAX_AGE=30
//	LOG_FILE_COMPRESS=true
//	LOG_DEVELOPMENT=true
//	LOG_SYSLOG_TAG=myapp
//	LOG_SAMPLING_INITIAL=100
//	LOG_SAMPLING_THEREAFTER=100
func ConfigFromEnv() Config {
	cfg := DefaultConfig()

	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if lvl, err := ParseLevel(v); err == nil {
			cfg.Level = lvl
		}
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		switch strings.ToLower(v) {
		case "json":
			cfg.Format = FormatJSON
		case "console":
			cfg.Format = FormatConsole
		}
	}
	if v := os.Getenv("LOG_OUTPUT"); v != "" {
		parts := strings.Split(v, ",")
		cfg.Outputs = cfg.Outputs[:0]
		for _, p := range parts {
			switch strings.TrimSpace(strings.ToLower(p)) {
			case "stdout":
				cfg.Outputs = append(cfg.Outputs, OutputStdout)
			case "stderr":
				cfg.Outputs = append(cfg.Outputs, OutputStderr)
			case "file":
				cfg.Outputs = append(cfg.Outputs, OutputFile)
			case "syslog":
				cfg.Outputs = append(cfg.Outputs, OutputSyslog)
			}
		}
		if len(cfg.Outputs) == 0 {
			cfg.Outputs = []Output{OutputStdout}
		}
	}
	if v := os.Getenv("LOG_FILE_PATH"); v != "" {
		cfg.File.Path = v
	}
	if v := os.Getenv("LOG_FILE_MAX_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.File.Rotation.MaxSize = n
		}
	}
	if v := os.Getenv("LOG_FILE_MAX_BACKUPS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.File.Rotation.MaxBackups = n
		}
	}
	if v := os.Getenv("LOG_FILE_MAX_AGE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.File.Rotation.MaxAge = n
		}
	}
	if v := os.Getenv("LOG_FILE_COMPRESS"); v != "" {
		cfg.File.Rotation.Compress = strings.EqualFold(v, "true") || v == "1"
	}
	if v := os.Getenv("LOG_DEVELOPMENT"); v != "" {
		cfg.Development = strings.EqualFold(v, "true") || v == "1"
	}
	if v := os.Getenv("LOG_SYSLOG_TAG"); v != "" {
		cfg.SyslogTag = v
	}
	if v := os.Getenv("LOG_SAMPLING_INITIAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Sampling.Initial = n
		}
	}
	if v := os.Getenv("LOG_SAMPLING_THEREAFTER"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Sampling.Thereafter = n
		}
	}

	return cfg
}
