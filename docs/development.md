# Development Guide

## Prerequisites

- Go 1.22+
- Docker and Docker Compose (for testing)
- Make

## Project Setup

1. Clone the repository
2. Run `go mod tidy` to ensure all dependencies are downloaded

## Building the Plugin

To build the WebAssembly module:

```bash
make build
```

This will compile the Go code to a WASM module named `graphql.wasm`.

## Testing Locally

To test the plugin locally with Envoy:

```bash
docker-compose up
```

This will start:
- An Envoy instance with the GraphQL plugin loaded
- An httpbin service for testing

The Envoy instance will be available at http://localhost:10000

## Project Structure

The project follows the standard Go project layout:

- `cmd/graphql/` - Main application entry point
- `pkg/` - Public library code
- `internal/` - Private application code
- `configs/` - Configuration files
- `docs/` - Documentation

## Code Organization

### Main Entry Point

The main entry point is in `cmd/graphql/main.go`. This file initializes the plugin and sets up the configuration parsing and request handling functions.

### Configuration

The plugin configuration is defined in the `MyConfig` struct and parsed in the `parseConfig` function.

### Request Handling

Request handling is done through several functions:
- `onHttpRequestHeaders` - Processes HTTP request headers
- `body` - Processes request body

## Extending the Plugin

To add new functionality to the plugin:

1. Modify the `MyConfig` struct to add new configuration options
2. Update the `parseConfig` function to parse new configuration options
3. Add new handler functions for different phases of request processing
4. Register the new handlers in the `init` function