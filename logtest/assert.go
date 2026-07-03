package logtest

import (
	"testing"

	"github.com/navms/go-log"
)

// AssertLevel asserts the last captured entry has the expected level.
func AssertLevel(t *testing.T, obs *Observer, level xlog.Level) {
	t.Helper()
	entries := obs.Entries()
	if len(entries) == 0 {
		t.Fatal("logtest: no entries captured")
	}
	last := entries[len(entries)-1]
	if last.Level != level {
		t.Fatalf("logtest: level = %v, want %v", last.Level, level)
	}
}
