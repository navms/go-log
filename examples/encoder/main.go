// 自定义格式：字段名、时间格式、隐藏 caller
package main

import "github.com/navms/go-log"

func main() {
	if err := xlog.Init(
		xlog.WithJSON(),
		xlog.WithEncoderConfig(xlog.EncoderConfig{
			TimeKey:       "time",
			MessageKey:    "msg",
			LevelKey:      "severity",
			TimeFormat:    "2006-01-02 15:04:05.000",
			DisableCaller: true,
		}),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	// 输出: {"time":"...","severity":"info","msg":"custom format","k":"v"}
	xlog.Info("custom format", "k", "v")
}
