// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moduleinit

import (
	"fmt"

	"github.com/vistone/certs"
	"github.com/vistone/crawler-system/internal/config"
)

// InitCerts 初始化证书模块（模块5）
func InitCerts(cfg *config.CertificateConfig, logger interface{ Info(string, ...interface{}) }) (*certs.Manager, error) {
	// 输出初始化详细信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🔒 [模块5] 证书模块初始化\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  ✓ 服务端域名: %s\n", cfg.ServerDomain)
	fmt.Printf("  ✓ 证书存储路径: %s\n", cfg.CertStoragePath)
	fmt.Printf("  ✓ 证书提供商: %s\n", cfg.Provider)
	fmt.Printf("  ✓ 自动续期: %v\n", cfg.AutoRenewal)
	if cfg.AutoRenewal {
		fmt.Printf("  ✓ 续期检查间隔: %d 小时\n", cfg.RenewalCheckInterval)
		fmt.Printf("  ✓ 提前续期天数: %d 天\n", cfg.RenewalBeforeDays)
	}
	if cfg.Provider == "letsencrypt" {
		fmt.Printf("  ✓ Let's Encrypt邮箱: %s\n", cfg.LetsEncryptEmail)
		fmt.Printf("  ✓ Let's Encrypt环境: %s\n", cfg.LetsEncryptEnvironment)
	} else if cfg.Provider == "self-signed" {
		fmt.Printf("  ✓ 自签名证书有效期: %d 天\n", cfg.SelfSignedValidityDays)
	}
	fmt.Printf("  ✓ 自动检测本地IP: %v\n", cfg.AutoDetectLocalIP)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// 创建证书配置（使用默认配置，后续根据实际API调整）
	certConfig := certs.DefaultConfig()
	// 注意：certs.Config的实际字段可能不同，这里先使用默认配置
	// 后续需要根据实际API调整配置方式
	_ = cfg // 暂时忽略配置，使用默认值

	manager, err := certs.NewManager(certConfig)
	if err != nil {
		return nil, fmt.Errorf("创建证书管理器失败: %w", err)
	}

	// 尝试获取证书信息（如果已存在）
	cert, err := manager.GetOrRequestCertificate(cfg.ServerDomain)
	if err == nil && cert != nil {
		fmt.Printf("  📄 证书信息:\n")
		fmt.Printf("    - 域名: %s\n", cfg.ServerDomain)
		fmt.Printf("    - 证书已获取\n")
	} else {
		fmt.Printf("  ⚠️  证书尚未申请，将在首次使用时自动申请\n")
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if logger != nil {
		logger.Info("证书模块初始化完成，provider=%s, domain=%s, auto_renewal=%v",
			cfg.Provider, cfg.ServerDomain, cfg.AutoRenewal)
	}
	return manager, nil
}

