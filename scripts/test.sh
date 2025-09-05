#!/bin/bash

# Test script for higress-graphql plugin

set -e

echo "Testing higress-graphql plugin..."

# Run go vet to check for potential issues
echo "Running go vet..."
go vet ./...

# Run go fmt to check formatting
echo "Checking code formatting..."
go fmt ./...

# Check if there are any uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
  echo "Warning: There are uncommitted changes in the repository"
else
  echo "No uncommitted changes found"
fi

echo "Basic tests completed!"