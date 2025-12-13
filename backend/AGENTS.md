# AGENTS.md - Rental Property Management App

## Build/Test Commands (run from `backend/`)
- `make run` - Run server (with swagger generation)
- `make dev` - Run with hot reload (requires air)
- `make test` - Run all tests: `go test -v -race -cover ./...`
- `go test -v -run TestName ./path/to/package` - Run a single test
- `make lint` - Run golangci-lint
- `make swagger` - Generate Swagger docs
- `make db-up` - Start PostgreSQL container
- `make db-down` - Stop PostgreSQL container
- `make migrate-up` - Run migrations
- `make migrate-create NAME=migration_name` - Create new migration

## Architecture
- **Stack**: Go 1.25 + Echo v4 + GORM + PostgreSQL + golang-migrate
- **Layers**: Handler → Service → Repository → Database
- **Structure**:
  - `cmd/api/` - Entrypoint
  - `internal/handler/` - HTTP handlers, route registration
  - `internal/service/` - Business logic, transactions
  - `internal/repository/` - Database operations
  - `internal/model/` - Domain models with GORM tags
  - `internal/database/` - DB connection, migrations
  - `pkg/response/` - Standardized API responses
  - `pkg/apperr/` - Structured error types with codes
  - `migrations/` - SQL migration files

## Code Style
- **Handlers**: Parse request, validate, call service, use `response.FromError(c, err)` for errors
- **Services**: Business logic, create domain objects (IDs, timestamps), return `*apperr.AppError` for expected errors
- **Repositories**: Pure DB operations, return domain errors (e.g., `ErrUserNotFound`)
- **Responses**: Use `pkg/response` helpers (Success, Created, FromError, BadRequest)
- **Errors**: Use `pkg/apperr` (NotFound, Conflict, Invalid, Internal, Unauthorized, Forbidden)
- **Transactions**: Use `services.Transaction(func(txServices *Services) error { ... })`
- **Swagger**: Add annotations (`@Summary`, `@Router`, etc.) to all handlers
- **Imports**: stdlib → external → internal (`backend/...`)
- **Validation**: Use `go-playground/validator` tags on request structs
- **Models**: Define GORM tags + `TableName()` method
- **Naming**: PascalCase exports, camelCase internal

## Documentation
- `ARCHITECTURE.md` - Layer responsibilities, patterns, and examples
- `docs/ADDING_NEW_ENTITY.md` - Step-by-step guide for new entities

## Adding New Entities
Files to touch:
1. `migrations/` - Create up/down SQL files
2. `internal/model/` - Domain struct + request DTOs
3. `internal/repository/` - Interface + implementation, register in `repository.go`
4. `internal/service/` - Business logic, register in `service.go`
5. `internal/handler/` - HTTP handlers, register routes in `routes.go`
