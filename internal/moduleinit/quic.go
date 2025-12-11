// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moduleinit

import (
	"fmt"
	"time"

	"github.com/vistone/quic"
	"github.com/vistone/crawler-system/internal/config"
)

// InitQUICPool 初始化QUIC连接池模块（模块9）
func InitQUICPool(cfg *config.QUICConfig, logger interface{ Info(string, ...interface{}) }) (*quic.Pool, error) {
	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("⚡ [模块9] QUIC连接池模块初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ 最大连接数: %d\n", cfg.MaxConnections)
	fmt.Printf("  ✓ 初始连接数: %d\n", cfg.InitialConnections)
	fmt.Printf("  ✓ 获取连接超时: %d 秒\n", cfg.AcquireTimeout)
	fmt.Printf("  ✓ 空闲连接超时: %d 秒\n", cfg.IdleTimeout)
	fmt.Printf("  ✓ 连接最大生存时间: %d 秒\n", cfg.MaxLifetime)
	fmt.Printf("  ✓ 健康检查间隔: %d 秒\n", cfg.HealthCheckInterval)
	fmt.Printf("  ✓ 健康检查超时: %d 秒\n", cfg.HealthCheckTimeout)
	fmt.Printf("  ✓ 握手超时: %d 秒\n", cfg.HandshakeTimeout)
	fmt.Printf("  ✓ 0-RTT支持: %v\n", cfg.Enable0RTT)
	if cfg.Enable0RTT {
		fmt.Printf("    ℹ️  说明: 0-RTT可以加速连接建立\n")
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// 创建QUIC客户端连接池
	minCap := cfg.InitialConnections
	if minCap < 1 {
		minCap = 1
	}
	maxCap := cfg.MaxConnections
	if maxCap < minCap {
		maxCap = minCap
	}

	fmt.Printf("  📊 连接池参数:\n")
	fmt.Printf("    - 最小容量: %d\n", minCap)
	fmt.Printf("    - 最大容量: %d\n", maxCap)
	fmt.Printf("    - 空闲超时: %d 秒\n", cfg.IdleTimeout)
	fmt.Printf("    - 最大生存时间: %d 秒\n", cfg.MaxLifetime)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	pool := quic.NewClientPool(
		minCap,
		maxCap,
		time.Duration(cfg.IdleTimeout)*time.Second,
		time.Duration(cfg.MaxLifetime)*time.Second,
		time.Duration(cfg.IdleTimeout)*time.Second,
		"", // tlsCode，后续从配置读取
		"", // hostname，后续从配置读取
		nil, // addrResolver，后续实现
	)

	if logger != nil {
		logger.Info("QUIC连接池模块初始化完成，max_connections=%d, enable_0rtt=%v", cfg.MaxConnections, cfg.Enable0RTT)
	}
	return pool, nil
}

