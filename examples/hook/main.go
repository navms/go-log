// Hook：日志写入后触发自定义逻辑（如上报监控）
package main

import (
	"fmt"

	"github.com/navms/go-log"
)

func main() {
	if err := xlog.Init(
		xlog.WithConsole(),
		xlog.WithHook(xlog.HookFunc(func(e xlog.Entry) error {
			if e.Level >= xlog.ErrorLevel {
				fmt.Printf("[hook] alert level=%s msg=%q fields=%v\n", e.Level, e.Message, e.Fields)
			}
			return nil
		})),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	xlog.Info("normal log")
	xlog.Error("something failed", "code", 500)
}
