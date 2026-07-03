# go-log

Structured logging facade for Go, built on [uber-go/zap](https://github.com/uber-go/zap).

- **Module:** `github.com/navms/go-log`
- **Package:** `xlog`

## Quick start

```go
import "github.com/navms/go-log"

func main() {
    _ = xlog.Init(
        xlog.WithJSON(),
        xlog.WithLevel(xlog.InfoLevel),
    )
    xlog.Info("server started", "port", 8080)
    defer xlog.Sync()
}
```

## Presets

```go
// Development: colored console, debug level
_ = xlog.Init(xlog.WithDevelopmentPreset())

// Production: JSON file with platform default path
_ = xlog.Init(xlog.WithProductionPreset("myapp"))

// Or build from config directly
l, _ := xlog.NewWithConfig(xlog.NewProductionConfig("myapp"))
```

`DefaultLogFilePath("myapp")` picks a platform-specific path:

- macOS: `~/Library/Logs/myapp/app.log`
- Windows: `%LOCALAPPDATA%/myapp/logs/app.log`
- Linux: `~/.local/state/myapp/logs/app.log`

## Sampling

High-volume trace/debug/info logs can be sampled; warn and above are always kept.

```go
xlog.Init(xlog.WithSampling(100, 100)) // first 100/s full rate, then 1/100
```

Environment: `LOG_SAMPLING_INITIAL`, `LOG_SAMPLING_THEREAFTER`.

## Fatal sync

By default, buffers are synced before `Fatal` exits (`SyncOnFatal: true`). Disable with:

```go
xlog.Init(xlog.WithSyncOnFatal(false))
```

## Custom encoder format

Use `WithEncoderConfig` to customize JSON/console field keys and encodings:

```go
xlog.Init(
    xlog.WithJSON(),
    xlog.WithEncoderConfig(xlog.EncoderConfig{
        TimeKey:    "time",
        MessageKey: "msg",
        LevelKey:   "severity",
        TimeFormat: "2006-01-02 15:04:05.000",
        DisableCaller: true,
    }),
)
```

| Field | Description |
|-------|-------------|
| `TimeKey` / `LevelKey` / `MessageKey` / ... | JSON/console field names |
| `TimeFormat` | Go time layout; empty = ISO8601 |
| `LevelFormat` | `lowercase`, `capital`, `capitalColor` |
| `CallerFormat` | `short`, `full` |
| `DurationFormat` | `millis`, `seconds`, `string` |
| `DisableTime` / `DisableCaller` / ... | Omit fields from output |

Defaults are available via `xlog.DefaultEncoderConfig()`.

## Configuration

### Functional options

| Option | Description |
|--------|-------------|
| `WithLevel` | Minimum log level |
| `WithJSON` / `WithConsole` | Output format |
| `WithOutputs` | stdout, stderr, file, syslog |
| `WithFile` | File path + rotation settings |
| `WithSyslog` | Enable syslog output |
| `WithHook` | Register a hook |
| `WithFields` | Initial structured fields |
| `WithDevelopment` | Development mode (colored console) |
| `WithSampling` | Sample trace/debug/info logs |
| `WithSyncOnFatal` | Sync buffers before fatal exit |
| `WithProductionPreset` | JSON file + platform log path |
| `WithDevelopmentPreset` | Colored console for local dev |
| `WithEncoderConfig` | Custom field keys and time/level encoding |

### Environment variables

| Variable | Description |
|----------|-------------|
| `LOG_LEVEL` | trace, debug, info, warn, error, fatal |
| `LOG_FORMAT` | json or console |
| `LOG_OUTPUT` | comma-separated: stdout, stderr, file, syslog |
| `LOG_FILE_PATH` | log file path |
| `LOG_FILE_MAX_SIZE` | rotation max size (MB) |
| `LOG_FILE_MAX_BACKUPS` | max backup files |
| `LOG_FILE_MAX_AGE` | max age (days) |
| `LOG_FILE_COMPRESS` | compress rotated files |
| `LOG_DEVELOPMENT` | true/false |
| `LOG_SYSLOG_TAG` | syslog tag |

```go
_ = xlog.InitFromEnv()
```

## Context trace fields

```go
ctx := xlog.WithTraceID(ctx, traceID)
ctx = xlog.WithSpanID(ctx, spanID)
xlog.WithContext(ctx).Info("handled")
```

## Child loggers

```go
userLog := xlog.With("module", "user")
userLog.Info("login", "user_id", uid)
```

## Typed API (hot paths)

```go
xlog.Typed().Info("order", xlog.String("id", id), xlog.Int("amount", 100))
```

## Dynamic log level

```go
http.Handle("/debug/loglevel", xlog.LevelHandler())
xlog.SetLevel(xlog.DebugLevel)
```

## Hooks

```go
xlog.Init(xlog.WithHook(xlog.HookFunc(func(e xlog.Entry) error {
    // e.g. send errors to Sentry
    return nil
})))
```

## Testing

```go
import "github.com/navms/go-log/logtest"

obs, logger := logtest.NewObserver(xlog.InfoLevel)
logger.Info("ping", "x", 1)
logtest.RequireLogContains(t, obs, "ping", map[string]any{"x": 1})
```

## Stdlib compatibility

```go
std := xlog.Default().StdLogger()
std.Println("via stdlib")
```

## Examples

```bash
go run ./examples/basic/          # 开发环境预设
go run ./examples/production/     # 生产 JSON + 文件
go run ./examples/json/           # JSON 输出
go run ./examples/context/        # trace_id / span_id
go run ./examples/child_logger/   # With / Named 子 logger
go run ./examples/typed/          # 零分配 Typed API
go run ./examples/encoder/        # 自定义字段格式
go run ./examples/file/           # 文件 + 轮转
go run ./examples/env/            # 环境变量配置
go run ./examples/sampling/       # 日志采样
go run ./examples/hook/           # Hook 回调
go run ./examples/level/          # 动态级别 + HTTP
go run ./examples/noop/           # Noop logger
```

## License

MIT
