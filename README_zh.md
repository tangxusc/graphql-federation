# Higress GraphQL Federation 插件

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

一个基于 WebAssembly 的 Higress GraphQL Federation 插件。

## 概述

本项目为 Higress 实现了一个 GraphQL Federation 插件，Higress 是一个基于 Envoy 的云原生 API 网关。该插件使用 Go 语言编写，并编译为 WebAssembly (WASM) 在 Envoy 代理中执行。

## 功能特性

- GraphQL Federation 查询解析与执行
- 与 Higress 网关集成
- 插件化架构，易于扩展
- 支持 GraphQL 请求的路由、限流和认证
- 基于 WASM 的实现，具有高性能和安全性

## 项目结构

```
.
├── cmd/
│   └── graphql/              # 插件主入口
├── graphql-go-tools-execution/ # GraphQL 执行引擎
├── graphql-go-tools-v2/      # GraphQL 工具 v2
├── pkg/
│   ├── config/               # 配置管理
│   └── process/              # 请求处理逻辑
├── scripts/                  # 测试和部署脚本
├── Makefile                  # 构建指令
└── README.md                 # 项目文档
```

## 环境要求

- Go 1.24+
- Make
- Docker (用于测试)

## 构建

构建 WASM 模块：

```bash
make build
```

这将生成一个 `build/graphql-federation.wasm` 文件，可以加载到 Envoy 中。

## 测试

本地测试插件：

```bash
# 使用 Docker Compose 启动服务
docker-compose -f scripts/docker-compose.yaml up -d

# 检查服务状态
docker-compose -f scripts/docker-compose.yaml ps

# 发送 GraphQL 请求
curl -X POST http://localhost:10000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ users { id name } products { id name price } }"
  }'

# 查看日志
docker logs scripts-envoy-1

# 停止服务
docker-compose -f scripts/docker-compose.yaml down
```

## 配置

可以通过 Envoy 配置来配置插件。有关示例配置，请参见 `scripts/envoy.yaml`。

示例配置：
```yaml
- name: envoy.filters.http.wasm
  config:
    config:
      name: "higress-graphql-federation"
      root_id: "higress-graphql-federation"
      vm_config:
        runtime: "envoy.wasm.runtime.v8"
        code:
          local:
            filename: "/etc/envoy/graphql-federation.wasm"
      configuration:
        "@type": "type.googleapis.com/google.protobuf.StringValue"
        value: |
          {
            "enableQueryPlanning": true,
            "enableCaching": true,
            "maxQueryDepth": 10,
            "queryTimeout": 5000,
            "enableIntrospection": true,
            "graphqlAddress": "http://graphql-service"
          }
```

## 许可证

本项目采用 Apache License 2.0 许可证。详见 [LICENSE](LICENSE)。