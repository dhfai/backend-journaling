.PHONY: help build run test clean keys migrate dev

help:
	@echo "Available commands:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make keys      - Generate JWT keys"
	@echo "  make dev       - Run in development mode with hot reload"

build:
	@echo "Building..."
	@go build -o bin/backend-journaling main.go

run:
	@echo "Running application..."
	@go run main.go

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out

keys:
	@echo "Generating JWT keys..."
	@mkdir -p keys
	@openssl genrsa -out keys/jwt_private.pem 2048
	@openssl rsa -in keys/jwt_private.pem -pubout -out keys/jwt_public.pem
	@echo "Keys generated successfully!"

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

dev:
	@echo "Running in development mode..."
	@go run main.go
