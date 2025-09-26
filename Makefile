# Makefile for Hexagonal Architecture Go Template
.PHONY: help install clean test test-unit test-integration test-coverage build run docker setup-dev lint fmt generate deps

# Variables
APP_NAME := hexago
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Detect OS for platform-specific commands
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RM := del /Q
    RMDIR := rmdir /S /Q
    EXE_EXT := .exe
else
    DETECTED_OS := $(shell uname -s)
    RM := rm -f
    RMDIR := rm -rf
    EXE_EXT :=
endif

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

setup-dev: install ## Setup development environment
	@echo "Setting up development environment..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

deps: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

generate: ## Generate code (protobuf, mocks, etc.)
	@echo "Generating code..."
	@protoc --go_out=. --go-grpc_out=. api/backoffice/grpc/ohlc.proto

clean: ## Clean build artifacts and temporary files
	@echo "Cleaning..."
ifeq ($(DETECTED_OS),Windows)
	@if exist $(COVERAGE_FILE) $(RM) $(COVERAGE_FILE)
	@if exist $(COVERAGE_HTML) $(RM) $(COVERAGE_HTML)  
	@if exist *.exe $(RM) *.exe
else
	@$(RM) $(COVERAGE_FILE) $(COVERAGE_HTML)
	@$(RM) *.exe
endif
	@go clean -cache -testcache

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

test: ## Run all tests
	@echo "Running all tests..."
	@go test ./tests/...

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	@go test ./tests/unit/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test ./tests/integration/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) ./tests/...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -race ./tests/...

build: clean generate ## Build all applications
	@echo "Building applications..."
	@go build -o all-in-one$(EXE_EXT) ./cmd/all-in-one
	@go build -o binance-tick-consumer$(EXE_EXT) ./cmd/binance-tick-consumer
	@go build -o backoffice-api$(EXE_EXT) ./cmd/backoffice-api
	@go build -o ohlc-generator$(EXE_EXT) ./cmd/ohlc-generator

run-all: ## Run all-in-one application
	@echo "Starting all-in-one application..."
	@./all-in-one$(EXE_EXT)

run-consumer: ## Run binance tick consumer
	@echo "Starting binance tick consumer..."
	@./binance-tick-consumer$(EXE_EXT)

run-api: ## Run backoffice API
	@echo "Starting backoffice API..."
	@./backoffice-api$(EXE_EXT)

run-generator: ## Run OHLC generator
	@echo "Starting OHLC generator..."
	@./ohlc-generator$(EXE_EXT)

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	@docker-compose build

docker-up: ## Start services with Docker Compose
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-down: ## Stop services with Docker Compose
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

docker-logs: ## Show Docker logs
	@echo "Showing Docker logs..."
	@docker-compose logs -f

# Template targets for new projects (Go-based, cross-platform)
template-init: ## Initialize new project from template
	@echo "Initializing new project from template..."
	@go run ./scripts/template-init

template-create-adapter: ## Create new adapter from template
	@echo "Creating new adapter..."
	@go run ./scripts/create-adapter

template-create-service: ## Create new service from template
	@echo "Creating new service..."
	@go run ./scripts/create-service