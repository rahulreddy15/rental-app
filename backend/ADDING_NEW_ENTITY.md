# Adding a New Entity - Step by Step Guide

This guide walks through adding a new entity (e.g., `Property`) to the codebase.

## Checklist

- [ ] 1. Create migration
- [ ] 2. Create model
- [ ] 3. Create repository
- [ ] 4. Register repository
- [ ] 5. Create service
- [ ] 6. Register service
- [ ] 7. Create handler
- [ ] 8. Register handler & routes
- [ ] 9. Run migration & test

---

## 1. Create Migration

```bash
make migrate-create NAME=create_properties
```

This creates two files in `migrations/`:
- `XXXXXX_create_properties.up.sql`
- `XXXXXX_create_properties.down.sql`

**Edit the UP migration:**

```sql
-- migrations/XXXXXX_create_properties.up.sql
CREATE TABLE properties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    pincode VARCHAR(10) NOT NULL,
    property_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_properties_owner_id ON properties(owner_id);
CREATE INDEX idx_properties_city ON properties(city);
```

**Edit the DOWN migration:**

```sql
-- migrations/XXXXXX_create_properties.down.sql
DROP INDEX IF EXISTS idx_properties_city;
DROP INDEX IF EXISTS idx_properties_owner_id;
DROP TABLE IF EXISTS properties;
```

---

## 2. Create Model

Create `internal/model/property.go`:

```go
package model

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Property struct {
    ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OwnerID      uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
    Name         string    `json:"name" gorm:"type:varchar(255);not null"`
    Address      string    `json:"address" gorm:"type:text;not null"`
    City         string    `json:"city" gorm:"type:varchar(100);not null"`
    State        string    `json:"state" gorm:"type:varchar(100);not null"`
    Pincode      string    `json:"pincode" gorm:"type:varchar(10);not null"`
    PropertyType string    `json:"property_type" gorm:"type:varchar(50);not null"`
    CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`

    // Relationships (optional, for preloading)
    Owner *User `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
}

func (p *Property) BeforeCreate(tx *gorm.DB) error {
    if p.ID == uuid.Nil {
        p.ID = uuid.New()
    }
    return nil
}

func (Property) TableName() string {
    return "properties"
}

// Request DTOs
type CreatePropertyRequest struct {
    Name         string `json:"name" validate:"required,min=2,max=255"`
    Address      string `json:"address" validate:"required"`
    City         string `json:"city" validate:"required,max=100"`
    State        string `json:"state" validate:"required,max=100"`
    Pincode      string `json:"pincode" validate:"required,len=6"`
    PropertyType string `json:"property_type" validate:"required,oneof=apartment flat villa house"`
}

type UpdatePropertyRequest struct {
    Name         string `json:"name" validate:"omitempty,min=2,max=255"`
    Address      string `json:"address" validate:"omitempty"`
    City         string `json:"city" validate:"omitempty,max=100"`
    State        string `json:"state" validate:"omitempty,max=100"`
    Pincode      string `json:"pincode" validate:"omitempty,len=6"`
    PropertyType string `json:"property_type" validate:"omitempty,oneof=apartment flat villa house"`
}
```

---

## 3. Create Repository

Create `internal/repository/property_repository.go`:

```go
package repository

import (
    "context"
    "errors"

    "backend/internal/model"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

var (
    ErrPropertyNotFound = errors.New("property not found")
)

type PropertyRepository interface {
    Create(ctx context.Context, property *model.Property) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.Property, error)
    ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]model.Property, int64, error)
    Update(ctx context.Context, property *model.Property) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type propertyRepository struct {
    db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) PropertyRepository {
    return &propertyRepository{db: db}
}

func (r *propertyRepository) Create(ctx context.Context, property *model.Property) error {
    return r.db.WithContext(ctx).Create(property).Error
}

func (r *propertyRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Property, error) {
    var property model.Property
    if err := r.db.WithContext(ctx).First(&property, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrPropertyNotFound
        }
        return nil, err
    }
    return &property, nil
}

func (r *propertyRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]model.Property, int64, error) {
    var properties []model.Property
    var total int64

    query := r.db.WithContext(ctx).Model(&model.Property{}).Where("owner_id = ?", ownerID)

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&properties).Error; err != nil {
        return nil, 0, err
    }

    return properties, total, nil
}

func (r *propertyRepository) Update(ctx context.Context, property *model.Property) error {
    result := r.db.WithContext(ctx).Save(property)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrPropertyNotFound
    }
    return nil
}

func (r *propertyRepository) Delete(ctx context.Context, id uuid.UUID) error {
    result := r.db.WithContext(ctx).Delete(&model.Property{}, "id = ?", id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrPropertyNotFound
    }
    return nil
}
```

---

## 4. Register Repository

Edit `internal/repository/repository.go`:

```go
package repository

import "gorm.io/gorm"

type Repositories struct {
    User     UserRepository
    Property PropertyRepository  // <-- ADD THIS
}

func NewRepositories(db *gorm.DB) *Repositories {
    return &Repositories{
        User:     NewUserRepository(db),
        Property: NewPropertyRepository(db),  // <-- ADD THIS
    }
}
```

---

## 5. Create Service

Create `internal/service/property_service.go`:

```go
package service

import (
    "context"
    "errors"
    "time"

    "backend/internal/model"
    "backend/internal/repository"
    "backend/pkg/apperr"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type PropertyService interface {
    Create(ctx context.Context, ownerID uuid.UUID, input CreatePropertyInput) (*model.Property, error)
    GetByID(ctx context.Context, id uuid.UUID) (*model.Property, error)
    ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]model.Property, int64, error)
    Update(ctx context.Context, id uuid.UUID, input UpdatePropertyInput) (*model.Property, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type CreatePropertyInput struct {
    Name         string
    Address      string
    City         string
    State        string
    Pincode      string
    PropertyType string
}

type UpdatePropertyInput struct {
    Name         *string
    Address      *string
    City         *string
    State        *string
    Pincode      *string
    PropertyType *string
}

type propertyService struct {
    db           *gorm.DB
    propertyRepo repository.PropertyRepository
    userRepo     repository.UserRepository
}

func NewPropertyService(db *gorm.DB, propertyRepo repository.PropertyRepository, userRepo repository.UserRepository) PropertyService {
    return &propertyService{
        db:           db,
        propertyRepo: propertyRepo,
        userRepo:     userRepo,
    }
}

func (s *propertyService) Create(ctx context.Context, ownerID uuid.UUID, input CreatePropertyInput) (*model.Property, error) {
    // Verify owner exists
    if _, err := s.userRepo.GetByID(ctx, ownerID); err != nil {
        if errors.Is(err, repository.ErrUserNotFound) {
            return nil, apperr.NotFound("Owner not found", err)
        }
        return nil, apperr.Internal("Failed to verify owner", err)
    }

    property := &model.Property{
        ID:           uuid.New(),
        OwnerID:      ownerID,
        Name:         input.Name,
        Address:      input.Address,
        City:         input.City,
        State:        input.State,
        Pincode:      input.Pincode,
        PropertyType: input.PropertyType,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if err := s.propertyRepo.Create(ctx, property); err != nil {
        return nil, apperr.Internal("Failed to create property", err)
    }

    return property, nil
}

func (s *propertyService) GetByID(ctx context.Context, id uuid.UUID) (*model.Property, error) {
    property, err := s.propertyRepo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrPropertyNotFound) {
            return nil, apperr.NotFound("Property not found", err)
        }
        return nil, apperr.Internal("Failed to fetch property", err)
    }
    return property, nil
}

func (s *propertyService) ListByOwner(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]model.Property, int64, error) {
    properties, total, err := s.propertyRepo.ListByOwner(ctx, ownerID, limit, offset)
    if err != nil {
        return nil, 0, apperr.Internal("Failed to fetch properties", err)
    }
    return properties, total, nil
}

func (s *propertyService) Update(ctx context.Context, id uuid.UUID, input UpdatePropertyInput) (*model.Property, error) {
    property, err := s.propertyRepo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrPropertyNotFound) {
            return nil, apperr.NotFound("Property not found", err)
        }
        return nil, apperr.Internal("Failed to fetch property", err)
    }

    if input.Name != nil {
        property.Name = *input.Name
    }
    if input.Address != nil {
        property.Address = *input.Address
    }
    if input.City != nil {
        property.City = *input.City
    }
    if input.State != nil {
        property.State = *input.State
    }
    if input.Pincode != nil {
        property.Pincode = *input.Pincode
    }
    if input.PropertyType != nil {
        property.PropertyType = *input.PropertyType
    }
    property.UpdatedAt = time.Now()

    if err := s.propertyRepo.Update(ctx, property); err != nil {
        return nil, apperr.Internal("Failed to update property", err)
    }

    return property, nil
}

func (s *propertyService) Delete(ctx context.Context, id uuid.UUID) error {
    if err := s.propertyRepo.Delete(ctx, id); err != nil {
        if errors.Is(err, repository.ErrPropertyNotFound) {
            return apperr.NotFound("Property not found", err)
        }
        return apperr.Internal("Failed to delete property", err)
    }
    return nil
}
```

---

## 6. Register Service

Edit `internal/service/service.go`:

```go
package service

import (
    "backend/internal/repository"

    "gorm.io/gorm"
)

type Services struct {
    User     UserService
    Property PropertyService  // <-- ADD THIS
    db       *gorm.DB
}

func NewServices(db *gorm.DB, repos *repository.Repositories) *Services {
    return &Services{
        User:     NewUserService(db, repos.User),
        Property: NewPropertyService(db, repos.Property, repos.User),  // <-- ADD THIS
        db:       db,
    }
}

func (s *Services) Transaction(fn func(txServices *Services) error) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        txRepos := repository.NewRepositories(tx)
        txServices := NewServices(tx, txRepos)
        return fn(txServices)
    })
}
```

---

## 7. Create Handler

Create `internal/handler/property.go`:

```go
package handler

import (
    "strconv"

    "backend/internal/model"
    "backend/internal/service"
    "backend/pkg/response"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
)

type PropertyHandler struct {
    propertyService service.PropertyService
}

func NewPropertyHandler(propertyService service.PropertyService) *PropertyHandler {
    return &PropertyHandler{propertyService: propertyService}
}

type ListPropertiesResponse struct {
    Properties []model.Property `json:"properties"`
    Total      int64            `json:"total"`
    Limit      int              `json:"limit"`
    Offset     int              `json:"offset"`
}

// ListProperties godoc
// @Summary List properties by owner
// @Description Get a paginated list of properties for an owner
// @Tags properties
// @Accept json
// @Produce json
// @Param owner_id query string true "Owner ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=ListPropertiesResponse}
// @Router /properties [get]
func (h *PropertyHandler) ListProperties(c echo.Context) error {
    ownerID, err := uuid.Parse(c.QueryParam("owner_id"))
    if err != nil {
        return response.BadRequest(c, "Invalid owner_id format", nil)
    }

    limit, _ := strconv.Atoi(c.QueryParam("limit"))
    if limit <= 0 || limit > 100 {
        limit = 20
    }

    offset, _ := strconv.Atoi(c.QueryParam("offset"))
    if offset < 0 {
        offset = 0
    }

    properties, total, err := h.propertyService.ListByOwner(c.Request().Context(), ownerID, limit, offset)
    if err != nil {
        return response.FromError(c, err)
    }

    return response.Success(c, ListPropertiesResponse{
        Properties: properties,
        Total:      total,
        Limit:      limit,
        Offset:     offset,
    })
}

// CreateProperty godoc
// @Summary Create a new property
// @Description Create a new property for an owner
// @Tags properties
// @Accept json
// @Produce json
// @Param owner_id query string true "Owner ID"
// @Param property body model.CreatePropertyRequest true "Property details"
// @Success 201 {object} response.Response{data=model.Property}
// @Failure 400 {object} response.ErrorResponse
// @Router /properties [post]
func (h *PropertyHandler) CreateProperty(c echo.Context) error {
    ownerID, err := uuid.Parse(c.QueryParam("owner_id"))
    if err != nil {
        return response.BadRequest(c, "Invalid owner_id format", nil)
    }

    req := new(model.CreatePropertyRequest)
    if err := c.Bind(req); err != nil {
        return response.BadRequest(c, "Invalid request body", nil)
    }

    if err := c.Validate(req); err != nil {
        return err
    }

    property, err := h.propertyService.Create(c.Request().Context(), ownerID, service.CreatePropertyInput{
        Name:         req.Name,
        Address:      req.Address,
        City:         req.City,
        State:        req.State,
        Pincode:      req.Pincode,
        PropertyType: req.PropertyType,
    })
    if err != nil {
        return response.FromError(c, err)
    }

    return response.Created(c, property)
}

// GetProperty godoc
// @Summary Get a property by ID
// @Description Get property details by ID
// @Tags properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success 200 {object} response.Response{data=model.Property}
// @Failure 404 {object} response.ErrorResponse
// @Router /properties/{id} [get]
func (h *PropertyHandler) GetProperty(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return response.BadRequest(c, "Invalid property ID format", nil)
    }

    property, err := h.propertyService.GetByID(c.Request().Context(), id)
    if err != nil {
        return response.FromError(c, err)
    }

    return response.Success(c, property)
}

// UpdateProperty godoc
// @Summary Update a property
// @Description Update property details by ID
// @Tags properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Param property body model.UpdatePropertyRequest true "Property update details"
// @Success 200 {object} response.Response{data=model.Property}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /properties/{id} [put]
func (h *PropertyHandler) UpdateProperty(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return response.BadRequest(c, "Invalid property ID format", nil)
    }

    req := new(model.UpdatePropertyRequest)
    if err := c.Bind(req); err != nil {
        return response.BadRequest(c, "Invalid request body", nil)
    }

    if err := c.Validate(req); err != nil {
        return err
    }

    input := service.UpdatePropertyInput{}
    if req.Name != "" {
        input.Name = &req.Name
    }
    if req.Address != "" {
        input.Address = &req.Address
    }
    if req.City != "" {
        input.City = &req.City
    }
    if req.State != "" {
        input.State = &req.State
    }
    if req.Pincode != "" {
        input.Pincode = &req.Pincode
    }
    if req.PropertyType != "" {
        input.PropertyType = &req.PropertyType
    }

    property, err := h.propertyService.Update(c.Request().Context(), id, input)
    if err != nil {
        return response.FromError(c, err)
    }

    return response.Success(c, property)
}

// DeleteProperty godoc
// @Summary Delete a property
// @Description Delete a property by ID
// @Tags properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.ErrorResponse
// @Router /properties/{id} [delete]
func (h *PropertyHandler) DeleteProperty(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return response.BadRequest(c, "Invalid property ID format", nil)
    }

    if err := h.propertyService.Delete(c.Request().Context(), id); err != nil {
        return response.FromError(c, err)
    }

    return response.NoContent(c)
}
```

---

## 8. Register Handler & Routes

Edit `internal/handler/routes.go`:

```go
package handler

import (
    "backend/internal/service"

    "github.com/labstack/echo/v4"
)

type Handlers struct {
    User     *UserHandler
    Property *PropertyHandler  // <-- ADD THIS
}

func NewHandlers(services *service.Services) *Handlers {
    return &Handlers{
        User:     NewUserHandler(services.User),
        Property: NewPropertyHandler(services.Property),  // <-- ADD THIS
    }
}

func RegisterRoutes(g *echo.Group, handlers *Handlers) {
    g.GET("/health", HealthCheck)

    // User routes
    users := g.Group("/users")
    {
        users.GET("", handlers.User.ListUsers)
        users.POST("", handlers.User.CreateUser)
        users.GET("/:id", handlers.User.GetUser)
        users.PUT("/:id", handlers.User.UpdateUser)
        users.DELETE("/:id", handlers.User.DeleteUser)
    }

    // Property routes  <-- ADD THIS BLOCK
    properties := g.Group("/properties")
    {
        properties.GET("", handlers.Property.ListProperties)
        properties.POST("", handlers.Property.CreateProperty)
        properties.GET("/:id", handlers.Property.GetProperty)
        properties.PUT("/:id", handlers.Property.UpdateProperty)
        properties.DELETE("/:id", handlers.Property.DeleteProperty)
    }
}
```

---

## 9. Run Migration & Test

```bash
# Start database (if not running)
make db-up

# Run migrations
make migrate-up

# Generate swagger docs
make swagger

# Build to verify no errors
go build ./...

# Run the server
make run
```

**Test the endpoints:**

```bash
# Create a user first (to be the owner)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","role":"user"}'

# Create a property (use the user ID from above)
curl -X POST "http://localhost:8080/api/v1/properties?owner_id=<USER_ID>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sunset Apartment",
    "address": "123 MG Road",
    "city": "Bangalore",
    "state": "Karnataka",
    "pincode": "560001",
    "property_type": "apartment"
  }'

# List properties
curl "http://localhost:8080/api/v1/properties?owner_id=<USER_ID>"
```

---

## Quick Reference: Files to Touch

| Step | File(s) |
|------|---------|
| Migration | `migrations/XXXXXX_*.sql` |
| Model | `internal/model/<entity>.go` |
| Repository | `internal/repository/<entity>_repository.go` |
| Register Repo | `internal/repository/repository.go` |
| Service | `internal/service/<entity>_service.go` |
| Register Service | `internal/service/service.go` |
| Handler | `internal/handler/<entity>.go` |
| Register Handler | `internal/handler/routes.go` |
