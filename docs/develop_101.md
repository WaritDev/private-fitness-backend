# 🛠️ Development Guide: Step-by-Step API Implementation

## ⚠️ IMPORTANT PRINCIPLES

### 🚫 NEVER Write Raw SQL in Repository
**Repository layer ต้อง call methods จาก sqlc เท่านั้น ห้ามเขียน raw SQL โดยเด็ดขาด!**

```go
// ❌ WRONG - ห้ามทำแบบนี้!
func (r *UserRepository) GetUser(ctx context.Context, id int) error {
    query := "SELECT * FROM users WHERE id = ?"
    return r.db.QueryRowContext(ctx, query, id).Scan(...)
}

// ✅ CORRECT - ต้องทำแบบนี้!
func (r *UserRepository) GetUser(ctx context.Context, id int) error {
    return r.q.GetUserById(ctx, id)  // เรียกใช้ method จาก sqlc
}
```

### ✅ The Correct Workflow

```
1. Write SQL in queries/*.sql
   ↓
2. Run: make gen-sqlc
   ↓
3. Use generated methods in repository: r.q.MethodName()
```

**Example**:
```sql
-- queries/users.sql
-- name: CheckUsernameExists :one
SELECT COUNT(*) FROM users WHERE username = ?;
```

```bash
make gen-sqlc  # Generate dbmodel/users.sql.go
```

```go
// repository
func (r *UserRepository) CheckExists(ctx context.Context, username string) (bool, error) {
    count, err := r.q.CheckUsernameExists(ctx, username)  // ✅ Use generated method
    return count > 0, err
}
```

---

## 📋 Table of Contents
1. [Creating a New API Feature](#1-creating-a-new-api-feature)
2. [Database Migrations](#2-database-migrations)
3. [SQL Schema & Queries](#3-sql-schema--queries)
4. [Repository Implementation](#4-repository-implementation)
5. [Use Case Implementation](#5-use-case-implementation)
6. [REST Handler](#6-rest-handler)
7. [Wire Integration](#7-wire-integration)
8. [Complete Example](#8-complete-example-post-users)

---

## 1. Creating a New API Feature

### Overview
เมื่อต้องการสร้าง API endpoint ใหม่ ให้ทำตามขั้นตอนนี้:

```
1. Create Migration → 2. Write Schema → 3. Write Queries → 4. Generate sqlc
                   ↓
5. Create Repository → 6. Create Use Case → 7. Create Handler → 8. Wire Integration
```

---

## 2. Database Migrations

### 2.1 Create New Migration

```bash
make migrate-make name=create_products
```

**Output**:
```
internal/infrastructure/migrations/20251030120000_create_products.go
```

### 2.2 Migration Template

**File**: `internal/infrastructure/migrations/20251030120000_create_products.go`

```go
package migrations

import (
	"context"
	"database/sql"
)

func init() {
	Migrations = append(Migrations, createProducts)
}

var createProducts = &Migration{
	Title: "20251030120000_create_products.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
			CREATE TABLE IF NOT EXISTS products (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				type VARCHAR(50) NOT NULL CHECK (type IN ('DURATION', 'SESSION')),
				category VARCHAR(50) NOT NULL CHECK (category IN ('ECONOMIC', 'BUSINESS', 'FIRST_CLASS')),
				list_price DECIMAL(10, 2) NOT NULL,
				duration_days INT,
				session_amount INT,
				is_active TINYINT(1) DEFAULT 1,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			);
		`)
		return err
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `DROP TABLE IF EXISTS products;`)
		return err
	},
}
```

### 2.3 Run Migration

```bash
# Run all pending migrations
make docker-migrate-up

# Rollback last migration
make docker-migrate-down step=1

# Reset all (drop everything and re-run)
make docker-migrate-reset
```

### 2.4 MariaDB Specific Notes

**DO**:
- ✅ Use inline `ENUM('A', 'B', ...)` for categorical fields (MariaDB ENUM is acceptable in this project)
- ✅ Use `TINYINT(1)` for boolean
- ✅ Use `INT AUTO_INCREMENT` for primary keys
- ✅ Use `DECIMAL(10,2)` for prices
- ✅ Use `TIMESTAMP` with `DEFAULT CURRENT_TIMESTAMP`

**DON'T**:
- ❌ Use `BOOLEAN` - use `TINYINT(1)` instead
- ❌ Use `SERIAL` - use `INT AUTO_INCREMENT` instead
- ❌ Use `now()` - use `CURRENT_TIMESTAMP` instead

---

## 3. SQL Schema & Queries

### 3.1 Write Schema (for sqlc reference)

**File**: `internal/infrastructure/db/schema/create_products.sql`

```sql
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('DURATION', 'SESSION')),
    category VARCHAR(50) NOT NULL CHECK (category IN ('ECONOMIC', 'BUSINESS', 'FIRST_CLASS')),
    list_price DECIMAL(10, 2) NOT NULL,
    duration_days INT,
    session_amount INT,
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**Note**: Schema ซ้ำกับ migration ได้ (migration ใช้ apply ลง DB, schema ใช้ให้ sqlc อ้างอิง)

### 3.2 Write Queries

**File**: `internal/infrastructure/db/queries/products.sql`

```sql
-- name: ListProducts :many
SELECT id, name, type, category, list_price, duration_days, session_amount, is_active, created_at, updated_at
FROM products
WHERE is_active = TRUE
ORDER BY created_at DESC;

-- name: GetProductById :one
SELECT id, name, type, category, list_price, duration_days, session_amount, is_active, created_at, updated_at
FROM products
WHERE id = ?
LIMIT 1;

-- name: CreateProduct :exec
INSERT INTO products (name, type, category, list_price, duration_days, session_amount)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateProduct :exec
UPDATE products
SET name = ?, type = ?, category = ?, list_price = ?, duration_days = ?, session_amount = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteProduct :exec
UPDATE products
SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
```

### 3.3 Query Annotations

```sql
-- name: QueryName :return_type

:many    -- Returns []T (slice of structs)
:one     -- Returns T (single struct) or error if not found
:exec    -- Returns error only (for INSERT/UPDATE/DELETE)
:execrows -- Returns (rows_affected int64, error)
```

### 3.4 Generate sqlc Code

```bash
make gen-sqlc
```

**Output**: `internal/infrastructure/db/dbmodel/`
- `models.go` - Struct definitions
- `querier.go` - Interface
- `products.sql.go` - Generated functions

**Example Generated Code**:
```go
// models.go
type Product struct {
    ID            int32
    Name          string
    Type          string
    Category      string
    ListPrice     string
    DurationDays  sql.NullInt32
    SessionAmount sql.NullInt32
    IsActive      int8
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// products.sql.go
func (q *Queries) ListProducts(ctx context.Context) ([]Product, error)
func (q *Queries) GetProductById(ctx context.Context, id int32) (Product, error)
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) error
```

---

## 4. Repository Implementation

### 4.1 Define Repository Interface

**File**: `domain/repositories/product_repo.go`

```go
package repositories

import (
	"context"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type ProductRepository interface {
	List(ctx context.Context) ([]dbmodel.Product, error)
	GetById(ctx context.Context, id int32) (dbmodel.Product, error)
	Create(ctx context.Context, arg dbmodel.CreateProductParams) error
	Update(ctx context.Context, arg dbmodel.UpdateProductParams) error
	Delete(ctx context.Context, id int32) error
}
```

### 4.2 Implement Repository

**File**: `internal/adapters/repositories/sql/product_sql.go`

```go
package sql

import (
	"context"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type ProductRepository struct {
	q *dbmodel.Queries
}

func ProvideProductRepository(q *dbmodel.Queries) *ProductRepository {
	return &ProductRepository{q: q}
}

func (r *ProductRepository) List(ctx context.Context) ([]dbmodel.Product, error) {
	return r.q.ListProducts(ctx)
}

func (r *ProductRepository) GetById(ctx context.Context, id int32) (dbmodel.Product, error) {
	return r.q.GetProductById(ctx, id)
}

func (r *ProductRepository) Create(ctx context.Context, arg dbmodel.CreateProductParams) error {
	return r.q.CreateProduct(ctx, arg)
}

func (r *ProductRepository) Update(ctx context.Context, arg dbmodel.UpdateProductParams) error {
	return r.q.UpdateProduct(ctx, arg)
}

func (r *ProductRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteProduct(ctx, id)
}
```

---

## 5. Use Case Implementation

### 5.1 Define Request/Response DTOs

**File**: `domain/requests/product_request.go`

```go
package requests

type CreateProductRequest struct {
	Name          string  `json:"name" validate:"required"`
	Type          string  `json:"type" validate:"required,oneof=DURATION SESSION"`
	Category      string  `json:"category" validate:"required,oneof=ECONOMIC BUSINESS FIRST_CLASS"`
	ListPrice     float64 `json:"list_price" validate:"required,gt=0"`
	DurationDays  *int32  `json:"duration_days"`
	SessionAmount *int32  `json:"session_amount"`
}

type UpdateProductRequest struct {
	ID            int32   `json:"id" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Type          string  `json:"type" validate:"required"`
	Category      string  `json:"category" validate:"required"`
	ListPrice     float64 `json:"list_price" validate:"required,gt=0"`
	DurationDays  *int32  `json:"duration_days"`
	SessionAmount *int32  `json:"session_amount"`
}
```

**File**: `domain/responses/product_response.go`

```go
package responses

import "time"

type ProductResponse struct {
	ID            int32     `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Category      string    `json:"category"`
	ListPrice     float64   `json:"list_price"`
	DurationDays  *int32    `json:"duration_days,omitempty"`
	SessionAmount *int32    `json:"session_amount,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
```

### 5.2 Implement Use Case

**File**: `domain/usecases/product_use_case.go`

```go
package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type ProductUseCase struct {
	repo repositories.ProductRepository
}

func ProvideProductUseCase(repo repositories.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (u *ProductUseCase) ListProducts(ctx context.Context) ([]responses.ProductResponse, error) {
	products, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]responses.ProductResponse, len(products))
	for i, p := range products {
		result[i] = u.mapToResponse(p)
	}
	return result, nil
}

func (u *ProductUseCase) GetProductById(ctx context.Context, id int32) (responses.ProductResponse, error) {
	product, err := u.repo.GetById(ctx, id)
	if err != nil {
		return responses.ProductResponse{}, err
	}
	return u.mapToResponse(product), nil
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, req requests.CreateProductRequest) error {
	// Validate business rules
	if req.Type == "DURATION" && req.DurationDays == nil {
		return fmt.Errorf("duration_days is required for DURATION type")
	}
	if req.Type == "SESSION" && req.SessionAmount == nil {
		return fmt.Errorf("session_amount is required for SESSION type")
	}

	arg := dbmodel.CreateProductParams{
		Name:      req.Name,
		Type:      req.Type,
		Category:  req.Category,
		ListPrice: fmt.Sprintf("%.2f", req.ListPrice),
	}

	if req.DurationDays != nil {
		arg.DurationDays = sql.NullInt32{Int32: *req.DurationDays, Valid: true}
	}
	if req.SessionAmount != nil {
		arg.SessionAmount = sql.NullInt32{Int32: *req.SessionAmount, Valid: true}
	}

	return u.repo.Create(ctx, arg)
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, req requests.UpdateProductRequest) error {
	// Check if product exists
	_, err := u.repo.GetById(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	arg := dbmodel.UpdateProductParams{
		ID:        req.ID,
		Name:      req.Name,
		Type:      req.Type,
		Category:  req.Category,
		ListPrice: fmt.Sprintf("%.2f", req.ListPrice),
	}

	if req.DurationDays != nil {
		arg.DurationDays = sql.NullInt32{Int32: *req.DurationDays, Valid: true}
	}
	if req.SessionAmount != nil {
		arg.SessionAmount = sql.NullInt32{Int32: *req.SessionAmount, Valid: true}
	}

	return u.repo.Update(ctx, arg)
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, id int32) error {
	return u.repo.Delete(ctx, id)
}

// Helper: Map dbmodel to response
func (u *ProductUseCase) mapToResponse(p dbmodel.Product) responses.ProductResponse {
	var durationDays, sessionAmount *int32
	if p.DurationDays.Valid {
		durationDays = &p.DurationDays.Int32
	}
	if p.SessionAmount.Valid {
		sessionAmount = &p.SessionAmount.Int32
	}

	return responses.ProductResponse{
		ID:            p.ID,
		Name:          p.Name,
		Type:          p.Type,
		Category:      p.Category,
		ListPrice:     parseFloat(p.ListPrice),
		DurationDays:  durationDays,
		SessionAmount: sessionAmount,
		IsActive:      p.IsActive == 1,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
```

---

## 6. REST Handler

### 6.1 Create Handler

**File**: `internal/adapters/rest/product_rest.go`

```go
package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
)

type ProductHandler struct {
	UC *usecases.ProductUseCase
}

func ProvideProductHandler(uc *usecases.ProductUseCase) *ProductHandler {
	return &ProductHandler{UC: uc}
}

// GET /api/products
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	products, err := h.UC.ListProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Products retrieved successfully",
		"result":      products,
	})
}

// GET /api/products/:id
func (h *ProductHandler) GetProductById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid product ID",
			"result":      nil,
		})
	}

	product, err := h.UC.GetProductById(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "Product not found",
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Product retrieved successfully",
		"result":      product,
	})
}

// POST /api/products
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req requests.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	if err := h.UC.CreateProduct(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Created",
		"status_code": fiber.StatusCreated,
		"message":     "Product created successfully",
		"result":      nil,
	})
}

// PUT /api/products/:id
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid product ID",
			"result":      nil,
		})
	}

	var req requests.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}
	req.ID = int32(id)

	if err := h.UC.UpdateProduct(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Product updated successfully",
		"result":      nil,
	})
}

// DELETE /api/products/:id
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid product ID",
			"result":      nil,
		})
	}

	if err := h.UC.DeleteProduct(c.Context(), int32(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Product deleted successfully",
		"result":      nil,
	})
}
```

### 6.2 Update Handler Aggregator

**File**: `internal/adapters/rest/rest.go`

```go
package rest

type Handler struct {
	User    *UserHandler
	Auth    *AuthHandler
	Product *ProductHandler  // ✅ Add this
}

func ProvideHandler(user *UserHandler, auth *AuthHandler, product *ProductHandler) *Handler {
	return &Handler{
		User:    user,
		Auth:    auth,
		Product: product,  // ✅ Add this
	}
}
```

### 6.3 Register Routes

**File**: `router/api_router.go`

```go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest"
)

func RegisterApiRouter(app *fiber.App, handler *rest.Handler) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handler.Auth.Login)

	// User routes
	users := api.Group("/users")
	users.Get("/", handler.User.List)
	users.Get("/:username", handler.User.GetByUsername)

	// Product routes ✅
	products := api.Group("/products")
	products.Get("/", handler.Product.ListProducts)
	products.Get("/:id", handler.Product.GetProductById)
	products.Post("/", handler.Product.CreateProduct)
	products.Put("/:id", handler.Product.UpdateProduct)
	products.Delete("/:id", handler.Product.DeleteProduct)
}
```

---

## 7. Wire Integration

### 7.1 Update Provider Sets

**File**: `internal/wire/provider.go`

```go
package wire

import (
	"github.com/google/wire"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/repositories/sql"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db"
)

// Infrastructure layer
var InfraSet = wire.NewSet(
	db.ProvideMariaDB,
	db.ProvideQueries,
)

// Repository layer
var RepositorySet = wire.NewSet(
	sql.ProvideUserRepository,
	sql.ProvideAuthRepository,
	sql.ProvideProductRepository,  // ✅ Add this
)

// Use case layer
var ServiceSet = wire.NewSet(
	usecases.ProvideUserUseCase,
	usecases.ProvideAuthUseCase,
	usecases.ProvideProductUseCase,  // ✅ Add this
)

// Handler layer
var HandlerSet = wire.NewSet(
	rest.ProvideUserHandler,
	rest.ProvideAuthHandler,
	rest.ProvideProductHandler,  // ✅ Add this
	rest.ProvideHandler,
)
```

### 7.2 Generate Wire Code

```bash
make gen-wire
```

This will regenerate `internal/wire/wire_gen.go` with all dependencies wired up.

---

## 8. Complete Example: POST /users

### Step-by-Step Implementation

#### 8.1 Migration

```bash
make migrate-make name=create_users
```

```go
// internal/infrastructure/migrations/20251030011453_create_users.go
var createUsers = &Migration{
	Title: "20251030011453_create_users.go",
	Up: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `
			CREATE TABLE IF NOT EXISTS users (
				username VARCHAR(255) PRIMARY KEY,
				password_hash VARCHAR(255) NOT NULL,
				role VARCHAR(50) NOT NULL CHECK (role IN ('CUSTOMER', 'TRAINER', 'SALES', 'ADMIN')),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			);
		`)
		return err
	},
	Down: func(db *sql.DB) error {
		_, err := db.ExecContext(context.Background(), `DROP TABLE IF EXISTS users;`)
		return err
	},
}
```

#### 8.2 Schema

```sql
-- internal/infrastructure/db/schema/create_users.sql
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(255) PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('CUSTOMER', 'TRAINER', 'SALES', 'ADMIN')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 8.3 Queries

```sql
-- internal/infrastructure/db/queries/users.sql

-- name: CreateUser :exec
INSERT INTO users (username, password_hash, role)
VALUES (?, ?, ?);

-- name: GetUserByUsername :one
SELECT username, password_hash, role, created_at, updated_at
FROM users
WHERE username = ?
LIMIT 1;

-- name: ListUsers :many
SELECT username, role, created_at, updated_at
FROM users
ORDER BY created_at DESC;
```

#### 8.4 Generate sqlc

```bash
make gen-sqlc
```

#### 8.5 Request/Response DTOs

```go
// domain/requests/user_request.go
package requests

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=4,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=CUSTOMER TRAINER SALES ADMIN"`
}
```

```go
// domain/responses/user_response.go
package responses

import "time"

type UserResponse struct {
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

#### 8.6 Repository

```go
// internal/adapters/repositories/sql/user_sql.go
package sql

import (
	"context"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserRepository struct {
	q *dbmodel.Queries
}

func ProvideUserRepository(q *dbmodel.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) Create(ctx context.Context, username, passwordHash, role string) error {
	return r.q.CreateUser(ctx, dbmodel.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
		Role:         role,
	})
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (dbmodel.GetUserByUsernameRow, error) {
	return r.q.GetUserByUsername(ctx, username)
}

func (r *UserRepository) List(ctx context.Context) ([]dbmodel.ListUsersRow, error) {
	return r.q.ListUsers(ctx)
}
```

#### 8.7 Use Case

```go
// domain/usecases/user_use_case.go
package usecases

import (
	"context"
	"fmt"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/adapters/repositories/sql"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo *sql.UserRepository
}

func ProvideUserUseCase(repo *sql.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) CreateUser(ctx context.Context, req requests.CreateUserRequest) error {
	// Check if user exists
	_, err := u.repo.GetByUsername(ctx, req.Username)
	if err == nil {
		return fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	return u.repo.Create(ctx, req.Username, string(hashedPassword), req.Role)
}

func (u *UserUseCase) ListUsers(ctx context.Context) ([]responses.UserResponse, error) {
	users, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]responses.UserResponse, len(users))
	for i, user := range users {
		result[i] = responses.UserResponse{
			Username:  user.Username,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return result, nil
}
```

#### 8.8 Handler

```go
// internal/adapters/rest/user_rest.go
package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
)

type UserHandler struct {
	UC *usecases.UserUseCase
}

func ProvideUserHandler(uc *usecases.UserUseCase) *UserHandler {
	return &UserHandler{UC: uc}
}

// POST /api/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req requests.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	if err := h.UC.CreateUser(c.Context(), req); err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "username already exists" {
			statusCode = fiber.StatusConflict
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"status":      "Error",
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Created",
		"status_code": fiber.StatusCreated,
		"message":     "User created successfully",
		"result":      nil,
	})
}

// GET /api/users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.UC.ListUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Users retrieved successfully",
		"result":      users,
	})
}
```

#### 8.9 Router

```go
// router/api_router.go
users := api.Group("/users")
users.Post("/", handler.User.CreateUser)
users.Get("/", handler.User.ListUsers)
```

#### 8.10 Wire Integration

```go
// internal/wire/provider.go
var RepositorySet = wire.NewSet(
	sql.ProvideUserRepository,  // ✅
)

var ServiceSet = wire.NewSet(
	usecases.ProvideUserUseCase,  // ✅
)

var HandlerSet = wire.NewSet(
	rest.ProvideUserHandler,  // ✅
	rest.ProvideHandler,
)
```

```bash
make gen-wire
```

#### 8.11 Test with curl

```bash
# Create user
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!",
    "role": "CUSTOMER"
  }'

# List users
curl http://localhost:8000/api/users
```

---

## 9. Common Patterns & Best Practices

### 9.1 Error Handling

```go
// domain/exceptions/errors.go
package exceptions

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrAlreadyExists  = errors.New("resource already exists")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidInput   = errors.New("invalid input")
)
```

**Usage in Use Case**:
```go
import "github.com/WaritDev/private-fitness-backend/domain/exceptions"

if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
		return exceptions.ErrNotFound
	}
	return err
}
```

**Usage in Handler**:
```go
import (
	"errors"
	"github.com/WaritDev/private-fitness-backend/domain/exceptions"
)

if err != nil {
	statusCode := fiber.StatusInternalServerError
	message := err.Error()

	switch {
	case errors.Is(err, exceptions.ErrNotFound):
		statusCode = fiber.StatusNotFound
	case errors.Is(err, exceptions.ErrAlreadyExists):
		statusCode = fiber.StatusConflict
	case errors.Is(err, exceptions.ErrUnauthorized):
		statusCode = fiber.StatusUnauthorized
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":      "Error",
		"status_code": statusCode,
		"message":     message,
		"result":      nil,
	})
}
```

### 9.2 Validation

```go
// Use struct tags with validator
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=4,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
```

**Usage in Handler**:
```go
if err := ValidateStruct(&req); err != nil {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": err.Error(),
	})
}
```

### 9.3 Transaction Pattern

**❌ WRONG - ห้ามใช้ raw SQL ใน repository:**
```go
// ❌ BAD: Using raw SQL query
func (r *UserRepository) CheckExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	return count > 0, err
}
```

**✅ CORRECT - เขียน SQL ใน queries/*.sql แล้วใช้ sqlc generate:**

**Step 1: Write SQL Query**
```sql
-- internal/infrastructure/db/queries/users.sql

-- name: CheckUsernameExists :one
SELECT COUNT(*) as count
FROM users
WHERE username = ?;
```

**Step 2: Generate with sqlc**
```bash
make gen-sqlc
```

**Step 3: Use generated method in repository**
```go
// internal/adapters/repositories/sql/user_sql.go
func (r *UserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := r.q.CheckUsernameExists(ctx, username)  // ✅ ใช้ method ที่ sqlc generate มาให้
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
```

---

**Transaction Pattern (Multi-table operations):**

**Important**: Repository ต้องมี `*sql.DB` เพื่อสร้าง transaction และใช้ `q.WithTx(tx)` เพื่อรัน queries ภายใน transaction

**File**: `internal/adapters/repositories/sql/user_sql.go`

```go
package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserRepository struct {
	q  *dbmodel.Queries
	db *sql.DB  // ✅ ต้องมี db เพื่อสร้าง transaction
}

func ProvideUserRepository(q *dbmodel.Queries, db *sql.DB) *UserRepository {
	return &UserRepository{q: q, db: db}
}

// CreateUserWithCustomer creates user and customer in a single transaction
func (r *UserRepository) CreateUserWithCustomer(ctx context.Context, userParams CreateUserParams, customerParams CreateCustomerParams) error {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create queries object with transaction
	qtx := r.q.WithTx(tx)

	// Execute first query (CreateUser) - all SQL is in queries/users.sql
	err = qtx.CreateUser(ctx, dbmodel.CreateUserParams{
		Username: userParams.Username,
		Password: userParams.Password,
		Role:     userParams.Role,
		// ... other fields
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Execute second query (CreateCustomer) - all SQL is in queries/customers.sql
	err = qtx.CreateCustomer(ctx, dbmodel.CreateCustomerParams{
		Username:    customerParams.Username,
		HealthInfo:  customerParams.HealthInfo,
		// ... other fields
	})
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
```

**Key Points**:
1. ✅ **ทุก SQL query ต้องอยู่ใน `queries/*.sql`** - ห้ามเขียน raw SQL ใน repository
2. ✅ **Repository ต้องมี `db *sql.DB`** เพื่อเรียก `db.BeginTx()`
3. ✅ **ใช้ `q.WithTx(tx)`** เพื่อสร้าง queries object ที่รัน queries ภายใน transaction
4. ✅ **defer tx.Rollback()** เพื่อ rollback อัตโนมัติถ้าเกิด error ก่อน commit
5. ✅ **tx.Commit()** เพื่อ save changes ถาวร

### 9.4 Pagination

```sql
-- queries/users.sql
-- name: ListUsersPaginated :many
SELECT username, role, created_at
FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
```

```go
type PaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalItems int64       `json:"total_items"`
	TotalPages int64       `json:"total_pages"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}
```

---

## 10. Testing

### 10.1 Unit Test (Use Case)

```go
// domain/usecases/user_use_case_test.go
package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, username, password, role string) error {
	args := m.Called(ctx, username, password, role)
	return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	uc := usecases.ProvideUserUseCase(mockRepo)
	
	mockRepo.On("Create", mock.Anything, "test_user", mock.Anything, "CUSTOMER").Return(nil)

	req := requests.CreateUserRequest{
		Username: "test_user",
		Password: "SecurePass123!",
		Role:     "CUSTOMER",
	}

	// Act
	err := uc.CreateUser(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
```

### 10.2 Integration Test (Handler)

```go
// internal/adapters/rest/user_rest_test.go
package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
	app := fiber.New()
	// Setup routes with test dependencies
	
	reqBody := map[string]string{
		"username": "test_user",
		"password": "SecurePass123!",
		"role":     "CUSTOMER",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}
```

---

## 11. Debugging Tips

### 11.1 Check Generated Code
```bash
cat internal/infrastructure/db/dbmodel/users.sql.go
```

### 11.2 Test SQL Queries Directly
```bash
docker-compose exec mariadb mysql -u app -p private_fitness

# Test query
SELECT * FROM users;
```

### 11.3 Check Wire Dependencies
```bash
cat internal/wire/wire_gen.go
```

### 11.4 View Logs
```bash
docker-compose logs -f api
```

---

## 12. Checklist for New API

- [ ] Create migration file
- [ ] Run migration (`make docker-migrate-up`)
- [ ] Create schema file for sqlc
- [ ] Write SQL queries
- [ ] Generate sqlc code (`make gen-sqlc`)
- [ ] Define repository interface (domain/repositories)
- [ ] Implement repository (internal/adapters/repositories/psql)
- [ ] Create request/response DTOs
- [ ] Implement use case (domain/usecases)
- [ ] Create REST handler (internal/adapters/rest)
- [ ] Update handler aggregator (rest.go)
- [ ] Register routes (router/api_router.go)
- [ ] Update Wire providers (internal/wire/provider.go)
- [ ] Generate Wire code (`make gen-wire`)
- [ ] Test API with curl/Postman

---

---

## 13. API Testing Examples (Auth APIs)

### 13.1 Signup API (POST /api/auth/signup)

**Complete Request with ALL Required Fields**:

```bash
curl -X POST http://localhost:8000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser01",
    "password": "TestPass123!",
    "firstName": "นายทดสอบ",
    "lastName": "ระบบ",
    "gender": "MALE",
    "dateOfBirth": "1995-05-15",
    "phone": "0812345678",
    "email": "testuser01@example.com",
    "healthInfo": "ไม่มีโรคประจำตัว แพ้อาหารทะเล",
    "address": "123 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
    "companyName": "บริษัท เทคโนโลยี จำกัด",
    "companyPosition": "Software Engineer",
    "maritalStatus": "SINGLE",
    "emergencyContactName": "นางสาวทดสอบ ระบบ",
    "emergencyContactRelationship": "พี่สาว",
    "emergencyContactPhone": "0898765432",
    "marketingSource": "Facebook Ads"
  }'
```

**Expected Response (201 Created)**:
```json
{
  "status": "Created",
  "status_code": 201,
  "message": "User registered successfully",
  "result": {
    "ok": true
  }
}
```

**Notes**:
- All fields from `SignupRequest` are required except those marked optional
- `username` must be 4-30 characters, alphanumeric only
- `password` minimum 8 characters (complex validation on frontend)
- `gender` must be: MALE, FEMALE, or OTHER
- `maritalStatus` must be: SINGLE, MARRIED, DIVORCED, or WIDOWED
- `dateOfBirth` format: YYYY-MM-DD
- `email` must be valid email format
- User will automatically be assigned role "CUSTOMER"

---

### 13.2 Login API (POST /api/auth/login)

**Request**:
```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "username": "testuser01",
    "password": "TestPass123!"
  }'
```

**Expected Response (200 OK)**:
```json
{
  "status": "OK",
  "status_code": 200,
  "message": "Login successful",
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "sub": "testuser01",
      "role": "CUSTOMER",
      "firstName": "นายทดสอบ",
      "lastName": "ระบบ"
    }
  }
}
```

**Cookie Set**:
- Cookie name: `pf_auth`
- HTTPOnly: true
- SameSite: Lax
- Max-Age: 604800 (7 days)
- Secure: true (if NODE_ENV=production)

**Notes**:
- JWT token returned in both response body AND HTTP-only cookie
- Token is valid for 7 days
- Use `-c cookies.txt` to save cookies for subsequent requests

---

### 13.3 Me API (GET /api/auth/me)

**Request (with cookie)**:
```bash
curl http://localhost:8000/api/auth/me \
  -b cookies.txt
```

**OR Request (with Authorization header)**:
```bash
curl http://localhost:8000/api/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Expected Response (Authenticated - 200 OK)**:
```json
{
  "status": "OK",
  "status_code": 200,
  "message": "User retrieved successfully",
  "result": {
    "authenticated": true,
    "user": {
      "sub": "testuser01",
      "role": "CUSTOMER",
      "firstName": "นายทดสอบ",
      "lastName": "ระบบ"
    }
  }
}
```

**Expected Response (Not Authenticated - 200 OK)**:
```json
{
  "status": "OK",
  "status_code": 200,
  "message": "User not authenticated",
  "result": {
    "authenticated": false
  }
}
```

**Notes**:
- Token can be provided via cookie OR Authorization header
- Priority: Cookie first, then Authorization header
- No error if not authenticated - returns `authenticated: false`

---

### 13.4 Logout API (POST /api/auth/logout)

**Request**:
```bash
curl -X POST http://localhost:8000/api/auth/logout \
  -b cookies.txt
```

**OR via GET**:
```bash
curl http://localhost:8000/api/auth/logout \
  -b cookies.txt
```

**Expected Response (200 OK)**:
```json
{
  "status": "OK",
  "status_code": 200,
  "message": "Logged out successfully",
  "result": {
    "ok": true
  }
}
```

**Cookie Cleared**:
- Cookie `pf_auth` is cleared by setting MaxAge=-1
- Expires date set to past time

**Notes**:
- Supports both POST and GET methods
- No authentication required (can logout anytime)
- Cookie is cleared on client side

---

### 13.5 Complete Testing Flow

**Sequential Test Script**:

```bash
#!/bin/bash

echo "=== 1. Signup New User ==="
curl -X POST http://localhost:8000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser99",
    "password": "SecurePass123!",
    "firstName": "Test",
    "lastName": "User",
    "gender": "MALE",
    "dateOfBirth": "1990-01-01",
    "phone": "0801234567",
    "email": "testuser99@example.com",
    "healthInfo": "No health issues",
    "address": "123 Test Street, Bangkok 10110",
    "companyName": "Test Company Ltd.",
    "companyPosition": "Tester",
    "maritalStatus": "SINGLE",
    "emergencyContactName": "Emergency Contact",
    "emergencyContactRelationship": "Sibling",
    "emergencyContactPhone": "0809876543",
    "marketingSource": "Google Search"
  }'

echo -e "\n\n=== 2. Login ==="
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "username": "testuser99",
    "password": "SecurePass123!"
  }'

echo -e "\n\n=== 3. Get Current User (Me) ==="
curl http://localhost:8000/api/auth/me -b cookies.txt

echo -e "\n\n=== 4. Logout ==="
curl -X POST http://localhost:8000/api/auth/logout -b cookies.txt

echo -e "\n\n=== 5. Verify Logged Out (Me again) ==="
curl http://localhost:8000/api/auth/me -b cookies.txt

echo -e "\n\nDone!"
```

**Expected Output**:
1. Signup: `status_code: 201`
2. Login: `status_code: 200` + token + cookie set
3. Me: `authenticated: true` + user info
4. Logout: `status_code: 200` + cookie cleared
5. Me (after logout): `authenticated: false`

---

**Happy Coding! 🚀**

````
