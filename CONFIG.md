# 配置文件说明

## 概述

系统使用 `config.toml` 作为主配置文件，包含9个核心模块的所有配置项。

## 配置文件结构

### 1. 日志配置 (logs)

控制系统的日志输出行为。

**关键配置项**：
- `level`: 日志级别（debug, info, warn, error）
- `file_path`: 日志文件路径
- `max_size`: 日志文件最大大小（MB）
- `format`: 日志格式（json, text）

### 2. 指纹配置 (fingerprint)

控制TLS指纹模拟行为。

**关键配置项**：
- `selection_strategy`: 指纹选择策略（random, round_robin, least_used）
- `enable_rotation`: 是否启用指纹轮换
- `rotation_interval`: 指纹轮换间隔（秒）
- `browsers`: 支持的浏览器列表（空表示使用所有）

### 3. DNS解析配置 (domaindns)

控制DNS解析行为。

**关键配置项**：
- `dns_servers`: DNS服务器列表
- `cache_enabled`: 是否启用DNS缓存
- `cache_ttl`: DNS缓存TTL（秒）
- `timeout`: DNS查询超时（秒）
- `pollution_detection`: 是否启用DNS污染检测

### 4. 本地IP池配置 (local_ip_pool)

控制本地出口IP池管理。

**关键配置项**：
- `ips`: 本地出口IP列表（用于绑定本地连接，与黑白名单无关）
- `selection_strategy`: IP选择策略（random, round_robin, least_used）
- `health_check_enabled`: 是否启用IP健康检查
- `health_check_interval`: IP健康检查间隔（秒）

**重要说明**：
- 本地IP池用于绑定本地出口IP，与黑白名单无关
- 黑白名单管理的是目标服务器IP（通过domaindns解析得到）

### 5. 连接配置 (conn)

控制网络连接的基本参数。

**关键配置项**：
- `connect_timeout`: 连接超时（秒）
- `read_timeout`: 读取超时（秒）
- `write_timeout`: 写入超时（秒）
- `keep_alive`: 是否启用Keep-Alive
- `max_idle_conns`: 最大空闲连接数

### 6. TCP连接池配置 (netconnpool)

控制TCP连接池行为。

**关键配置项**：
- `max_connections`: 最大连接数
- `initial_connections`: 初始连接数
- `idle_timeout`: 连接空闲超时（秒）
- `max_lifetime`: 连接最大生存时间（秒）

### 7. QUIC连接池配置 (quic)

控制QUIC连接池行为。

**关键配置项**：
- `max_connections`: 最大连接数
- `enable_0rtt`: 是否启用0-RTT
- `handshake_timeout`: QUIC握手超时（秒）

### 8. 证书配置 (certificate)

控制TLS证书的申请和管理。

**关键配置项**：
- `server_domain`: 服务端域名（VPS域名，用于QUIC服务端）
- `provider`: 证书提供商（letsencrypt, self-signed）
- `auto_renewal`: 是否自动续期
- `letsencrypt_email`: Let's Encrypt邮箱
- `auto_detect_local_ip`: 是否自动检测本地IP并加入证书

### 9. 黑白名单配置 (ip_status)

控制目标服务器IP的黑白名单管理。

**关键配置项**：
- `min_whitelist_count`: 白名单最小数量（低于此值告警）
- `allow_start_when_empty`: 白名单为空时是否允许启动
- `whitelist_monitoring`: 是否启用白名单监控

**重要说明**：
- 黑白名单管理的是目标服务器的IP地址（通过domaindns解析得到）
- 不是本地IP池（localippool）的IP

### 10. IP池测试配置 (ip_pool_test)

控制目标服务器IP的测试行为。

**关键配置项**：
- `target_domains`: 目标域名列表（系统会解析这些域名，测试解析出的IP）
- `test_url`: 测试URL（用于测试IP可用性，{domain}会被替换为实际域名）
- `test_method`: 测试方法（GET, HEAD）
- `max_concurrent`: 最大并发测试数
- `test_timeout`: 测试超时（秒）
- `success_status_codes`: 测试成功状态码列表
- `forbidden_status_codes`: 测试失败状态码列表（会被加入黑名单）

### 11. 黑名单恢复配置 (blacklist_recovery)

控制黑名单IP的自动恢复行为。

**关键配置项**：
- `enabled`: 是否启用黑名单恢复
- `check_interval`: 检查间隔（秒）
- `ip_test_interval`: 每个IP的测试间隔（秒，避免频繁测试）
- `max_concurrent`: 最大并发恢复测试数

### 12. 状态报告配置 (status_report)

控制服务器端状态报告行为。

**关键配置项**：
- `report_interval`: 报告间隔（秒）
- `report_on_change`: 状态变化时是否立即报告
- `report_ip_list`: 是否报告IP列表
- `max_report_ips`: 最大报告IP数量

### 13. 服务端配置 (server)

控制QUIC服务端行为。

**关键配置项**：
- `listen_address`: QUIC服务端监听地址
- `quic_enabled`: 是否启用QUIC服务端
- `max_clients`: 最大客户端连接数
- `client_timeout`: 客户端连接超时（秒）

### 14. 爬虫配置 (crawler)

控制爬虫请求行为。

**关键配置项**：
- `default_timeout`: 默认请求超时（秒）
- `max_retries`: 最大重试次数
- `protocol_priority`: 协议优先级（http3, http2, http1.1）
- `protocol_fallback`: 是否启用协议降级
- `concurrency`: 并发请求数
- `rate_limit`: 请求速率限制（每秒请求数，0表示不限制）

### 15. 系统配置 (system)

控制系统基本信息。

**关键配置项**：
- `name`: 系统名称
- `version`: 系统版本
- `work_dir`: 工作目录
- `data_dir`: 数据目录
- `health_check_enabled`: 是否启用健康检查
- `metrics_enabled`: 是否启用指标收集

## 配置加载

使用 `config.go` 中的 `LoadConfig()` 函数加载配置文件：

```go
import (
    "github.com/pelletier/go-toml/v2"
    "os"
)

func LoadConfig(path string) (*SystemConfig, error) {
    config := &SystemConfig{}
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    err = toml.Unmarshal(data, config)
    if err != nil {
        return nil, err
    }
    return config, nil
}
```

## 配置验证

在加载配置后，应该验证关键配置项：

1. **必需配置项检查**：
   - 证书域名（如果启用QUIC服务端）
   - 目标域名列表（如果启用IP池测试）
   - DNS服务器列表

2. **配置合理性检查**：
   - 超时时间应该大于0
   - 连接池大小应该合理
   - 重试次数应该合理

## 配置示例

### 最小配置示例

```toml
[logs]
level = "info"

[domaindns]
dns_servers = ["8.8.8.8", "1.1.1.1"]

[ip_pool_test]
target_domains = ["example.com"]

[certificate]
server_domain = "crawler.example.com"
provider = "letsencrypt"
```

### 完整配置示例

参考 `config.toml` 文件中的完整配置示例。

## 注意事项

1. **IP池理解**：
   - `local_ip_pool.ips`: 本地出口IP列表（用于绑定本地连接）
   - `ip_pool_test.target_domains`: 目标域名列表（解析后测试IP，加入黑白名单）

2. **证书配置**：
   - 如果使用Let's Encrypt，需要确保域名可以访问
   - 证书会自动检测本地IP并加入证书

3. **黑白名单**：
   - 黑白名单管理的是目标服务器IP，不是本地IP
   - 系统启动时会测试目标域名解析出的IP

4. **配置热更新**：
   - 部分配置支持热更新（如IP池测试的目标域名）
   - 部分配置需要重启才能生效（如服务端监听地址）

