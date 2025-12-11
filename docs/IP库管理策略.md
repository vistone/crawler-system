# IP库管理策略详解

## 设计目标

1. **减少DNS解析**：通过IP库缓存，避免频繁DNS查询
2. **提高性能**：直接使用IP连接，减少域名解析延迟
3. **支持降级**：IP不可用时自动降级到域名
4. **健康管理**：自动检测和维护IP健康状态

## IP库架构设计

### 1. 存储结构

```go
// IPLibrary IP库核心结构
type IPLibrary struct {
    // 主索引：域名 -> IP列表
    domainIndex *DomainIndex
    
    // 反向索引：IP -> 域名列表
    ipIndex *IPIndex
    
    // IP健康状态索引
    healthIndex *HealthIndex
    
    // DNS记录缓存
    dnsCache *DNSCache
    
    // 统计信息
    stats *IPLibraryStats
}

// DomainIndex 域名索引
type DomainIndex struct {
    // 域名 -> IP记录列表
    records map[string]*DomainRecord
    
    // 读写锁
    mutex sync.RWMutex
}

// DomainRecord 域名记录
type DomainRecord struct {
    Domain      string
    IPs         []*IPRecord
    LastUpdated time.Time
    UpdateCount int64
    TTL         time.Duration
    ExpiresAt   time.Time
}

// IPRecord IP记录
type IPRecord struct {
    IP          string
    Port        int
    Protocols   []string      // 支持的协议
    IsHealthy   bool
    SuccessRate float64
    Latency     time.Duration
    UseCount    int64
    LastUsed    time.Time
    LastChecked time.Time
    CreatedAt   time.Time
    TTL         time.Duration
    ExpiresAt   time.Time
    Metadata    map[string]interface{}
}

// IPIndex IP反向索引
type IPIndex struct {
    // IP -> 域名列表
    domains map[string][]string
    mutex   sync.RWMutex
}

// HealthIndex 健康状态索引
type HealthIndex struct {
    // IP -> 健康状态
    status map[string]*HealthStatus
    mutex  sync.RWMutex
}

// HealthStatus 健康状态
type HealthStatus struct {
    IP            string
    IsHealthy     bool
    IsBlocked     bool
    LastChecked   time.Time
    CheckInterval time.Duration
    ErrorCount    int
    SuccessCount  int
    Protocols     map[string]*ProtocolHealth
}

// ProtocolHealth 协议健康状态
type ProtocolHealth struct {
    Protocol    string
    IsSupported bool
    SuccessRate float64
    Latency     time.Duration
    LastChecked time.Time
}
```

### 2. IP库操作接口

```go
// IPLibraryManager IP库管理器
type IPLibraryManager interface {
    // 获取IP
    GetIPs(domain string, protocol string) ([]*IPRecord, error)
    GetBestIP(domain string, protocol string) (*IPRecord, error)
    GetIPsWithFallback(domain string, protocol string) ([]*IPRecord, error)
    
    // 更新IP
    UpdateIPs(domain string) error
    AddIP(domain string, ip string, port int) error
    RemoveIP(domain string, ip string) error
    
    // 健康管理
    MarkIPHealthy(ip string, protocol string)
    MarkIPBlocked(ip string, protocol string)
    MarkIPError(ip string, protocol string, err error)
    
    // 统计
    GetStats() *IPLibraryStats
    GetDomainStats(domain string) *DomainStats
    GetIPStats(ip string) *IPStats
}
```

## IP库更新策略

### 1. 更新触发机制

```go
// UpdateTrigger 更新触发器
type UpdateTrigger struct {
    // 定时更新
    TimeBased *TimeBasedTrigger
    
    // 失败触发
    FailureBased *FailureBasedTrigger
    
    // 使用频率触发
    UsageBased *UsageBasedTrigger
    
    // 手动触发
    Manual chan struct{}
}

// TimeBasedTrigger 定时触发器
type TimeBasedTrigger struct {
    Interval time.Duration
    Ticker   *time.Ticker
}

// FailureBasedTrigger 失败触发器
type FailureBasedTrigger struct {
    FailureThreshold int
    FailureWindow    time.Duration
}

// UsageBasedTrigger 使用频率触发器
type UsageBasedTrigger struct {
    UsageThreshold int64
    UsageWindow     time.Duration
}
```

### 2. 更新策略实现

```go
// UpdateStrategy 更新策略
type UpdateStrategy struct {
    // 更新方法
    Method UpdateMethod
    
    // 并发控制
    MaxConcurrent int
    
    // 重试配置
    RetryConfig *RetryConfig
}

// UpdateMethod 更新方法
type UpdateMethod string

const (
    // 全量更新：更新所有域名的IP
    MethodFull UpdateMethod = "full"
    
    // 增量更新：只更新变化的域名
    MethodIncremental UpdateMethod = "incremental"
    
    // 按需更新：只更新请求的域名
    MethodOnDemand UpdateMethod = "on_demand"
    
    // 混合更新：结合多种方法
    MethodHybrid UpdateMethod = "hybrid"
)

// 更新实现
func (m *IPLibraryManager) UpdateIPs(domain string) error {
    // 1. 解析域名
    ips, err := m.resolver.Resolve(context.Background(), domain)
    if err != nil {
        return err
    }
    
    // 2. 获取现有记录
    record, exists := m.library.domainIndex.Get(domain)
    
    // 3. 更新IP记录
    newIPs := make([]*IPRecord, 0, len(ips))
    for _, ip := range ips {
        ipRecord := m.createOrUpdateIPRecord(domain, ip, record)
        newIPs = append(newIPs, ipRecord)
    }
    
    // 4. 更新域名记录
    m.library.domainIndex.Update(domain, &DomainRecord{
        Domain:      domain,
        IPs:         newIPs,
        LastUpdated: time.Now(),
        UpdateCount: record.UpdateCount + 1,
        TTL:         m.config.DefaultTTL,
        ExpiresAt:   time.Now().Add(m.config.DefaultTTL),
    })
    
    // 5. 更新反向索引
    m.updateIPIndex(domain, newIPs)
    
    return nil
}
```

## IP选择策略

### 1. 选择算法

```go
// IPSelector IP选择器
type IPSelector struct {
    // 选择策略
    Strategy SelectionStrategy
    
    // 权重配置
    Weights *SelectionWeights
}

// SelectionStrategy 选择策略
type SelectionStrategy string

const (
    // 随机选择
    StrategyRandom SelectionStrategy = "random"
    
    // 轮询选择
    StrategyRoundRobin SelectionStrategy = "round_robin"
    
    // 最少使用
    StrategyLeastUsed SelectionStrategy = "least_used"
    
    // 最高成功率
    StrategyHighestSuccess StrategySelectionStrategy = "highest_success"
    
    // 最低延迟
    StrategyLowestLatency SelectionStrategy = "lowest_latency"
    
    // 加权选择
    StrategyWeighted SelectionStrategy = "weighted"
)

// SelectionWeights 选择权重
type SelectionWeights struct {
    SuccessRate float64 // 成功率权重
    Latency     float64 // 延迟权重
    UseCount    float64 // 使用次数权重
    Health      float64 // 健康状态权重
}

// SelectBestIP 选择最佳IP
func (s *IPSelector) SelectBestIP(
    ips []*IPRecord,
    protocol string,
) *IPRecord {
    // 过滤健康的IP
    healthyIPs := s.filterHealthyIPs(ips, protocol)
    if len(healthyIPs) == 0 {
        return nil
    }
    
    // 根据策略选择
    switch s.Strategy {
    case StrategyRandom:
        return s.selectRandom(healthyIPs)
    case StrategyRoundRobin:
        return s.selectRoundRobin(healthyIPs)
    case StrategyLeastUsed:
        return s.selectLeastUsed(healthyIPs)
    case StrategyHighestSuccess:
        return s.selectHighestSuccess(healthyIPs)
    case StrategyLowestLatency:
        return s.selectLowestLatency(healthyIPs)
    case StrategyWeighted:
        return s.selectWeighted(healthyIPs, s.Weights)
    default:
        return healthyIPs[0]
    }
}

// selectWeighted 加权选择
func (s *IPSelector) selectWeighted(
    ips []*IPRecord,
    weights *SelectionWeights,
) *IPRecord {
    best := ips[0]
    bestScore := s.calculateScore(best, weights)
    
    for _, ip := range ips[1:] {
        score := s.calculateScore(ip, weights)
        if score > bestScore {
            best = ip
            bestScore = score
        }
    }
    
    return best
}

// calculateScore 计算IP得分
func (s *IPSelector) calculateScore(
    ip *IPRecord,
    weights *SelectionWeights,
) float64 {
    score := 0.0
    
    // 成功率得分
    score += ip.SuccessRate * weights.SuccessRate
    
    // 延迟得分（延迟越低得分越高）
    latencyScore := 1.0 / (1.0 + float64(ip.Latency)/float64(time.Second))
    score += latencyScore * weights.Latency
    
    // 使用次数得分（使用越少得分越高）
    useScore := 1.0 / (1.0 + float64(ip.UseCount))
    score += useScore * weights.UseCount
    
    // 健康状态得分
    if ip.IsHealthy {
        score += weights.Health
    }
    
    return score
}
```

## 健康检查策略

### 1. 健康检查配置

```go
// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
    // 检查间隔
    Interval time.Duration
    
    // 检查超时
    Timeout time.Duration
    
    // 检查方法
    Method HealthCheckMethod
    
    // 健康阈值
    HealthyThreshold float64
    
    // 失败阈值
    FailureThreshold int
}

// HealthCheckMethod 健康检查方法
type HealthCheckMethod string

const (
    // TCP连接检查
    MethodTCP HealthCheckMethod = "tcp"
    
    // HTTP HEAD请求
    MethodHTTPHead HealthCheckMethod = "http_head"
    
    // HTTP GET请求
    MethodHTTPGet HealthCheckMethod = "http_get"
    
    // 自定义检查
    MethodCustom HealthCheckMethod = "custom"
)
```

### 2. 健康检查实现

```go
// HealthChecker 健康检查器
type HealthChecker struct {
    config  *HealthCheckConfig
    library *IPLibrary
    logger  *logs.Logger
}

// CheckIP 检查IP健康状态
func (c *HealthChecker) CheckIP(
    ctx context.Context,
    ip string,
    protocol string,
) bool {
    switch c.config.Method {
    case MethodTCP:
        return c.checkTCP(ctx, ip)
    case MethodHTTPHead:
        return c.checkHTTPHead(ctx, ip, protocol)
    case MethodHTTPGet:
        return c.checkHTTPGet(ctx, ip, protocol)
    case MethodCustom:
        return c.checkCustom(ctx, ip, protocol)
    default:
        return c.checkTCP(ctx, ip)
    }
}

// checkTCP TCP连接检查
func (c *HealthChecker) checkTCP(
    ctx context.Context,
    ip string,
) bool {
    conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "443"), c.config.Timeout)
    if err != nil {
        return false
    }
    conn.Close()
    return true
}

// checkHTTPHead HTTP HEAD检查
func (c *HealthChecker) checkHTTPHead(
    ctx context.Context,
    ip string,
    protocol string,
) bool {
    // 实现HTTP HEAD请求检查
    return true
}
```

## IP库持久化

### 1. 持久化策略

```go
// PersistenceStrategy 持久化策略
type PersistenceStrategy struct {
    // 存储后端
    Backend StorageBackend
    
    // 持久化间隔
    Interval time.Duration
    
    // 持久化触发
    Trigger PersistenceTrigger
}

// StorageBackend 存储后端
type StorageBackend string

const (
    // 内存（不持久化）
    BackendMemory StorageBackend = "memory"
    
    // 文件系统
    BackendFile StorageBackend = "file"
    
    // 数据库
    BackendDatabase StorageBackend = "database"
    
    // Redis
    BackendRedis StorageBackend = "redis"
)

// IPLibraryStorage IP库存储接口
type IPLibraryStorage interface {
    // 保存
    Save(library *IPLibrary) error
    
    // 加载
    Load() (*IPLibrary, error)
    
    // 增量更新
    Update(domain string, record *DomainRecord) error
}
```

### 2. 文件存储实现

```go
// FileStorage 文件存储
type FileStorage struct {
    filePath string
    logger   *logs.Logger
}

// Save 保存到文件
func (s *FileStorage) Save(library *IPLibrary) error {
    data, err := json.Marshal(library)
    if err != nil {
        return err
    }
    
    return os.WriteFile(s.filePath, data, 0644)
}

// Load 从文件加载
func (s *FileStorage) Load() (*IPLibrary, error) {
    data, err := os.ReadFile(s.filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return NewIPLibrary(), nil
        }
        return nil, err
    }
    
    var library IPLibrary
    if err := json.Unmarshal(data, &library); err != nil {
        return nil, err
    }
    
    return &library, nil
}
```

## 性能优化

### 1. 缓存优化

- **多级缓存**：内存缓存 + 持久化缓存
- **LRU缓存**：最近使用的域名优先
- **预加载**：启动时预加载常用域名

### 2. 并发优化

- **读写分离**：使用读写锁
- **批量更新**：批量更新减少锁竞争
- **异步更新**：后台异步更新IP库

### 3. 内存优化

- **压缩存储**：IP地址压缩存储
- **过期清理**：定期清理过期IP
- **分片存储**：大库分片存储

## 使用示例

```go
// 初始化IP库管理器
ipLibrary := iplibrary.NewIPLibraryManager(
    dnsResolver,
    logger,
    &iplibrary.IPLibraryConfig{
        UpdateInterval:       1 * time.Hour,
        IPExpiration:         24 * time.Hour,
        HealthCheckInterval:  5 * time.Minute,
        MaxIPsPerDomain:      10,
        MinHealthyIPs:        3,
    },
)

// 获取IP
ips, err := ipLibrary.GetIPs("example.com", "http2")
if err != nil {
    // 降级到域名
    target = "example.com"
} else {
    // 使用IP
    target = ips[0].IP
}

// 标记IP状态
ipLibrary.MarkIPHealthy(ip, "http2")
ipLibrary.MarkIPBlocked(ip, "http2")
```

