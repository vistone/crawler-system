# 爬虫反审查系统

[![Go Version](https://img.shields.io/badge/go-1.25.4-blue.svg)](https://golang.org)
[![Go Dev](https://pkg.go.dev/badge/github.com/vistone/crawler-system.svg)](https://pkg.go.dev/github.com/vistone/crawler-system)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-v1.0.0-green.svg)](https://github.com/vistone/crawler-system/releases/tag/v1.0.0)

> 基于9个核心库构建的精细化控制反审查爬虫系统，支持双重角色（QUIC服务器 + 爬虫客户端），具备自动IP管理、证书管理、实时状态同步等核心功能。

## 📋 项目概述

本项目是一个部署在公网VPS上的反审查爬虫系统，同时扮演**服务器**和**爬虫客户端**两个角色：

- **服务器角色**：提供QUIC服务端，接收客户端连接，实时报告系统状态
- **爬虫客户端角色**：执行爬虫任务，管理目标服务器IP的黑白名单，处理反爬虫机制

系统通过多层架构设计，实现了指纹模拟、IP管理、DNS解析、连接池管理、证书管理、黑白名单管理等核心功能，能够有效规避目标网站的反爬虫机制。

## 🎯 核心特性

- ✅ **多协议支持**：HTTP/1.1、HTTP/2、HTTP/3（QUIC）智能选择和降级
- ✅ **TLS指纹模拟**：支持44+主流浏览器指纹，随机轮换
- ✅ **智能IP管理**：自动测试目标服务器IP，黑白名单管理，黑名单自动恢复
- ✅ **DNS解析优化**：自定义DNS服务器，DNS缓存，规避污染
- ✅ **连接池复用**：TCP、HTTP/2和QUIC连接池，提高性能
- ✅ **证书自动管理**：Let's Encrypt自动申请和续期
- ✅ **实时状态同步**：服务器端实时向客户端报告系统状态
- ✅ **高可用性**：白名单为空时进入待机状态，自动恢复机制

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    应用层 (Application)                   │
│  - 爬虫逻辑 - 请求调度 - 数据解析 - 反审查策略           │
│  - QUIC服务端（接收客户端请求）                          │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   业务整合层 (Integration)                │
│  - IP状态管理 - 证书管理 - 客户端管理 - 状态报告         │
│  - 请求封装 - 指纹管理 - IP管理 - 连接管理              │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   网络通信层 (Network)                    │
│  - conn (连接抽象) - quic (QUIC连接池) - netconnpool     │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   基础设施层 (Infrastructure)             │
│  - fingerprint (指纹模拟) - domaindns (DNS解析)        │
│  - localippool (本地IP池管理)                            │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   支撑层 (Support)                        │
│  - logs (日志记录与监控)                                 │
└─────────────────────────────────────────────────────────┘
```

## 📦 核心库（9个）

### 基础库（7个）

1. **[fingerprint](https://github.com/vistone/fingerprint)** - TLS指纹模拟库
   - 生成和管理浏览器TLS指纹
   - 支持44+主流浏览器指纹配置

2. **[logs](https://github.com/vistone/logs)** - 日志管理库
   - 多级别日志记录
   - 线程安全的日志输出

3. **[quic](https://github.com/vistone/quic)** - QUIC连接池
   - QUIC协议连接管理
   - 多路复用支持

4. **[conn](https://github.com/vistone/conn)** - 增强连接库
   - 通用网络连接抽象
   - 多协议连接封装

5. **[netconnpool](https://github.com/vistone/netconnpool)** - TCP连接池
   - TCP连接复用
   - 连接池管理

6. **[domaindns](https://github.com/vistone/domaindns)** - DNS解析库
   - 自定义DNS解析
   - DNS缓存和污染规避

7. **[localippool](https://github.com/vistone/localippool)** - 本地IP池管理
   - 本地IP地址池管理
   - IP轮换策略

### 新增库（2个）

8. **[certs](https://github.com/vistone/certs)** - 证书申请和管理库
   - 自动申请和管理TLS证书
   - 支持Let's Encrypt和自签名证书
   - 证书自动续期

9. **[whitelist-blacklist-manager](https://github.com/vistone/whitelist-blacklist-manager)** - 黑白名单管理器
   - IP地址的黑白名单管理
   - 线程安全的并发管理
   - 状态查询和移动操作

## 🚀 快速开始

### 安装依赖

```bash
go mod init crawler-system
go get github.com/vistone/fingerprint
go get github.com/vistone/logs
go get github.com/vistone/quic
go get github.com/vistone/conn
go get github.com/vistone/netconnpool
go get github.com/vistone/domaindns
go get github.com/vistone/localippool
go get github.com/vistone/certs
go get github.com/vistone/whitelist-blacklist-manager
```

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/vistone/certs"
    "github.com/vistone/logs"
    // ... 其他导入
)

func main() {
    // 初始化系统
    system, err := NewCrawlerSystem()
    if err != nil {
        log.Fatal(err)
    }
    defer system.Close()
    
    // 执行爬取请求
    executor := NewRequestExecutor(system)
    response, err := executor.Execute(context.Background(), RequestOptions{
        URL:     "https://example.com",
        Method:  "GET",
        UseQUIC: false,
    })
    
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    
    fmt.Printf("响应状态: %d\n", response.StatusCode)
}
```

## 📚 文档说明

### 📖 核心文档

**所有文档内容已整合到开发可行性报告中，建议直接阅读：**

- **[开发可行性报告](./docs/开发可行性报告.md)** ⭐ **推荐首先阅读**
  
  这是一份完整的项目文档，包含：
  - **执行摘要**：项目概述、核心价值、技术栈
  - **项目背景与需求分析**：业务背景、核心需求、目标用户
  - **系统架构设计**：整体架构、系统角色、核心设计原则
  - **核心库详细说明**：9个库的功能点、API设计、架构配合使用方法（每个库单独章节）
  - **核心组件设计**：IP状态管理器、IP池测试器、黑名单恢复监控器等
  - **系统工作流程**：启动流程、爬取请求处理流程、状态同步流程
  - **系统运行状态**：正常运行状态、待机状态、状态转换
  - **配置设计**：系统配置结构、关键配置项
  - **技术可行性分析**：技术栈可行性、架构设计可行性、性能可行性
  - **实施计划**：6个开发阶段、时间估算、资源需求
  - **风险评估与应对**：技术风险、业务风险、运维风险
  - **成功标准**：功能标准、性能标准、可靠性标准
  - **结论与建议**：可行性结论、关键成功因素、下一步行动

### 📝 其他文档（参考）

以下文档的内容已整合到可行性报告中，保留作为参考：

- [新包整合架构探讨](./docs/新包整合架构探讨.md) - 已整合到可行性报告
- [系统架构分析](./docs/系统架构分析.md) - 已整合到可行性报告
- [库配合使用详解](./docs/库配合使用详解.md) - 已整合到可行性报告第4章
- [多协议支持策略](./docs/多协议支持策略.md) - 已整合到可行性报告
- [IP库管理策略](./docs/IP库管理策略.md) - 已整合到可行性报告
- [API接口设计](./docs/API接口设计.md) - 参考文档
- [集成示例代码](./docs/集成示例代码.md) - 参考文档
- [最佳实践建议](./docs/最佳实践建议.md) - 已整合到可行性报告
- [系统设计总结](./docs/系统设计总结.md) - 已整合到可行性报告

## 🔄 核心工作流程

### 系统启动流程

1. 初始化日志系统
2. 初始化证书管理器（申请服务端证书）
3. 初始化IP状态管理器（创建空的黑白名单）
4. 初始化IP池测试器（解析目标域名，测试目标服务器IP）
5. 根据测试结果更新黑白名单
6. 检查白名单数量，决定系统运行状态
7. 初始化其他组件（DNS、指纹、连接池等）
8. 启动黑名单恢复监控
9. 启动QUIC服务端和状态报告服务
10. 启动爬虫服务

### 爬取请求处理流程

1. **DNS解析**：通过`domaindns`解析目标域名，得到目标服务器IP列表
2. **IP选择**：从目标服务器IP列表中选择，只使用白名单中的IP
3. **本地IP绑定**：从`localippool`获取本地出口IP，用于绑定本地连接
4. **指纹获取**：从`fingerprint`获取随机指纹
5. **协议选择**：选择协议（HTTP/3 → HTTP/2 → HTTP/1.1）
6. **连接获取**：从对应协议连接池获取连接
7. **发送请求**：使用目标服务器IP和本地出口IP发送请求
8. **响应处理**：
   - 200/其他 → 正常处理
   - 403 → 目标服务器IP移入黑名单，从白名单移除
9. **状态更新**：更新IP状态，归还连接

### 状态同步流程

- 服务器端实时收集系统状态（白名单、黑名单、客户端列表）
- 定期向所有连接的客户端广播状态更新
- 状态变化时立即推送更新
- 客户端接收状态并更新本地状态

## ⚙️ 配置说明

### 配置文件

系统使用 `config.toml` 作为主配置文件，包含9个核心模块的所有配置项。

**配置文件位置**：`config.toml`

**配置说明文档**：[CONFIG.md](./CONFIG.md)

### 快速配置

1. **复制配置文件**：
   ```bash
   cp config.toml config.toml.example
   ```

2. **编辑配置文件**：
   ```bash
   vim config.toml
   ```

3. **关键配置项**：

```toml
# 证书配置
[certificate]
server_domain = "crawler.example.com"  # VPS域名
provider = "letsencrypt"
auto_renewal = true

# IP池测试配置（目标域名列表）
[ip_pool_test]
target_domains = [
    "example.com",
    "target-site.com"
]
max_concurrent = 10
test_timeout = 10

# 本地IP池配置（本地出口IP）
[local_ip_pool]
ips = [
    "192.168.1.100",
    "192.168.1.101"
]

# 黑名单恢复配置
[blacklist_recovery]
enabled = true
check_interval = 1800  # 30分钟
ip_test_interval = 3600  # 1小时

# 状态报告配置
[status_report]
report_interval = 5
report_on_change = true
```

**重要说明**：
- **目标服务器IP**：通过`domaindns`解析目标域名得到，这些IP会被测试并加入黑白名单
- **本地IP池**：`localippool`管理的本地出口IP，用于绑定本地连接，与黑白名单无关

### 配置模块

配置文件包含15个配置模块：

1. **logs** - 日志配置
2. **fingerprint** - 指纹配置
3. **domaindns** - DNS解析配置
4. **local_ip_pool** - 本地IP池配置
5. **conn** - 连接配置
6. **netconnpool** - TCP连接池配置
7. **quic** - QUIC连接池配置
8. **certificate** - 证书配置
9. **ip_status** - 黑白名单配置
10. **ip_pool_test** - IP池测试配置
11. **blacklist_recovery** - 黑名单恢复配置
12. **status_report** - 状态报告配置
13. **server** - 服务端配置
14. **crawler** - 爬虫配置
15. **system** - 系统配置

详细配置说明请参考 [CONFIG.md](./CONFIG.md)

## 🎯 系统运行状态

### 正常运行状态（白名单有IP）
- 系统完全运行
- 参与爬取工作
- 提供QUIC服务端
- 运行监控和恢复机制

### 待机状态（白名单为空）
- 系统继续运行
- **不参与爬取工作**（因为没有可用IP）
- 继续提供QUIC服务端
- 继续运行监控和恢复机制
- 等待IP恢复或管理员更新配置

## 📊 核心设计原则

1. **白名单生命线原则**
   - 白名单是爬虫工作的生命线
   - 白名单为空时，系统不参与爬取工作，但继续运行

2. **403错误自动处理**
   - 检测到403响应 → 目标服务器IP移入黑名单 → 从白名单移除

3. **内存存储原则**
   - 黑白名单不持久化，只存在内存中
   - 系统重启后需要重新测试目标服务器IP

4. **实时状态同步**
   - 服务器端实时向所有客户端报告状态
   - 客户端实时了解服务器情况

## 🔍 示例代码

- [多协议爬虫示例](./examples/multi_protocol_crawler.go) - 完整的多协议爬虫实现

## 🛠️ 开发计划

详细开发计划请参考 [开发可行性报告 - 实施计划](./docs/开发可行性报告.md#9-实施计划)

### 开发阶段

1. **第一阶段**：IP池测试和黑白名单管理（2-3周）
2. **第二阶段**：黑名单恢复机制（1-2周）
3. **第三阶段**：证书管理（1-2周）
4. **第四阶段**：状态同步机制（2-3周）
5. **第五阶段**：QUIC服务端（3-4周）
6. **第六阶段**：完善和优化（2-3周）

**总开发时间**：11-17周（约3-4个月）

## ⚠️ 注意事项

1. **IP池理解**：
   - 黑白名单管理的是**目标服务器的IP地址**（通过`domaindns`解析域名得到）
   - **不是**本地IP池（`localippool`）的IP
   - 本地IP池用于绑定本地出口IP，与黑白名单无关

2. **系统状态**：
   - 白名单为空时，系统不参与爬取工作，但继续运行
   - 系统可以继续提供QUIC服务端和运行监控恢复机制

3. **证书管理**：
   - 证书域名从配置文件读取
   - 支持Let's Encrypt自动申请和续期

## 📝 许可证

本项目基于各库的许可证。

## 🤝 贡献

欢迎提交Issue和Pull Request。

## 📞 相关链接

### 核心库

- [fingerprint](https://github.com/vistone/fingerprint) - TLS指纹模拟
- [logs](https://github.com/vistone/logs) - 日志管理
- [quic](https://github.com/vistone/quic) - QUIC连接池
- [conn](https://github.com/vistone/conn) - 增强连接库
- [netconnpool](https://github.com/vistone/netconnpool) - TCP连接池
- [domaindns](https://github.com/vistone/domaindns) - DNS解析
- [localippool](https://github.com/vistone/localippool) - 本地IP池管理
- [certs](https://github.com/vistone/certs) - 证书申请和管理
- [whitelist-blacklist-manager](https://github.com/vistone/whitelist-blacklist-manager) - 黑白名单管理

---

**最后更新**：2025-01-XX

**推荐阅读顺序**：
1. [开发可行性报告](./docs/开发可行性报告.md) - **完整阅读，了解项目全貌**
   - 第1-3章：项目概述、背景、架构设计
   - 第4章：**核心库详细说明**（9个库的功能点和架构配合使用方法）
   - 第5-7章：核心组件、工作流程、运行状态
   - 第8-13章：配置、可行性分析、实施计划、风险评估、结论
2. [API接口设计](./docs/API接口设计.md) - 了解接口使用（参考）
3. [集成示例代码](./docs/集成示例代码.md) - 查看代码示例（参考）
