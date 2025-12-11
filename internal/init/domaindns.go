// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package init

import (
	"fmt"

	"github.com/vistone/domaindns"
	"github.com/vistone/crawler-system/internal/config"
)

// InitDomainDNS 初始化DNS解析模块（模块3）
func InitDomainDNS(cfg *config.DomainDNSConfig, targetDomains []string, logger interface{ Info(string, ...interface{}); Warn(string, ...interface{}) }) (domaindns.DomainMonitor, error) {
	if len(targetDomains) == 0 {
		if logger != nil {
			logger.Warn("未配置目标域名，跳过DNS监控器创建")
		}
		return nil, nil
	}

	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🌐 [模块3] DNS解析模块初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ DNS服务器数量: %d\n", len(cfg.DNSServers))
	if len(cfg.DNSServers) > 0 {
		fmt.Printf("  ✓ DNS服务器列表:\n")
		for i, server := range cfg.DNSServers {
			if i < 10 {
				fmt.Printf("    - %s\n", server)
			} else if i == 10 {
				fmt.Printf("    ... 还有 %d 个服务器\n", len(cfg.DNSServers)-10)
				break
			}
		}
	}
	fmt.Printf("  ✓ DNS缓存: %v\n", cfg.CacheEnabled)
	if cfg.CacheEnabled {
		fmt.Printf("  ✓ 缓存TTL: %d 秒\n", cfg.CacheTTL)
	}
	fmt.Printf("  ✓ 查询超时: %d 秒\n", cfg.Timeout)
	fmt.Printf("  ✓ 最大重试: %d 次\n", cfg.MaxRetries)
	fmt.Printf("  ✓ 重试间隔: %d 秒\n", cfg.RetryInterval)
	fmt.Printf("  ✓ DNS污染检测: %v\n", cfg.PollutionDetection)
	fmt.Printf("  ✓ IPv6支持: %v\n", cfg.IPv6Enabled)
	fmt.Printf("  ✓ IPInfo Token: %s\n", getDisplayValue(cfg.IPInfoToken, "未配置"))
	fmt.Printf("  ✓ 目标域名数量: %d\n", len(targetDomains))
	fmt.Printf("  ✓ 目标域名列表:\n")
	for _, domain := range targetDomains {
		fmt.Printf("    - %s\n", domain)
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// 使用NewMonitorWithGlobalDNSServers创建监控器
	ipInfoToken := cfg.IPInfoToken
	if ipInfoToken == "" && logger != nil {
		logger.Warn("未配置IPInfo Token，IP详细信息获取功能将不可用")
	}

	monitor, err := domaindns.NewMonitorWithGlobalDNSServers(
		targetDomains,
		ipInfoToken,
		"", // dnsServerFile，空表示使用默认的dnsservernames.json
		0,  // maxServers，0表示使用全部DNS服务器
	)
	if err != nil {
		return nil, fmt.Errorf("创建DNS监控器失败: %w", err)
	}

	// 启动DomainMonitor
	monitor.Start()

	if logger != nil {
		logger.Info("DNS解析模块初始化完成，dns_servers=%d, target_domains=%d, domains=%v",
			len(cfg.DNSServers), len(targetDomains), targetDomains)
	}

	return monitor, nil
}


