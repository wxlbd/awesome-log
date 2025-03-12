package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	logger "github.com/wxlbd/awesome-log"
	"go.uber.org/zap"
)

// User 用户信息
type User struct {
	ID       int
	Username string
	Email    string
}

// Order 订单信息
type Order struct {
	ID        int
	UserID    int
	Amount    float64
	Status    string
	CreatedAt time.Time
}

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
			filepath.Join(currentDir, "logs", "structured.log"),
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

	// 模拟用户注册
	user := &User{
		ID:       1001,
		Username: "alice",
		Email:    "alice@example.com",
	}

	logger.Info("用户注册成功",
		zap.Int("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("email", user.Email),
		zap.String("ip", "192.168.1.100"),
		zap.String("user_agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"),
	)

	// 模拟订单创建
	order := &Order{
		ID:        2001,
		UserID:    user.ID,
		Amount:    99.99,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	logger.Info("订单创建成功",
		zap.Int("order_id", order.ID),
		zap.Int("user_id", order.UserID),
		zap.Float64("amount", order.Amount),
		zap.String("status", order.Status),
		zap.Time("created_at", order.CreatedAt),
	)

	// 模拟支付处理
	startTime := time.Now()
	time.Sleep(100 * time.Millisecond) // 模拟处理时间

	logger.Info("支付处理完成",
		zap.Int("order_id", order.ID),
		zap.String("payment_method", "alipay"),
		zap.String("transaction_id", "T20240312123456"),
		zap.Duration("process_time", time.Since(startTime)),
	)

	// 模拟错误处理
	if err := processRefund(order); err != nil {
		logger.Error("退款处理失败",
			zap.Int("order_id", order.ID),
			zap.Float64("amount", order.Amount),
			zap.Error(err),
			zap.String("refund_id", "R20240312123456"),
		)
	}

	// 等待一秒，确保日志写入
	time.Sleep(time.Second)
}

func processRefund(order *Order) error {
	return fmt.Errorf("退款金额超过订单金额：请求退款金额 %.2f，订单金额 %.2f", 199.99, order.Amount)
}
