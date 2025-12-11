// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package init

import (
	"fmt"

	"github.com/vistone/crawler-system/internal/config"
	"github.com/vistone/logs"
)

// InitLogs 初始化日志系统（模块1）
func InitLogs(cfg *config.LogsConfig) (*logs.Logger, error) {
	// 转换日志级别
	var logLevel logs.LogLevel
	switch cfg.Level {
	case "debug":
		logLevel = logs.Debug
	case "info":
		logLevel = logs.Info
	case "warn":
		logLevel = logs.Warn
	case "error":
		logLevel = logs.Error
	default:
		logLevel = logs.Info
	}

	// 创建日志器（启用彩色输出）
	logger := logs.NewLogger(logLevel, true)

	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📝 [模块1] 日志系统初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ 日志级别: %s (LogLevel=%d)\n", cfg.Level, logLevel)
	fmt.Printf("  ✓ 文件输出: %v\n", cfg.FileEnabled)
	if cfg.FileEnabled {
		fmt.Printf("  ✓ 文件路径: %s\n", cfg.FilePath)
		fmt.Printf("  ✓ 最大大小: %d MB\n", cfg.MaxSize)
		fmt.Printf("  ✓ 保留数量: %d\n", cfg.MaxBackups)
		fmt.Printf("  ✓ 压缩: %v\n", cfg.Compress)
	}
	fmt.Printf("  ✓ 日志格式: %s\n", cfg.Format)
	fmt.Printf("  ✓ 显示调用位置: %v\n", cfg.ShowCaller)
	fmt.Printf("  ✓ 彩色输出: 已启用\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	logger.Info("日志系统初始化完成，level=%s", cfg.Level)
	return logger, nil
}
