# higress-graphql

A GraphQL plugin for Higress based on WebAssembly.

## Overview

This project implements a GraphQL plugin for Higress, which is a cloud-native API gateway based on Envoy. The plugin is written in Go and compiled to WebAssembly (WASM) for execution within the Envoy proxy.

## Features

- GraphQL query parsing and execution
- Integration with Higress gateway
- Plugin-based architecture for easy extension
- Support for routing, rate limiting, and authentication of GraphQL requests

## Project Structure

```
.
├── cmd/                 # Main applications
│   └── graphql/         # GraphQL plugin main package
├── pkg/                 # Public library code
│   ├── config/          # Configuration types
│   └── graphql/         # GraphQL types
├── internal/            # Private application code
│   ├── handler/         # Request handlers
│   └── middleware/      # Middleware functions
├── configs/             # Configuration files
│   ├── config.yaml      # Example configuration
│   └── schema.graphql   # Example GraphQL schema
├── docs/                # Documentation
├── Makefile             # Build instructions
├── docker-compose.yaml  # Docker Compose configuration for testing
├── envoy.yaml           # Envoy configuration
└── go.mod               # Go module definition
```

## Building

To build the WASM module:

```bash
make build
```

This will produce a `graphql-federation.wasm` file that can be loaded into Envoy.

## Testing

To test the plugin locally:

```bash
docker-compose -f ./scripts/docker-compose.yaml up -d
# 检查服务状态
docker-compose -f ./scripts/docker-compose.yaml ps
```
send graphql request
```shell
curl -X POST http://localhost:10000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ users { id name } products { id name price } }"
  }'
  
docker logs scripts-envoy-1
```
close docker container
```shell
docker-compose -f ./scripts/docker-compose.yaml down 
```


This will start an Envoy instance with the GraphQL plugin loaded, along with an httpbin service for testing.

## Configuration

The plugin can be configured through the Envoy configuration. See `envoy.yaml` for an example configuration.