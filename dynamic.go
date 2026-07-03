package xlog

import "net/http"

// DynamicLevelHandler LevelHandler is an alias for the global level HTTP handler.
func DynamicLevelHandler() http.Handler { return LevelHandler() }
