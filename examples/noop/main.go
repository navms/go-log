// Noop Logger：测试或关闭日志输出
package main

import "github.com/navms/go-log"

func main() {
	xlog.SetDefault(xlog.NewNoop())

	xlog.Info("this is discarded")
	xlog.Error("this too")

	// 恢复为正常 logger
	if err := xlog.Init(xlog.WithDevelopmentPreset()); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	xlog.Info("back to normal logging")
}
