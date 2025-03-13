package main

import (
	"errors"
	"fmt"
	"time"

	logger "github.com/wxlbd/awesome-log"
	"go.uber.org/zap"
)

// User 用户模型
type User struct {
	ID       string
	Username string
	Email    string
}

// Order 订单模型
type Order struct {
	ID        int
	UserID    string
	Amount    float64
	Status    string
	CreatedAt time.Time
}

// Payment 支付模型
type Payment struct {
	ID        string
	OrderID   int
	Amount    float64
	Status    string
	CreatedAt time.Time
}

// UserService 用户服务
type UserService struct {
	log *logger.Logger
}

// OrderService 订单服务
type OrderService struct {
	log         *logger.Logger
	userService *UserService
}

// PaymentService 支付服务
type PaymentService struct {
	log          *logger.Logger
	orderService *OrderService
}

// NewUserService 创建用户服务
func NewUserService(baseLogger *logger.Logger) *UserService {
	return &UserService{
		log: baseLogger.WithName("user-service"),
	}
}

// NewOrderService 创建订单服务
func NewOrderService(baseLogger *logger.Logger, userService *UserService) *OrderService {
	return &OrderService{
		log:         baseLogger.WithName("order-service"),
		userService: userService,
	}
}

// NewPaymentService 创建支付服务
func NewPaymentService(baseLogger *logger.Logger, orderService *OrderService) *PaymentService {
	return &PaymentService{
		log:          baseLogger.WithName("payment-service"),
		orderService: orderService,
	}
}

// Register 用户注册
func (s *UserService) Register(username, email string) (*User, error) {
	s.log.Info("开始用户注册",
		zap.String("username", username),
		zap.String("email", email),
	)

	// 模拟用户注册逻辑
	user := &User{
		ID:       fmt.Sprintf("u_%d", time.Now().UnixNano()),
		Username: username,
		Email:    email,
	}

	s.log.Info("用户注册成功",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
	)

	return user, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(userID string) (*User, error) {
	s.log.Debug("查询用户信息", zap.String("user_id", userID))

	// 模拟数据库查询
	if userID == "" {
		s.log.Error("用户ID不能为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 模拟找到用户
	user := &User{
		ID:       userID,
		Username: "测试用户",
		Email:    "test@example.com",
	}

	s.log.Debug("查询用户成功",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
	)

	return user, nil
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID string, amount float64) (*Order, error) {
	s.log.Info("开始创建订单",
		zap.String("user_id", userID),
		zap.Float64("amount", amount),
	)

	// 验证用户
	user, err := s.userService.GetUser(userID)
	if err != nil {
		s.log.Error("创建订单失败：用户不存在",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 创建订单
	order := &Order{
		ID:        int(time.Now().UnixNano() % 1000000),
		UserID:    user.ID,
		Amount:    amount,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	s.log.Info("订单创建成功",
		zap.Int("order_id", order.ID),
		zap.String("user_id", order.UserID),
		zap.Float64("amount", order.Amount),
	)

	return order, nil
}

// ProcessPayment 处理支付
func (s *PaymentService) ProcessPayment(orderID int, amount float64) (*Payment, error) {
	s.log.Info("开始处理支付",
		zap.Int("order_id", orderID),
		zap.Float64("amount", amount),
	)

	// 模拟支付处理
	payment := &Payment{
		ID:        fmt.Sprintf("p_%d", time.Now().UnixNano()),
		OrderID:   orderID,
		Amount:    amount,
		Status:    "success",
		CreatedAt: time.Now(),
	}

	// 模拟随机支付失败
	if time.Now().UnixNano()%2 == 0 {
		s.log.Warn("支付处理失败",
			zap.Int("order_id", orderID),
			zap.Float64("amount", amount),
			zap.String("reason", "余额不足"),
		)
		payment.Status = "failed"
		return payment, errors.New("支付失败：余额不足")
	}

	s.log.Info("支付处理成功",
		zap.String("payment_id", payment.ID),
		zap.Int("order_id", payment.OrderID),
		zap.Float64("amount", payment.Amount),
	)

	return payment, nil
}

func main() {
	// 1. 初始化日志配置
	logger.Init(
		logger.WithLevel("debug"),
		logger.WithColor(true),
		logger.WithCaller(true),
		logger.WithStackLevel("error"),
		logger.WithFileRotation(
			"logs/app.log",
			100,
			7,
			10,
			true,
		),
		logger.WithFileFormat("json"),
	)
	defer logger.Sync()

	// 2. 创建基础 logger
	baseLogger := logger.GetLogger("")

	// 3. 创建服务实例（依赖注入）
	userService := NewUserService(baseLogger)
	orderService := NewOrderService(baseLogger, userService)
	paymentService := NewPaymentService(baseLogger, orderService)

	// 4. 模拟业务流程
	// 4.1 用户注册
	user, err := userService.Register("alice", "alice@example.com")
	if err != nil {
		panic(err)
	}

	// 4.2 创建订单
	order, err := orderService.CreateOrder(user.ID, 99.99)
	if err != nil {
		panic(err)
	}

	// 4.3 处理支付
	payment, err := paymentService.ProcessPayment(order.ID, order.Amount)
	if err != nil {
		fmt.Printf("支付失败: %v\n", err)
	} else {
		fmt.Printf("支付成功: %s\n", payment.ID)
	}
}
