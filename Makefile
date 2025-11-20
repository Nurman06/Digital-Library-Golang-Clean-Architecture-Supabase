.PHONY: help build run test clean migrate-up migrate-down docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go
	@echo "Build complete! Binary available at bin/api"

run: ## Run the application
	@echo "Starting application..."
	@go run cmd/api/main.go

test: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -cover -short ./...

test-all: ## Run all tests (unit + integration)
	@echo "Running all tests..."
	@go test -v -cover -tags=integration ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v -tags=integration ./internal/repository/...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration-coverage: ## Run integration tests with coverage
	@echo "Running integration tests with coverage..."
	@go test -v -tags=integration -coverprofile=coverage-integration.out ./internal/repository/...
	@go tool cover -html=coverage-integration.out -o coverage-integration.html
	@echo "Integration coverage report generated: coverage-integration.html"

test-db-up: ## Start test database
	@echo "Starting test database..."
	@docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@docker-compose -f docker-compose.test.yml exec -T test-db pg_isready -U postgres || true

test-db-down: ## Stop test database
	@echo "Stopping test database..."
	@docker-compose -f docker-compose.test.yml down

test-db-clean: ## Stop test database and remove volumes
	@echo "Cleaning test database..."
	@docker-compose -f docker-compose.test.yml down -v

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated!"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete!"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run
	@echo "Lint complete!"

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@psql $(DATABASE_URL) -f migrations/001_create_tables.sql
	@echo "Migrations complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t digital-library-api:latest .
	@echo "Docker image built!"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env digital-library-api:latest

dev: ## Run in development mode with hot reload (requires air)
	@echo "Starting development server..."
	@air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed!"