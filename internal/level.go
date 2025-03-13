package internal

import (
	"fmt"

	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
)

// LevelMap 日志级别映射
var LevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

// ColorFuncs 颜色输出函数映射
var ColorFuncs = map[zapcore.Level]func(format string, a ...interface{}) string{
	zapcore.DebugLevel: color.New(color.FgBlue).SprintfFunc(),
	zapcore.InfoLevel:  color.New(color.FgGreen).SprintfFunc(),
	zapcore.WarnLevel:  color.New(color.FgYellow).SprintfFunc(),
	zapcore.ErrorLevel: color.New(color.FgRed).SprintfFunc(),
	zapcore.FatalLevel: color.New(color.FgRed, color.Bold).SprintfFunc(),
}

// GetZapLevel 获取日志级别
func GetZapLevel(level string) zapcore.Level {
	if zapLevel, ok := LevelMap[level]; ok {
		return zapLevel
	}
	return zapcore.InfoLevel
}

// GetColorLevelEncoder 获取带颜色的级别编码器
func GetColorLevelEncoder() zapcore.LevelEncoder {
	return func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		colorFunc, ok := ColorFuncs[level]
		if !ok {
			enc.AppendString(fmt.Sprintf("%-5s", level.CapitalString()))
			return
		}
		enc.AppendString(colorFunc("%-5s", level.CapitalString()))
	}
}

// GetPlainLevelEncoder 获取普通的级别编码器
func GetPlainLevelEncoder() zapcore.LevelEncoder {
	return func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%-5s", level.CapitalString()))
	}
}
