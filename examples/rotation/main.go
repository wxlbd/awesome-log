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
			filepath.Join(currentDir, "logs", "rotation.log"),
			1,    // 1MB，方便测试
			1,    // 1天，方便测试
			3,    // 保留3个备份
			true, // 启用压缩
		),
		logger.WithFileFormat("json"),
		logger.WithTimeFormat("2006-01-02 15:04:05.000"),
	); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 生成大量日志以触发轮转
	for i := 0; i < 10000; i++ {
		logger.Info("测试日志轮转",
			zap.Int("index", i),
			zap.String("data", generateLargeString(100)), // 生成100字节的数据
			zap.Time("timestamp", time.Now()),
		)

		if i%1000 == 0 {
			// 每1000条日志暂停一下，方便观察
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("已写入 %d 条日志\n", i)
		}
	}

	// 等待一秒，确保日志写入
	time.Sleep(time.Second)

	// 列出日志文件
	files, err := filepath.Glob(filepath.Join(currentDir, "logs", "rotation*.log*"))
	if err != nil {
		fmt.Printf("列出日志文件失败: %v\n", err)
		return
	}

	fmt.Println("\n生成的日志文件:")
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		fmt.Printf("- %s (%.2f MB)\n", filepath.Base(file), float64(info.Size())/1024/1024)
	}
}

// generateLargeString 生成指定大小的字符串
func generateLargeString(size int) string {
	chars := make([]byte, size)
	for i := 0; i < size; i++ {
		chars[i] = 'a' + byte(i%26)
	}
	return string(chars)
}
