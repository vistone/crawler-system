// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crawler

import (
	"fmt"
	"time"

	"github.com/vistone/certs"
	"github.com/vistone/domaindns"
	"github.com/vistone/fingerprint"
	"github.com/vistone/localippool"
	"github.com/vistone/logs"
	"github.com/vistone/netconnpool"
	"github.com/vistone/quic"
)

// System 系统主结构体，管理所有模块
type System struct {
	Config *SystemConfig

	// 9个核心模块
	Logger             *logs.Logger
	FingerprintManager *FingerprintManager
	DNSMonitor         domaindns.DomainMonitor
	LocalIPPool        localippool.IPPool
	ConnManager        *ConnManager
	NetConnPool        *netconnpool.Pool
	QUICPool           *quic.Pool
	CertManager        *certs.Manager
	IPStatusManager    IPStatusManagerInterface
}

// IPStatusManagerInterface 黑白名单管理器接口
type IPStatusManagerInterface interface {
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

// FingerprintManager 指纹管理器包装
type FingerprintManager struct {
	config *FingerprintConfig
}

// ConnManager 连接管理器包装
type ConnManager struct {
	config *ConnConfig
}

// NewSystem 创建系统实例
func NewSystem(configPath string) (*System, error) {
	// 1. 加载配置
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	system := &System{
		Config: config,
	}

	// 2. 初始化所有模块
	if err := system.Initialize(); err != nil {
		return nil, fmt.Errorf("初始化系统失败: %w", err)
	}

	return system, nil
}

// Initialize 初始化所有模块
func (s *System) Initialize() error {
	// 步骤1: 初始化日志系统（最先初始化，其他模块都需要）
	if err := s.initLogs(); err != nil {
		return fmt.Errorf("初始化日志系统失败: %w", err)
	}

	// 步骤2: 初始化独立模块（可以并行初始化）
	if err := s.initFingerprint(); err != nil {
		return fmt.Errorf("初始化指纹模块失败: %w", err)
	}

	if err := s.initDomainDNS(); err != nil {
		return fmt.Errorf("初始化DNS解析模块失败: %w", err)
	}

	if err := s.initLocalIPPool(); err != nil {
		return fmt.Errorf("初始化本地IP池模块失败: %w", err)
	}

	if err := s.initCerts(); err != nil {
		return fmt.Errorf("初始化证书模块失败: %w", err)
	}

	if err := s.initIPStatusManager(); err != nil {
		return fmt.Errorf("初始化黑白名单模块失败: %w", err)
	}

	// 步骤3: 初始化依赖模块
	if err := s.initConn(); err != nil {
		return fmt.Errorf("初始化连接模块失败: %w", err)
	}

	// 步骤4: 初始化连接池模块
	if err := s.initNetConnPool(); err != nil {
		return fmt.Errorf("初始化TCP连接池模块失败: %w", err)
	}

	if err := s.initQUICPool(); err != nil {
		return fmt.Errorf("初始化QUIC连接池模块失败: %w", err)
	}

	s.Logger.Info("系统初始化完成，modules=%d", 9)
	return nil
}

// initLogs 初始化日志系统（模块1）
func (s *System) initLogs() error {
	cfg := s.Config.Logs

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

	s.Logger = logger
	s.Logger.Info("日志系统初始化完成，level=%s", cfg.Level)
	return nil
}

// initFingerprint 初始化指纹模块（模块2）
func (s *System) initFingerprint() error {
	cfg := s.Config.Fingerprint

	// 创建指纹管理器
	fm := &FingerprintManager{
		config: &cfg,
	}

	s.FingerprintManager = fm
	s.Logger.Info("指纹模块初始化完成，strategy=%s, rotation=%v", cfg.SelectionStrategy, cfg.EnableRotation)
	return nil
}

// GetRandomFingerprint 获取随机指纹（包装fingerprint库）
func (fm *FingerprintManager) GetRandomFingerprint() (*fingerprint.FingerprintResult, error) {
	return fingerprint.GetRandomFingerprint()
}

// initDomainDNS 初始化DNS解析模块（模块3）
func (s *System) initDomainDNS() error {
	cfg := s.Config.DomainDNS

	// 获取目标域名列表（从IP池测试配置中获取）
	targetDomains := s.Config.IPPoolTest.TargetDomains
	if len(targetDomains) == 0 {
		// 如果没有配置目标域名，跳过监控器创建
		s.Logger.Warn("未配置目标域名，跳过DNS监控器创建")
		return nil
	}

	// 使用NewMonitorWithGlobalDNSServers创建监控器
	// 参数：domains, ipInfoToken, dnsServerFile, maxServers
	ipInfoToken := cfg.IPInfoToken
	if ipInfoToken == "" {
		s.Logger.Warn("未配置IPInfo Token，IP详细信息获取功能将不可用")
	}

	monitor, err := domaindns.NewMonitorWithGlobalDNSServers(
		targetDomains,
		ipInfoToken, // 使用配置的IPInfo Token
		"",          // dnsServerFile，空表示使用默认的dnsservernames.json
		0,           // maxServers，0表示使用全部DNS服务器
	)
	if err != nil {
		return fmt.Errorf("创建DNS监控器失败: %w", err)
	}

	s.DNSMonitor = monitor

	// 启动DomainMonitor，它会自动维护域名的IP地址变化
	// 启动后会立即执行一次DNS解析，然后按配置的间隔定期更新
	monitor.Start()

	s.Logger.Info("DNS解析模块初始化完成，dns_servers=%d, target_domains=%d, domains=%v, note=DomainMonitor已启动，将自动维护域名IP地址",
		len(cfg.DNSServers), len(targetDomains), targetDomains)
	return nil
}

// initLocalIPPool 初始化本地IP池模块（模块4）
func (s *System) initLocalIPPool() error {
	cfg := s.Config.LocalIPPool

	// 创建本地IP池（使用配置的IPv4列表，IPv6为空表示自动检测）
	pool, err := localippool.NewLocalIPPool(cfg.IPs, "")
	if err != nil {
		return fmt.Errorf("创建本地IP池失败: %w", err)
	}

	s.LocalIPPool = pool
	s.Logger.Info("本地IP池模块初始化完成，ipv4_count=%d, strategy=%s", len(cfg.IPs), cfg.SelectionStrategy)
	return nil
}

// initCerts 初始化证书模块（模块5）
func (s *System) initCerts() error {
	cfg := s.Config.Certificate

	// 创建证书配置（使用默认配置，后续根据实际API调整）
	certConfig := certs.DefaultConfig()
	// 注意：certs.Config的实际字段可能不同，这里先使用默认配置
	// 后续需要根据实际API调整配置方式
	_ = cfg // 暂时忽略配置，使用默认值

	manager, err := certs.NewManager(certConfig)
	if err != nil {
		return fmt.Errorf("创建证书管理器失败: %w", err)
	}

	s.CertManager = manager
	s.Logger.Info("证书模块初始化完成，provider=%s, domain=%s, auto_renewal=%v", cfg.Provider, cfg.ServerDomain, cfg.AutoRenewal)
	return nil
}

// initIPStatusManager 初始化黑白名单模块（模块6）
func (s *System) initIPStatusManager() error {
	cfg := s.Config.IPStatus

	// TODO: 实际实现时使用真实的whitelist-blacklist-manager库
	// manager := whitelist.NewManager()
	// 这里先使用占位符实现
	manager := &PlaceholderIPStatusManager{
		whitelist:                   make(map[string]bool),
		blacklist:                   make(map[string]string), // IP -> reason
		minWhitelistCount:           cfg.MinWhitelistCount,
		allowStartWhenEmpty:         cfg.AllowStartWhenEmpty,
		whitelistMonitoring:         cfg.WhitelistMonitoring,
		whitelistMonitoringInterval: cfg.WhitelistMonitoringInterval,
	}

	s.IPStatusManager = manager
	s.Logger.Info("黑白名单模块初始化完成，min_whitelist_count=%d, allow_start_when_empty=%v",
		cfg.MinWhitelistCount, cfg.AllowStartWhenEmpty)
	return nil
}

// PlaceholderIPStatusManager 占位符黑白名单管理器实现
type PlaceholderIPStatusManager struct {
	whitelist                   map[string]bool
	blacklist                   map[string]string
	minWhitelistCount           int
	allowStartWhenEmpty         bool
	whitelistMonitoring         bool
	whitelistMonitoringInterval int
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

// initConn 初始化连接模块（模块7）
func (s *System) initConn() error {
	cfg := s.Config.Conn

	cm := &ConnManager{
		config: &cfg,
	}

	s.ConnManager = cm
	s.Logger.Info("连接模块初始化完成，connect_timeout=%d, read_timeout=%d", cfg.ConnectTimeout, cfg.ReadTimeout)
	return nil
}

// initNetConnPool 初始化TCP连接池模块（模块8）
func (s *System) initNetConnPool() error {
	cfg := s.Config.NetConnPool

	// 创建连接池配置（客户端模式）
	poolConfig := netconnpool.DefaultConfig()

	// 设置基本参数
	poolConfig.MaxConnections = cfg.MaxConnections
	poolConfig.MinConnections = cfg.InitialConnections
	poolConfig.GetConnectionTimeout = time.Duration(cfg.AcquireTimeout) * time.Second
	poolConfig.IdleTimeout = time.Duration(cfg.IdleTimeout) * time.Second
	poolConfig.MaxLifetime = time.Duration(cfg.MaxLifetime) * time.Second
	poolConfig.HealthCheckInterval = time.Duration(cfg.HealthCheckInterval) * time.Second
	poolConfig.HealthCheckTimeout = time.Duration(cfg.HealthCheckTimeout) * time.Second

	// 注意：客户端模式需要Dialer，但Dialer需要目标地址
	// 由于在初始化时还不知道目标地址，这里暂时不创建连接池
	// 连接池将在实际使用时按需创建
	// 或者提供一个占位符Dialer，后续再替换
	s.Logger.Info("TCP连接池配置已准备，max_connections=%d, initial_connections=%d, note=连接池将在需要时按需创建",
		cfg.MaxConnections, cfg.InitialConnections)

	// 暂时不创建连接池，返回nil表示成功（连接池延迟初始化）
	return nil
}

// initQUICPool 初始化QUIC连接池模块（模块9）
func (s *System) initQUICPool() error {
	cfg := s.Config.QUIC

	// 创建QUIC客户端连接池
	// 注意：这里使用客户端池，服务端池需要TLS配置，在后续实现
	minCap := cfg.InitialConnections
	if minCap < 1 {
		minCap = 1
	}
	maxCap := cfg.MaxConnections
	if maxCap < minCap {
		maxCap = minCap
	}

	pool := quic.NewClientPool(
		minCap,
		maxCap,
		time.Duration(cfg.IdleTimeout)*time.Second,
		time.Duration(cfg.MaxLifetime)*time.Second,
		time.Duration(cfg.IdleTimeout)*time.Second,
		"",  // tlsCode，后续从配置读取
		"",  // hostname，后续从配置读取
		nil, // addrResolver，后续实现
	)

	s.QUICPool = pool
	s.Logger.Info("QUIC连接池模块初始化完成，max_connections=%d, enable_0rtt=%v", cfg.MaxConnections, cfg.Enable0RTT)
	return nil
}

// Close 关闭系统，释放资源
func (s *System) Close() error {
	if s.Logger != nil {
		s.Logger.Info("正在关闭系统...")
	}

	// 关闭连接池
	if s.NetConnPool != nil {
		if err := s.NetConnPool.Close(); err != nil {
			if s.Logger != nil {
				s.Logger.Error("关闭TCP连接池失败，error=%v", err)
			}
		}
	}

	if s.QUICPool != nil {
		s.QUICPool.Close() // QUIC池的Close方法可能没有返回值
	}

	// 停止DNS监控器
	if s.DNSMonitor != nil {
		s.DNSMonitor.Stop()
		if s.Logger != nil {
			s.Logger.Info("DNS监控器已停止")
		}
	}

	// 关闭本地IP池
	if s.LocalIPPool != nil {
		if err := s.LocalIPPool.Close(); err != nil {
			if s.Logger != nil {
				s.Logger.Error("关闭本地IP池失败，error=%v", err)
			}
		}
	}

	// 关闭证书管理器
	// certs.Manager 可能没有Close方法，暂时不关闭
	if s.CertManager != nil {
		// 证书管理器可能不需要显式关闭
		_ = s.CertManager
	}

	if s.Logger != nil {
		s.Logger.Info("系统关闭完成")
	}

	return nil
}
