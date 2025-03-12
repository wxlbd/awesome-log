package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 封装的日志结构
type Logger struct {
	config *Config
	zap    *zap.Logger
	sugar  *zap.SugaredLogger
}

var (
	// 全局日志实例
	globalLogger *Logger
	// 日志实例映射，用于管理多个命名日志实例
	loggerMap   = make(map[string]*Logger)
	loggerMutex sync.RWMutex
	// 日志级别映射
	levelMap = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}
	// 颜色输出函数映射
	colorFuncs = map[zapcore.Level]func(format string, a ...interface{}) string{
		zapcore.DebugLevel: color.New(color.FgBlue).SprintfFunc(),
		zapcore.InfoLevel:  color.New(color.FgGreen).SprintfFunc(),
		zapcore.WarnLevel:  color.New(color.FgYellow).SprintfFunc(),
		zapcore.ErrorLevel: color.New(color.FgRed).SprintfFunc(),
		zapcore.FatalLevel: color.New(color.FgRed, color.Bold).SprintfFunc(),
	}
)

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(color.New(color.FgWhite, color.Bold).Sprintf(
		"%d-%02d-%02d %02d:%02d:%02d.%03d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond()/1e6,
	))
}

// customCallerEncoder 自定义调用者编码器
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// 获取调用文件的绝对路径
	path := caller.File
	// 获取项目根目录
	rootDir := filepath.Dir(filepath.Dir(path))
	// 将绝对路径转换为相对于项目根目录的路径
	if rel, err := filepath.Rel(rootDir, path); err == nil {
		path = rel
	}
	// 格式化输出：文件名:行号
	enc.AppendString(fmt.Sprintf("%s:%d", path, caller.Line))
}

// NewLogger 创建一个新的命名日志实例
func NewLogger(name string, opts ...Option) *Logger {
	loggerMutex.Lock()
	// 双重检查，确保在获取锁的过程中没有其他goroutine创建了logger
	if logger, exists := loggerMap[name]; exists {
		loggerMutex.Unlock()
		return logger
	}

	// 使用默认配置
	config := DefaultConfig()

	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	// 打印配置信息
	fmt.Printf("Logger配置: WriteToFile=%v, Filename=%s\n", config.WriteToFile, config.FileConfig.Filename)

	logger := &Logger{
		config: config,
	}

	// 创建控制台编码器配置
	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    logger.getEncodeLevel(),
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   customCallerEncoder,
	}

	// 创建文件编码器配置（无颜色）
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(config.TimeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// 控制台输出
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		getZapLevel(config.Level),
	)
	cores = append(cores, consoleCore)

	// 文件输出
	if config.WriteToFile {
		fmt.Printf("正在配置文件输出...\n")

		// 为每个命名logger创建独立的日志文件
		filename := config.FileConfig.Filename
		if name != "" {
			ext := filepath.Ext(filename)
			filename = filename[:len(filename)-len(ext)] + "." + name + ext
		}

		// 确保日志目录存在
		logDir := filepath.Dir(filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
			loggerMutex.Unlock()
			return nil
		}

		// 尝试创建日志文件
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("创建日志文件失败: %v\n", err)
			loggerMutex.Unlock()
			return nil
		}
		file.Close()

		// 创建日志写入器
		writer := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    config.FileConfig.MaxSize,
			MaxBackups: config.FileConfig.MaxBackups,
			MaxAge:     config.FileConfig.MaxAge,
			Compress:   config.FileConfig.Compress,
			LocalTime:  true,
		}

		fmt.Printf("创建日志写入器: %+v\n", writer)

		fileCore := zapcore.NewCore(
			getEncoder(config.FileConfig.Format, fileEncoderConfig),
			zapcore.AddSync(writer),
			getZapLevel(config.Level),
		)
		cores = append(cores, fileCore)
		fmt.Printf("文件输出配置完成\n")
	} else {
		fmt.Printf("未启用文件输出\n")
	}

	// 创建Logger
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core).Named(name)
	if config.RecordCaller {
		zapLogger = zapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	// 添加堆栈跟踪
	zapLogger = zapLogger.WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	logger.zap = zapLogger
	logger.sugar = zapLogger.Sugar()

	// 保存到映射中
	loggerMap[name] = logger
	loggerMutex.Unlock()
	return logger
}

// GetLogger 获取指定名称的日志实例
func GetLogger(name string) *Logger {
	loggerMutex.RLock()
	if logger, exists := loggerMap[name]; exists {
		loggerMutex.RUnlock()
		return logger
	}
	loggerMutex.RUnlock()

	// 如果全局logger已经初始化，使用其配置
	var opts []Option
	if globalLogger != nil {
		opts = append(opts, WithFullConfig(globalLogger.config))
	}
	return NewLogger(name, opts...)
}

// Init 初始化全局日志实例
func Init(opts ...Option) error {
	// 使用默认配置
	config := DefaultConfig()

	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	// 创建全局logger
	logger := NewLogger("", WithFullConfig(config))
	if logger == nil {
		return fmt.Errorf("初始化日志失败")
	}
	globalLogger = logger
	return nil
}

// getEncoder 获取编码器
func getEncoder(format string, config zapcore.EncoderConfig) zapcore.Encoder {
	if format == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)
}

// getZapLevel 获取日志级别
func getZapLevel(level string) zapcore.Level {
	if zapLevel, ok := levelMap[level]; ok {
		return zapLevel
	}
	return zapcore.InfoLevel
}

// getEncodeLevel 获取级别编码器
func (l *Logger) getEncodeLevel() zapcore.LevelEncoder {
	if l.config.EnableColor {
		return func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			colorFunc, ok := colorFuncs[level]
			if !ok {
				colorFunc = fmt.Sprintf
			}
			enc.AppendString(colorFunc("%-7s", level.CapitalString()))
		}
	}
	return func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%-7s", level.CapitalString()))
	}
}

// Debug 输出Debug级别日志
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debugf 输出Debug级别日志（格式化）
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Info 输出Info级别日志
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Infof 输出Info级别日志（格式化）
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warn 输出Warn级别日志
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Warnf 输出Warn级别日志（格式化）
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Error 输出Error级别日志
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Errorf 输出Error级别日志（格式化）
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Fatal 输出Fatal级别日志
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Fatalf 输出Fatal级别日志（格式化）
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

// Sync 同步缓存的日志
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// 以下是全局函数，使用全局logger实例

// Debug 输出Debug级别日志
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

// Debugf 输出Debug级别日志（格式化）
func Debugf(template string, args ...interface{}) {
	globalLogger.Debugf(template, args...)
}

// Info 输出Info级别日志
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

// Infof 输出Info级别日志（格式化）
func Infof(template string, args ...interface{}) {
	globalLogger.Infof(template, args...)
}

// Warn 输出Warn级别日志
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

// Warnf 输出Warn级别日志（格式化）
func Warnf(template string, args ...interface{}) {
	globalLogger.Warnf(template, args...)
}

// Error 输出Error级别日志
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

// Errorf 输出Error级别日志（格式化）
func Errorf(template string, args ...interface{}) {
	globalLogger.Errorf(template, args...)
}

// Fatal 输出Fatal级别日志
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Fatalf 输出Fatal级别日志（格式化）
func Fatalf(template string, args ...interface{}) {
	globalLogger.Fatalf(template, args...)
}

// Sync 同步缓存的日志
func Sync() error {
	return globalLogger.Sync()
}
