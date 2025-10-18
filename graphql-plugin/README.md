# Higress GraphQL Federation Plugin

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](../LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/wundergraph/graphql-go-tools)

A GraphQL Federation plugin for Higress based on Go filter.

## Overview

This project implements a GraphQL Federation plugin for Higress, which is a cloud-native API gateway based on Envoy. The
plugin is written in Go and compiled as a shared library (.so) for execution within the Envoy proxy using the Go filter extension.

## Features

- GraphQL Federation query parsing and execution
- Integration with Higress gateway
- Plugin-based architecture for easy extension
- Support for routing, rate limiting, and authentication of GraphQL requests
- Go-based implementation for high performance and security
- Support for GraphQL subscriptions via WebSocket

## Project Structure

```
graphql-plugin/
├── cmd/
│   └── graphql/              # Plugin main entry point
├── pkg/
│   └── filter/               # Filter implementation
│       ├── config.go         # Configuration management
│       ├── filter.go         # Main filter logic
│       ├── engine.go         # GraphQL execution engine
│       └── types.go          # Type definitions
├── scripts/                  # Testing and deployment scripts
│   ├── envoy.yaml           # Envoy configuration example
│   └── docker-compose.yaml  # Test environment setup
└── README.md                 # Project documentation
```

## Prerequisites

- Go 1.21+
- Make
- Docker (for testing)

## Building

To build the Go filter shared library:

```bash
make build
```

This will produce a `build/graphql-federation_arm64.so` file that can be loaded into Envoy.

## Testing

To test the plugin locally:

```bash
# Start services using Docker Compose
make test-up

# Send a GraphQL request
curl -X POST -H "Content-Type: application/json" \
  --data '{"query":"{ users { id name } products { id name price } }"}' \
  http://localhost:10000/graphql

# Check logs
docker logs -f scripts-envoy-1

# Stop services
make test-down
```

## Configuration

The plugin can be configured through the Envoy configuration using the Go filter extension. See `scripts/envoy.yaml` for an example configuration.

Example configuration:

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

### Configuration Fields

- `sub_graphql_config`: Array of subgraph configurations
  - `service_name`: Name of the GraphQL service
  - `graphql_url`: HTTP endpoint for GraphQL queries
  - `subscription_url`: WebSocket endpoint for GraphQL subscriptions
- `schema_refresh_interval`: How often to refresh the federated schema (default: 5m)
- `schema_refresh_timeout`: Timeout for schema refresh operations (default: 1m)

## License

This project is licensed under the Apache License 2.0. See [LICENSE](../LICENSE) for details.