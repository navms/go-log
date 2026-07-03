package zap

import (
	"fmt"
	"io"
	"log/syslog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap/zapcore"
)

// OutputSpec describes a single output target for the builder.
type OutputSpec struct {
	Stdout bool
	Stderr bool
	File   *FileSpec
	Syslog *SyslogSpec
}

// FileSpec holds file output configuration.
type FileSpec struct {
	Path     string
	MaxSize  int
	MaxBackups int
	MaxAge   int
	Compress bool
}

// SyslogSpec holds syslog output configuration.
type SyslogSpec struct {
	Tag     string
	Network string
}

func buildWriteSyncer(spec OutputSpec) (zapcore.WriteSyncer, error) {
	var writers []io.Writer

	if spec.Stdout {
		writers = append(writers, os.Stdout)
	}
	if spec.Stderr {
		writers = append(writers, os.Stderr)
	}
	if spec.File != nil {
		if spec.File.Path == "" {
			return nil, fmt.Errorf("zap backend: file output requires a path")
		}
		writers = append(writers, &lumberjack.Logger{
			Filename:   spec.File.Path,
			MaxSize:    spec.File.MaxSize,
			MaxBackups: spec.File.MaxBackups,
			MaxAge:     spec.File.MaxAge,
			Compress:   spec.File.Compress,
		})
	}
	if spec.Syslog != nil {
		tag := spec.Syslog.Tag
		if tag == "" {
			tag = "xlog"
		}
		w, err := syslogDial(spec.Syslog.Network, tag)
		if err != nil {
			return nil, fmt.Errorf("zap backend: syslog: %w", err)
		}
		writers = append(writers, w)
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	if len(writers) == 1 {
		return zapcore.AddSync(writers[0]), nil
	}

	syncers := make([]zapcore.WriteSyncer, len(writers))
	for i, w := range writers {
		syncers[i] = zapcore.AddSync(w)
	}
	return zapcore.NewMultiWriteSyncer(syncers...), nil
}

var syslogDial = func(network, tag string) (io.Writer, error) {
	if network == "" {
		return syslog.New(syslog.LOG_INFO, tag)
	}
	return syslog.Dial(network, "", syslog.LOG_INFO, tag)
}

// SetSyslogDial replaces syslog dial for testing.
func SetSyslogDial(fn func(network, tag string) (io.Writer, error)) {
	syslogDial = fn
}
