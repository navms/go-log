// 子 Logger：公共字段继承 + 命名空间
package main

import "github.com/navms/go-log"

func main() {
	if err := xlog.Init(
		xlog.WithConsole(),
		xlog.WithFields("service", "user-api", "env", "dev"),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	// With：附加固定字段，子 logger 每条日志都会带上
	userLog := xlog.With("module", "user")
	userLog.Info("login success", "user_id", 1001)

	// Named：为 logger 添加命名空间（对应 zap Named）
	httpLog := xlog.With("component", "http").Named("http")
	httpLog.Info("listening", "port", 8080)
}
