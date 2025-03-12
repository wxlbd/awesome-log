package internal

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
)

// CustomTimeEncoder 自定义时间编码器
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
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

// CustomCallerEncoder 自定义调用者编码器
func CustomCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
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

// GetConsoleEncoder 获取控制台编码器配置
func GetConsoleEncoder(enableColor bool, timeFormat string) zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   CustomCallerEncoder,
	}

	if enableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	return encoderConfig
}

// GetFileEncoder 获取文件编码器配置
func GetFileEncoder(timeFormat string) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
