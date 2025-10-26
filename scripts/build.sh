#!/bin/bash

set -e

echo "Building application..."
go build -o bin/backend-journaling main.go

echo "Build complete! Binary: ./bin/backend-journaling"
