// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package init

import (
	"fmt"

	"github.com/vistone/netconnpool"
	"github.com/vistone/crawler-system/internal/config"
)

// InitNetConnPool 初始化TCP连接池模块（模块8）
func InitNetConnPool(cfg *config.NetConnPoolConfig, logger interface{ Info(string, ...interface{}) }) (*netconnpool.Pool, error) {
	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🌐 [模块8] TCP连接池模块初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ 最大连接数: %d\n", cfg.MaxConnections)
	fmt.Printf("  ✓ 初始连接数: %d\n", cfg.InitialConnections)
	fmt.Printf("  ✓ 获取连接超时: %d 秒\n", cfg.AcquireTimeout)
	fmt.Printf("  ✓ 空闲连接超时: %d 秒\n", cfg.IdleTimeout)
	fmt.Printf("  ✓ 连接最大生存时间: %d 秒\n", cfg.MaxLifetime)
	fmt.Printf("  ✓ 健康检查间隔: %d 秒\n", cfg.HealthCheckInterval)
	fmt.Printf("  ✓ 健康检查超时: %d 秒\n", cfg.HealthCheckTimeout)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// 注意：客户端模式需要Dialer，但Dialer需要目标地址
	// 由于在初始化时还不知道目标地址，这里暂时不创建连接池
	// 连接池将在实际使用时按需创建
	fmt.Printf("  ℹ️  说明: 连接池将在需要时按需创建（需要Dialer时）\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if logger != nil {
		logger.Info("TCP连接池配置已准备，max_connections=%d, initial_connections=%d, note=连接池将在需要时按需创建",
			cfg.MaxConnections, cfg.InitialConnections)
	}

	// 暂时返回nil，连接池延迟初始化
	return nil, nil
}

