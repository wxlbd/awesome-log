# Awesome Log

一个基于 [zap](https://github.com/uber-go/zap) 的高性能日志库，提供了简单易用的 API 和丰富的功能特性。

## 特性

- 🚀 **高性能**: 基于 uber 的 zap 日志库，提供极致的日志性能
- 🎨 **彩色输出**: 支持终端彩色日志输出，提升日志可读性
- 📁 **文件轮转**: 支持日志文件自动轮转，包括大小限制、时间限制和备份数量限制
- 🔄 **多格式支持**: 支持 JSON 和 Console 两种日志格式
- 🎯 **调用者信息**: 自动记录日志调用的文件名和行号
- 🔍 **堆栈跟踪**: Fatal 级别日志自动记录完整堆栈信息
- 🎪 **灵活配置**: 支持函数式选项模式进行配置
- 💉 **依赖注入**: 支持命名日志实例，便于服务注入使用

## 安装

```bash
go get github.com/wxlbd/awesome-log
```

## 快速开始

### 基本使用

```go
package main

import "github.com/wxlbd/awesome-log/pkg/logger"

func main() {
    // 使用默认配置初始化
    logger.Init()
    defer logger.Sync()

    // 使用日志
    logger.Debug("这是一条调试日志")
    logger.Info("这是一条信息日志")
    logger.Warn("这是一条警告日志")
    logger.Error("这是一条错误日志")
}
```

### 自定义配置

```go
package main

import "github.com/wxlbd/awesome-log/pkg/logger"

func main() {
    // 使用自定义选项初始化
    logger.Init(
        logger.WithLevel("debug"),
        logger.WithColor(true),
        logger.WithFileRotation(
            "logs/app.log",
            100,  // 100MB
            7,    // 7天
            10,   // 保留10个备份
            true, // 启用压缩
        ),
        logger.WithFileFormat("json"),
        logger.WithTimeFormat("2006-01-02 15:04:05.000"),
    )
    defer logger.Sync()

    // 使用日志
    logger.Info("日志系统初始化完成")
}
```

### 依赖注入使用

```go
package main

import "github.com/wxlbd/awesome-log/pkg/logger"

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
func (s *UserService) Login(username string) {
    s.log.Infof("用户 %s 尝试登录", username)
    // 业务逻辑...
    s.log.Info("登录成功")
}

func main() {
    logger.Init()
    defer logger.Sync()

    userService := NewUserService()
    userService.Login("admin")
}
```

## 配置选项

### 日志级别

支持的日志级别：
- debug
- info
- warn
- error
- fatal

```go
logger.WithLevel("debug")
```

### 日志格式

支持的日志格式：
- console (默认)
- json

```go
logger.WithFormat("json")
```

### 文件输出

```go
logger.WithFileRotation(
    "logs/app.log", // 文件路径
    100,           // 单个文件最大大小（MB）
    7,             // 文件保留天数
    10,            // 最大保留文件数
    true,          // 是否压缩
)
```

### 时间格式

```go
logger.WithTimeFormat("2006-01-02 15:04:05.000")
```

### 颜色输出

```go
logger.WithColor(true)
```

### 调用者信息

```go
logger.WithCaller(true)
```

## 日志输出示例

### 控制台格式
```
[2024-03-20 10:30:00.000] INFO    user-service/user.go:25 用户 admin 尝试登录
[2024-03-20 10:30:00.001] INFO    user-service/user.go:27 登录成功
```

### JSON 格式
```json
{
    "level": "info",
    "time": "2024-03-20 10:30:00.000",
    "caller": "user-service/user.go:25",
    "msg": "用户 admin 尝试登录",
    "logger": "user-service"
}
```

## 性能考虑

- 使用 `Sync()` 确保日志完全写入
- JSON 格式相比 Console 格式有更好的性能
- 适当配置日志级别，避免过多的 Debug 日志
- 合理使用日志轮转配置，避免日志文件过大

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License 