// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crawler

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// SystemConfig 系统配置
type SystemConfig struct {
	Logs              LogsConfig              `toml:"logs"`
	Fingerprint       FingerprintConfig       `toml:"fingerprint"`
	DomainDNS         DomainDNSConfig         `toml:"domaindns"`
	LocalIPPool       LocalIPPoolConfig       `toml:"local_ip_pool"`
	Conn              ConnConfig              `toml:"conn"`
	NetConnPool       NetConnPoolConfig       `toml:"netconnpool"`
	QUIC              QUICConfig              `toml:"quic"`
	Certificate       CertificateConfig       `toml:"certificate"`
	IPStatus          IPStatusConfig          `toml:"ip_status"`
	IPPoolTest        IPPoolTestConfig        `toml:"ip_pool_test"`
	BlacklistRecovery BlacklistRecoveryConfig `toml:"blacklist_recovery"`
	StatusReport      StatusReportConfig      `toml:"status_report"`
	Server            ServerConfig            `toml:"server"`
	Crawler           CrawlerConfig           `toml:"crawler"`
	System            SystemInfoConfig        `toml:"system"`
}

// LogsConfig 日志配置
type LogsConfig struct {
	Level       string `toml:"level"`       // debug, info, warn, error
	OutputPath  string `toml:"output_path"` // 空表示输出到标准输出
	FileEnabled bool   `toml:"file_enabled"`
	FilePath    string `toml:"file_path"`
	MaxSize     int    `toml:"max_size"` // MB
	MaxBackups  int    `toml:"max_backups"`
	Compress    bool   `toml:"compress"`
	Format      string `toml:"format"` // json, text
	ShowCaller  bool   `toml:"show_caller"`
}

// FingerprintConfig 指纹配置
type FingerprintConfig struct {
	SelectionStrategy string   `toml:"selection_strategy"` // random, round_robin, least_used
	EnableRotation    bool     `toml:"enable_rotation"`
	RotationInterval  int      `toml:"rotation_interval"` // 秒
	LibraryPath       string   `toml:"library_path"`
	Browsers          []string `toml:"browsers"`
	OSRandomization   bool     `toml:"os_randomization"`
	UARandomization   bool     `toml:"ua_randomization"`
}

// DomainDNSConfig DNS解析配置
type DomainDNSConfig struct {
	DNSServers         []string `toml:"dns_servers"`
	CacheEnabled       bool     `toml:"cache_enabled"`
	CacheTTL           int      `toml:"cache_ttl"` // 秒
	Timeout            int      `toml:"timeout"`   // 秒
	MaxRetries         int      `toml:"max_retries"`
	RetryInterval      int      `toml:"retry_interval"` // 秒
	PollutionDetection bool     `toml:"pollution_detection"`
	IPv6Enabled        bool     `toml:"ipv6_enabled"`
	IPInfoToken        string   `toml:"ipinfo_token"` // IPInfo.io API Token
}

// LocalIPPoolConfig 本地IP池配置
type LocalIPPoolConfig struct {
	IPs                   []string `toml:"ips"`
	SelectionStrategy     string   `toml:"selection_strategy"` // random, round_robin, least_used
	HealthCheckEnabled    bool     `toml:"health_check_enabled"`
	HealthCheckInterval   int      `toml:"health_check_interval"` // 秒
	HealthCheckTimeout    int      `toml:"health_check_timeout"`  // 秒
	MaxFailures           int      `toml:"max_failures"`
	RecoveryCheckInterval int      `toml:"recovery_check_interval"` // 秒
}

// ConnConfig 连接配置
type ConnConfig struct {
	ConnectTimeout      int  `toml:"connect_timeout"` // 秒
	ReadTimeout         int  `toml:"read_timeout"`    // 秒
	WriteTimeout        int  `toml:"write_timeout"`   // 秒
	KeepAlive           bool `toml:"keep_alive"`
	KeepAliveTime       int  `toml:"keep_alive_time"` // 秒
	MaxIdleConns        int  `toml:"max_idle_conns"`
	MaxConnsPerHost     int  `toml:"max_conns_per_host"`
	TLSHandshakeTimeout int  `toml:"tls_handshake_timeout"` // 秒
	InsecureSkipVerify  bool `toml:"insecure_skip_verify"`
}

// NetConnPoolConfig TCP连接池配置
type NetConnPoolConfig struct {
	MaxConnections      int `toml:"max_connections"`
	InitialConnections  int `toml:"initial_connections"`
	AcquireTimeout      int `toml:"acquire_timeout"`       // 秒
	IdleTimeout         int `toml:"idle_timeout"`          // 秒
	MaxLifetime         int `toml:"max_lifetime"`          // 秒
	HealthCheckInterval int `toml:"health_check_interval"` // 秒
	HealthCheckTimeout  int `toml:"health_check_timeout"`  // 秒
}

// QUICConfig QUIC连接池配置
type QUICConfig struct {
	MaxConnections      int  `toml:"max_connections"`
	InitialConnections  int  `toml:"initial_connections"`
	AcquireTimeout      int  `toml:"acquire_timeout"`       // 秒
	IdleTimeout         int  `toml:"idle_timeout"`          // 秒
	MaxLifetime         int  `toml:"max_lifetime"`          // 秒
	HealthCheckInterval int  `toml:"health_check_interval"` // 秒
	HealthCheckTimeout  int  `toml:"health_check_timeout"`  // 秒
	HandshakeTimeout    int  `toml:"handshake_timeout"`     // 秒
	Enable0RTT          bool `toml:"enable_0rtt"`
}

// CertificateConfig 证书配置
type CertificateConfig struct {
	ServerDomain           string `toml:"server_domain"`
	CertStoragePath        string `toml:"cert_storage_path"`
	Provider               string `toml:"provider"` // letsencrypt, self-signed
	AutoRenewal            bool   `toml:"auto_renewal"`
	RenewalCheckInterval   int    `toml:"renewal_check_interval"` // 小时
	RenewalBeforeDays      int    `toml:"renewal_before_days"`
	LetsEncryptEmail       string `toml:"letsencrypt_email"`
	LetsEncryptEnvironment string `toml:"letsencrypt_environment"` // production, staging
	AutoDetectLocalIP      bool   `toml:"auto_detect_local_ip"`
	SelfSignedValidityDays int    `toml:"self_signed_validity_days"`
}

// IPStatusConfig 黑白名单配置
type IPStatusConfig struct {
	MinWhitelistCount           int  `toml:"min_whitelist_count"`
	AllowStartWhenEmpty         bool `toml:"allow_start_when_empty"`
	WhitelistMonitoring         bool `toml:"whitelist_monitoring"`
	WhitelistMonitoringInterval int  `toml:"whitelist_monitoring_interval"` // 秒
}

// IPPoolTestConfig IP池测试配置
type IPPoolTestConfig struct {
	TargetDomains        []string `toml:"target_domains"`
	TestURL              string   `toml:"test_url"`
	TestMethod           string   `toml:"test_method"` // GET, HEAD
	MaxConcurrent        int      `toml:"max_concurrent"`
	TestTimeout          int      `toml:"test_timeout"` // 秒
	RetryCount           int      `toml:"retry_count"`
	RetryInterval        int      `toml:"retry_interval"` // 秒
	TestInterval         int      `toml:"test_interval"`  // 秒
	UseFingerprint       bool     `toml:"use_fingerprint"`
	SuccessStatusCodes   []int    `toml:"success_status_codes"`
	ForbiddenStatusCodes []int    `toml:"forbidden_status_codes"`
}

// BlacklistRecoveryConfig 黑名单恢复配置
type BlacklistRecoveryConfig struct {
	Enabled        bool   `toml:"enabled"`
	CheckInterval  int    `toml:"check_interval"`   // 秒
	IPTestInterval int    `toml:"ip_test_interval"` // 秒
	MaxConcurrent  int    `toml:"max_concurrent"`
	TestTimeout    int    `toml:"test_timeout"` // 秒
	TestURL        string `toml:"test_url"`
	TestMethod     string `toml:"test_method"` // GET, HEAD
	UseFingerprint bool   `toml:"use_fingerprint"`
}

// StatusReportConfig 状态报告配置
type StatusReportConfig struct {
	ReportInterval      int  `toml:"report_interval"` // 秒
	ReportOnChange      bool `toml:"report_on_change"`
	ReportClientDetails bool `toml:"report_client_details"`
	ReportIPList        bool `toml:"report_ip_list"`
	MaxReportIPs        int  `toml:"max_report_ips"`
	CompressData        bool `toml:"compress_data"`
}

// ServerConfig 服务端配置
type ServerConfig struct {
	ListenAddress     string `toml:"listen_address"`
	QUICEnabled       bool   `toml:"quic_enabled"`
	MaxClients        int    `toml:"max_clients"`
	ClientTimeout     int    `toml:"client_timeout"` // 秒
	ClientAuthEnabled bool   `toml:"client_auth_enabled"`
	ClientCertPath    string `toml:"client_cert_path"`
	AccessLogEnabled  bool   `toml:"access_log_enabled"`
	AccessLogPath     string `toml:"access_log_path"`
}

// CrawlerConfig 爬虫配置
type CrawlerConfig struct {
	DefaultTimeout       int      `toml:"default_timeout"` // 秒
	MaxRetries           int      `toml:"max_retries"`
	RetryInterval        int      `toml:"retry_interval"` // 秒
	ProtocolPriority     []string `toml:"protocol_priority"`
	ProtocolFallback     bool     `toml:"protocol_fallback"`
	Concurrency          int      `toml:"concurrency"`
	RateLimit            int      `toml:"rate_limit"` // 每秒请求数，0表示不限制
	QueueEnabled         bool     `toml:"queue_enabled"`
	QueueSize            int      `toml:"queue_size"`
	DeduplicationEnabled bool     `toml:"deduplication_enabled"`
	DeduplicationTTL     int      `toml:"deduplication_ttl"` // 秒
}

// SystemInfoConfig 系统信息配置
type SystemInfoConfig struct {
	Name                  string `toml:"name"`
	Version               string `toml:"version"`
	WorkDir               string `toml:"work_dir"`
	DataDir               string `toml:"data_dir"`
	PerformanceMonitoring bool   `toml:"performance_monitoring"`
	PerformanceInterval   int    `toml:"performance_interval"` // 秒
	HealthCheckEnabled    bool   `toml:"health_check_enabled"`
	HealthCheckPort       int    `toml:"health_check_port"`
	MetricsEnabled        bool   `toml:"metrics_enabled"`
	MetricsPort           int    `toml:"metrics_port"`
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*SystemConfig, error) {
	// 如果文件不存在，返回默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析TOML配置
	config := &SystemConfig{}
	if err := toml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *SystemConfig {
	return &SystemConfig{
		Logs: LogsConfig{
			Level:       "info",
			OutputPath:  "",
			FileEnabled: true,
			FilePath:    "./logs/crawler.log",
			MaxSize:     100,
			MaxBackups:  10,
			Compress:    true,
			Format:      "text",
			ShowCaller:  true,
		},
		Fingerprint: FingerprintConfig{
			SelectionStrategy: "random",
			EnableRotation:    true,
			RotationInterval:  300,
			LibraryPath:       "",
			Browsers:          []string{},
			OSRandomization:   true,
			UARandomization:   true,
		},
		DomainDNS: DomainDNSConfig{
			DNSServers:         []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "1.0.0.1"},
			CacheEnabled:       true,
			CacheTTL:           3600,
			Timeout:            5,
			MaxRetries:         3,
			RetryInterval:      1,
			PollutionDetection: true,
			IPv6Enabled:        true,
		},
		LocalIPPool: LocalIPPoolConfig{
			IPs:                   []string{},
			SelectionStrategy:     "round_robin",
			HealthCheckEnabled:    true,
			HealthCheckInterval:   60,
			HealthCheckTimeout:    5,
			MaxFailures:           3,
			RecoveryCheckInterval: 300,
		},
		Conn: ConnConfig{
			ConnectTimeout:      10,
			ReadTimeout:         30,
			WriteTimeout:        30,
			KeepAlive:           true,
			KeepAliveTime:       60,
			MaxIdleConns:        100,
			MaxConnsPerHost:     10,
			TLSHandshakeTimeout: 10,
			InsecureSkipVerify:  false,
		},
		NetConnPool: NetConnPoolConfig{
			MaxConnections:      100,
			InitialConnections:  10,
			AcquireTimeout:      5,
			IdleTimeout:         300,
			MaxLifetime:         3600,
			HealthCheckInterval: 60,
			HealthCheckTimeout:  5,
		},
		QUIC: QUICConfig{
			MaxConnections:      50,
			InitialConnections:  5,
			AcquireTimeout:      5,
			IdleTimeout:         300,
			MaxLifetime:         3600,
			HealthCheckInterval: 60,
			HealthCheckTimeout:  5,
			HandshakeTimeout:    10,
			Enable0RTT:          true,
		},
		Certificate: CertificateConfig{
			ServerDomain:           "crawler.example.com",
			CertStoragePath:        "./certs",
			Provider:               "letsencrypt",
			AutoRenewal:            true,
			RenewalCheckInterval:   24,
			RenewalBeforeDays:      30,
			LetsEncryptEmail:       "admin@example.com",
			LetsEncryptEnvironment: "production",
			AutoDetectLocalIP:      true,
			SelfSignedValidityDays: 365,
		},
		IPStatus: IPStatusConfig{
			MinWhitelistCount:           1,
			AllowStartWhenEmpty:         true,
			WhitelistMonitoring:         true,
			WhitelistMonitoringInterval: 60,
		},
		IPPoolTest: IPPoolTestConfig{
			TargetDomains:        []string{},
			TestURL:              "https://{domain}/",
			TestMethod:           "HEAD",
			MaxConcurrent:        10,
			TestTimeout:          10,
			RetryCount:           2,
			RetryInterval:        5,
			TestInterval:         300,
			UseFingerprint:       true,
			SuccessStatusCodes:   []int{200, 201, 202, 204},
			ForbiddenStatusCodes: []int{403},
		},
		BlacklistRecovery: BlacklistRecoveryConfig{
			Enabled:        true,
			CheckInterval:  1800, // 30分钟
			IPTestInterval: 3600, // 1小时
			MaxConcurrent:  5,
			TestTimeout:    10,
			TestURL:        "https://{domain}/",
			TestMethod:     "HEAD",
			UseFingerprint: true,
		},
		StatusReport: StatusReportConfig{
			ReportInterval:      5,
			ReportOnChange:      true,
			ReportClientDetails: true,
			ReportIPList:        true,
			MaxReportIPs:        1000,
			CompressData:        false,
		},
		Server: ServerConfig{
			ListenAddress:     "0.0.0.0:8443",
			QUICEnabled:       true,
			MaxClients:        1000,
			ClientTimeout:     300,
			ClientAuthEnabled: false,
			ClientCertPath:    "",
			AccessLogEnabled:  true,
			AccessLogPath:     "./logs/access.log",
		},
		Crawler: CrawlerConfig{
			DefaultTimeout:       30,
			MaxRetries:           3,
			RetryInterval:        2,
			ProtocolPriority:     []string{"http3", "http2", "http1.1"},
			ProtocolFallback:     true,
			Concurrency:          10,
			RateLimit:            0,
			QueueEnabled:         true,
			QueueSize:            1000,
			DeduplicationEnabled: false,
			DeduplicationTTL:     3600,
		},
		System: SystemInfoConfig{
			Name:                  "crawler-system",
			Version:               "1.0.0",
			WorkDir:               "./",
			DataDir:               "./data",
			PerformanceMonitoring: true,
			PerformanceInterval:   60,
			HealthCheckEnabled:    true,
			HealthCheckPort:       8080,
			MetricsEnabled:        true,
			MetricsPort:           9090,
		},
	}
}

// GetDuration 辅助函数：将秒转换为time.Duration
func (c *ConnConfig) GetConnectTimeout() time.Duration {
	return time.Duration(c.ConnectTimeout) * time.Second
}

func (c *ConnConfig) GetReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeout) * time.Second
}

func (c *ConnConfig) GetWriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeout) * time.Second
}

func (c *ConnConfig) GetKeepAliveTime() time.Duration {
	return time.Duration(c.KeepAliveTime) * time.Second
}

func (c *ConnConfig) GetTLSHandshakeTimeout() time.Duration {
	return time.Duration(c.TLSHandshakeTimeout) * time.Second
}
