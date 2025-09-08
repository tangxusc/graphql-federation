# Higress GraphQL Federation Plugin

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/wundergraph/graphql-go-tools)


A GraphQL Federation plugin for Higress based on WebAssembly.

## Overview

This project implements a GraphQL Federation plugin for Higress, which is a cloud-native API gateway based on Envoy. The plugin is written in Go and compiled to WebAssembly (WASM) for execution within the Envoy proxy.

## Features

- GraphQL Federation query parsing and execution
- Integration with Higress gateway
- Plugin-based architecture for easy extension
- Support for routing, rate limiting, and authentication of GraphQL requests
- WASM-based implementation for high performance and security

## Project Structure

```
.
├── cmd/
│   └── graphql/              # Plugin main entry point
├── graphql-go-tools-execution/ # GraphQL execution engine
├── graphql-go-tools-v2/      # GraphQL tools v2
├── pkg/
│   ├── config/               # Configuration management
│   └── process/              # Request processing logic
├── scripts/                  # Testing and deployment scripts
├── Makefile                  # Build instructions
└── README.md                 # Project documentation
```

## Prerequisites

- Go 1.24+
- Make
- Docker (for testing)

## Building

To build the WASM module:

```bash
make build
```

This will produce a `build/graphql-federation.wasm` file that can be loaded into Envoy.

## Testing

To test the plugin locally:

```bash
# Start services using Docker Compose
docker-compose -f scripts/docker-compose.yaml up -d

# Check service status
docker-compose -f scripts/docker-compose.yaml ps

# Send a GraphQL request
curl -X POST http://localhost:10000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ users { id name } products { id name price } }"
  }'

# Check logs
docker logs scripts-envoy-1

# Stop services
docker-compose -f scripts/docker-compose.yaml down
```

## Configuration

The plugin can be configured through the Envoy configuration. See `scripts/envoy.yaml` for an example configuration.

Example configuration:
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

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.