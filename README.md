# GraphQL Federation

<div align="center">

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org/)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/tangxusc/graphql-federation)
[![Go Report Card](https://goreportcard.com/badge/github.com/tangxusc/graphql-federation)](https://goreportcard.com/report/github.com/tangxusc/graphql-federation)
[![GitHub stars](https://img.shields.io/github/stars/tangxusc/graphql-federation.svg)](https://github.com/tangxusc/graphql-federation/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tangxusc/graphql-federation.svg)](https://github.com/tangxusc/graphql-federation/network)

**A high-performance GraphQL Federation Envoy Go filter for Istio proxyv2**

[English](#english) | [中文](#中文)

</div>

---

## English

### Overview

GraphQL Federation is a comprehensive solution that provides GraphQL Federation capabilities for Istio proxyv2 (Envoy). Built on top of the powerful [graphql-go-tools](https://github.com/wundergraph/graphql-go-tools) library from WunderGraph, this project enables seamless integration of multiple GraphQL services into a unified federated GraphQL API using Apollo Federation protocol.

### Key Features

- 🚀 **High Performance**: Built with Go for optimal performance and low latency
- 🔗 **GraphQL Federation**: Seamlessly combine multiple GraphQL services using Apollo Federation
- 🌐 **Istio proxyv2 Integration**: Runs as an Envoy Go filter within Istio proxyv2
- 🔌 **Plugin Architecture**: Extensible plugin-based design
- 📡 **Automatic Schema Refresh**: Periodic schema synchronization from subgraphs
- 📊 **Request Logging**: Comprehensive request/response logging
- 🐳 **Docker Ready**: Containerized deployment support

### Architecture

This project consists of several key components:

```
graphql-federation/
├── cmd/
│   └── graphql/
│       └── graphql-federation.go
├── pkg/
│   └── filter/
│       ├── config.go
│       ├── engine.go
│       ├── filter.go
│       ├── logger_adapter.go
│       ├── register_filter.go
│       └── types.go
├── scripts/
│   ├── docker-compose.yaml
│   └── envoy.yaml
├── Dockerfile
├── Dockerfile_local
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

### Quick Start

#### Prerequisites

- Go 1.25+
- Docker (for testing)
- Make

#### Building the Plugin

```bash
# Clone the repository
git clone https://github.com/tangxusc/graphql-federation.git
cd graphql-federation

# Build the GraphQL Federation plugin
make build-local
```

This will generate a shared library (`graphql-federation_arm64.so` or `graphql-federation_amd64.so` depending on `GOARCH`) in the `build/` directory that can be loaded into Envoy.

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
    library_path: "/var/lib/istio/envoy/graphql-federation_arm64.so"
    plugin_name: graphql-federation
    plugin_config:
      "@type": type.googleapis.com/xds.type.v3.TypedStruct
      value:
        sub_graphql_config:
          - service_name: 'users'
            graphql_url: 'http://users-service:4001/graphql'
          - service_name: 'products'
            graphql_url: 'http://products-service:4002/graphql'
        schema_refresh_interval: "5m"
        schema_refresh_timeout: "1m"
```

### Configuration Options

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `sub_graphql_config` | Array | Configuration for subgraph services | Required |
| `service_name` | String | Name of the GraphQL service | Required |
| `graphql_url` | String | HTTP endpoint for GraphQL queries (also used for subscriptions) | Required |
| `schema_refresh_interval` | Duration | How often to refresh the federated schema | `5m` |
| `schema_refresh_timeout` | Duration | Timeout for schema refresh operations | `1m` |

**Note**: The filter only handles HTTP POST requests to `/graphql` endpoint. GraphQL subscriptions are currently not supported.

### Development

#### Project Structure

- **cmd/graphql/**: Entry point for building the Envoy Go filter plugin
- **pkg/filter/**: Core implementation including filter logic, configuration, and engine management
  - `filter.go`: Main filter implementation handling HTTP requests
  - `config.go`: Configuration parsing and plugin factory
  - `engine.go`: GraphQL federation engine initialization and schema refresh
  - `register_filter.go`: Filter registration with Envoy
  - `logger_adapter.go`: Logger adapter implementation
  - `types.go`: Type definitions

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

GraphQL Federation 是一个基于 Istio proxyv2（Envoy）的 GraphQL 联邦能力解决方案。依托强大的 [graphql-go-tools](https://github.com/wundergraph/graphql-go-tools) 库（由 WunderGraph 开发）构建，该项目能够将多个 GraphQL 服务无缝集成到统一的联邦 GraphQL API 中，遵循 Apollo Federation 协议。

### 核心特性

- 🚀 **高性能**: 使用 Go 构建，具有最佳性能和低延迟
- 🔗 **GraphQL 联邦**: 使用 Apollo Federation 无缝组合多个 GraphQL 服务
- 🌐 **Istio proxyv2 集成**: 作为 Envoy Go 过滤器运行于 Istio proxyv2 中
- 🔌 **插件架构**: 可扩展的基于插件的设计
- 📡 **自动 Schema 刷新**: 定期从子图同步 Schema
- 📊 **请求日志**: 全面的请求/响应日志记录
- 🐳 **Docker 就绪**: 容器化部署支持

### 架构

该项目由几个关键组件组成：

```
graphql-federation/
├── cmd/
│   └── graphql/
│       └── graphql-federation.go
├── pkg/
│   └── filter/
│       ├── config.go
│       ├── engine.go
│       ├── filter.go
│       ├── logger_adapter.go
│       ├── register_filter.go
│       └── types.go
├── scripts/
│   ├── docker-compose.yaml
│   └── envoy.yaml
├── Dockerfile
├── Dockerfile_local
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

### 快速开始

#### 环境要求

- Go 1.25+
- Docker (用于测试)
- Make

#### 构建插件

```bash
# 克隆仓库
git clone https://github.com/tangxusc/graphql-federation.git
cd graphql-federation

# 构建 GraphQL Federation 插件
make build-local
```

这将根据 `GOARCH` 环境变量（默认为 `arm64`）生成对应的共享库（`graphql-federation_arm64.so` 或 `graphql-federation_amd64.so`），文件位于 `build/` 目录中，可以加载到 Envoy 中。

#### 本地测试

```bash
# 启动测试环境
make test-up

# 发送 GraphQL 联邦查询
curl -X POST -H "Content-Type: application/json" \
  --data '{"query":"{ users { id name } products { id name price } }"}' \
  http://localhost:10000/graphql

# 查看日志
docker logs -f scripts_envoy_1

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
    library_path: "/var/lib/istio/envoy/graphql-federation_arm64.so"
    plugin_name: graphql-federation
    plugin_config:
      "@type": type.googleapis.com/xds.type.v3.TypedStruct
      value:
        sub_graphql_config:
          - service_name: 'users'
            graphql_url: 'http://users-service:4001/graphql'
          - service_name: 'products'
            graphql_url: 'http://products-service:4002/graphql'
        schema_refresh_interval: "5m"
        schema_refresh_timeout: "1m"
```

### 配置选项

| 字段 | 类型 | 描述 | 默认值 |
|------|------|------|--------|
| `sub_graphql_config` | 数组 | 子图服务配置 | 必需 |
| `service_name` | 字符串 | GraphQL 服务名称 | 必需 |
| `graphql_url` | 字符串 | GraphQL 查询的 HTTP 端点（订阅也使用此端点） | 必需 |
| `schema_refresh_interval` | 持续时间 | 刷新联邦模式的频率 | `5m` |
| `schema_refresh_timeout` | 持续时间 | 模式刷新操作的超时时间 | `1m` |

**注意**: 过滤器目前仅处理发送到 `/graphql` 端点的 HTTP POST 请求。GraphQL 订阅功能目前尚未支持。

### 开发

#### 项目结构

- **cmd/graphql/**: 构建 Envoy Go 过滤器插件的入口点
- **pkg/filter/**: 核心实现，包括过滤器逻辑、配置和引擎管理
  - `filter.go`: 处理 HTTP 请求的主要过滤器实现
  - `config.go`: 配置解析和插件工厂
  - `engine.go`: GraphQL 联邦引擎初始化和 Schema 刷新
  - `register_filter.go`: 在 Envoy 中注册过滤器
  - `logger_adapter.go`: 日志适配器实现
  - `types.go`: 类型定义

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

[GitHub](https://github.com/tangxusc/graphql-federation) • [Issues](https://github.com/tangxusc/graphql-federation/issues) • [Discussions](https://github.com/tangxusc/graphql-federation/discussions)

</div>
