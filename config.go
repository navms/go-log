package logger

// Config 日志配置
type Config struct {
	Env                Env    `mapstructure:"env" yaml:"env"`                                 // 运行环境
	Level              Level  `mapstructure:"level" yaml:"level"`                             // 日志级别
	Output             Output `mapstructure:"output" yaml:"output"`                           // 输出方式
	RedirectStdLog     bool   `mapstructure:"redirect-std-log" yaml:"redirect-std-log"`       // 是否将标准库 log 重定向到 zap
	Filename           string `mapstructure:"filename" yaml:"filename"`                       // 日志文件路径
	MaxSize            int    `mapstructure:"max-size" yaml:"max-size"`                       // 单个日志文件最大尺寸(MB)
	MaxBackups         int    `mapstructure:"max-backups" yaml:"max-backups"`                 // 保留旧文件最大数量
	MaxAge             int    `mapstructure:"max-age" yaml:"max-age"`                         // 保留旧文件最大天数
	Compress           bool   `mapstructure:"compress" yaml:"compress"`                       // 是否压缩旧文件
	ShowCaller         bool   `mapstructure:"show-caller" yaml:"show-caller"`                 // 是否显示调用位置
	SamplingInitial    int    `mapstructure:"sampling-initial" yaml:"sampling-initial"`       // 采样：每秒前N条全记录
	SamplingThereafter int    `mapstructure:"sampling-thereafter" yaml:"sampling-thereafter"` // 采样：之后每M条记录1条
}

func NewDevelopConfig() *Config {
	return &Config{
		Env:            Dev,
		Level:          DebugL,
		Output:         Console,
		ShowCaller:     true,
		RedirectStdLog: false,
	}
}

func NewProductConfig() *Config {
	return &Config{
		Env:            Prod,
		Level:          InfoL,
		Output:         File,
		Filename:       "./logs/app.log",
		MaxSize:        128,
		MaxBackups:     3,
		MaxAge:         7,
		Compress:       true,
		ShowCaller:     true,
		RedirectStdLog: false,
	}
}
