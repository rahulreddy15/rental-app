# Architecture Guide

This document explains the layered architecture used in this codebase and the responsibilities of each layer.

## Overview

```
┌─────────────────────────────────────────────────────────┐
│                      HTTP Request                        │
└─────────────────────────┬───────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────┐
│                      HANDLER                             │
│  • Parse HTTP request (path params, query, body)         │
│  • Validate input                                        │
│  • Call service                                          │
│  • Format HTTP response                                  │
└─────────────────────────┬───────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────┐
│                      SERVICE                             │
│  • Business logic and rules                              │
│  • Create domain objects (IDs, timestamps)               │
│  • Orchestrate multiple repositories                     │
│  • Handle transactions                                   │
│  • Return structured errors (apperr)                     │
└─────────────────────────┬───────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────┐
│                     REPOSITORY                           │
│  • Pure database operations (CRUD)                       │
│  • No business logic                                     │
│  • Return domain errors (ErrNotFound, etc.)              │
└─────────────────────────┬───────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────┐
│                      DATABASE                            │
└─────────────────────────────────────────────────────────┘
```

---

## Layer Details

### 1. Handler Layer (`internal/handler/`)

**Purpose**: Translate between HTTP and the application.

**Responsibilities**:
- Parse path parameters, query strings, request bodies
- Validate input using struct tags
- Call the appropriate service method
- Convert service errors to HTTP responses using `response.FromError()`
- Format successful responses

**Does NOT**:
- Contain business logic
- Directly access the database
- Create domain objects (IDs, timestamps)

**Example**:
```go
func (h *UserHandler) CreateUser(c echo.Context) error {
    // 1. Parse request
    req := new(model.CreateUserRequest)
    if err := c.Bind(req); err != nil {
        return response.BadRequest(c, "Invalid request body", nil)
    }

    // 2. Validate
    if err := c.Validate(req); err != nil {
        return err
    }

    // 3. Call service (no business logic here!)
    user, err := h.userService.Create(c.Request().Context(), service.CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
        Role:  req.Role,
    })

    // 4. Handle errors uniformly
    if err != nil {
        return response.FromError(c, err)
    }

    // 5. Return response
    return response.Created(c, user)
}
```

---

### 2. Service Layer (`internal/service/`)

**Purpose**: Implement business logic and orchestrate operations.

**Responsibilities**:
- Enforce business rules and validations
- Create domain objects with IDs and timestamps
- Coordinate operations across multiple repositories
- Manage transactions for multi-step operations
- Convert repository errors to application errors (`apperr`)

**Does NOT**:
- Know about HTTP (no echo.Context, no status codes)
- Execute raw SQL queries
- Parse or format HTTP requests/responses

**Example**:
```go
func (s *leaseService) CreateLease(ctx context.Context, input CreateLeaseInput) (*model.Lease, error) {
    // Business rule: End date must be after start date
    if !input.EndDate.After(input.StartDate) {
        return nil, apperr.Invalid("End date must be after start date", nil)
    }

    // Business rule: Property must not have active lease
    existingLease, err := s.leaseRepo.GetActiveByProperty(ctx, input.PropertyID)
    if err != nil && !errors.Is(err, repository.ErrLeaseNotFound) {
        return nil, apperr.Internal("Failed to check existing leases", err)
    }
    if existingLease != nil {
        return nil, apperr.Conflict("Property already has an active lease", nil)
    }

    // Create domain object (service owns ID and timestamps)
    lease := &model.Lease{
        ID:         uuid.New(),
        PropertyID: input.PropertyID,
        TenantID:   input.TenantID,
        StartDate:  input.StartDate,
        EndDate:    input.EndDate,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    if err := s.leaseRepo.Create(ctx, lease); err != nil {
        return nil, apperr.Internal("Failed to create lease", err)
    }

    return lease, nil
}
```

**Transactions** (when operations must succeed or fail together):
```go
func (s *leaseService) CreateLeaseWithPayment(ctx context.Context, input CreateLeaseInput) error {
    return s.services.Transaction(func(tx *Services) error {
        // All operations use tx.* services, sharing the same DB transaction
        
        lease, err := tx.Lease.Create(ctx, input)
        if err != nil {
            return err // Rolls back entire transaction
        }

        _, err = tx.Payment.CreateSchedule(ctx, lease.ID, input.MonthlyRent)
        if err != nil {
            return err // Rolls back entire transaction
        }

        return nil // Commits transaction
    })
}
```

---

### 3. Repository Layer (`internal/repository/`)

**Purpose**: Abstract database operations.

**Responsibilities**:
- Execute CRUD operations against the database
- Map between Go structs and database rows
- Return domain-specific errors (`ErrUserNotFound`, `ErrPropertyNotFound`)
- Handle database-specific concerns (query building, pagination)

**Does NOT**:
- Contain business logic or rules
- Know about other entities (no cross-repository calls)
- Create IDs or timestamps (that's the service's job)
- Know about HTTP or application errors

**Example**:
```go
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound  // Domain error, not GORM error
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    // Just save - no ID generation, no business rules
    return r.db.WithContext(ctx).Create(user).Error
}
```

---

## Supporting Packages

### `pkg/apperr/` - Application Errors

Structured errors with codes that map to HTTP statuses:

| Error | HTTP Status | Usage |
|-------|-------------|-------|
| `NotFound` | 404 | Resource doesn't exist |
| `Conflict` | 409 | Duplicate, already exists |
| `Invalid` | 400 | Validation failed, bad input |
| `Unauthorized` | 401 | Not authenticated |
| `Forbidden` | 403 | Not authorized |
| `Internal` | 500 | Unexpected server error |

```go
// In service layer
if existingUser != nil {
    return nil, apperr.Conflict("User with this email already exists", nil)
}
```

### `pkg/response/` - HTTP Responses

Standardized JSON response format:

```go
// Success responses
response.Success(c, data)      // 200 OK
response.Created(c, data)      // 201 Created
response.NoContent(c)          // 204 No Content

// Error responses (manual)
response.BadRequest(c, "message", details)
response.NotFound(c, "message")

// Error responses (from apperr - preferred)
response.FromError(c, err)     // Automatically maps apperr.Code to HTTP status
```

---

## Why This Architecture?

### 1. Testability
Each layer can be tested independently:
- **Handlers**: Mock the service interface
- **Services**: Mock the repository interface
- **Repositories**: Use a test database

### 2. Maintainability
Changes are isolated:
- Change the database? Only repositories change.
- Change HTTP framework? Only handlers change.
- Change business rules? Only services change.

### 3. Reusability
Services can be called from:
- HTTP handlers (REST API)
- gRPC handlers (future)
- Background jobs (future)
- CLI commands (future)

### 4. Clarity
Each layer has a single responsibility:
- "Where does this validation go?" → Service
- "Where do I parse the request?" → Handler
- "Where do I write the SQL?" → Repository

---

## Common Mistakes to Avoid

### ❌ Business logic in handlers
```go
// BAD: Handler checking business rules
func (h *Handler) CreateLease(c echo.Context) error {
    if req.EndDate.Before(req.StartDate) {  // ❌ Business logic!
        return response.BadRequest(c, "Invalid dates", nil)
    }
}
```

### ❌ HTTP concerns in services
```go
// BAD: Service returning HTTP status
func (s *Service) Create(ctx context.Context) (int, error) {
    return 409, errors.New("conflict")  // ❌ HTTP status code!
}
```

### ❌ Cross-repository calls in repository
```go
// BAD: Repository calling another repository
func (r *leaseRepo) Create(ctx context.Context, lease *model.Lease) error {
    user, _ := r.userRepo.GetByID(ctx, lease.TenantID)  // ❌ Wrong layer!
}
```

### ✅ Correct approach
```go
// Service orchestrates, handler translates, repository persists
func (s *leaseService) Create(ctx context.Context, input Input) (*model.Lease, error) {
    // Business validation
    if input.EndDate.Before(input.StartDate) {
        return nil, apperr.Invalid("End date must be after start date", nil)
    }

    // Cross-entity check (service coordinates repositories)
    if _, err := s.userRepo.GetByID(ctx, input.TenantID); err != nil {
        return nil, apperr.NotFound("Tenant not found", err)
    }

    // Create and persist
    lease := &model.Lease{...}
    if err := s.leaseRepo.Create(ctx, lease); err != nil {
        return nil, apperr.Internal("Failed to create lease", err)
    }

    return lease, nil
}
```

---

## Quick Reference

| Question | Answer |
|----------|--------|
| Where do I parse the request body? | Handler |
| Where do I validate business rules? | Service |
| Where do I create UUIDs and timestamps? | Service |
| Where do I call multiple repositories? | Service |
| Where do I handle transactions? | Service |
| Where do I write database queries? | Repository |
| Where do I format HTTP responses? | Handler |
| Where do I convert errors to HTTP status? | Handler (via `response.FromError`) |
