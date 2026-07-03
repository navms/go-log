// 采样：trace/debug/info 采样，warn 及以上全量
package main

import "github.com/navms/go-log"

func main() {
	if err := xlog.Init(
		xlog.WithJSON(),
		xlog.WithLevel(xlog.DebugLevel),
		xlog.WithSampling(10, 10), // 每秒前 10 条全记，之后 1/10
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	for i := 0; i < 100; i++ {
		xlog.Debug("high volume", "i", i)
	}
	xlog.Warn("warn always logged", "i", 100)
}
