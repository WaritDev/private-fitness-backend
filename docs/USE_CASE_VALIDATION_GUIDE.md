# 📋 Use Case Validation Guide for System Analysts

> **สำหรับ Senior System Analyst**  
> แนวทางการตรวจสอบ Use Case Description ให้ตรงกับการ Implement จริง  
> โดยเน้นที่ Backend Logic และ SQL Queries  
> Updated: October 31, 2025

---

## 📝 Table of Contents

1. [Overview](#overview)
2. [Validation Checklist](#validation-checklist)
3. [SQL Query Validation Process](#sql-query-validation-process)
4. [SQL Writing Standards](#sql-writing-standards)
5. [Example: Use Case 0S Validation](#example-use-case-0s-validation)
6. [Common Mistakes to Avoid](#common-mistakes-to-avoid)
7. [Validation Report Template](#validation-report-template)

---

## Overview

### หน้าที่ของ System Analyst ในการตรวจสอบ Use Case

1. **ตรวจสอบความสมบูรณ์ของ Use Case Description**
   - Pre-conditions ครบถ้วน
   - Steps ชัดเจน ไม่กำกวม
   - Business rules ถูกต้อง

2. **ตรวจสอบ SQL Queries ที่ใช้ใน Use Case**
   - Query ตรงกับที่เขียนไว้ใน `internal/infrastructure/db/queries/*.sql`
   - Logic ถูกต้องตาม Business rules
   - Performance เหมาะสม (ไม่มี N+1 query problem)

3. **ตรวจสอบ Backend Implementation**
   - Use case flow ตรงกับที่ระบุ
   - Error handling ครบถ้วน
   - Validation logic ถูกต้อง

---

## Validation Checklist

### ✅ Pre-Implementation Validation

- [ ] **Use Case Description สมบูรณ์**
  - มี Use Case Name, Description, Actor, Pre-conditions, Post-conditions
  - Steps มี numbering ชัดเจน
  - แยก Actor actions และ System actions

- [ ] **SQL Queries ถูกต้อง**
  - Query ID (เช่น Q0S.1) มีการระบุชัดเจน
  - SQL syntax ถูกต้อง
  - ใช้ `${variable_name}` แทน `?` หรือ `$1` ในเอกสาร
  - Index ครอบคลุม WHERE clause
  - JOIN logic ถูกต้อง

- [ ] **Business Rules ครบถ้วน**
  - Validation rules ระบุชัดเจน
  - Error messages ระบุครบ
  - Alternative flows (error cases) ครบถ้วน

### ✅ Post-Implementation Validation

- [ ] **SQL Queries Implementation**
  - Query ถูก implement ใน `internal/infrastructure/db/queries/*.sql`
  - Query name ตรงกับ Use Case (เช่น `-- name: GetUserByUsername :one`)
  - Placeholder ใช้ `?` สำหรับ MySQL/MariaDB
  - Generated code ใน `dbmodel/*.sql.go` ถูกต้อง

- [ ] **Repository Implementation**
  - Repository method เรียกใช้ sqlc generated method
  - ไม่มี raw SQL ใน repository
  - Error handling ถูกต้อง

- [ ] **Use Case Implementation**
  - Business logic ตรงกับ Use Case Description
  - Validation ครบถ้วน
  - Error messages ตรงกับที่ระบุ

- [ ] **Handler Implementation**
  - Request/Response DTOs ถูกต้อง
  - HTTP status codes เหมาะสม
  - CORS/Authentication ถูกต้อง

---

## SQL Query Validation Process

### Step 1: อ่าน Use Case Description และจดบันทึก SQL Queries

จาก Use Case Description ให้จดบันทึก:
- Query ID (เช่น Q0S.1, Q0S.2)
- จุดประสงค์ของ query (SELECT, INSERT, UPDATE, DELETE)
- ตารางที่เกี่ยวข้อง
- Conditions และ Validations

**ตัวอย่างจาก Use Case 0S:**

| Query ID | Purpose | Tables | Conditions |
|----------|---------|--------|------------|
| Q0S.1 | Verify username & password | users | WHERE username = ${Username} |
| Q0S.2 | Update last login time | users | WHERE username = ${Username} |

### Step 2: ตรวจสอบว่า Query ถูก Implement แล้วหรือยัง

เปิดไฟล์ `internal/infrastructure/db/queries/*.sql` ที่เกี่ยวข้อง

**ตัวอย่าง: Use Case 0S → ตรวจสอบ `queries/users.sql`**

```bash
# ค้นหา query ที่เกี่ยวข้อง
cat internal/infrastructure/db/queries/users.sql | grep -A 10 "GetUserByUsername"
```

### Step 3: เปรียบเทียบ Logic ระหว่าง Use Case และ Implementation

**Use Case Description (Q0S.1):**
```sql
SELECT 
  Username,
  Role,
  Is_Active
FROM USER
WHERE Username = ${Username}
  AND Password = crypt(${Password}, Password);
```

**Actual Implementation (`queries/users.sql`):**
```sql
-- name: GetUserByUsername :one
SELECT username, password, role, first_name, last_name
FROM users
WHERE username = ?
LIMIT 1;
```

**⚠️ ความแตกต่าง:**
1. Use Case ระบุตรวจสอบ password ใน SQL → **Implementation แยกตรวจสอบใน Go (bcrypt)**
2. Use Case ระบุ Is_Active → **Implementation ไม่มีใน query นี้ (ตรวจสอบใน use case layer)**

**✅ การตรวจสอบ:**
- ตรวจสอบ `auth_use_case.go` ว่ามีการตรวจสอบ password ด้วย bcrypt
- ตรวจสอบว่ามีการตรวจสอบ `is_active` ใน use case

### Step 4: ตรวจสอบ Generated Code

```bash
# ตรวจสอบว่า sqlc generate แล้ว
cat internal/infrastructure/db/dbmodel/users.sql.go | grep "GetUserByUsername"
```

**Expected Output:**
```go
func (q *Queries) GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
```

### Step 5: ตรวจสอบ Repository Implementation

**ไฟล์:** `internal/adapters/repositories/sql/user_sql.go`

**✅ ถูกต้อง:**
```go
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (dbmodel.GetUserByUsernameRow, error) {
    return r.q.GetUserByUsername(ctx, username)  // ✅ ใช้ sqlc method
}
```

**❌ ผิด:**
```go
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (dbmodel.User, error) {
    var user dbmodel.User
    err := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE username = ?", username).Scan(...)  // ❌ Raw SQL
    return user, err
}
```

---

## SQL Writing Standards

### 📌 Use Case Documentation Format

**ใช้ `${variable_name}` สำหรับ parameters ในเอกสาร**

```sql
-- ✅ CORRECT (ใช้ในเอกสาร Use Case)
SELECT username, role, is_active
FROM users
WHERE username = ${Username}
  AND email = ${Email};

-- ❌ WRONG (อย่าใช้ในเอกสาร)
SELECT username, role, is_active
FROM users
WHERE username = ?
  AND email = ?;
```

### 📌 Implementation Format (sqlc queries)

**ใช้ `?` สำหรับ MySQL/MariaDB ใน `queries/*.sql`**

```sql
-- ✅ CORRECT (ใช้ในไฟล์ queries/*.sql)
-- name: GetUserByUsernameAndEmail :one
SELECT username, role, is_active
FROM users
WHERE username = ?
  AND email = ?
LIMIT 1;

-- ❌ WRONG (อย่าใช้ใน MySQL)
SELECT username, role, is_active
FROM users
WHERE username = $1
  AND email = $2;
```

### 📌 SQL Style Guide

**1. Table/Column Naming:**
- Use Case: `USER`, `Username`, `Is_Active` (ตามเอกสาร Business)
- Implementation: `users`, `username`, `is_active` (lowercase, snake_case)

**2. Keywords:**
- ใช้ UPPERCASE: `SELECT`, `FROM`, `WHERE`, `JOIN`, `ORDER BY`

**3. Formatting:**
```sql
-- ✅ CORRECT: Easy to read
SELECT 
  u.username,
  u.role,
  u.is_active,
  c.first_name,
  c.last_name
FROM users u
LEFT JOIN customers c ON u.username = c.username
WHERE u.username = ${Username}
  AND u.is_active = TRUE
ORDER BY u.created_at DESC;

-- ❌ WRONG: Hard to read
SELECT u.username,u.role,u.is_active,c.first_name,c.last_name FROM users u LEFT JOIN customers c ON u.username=c.username WHERE u.username=${Username} AND u.is_active=TRUE ORDER BY u.created_at DESC;
```

---

## Example: Use Case 0S Validation

### 📄 Use Case Description

**Use Case Name:** เข้าสู่ระบบ  
**Use Case ID:** 0S  
**Actor:** Sales (และ CUSTOMER, TRAINER, MANAGER, ADMIN)  
**Pre-Conditions:** 
- ผู้ใช้งานมีบัญชีในระบบ (users table)
- สถานะ `is_active` เป็น `TRUE` (1)

### 📊 Normal Flow

| Step | Actor | System |
|------|-------|--------|
| 1 | ผู้ใช้งานคลิกปุ่ม "Sign In" | |
| 2 | | ระบบ Redirect ไปยังหน้า "Sign In" |
| 3 | ผู้ใช้งานกรอก Username, Password และคลิก "Sign In" | |
| 4 | | ระบบตรวจสอบ Model Validation:<br>• Username ห้ามค่าว่าง<br>• Password ห้ามค่าว่าง |
| 5 | | เมื่อผ่านเงื่อนไข (4) ระบบเข้ารหัสรหัสผ่านด้วย bcrypt<br>❌ หากไม่ผ่าน แสดงข้อความในฟอร์ม (สีแดง) |
| 6 | | ระบบตรวจสอบ Username และ Password ตาม Query **Q0S.1** |
| 7 | | เมื่อผ่านเงื่อนไข (6) ตรวจสอบ `is_active = TRUE`<br>❌ หากไม่ผ่าน แสดง "Username or Password is Incorrect" |
| 8 | | เมื่อผ่านเงื่อนไข (7) ทำการ:<br>• อัปเดต `updated_at` ตาม Query **Q0S.2**<br>• สร้าง JWT Token (7 days expiry)<br>❌ หากไม่ผ่าน แสดง "This account has been suspended." |
| 9 | | ระบบ Redirect ไปยังหน้า Landing Page ตามบทบาท |

### 🔍 SQL Queries Validation

#### **Q0S.1: Verify Username & Password**

**Use Case Description:**
```sql
-- Q0S.1: ตรวจสอบ Username และ Password
SELECT 
  Username,
  Role,
  Is_Active
FROM USER
WHERE Username = ${Username}
  AND Password = crypt(${Password}, Password);
```

**Actual Implementation:** `internal/infrastructure/db/queries/users.sql`
```sql
-- name: GetUserByUsername :one
SELECT username, password, role, first_name, last_name
FROM users
WHERE username = ?
LIMIT 1;
```

**✅ Validation Result:**

| Aspect | Use Case | Implementation | Status |
|--------|----------|----------------|--------|
| Query Name | Q0S.1 | GetUserByUsername | ✅ Match (renamed for clarity) |
| Table | USER | users | ✅ Match (lowercase convention) |
| Username Check | ✓ | ✓ | ✅ Match |
| Password Check | crypt() in SQL | bcrypt in Go | ⚠️ Different approach |
| Is_Active Check | In SQL | In Use Case layer | ⚠️ Different layer |
| Additional Fields | - | first_name, last_name | ℹ️ Extra fields for response |

**📝 Notes:**
1. **Password Verification:** Use Case แนะนำใช้ `crypt()` ใน SQL แต่ implementation ใช้ `bcrypt.CompareHashAndPassword()` ใน Go → **ถูกต้อง** (best practice สำหรับ bcrypt)
2. **Is_Active Check:** Use Case แนะนำตรวจสอบใน SQL แต่ implementation ตรวจสอบใน use case layer → **ยอมรับได้** (แยก business logic ออกจาก data access)

**Go Implementation:** `domain/usecases/auth_use_case.go`
```go
func (u *AuthUseCase) Login(ctx context.Context, req *requests.LoginRequest) (responses.LoginResponse, error) {
    // Step 6: Get user from database
    user, err := u.userRepo.GetByUsername(ctx, req.Username)
    if err != nil {
        return responses.LoginResponse{}, fmt.Errorf("invalid credentials")
    }

    // Step 6: Verify password with bcrypt
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return responses.LoginResponse{}, fmt.Errorf("invalid credentials")
    }

    // Step 7: Check is_active (ตรวจสอบที่นี่แทนใน SQL)
    // Note: ในโค้ดปัจจุบันยังไม่มี - ต้องเพิ่ม!
    
    // Step 8: Sign JWT token
    token, err := u.authRepo.SignUsersAccessToken(&dbmodel.UsersPassport{
        Username: user.Username,
    })
    
    return responses.LoginResponse{Token: token, ...}, nil
}
```

**⚠️ Issue Found:** ยังไม่มีการตรวจสอบ `is_active` ใน use case → **ต้องแก้ไข**

---

#### **Q0S.2: Update Last Login Time**

**Use Case Description:**
```sql
-- Q0S.2: อัปเดต Updated_At เพื่อ track last login
UPDATE USER
SET Updated_At = CURRENT_TIMESTAMP
WHERE Username = ${Username};
```

**Actual Implementation:** `internal/infrastructure/db/queries/users.sql`
```sql
-- name: UpdateUserLastLogin :exec
UPDATE users
SET updated_at = CURRENT_TIMESTAMP
WHERE username = ?;
```

**✅ Validation Result:**

| Aspect | Use Case | Implementation | Status |
|--------|----------|----------------|--------|
| Query Name | Q0S.2 | UpdateUserLastLogin | ✅ Match |
| Operation | UPDATE | UPDATE | ✅ Match |
| Table | USER | users | ✅ Match |
| SET Clause | Updated_At = CURRENT_TIMESTAMP | updated_at = CURRENT_TIMESTAMP | ✅ Match |
| WHERE Clause | Username = ${Username} | username = ? | ✅ Match |

**📝 Notes:**
- Implementation ตรงตาม Use Case Description ✅
- Query name เป็น descriptive (UpdateUserLastLogin) แทน ID (Q0S.2) → **ดีขึ้น**

**Go Implementation:** `internal/adapters/repositories/sql/user_sql.go`
```go
func (r *UserRepository) UpdateLastLogin(ctx context.Context, username string) error {
    return r.q.UpdateUserLastLogin(ctx, username)  // ✅ ใช้ sqlc method
}
```

---

### 📋 Validation Summary for Use Case 0S

| Item | Status | Notes |
|------|--------|-------|
| Use Case Description | ✅ Complete | ครบถ้วน มี pre-conditions, normal flow, queries |
| Q0S.1 Implementation | ⚠️ Partial | Query ถูกต้อง แต่ขาดการตรวจสอบ `is_active` |
| Q0S.2 Implementation | ✅ Complete | ตรงตาม Use Case Description |
| Repository | ✅ Correct | ใช้ sqlc methods ไม่มี raw SQL |
| Use Case Logic | ⚠️ Missing | ขาดการตรวจสอบ `is_active` |
| Error Handling | ✅ Good | Error messages ตรงตาม Use Case |

**🔧 Action Items:**
1. เพิ่มการตรวจสอบ `is_active` ใน `auth_use_case.go`
2. เพิ่มการเรียก `UpdateLastLogin()` หลัง login สำเร็จ
3. เพิ่ม integration test สำหรับ suspended account

---

## Common Mistakes to Avoid

### ❌ Mistake 1: ใช้ Raw SQL ใน Repository

**Wrong:**
```go
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (User, error) {
    var user User
    err := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE username = ?", username).Scan(...)
    return user, err
}
```

**Correct:**
```go
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (dbmodel.GetUserByUsernameRow, error) {
    return r.q.GetUserByUsername(ctx, username)  // ✅ Use sqlc method
}
```

### ❌ Mistake 2: ใช้ Placeholder ผิด

**Wrong (ใน queries/*.sql สำหรับ MySQL):**
```sql
-- name: GetUser :one
SELECT * FROM users WHERE username = $1;  -- ❌ PostgreSQL style
```

**Correct:**
```sql
-- name: GetUser :one
SELECT * FROM users WHERE username = ?;  -- ✅ MySQL style
```

### ❌ Mistake 3: Query ไม่มี LIMIT

**Wrong:**
```sql
-- name: GetUser :one
SELECT * FROM users WHERE username = ?;  -- ❌ ไม่มี LIMIT
```

**Correct:**
```sql
-- name: GetUser :one
SELECT * FROM users WHERE username = ? LIMIT 1;  -- ✅ มี LIMIT
```

### ❌ Mistake 4: Missing Index

**Wrong:**
```sql
-- Query ที่ใช้บ่อย แต่ไม่มี index
SELECT * FROM training_schedules
WHERE customer_username = ${Username}
  AND start_time >= ${StartDate}
ORDER BY start_time;
```

**Correct:**
```sql
-- ต้องมี index
CREATE INDEX idx_training_schedules_customer_time 
ON training_schedules(customer_username, start_time);
```

---

## Validation Report Template

```markdown
# Use Case Validation Report

**Use Case ID:** [เช่น 0S, 1P, 2C]  
**Use Case Name:** [เช่น เข้าสู่ระบบ]  
**Validated By:** [ชื่อ System Analyst]  
**Date:** [วันที่ตรวจสอบ]

---

## 1. Use Case Description Review

- [ ] Use Case Description สมบูรณ์
- [ ] Pre-conditions ชัดเจน
- [ ] Normal flow ครบถ้วน
- [ ] Alternative flows ครบถ้วน
- [ ] SQL Queries ระบุชัดเจน

**Issues Found:**
- [ระบุปัญหาที่พบ]

---

## 2. SQL Queries Validation

### Query Q[ID].1: [ชื่อ Query]

**Use Case Description:**
```sql
[SQL จาก Use Case]
```

**Actual Implementation:** `queries/[table].sql`
```sql
[SQL จริง]
```

**Validation Result:**

| Aspect | Use Case | Implementation | Status |
|--------|----------|----------------|--------|
| [aspect] | [expected] | [actual] | [✅/⚠️/❌] |

**Notes:**
- [บันทึกข้อสังเกต]

---

## 3. Backend Implementation Review

### Repository Layer

- [ ] ใช้ sqlc generated methods
- [ ] ไม่มี raw SQL
- [ ] Error handling ถูกต้อง

**Issues Found:**
- [ระบุปัญหา]

### Use Case Layer

- [ ] Business logic ตรงตาม Use Case
- [ ] Validation ครบถ้วน
- [ ] Error messages ตรงตาม Use Case

**Issues Found:**
- [ระบุปัญหา]

### Handler Layer

- [ ] Request/Response DTOs ถูกต้อง
- [ ] HTTP status codes เหมาะสม
- [ ] Authentication/Authorization ถูกต้อง

**Issues Found:**
- [ระบุปัญหา]

---

## 4. Summary

**Overall Status:** [✅ Pass / ⚠️ Pass with Issues / ❌ Fail]

**Action Items:**
1. [รายการสิ่งที่ต้องแก้ไข]
2. [...]

**Approved By:** [ชื่อ]  
**Date:** [วันที่]
```

---

## Next Steps

1. **อ่าน Use Case Description** ที่ต้องการตรวจสอบ
2. **จดบันทึก SQL Queries** และ Business Logic
3. **ตรวจสอบ Implementation** ตาม Checklist
4. **เขียน Validation Report** ตาม Template
5. **สร้าง Action Items** สำหรับสิ่งที่ต้องแก้ไข

---

**Happy Validating! 🔍**
