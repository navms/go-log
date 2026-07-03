package logtest_test

import (
	"testing"

	"github.com/navms/go-log"
	"github.com/navms/go-log/logtest"
)

func TestObserver(t *testing.T) {
	obs, l := logtest.NewObserver(xlog.InfoLevel)
	l.Info("ping", "x", 1)
	logtest.RequireLogContains(t, obs, "ping", map[string]any{"x": 1})
	logtest.AssertLevel(t, obs, xlog.InfoLevel)
}
