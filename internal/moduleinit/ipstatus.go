// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moduleinit

import (
	"fmt"

	"github.com/vistone/crawler-system/internal/config"
)

// IPStatusManager 黑白名单管理器接口
type IPStatusManager interface {
	AddToWhitelist(ip string) error
	RemoveFromWhitelist(ip string, reason string) error
	AddToBlacklist(ip string, reason string) error
	GetStatus(ip string) string
	GetWhitelistIPs() []string
	GetWhitelistCount() int
	CheckSystemHealth() error
	SetMinWhitelistCount(count int)
	SetAllowStartWhenEmpty(allow bool)
	SetWhitelistMonitoring(enabled bool)
	SetWhitelistMonitoringInterval(interval int)
}

// PlaceholderIPStatusManager 占位符黑白名单管理器实现
type PlaceholderIPStatusManager struct {
	whitelist                    map[string]bool
	blacklist                    map[string]string
	minWhitelistCount            int
	allowStartWhenEmpty          bool
	whitelistMonitoring          bool
	whitelistMonitoringInterval  int
}

// InitIPStatusManager 初始化黑白名单模块（模块6）
func InitIPStatusManager(cfg *config.IPStatusConfig, logger interface{ Info(string, ...interface{}) }) (IPStatusManager, error) {
	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📋 [模块6] 黑白名单模块初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ 白名单最小数量: %d\n", cfg.MinWhitelistCount)
	fmt.Printf("  ✓ 允许空白名单启动: %v\n", cfg.AllowStartWhenEmpty)
	if !cfg.AllowStartWhenEmpty {
		fmt.Printf("    ⚠️  警告: 白名单为空时系统将无法启动\n")
	} else {
		fmt.Printf("    ℹ️  说明: 白名单为空时系统进入待机状态，不参与爬取\n")
	}
	fmt.Printf("  ✓ 白名单监控: %v\n", cfg.WhitelistMonitoring)
	if cfg.WhitelistMonitoring {
		fmt.Printf("  ✓ 监控间隔: %d 秒\n", cfg.WhitelistMonitoringInterval)
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// TODO: 实际实现时使用真实的whitelist-blacklist-manager库
	manager := &PlaceholderIPStatusManager{
		whitelist:                    make(map[string]bool),
		blacklist:                    make(map[string]string),
		minWhitelistCount:            cfg.MinWhitelistCount,
		allowStartWhenEmpty:          cfg.AllowStartWhenEmpty,
		whitelistMonitoring:          cfg.WhitelistMonitoring,
		whitelistMonitoringInterval:  cfg.WhitelistMonitoringInterval,
	}

	fmt.Printf("  📊 当前状态:\n")
	fmt.Printf("    - 白名单IP数量: %d\n", manager.GetWhitelistCount())
	fmt.Printf("    - 黑名单IP数量: %d\n", len(manager.blacklist))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if logger != nil {
		logger.Info("黑白名单模块初始化完成，min_whitelist_count=%d, allow_start_when_empty=%v",
			cfg.MinWhitelistCount, cfg.AllowStartWhenEmpty)
	}
	return manager, nil
}

func (m *PlaceholderIPStatusManager) AddToWhitelist(ip string) error {
	m.whitelist[ip] = true
	delete(m.blacklist, ip)
	return nil
}

func (m *PlaceholderIPStatusManager) RemoveFromWhitelist(ip string, reason string) error {
	delete(m.whitelist, ip)
	return nil
}

func (m *PlaceholderIPStatusManager) AddToBlacklist(ip string, reason string) error {
	m.blacklist[ip] = reason
	delete(m.whitelist, ip)
	return nil
}

func (m *PlaceholderIPStatusManager) GetStatus(ip string) string {
	if m.whitelist[ip] {
		return "whitelist"
	}
	if _, ok := m.blacklist[ip]; ok {
		return "blacklist"
	}
	return "unknown"
}

func (m *PlaceholderIPStatusManager) GetWhitelistIPs() []string {
	ips := make([]string, 0, len(m.whitelist))
	for ip := range m.whitelist {
		ips = append(ips, ip)
	}
	return ips
}

func (m *PlaceholderIPStatusManager) GetWhitelistCount() int {
	return len(m.whitelist)
}

func (m *PlaceholderIPStatusManager) CheckSystemHealth() error {
	if len(m.whitelist) == 0 && !m.allowStartWhenEmpty {
		return fmt.Errorf("白名单为空且不允许启动")
	}
	return nil
}

func (m *PlaceholderIPStatusManager) SetMinWhitelistCount(count int) {
	m.minWhitelistCount = count
}

func (m *PlaceholderIPStatusManager) SetAllowStartWhenEmpty(allow bool) {
	m.allowStartWhenEmpty = allow
}

func (m *PlaceholderIPStatusManager) SetWhitelistMonitoring(enabled bool) {
	m.whitelistMonitoring = enabled
}

func (m *PlaceholderIPStatusManager) SetWhitelistMonitoringInterval(interval int) {
	m.whitelistMonitoringInterval = interval
}

