# 🔄 การอัปเดต CheckBookingPermission Logic

**วันที่:** October 31, 2025  
**ผู้แก้ไข:** Development Team  
**Ticket/Issue:** Update Q2C.1 - Check Permission Logic

---

## 📝 สรุปการเปลี่ยนแปลง

### Logic เดิม:
```sql
-- ตรวจสอบว่ามี Session package ACTIVE และยังมี sessions คงเหลือ
SELECT COUNT(*) as has_permission
FROM customer_sessions
WHERE customer_username = ?
  AND status = 'ACTIVE'
  AND used_sessions < total_sessions;
```

### Logic ใหม่:
```sql
-- Q2C.1: ตรวจสอบเฉพาะว่ามี Session package ACTIVE หรือไม่
-- หมายเหตุ: ถ้าทำครบแล้วจะเปลี่ยน status เป็น 'COMPLETED' โดยอัตโนมัติ
SELECT COUNT(id) as has_permission
FROM customer_sessions
WHERE customer_username = ?
  AND status = 'ACTIVE';
```

---

## 🎯 เหตุผลในการเปลี่ยนแปลง

1. **ไม่ต้องตรวจสอบ remaining sessions ที่ CheckPermission:**
   - ตรวจสอบเฉพาะว่ามีแพ็กเกจ ACTIVE หรือไม่ก็พอ
   - การตรวจสอบ sessions คงเหลือจะทำตอนจองจริง (Q2C.4)

2. **ลดความซับซ้อน:**
   - แค่ต้องการรู้ว่า Customer มีสิทธิ์เข้าหน้า Calendar หรือไม่
   - ถ้ามีแพ็กเกจ ACTIVE = มีสิทธิ์เข้า

3. **Business Rule:**
   - ถ้า sessions ใช้หมดแล้ว (used_sessions = total_sessions)
   - ระบบจะเปลี่ยน status เป็น 'COMPLETED' โดยอัตโนมัติ
   - ดังนั้น status = 'ACTIVE' หมายถึงยังใช้ไม่หมดแน่นอน

4. **ไม่ใช้ COUNT(*):**
   - เปลี่ยนเป็น COUNT(id) เพื่อความชัดเจน
   - Best practice: ควรระบุคอลัมน์แทนการใช้ *

---

## 📂 ไฟล์ที่แก้ไข

### 1. SQL Query
- **File:** `/internal/infrastructure/db/queries/customer_sessions.sql`
- **Query Name:** `CheckBookingPermission`
- **Changes:**
  - ✅ เปลี่ยน `COUNT(*)` → `COUNT(id)`
  - ✅ ลบเงื่อนไข `AND used_sessions < total_sessions`
  - ✅ เพิ่ม comment อธิบาย Q2C.1

### 2. Repository Interface
- **File:** `/domain/repositories/customer_session_repo.go`
- **Function:** `CheckBookingPermission`
- **Changes:**
  - ✅ อัปเดต comment อธิบาย logic ใหม่
  - ✅ เพิ่มหมายเหตุเรื่อง COMPLETED status

### 3. Use Case
- **File:** `/domain/usecases/customer_session_use_case.go`
- **Function:** `CheckBookingPermission`
- **Changes:**
  - ✅ อัปเดต comment อธิบาย logic ใหม่

### 4. Repository Implementation
- **File:** `/internal/adapters/repositories/sql/customer_session_sql.go`
- **Function:** `CheckBookingPermission`
- **Changes:**
  - ✅ อัปเดต comment อธิบาย logic ใหม่
  - ✅ เปลี่ยน comment ใน code จาก "มีสิทธิ์" → "มีแพ็กเกจ ACTIVE"

### 5. REST Handler
- **File:** `/internal/adapters/rest/customer_session_rest.go`
- **Function:** `CheckPermission`
- **Changes:**
  - ✅ อัปเดต comment อธิบาย Q2C.1

### 6. Documentation
- **File:** `/docs/BOOKING_FLOW_FRONTEND_GUIDE.md`
- **Changes:**
  - ✅ อัปเดต Business Logic section (Step 2)
  - ✅ อัปเดต Sequence Diagram
  - ✅ เปลี่ยน SQL query ใน documentation

---

## 🔄 ขั้นตอนต่อไป

### 1. Regenerate SQLC Code
```bash
cd /Users/pleng/cs-ku/year-3/sa/private-fitness-backend
sqlc generate
```

### 2. Build และ Test
```bash
# Build project
go build -o tmp/bin/server cmd/app/main.go

# Run tests
go test ./domain/usecases/... -v
go test ./internal/adapters/... -v
```

### 3. ทดสอบ API
```bash
# Test check permission
curl --location 'http://localhost:8000/api/customers/sessions/check-permission?username=cust01' \
  --header 'Authorization: Bearer YOUR_TOKEN'

# Expected response:
# {
#   "status": "success",
#   "status_code": 200,
#   "message": "Permission check completed",
#   "result": {
#     "hasPermission": true,
#     "canBook": true
#   }
# }
```

### 4. ทดสอบ Edge Cases

**Test Case 1: Customer มีแพ็กเกจ ACTIVE**
```sql
-- Setup
INSERT INTO customer_sessions (customer_username, status, total_sessions, used_sessions)
VALUES ('test_user', 'ACTIVE', 10, 5);

-- Test
SELECT COUNT(id) FROM customer_sessions 
WHERE customer_username = 'test_user' AND status = 'ACTIVE';

-- Expected: 1 (has permission)
```

**Test Case 2: Customer ไม่มีแพ็กเกจ**
```sql
-- Test
SELECT COUNT(id) FROM customer_sessions 
WHERE customer_username = 'new_user' AND status = 'ACTIVE';

-- Expected: 0 (no permission)
```

**Test Case 3: Customer มีแพ็กเกจ COMPLETED**
```sql
-- Setup
INSERT INTO customer_sessions (customer_username, status, total_sessions, used_sessions)
VALUES ('completed_user', 'COMPLETED', 10, 10);

-- Test
SELECT COUNT(id) FROM customer_sessions 
WHERE customer_username = 'completed_user' AND status = 'ACTIVE';

-- Expected: 0 (no permission - status is COMPLETED)
```

**Test Case 4: Customer มีแพ็กเกจ ACTIVE แต่ใช้หมดแล้ว**
```sql
-- Setup (จำลองกรณีที่ยัง ACTIVE แต่ใช้หมด - ไม่ควรเกิดใน production)
INSERT INTO customer_sessions (customer_username, status, total_sessions, used_sessions)
VALUES ('edge_case_user', 'ACTIVE', 10, 10);

-- Test
SELECT COUNT(id) FROM customer_sessions 
WHERE customer_username = 'edge_case_user' AND status = 'ACTIVE';

-- Expected: 1 (has permission - จะตรวจ remaining ตอนจองจริง Q2C.4)
```

---

## 🎯 Impact Analysis

### ✅ Positive Impacts:
1. **Simpler Logic:** Logic ง่ายขึ้น ตรวจเฉพาะ status
2. **Consistent with Business Rule:** สอดคล้องกับกฎที่ COMPLETED = ใช้หมดแล้ว
3. **Better Code Quality:** ไม่ใช้ COUNT(*) ใช้ COUNT(id) แทน
4. **Clear Separation:** แยก permission check กับ remaining sessions check

### ⚠️ Considerations:
1. **ต้องมั่นใจว่า:** ระบบมีกลไกเปลี่ยน status → COMPLETED เมื่อใช้หมด
2. **Q2C.4 ต้องทำงาน:** ตรวจสอบ remaining sessions ตอนจองจริง
3. **Frontend Display:** Frontend อาจต้องโชว์ remaining sessions แยก

---

## 📚 Related Queries

### Q2C.4: Check Remaining Sessions (Used in booking)
```sql
SELECT COUNT(id) AS remaining_session
FROM customer_sessions
WHERE customer_username = ?
  AND status = 'ACTIVE'
  AND (total_sessions - used_sessions) > 0;
```

### Q2C.5: Check Time Slot Available (Used in booking)
```sql
SELECT COUNT(id) AS overlapped_count
FROM training_schedules
WHERE trainer_username = ?
  AND (
    (? < end_time AND ? > start_time)
  );
```

---

## ✅ Checklist

- [x] แก้ SQL query (customer_sessions.sql)
- [x] อัปเดต comment ใน repository interface
- [x] อัปเดต comment ใน use case
- [x] อัปเดต comment ใน repository implementation
- [x] อัปเดต comment ใน REST handler
- [x] อัปเดต documentation (BOOKING_FLOW_FRONTEND_GUIDE.md)
- [ ] Run `sqlc generate`
- [ ] Build และ test
- [ ] ทดสอบ API endpoint
- [ ] ทดสอบ edge cases
- [ ] Review code
- [ ] Deploy to staging

---

**Status:** ✅ Code Changes Complete - Ready for Testing
