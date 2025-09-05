# Project Structure

## Overview

The higress-graphql project follows a standard Go project layout with some WebAssembly-specific considerations.

## Directory Layout

```
.
├── cmd/                 # Main applications
│   └── graphql/         # GraphQL plugin main package
├── pkg/                 # Public library code
├── internal/            # Private application code
├── docs/                # Documentation
├── configs/             # Configuration files
├── Makefile             # Build instructions
├── docker-compose.yaml  # Docker Compose configuration for testing
├── envoy.yaml           # Envoy configuration
├── go.mod               # Go module definition
└── README.md            # Project README
```

## Key Files

- `cmd/graphql/main.go`: Entry point for the GraphQL plugin
- `Makefile`: Contains build instructions for compiling to WASM
- `docker-compose.yaml`: For local testing with Envoy and httpbin
- `envoy.yaml`: Envoy configuration with the GraphQL plugin loaded
- `go.mod`: Go module definition with dependencies

## Build Output

The build process produces a `graphql.wasm` file which is the compiled WebAssembly module that can be loaded into Envoy.