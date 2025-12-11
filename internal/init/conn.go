// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package init

import (
	"fmt"

	"github.com/vistone/crawler-system/internal/config"
)

// ConnManager 连接管理器
type ConnManager struct {
	Config *config.ConnConfig
}

// InitConn 初始化连接模块（模块7）
func InitConn(cfg *config.ConnConfig, logger interface{ Info(string, ...interface{}) }) (*ConnManager, error) {
	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🔌 [模块7] 连接模块初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ 连接超时: %d 秒\n", cfg.ConnectTimeout)
	fmt.Printf("  ✓ 读取超时: %d 秒\n", cfg.ReadTimeout)
	fmt.Printf("  ✓ 写入超时: %d 秒\n", cfg.WriteTimeout)
	fmt.Printf("  ✓ Keep-Alive: %v\n", cfg.KeepAlive)
	if cfg.KeepAlive {
		fmt.Printf("  ✓ Keep-Alive时间: %d 秒\n", cfg.KeepAliveTime)
	}
	fmt.Printf("  ✓ 最大空闲连接数: %d\n", cfg.MaxIdleConns)
	fmt.Printf("  ✓ 每个主机最大连接数: %d\n", cfg.MaxConnsPerHost)
	fmt.Printf("  ✓ TLS握手超时: %d 秒\n", cfg.TLSHandshakeTimeout)
	fmt.Printf("  ✓ 跳过TLS验证: %v\n", cfg.InsecureSkipVerify)
	if cfg.InsecureSkipVerify {
		fmt.Printf("    ⚠️  警告: TLS证书验证已禁用（仅用于测试）\n")
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	cm := &ConnManager{
		Config: cfg,
	}

	if logger != nil {
		logger.Info("连接模块初始化完成，connect_timeout=%d, read_timeout=%d", cfg.ConnectTimeout, cfg.ReadTimeout)
	}
	return cm, nil
}
