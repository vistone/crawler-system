// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crawler

import (
	"fmt"

	"github.com/vistone/certs"
	"github.com/vistone/domaindns"
	"github.com/vistone/fingerprint"
	"github.com/vistone/localippool"
	"github.com/vistone/logs"
	"github.com/vistone/netconnpool"
	"github.com/vistone/quic"

	"github.com/vistone/crawler-system/internal/config"
	"github.com/vistone/crawler-system/internal/moduleinit"
)

// System 系统主结构体，管理所有模块
type System struct {
	Config *SystemConfig

	// 9个核心模块
	Logger             *logs.Logger
	FingerprintManager *moduleinit.FingerprintManager
	DNSMonitor         domaindns.DomainMonitor
	LocalIPPool        localippool.IPPool
	ConnManager        *moduleinit.ConnManager
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

// NewSystem 创建系统实例
func NewSystem(configPath string) (*System, error) {
	// 1. 加载配置
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	system := &System{
		Config: cfg,
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

// 配置转换函数
func convertLogsConfig(cfg *LogsConfig) *config.LogsConfig {
	return &config.LogsConfig{
		Level:       cfg.Level,
		FileEnabled: cfg.FileEnabled,
		FilePath:    cfg.FilePath,
		MaxSize:     cfg.MaxSize,
		MaxBackups:  cfg.MaxBackups,
		Compress:    cfg.Compress,
		Format:      cfg.Format,
		ShowCaller:  cfg.ShowCaller,
	}
}

func convertFingerprintConfig(cfg *FingerprintConfig) *config.FingerprintConfig {
	return &config.FingerprintConfig{
		SelectionStrategy: cfg.SelectionStrategy,
		EnableRotation:    cfg.EnableRotation,
		RotationInterval:  cfg.RotationInterval,
		LibraryPath:       cfg.LibraryPath,
		Browsers:          cfg.Browsers,
		OSRandomization:   cfg.OSRandomization,
		UARandomization:   cfg.UARandomization,
	}
}

func convertDomainDNSConfig(cfg *DomainDNSConfig) *config.DomainDNSConfig {
	return &config.DomainDNSConfig{
		DNSServers:         cfg.DNSServers,
		CacheEnabled:       cfg.CacheEnabled,
		CacheTTL:           cfg.CacheTTL,
		Timeout:            cfg.Timeout,
		MaxRetries:         cfg.MaxRetries,
		RetryInterval:      cfg.RetryInterval,
		PollutionDetection: cfg.PollutionDetection,
		IPv6Enabled:        cfg.IPv6Enabled,
		IPInfoToken:        cfg.IPInfoToken,
	}
}

func convertLocalIPPoolConfig(cfg *LocalIPPoolConfig) *config.LocalIPPoolConfig {
	return &config.LocalIPPoolConfig{
		IPs:                   cfg.IPs,
		SelectionStrategy:     cfg.SelectionStrategy,
		HealthCheckEnabled:    cfg.HealthCheckEnabled,
		HealthCheckInterval:   cfg.HealthCheckInterval,
		HealthCheckTimeout:    cfg.HealthCheckTimeout,
		MaxFailures:           cfg.MaxFailures,
		RecoveryCheckInterval: cfg.RecoveryCheckInterval,
	}
}

func convertCertificateConfig(cfg *CertificateConfig) *config.CertificateConfig {
	return &config.CertificateConfig{
		ServerDomain:           cfg.ServerDomain,
		CertStoragePath:        cfg.CertStoragePath,
		Provider:               cfg.Provider,
		AutoRenewal:            cfg.AutoRenewal,
		RenewalCheckInterval:   cfg.RenewalCheckInterval,
		RenewalBeforeDays:      cfg.RenewalBeforeDays,
		LetsEncryptEmail:       cfg.LetsEncryptEmail,
		LetsEncryptEnvironment: cfg.LetsEncryptEnvironment,
		SelfSignedValidityDays: cfg.SelfSignedValidityDays,
		AutoDetectLocalIP:      cfg.AutoDetectLocalIP,
	}
}

func convertIPStatusConfig(cfg *IPStatusConfig) *config.IPStatusConfig {
	return &config.IPStatusConfig{
		MinWhitelistCount:           cfg.MinWhitelistCount,
		AllowStartWhenEmpty:         cfg.AllowStartWhenEmpty,
		WhitelistMonitoring:         cfg.WhitelistMonitoring,
		WhitelistMonitoringInterval: cfg.WhitelistMonitoringInterval,
	}
}

func convertConnConfig(cfg *ConnConfig) *config.ConnConfig {
	return &config.ConnConfig{
		ConnectTimeout:      cfg.ConnectTimeout,
		ReadTimeout:         cfg.ReadTimeout,
		WriteTimeout:        cfg.WriteTimeout,
		KeepAlive:           cfg.KeepAlive,
		KeepAliveTime:       cfg.KeepAliveTime,
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxConnsPerHost:     cfg.MaxConnsPerHost,
		TLSHandshakeTimeout: cfg.TLSHandshakeTimeout,
		InsecureSkipVerify:  cfg.InsecureSkipVerify,
	}
}

func convertNetConnPoolConfig(cfg *NetConnPoolConfig) *config.NetConnPoolConfig {
	return &config.NetConnPoolConfig{
		MaxConnections:      cfg.MaxConnections,
		InitialConnections:  cfg.InitialConnections,
		AcquireTimeout:      cfg.AcquireTimeout,
		IdleTimeout:         cfg.IdleTimeout,
		MaxLifetime:         cfg.MaxLifetime,
		HealthCheckInterval: cfg.HealthCheckInterval,
		HealthCheckTimeout:  cfg.HealthCheckTimeout,
	}
}

func convertQUICConfig(cfg *QUICConfig) *config.QUICConfig {
	return &config.QUICConfig{
		MaxConnections:      cfg.MaxConnections,
		InitialConnections:  cfg.InitialConnections,
		AcquireTimeout:      cfg.AcquireTimeout,
		IdleTimeout:         cfg.IdleTimeout,
		MaxLifetime:         cfg.MaxLifetime,
		HealthCheckInterval: cfg.HealthCheckInterval,
		HealthCheckTimeout:  cfg.HealthCheckTimeout,
		HandshakeTimeout:    cfg.HandshakeTimeout,
		Enable0RTT:          cfg.Enable0RTT,
	}
}

// initLogs 初始化日志系统（模块1）
func (s *System) initLogs() error {
	cfg := convertLogsConfig(&s.Config.Logs)
	logger, err := moduleinit.InitLogs(cfg)
	if err != nil {
		return err
	}
	s.Logger = logger
	return nil
}

// initFingerprint 初始化指纹模块（模块2）
func (s *System) initFingerprint() error {
	cfg := convertFingerprintConfig(&s.Config.Fingerprint)
	fm, err := moduleinit.InitFingerprint(cfg, s.Logger)
	if err != nil {
		return err
	}
	s.FingerprintManager = fm
	return nil
}

// GetRandomFingerprint 获取随机指纹（包装fingerprint库）
// 注意：这个方法应该在FingerprintManager中实现，而不是在System中
// 这里暂时保留作为占位符
func GetRandomFingerprint(fm *moduleinit.FingerprintManager) (*fingerprint.FingerprintResult, error) {
	return fingerprint.GetRandomFingerprint()
}

// initDomainDNS 初始化DNS解析模块（模块3）
func (s *System) initDomainDNS() error {
	cfg := convertDomainDNSConfig(&s.Config.DomainDNS)
	targetDomains := s.Config.IPPoolTest.TargetDomains
	monitor, err := moduleinit.InitDomainDNS(cfg, targetDomains, s.Logger)
	if err != nil {
		return err
	}
	s.DNSMonitor = monitor
	return nil
}

// initLocalIPPool 初始化本地IP池模块（模块4）
func (s *System) initLocalIPPool() error {
	cfg := convertLocalIPPoolConfig(&s.Config.LocalIPPool)
	pool, err := moduleinit.InitLocalIPPool(cfg, s.Logger)
	if err != nil {
		return err
	}
	s.LocalIPPool = pool
	return nil
}

// initCerts 初始化证书模块（模块5）
func (s *System) initCerts() error {
	cfg := convertCertificateConfig(&s.Config.Certificate)
	manager, err := moduleinit.InitCerts(cfg, s.Logger)
	if err != nil {
		return err
	}
	s.CertManager = manager
	return nil
}

// initIPStatusManager 初始化黑白名单模块（模块6）
func (s *System) initIPStatusManager() error {
	cfg := convertIPStatusConfig(&s.Config.IPStatus)
	manager, err := moduleinit.InitIPStatusManager(cfg, s.Logger)
	if err != nil {
		return err
	}
	s.IPStatusManager = manager
	return nil
}

// initConn 初始化连接模块（模块7）
func (s *System) initConn() error {
	cfg := convertConnConfig(&s.Config.Conn)
	cm, err := moduleinit.InitConn(cfg, s.Logger)
	if err != nil {
		return err
	}
	s.ConnManager = cm
	return nil
}

// initNetConnPool 初始化TCP连接池模块（模块8）
func (s *System) initNetConnPool() error {
	cfg := convertNetConnPoolConfig(&s.Config.NetConnPool)
	pool, err := moduleinit.InitNetConnPool(cfg, s.Logger)
	if err != nil {
		return err
	}
	s.NetConnPool = pool
	return nil
}

// initQUICPool 初始化QUIC连接池模块（模块9）
func (s *System) initQUICPool() error {
	cfg := convertQUICConfig(&s.Config.QUIC)
	pool, err := moduleinit.InitQUICPool(cfg, s.Logger)
	if err != nil {
		return err
	}
	s.QUICPool = pool
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
