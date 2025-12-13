.PHONY: help build run test clean swagger deps migrate-up migrate-down migrate-create

# Variables
APP_NAME := backend
MAIN_PATH := ./cmd/api
MIGRATIONS_PATH := ./migrations
DB_URL ?= postgres://postgres:postgres@localhost:5432/rental_app?sslmode=disable

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	go mod download
	go mod tidy

swagger: ## Generate swagger documentation
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go -o docs

build: swagger ## Build the application
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

run: swagger ## Run the application
	go run $(MAIN_PATH)/main.go

dev: ## Run with hot reload (requires air)
	@which air > /dev/null || go install github.com/air-verse/air@latest
	air

test: ## Run tests
	go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

# Migration commands
migrate-install: ## Install golang-migrate CLI
	@which migrate > /dev/null || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: migrate-install ## Run all up migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down: migrate-install ## Rollback last migration
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-down-all: migrate-install ## Rollback all migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all

migrate-force: migrate-install ## Force set migration version (usage: make migrate-force V=1)
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(V)

migrate-version: migrate-install ## Show current migration version
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-create: migrate-install ## Create new migration (usage: make migrate-create NAME=create_properties)
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

# Database commands (Docker)
db-up: ## Start PostgreSQL container
	docker compose up -d postgres

db-down: ## Stop PostgreSQL container
	docker compose down

db-logs: ## View PostgreSQL logs
	docker compose logs -f postgres

db-reset: db-down ## Reset database (destroy volume and recreate)
	docker volume rm backend_postgres_data || true
	docker compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	$(MAKE) migrate-up

db-shell: ## Open psql shell in container
	docker compose exec postgres psql -U postgres -d rental_app

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf docs/
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t $(APP_NAME):latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(APP_NAME):latest
