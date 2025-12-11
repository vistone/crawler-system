# 系统初始化说明

## 概述

系统已经实现了9个核心模块的初始化框架，所有模块都可以通过配置文件进行初始化。

## 初始化顺序

系统按照以下顺序初始化9个模块：

```
1. logs (日志系统) - 最先初始化，其他模块都需要
   ↓
2. fingerprint (指纹模块) - 独立模块
3. domaindns (DNS解析模块) - 独立模块
4. localippool (本地IP池模块) - 独立模块
5. certs (证书模块) - 独立模块
6. whitelist-blacklist-manager (黑白名单模块) - 独立模块
   ↓
7. conn (连接模块) - 依赖前面模块
   ↓
8. netconnpool (TCP连接池模块) - 依赖conn模块
9. quic (QUIC连接池模块) - 依赖conn模块
```

## 使用方法

### 1. 基本使用

```go
package main

import (
    "log"
)

func main() {
    // 创建系统实例（会自动加载配置并初始化所有模块）
    system, err := NewSystem("config.toml")
    if err != nil {
        log.Fatal(err)
    }
    defer system.Close()

    // 系统已初始化完成，可以使用各个模块
    // system.Logger.Info("系统运行中...")
    // system.DNSResolver.Resolve("example.com")
    // ...
}
```

### 2. 命令行使用

```bash
# 使用默认配置文件
go run main.go

# 指定配置文件
go run main.go -config /path/to/config.toml
```

## 模块初始化详情

### 模块1: logs (日志系统)

**初始化函数**: `initLogs()`

**配置项**: `[logs]`

**功能**:
- 设置日志级别
- 配置日志输出（文件/标准输出）
- 设置日志格式（json/text）

**当前实现**: 占位符实现（输出到标准输出）

**TODO**: 集成真实的logs库

### 模块2: fingerprint (指纹模块)

**初始化函数**: `initFingerprint()`

**配置项**: `[fingerprint]`

**功能**:
- 设置指纹选择策略
- 配置指纹轮换
- 管理指纹库

**当前实现**: 占位符实现

**TODO**: 集成真实的fingerprint库

### 模块3: domaindns (DNS解析模块)

**初始化函数**: `initDomainDNS()`

**配置项**: `[domaindns]`

**功能**:
- 设置DNS服务器
- 配置DNS缓存
- 设置DNS查询超时和重试

**当前实现**: 占位符实现

**TODO**: 集成真实的domaindns库

### 模块4: localippool (本地IP池模块)

**初始化函数**: `initLocalIPPool()`

**配置项**: `[local_ip_pool]`

**功能**:
- 添加本地出口IP
- 设置IP选择策略
- 配置IP健康检查

**当前实现**: 占位符实现（内存存储）

**TODO**: 集成真实的localippool库

### 模块5: certs (证书模块)

**初始化函数**: `initCerts()`

**配置项**: `[certificate]`

**功能**:
- 配置证书提供商（Let's Encrypt/自签名）
- 设置证书存储路径
- 配置证书自动续期

**当前实现**: 占位符实现

**TODO**: 集成真实的certs库

### 模块6: whitelist-blacklist-manager (黑白名单模块)

**初始化函数**: `initIPStatusManager()`

**配置项**: `[ip_status]`

**功能**:
- 管理目标服务器IP的黑白名单
- 设置白名单最小数量
- 配置白名单监控

**当前实现**: 占位符实现（内存存储）

**TODO**: 集成真实的whitelist-blacklist-manager库

### 模块7: conn (连接模块)

**初始化函数**: `initConn()`

**配置项**: `[conn]`

**功能**:
- 设置连接超时
- 配置Keep-Alive
- 设置TLS参数

**当前实现**: 占位符实现

**TODO**: 集成真实的conn库

### 模块8: netconnpool (TCP连接池模块)

**初始化函数**: `initNetConnPool()`

**配置项**: `[netconnpool]`

**功能**:
- 设置连接池大小
- 配置连接生命周期
- 设置健康检查

**依赖**: conn, fingerprint, domaindns, localippool

**当前实现**: 占位符实现

**TODO**: 集成真实的netconnpool库

### 模块9: quic (QUIC连接池模块)

**初始化函数**: `initQUICPool()`

**配置项**: `[quic]`

**功能**:
- 设置QUIC连接池大小
- 配置0-RTT
- 设置握手超时

**依赖**: conn, fingerprint, domaindns, localippool

**当前实现**: 占位符实现

**TODO**: 集成真实的quic库

## 接口设计

为了便于后续替换实际实现，所有模块都使用接口：

- `LoggerInterface` - 日志接口
- `DNSResolverInterface` - DNS解析接口
- `LocalIPPoolInterface` - 本地IP池接口
- `NetConnPoolInterface` - TCP连接池接口
- `QUICPoolInterface` - QUIC连接池接口
- `CertManagerInterface` - 证书管理接口
- `IPStatusManagerInterface` - 黑白名单管理接口

## 占位符实现

当前所有模块都使用占位符实现，这些实现：

1. **满足接口要求**：所有占位符实现都实现了对应的接口
2. **存储配置**：保存了配置信息，便于后续替换
3. **基本功能**：实现了基本的接口方法
4. **TODO标记**：标记了需要集成真实库的位置

## 下一步

1. **安装依赖库**：
   ```bash
   go get github.com/pelletier/go-toml/v2
   go get github.com/vistone/logs
   go get github.com/vistone/fingerprint
   # ... 其他库
   ```

2. **实现TOML解析**：
   - 在`config.go`中实现真实的TOML解析
   - 使用`github.com/pelletier/go-toml/v2`

3. **替换占位符实现**：
   - 逐步替换各个模块的占位符实现
   - 使用真实的库API

4. **测试初始化**：
   - 测试每个模块的初始化
   - 验证配置加载正确性

## 注意事项

1. **初始化顺序很重要**：必须按照依赖关系顺序初始化
2. **错误处理**：如果某个模块初始化失败，整个系统初始化会失败
3. **资源清理**：使用`system.Close()`确保资源正确释放
4. **配置验证**：在初始化前应该验证关键配置项

