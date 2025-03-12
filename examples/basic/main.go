package main

import (
	"fmt"
	"time"

	log "github.com/wxlbd/awesome-log"
	"go.uber.org/zap"
)

func main() {
	logger := log.NewLogger("")
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
	// 等待一秒，确保日志写入
	time.Sleep(time.Second)
	logger1 := logger.WithName("new name")
	fmt.Println("logger1:", logger, logger1)

	logger1.Info("这是一个新的日志实例")
	logger.Fatal("这是一条致命日志", zap.String("error", "致命错误"))

}
