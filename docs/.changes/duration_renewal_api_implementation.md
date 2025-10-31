# Duration Renewal API Implementation Summary

## 📅 Date: October 31, 2025

## 🎯 Objective
Implement **Customer Self-Purchase Duration Package API** allowing customers to renew/purchase additional duration packages by themselves without Sales assistance.

---

## ✅ Changes Made

### 1. SQL Query (`customer_durations.sql`) ✅ UPDATED
**File**: `/internal/infrastructure/db/queries/customer_durations.sql`

**Added Query** (INSERT...SELECT pattern):
```sql
-- name: RenewCustomerDuration :exec
INSERT INTO customer_durations (
  customer_username,
  product_id,
  sales_username,      -- NULL (self-purchase)
  purchase_date,       -- NOW() (TIMESTAMP)
  start_date,          -- NOW() (TIMESTAMP, auto-calculated)
  end_date,            -- NOW() + duration_days (TIMESTAMP, from products)
  price_paid,          -- list_price (full price)
  discount_amount,     -- 0 (no discount)
  status               -- 'ACTIVE'
)
SELECT 
  ?,                   -- customer_username (from JWT)
  ?,                   -- product_id
  NULL,                -- sales_username (self-purchase)
  NOW(),               -- purchase_date
  NOW(),               -- start_date (auto)
  DATE_ADD(NOW(), INTERVAL p.duration_days DAY), -- end_date (auto)
  ?,                   -- price_paid (list_price)
  '0.00',              -- discount_amount
  'ACTIVE'             -- status
FROM products p
WHERE p.id = ?         -- product_id (validate + retrieve duration_days)
  AND p.type = 'DURATION'
  AND p.is_active = 1
  AND p.duration_days IS NOT NULL
LIMIT 1;
```

**Key Improvements:**
- ✅ Changed to `INSERT...SELECT` pattern for atomic operations
- ✅ Auto-calculates `start_date = NOW()` (no user input needed)
- ✅ Auto-calculates `end_date = NOW() + duration_days`
- ✅ Validates product in same query (single atomic operation)
- ✅ Retrieves `duration_days` from products table via JOIN

---

### 2. Request DTO (`customer_duration_request.go`) ✅ UPDATED
**File**: `/domain/requests/customer_duration_request.go`

**Added** (simplified):
```go
type RenewDurationRequest struct {
    ProductID int32 `json:"productId" validate:"required,gt=0"`
    // StartDate removed - auto-calculated by backend as NOW()
}
```

**Changes:**
- ✅ Removed `StartDate` field (backend auto-calculates with NOW())

---

### 3. Response DTO (`customer_duration_response.go`)
**File**: `/domain/responses/customer_duration_response.go`

**Added**:
```go
type RenewDurationResponse struct {
    ID               int32     `json:"id"`
    CustomerUsername string    `json:"customerUsername"`
    ProductID        int32     `json:"productId"`
    ProductName      string    `json:"productName"`
    DurationDays     int32     `json:"durationDays"`
    SalesUsername    *string   `json:"salesUsername"` // NULL
    PurchaseDate     time.Time `json:"purchaseDate"`
    StartDate        time.Time `json:"startDate"`
    EndDate          time.Time `json:"endDate"`
    DaysRemaining    int32     `json:"daysRemaining"`
    PricePaid        float64   `json:"pricePaid"`
    DiscountAmount   float64   `json:"discountAmount"`
    Status           string    `json:"status"`
    Message          string    `json:"message"`
}
```

---

### 4. Repository Interface (`duration_repo.go`) ✅ UPDATED
**File**: `/domain/repositories/duration_repo.go`

**Added Interface Method** (simplified params):
```go
RenewDuration(ctx context.Context, params RenewDurationParams) error

type RenewDurationParams struct {
    CustomerUsername string
    ProductID        int32
    PricePaid        string // DECIMAL string
    // StartDate and DurationDays removed - auto-calculated in SQL
}
```

**Changes:**
- ✅ Removed `StartDate` field (SQL uses NOW())
- ✅ Removed `DurationDays` field (SQL retrieves from products table)

---

### 5. Repository Implementation (`duration_sql.go`) ✅ UPDATED
**File**: `/internal/adapters/repositories/sql/duration_sql.go`

**Added Method** (simplified params):
```go
func (r *CustomerDurationRepository) RenewDuration(ctx context.Context, params repositories.RenewDurationParams) error {
    return r.q.RenewCustomerDuration(ctx, dbmodel.RenewCustomerDurationParams{
        CustomerUsername: sql.NullString{String: params.CustomerUsername, Valid: true},
        ProductID:        sql.NullInt32{Int32: params.ProductID, Valid: true},
        PricePaid:        params.PricePaid,
        ID:               params.ProductID, // WHERE p.id = ? (validate product)
    })
}
```

**Changes:**
- ✅ Removed `StartDate` and `StartDate_2` (SQL uses NOW())
- ✅ Removed `DurationDays` (SQL retrieves from products via JOIN)
- ✅ Added `ID` field for WHERE clause validation

---

### 6. Use Case (`duration_use_case.go`)
**File**: `/domain/usecases/duration_use_case.go`

**Added Method**:
```go
func (u *CustomerDurationUseCase) RenewDuration(
    ctx context.Context, 
    customerUsername string, 
    req requests.RenewDurationRequest,
) (*responses.RenewDurationResponse, error)
```

**Business Logic** (simplified):
1. ✅ Validate product exists and is DURATION type
2. ✅ Validate duration_days > 0
3. ✅ Parse list_price from DECIMAL string
4. ✅ Calculate dates using `time.Now()` for response
5. ✅ Calculate days_remaining
6. ✅ INSERT new duration package (SQL handles date calculations):
   - `sales_username = NULL` (self-purchase)
   - `discount_amount = 0` (no discount)
   - `price_paid = list_price` (full price)
   - `start_date = NOW()` (TIMESTAMP, auto-calculated by SQL)
   - `end_date = NOW() + duration_days` (TIMESTAMP, auto-calculated by SQL)
   - `status = ACTIVE`
7. ✅ Return comprehensive response

**Changes:**
- ❌ Removed start_date parsing from request
- ❌ Removed end_date calculation in Go code
- ✅ SQL now calculates both dates automatically

---

### 7. REST Handler (`duration_rest.go`)
**File**: `/internal/adapters/rest/duration_rest.go`

**Added Handler**:
```go
func (h *CustomerDurationHandler) RenewDuration(c *fiber.Ctx) error
```

**Features** (simplified):
- ✅ JWT Authentication (c.Locals("username"))
- ✅ Request validation (productId only)
- ✅ Error handling (PRODUCT_NOT_FOUND, INVALID_PRODUCT)
- ✅ HTTP 201 Created on success

**Changes:**
- ❌ Removed startDate validation (not needed anymore)

---

### 8. Route Registration (`api_router.go`)
**File**: `/router/api_router.go`

**Added Route**:
```go
durations.Post("/renew", handler.CustomerDuration.RenewDuration) // ✅ ต่ออายุ Duration (ลูกค้าซื้อเอง)
```

**Full Path**: `POST /api/customer-durations/renew`

---

### 9. Documentation Update
**File**: `/docs/CUSTOMER_SELF_PURCHASE_API.md`

**Changes**:
- ✅ Updated Duration API status: ⏳ → ✅ (พร้อมใช้งาน)
- ✅ Added Next Steps section with commands

---

## 📋 Checklist (Following develop_101.md)

- [x] ✅ Step 1: Write SQL Query (`RenewCustomerDuration`)
- [x] ✅ Step 2: Create Request/Response DTOs
- [x] ✅ Step 3: Update Repository Interface
- [x] ✅ Step 4: Implement Repository
- [ ] ⏳ Step 5: Generate sqlc (`make gen-sqlc`)
- [x] ✅ Step 6: Implement Use Case
- [x] ✅ Step 7: Create REST Handler
- [x] ✅ Step 8: Register Route
- [ ] ⏳ Step 9: Generate Wire (`make gen-wire`)
- [ ] ⏳ Step 10: Build & Test

---

## 🔄 Next Steps

### 1. Generate sqlc Code
```bash
cd /Users/pleng/cs-ku/year-3/sa/private-fitness-backend
make gen-sqlc
```

This will generate:
- `dbmodel.RenewCustomerDuration()` method
- `dbmodel.RenewCustomerDurationParams` struct

### 2. Generate Wire Dependencies
```bash
make gen-wire
```

This will update `wire_gen.go` with all dependency injections.

### 3. Build Application
```bash
go build -o tmp/bin/server cmd/app/main.go
```

### 4. Test API

**Login First:**
```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "username": "cust01",
    "password": "password123"
  }'
```

**Test Duration Renewal** (simplified):
```bash
curl -X POST http://localhost:8000/api/customer-durations/renew \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{
    "productId": 1
  }'
```

**Note**: No `startDate` needed - backend auto-calculates as `NOW()`

**Expected Response (201 Created):**
```json
{
  "status": "success",
  "status_code": 201,
  "message": "Duration package renewed successfully",
  "result": {
    "id": 0,
    "customerUsername": "cust01",
    "productId": 1,
    "productName": "Duration 30 days",
    "durationDays": 30,
    "salesUsername": null,
    "purchaseDate": "2025-10-31T15:30:00+07:00",
    "startDate": "2025-10-31T15:30:00+07:00",
    "endDate": "2025-11-30T15:30:00+07:00",
    "daysRemaining": 30,
    "pricePaid": 2500.00,
    "discountAmount": 0,
    "status": "ACTIVE",
    "message": "Duration package renewed successfully"
  }
}
```

**Note**: 
- `startDate` and `endDate` now use TIMESTAMP (includes time component)
- Dates are auto-calculated by SQL as `NOW()` and `NOW() + 30 days`

---

## 🎯 Business Rules Implemented

| Rule | Implementation |
|------|----------------|
| Self-Purchase | `sales_username = NULL` |
| No Discount | `discount_amount = 0.00` |
| Full Price | `price_paid = list_price` |
| Active Status | `status = 'ACTIVE'` |
| Auto Purchase Date | `purchase_date = NOW()` |
| Date Calculation | `end_date = start_date + duration_days` |
| Authentication | JWT token required |

---

## 🔍 Code Pattern Followed

This implementation follows the **exact same pattern** as:
- `POST /api/customer-sessions/renew` (Session Renewal)

**Consistency across codebase:**
1. SQL Query with `:exec` annotation
2. Request/Response DTOs with proper validation
3. Repository interface + implementation
4. Use case with business logic validation
5. REST handler with JWT auth + error handling
6. Route registration in api_router.go

---

## 📊 Database Impact

**Table**: `customer_durations`

**Insert Behavior** (INSERT...SELECT pattern):
```sql
INSERT INTO customer_durations (...)
SELECT 
  ?,                   -- customer_username (from JWT token)
  ?,                   -- product_id (from request)
  NULL,                -- sales_username (self-purchase indicator)
  NOW(),               -- purchase_date (TIMESTAMP)
  NOW(),               -- start_date (TIMESTAMP, auto-calculated)
  DATE_ADD(NOW(), INTERVAL p.duration_days DAY), -- end_date (auto)
  ?,                   -- price_paid (from product.list_price)
  '0.00',              -- discount_amount (no discount)
  'ACTIVE'             -- status
FROM products p
WHERE p.id = ? AND p.type = 'DURATION' AND p.is_active = 1 
  AND p.duration_days IS NOT NULL
```

**Key Features**:
- ✅ **Atomic Operation**: Validation + insertion in single query
- ✅ **Auto-calculated Dates**: Uses NOW() for TIMESTAMP fields
- ✅ **JOIN with Products**: Retrieves duration_days automatically
- ✅ **Product Validation**: Checks type, is_active, duration_days in same query

**No Updates**: This is an INSERT-only operation. Multiple duration packages can coexist for the same customer.

---

## 🧪 Test Scenarios

### Success Cases:
- ✅ Customer purchases valid duration package
- ✅ Multiple packages for same customer
- ✅ Different start dates

### Error Cases:
- ❌ 401 Unauthorized: No JWT token
- ❌ 400 Bad Request: Invalid productId (≤ 0)
- ❌ 404 Not Found: Product not found
- ❌ 400 Bad Request: Product not DURATION type
- ❌ 400 Bad Request: duration_days ≤ 0
- ~~❌ 400 Bad Request: Missing startDate~~ (removed)
- ~~❌ 400 Bad Request: Invalid date format~~ (removed)

---

## 📝 Notes

1. **Pattern Consistency**: This API follows the exact same pattern as Session Renewal for maintainability.

2. **No ID Returned**: The `id` field in response is currently 0. Can be updated to return actual ID using `LAST_INSERT_ID()` if needed.

3. **Frontend Behavior**: Frontend should:
   - Convert price from response (already float64)
   - Display as THB currency
   - Show days_remaining countdown
   - Allow multiple active packages

4. **Database Design**: Customer can have multiple ACTIVE duration packages simultaneously. System doesn't check for existing ACTIVE packages before INSERT.

---

## 🚀 Status

**Implementation**: ✅ Complete  
**Testing**: ⏳ Pending (needs `make gen-sqlc` and `make gen-wire`)  
**Documentation**: ✅ Updated

**Ready for**:
- sqlc generation
- Wire generation
- Build & testing
- Frontend integration

---

**Implemented by**: GitHub Copilot  
**Date**: October 31, 2025  
**Follows**: `/docs/develop_101.md` methodology

---

## 🔄 Update History

### October 31, 2025 - SQL Optimization
**Reason**: Database changed from DATE to TIMESTAMP fields + simplify UX

**Changes Made**:
1. ✅ **SQL Query**: Changed to INSERT...SELECT pattern
   - Auto-calculates `start_date = NOW()`
   - Auto-calculates `end_date = NOW() + duration_days`
   - Validates product in same atomic operation
   - Retrieves `duration_days` from products via JOIN

2. ✅ **Request DTO**: Removed `StartDate` field
   - Frontend only sends `productId`
   - Backend calculates dates automatically

3. ✅ **Repository Params**: Removed `StartDate` and `DurationDays`
   - Only needs: `CustomerUsername`, `ProductID`, `PricePaid`
   - Reduced parameters from 5 to 3

4. ✅ **Repository Implementation**: Updated parameter mapping
   - Added `ID` field for WHERE clause validation
   - Removed `StartDate` and `DurationDays` fields

5. ✅ **Use Case**: Simplified date logic
   - Uses `time.Now()` for response calculations
   - Removed request date parsing
   - SQL handles actual date insertion

6. ✅ **Handler**: Removed startDate validation
   - Only validates `productId > 0`
   - Removed date format error handling

7. ✅ **Documentation**: Updated all examples
   - `/docs/CUSTOMER_SELF_PURCHASE_API.md`
   - `/docs/.changes/duration_renewal_api_implementation.md`

**Benefits**:
- ✅ Simpler frontend (no date picker needed)
- ✅ Fewer parameters (5 → 3)
- ✅ Atomic SQL operation (validation + insertion)
- ✅ Better performance (single query)
- ✅ TIMESTAMP support (includes time component)
- ✅ Consistent with database schema changes
