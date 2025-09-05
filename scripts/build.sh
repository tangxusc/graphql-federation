#!/bin/bash

# Build script for higress-graphql plugin

set -e

echo "Building higress-graphql WASM module..."

# Build the WASM module
make build

echo "Build completed successfully!"

# Check if the WASM file was created
if [ -f "graphql.wasm" ]; then
    echo "graphql.wasm file created successfully"
    ls -lh graphql.wasm
else
    echo "Error: graphql.wasm file not found"
    exit 1
fi

echo "Build process completed!"