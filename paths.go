package xlog

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DefaultLogFilePath returns a platform-appropriate default log file path.
//
//	macOS:   ~/Library/Logs/{appName}/app.log
//	Windows: %LOCALAPPDATA%/{appName}/logs/app.log
//	Linux:   $XDG_STATE_HOME/{appName}/logs/app.log or ~/.local/state/...
func DefaultLogFilePath(appName string) string {
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, "Library", "Logs", appName, "app.log")
		}
	case "windows":
		base := os.Getenv("LOCALAPPDATA")
		if strings.TrimSpace(base) == "" {
			base, _ = os.UserConfigDir()
		}
		if strings.TrimSpace(base) != "" {
			return filepath.Join(base, appName, "logs", "app.log")
		}
	default:
		base := os.Getenv("XDG_STATE_HOME")
		if strings.TrimSpace(base) == "" {
			home, err := os.UserHomeDir()
			if err == nil {
				base = filepath.Join(home, ".local", "state")
			}
		}
		if strings.TrimSpace(base) != "" {
			return filepath.Join(base, appName, "logs", "app.log")
		}
	}
	return filepath.Join(os.TempDir(), appName, "logs", "app.log")
}

// NewDevelopmentConfig returns a console logger preset for local development.
func NewDevelopmentConfig() Config {
	cfg := DefaultConfig()
	cfg.Level = DebugLevel
	cfg.Format = FormatConsole
	cfg.Outputs = []Output{OutputStdout}
	cfg.Development = true
	cfg.SyncOnFatal = true
	return cfg
}

// NewProductionConfig returns a JSON file logger preset for production.
func NewProductionConfig(appName string) Config {
	cfg := DefaultConfig()
	cfg.Level = InfoLevel
	cfg.Format = FormatJSON
	cfg.Outputs = []Output{OutputFile}
	cfg.Development = false
	cfg.SyncOnFatal = true
	cfg.File = FileConfig{
		Path: DefaultLogFilePath(appName),
		Rotation: Rotation{
			MaxSize:    128,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		},
	}
	return cfg
}
