// 基础用法：开发环境预设（彩色 Console + Debug 级别）
package main

import "github.com/navms/go-log"

func main() {
	if err := xlog.Init(xlog.WithDevelopmentPreset()); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	xlog.Info("server started", "port", 8080)
	xlog.Debug("debug details", "version", "1.0.0")
}
