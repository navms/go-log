// Typed API：热路径零分配，适合高频日志点
package main

import "github.com/navms/go-log"

func main() {
	if err := xlog.Init(xlog.WithJSON()); err != nil {
		panic(err)
	}
	defer xlog.Sync()

	typed := xlog.Typed()
	typed.Info("order created",
		xlog.String("order_id", "ORD-001"),
		xlog.Int("amount", 100),
		xlog.Bool("paid", true),
	)

	// 子 typed logger
	orderLog := typed.With(xlog.String("module", "order"))
	orderLog.Debug("inventory checked", xlog.Int("stock", 50))
}
