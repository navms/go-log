// 文件输出 + 轮转：同时写 stdout 和文件
package main

import (
	"os"
	"path/filepath"

	"github.com/navms/go-log"
)

func main() {
	logDir := filepath.Join(os.TempDir(), "xlog-example")
	_ = os.MkdirAll(logDir, 0o755)
	logFile := filepath.Join(logDir, "app.log")

	if err := xlog.Init(
		xlog.WithJSON(),
		xlog.WithOutputs(xlog.OutputStdout, xlog.OutputFile),
		xlog.WithFile(logFile, xlog.Rotation{
			MaxSize:    10, // MB
			MaxBackups: 3,
			MaxAge:     7, // days
			Compress:   true,
		}),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	xlog.Info("written to stdout and file", "path", logFile)
}
