# Cancel Appointment Logic (แก้ไขแล้ว ✅)

## Schema Analysis (จากไฟล์ schema)

### 1. training_schedules (การจอง/นัด)
```sql
CREATE TABLE `training_schedules` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,                    -- PK: unique identifier
  `trainer_username` VARCHAR(100),                        -- FK: users.username
  `customer_username` VARCHAR(100),                       -- FK: customers.username
  `session_id` INT,                                       -- FK: customer_sessions.id
  `start_time` TIMESTAMP NOT NULL,
  `end_time` TIMESTAMP NOT NULL,
  `schedule_type` ENUM('APPOINTMENT','DAY_OFF') NOT NULL,
  ...
)
```

### 2. customer_sessions (แพ็กเกจ Sessions)
```sql
CREATE TABLE `customer_sessions` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,                    -- PK: unique identifier
  `customer_username` VARCHAR(100),                       -- FK: customers.username
  `trainer_username` VARCHAR(100),
  `total_sessions` INT NOT NULL,                          -- จำนวน sessions ทั้งหมด
  `used_sessions` INT DEFAULT 0,                          -- จำนวนครั้งที่ใช้ไป
  `status` ENUM('ACTIVE','EXPIRED','CANCELLED','COMPLETED') NOT NULL,
  ...
)
```

### 3. customer_logs (ประวัติการใช้งาน)
```sql
CREATE TABLE `customer_logs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `customer_username` VARCHAR(100),                       -- FK: customers.username
  `log_type` ENUM('CHECK_IN','CHECK_OUT','BOOK_SESSION','CANCEL_SESSION') NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  ...
)
```

---

## ❌ Logic เดิม (ผิด)

```sql
-- ผิด: DELETE ใช้ customer_username ซ้ำซ้อน (id เป็น PK unique อยู่แล้ว)
DELETE FROM training_schedules
WHERE id = ${TargetScheduleId}
  AND customer_username = ${CustomerUsername};

-- ผิด: UPDATE ใช้ customer_username ซ้ำซ้อน (id เป็น PK unique อยู่แล้ว)
UPDATE customer_sessions
SET used_sessions = used_sessions - 1
WHERE session_id = ${SessionId}              -- ⚠️ session_id ไม่มีใน customer_sessions
  AND customer_username = ${CustomerUsername}  -- ⚠️ ไม่จำเป็น id เป็น PK แล้ว
  AND used_sessions > 0;

-- ถูกต้อง: INSERT มี created_at AUTO DEFAULT
INSERT INTO customer_logs (
  customer_username,
  log_type,
  created_at                                 -- ⚠️ ไม่ต้องระบุ มี DEFAULT CURRENT_TIMESTAMP
) VALUES (
  ${Username},
  'CANCEL_SESSION',
  NOW()                                      -- ⚠️ ไม่ต้องใส่ DB handle เอง
);
```

---

## ✅ Logic ที่ถูกต้อง (แก้ไขแล้ว)

### Query 1: Delete Appointment (training_schedules.sql)
```sql
-- name: DeleteAppointment :exec
-- ยกเลิกการจอง - ลบ training_schedule record
-- Logic: DELETE FROM training_schedules WHERE id = ? AND schedule_type = 'APPOINTMENT'
-- Note: ใช้ id เพียงอย่างเดียว เพราะ id เป็น Primary Key (unique)
-- การตรวจสอบ customer_username ทำใน use case layer ก่อนเรียก query นี้
DELETE FROM training_schedules
WHERE id = ?
  AND schedule_type = 'APPOINTMENT';
```

**เหตุผล:**
- `id` เป็น **PRIMARY KEY** → unique อยู่แล้ว ไม่ต้องใช้ `customer_username`
- ตรวจสอบ ownership ใน **Use Case layer** ก่อนเรียก query (Validation 2)
- เพิ่ม `schedule_type = 'APPOINTMENT'` เพื่อป้องกันลบ DAY_OFF โดยไม่ตั้งใจ

---

### Query 2: Decrement Used Sessions (customer_sessions.sql)
```sql
-- name: DecrementUsedSessions :exec
-- ยกเลิกการจอง - ลดจำนวนครั้งที่ใช้ไป (used_sessions - 1)
-- Logic: UPDATE customer_sessions SET used_sessions = used_sessions - 1 WHERE id = ? AND used_sessions > 0
-- Note: id คือ customer_sessions.id (ได้จาก training_schedules.session_id)
-- ต้องมี used_sessions > 0 เพื่อป้องกันค่าติดลบ
UPDATE customer_sessions
SET used_sessions = used_sessions - 1
WHERE id = ?
  AND used_sessions > 0;
```

**เหตุผล:**
- ใช้ `customer_sessions.id` (ได้จาก `training_schedules.session_id`)
- ไม่ต้องใช้ `customer_username` เพราะ `id` เป็น PRIMARY KEY
- เพิ่ม `used_sessions > 0` เพื่อป้องกัน negative value
- **Parameter `id` มาจาก:** `appointment.SessionID` (ดึงจาก GetAppointmentById)

---

### Query 3: Create Customer Log (customer_logs.sql)
```sql
-- name: CreateCustomerLog :exec
INSERT INTO customer_logs (
  customer_username,
  log_type
) VALUES (
  ?, ?
);
```

**เหตุผล:**
- ไม่ต้องระบุ `created_at` เพราะมี `DEFAULT CURRENT_TIMESTAMP`
- DB จะ auto-generate timestamp เอง

---

## Transaction Flow (Use Case Layer)

```go
func (u *BookingUseCase) CancelAppointment(ctx, req) (*CancelAppointmentResponse, error) {
    // ============================================
    // VALIDATION PHASE (ก่อน Transaction)
    // ============================================
    
    // Validation 1: ตรวจสอบว่าการจองนี้มีอยู่จริง
    appointment, err := u.scheduleRepo.GetAppointmentById(ctx, req.AppointmentID)
    if appointment == nil {
        return &Response{Success: false, Message: "Appointment not found"}, nil
    }

    // Validation 2: ตรวจสอบว่าเป็นเจ้าของการจองจริง ⭐ สำคัญ!
    if appointment.CustomerUsername != req.CustomerUsername {
        return &Response{Success: false, Message: "You are not authorized"}, nil
    }

    // Validation 3: ตรวจสอบว่าการจองยังไม่ผ่านไปแล้ว
    if time.Now().After(appointment.StartTime) {
        return &Response{Success: false, Message: "Cannot cancel past appointments"}, nil
    }

    // ============================================
    // TRANSACTION PHASE (3 Operations)
    // ============================================
    tx, _ := u.db.BeginTx(ctx, nil)
    defer tx.Rollback()

    // 1. DELETE training_schedules (ลบการจอง)
    err = u.scheduleRepo.DeleteAppointment(ctx, tx, req.AppointmentID)

    // 2. DECREMENT customer_sessions.used_sessions (คืนสิทธิ์)
    err = u.sessionRepo.DecrementUsedSessions(ctx, tx, appointment.SessionID)

    // 3. INSERT customer_logs (บันทึกประวัติ)
    err = u.customerLogRepo.CreateCustomerLog(ctx, tx, req.CustomerUsername, "CANCEL_SESSION")

    tx.Commit()

    return &Response{Success: true, Message: "Canceled successfully"}, nil
}
```

---

## Key Points (สิ่งที่สำคัญ)

### 1. ✅ ใช้ Primary Key เพียงอย่างเดียว
- `training_schedules.id` → unique
- `customer_sessions.id` → unique
- ไม่ต้องใช้ `customer_username` ซ้ำซ้อน

### 2. ✅ Ownership Validation ใน Use Case
- **ไม่ใส่** `customer_username` ใน WHERE clause ของ DELETE
- **ตรวจสอบ** ownership ใน Use Case layer ก่อนเรียก query
- Separation of Concerns: SQL ทำหน้าที่ DELETE, Use Case ทำหน้าที่ validate

### 3. ✅ Safety Checks
- `AND used_sessions > 0` → ป้องกันค่าติดลบ
- `AND schedule_type = 'APPOINTMENT'` → ป้องกันลบ DAY_OFF

### 4. ✅ Transaction Atomicity
- ทั้ง 3 operations ต้อง **succeed หมด** หรือ **rollback หมด**
- ไม่มีทางที่จะลบ appointment แล้วไม่คืนสิทธิ์ session

### 5. ✅ Auto-Generated Fields
- `customer_logs.created_at` → มี DEFAULT CURRENT_TIMESTAMP
- ไม่ต้องส่ง `NOW()` เข้าไป DB handle เอง

---

## API Endpoint

```http
DELETE /api/bookings/cancel/:id
Content-Type: application/json

{
  "customerUsername": "cust01"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Appointment canceled successfully",
  "appointmentId": 1,
  "customerUsername": "cust01",
  "startTime": "2025-10-31T10:00:00Z",
  "endTime": "2025-10-31T11:00:00Z",
  "sessionId": 5,
  "remainingSessions": 9
}
```

**Response (Error):**
```json
{
  "success": false,
  "message": "You are not authorized to cancel this appointment"
}
```

---

## Next Steps

1. ✅ **แก้ไข SQL queries** (เสร็จแล้ว)
2. 🔄 **Run code generation:**
   ```bash
   cd /Users/pleng/cs-ku/year-3/sa/private-fitness-backend
   make gen-sqlc  # Generate GetAppointmentById, DeleteAppointment, DecrementUsedSessions
   make gen-wire  # Update dependency injection
   docker compose restart api
   ```
3. 🧪 **Test API:**
   ```bash
   # Test valid cancel
   curl -X DELETE http://localhost:8000/api/bookings/cancel/1 \
     -H "Content-Type: application/json" \
     -d '{"customerUsername": "cust01"}'
   
   # Test unauthorized cancel
   curl -X DELETE http://localhost:8000/api/bookings/cancel/1 \
     -H "Content-Type: application/json" \
     -d '{"customerUsername": "wrong_user"}'
   ```

---

## Summary

| Aspect | ❌ เดิม (ผิด) | ✅ ใหม่ (ถูก) |
|--------|--------------|--------------|
| DELETE WHERE | `id = ? AND customer_username = ?` | `id = ? AND schedule_type = 'APPOINTMENT'` |
| UPDATE WHERE | `session_id = ? AND customer_username = ?` | `id = ? AND used_sessions > 0` |
| Ownership Check | ใน SQL WHERE clause | ใน Use Case validation |
| created_at | ใส่ `NOW()` | ไม่ต้องใส่ (AUTO DEFAULT) |
| Parameter | `session_id` (ไม่มีใน table) | `id` (customer_sessions.id) |

✅ **แก้ไขเสร็จแล้ว** - Logic ถูกต้องตาม Schema และ Best Practices!
