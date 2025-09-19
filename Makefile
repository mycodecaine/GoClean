# Build variables
BINARY_NAME=goclean
DOCKER_TAG=goclean:latest
PROTO_PATH=api/proto
PROTO_FILES=$(wildcard $(PROTO_PATH)/*.proto)

# Go variables
GOBASE=$(shell pwd)
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOBASE)/bin
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Docker variables
DOCKER_COMPOSE=docker-compose

.PHONY: help build clean test run dev proto deps docker

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	@mkdir -p api/proto/v1
	@protoc --go_out=api/proto/v1 --go_opt=paths=source_relative \
		--go-grpc_out=api/proto/v1 --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

build-http: deps ## Build HTTP server binary
	@echo "Building HTTP server..."
	@mkdir -p $(GOBIN)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-w -s" -o $(GOBIN)/$(BINARY_NAME)-http cmd/http/main.go

build-grpc: deps ## Build gRPC server binary
	@echo "Building gRPC server..."
	@mkdir -p $(GOBIN)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-w -s" -o $(GOBIN)/$(BINARY_NAME)-grpc cmd/grpc/main.go

build: build-http build-grpc ## Build all binaries

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(GOBIN)
	@go clean

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run-http: ## Run HTTP server in development mode
	@echo "Starting HTTP server..."
	@go run cmd/http/main.go

run-grpc: ## Run gRPC server in development mode
	@echo "Starting gRPC server..."
	@go run cmd/grpc/main.go

dev: ## Start development environment with hot reload
	@echo "Starting development environment..."
	@$(DOCKER_COMPOSE) up -d postgres redis keycloak
	@echo "Waiting for services to be ready..."
	@sleep 10
	@go run cmd/http/main.go

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_TAG) .

docker-up: ## Start all services with Docker Compose
	@echo "Starting all services..."
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all services
	@echo "Stopping all services..."
	@$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker Compose logs
	@$(DOCKER_COMPOSE) logs -f

docker-clean: ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) down -v --remove-orphans
	@docker system prune -f

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/http/main.go -o api/swagger

migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go up

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@go run cmd/migrate/main.go down

migrate-create: ## Create new migration file
	@read -p "Enter migration name: " name; \
	go run cmd/migrate/main.go create $$name

setup: deps docker-up ## Setup development environment
	@echo "Setting up development environment..."
	@sleep 15
	@echo "Development environment ready!"
	@echo "Services:"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379" 
	@echo "  - Keycloak: http://localhost:8081"
	@echo ""
	@echo "Run 'make run-http' to start the API server"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

security: ## Run security scan
	@echo "Running security scan..."
	@gosec ./...

mod-update: ## Update Go modules
	@echo "Updating Go modules..."
	@go get -u ./...
	@go mod tidy

all: clean deps proto build test ## Build everything from scratch

release: clean deps proto build test docker-build ## Prepare release build

.DEFAULT_GOAL := help