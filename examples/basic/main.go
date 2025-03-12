package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	logger "github.com/wxlbd/awesome-log"
	"go.uber.org/zap"
)

func main() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		os.Exit(1)
	}

	// 使用自定义选项初始化
	if err := logger.Init(
		logger.WithLevel("debug"),
		logger.WithColor(true),
		logger.WithFileRotation(
			filepath.Join(currentDir, "logs", "app.log"),
			100,  // 100MB
			7,    // 7天
			10,   // 保留10个备份
			true, // 启用压缩
		),
		logger.WithFileFormat("json"),
		logger.WithTimeFormat("2006-01-02 15:04:05.000"),
	); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 基本日志示例
	logger.Debug("这是一条调试日志")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")

	// 格式化日志示例
	logger.Debugf("调试信息: %s", "详细内容")
	logger.Infof("用户 %s 登录系统", "admin")
	logger.Warnf("磁盘使用率达到 %d%%", 85)
	logger.Errorf("数据库连接失败: %v", "connection refused")

	// 结构化日志示例
	logger.Info("HTTP请求完成",
		zap.String("method", "GET"),
		zap.String("url", "http://example.com/api"),
		zap.Int("status", 200),
		zap.Duration("latency", 50*time.Millisecond),
	)

	logger.Error("数据库操作失败",
		zap.String("operation", "INSERT"),
		zap.String("table", "users"),
		zap.Error(fmt.Errorf("主键冲突")),
	)
	logger.Fatal("这是一条致命日志", zap.String("error", "致命错误"))
	// 等待一秒，确保日志写入
	time.Sleep(time.Second)
}
