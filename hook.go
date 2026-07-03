package xlog

import "time"

// Entry is passed to hooks after a log line is written.
type Entry struct {
	Level   Level
	Time    time.Time
	Message string
	Fields  map[string]any
	Caller  string
}

// Hook runs custom logic when a log entry is written.
// Errors from Handle do not prevent the log from being written.
type Hook interface {
	Handle(entry Entry) error
}

// HookFunc adapts a function to Hook.
type HookFunc func(entry Entry) error

func (f HookFunc) Handle(entry Entry) error { return f(entry) }
