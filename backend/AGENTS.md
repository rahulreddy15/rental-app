# AGENTS.md - Rental Property Management App

## Build/Test Commands (run from `backend/`)
- `make run` - Run server (with swagger generation)
- `make dev` - Run with hot reload (requires air)
- `make test` - Run all tests: `go test -v -race -cover ./...`
- `go test -v -run TestName ./path/to/package` - Run a single test
- `make lint` - Run golangci-lint
- `make swagger` - Generate Swagger docs

## Architecture
- **Backend**: Go 1.25 + Echo v4 framework, Swagger/OpenAPI docs
- **Structure**: `cmd/api/` (entrypoint), `internal/` (private: handler, config, middleware, model, validator), `pkg/` (public: response helpers)
- **API**: RESTful, base path `/api/v1`, Swagger at `/swagger/*`

## Code Style
- Use `pkg/response` helpers (Success, Created, BadRequest, etc.) for consistent API responses
- Add Swagger annotations (`@Summary`, `@Router`, etc.) to all handlers
- Imports: stdlib first, then external packages, then internal (`backend/...`)
- Use `go-playground/validator` for request validation
- Naming: PascalCase for exports, camelCase for internal; descriptive handler/model names
