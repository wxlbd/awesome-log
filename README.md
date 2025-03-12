# Awesome Log

一个基于 [zap](https://github.com/uber-go/zap) 的高性能日志库，提供了简单易用的 API 和丰富的功能特性。

## 特性

- 🚀 **高性能**: 基于 zap 实现，提供极致的日志性能
- 🎨 **彩色输出**: 支持终端彩色输出，提高日志可读性
- 📝 **结构化日志**: 支持字段化的结构化日志记录
- 🔄 **日志轮转**: 支持基于大小、时间的日志文件轮转
- 📊 **多种格式**: 支持 JSON 和 Console 两种输出格式
- 🎯 **灵活配置**: 提供函数式选项模式的配置方式
- 🔍 **调用信息**: 自动记录日志调用的文件和行号
- 📚 **堆栈跟踪**: 支持可配置的错误堆栈跟踪
- 🌈 **命名日志器**: 支持创建多个命名日志实例
- 💉 **依赖注入**: 支持通过 WithName 方法创建命名日志实例

## 安装

```bash
go get github.com/wxlbd/awesome-log
```

## 快速开始

### 最佳实践

```go
package main

import (
    "github.com/wxlbd/awesome-log"
)

// UserService 用户服务
type UserService struct {
    log *logger.Logger
}

// OrderService 订单服务
type OrderService struct {
    log *logger.Logger
}

// NewUserService 创建用户服务实例
func NewUserService(baseLogger *logger.Logger) *UserService {
    return &UserService{
        log: baseLogger.WithName("user-service"),
    }
}

// NewOrderService 创建订单服务实例
func NewOrderService(baseLogger *logger.Logger) *OrderService {
    return &OrderService{
        log: baseLogger.WithName("order-service"),
    }
}

func main() {
    // 1. 在程序入口处初始化全局配置
    logger.Init(
        logger.WithLevel("debug"),
        logger.WithFileRotation("logs/app.log", 100, 7, 10, true),
        logger.WithColor(true),
        logger.WithStackLevel("error"),
    )
    defer logger.Sync()

    // 2. 创建基础 logger 实例
    baseLogger := logger.GetLogger("")

    // 3. 创建服务实例，注入命名 logger
    userService := NewUserService(baseLogger)
    orderService := NewOrderService(baseLogger)

    // 4. 使用服务
    userService.Login("12345")
    orderService.CreateOrder(1001, 99.99)
}

// Login 用户登录
func (s *UserService) Login(userID string) {
    s.log.Info("用户登录",
        zap.String("user_id", userID),
        zap.String("ip", "192.168.1.100"),
    )
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(orderID int, amount float64) {
    s.log.Info("创建订单",
        zap.Int("order_id", orderID),
        zap.Float64("amount", amount),
    )
}
```

提示：
1. 全局配置只需要在程序入口处初始化一次
2. 使用 `WithName` 方法从基础 logger 创建命名 logger
3. 每个命名 logger 会自动创建独立的日志文件（格式：原文件名.服务名.扩展名）
4. 所有命名 logger 共享相同的基础配置（级别、格式、轮转策略等）
5. 适合依赖注入场景，可以在服务初始化时注入 logger

### 基本使用

```go
package main

import (
    "github.com/wxlbd/awesome-log"
)

func main() {
    // 使用默认配置初始化
    logger.Init()
    defer logger.Sync()

    // 输出不同级别的日志
    logger.Debug("这是一条调试日志")
    logger.Info("这是一条信息日志")
    logger.Warn("这是一条警告日志")
    logger.Error("这是一条错误日志")
}
```

### 结构化日志

```go
logger.Info("用户登录",
    zap.String("username", "alice"),
    zap.String("ip", "192.168.1.100"),
    zap.Duration("latency", 50*time.Millisecond),
)

logger.Error("数据库错误",
    zap.String("operation", "insert"),
    zap.String("table", "users"),
    zap.Error(err),
)
```

### 自定义配置

```go
logger.Init(
    logger.WithLevel("debug"),                // 设置日志级别
    logger.WithColor(true),                   // 启用彩色输出
    logger.WithTimeFormat("2006-01-02 15:04:05.000"), // 自定义时间格式
    logger.WithCaller(true),                  // 记录调用者信息
    logger.WithStackLevel("error"),           // error 及以上级别输出堆栈
    logger.WithFileRotation(
        "logs/app.log", // 日志文件路径
        100,           // 单个文件最大大小（MB）
        7,             // 保留天数
        10,            // 保留文件数
        true,          // 启用压缩
    ),
    logger.WithFileFormat("json"), // 文件输出格式
)
```

### 命名日志器

```go
// 方式一：通过 GetLogger 获取
userLogger := logger.GetLogger("user-service")
orderLogger := logger.GetLogger("order-service")

// 方式二：通过 WithName 创建（推荐）
baseLogger := logger.GetLogger("")
userLogger := baseLogger.WithName("user-service")
orderLogger := baseLogger.WithName("order-service")

// 方式三：通过 NewLogger 创建
baseLogger := logger.NewLogger("")
userLogger := baseLogger.WithName("user-service")
orderLogger := baseLogger.WithName("order-service")

// 使用命名日志器
userLogger.Info("用户注册成功", 
    zap.String("username", "alice"),
    zap.String("email", "alice@example.com"),
)

orderLogger.Info("订单创建成功",
    zap.Int("order_id", 12345),
    zap.Float64("amount", 99.99),
)
```

## 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| WithLevel | 设置日志级别 (debug/info/warn/error/fatal) | "info" |
| WithColor | 启用/禁用彩色输出 | true |
| WithTimeFormat | 设置时间格式 | "2006-01-02 15:04:05.000" |
| WithCaller | 是否记录调用者信息 | true |
| WithStackLevel | 设置堆栈跟踪级别 | "fatal" |
| WithFileRotation | 配置日志文件轮转 | 未启用 |
| WithFileFormat | 设置文件输出格式 (json/console) | "json" |

## 日志格式示例

### 控制台输出
```
2024-03-12 15:04:05.000 INFO    user-service    main.go:28  用户登录成功  {"username": "alice", "ip": "192.168.1.100"}
2024-03-12 15:04:05.000 ERROR   order-service   main.go:35  订单创建失败  {"order_id": 12345, "error": "余额不足"}
```

### JSON 输出
```json
{
    "level": "INFO",
    "time": "2024-03-12 15:04:05.000",
    "logger": "user-service",
    "caller": "main.go:28",
    "msg": "用户登录成功",
    "username": "alice",
    "ip": "192.168.1.100"
}
```

## 示例

查看 [examples](./examples) 目录获取更多示例：

- [基础使用](./examples/basic/main.go)
- [结构化日志](./examples/structured/main.go)
- [日志轮转](./examples/rotation/main.go)
- [依赖注入](./examples/di/main.go)

## 性能

基于 zap 的高性能实现，在启用所有功能（调用者信息、时间格式化、日志轮转）的情况下：

- 结构化日志: ~2800ns/op
- 格式化日志: ~3200ns/op
- 并发写入: 支持每秒数百万条日志

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

[MIT License](LICENSE) 