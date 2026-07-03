package zap

import "go.uber.org/zap/zapcore"

// syncFatalHook flushes buffers before the process exits on fatal logs.
type syncFatalHook struct {
	syncFn func() error
}

func (h syncFatalHook) OnWrite(ce *zapcore.CheckedEntry, fields []zapcore.Field) {
	if ce.Entry.Level == zapcore.FatalLevel && h.syncFn != nil {
		_ = h.syncFn()
	}
	zapcore.WriteThenFatal.OnWrite(ce, fields)
}
