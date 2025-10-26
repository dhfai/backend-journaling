#!/bin/bash

set -e

echo "Running tests..."
go test -v -race -coverprofile=coverage.out ./...

echo "Test coverage:"
go tool cover -func=coverage.out

echo "Done!"
