# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=qris-pos-api
BINARY_UNIX=$(BINARY_NAME)_unix

# Development
.PHONY: dev
dev:
	@echo "Starting development server..."
	@$(GOCMD) run cmd/api/main.go

.PHONY: build
build:
	@echo "Building binary..."
	@$(GOBUILD) -o $(BINARY_NAME) -v cmd/api/main.go

.PHONY: build-linux
build-linux:
	@echo "Building binary for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v cmd/api/main.go

# Testing
.PHONY: test
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@$(GOTEST) -v ./tests/unit/...

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@$(GOTEST) -v ./tests/integration/...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Database
.PHONY: db-migrate-up
db-migrate-up:
	@echo "Running database migrations..."
	@$(GOCMD) run cmd/api/main.go migrate

.PHONY: db-migrate-down
db-migrate-down:
	@echo "Rolling back database migrations..."
	@echo "Not implemented yet"

.PHONY: db-seed
db-seed:
	@echo "Seeding database..."
	@echo "Not implemented yet"

.PHONY: db-reset
db-reset:
	@echo "Resetting database..."
	@echo "Not implemented yet"

# Code Quality
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GOCMD) fmt ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	@$(GOCMD) vet ./...

.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying go modules..."
	@$(GOMOD) tidy

# Docker
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t qris-pos-backend .

.PHONY: docker-run
docker-run:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

# Utilities
.PHONY: gen-docs
gen-docs:
	@echo "Generating API documentation..."
	@echo "Not implemented yet"

.PHONY: gen-mocks
gen-mocks:
	@echo "Generating test mocks..."
	@echo "Not implemented yet"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_UNIX)
	@rm -f coverage.out coverage.html

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  dev              - Start development server"
	@echo "  build            - Build production binary"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Start with Docker Compose"
	@echo "  clean            - Clean build artifacts"