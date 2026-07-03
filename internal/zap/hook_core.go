package zap

import (
	"go.uber.org/zap/zapcore"
)

type hookCore struct {
	zapcore.Core
	hooks     []HookFunc
	ctxFields []zapcore.Field
}

func newHookCore(core zapcore.Core, hooks []HookFunc) zapcore.Core {
	if len(hooks) == 0 {
		return core
	}
	return &hookCore{Core: core, hooks: hooks}
}

func (c *hookCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Core.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *hookCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	err := c.Core.Write(ent, fields)
	all := append(append([]zapcore.Field{}, c.ctxFields...), fields...)
	fieldMap := fieldsToMap(all)
	for _, h := range c.hooks {
		_ = h(HookEntry{
			Level:   int8(ent.Level),
			Time:    ent.Time.UnixNano(),
			Message: ent.Message,
			Fields:  fieldMap,
			Caller:  ent.Caller.String(),
		})
	}
	return err
}

func (c *hookCore) With(fields []zapcore.Field) zapcore.Core {
	return &hookCore{
		Core:      c.Core.With(fields),
		hooks:     c.hooks,
		ctxFields: append(append([]zapcore.Field{}, c.ctxFields...), fields...),
	}
}

func fieldsToMap(fields []zapcore.Field) map[string]any {
	if len(fields) == 0 {
		return nil
	}
	m := make(map[string]any, len(fields))
	for _, f := range fields {
		m[f.Key] = fieldValue(f)
	}
	return m
}

func fieldValue(f zapcore.Field) any {
	if f.Interface != nil {
		return f.Interface
	}
	switch f.Type {
	case zapcore.StringType:
		return f.String
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return f.Integer
	case zapcore.BoolType:
		return f.Integer == 1
	case zapcore.ErrorType:
		if f.Interface != nil {
			return f.Interface
		}
		return nil
	default:
		if f.Interface != nil {
			return f.Interface
		}
		return f.String
	}
}
