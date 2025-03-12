package logger

// Option 定义日志选项函数类型
type Option func(*Config)

// Config 日志配置结构
type Config struct {
	// 日志级别: debug, info, warn, error, fatal
	Level string `json:"level" yaml:"level"`
	// 日志格式: json, console
	Format string `json:"format" yaml:"format"`
	// 是否输出到文件
	WriteToFile bool `json:"writeToFile" yaml:"writeToFile"`
	// 是否启用彩色输出（仅在控制台格式下有效）
	EnableColor bool `json:"enableColor" yaml:"enableColor"`
	// 日志文件配置
	FileConfig FileConfig `json:"fileConfig" yaml:"fileConfig"`
	// 是否记录调用者信息
	RecordCaller bool `json:"recordCaller" yaml:"recordCaller"`
	// 时间格式
	TimeFormat string `json:"timeFormat" yaml:"timeFormat"`
}

// FileConfig 文件输出配置
type FileConfig struct {
	// 日志文件路径
	Filename string `json:"filename" yaml:"filename"`
	// 单个日志文件最大大小（MB）
	MaxSize int `json:"maxSize" yaml:"maxSize"`
	// 最大保留天数
	MaxAge int `json:"maxAge" yaml:"maxAge"`
	// 最大保留文件数
	MaxBackups int `json:"maxBackups" yaml:"maxBackups"`
	// 是否压缩
	Compress bool `json:"compress" yaml:"compress"`
	// 日志格式: json, console
	Format string `json:"format" yaml:"format"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:       "info",
		Format:      "console",
		WriteToFile: false,
		EnableColor: true,
		FileConfig: FileConfig{
			Filename:   "logs/app.log",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 10,
			Compress:   true,
			Format:     "json",
		},
		RecordCaller: true,
		TimeFormat:   "2006-01-02 15:04:05.000",
	}
}

// WithLevel 设置日志级别
func WithLevel(level string) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithFormat 设置日志格式
func WithFormat(format string) Option {
	return func(c *Config) {
		c.Format = format
	}
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(format string) Option {
	return func(c *Config) {
		c.TimeFormat = format
	}
}

// WithColor 设置是否启用彩色输出
func WithColor(enable bool) Option {
	return func(c *Config) {
		c.EnableColor = enable
	}
}

// WithCaller 设置是否记录调用者信息
func WithCaller(enable bool) Option {
	return func(c *Config) {
		c.RecordCaller = enable
	}
}

// WithFileRotation 启用文件输出并设置文件轮转配置
func WithFileRotation(filename string, maxSize, maxAge, maxBackups int, compress bool) Option {
	return func(c *Config) {
		c.WriteToFile = true
		c.FileConfig = FileConfig{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
			Compress:   compress,
			Format:     "json",
		}
	}
}

// WithFileFormat 设置文件输出格式
func WithFileFormat(format string) Option {
	return func(c *Config) {
		c.FileConfig.Format = format
	}
}

// WithFullConfig 使用完整配置
func WithFullConfig(config *Config) Option {
	return func(c *Config) {
		*c = *config
	}
}
