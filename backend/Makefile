.PHONY: help build run test clean swagger deps

# Variables
APP_NAME := backend
MAIN_PATH := ./cmd/api

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

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

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf docs/
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t $(APP_NAME):latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(APP_NAME):latest
