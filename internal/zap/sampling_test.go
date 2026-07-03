package zap

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestSamplingCore(t *testing.T) {
	b, err := New(BuildConfig{
		Level:  zapcore.DebugLevel,
		Format: "json",
		Output: OutputSpec{Stdout: true},
		Sampling: SamplingSpec{
			Initial:    2,
			Thereafter: 100,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		b.Sugar().Debug("sampled")
	}
	b.Sugar().Warn("unsampled")
}

func TestSyncFatalHookSyncsBeforeExit(t *testing.T) {
	var synced bool
	h := syncFatalHook{syncFn: func() error {
		synced = true
		return nil
	}}
	if h.syncFn != nil {
		_ = h.syncFn()
	}
	if !synced {
		t.Fatal("expected sync")
	}
}

func TestBuildWithSyncOnFatal(t *testing.T) {
	b, err := New(BuildConfig{
		Level:       zapcore.InfoLevel,
		Format:      "json",
		Output:      OutputSpec{Stdout: true},
		SyncOnFatal: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if b.Zap() == nil {
		t.Fatal("nil logger")
	}
}
