// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

// 配置类型定义，避免循环导入

// LogsConfig 日志配置
type LogsConfig struct {
	Level       string
	FileEnabled bool
	FilePath    string
	MaxSize     int
	MaxBackups  int
	Compress    bool
	Format      string
	ShowCaller  bool
}

// FingerprintConfig 指纹配置
type FingerprintConfig struct {
	SelectionStrategy string
	EnableRotation    bool
	RotationInterval  int
	LibraryPath       string
	Browsers          []string
	OSRandomization   bool
	UARandomization   bool
}

// DomainDNSConfig DNS配置
type DomainDNSConfig struct {
	DNSServers         []string
	CacheEnabled       bool
	CacheTTL           int
	Timeout            int
	MaxRetries         int
	RetryInterval      int
	PollutionDetection bool
	IPv6Enabled        bool
	IPInfoToken        string
}

// LocalIPPoolConfig 本地IP池配置
type LocalIPPoolConfig struct {
	IPs                   []string
	SelectionStrategy     string
	HealthCheckEnabled    bool
	HealthCheckInterval   int
	HealthCheckTimeout    int
	MaxFailures           int
	RecoveryCheckInterval int
}

// CertificateConfig 证书配置
type CertificateConfig struct {
	ServerDomain           string
	CertStoragePath        string
	Provider               string
	AutoRenewal            bool
	RenewalCheckInterval   int
	RenewalBeforeDays      int
	LetsEncryptEmail       string
	LetsEncryptEnvironment string
	SelfSignedValidityDays int
	AutoDetectLocalIP      bool
}

// IPStatusConfig 黑白名单配置
type IPStatusConfig struct {
	MinWhitelistCount           int
	AllowStartWhenEmpty         bool
	WhitelistMonitoring         bool
	WhitelistMonitoringInterval int
}

// ConnConfig 连接配置
type ConnConfig struct {
	ConnectTimeout      int
	ReadTimeout         int
	WriteTimeout        int
	KeepAlive           bool
	KeepAliveTime       int
	MaxIdleConns        int
	MaxConnsPerHost     int
	TLSHandshakeTimeout int
	InsecureSkipVerify  bool
}

// NetConnPoolConfig TCP连接池配置
type NetConnPoolConfig struct {
	MaxConnections      int
	InitialConnections  int
	AcquireTimeout      int
	IdleTimeout         int
	MaxLifetime         int
	HealthCheckInterval int
	HealthCheckTimeout  int
}

// QUICConfig QUIC配置
type QUICConfig struct {
	MaxConnections      int
	InitialConnections  int
	AcquireTimeout      int
	IdleTimeout         int
	MaxLifetime         int
	HealthCheckInterval int
	HealthCheckTimeout  int
	HandshakeTimeout    int
	Enable0RTT          bool
}
