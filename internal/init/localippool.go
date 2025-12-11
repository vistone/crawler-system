// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package init

import (
	"fmt"

	"github.com/vistone/crawler-system/internal/config"
	"github.com/vistone/localippool"
)

// InitLocalIPPool 初始化本地IP池模块（模块4）
func InitLocalIPPool(cfg *config.LocalIPPoolConfig, logger interface{ Info(string, ...interface{}) }) (localippool.IPPool, error) {
	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🌍 [模块4] 本地IP池模块初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ IPv4地址数量: %d\n", len(cfg.IPs))
	if len(cfg.IPs) > 0 {
		fmt.Printf("  ✓ IPv4地址列表:\n")
		for _, ip := range cfg.IPs {
			fmt.Printf("    - %s\n", ip)
		}
	} else {
		fmt.Printf("  ✓ IPv4地址: 自动检测\n")
	}
	fmt.Printf("  ✓ 选择策略: %s\n", cfg.SelectionStrategy)
	fmt.Printf("  ✓ 健康检查: %v\n", cfg.HealthCheckEnabled)
	if cfg.HealthCheckEnabled {
		fmt.Printf("  ✓ 健康检查间隔: %d 秒\n", cfg.HealthCheckInterval)
		fmt.Printf("  ✓ 健康检查超时: %d 秒\n", cfg.HealthCheckTimeout)
		fmt.Printf("  ✓ 最大失败次数: %d\n", cfg.MaxFailures)
		fmt.Printf("  ✓ 恢复检查间隔: %d 秒\n", cfg.RecoveryCheckInterval)
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// 创建本地IP池（使用配置的IPv4列表，IPv6为空表示自动检测）
	pool, err := localippool.NewLocalIPPool(cfg.IPs, "")
	if err != nil {
		return nil, fmt.Errorf("创建本地IP池失败: %w", err)
	}

	// 获取实际检测到的IP信息
	ipv4s := pool.GetIPv4Addresses()
	ipv6s := pool.GetActiveIPv6Addresses()
	supportsDynamic := pool.SupportsDynamicPool()

	fmt.Printf("  📊 实际检测结果:\n")
	fmt.Printf("    - IPv4地址: %v\n", ipv4s)
	if supportsDynamic {
		fmt.Printf("    - IPv6动态池: 已启用\n")
		if len(ipv6s) > 0 {
			fmt.Printf("    - 活跃IPv6地址数量: %d\n", len(ipv6s))
			if len(ipv6s) <= 5 {
				fmt.Printf("    - IPv6地址列表: %v\n", ipv6s)
			} else {
				fmt.Printf("    - IPv6地址列表: %v ... (共%d个)\n", ipv6s[:5], len(ipv6s))
			}
		}
	} else {
		fmt.Printf("    - IPv6动态池: 未启用\n")
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if logger != nil {
		logger.Info("本地IP池模块初始化完成，ipv4_count=%d, strategy=%s", len(cfg.IPs), cfg.SelectionStrategy)
	}
	return pool, nil
}
