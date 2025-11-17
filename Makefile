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

test: ## Run tests
	@echo "Running tests..."
	@go test -v -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

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