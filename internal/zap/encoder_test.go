package zap

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCustomEncoderKeys(t *testing.T) {
	var buf bytes.Buffer
	ws := zapcore.AddSync(&buf)
	spec := &EncoderSpec{
		MessageKey: "msg",
		TimeKey:    "time",
		LevelKey:   "severity",
		TimeFormat: "2006-01-02",
		DisableCaller: true,
	}
	enc := newEncoder("json", false, spec)
	core := zapcore.NewCore(enc, ws, zapcore.InfoLevel)
	l := zap.New(core)
	l.Info("hello")

	line := strings.TrimSpace(buf.String())
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("json: %v line=%q", err, line)
	}
	if m["msg"] != "hello" {
		t.Fatalf("msg = %v", m["msg"])
	}
	if _, ok := m["message"]; ok {
		t.Fatal("default message key should not appear")
	}
	if _, ok := m["caller"]; ok {
		t.Fatal("caller should be omitted")
	}
	if m["severity"] != "info" {
		t.Fatalf("severity = %v", m["severity"])
	}
}

func TestDisableTime(t *testing.T) {
	var buf bytes.Buffer
	spec := &EncoderSpec{DisableTime: true, DisableCaller: true}
	enc := newEncoder("json", false, spec)
	core := zapcore.NewCore(enc, zapcore.AddSync(&buf), zapcore.InfoLevel)
	zap.New(core).Info("x")
	var m map[string]any
	_ = json.Unmarshal(buf.Bytes(), &m)
	if _, ok := m["timestamp"]; ok {
		t.Fatal("time should be omitted")
	}
}
