package xlog_test

import (
	"testing"

	"github.com/navms/go-log"
)

func TestEncoderConfigDefaults(t *testing.T) {
	d := xlog.DefaultEncoderConfig()
	if d.TimeKey != "timestamp" || d.MessageKey != "message" {
		t.Fatalf("defaults = %+v", d)
	}
}

func TestCustomEncoderViaNewWithConfig(t *testing.T) {
	c := xlog.DefaultConfig()
	c.Format = xlog.FormatJSON
	c.Encoder = &xlog.EncoderConfig{
		MessageKey:    "msg",
		TimeKey:         "time",
		LevelKey:        "severity",
		TimeFormat:      "2006-01-02",
		DisableCaller:   true,
	}
	l, err := xlog.NewWithConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	l.Info("formatted")
}

func TestEncoderConfigLevelFormat(t *testing.T) {
	formats := []xlog.LevelFormat{
		xlog.LevelFormatLowercase,
		xlog.LevelFormatCapital,
		xlog.LevelFormatCapitalColor,
	}
	for _, f := range formats {
		l, err := xlog.New(
			xlog.WithConsole(),
			xlog.WithEncoderConfig(xlog.EncoderConfig{LevelFormat: f}),
		)
		if err != nil {
			t.Fatalf("format %s: %v", f, err)
		}
		l.Info("level format test")
	}
}

func TestWithEncoderConfigOption(t *testing.T) {
	l, err := xlog.New(
		xlog.WithJSON(),
		xlog.WithEncoderConfig(xlog.EncoderConfig{
			MessageKey: "msg",
			TimeFormat: "2006-01-02 15:04:05.000",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	l.Info("custom keys")
}
