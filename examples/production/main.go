// 生产环境：JSON 格式 + 平台默认日志路径 + 文件轮转
package main

import (
	"os"

	"github.com/navms/go-log"
)

func main() {
	appName := "myapp"
	if v := os.Getenv("APP_NAME"); v != "" {
		appName = v
	}

	if err := xlog.Init(
		xlog.WithProductionPreset(appName),
		xlog.WithFields("service", appName, "env", "prod"),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	xlog.Info("production logger ready")
	xlog.Warn("disk usage high", "percent", 85)
}
