package logger

import "go.uber.org/zap/zapcore"

// SyncFatalHook 实现 zapcore.CheckWriteHook，用于在 Fatal 级别日志写入前强制 Sync。
// 该钩子保证了进程退出前所有缓冲日志都被刷盘。
type SyncFatalHook struct {
	logger *Logger
}

// OnWrite 是 zapcore.CheckWriteHook 接口要求的方法。
// 当日志被写入后调用此方法，如果写入级别是 Fatal，则先 Sync logger，再调用默认的 Fatal 行为。
func (h SyncFatalHook) OnWrite(entry *zapcore.CheckedEntry, fields []zapcore.Field) {
	if entry.Level == zapcore.FatalLevel {
		if h.logger != nil {
			_ = h.logger.Sync()
		}
	}
	// 调用默认的 Fatal 处理
	zapcore.WriteThenFatal.OnWrite(entry, fields)
}
