# Dongome Marketplace Makefile

# Variables
APP_NAME = dongome
API_BINARY = dongome-api
WORKER_BINARY = dongome-worker
DOCKER_COMPOSE_FILE = docker-compose.yml

# Go related variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GORUN = $(GOCMD) run

# Build variables
BUILD_DIR = ./build
API_SOURCE = ./cmd/api
WORKER_SOURCE = ./cmd/worker

# Docker variables
DOCKER_API_TAG = $(APP_NAME)-api:latest
DOCKER_WORKER_TAG = $(APP_NAME)-worker:latest

# Migration variables
MIGRATE = migrate
MIGRATION_DIR = ./migrations
DB_URL = postgres://dongome:password@localhost:5432/dongome_db?sslmode=disable

.PHONY: help build build-api build-worker clean test lint run run-api run-worker run-docker stop-docker migrate-up migrate-down migrate-create docker-build docker-push

# Default target
help: ## Show this help message
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: build-api build-worker ## Build all binaries

build-api: ## Build API binary
	@echo "Building API binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(BUILD_DIR)/$(API_BINARY) $(API_SOURCE)
	@echo "API binary built successfully!"

build-worker: ## Build worker binary
	@echo "Building worker binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(BUILD_DIR)/$(WORKER_BINARY) $(WORKER_SOURCE)
	@echo "Worker binary built successfully!"

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean completed!"

# Test targets
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Linting
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

# Development run targets
run: run-infrastructure ## Start all services for development
	@echo "Starting all services..."

run-api: ## Run API server locally
	@echo "Starting API server..."
	$(GORUN) $(API_SOURCE)/main.go

run-worker: ## Run worker locally
	@echo "Starting worker..."
	$(GORUN) $(WORKER_SOURCE)/main.go

run-infrastructure: ## Start infrastructure services (DB, Redis, NATS)
	@echo "Starting infrastructure services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d postgres redis nats

# Docker targets
run-docker: ## Start all services with Docker Compose
	@echo "Starting services with Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

stop-docker: ## Stop all Docker services
	@echo "Stopping Docker services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

restart-docker: stop-docker run-docker ## Restart all Docker services

logs-docker: ## View Docker logs
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker build -f ./docker/Dockerfile.api -t $(DOCKER_API_TAG) .
	docker build -f ./docker/Dockerfile.worker -t $(DOCKER_WORKER_TAG) .
	@echo "Docker images built successfully!"

# Database migration targets
migrate-up: ## Run all up migrations
	@echo "Running database migrations up..."
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down: ## Run all down migrations
	@echo "Running database migrations down..."
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DB_URL)" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@echo "Creating new migration: $(name)"
	$(MIGRATE) create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

migrate-force: ## Force migration to specific version (usage: make migrate-force version=VERSION)
	@echo "Forcing migration to version: $(version)"
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DB_URL)" force $(version)

migrate-version: ## Show current migration version
	$(MIGRATE) -path $(MIGRATION_DIR) -database "$(DB_URL)" version

# Go module management
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# Quality assurance
format: ## Format code
	@echo "Formatting code..."
	gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Development setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	$(GOMOD) download
	@echo "Installing development tools..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing migrate tool..."
	$(GOGET) -u -d github.com/golang-migrate/migrate/cmd/migrate
	@echo "Setup complete!"

# Database commands
db-reset: migrate-down migrate-up ## Reset database (down then up)

db-seed: ## Seed database with initial data
	@echo "Seeding database..."
	$(GORUN) ./scripts/seed/main.go

# Production targets
deploy-staging: ## Deploy to staging environment
	@echo "Deploying to staging..."
	# Add staging deployment commands here

deploy-prod: ## Deploy to production environment
	@echo "Deploying to production..."
	# Add production deployment commands here

# Health checks
health-check: ## Check if services are healthy
	@echo "Checking service health..."
	curl -f http://localhost:8080/health || exit 1

# Local development
dev: run-infrastructure ## Start development environment
	@echo "Development environment started!"
	@echo "API will be available at: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "NATS: localhost:4222"
	@echo "NATS Monitor: http://localhost:8222"

# Security
security-scan: ## Run security scan
	@echo "Running security scan..."
	gosec ./...

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	godoc -http=:6060
	@echo "Documentation available at: http://localhost:6060"

# All-in-one targets
all: clean deps lint test build ## Run full CI pipeline

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Tools installed successfully!"