# AGENTS.md - Rental Property Management App

## Build/Test Commands (run from `backend/`)
- `make run` - Run server (with swagger generation)
- `make dev` - Run with hot reload (requires air)
- `make test` - Run all tests: `go test -v -race -cover ./...`
- `go test -v -run TestName ./path/to/package` - Run a single test
- `make lint` - Run golangci-lint
- `make swagger` - Generate Swagger docs
- `make migrate-up` - Run migrations; `make migrate-create NAME=migration_name` - Create new migration

## Architecture
- **Backend**: Go 1.25 + Echo v4 + GORM + PostgreSQL + golang-migrate
- **Structure**: `cmd/api/` (entrypoint), `internal/` (database, handler, repository, config, middleware, model, validator), `pkg/` (response helpers), `migrations/` (SQL migrations)
- **Patterns**: Repository pattern for DB access; handlers receive repository interfaces via dependency injection
- **API**: RESTful, base path `/api/v1`, Swagger at `/swagger/*`

## Code Style
- Use `pkg/response` helpers (Success, Created, BadRequest, etc.) for consistent API responses
- Add Swagger annotations (`@Summary`, `@Router`, etc.) to all handlers
- Imports: stdlib → external → internal (`backend/...`)
- Use `go-playground/validator` for request validation
- Use repository interfaces for testability; pass context.Context for DB operations
- Naming: PascalCase exports, camelCase internal; models define GORM tags + TableName()
