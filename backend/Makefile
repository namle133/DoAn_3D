.PHONY: help build run test clean docker-up docker-down dev deps

# Variables
BINARY_NAME=vinhomes
DOCKER_COMPOSE_FILE=docker-compose.yml

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download and verify Go dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	go mod tidy
	@echo "✓ Dependencies ready"

build: ## Build the application binary
	@echo "Building application..."
	go build -o $(BINARY_NAME) ./cmd
	@echo "✓ Build complete: $(BINARY_NAME)"

run: ## Run the application locally
	@echo "Starting application..."
	go run ./cmd/main.go

dev: ## Run with hot-reload (requires air)
	@echo "Starting with hot-reload..."
	air

test: ## Run unit tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run ./...

format: ## Format code with gofmt
	@echo "Formatting code..."
	gofmt -s -w .
	@echo "✓ Code formatted"

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "✓ Containers running"
	@echo "Database: localhost:5432"
	@echo "Backend: http://localhost:8080"

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "✓ Containers stopped"

docker-logs: ## View Docker container logs
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) build
	@echo "✓ Image built"

clean: ## Clean build artifacts and temporary files
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "✓ Clean complete"

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	go run ./cmd/main.go migrate up
	@echo "✓ Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	go run ./cmd/main.go migrate down
	@echo "✓ Rollback complete"

seed-db: ## Seed database with sample data
	@echo "Seeding database..."
	@echo "TODO: Implement database seeding"

health-check: ## Check API health
	@echo "Checking API health..."
	curl -s http://localhost:8080/health | jq .
	@echo ""

docs: ## Generate API documentation
	@echo "Generating API documentation..."
	@echo "TODO: Implement documentation generation"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Tools installed"

.DEFAULT_GOAL := help
