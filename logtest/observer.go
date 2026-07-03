package logtest

import (
	"sync"
	"testing"

	"github.com/navms/go-log"
)

type recorder struct {
	mu      sync.Mutex
	entries []xlog.Entry
}

func (r *recorder) Handle(e xlog.Entry) error {
	r.mu.Lock()
	r.entries = append(r.entries, e)
	r.mu.Unlock()
	return nil
}

func (r *recorder) Entries() []xlog.Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]xlog.Entry, len(r.entries))
	copy(out, r.entries)
	return out
}

// Observer captures log entries via hooks for test assertions.
type Observer struct {
	rec *recorder
}

// Entries returns captured log entries.
func (o *Observer) Entries() []xlog.Entry { return o.rec.Entries() }

// NewObserver creates a logger that records entries to obs.
func NewObserver(level xlog.Level) (*Observer, xlog.Logger) {
	rec := &recorder{}
	l, err := xlog.New(
		xlog.WithLevel(level),
		xlog.WithHook(rec),
		xlog.WithConsole(),
	)
	if err != nil {
		panic(err)
	}
	return &Observer{rec: rec}, l
}

// RequireLogContains asserts that a log entry with msg and fields exists.
func RequireLogContains(t *testing.T, obs *Observer, msg string, fields map[string]any) {
	t.Helper()
	for _, e := range obs.Entries() {
		if e.Message != msg {
			continue
		}
		if fields == nil || fieldsMatch(e.Fields, fields) {
			return
		}
	}
	t.Fatalf("logtest: no entry with message %q and fields %v; got %v", msg, fields, obs.Entries())
}

func fieldsMatch(got map[string]any, want map[string]any) bool {
	for k, v := range want {
		gv, ok := got[k]
		if !ok || !valuesEqual(gv, v) {
			return false
		}
	}
	return true
}

func valuesEqual(a, b any) bool {
	switch av := a.(type) {
	case int64:
		switch bv := b.(type) {
		case int:
			return av == int64(bv)
		case int64:
			return av == bv
		}
	case int:
		switch bv := b.(type) {
		case int64:
			return int64(av) == bv
		case int:
			return av == bv
		}
	}
	return a == b
}
