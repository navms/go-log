package xlog_test

import (
	"testing"

	"github.com/navms/go-log"
	"github.com/navms/go-log/logtest"
)

func TestStdLogger(t *testing.T) {
	obs, l := logtest.NewObserver(xlog.InfoLevel)
	sl := l.StdLogger()
	sl.Println("stdlib line")
	logtest.RequireLogContains(t, obs, "stdlib line\n", nil)
}

func TestNoop(t *testing.T) {
	l := xlog.NewNoop()
	l.Info("ignored")
	if l.Sync() != nil {
		t.Fatal("sync")
	}
	_ = l.StdLogger()
	_ = l.Typed().Zap()
}

func TestTypedLogger(t *testing.T) {
	obs, l := logtest.NewObserver(xlog.InfoLevel)
	l.Typed().Info("typed", xlog.String("k", "v"))
	logtest.RequireLogContains(t, obs, "typed", map[string]any{"k": "v"})
}

func TestHook(t *testing.T) {
	var called bool
	l, err := xlog.New(
		xlog.WithHook(xlog.HookFunc(func(e xlog.Entry) error {
			called = true
			return nil
		})),
		xlog.WithConsole(),
	)
	if err != nil {
		t.Fatal(err)
	}
	l.Info("hook test")
	if !called {
		t.Fatal("hook not called")
	}
}
