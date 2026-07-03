package xlog

import "go.uber.org/zap"

// Field is a structured log field for the typed (zero-allocation) API.
type Field struct {
	zap zap.Field
}

func (f Field) toZap() zap.Field { return f.zap }

// String constructs a string field.
func String(key, val string) Field { return Field{zap: zap.String(key, val)} }

// Int constructs an int field.
func Int(key string, val int) Field { return Field{zap: zap.Int(key, val)} }

// Int64 constructs an int64 field.
func Int64(key string, val int64) Field { return Field{zap: zap.Int64(key, val)} }

// Bool constructs a bool field.
func Bool(key string, val bool) Field { return Field{zap: zap.Bool(key, val)} }

// Float64 constructs a float64 field.
func Float64(key string, val float64) Field { return Field{zap: zap.Float64(key, val)} }

// Any constructs a field with arbitrary value.
func Any(key string, val any) Field { return Field{zap: zap.Any(key, val)} }

// Err constructs an error field.
func Err(err error) Field { return Field{zap: zap.Error(err)} }
