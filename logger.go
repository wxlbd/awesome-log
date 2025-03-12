package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wxlbd/awesome-log/internal"
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
)

// NewLogger 创建一个新的命名日志实例
func NewLogger(name string, opts ...Option) *Logger {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// 双重检查，确保在获取锁的过程中没有其他goroutine创建了logger
	if logger, exists := loggerMap[name]; exists {
		return logger
	}

	// 使用默认配置
	config := DefaultConfig()

	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	logger := &Logger{
		config: config,
	}

	var cores []zapcore.Core

	// 控制台输出
	consoleEncoderConfig := internal.GetConsoleEncoder(config.EnableColor, config.TimeFormat)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stdout),
		internal.GetZapLevel(config.Level),
	)
	cores = append(cores, consoleCore)

	// 文件输出
	if config.WriteToFile {
		// 为每个命名logger创建独立的日志文件
		filename := config.FileConfig.Filename
		if name != "" {
			ext := filepath.Ext(filename)
			filename = filename[:len(filename)-len(ext)] + "." + name + ext
		}

		// 确保日志目录存在
		logDir := filepath.Dir(filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil
		}

		// 创建日志写入器
		writer := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    config.FileConfig.MaxSize,
			MaxBackups: config.FileConfig.MaxBackups,
			MaxAge:     config.FileConfig.MaxAge,
			Compress:   config.FileConfig.Compress,
			LocalTime:  true,
		}

		fileEncoderConfig := internal.GetFileEncoder(config.TimeFormat)
		var fileEncoder zapcore.Encoder
		if config.FileConfig.Format == "json" {
			fileEncoder = zapcore.NewJSONEncoder(fileEncoderConfig)
		} else {
			fileEncoder = zapcore.NewConsoleEncoder(fileEncoderConfig)
		}

		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(writer),
			internal.GetZapLevel(config.Level),
		)
		cores = append(cores, fileCore)
	}

	// 创建Logger
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core).Named(name)
	if config.RecordCaller {
		zapLogger = zapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}

	// 设置堆栈跟踪级别
	stackLevel := internal.GetZapLevel(config.StackLevel)
	zapLogger = zapLogger.WithOptions(zap.AddStacktrace(stackLevel))

	logger.zap = zapLogger
	logger.sugar = zapLogger.Sugar()

	// 保存到映射中
	loggerMap[name] = logger
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

// Debug 输出Debug级别日志
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

// Info 输出Info级别日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

// Warn 输出Warn级别日志
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

// Error 输出Error级别日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

// Fatal 输出Fatal级别日志
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// Debugf 输出Debug级别日志（格式化）
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Infof 输出Info级别日志（格式化）
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warnf 输出Warn级别日志（格式化）
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Errorf 输出Error级别日志（格式化）
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
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
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info 输出Info级别日志
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn 输出Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error 输出Error级别日志
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal 输出Fatal级别日志
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// Debugf 输出Debug级别日志（格式化）
func Debugf(template string, args ...interface{}) {
	globalLogger.Debugf(template, args...)
}

// Infof 输出Info级别日志（格式化）
func Infof(template string, args ...interface{}) {
	globalLogger.Infof(template, args...)
}

// Warnf 输出Warn级别日志（格式化）
func Warnf(template string, args ...interface{}) {
	globalLogger.Warnf(template, args...)
}

// Errorf 输出Error级别日志（格式化）
func Errorf(template string, args ...interface{}) {
	globalLogger.Errorf(template, args...)
}

// Fatalf 输出Fatal级别日志（格式化）
func Fatalf(template string, args ...interface{}) {
	globalLogger.Fatalf(template, args...)
}

// Sync 同步缓存的日志
func Sync() error {
	return globalLogger.Sync()
}
