package xlog_test

import (
	"testing"

	"github.com/navms/go-log"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in   string
		want xlog.Level
	}{
		{"trace", xlog.TraceLevel},
		{"debug", xlog.DebugLevel},
		{"info", xlog.InfoLevel},
		{"warn", xlog.WarnLevel},
		{"error", xlog.ErrorLevel},
		{"fatal", xlog.FatalLevel},
	}
	for _, tt := range tests {
		got, err := xlog.ParseLevel(tt.in)
		if err != nil {
			t.Fatalf("ParseLevel(%q): %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("ParseLevel(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSetLevel(t *testing.T) {
	l, err := xlog.New(xlog.WithLevel(xlog.InfoLevel), xlog.WithConsole())
	if err != nil {
		t.Fatal(err)
	}
	l.SetLevel(xlog.DebugLevel)
	if l.Level() != xlog.DebugLevel {
		t.Fatalf("level = %v, want debug", l.Level())
	}
}
