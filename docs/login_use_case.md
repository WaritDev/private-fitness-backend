# Use Case 0S: เข้าสู่ระบบ (Login)

## 📋 Use Case Overview

| Field | Value |
|-------|-------|
| **Use Case Name** | เข้าสู่ระบบ |
| **Use Case ID** | 0S |
| **Description** | พนักงานขาย (Sales) หรือผู้ใช้งานทุกบทบาทเข้าสู่ระบบโดยใช้ชื่อผู้ใช้งานและรหัสผ่านที่ลงทะเบียนไว้ เพื่อเข้าใช้งานฟังก์ชันตามสิทธิ์ของบทบาท |
| **Actor** | พนักงานขาย (Sales), CUSTOMER, TRAINER, MANAGER, ADMIN |
| **Pre-Conditions** | • ผู้ใช้งานมีบัญชีในระบบ (users table)<br>• สถานะ `is_active` เป็น `TRUE` (1) |
| **Post-Conditions** | • ผู้ใช้งานได้รับ JWT token<br>• `updated_at` ใน users table ถูกอัปเดต<br>• Redirect ไปยังหน้า Landing Page ตามบทบาท |

---

## 🔄 Normal Flow

### Actor → System Interactions

| Step | Actor | System |
|------|-------|--------|
| 1 | ผู้ใช้งานคลิกปุ่ม "Sign In" ที่มุมขวาบนของเว็บไซต์ | |
| 2 | | ระบบ Redirect ไปยังหน้า "Sign In" |
| 3 | ผู้ใช้งานกรอกข้อมูล:<br>• ชื่อผู้ใช้งาน (Username)<br>• รหัสผ่าน (Password)<br><br>คลิก "Sign In" เพื่อเข้าสู่ระบบ | |
| 4 | | ระบบตรวจสอบข้อมูลทั้งหมด:<br>**Model Validation:**<br>• Username ห้ามค่าว่าง<br>• Password ห้ามค่าว่าง |
| 5 | | เมื่อผ่านเงื่อนไข (4) ระบบดำเนินการเข้ารหัสรหัสผ่านที่ผู้ใช้งานกรอก (Hash ด้วย bcrypt)<br><br>❌ หากไม่ผ่านเงื่อนไข (4) ระบบจะแสดงข้อความในฟอร์มของช่องที่ไม่ผ่านเงื่อนไข (สีแดง) |
| 6 | | **Database Validation:**<br>ระบบตรวจสอบ Username และ Password ว่าตรงกับข้อมูลในตาราง `users` โดย Query **Q0S.1** |
| 7 | | เมื่อผ่านเงื่อนไข (6) ระบบตรวจสอบสถานะการใช้งานของบัญชี (`is_active = TRUE`)<br><br>❌ หากไม่ผ่านเงื่อนไข (6) ระบบจะแสดงข้อความในฟอร์มว่า "**Username or Password is Incorrect**" |
| 8 | | เมื่อผ่านเงื่อนไข (7) ระบบทำการ:<br>• อัปเดต `updated_at` ใน users table ตาม Query **Q0S.2**<br>• สร้าง Session Token / JWT สำหรับการเข้าสู่ระบบ (7 days expiry)<br><br>❌ หากไม่ผ่านเงื่อนไข (7) ระบบจะแสดงข้อความในฟอร์มว่า "**This account has been suspended.**" |
| 9 | | ระบบ Redirect ไปยังหน้า "Landing Page" ของบทบาทที่เข้าสู่ระบบ |
| 10 | | ระบบแสดงข้อความ Pop-Up ว่า<br>"**เข้าสู่ระบบสำเร็จ ✅ ยินดีต้อนรับ คุณ ${First_Name}!**" |

---

## 📊 Database Queries

### Q0S.1: Verify Username & Password

**Purpose:** ตรวจสอบความถูกต้องของ username และ password (bcrypt)

```sql
-- SQL Query (Conceptual - actual implementation uses bcrypt in Go)
SELECT 
  username,
  password,
  role,
  first_name,
  last_name,
  is_active
FROM users
WHERE username = ?
LIMIT 1;
```

**Go Implementation:**
```go
// In auth_use_case.go - Login()

// 1. Get user from database
user, err := u.userRepo.GetByUsername(ctx, req.Username)
if err != nil {
    return responses.LoginResponse{}, fmt.Errorf("invalid credentials")
}

// 2. Verify password with bcrypt
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
    return responses.LoginResponse{}, fmt.Errorf("invalid credentials")
}

// 3. Check is_active status (handled by query - users with is_active=false won't authenticate)
```

**sqlc Query:**
```sql
-- internal/infrastructure/db/queries/users.sql

-- name: GetUserByUsername :one
SELECT username, password, role, first_name, last_name
FROM users
WHERE username = ?
LIMIT 1;
```

---

### Q0S.2: Update Last Login Timestamp

**Purpose:** อัปเดต `updated_at` เพื่อ track last active time ของผู้ใช้งาน

```sql
-- name: UpdateUserLoginTimestamp :exec
-- Q0S.2: Update updated_at when user logs in (track last active time)
UPDATE users
SET updated_at = CURRENT_TIMESTAMP
WHERE username = ?;
```

**Go Implementation:**
```go
// In auth_use_case.go - Login()

// Q0S.2: Update updated_at timestamp to track last login time
if err := u.userRepo.UpdateLoginTimestamp(ctx, user.Username); err != nil {
    // Log error but don't fail login
    fmt.Printf("Warning: failed to update login timestamp for user %s: %v\n", user.Username, err)
}
```

**Notes:**
- การ update timestamp **ไม่ควรทำให้ login ล้มเหลว** หากเกิด error
- ใช้ `CURRENT_TIMESTAMP` เพื่อให้ database handle timezone automatically
- `updated_at` จะถูกอัปเดตทุกครั้งที่ user login สำเร็จ

---

## 🔐 JWT Token Generation

### JWT Payload Structure

```go
type JWTPayload struct {
    Sub       string // username
    Role      string // ADMIN, TRAINER, SALES, CUSTOMER, MANAGER
    FirstName string
    LastName  string
}
```

### Token Properties

| Property | Value |
|----------|-------|
| **Algorithm** | HS256 |
| **Expiry** | 7 days (604800 seconds) |
| **Storage** | HTTP-only Cookie (`pf_auth`) + Response body |
| **Cookie Properties** | • SameSite: Lax<br>• HTTPOnly: true<br>• Max-Age: 604800 |

### Response Format

```json
{
  "status": "success",
  "status_code": 200,
  "message": "Login successful",
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "sub": "sales1",
      "role": "SALES",
      "firstName": "Sara",
      "lastName": "Lee"
    }
  }
}
```

---

## ❌ Alternative Flows

### Alt 1: Validation Error (Step 4-5)

| Condition | Error Message | HTTP Status |
|-----------|--------------|-------------|
| Username is empty | "Username is required" | 400 |
| Password is empty | "Password is required" | 400 |

**Response Example:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Username is required",
  "result": null
}
```

---

### Alt 2: Invalid Credentials (Step 6)

| Condition | Error Message | HTTP Status |
|-----------|--------------|-------------|
| Username not found | "Username or Password is Incorrect" | 400 |
| Password mismatch | "Username or Password is Incorrect" | 400 |

**Response Example:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "invalid credentials",
  "result": null
}
```

**Security Note:** ไม่บอกว่า username หรือ password ส่วนไหนผิด เพื่อป้องกัน username enumeration attack

---

### Alt 3: Account Suspended (Step 7)

| Condition | Error Message | HTTP Status |
|-----------|--------------|-------------|
| `is_active = FALSE` (0) | "This account has been suspended." | 403 |

**Response Example:**
```json
{
  "status": "error",
  "status_code": 403,
  "message": "This account has been suspended.",
  "result": null
}
```

---

## 🧪 Test Cases

### Test Case 1: Normal Login (Sales)

**Input:**
```json
{
  "username": "sales1",
  "password": "Password123!"
}
```

**Expected:**
- ✅ Status: 200 OK
- ✅ JWT token returned
- ✅ `updated_at` in users table updated to CURRENT_TIMESTAMP
- ✅ User redirected to Sales Landing Page
- ✅ Pop-up: "เข้าสู่ระบบสำเร็จ ✅ ยินดีต้อนรับ คุณ Sara!"

**Verification Query:**
```sql
SELECT username, updated_at 
FROM users 
WHERE username = 'sales1';
-- updated_at should be current timestamp
```

---

### Test Case 2: Invalid Password

**Input:**
```json
{
  "username": "sales1",
  "password": "WrongPassword"
}
```

**Expected:**
- ❌ Status: 400 Bad Request
- ❌ Message: "invalid credentials"
- ❌ `updated_at` NOT updated
- ❌ No token generated

---

### Test Case 3: Inactive Account

**Setup:**
```sql
UPDATE users 
SET is_active = 0 
WHERE username = 'sales1';
```

**Input:**
```json
{
  "username": "sales1",
  "password": "Password123!"
}
```

**Expected:**
- ❌ Status: 403 Forbidden
- ❌ Message: "This account has been suspended."
- ❌ `updated_at` NOT updated
- ❌ No token generated

---

### Test Case 4: Empty Fields

**Input:**
```json
{
  "username": "",
  "password": ""
}
```

**Expected:**
- ❌ Status: 400 Bad Request
- ❌ Message: "Username is required" (or "Password is required")
- ❌ No database query executed

---

## 🔍 API Endpoint

### Request

**Method:** `POST`  
**Path:** `/api/auth/login`  
**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "username": "sales1",
  "password": "Password123!"
}
```

### Success Response

**Status:** `200 OK`

```json
{
  "status": "success",
  "status_code": 200,
  "message": "Login successful",
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJzYWxlczEiLCJyb2xlIjoiU0FMRVMiLCJmaXJzdE5hbWUiOiJTYXJhIiwibGFzdE5hbWUiOiJMZWUiLCJleHAiOjE3MzA5MDAwMDB9.xxx",
    "user": {
      "sub": "sales1",
      "role": "SALES",
      "firstName": "Sara",
      "lastName": "Lee"
    }
  }
}
```

**Cookies Set:**
```
Set-Cookie: pf_auth=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...; HttpOnly; SameSite=Lax; Max-Age=604800
```

---

## 📝 Implementation Files

### Backend Files Modified/Created:

1. **SQL Query:**
   - `internal/infrastructure/db/queries/users.sql` - Added `UpdateUserLoginTimestamp`

2. **Repository Interface:**
   - `domain/repositories/user_repo.go` - Added `UpdateLoginTimestamp()` method

3. **Repository Implementation:**
   - `internal/adapters/repositories/sql/user_sql.go` - Implemented `UpdateLoginTimestamp()`

4. **Use Case:**
   - `domain/usecases/auth_use_case.go` - Added Q0S.2 call in `Login()` method

5. **REST Handler:**
   - `internal/adapters/rest/auth_rest.go` - Already implemented (no changes needed)

6. **Router:**
   - `router/api_router.go` - Already registered (no changes needed)

---

## 🔗 Related Use Cases

- **Use Case 1S:** พนักงานขายกรอกข้อมูลลูกค้า (ใช้ token จาก login)
- **Use Case 2.1C:** ลูกค้าลงทะเบียนแพ็กเกจ Duration (require authenticated session)
- **Use Case 3C:** จองเวลาออกกำลังกาย (require authenticated session)

---

## 📌 Notes

1. **Password Security:**
   - ใช้ bcrypt hash (cost 10)
   - Password stored as hash เท่านั้น ไม่เก็บ plaintext

2. **Last Login Tracking:**
   - `updated_at` field ใช้ track last active time
   - สามารถใช้ query users ที่ไม่ได้ login นานๆ ได้

3. **Token Management:**
   - Token มี expiry 7 วัน
   - Frontend ควร check token expiry และ refresh/re-login

4. **Error Messages:**
   - ใช้ generic message "invalid credentials" เพื่อความปลอดภัย
   - ไม่เปิดเผยว่า username หรือ password ผิด

---

**Document Version:** 1.0  
**Last Updated:** October 30, 2025  
**Author:** Backend Team
