package zap

// HookEntry is the internal hook payload (no dependency on xlog).
type HookEntry struct {
	Level   int8
	Time    int64 // unix nano
	Message string
	Fields  map[string]any
	Caller  string
}

// HookFunc is a hook callback used by the internal core wrapper.
type HookFunc func(entry HookEntry) error
