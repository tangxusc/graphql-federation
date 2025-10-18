# Higress GraphQL Federation 插件

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](../LICENSE) 
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/wundergraph/graphql-go-tools)

一个基于 Go Filter 的 Higress GraphQL Federation 插件。

## 概述

本项目为 Higress 实现了一个 GraphQL Federation 插件，Higress 是一个基于 Envoy 的云原生 API 网关。该插件使用 Go 语言编写，并编译为共享库 (.so) 在 Envoy 代理中使用 Go Filter 扩展执行。

## 功能特性

- GraphQL Federation 查询解析与执行
- 与 Higress 网关集成
- 插件化架构，易于扩展
- 支持 GraphQL 请求的路由、限流和认证
- 基于 Go 的实现，具有高性能和安全性
- 支持通过 WebSocket 进行 GraphQL 订阅

## 项目结构

```
graphql-plugin/
├── cmd/
│   └── graphql/              # 插件主入口
├── pkg/
│   └── filter/               # Filter 实现
│       ├── config.go         # 配置管理
│       ├── filter.go         # 主要过滤逻辑
│       ├── engine.go         # GraphQL 执行引擎
│       └── types.go          # 类型定义
├── scripts/                  # 测试和部署脚本
│   ├── envoy.yaml           # Envoy 配置示例
│   └── docker-compose.yaml  # 测试环境设置
└── README.md                 # 项目文档
```

## 环境要求

- Go 1.21+
- Make
- Docker (用于测试)

## 构建

构建 Go Filter 共享库：

```bash
make build
```

这将生成一个 `build/graphql-federation_arm64.so` 文件，可以加载到 Envoy 中。

## 测试

本地测试插件：

```bash
# 使用 Docker Compose 启动服务
make test-up

# 发送 GraphQL 请求
curl -X POST -H "Content-Type: application/json" \
  --data '{"query":"{ users { id name } products { id name price } }"}' \
  http://localhost:10000/graphql

# 查看日志
docker logs -f scripts-envoy-1

# 停止服务
make test-down
```

## 配置

可以通过 Envoy 配置使用 Go Filter 扩展来配置插件。有关示例配置，请参见 `scripts/envoy.yaml`。

示例配置：

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
          - service_name: 'service1'
            graphql_url: 'http://service1:4001/graphql'
            subscription_url: 'ws://service1:4001/graphql'
          - service_name: 'service2'
            graphql_url: 'http://service2:4002/graphql'
            subscription_url: 'ws://service2:4002/graphql'
        schema_refresh_interval: "5m"
        schema_refresh_timeout: "1m"
```

### 配置字段

- `sub_graphql_config`: 子图配置数组
  - `service_name`: GraphQL 服务名称
  - `graphql_url`: GraphQL 查询的 HTTP 端点
  - `subscription_url`: GraphQL 订阅的 WebSocket 端点
- `schema_refresh_interval`: 刷新联邦模式的频率（默认：5m）
- `schema_refresh_timeout`: 模式刷新操作的超时时间（默认：1m）

## 许可证

本项目采用 Apache License 2.0 许可证。详见 [LICENSE](../LICENSE)。