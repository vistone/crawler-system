# 爬虫反审查系统 - API接口设计

## 核心接口定义

### 1. FingerprintManager 接口

```go
package fingerprint

// FingerprintManager 指纹管理器接口
type FingerprintManager interface {
    // GetRandomFingerprint 获取随机指纹
    GetRandomFingerprint() (*FingerprintResult, error)
    
    // GetRandomFingerprintWithOS 获取随机指纹（指定操作系统）
    GetRandomFingerprintWithOS(os OperatingSystem) (*FingerprintResult, error)
    
    // GetRandomFingerprintByBrowser 根据浏览器类型获取随机指纹
    GetRandomFingerprintByBrowser(browserType string) (*FingerprintResult, error)
    
    // GetRandomFingerprintByBrowserWithOS 根据浏览器类型和操作系统获取随机指纹
    GetRandomFingerprintByBrowserWithOS(
        browserType string,
        os OperatingSystem,
    ) (*FingerprintResult, error)
    
    // GetFingerprintByName 根据名称获取指纹
    GetFingerprintByName(name string) (*FingerprintResult, error)
}

// FingerprintResult 指纹结果
type FingerprintResult struct {
    Profile   ClientProfile  // TLS指纹配置
    UserAgent string         // User-Agent
    Name      string         // 指纹名称
    Headers   *HTTPHeaders   // HTTP请求头
}
```

### 2. IPPool 接口

```go
package localippool

// IPPool IP池接口
type IPPool interface {
    // AddIP 添加IP地址
    AddIP(ip string) error
    
    // RemoveIP 移除IP地址
    RemoveIP(ip string) error
    
    // GetAvailableIP 获取可用IP
    GetAvailableIP() (string, error)
    
    // GetAvailableIPWithStrategy 根据策略获取IP
    GetAvailableIPWithStrategy(strategy IPStrategy) (string, error)
    
    // MarkHealthy 标记IP为健康
    MarkHealthy(ip string)
    
    // MarkSuspicious 标记IP为可疑
    MarkSuspicious(ip string)
    
    // MarkBlocked 标记IP为被封锁
    MarkBlocked(ip string)
    
    // GetIPStatus 获取IP状态
    GetIPStatus(ip string) (IPStatus, error)
    
    // Count 获取IP数量
    Count() int
    
    // ListIPs 列出所有IP
    ListIPs() []string
}

// IPStrategy IP选择策略
type IPStrategy string

const (
    StrategyRoundRobin IPStrategy = "round_robin"  // 轮询
    StrategyRandom     IPStrategy = "random"       // 随机
    StrategyLeastUsed  IPStrategy = "least_used"   // 最少使用
    StrategyHealthy    IPStrategy = "healthy"       // 仅健康IP
)

// IPStatus IP状态
type IPStatus struct {
    IP            string
    IsHealthy     bool
    IsBlocked     bool
    UseCount      int
    LastUsed      time.Time
    LastChecked   time.Time
    ErrorCount    int
}
```

### 3. DNSResolver 接口

```go
package domaindns

// DNSResolver DNS解析器接口
type DNSResolver interface {
    // Resolve 解析域名
    Resolve(ctx context.Context, domain string) (string, error)
    
    // ResolveWithType 根据类型解析域名
    ResolveWithType(ctx context.Context, domain string, recordType string) ([]string, error)
    
    // SetDNSServers 设置DNS服务器
    SetDNSServers(servers []string)
    
    // GetDNSServers 获取DNS服务器列表
    GetDNSServers() []string
    
    // EnableCache 启用/禁用缓存
    EnableCache(enabled bool)
    
    // SetCacheTTL 设置缓存TTL
    SetCacheTTL(ttl time.Duration)
    
    // ClearCache 清空缓存
    ClearCache()
    
    // SetTimeout 设置解析超时
    SetTimeout(timeout time.Duration)
    
    // SetRetryCount 设置重试次数
    SetRetryCount(count int)
}

// NewResolver 创建DNS解析器
func NewResolver() DNSResolver {
    // 实现
}
```

### 4. Connection 接口

```go
package conn

// Connection 连接接口
type Connection interface {
    // SendRequest 发送HTTP请求
    SendRequest(ctx context.Context, req *HTTPRequest) (*HTTPResponse, error)
    
    // IsHealthy 检查连接是否健康
    IsHealthy() bool
    
    // Close 关闭连接
    Close() error
    
    // GetLocalAddr 获取本地地址
    GetLocalAddr() net.Addr
    
    // GetRemoteAddr 获取远程地址
    GetRemoteAddr() net.Addr
    
    // GetStats 获取连接统计信息
    GetStats() *ConnectionStats
}

// HTTPRequest HTTP请求
type HTTPRequest struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    []byte
}

// HTTPResponse HTTP响应
type HTTPResponse struct {
    StatusCode int
    Headers    map[string]string
    Body       []byte
}

// ConnectionStats 连接统计信息
type ConnectionStats struct {
    RequestCount  int
    ResponseCount int
    ErrorCount    int
    BytesSent     int64
    BytesReceived int64
    CreatedAt     time.Time
    LastUsedAt    time.Time
}
```

### 5. ConnectionPool 接口

```go
package netconnpool

// ConnectionPool 连接池接口
type ConnectionPool interface {
    // GetConnection 获取连接
    GetConnection(
        ctx context.Context,
        opts *ConnectionOptions,
    ) (conn.Connection, error)
    
    // PutConnection 归还连接
    PutConnection(conn conn.Connection)
    
    // CloseConnection 关闭连接
    CloseConnection(conn conn.Connection)
    
    // Close 关闭连接池
    Close() error
    
    // GetStats 获取连接池统计信息
    GetStats() *PoolStats
}

// ConnectionOptions 连接选项
type ConnectionOptions struct {
    TargetIP    string
    LocalIP     string
    Fingerprint fingerprint.ClientProfile
    Headers     map[string]string
    Timeout     time.Duration
}

// PoolStats 连接池统计信息
type PoolStats struct {
    TotalConnections    int
    ActiveConnections   int
    IdleConnections     int
    WaitingRequests     int
    TotalRequests       int64
    SuccessfulRequests  int64
    FailedRequests      int64
}
```

### 6. QUICPool 接口

```go
package quic

// QUICPool QUIC连接池接口
type QUICPool interface {
    // GetConnection 获取QUIC连接
    GetConnection(
        ctx context.Context,
        opts *QUICConnectionOptions,
    ) (conn.Connection, error)
    
    // PutConnection 归还连接
    PutConnection(conn conn.Connection)
    
    // CloseConnection 关闭连接
    CloseConnection(conn conn.Connection)
    
    // Close 关闭连接池
    Close() error
    
    // GetStats 获取连接池统计信息
    GetStats() *QUICPoolStats
}

// QUICConnectionOptions QUIC连接选项
type QUICConnectionOptions struct {
    TargetIP    string
    LocalIP     string
    Fingerprint fingerprint.ClientProfile
    Headers     map[string]string
    MaxStreams  int
    Timeout     time.Duration
}

// QUICPoolStats QUIC连接池统计信息
type QUICPoolStats struct {
    TotalConnections    int
    ActiveConnections   int
    IdleConnections     int
    TotalStreams        int
    ActiveStreams       int
    TotalRequests       int64
    SuccessfulRequests  int64
    FailedRequests      int64
}
```

### 7. Logger 接口

```go
package logs

// Logger 日志接口
type Logger interface {
    // Debug 记录调试日志
    Debug(msg string, fields ...Field)
    
    // Info 记录信息日志
    Info(msg string, fields ...Field)
    
    // Warn 记录警告日志
    Warn(msg string, fields ...Field)
    
    // Error 记录错误日志
    Error(msg string, fields ...Field)
    
    // Fatal 记录致命错误日志并退出
    Fatal(msg string, fields ...Field)
    
    // WithFields 添加字段
    WithFields(fields ...Field) Logger
    
    // SetLevel 设置日志级别
    SetLevel(level Level)
}

// Level 日志级别
type Level int

const (
    DebugLevel Level = iota
    InfoLevel
    WarnLevel
    ErrorLevel
    FatalLevel
)

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// NewField 创建日志字段
func NewField(key string, value interface{}) Field {
    return Field{Key: key, Value: value}
}
```

## 业务整合层接口

### CrawlerSystem 接口

```go
package crawler

// CrawlerSystem 爬虫系统接口
type CrawlerSystem interface {
    // ExecuteRequest 执行HTTP请求
    ExecuteRequest(ctx context.Context, req *Request) (*Response, error)
    
    // ExecuteRequestWithRetry 执行HTTP请求（带重试）
    ExecuteRequestWithRetry(
        ctx context.Context,
        req *Request,
        retryOpts *RetryOptions,
    ) (*Response, error)
    
    // GetStats 获取系统统计信息
    GetStats() *SystemStats
    
    // Close 关闭系统
    Close() error
}

// Request 请求结构
type Request struct {
    URL         string
    Method      string
    Headers     map[string]string
    Body        []byte
    UseQUIC     bool
    BrowserType string
    Timeout     time.Duration
}

// Response 响应结构
type Response struct {
    StatusCode int
    Headers    map[string]string
    Body       []byte
    Latency    time.Duration
}

// RetryOptions 重试选项
type RetryOptions struct {
    MaxRetries int
    RetryDelay time.Duration
    RetryOn    []int // HTTP状态码，遇到这些状态码时重试
}

// SystemStats 系统统计信息
type SystemStats struct {
    TotalRequests      int64
    SuccessfulRequests int64
    FailedRequests     int64
    AverageLatency     time.Duration
    IPPoolStats        *IPPoolStats
    ConnectionPoolStats *ConnectionPoolStats
    DNSStats           *DNSStats
}

// IPPoolStats IP池统计信息
type IPPoolStats struct {
    TotalIPs    int
    HealthyIPs  int
    BlockedIPs  int
    SuspiciousIPs int
}

// ConnectionPoolStats 连接池统计信息
type ConnectionPoolStats struct {
    TCPPool  *netconnpool.PoolStats
    QUICPool *quic.QUICPoolStats
}

// DNSStats DNS统计信息
type DNSStats struct {
    TotalResolves     int64
    SuccessfulResolves int64
    FailedResolves    int64
    CacheHits         int64
    CacheMisses       int64
}
```

## 配置接口

### SystemConfig 系统配置

```go
package crawler

// SystemConfig 系统配置
type SystemConfig struct {
    // 日志配置
    LogLevel logs.Level
    LogFile  string
    
    // IP池配置
    IPPoolConfig IPPoolConfig
    
    // DNS配置
    DNSConfig DNSConfig
    
    // 连接池配置
    TCPPoolConfig  TCPPoolConfig
    QUICPoolConfig QUICPoolConfig
    
    // 指纹配置
    FingerprintConfig FingerprintConfig
    
    // 请求配置
    RequestConfig RequestConfig
}

// IPPoolConfig IP池配置
type IPPoolConfig struct {
    IPs              []string
    HealthCheckInterval time.Duration
    RateLimit         int
    RateLimitWindow   time.Duration
    Strategy          localippool.IPStrategy
}

// DNSConfig DNS配置
type DNSConfig struct {
    DNSServers []string
    CacheEnabled bool
    CacheTTL     time.Duration
    Timeout      time.Duration
    RetryCount   int
}

// TCPPoolConfig TCP连接池配置
type TCPPoolConfig struct {
    MaxConnections     int
    MaxIdleConnections int
    IdleTimeout        time.Duration
    ConnectTimeout     time.Duration
    ReadTimeout        time.Duration
    WriteTimeout       time.Duration
}

// QUICPoolConfig QUIC连接池配置
type QUICPoolConfig struct {
    MaxConnections     int
    MaxIdleConnections int
    IdleTimeout        time.Duration
    ConnectTimeout     time.Duration
    MaxStreams         int
    KeepAliveInterval  time.Duration
}

// FingerprintConfig 指纹配置
type FingerprintConfig struct {
    BrowserTypes      []string
    OperatingSystems  []fingerprint.OperatingSystem
    RotationStrategy  string
}

// RequestConfig 请求配置
type RequestConfig struct {
    DefaultTimeout    time.Duration
    DefaultRetryCount int
    DefaultRetryDelay time.Duration
    UserAgentRotation bool
}
```

