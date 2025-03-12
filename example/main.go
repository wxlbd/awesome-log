package main

import (
	"fmt"
	logger "github.com/wxlbd/awesome-log"
	"os"
)

// UserService 用户服务
type UserService struct {
	log *logger.Logger
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		log: logger.GetLogger("user-service"),
	}
}

// Login 用户登录
func (s *UserService) Login(username, password string) error {
	s.log.Infof("用户 %s 尝试登录", username)
	// 登录逻辑...
	s.log.Info("登录成功")
	return nil
}

// OrderService 订单服务
type OrderService struct {
	log *logger.Logger
}

// NewOrderService 创建订单服务实例
func NewOrderService() *OrderService {
	return &OrderService{
		log: logger.GetLogger("order-service"),
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userId string, amount float64) error {
	s.log.Infof("用户 %s 创建订单，金额: %.2f", userId, amount)
	// 创建订单逻辑...
	s.log.Info("订单创建成功")
	return nil
}

func main() {
	// 使用自定义选项初始化
	if err := logger.Init(
		logger.WithLevel("debug"),
		logger.WithColor(true),
		logger.WithFileRotation(
			"logs/app.log", // 使用绝对路径
			100,            // 100MB
			7,              // 7天
			10,             // 保留10个备份
			true,           // 启用压缩
		),
		logger.WithFileFormat("json"),
		logger.WithTimeFormat("2006-01-02 15:04:05.000"),
	); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 创建服务实例
	userService := NewUserService()
	orderService := NewOrderService()

	// 使用服务
	userService.Login("admin", "password")
	orderService.CreateOrder("admin", 99.99)

	// 全局日志使用示例
	logger.Debug("这是一条全局调试日志")
	logger.Info("这是一条全局信息日志")
	logger.Warn("这是一条全局警告日志")
	logger.Error("这是一条全局错误日志")

	// 格式化日志示例
	logger.Debugf("调试信息: %s", "详细内容")
	logger.Infof("用户 %s 登录系统", "admin")
	logger.Warnf("磁盘使用率达到 %d%%", 85)
	logger.Errorf("数据库连接失败: %v", "connection refused")

	// Fatal级别日志会输出堆栈信息并终止程序
	logger.Fatal("发生致命错误，程序即将退出")
}
