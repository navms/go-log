// 动态级别：运行时调整 + HTTP 端点
//
// 启动后访问: curl -X PUT localhost:8080/debug/loglevel -d '{"level":"debug"}'
package main

import (
	"fmt"
	"net/http"

	"github.com/navms/go-log"
)

func main() {
	if err := xlog.Init(
		xlog.WithConsole(),
		xlog.WithLevel(xlog.InfoLevel),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	http.Handle("/debug/loglevel", xlog.LevelHandler())

	xlog.Info("only info and above")
	xlog.SetLevel(xlog.DebugLevel)
	xlog.Debug("now debug is visible")

	fmt.Println("level handler on :8080/debug/loglevel")
	_ = http.ListenAndServe(":8080", nil)
}
