// 环境变量配置：配合 InitFromEnv 使用
//
//	LOG_LEVEL=debug LOG_FORMAT=json go run .
package main

import "github.com/navms/go-log"

func main() {
	if err := xlog.InitFromEnv(
		xlog.WithFields("service", "demo"),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	xlog.Debug("visible when LOG_LEVEL=debug")
	xlog.Info("configured from environment")
}
