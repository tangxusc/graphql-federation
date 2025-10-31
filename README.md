# GraphQL Federation

<div align="center">

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/tangxusc/higress-graphql-federation)
[![Go Report Card](https://goreportcard.com/badge/github.com/tangxusc/higress-graphql-federation)](https://goreportcard.com/report/github.com/tangxusc/higress-graphql-federation)
[![GitHub stars](https://img.shields.io/github/stars/tangxusc/higress-graphql-federation.svg)](https://github.com/tangxusc/higress-graphql-federation/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tangxusc/higress-graphql-federation.svg)](https://github.com/tangxusc/higress-graphql-federation/network)

**A high-performance GraphQL Federation Envoy Go filter for Istio proxyv2**

[English](#english) | [中文](#中文)

</div>

---

## English

### Overview

GraphQL Federation is a comprehensive solution that provides GraphQL Federation capabilities for Istio proxyv2 (Envoy). Built on top of the powerful `graphql-go-tools` library, this project enables seamless integration of multiple GraphQL services into a unified federated GraphQL API.

### Key Features

- 🚀 **High Performance**: Built with Go for optimal performance and low latency
- 🔗 **GraphQL Federation**: Seamlessly combine multiple GraphQL services
- 🌐 **Istio proxyv2 Integration**: Runs as an Envoy Go filter within Istio proxyv2
- 🔌 **Plugin Architecture**: Extensible plugin-based design
- 📡 **WebSocket Support**: Full support for GraphQL subscriptions
- 🛡️ **Security**: Built-in authentication and authorization support
- 📊 **Monitoring**: Comprehensive logging and metrics
- 🐳 **Docker Ready**: Containerized deployment support

### Architecture

This project consists of several key components:

```
graphql-federation/
├── graphql-plugin/              # Envoy Golang HTTP filter implementation
│   ├── cmd/graphql/             # Plugin entry point
│   ├── pkg/filter/              # Core filter logic
│   └── scripts/                 # Deployment and testing scripts
├── graphql-go-tools-v2/         # Core GraphQL tools library
├── graphql-go-tools-execution/  # GraphQL execution engine
├── composition-go/              # Federation composition utilities
└── Dockerfile                   # Container build configuration
```

### Quick Start

#### Prerequisites

- Go 1.22+
- Docker (for testing)
- Make

#### Building the Plugin

```bash
# Clone the repository
git clone https://github.com/tangxusc/higress-graphql-federation.git
cd higress-graphql-federation

# Build the GraphQL Federation plugin
make build
```

This will generate a shared library (`graphql-federation_arm64.so`) that can be loaded into Envoy.

#### Testing Locally

```bash
# Start the test environment
make test-up

# Send a GraphQL federation query
curl -X POST -H "Content-Type: application/json" \
  --data '{"query":"{ users { id name } products { id name price } }"}' \
  http://localhost:10000/graphql

# View logs
docker logs -f scripts-envoy-1

# Stop the test environment
make test-down
```

### Configuration

Configure the plugin through Envoy's Go filter extension:

```yaml
- name: envoy.filters.http.golang
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.filters.http.golang.v3alpha.Config
    library_id: graphql-federation
    library_path: "/etc/envoy/graphql-federation_arm64.so"
    plugin_name: graphql-federation
    plugin_config:
      "@type": type.googleapis.com/xds.type.v3.TypedStruct
      value:
        sub_graphql_config:
          - service_name: 'users'
            graphql_url: 'http://users-service:4001/graphql'
            subscription_url: 'ws://users-service:4001/graphql'
          - service_name: 'products'
            graphql_url: 'http://products-service:4002/graphql'
            subscription_url: 'ws://products-service:4002/graphql'
        schema_refresh_interval: "5m"
        schema_refresh_timeout: "1m"
```

### Configuration Options

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `sub_graphql_config` | Array | Configuration for subgraph services | Required |
| `service_name` | String | Name of the GraphQL service | Required |
| `graphql_url` | String | HTTP endpoint for GraphQL queries | Required |
| `subscription_url` | String | WebSocket endpoint for subscriptions | Optional |
| `schema_refresh_interval` | Duration | How often to refresh the federated schema | `5m` |
| `schema_refresh_timeout` | Duration | Timeout for schema refresh operations | `1m` |

### Development

#### Project Structure

- **graphql-plugin/**: Envoy Golang HTTP filter implementation for Istio proxyv2
- **graphql-go-tools-v2/**: Core GraphQL parsing, validation, and execution library
- **graphql-go-tools-execution/**: GraphQL execution engine with federation support
- **composition-go/**: Federation composition and router configuration utilities

#### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## 中文

### 概述

GraphQL Federation 是一个基于 Istio proxyv2（Envoy）的 GraphQL 联邦能力解决方案。依托强大的 `graphql-go-tools` 库构建，该项目能够将多个 GraphQL 服务无缝集成到统一的联邦 GraphQL API 中。

### 核心特性

- 🚀 **高性能**: 使用 Go 构建，具有最佳性能和低延迟
- 🔗 **GraphQL 联邦**: 无缝组合多个 GraphQL 服务
- 🌐 **Istio proxyv2 集成**: 作为 Envoy Go 过滤器运行于 Istio proxyv2 中
- 🔌 **插件架构**: 可扩展的基于插件的设计
- 📡 **WebSocket 支持**: 完整支持 GraphQL 订阅
- 🛡️ **安全性**: 内置认证和授权支持
- 📊 **监控**: 全面的日志记录和指标
- 🐳 **Docker 就绪**: 容器化部署支持

### 架构

该项目由几个关键组件组成：

```
graphql-federation/
├── graphql-plugin/              # Envoy Golang HTTP 过滤器实现
│   ├── cmd/graphql/             # 插件入口点
│   ├── pkg/filter/              # 核心过滤器逻辑
│   └── scripts/                 # 部署和测试脚本
├── graphql-go-tools-v2/         # 核心 GraphQL 工具库
├── graphql-go-tools-execution/  # GraphQL 执行引擎
├── composition-go/              # 联邦组合工具
└── Dockerfile                   # 容器构建配置
```

### 快速开始

#### 环境要求

- Go 1.22+
- Docker (用于测试)
- Make

#### 构建插件

```bash
# 克隆仓库
git clone https://github.com/tangxusc/higress-graphql-federation.git
cd higress-graphql-federation

# 构建 GraphQL Federation 插件
make build
```

这将生成一个共享库 (`graphql-federation_arm64.so`)，可以加载到 Envoy 中。

#### 本地测试

```bash
# 启动测试环境
make test-up

# 发送 GraphQL 联邦查询
curl -X POST -H "Content-Type: application/json" \
  --data '{"query":"{ users { id name } products { id name price } }"}' \
  http://localhost:10000/graphql

# 查看日志
docker logs -f scripts-envoy-1

# 停止测试环境
make test-down
```

### 配置

通过 Envoy 的 Go 过滤器扩展配置插件：

```yaml
- name: envoy.filters.http.golang
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.filters.http.golang.v3alpha.Config
    library_id: graphql-federation
    library_path: "/etc/envoy/graphql-federation_arm64.so"
    plugin_name: graphql-federation
    plugin_config:
      "@type": type.googleapis.com/xds.type.v3.TypedStruct
      value:
        sub_graphql_config:
          - service_name: 'users'
            graphql_url: 'http://users-service:4001/graphql'
            subscription_url: 'ws://users-service:4001/graphql'
          - service_name: 'products'
            graphql_url: 'http://products-service:4002/graphql'
            subscription_url: 'ws://products-service:4002/graphql'
        schema_refresh_interval: "5m"
        schema_refresh_timeout: "1m"
```

### 配置选项

| 字段 | 类型 | 描述 | 默认值 |
|------|------|------|--------|
| `sub_graphql_config` | 数组 | 子图服务配置 | 必需 |
| `service_name` | 字符串 | GraphQL 服务名称 | 必需 |
| `graphql_url` | 字符串 | GraphQL 查询的 HTTP 端点 | 必需 |
| `subscription_url` | 字符串 | 订阅的 WebSocket 端点 | 可选 |
| `schema_refresh_interval` | 持续时间 | 刷新联邦模式的频率 | `5m` |
| `schema_refresh_timeout` | 持续时间 | 模式刷新操作的超时时间 | `1m` |

### 开发

#### 项目结构

- **graphql-plugin/**: 面向 Istio proxyv2 的 Envoy Golang HTTP 过滤器实现
- **graphql-go-tools-v2/**: 核心 GraphQL 解析、验证和执行库
- **graphql-go-tools-execution/**: 支持联邦的 GraphQL 执行引擎
- **composition-go/**: 联邦组合和路由器配置工具

#### 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

### 许可证

本项目采用 Apache License 2.0 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

<div align="center">

**Made with ❤️ for the GraphQL community**

[GitHub](https://github.com/tangxusc/higress-graphql-federation) • [Issues](https://github.com/tangxusc/higress-graphql-federation/issues) • [Discussions](https://github.com/tangxusc/higress-graphql-federation/discussions)

</div>
