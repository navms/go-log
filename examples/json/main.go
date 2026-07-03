// JSON 输出：适合容器 / K8s / ELK 收集
package main

import "github.com/navms/go-log"

func main() {
	if err := xlog.Init(
		xlog.WithJSON(),
		xlog.WithLevel(xlog.InfoLevel),
		xlog.WithOutputs(xlog.OutputStdout),
	); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	xlog.Info("order created", "order_id", "ORD-001", "amount", 99.9)
}
