.PHONY: help build run dev test clean docker-build docker-run docker-stop seed migrate

# Variables
APP_NAME := pos-api
MAIN_PATH := ./cmd/api
BUILD_DIR := ./bin

# Colors for terminal output
GREEN := \033[0;32m
NC := \033[0m

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run: build ## Build and run the application
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

dev: ## Run the application in development mode
	@echo "Running in development mode..."
	@go run $(MAIN_PATH)/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

seed: ## Seed the database with dummy data
	@echo "Seeding database..."
	@go run ./scripts/seed.go
	@echo "Seeding complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .
	@echo "Docker image built: $(APP_NAME):latest"

docker-run: ## Run with Docker Compose
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Containers started"

docker-stop: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped"

docker-logs: ## View Docker logs
	@docker-compose logs -f

docker-seed: ## Seed database in Docker container
	@echo "Seeding database in container..."
	@docker-compose exec api sh -c "cd /app && ./main seed"
