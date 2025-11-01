# Review: Use Case 4C และ 3P - สรุปการแก้ไข

## Use Case 4C: สแกนเข้า Fitness

### ❌ ปัญหาที่พบ

**Q4C.2 - Session Package Check-in:**

ใน Use Case Description คุณเขียนว่า:

```sql
INSERT INTO customer_logs (...)
VALUES (...);

UPDATE customer_sessions
SET used_sessions = used_sessions + 1,
    updated_at = NOW()
WHERE customer_username = ${customer_username}
  AND status = 'ACTIVE'
  AND used_sessions < total_sessions
ORDER BY created_at DESC
LIMIT 1;
```

**ปัญหาคือ:** Flow นี้จะหัก session ทันทีเมื่อสแกน QR Code ซึ่งไม่ตรงกับ Hybrid Flow ที่เราพัฒนาไว้

**Flow ที่ถูกต้อง (ตามโค้ดจริง):**
- เมื่อสแกน QR Code สำหรับ Session → สร้าง **PENDING log** เท่านั้น
- ไม่หัก session ทันที
- รอ Trainer ยืนยันก่อน (Use Case 3P)

### ✅ SQL ที่ใช้จริง

**Q4C.1 - Duration Package (ถูกต้องแล้ว):**
```sql
INSERT INTO customer_logs (customer_username, log_type)
VALUES (?, 'CHECK_IN');
```

**Q4C.2 - Session Package (ต้องแก้ไข):**
```sql
-- สร้าง PENDING log เท่านั้น (ไม่หัก session)
INSERT INTO customer_logs (
  customer_username,
  log_type,
  status,
  schedule_id
) VALUES (
  ?, 'CHECK_IN', 'PENDING', ?
);
```

**หมายเหตุ:** 
- ต้องหา `schedule_id` จาก `training_schedules` ก่อน (โดยใช้ `GetCustomerScheduleForToday`)
- ถ้าไม่มี schedule สำหรับวันนี้ → ไม่สามารถ check-in ได้

### ✅ แก้ไข Use Case Description

**Step 6 สำหรับ Session Package:**

```
กรณีลูกค้าใช้สิทธิ์ Sessions:

ระบบตรวจสอบว่ามี schedule สำหรับวันนี้หรือไม่:
- ถ้ามี → บันทึก pending check-in log (Query Q4C.2)
- ถ้าไม่มี → แสดงข้อความ "No appointment scheduled for today"

Q4C.2
INSERT INTO customer_logs (
  customer_username,
  log_type,
  status,
  schedule_id
) VALUES (
  ${customer_username}, 'CHECK_IN', 'PENDING', ${schedule_id}
);
```

**Step 7:**
```
7. ระบบแสดงข้อความ:
   - Duration: "Welcome, [First_Name]!"
   - Session: "Check-in pending. Waiting for trainer confirmation, [First_Name]!"
```

---

## Use Case 3P: ยืนยันการเข้าเรียน

### ❌ ปัญหาที่พบ

#### 1. **Q3P.1 - ถูกต้องแล้ว ✅**

Query ตรงกับโค้ดจริงทุกส่วน

#### 2. **Q3P.2 - ต้องแก้ไข**

**ใน Use Case Description คุณเขียน:**
```sql
SELECT COUNT(cs."Session_Id") AS Remaining_Session
FROM "CUSTOMER_SESSION" cs
WHERE cs."Session_Id" = ${SessionId}
AND cs."Customer_Username" = ${CustomerUsername}
AND cs."Status" = 'ACTIVE'
AND (cs."Total_Sessions" - cs."Used_Sessions") > 0;
```

**ปัญหาคือ:**
1. ❌ Table name ผิด: `"CUSTOMER_SESSION"` → ควรเป็น `customer_sessions` (lowercase, snake_case)
2. ❌ Column names ผิด: `"Session_Id"`, `"Customer_Username"` → ควรเป็น `id`, `customer_username` (lowercase)
3. ❌ Column name ผิด: `"Status"` → ควรเป็น `status` (lowercase)
4. ❌ Logic ผิด: ใช้ `COUNT(Session_Id)` แต่ควรตรวจสอบจากข้อมูลที่ได้จาก Q3P.1 แทน
5. ❌ Parameter ผิด: ใช้ `${SessionId}` แต่ควรใช้ `session_id` จาก appointment ที่เลือก

**Logic ที่ใช้จริง:**
- ข้อมูล `used_sessions` และ `total_sessions` ได้มาจาก Q3P.1 แล้ว
- ตรวจสอบว่า `used_sessions < total_sessions` ใน use case layer (Go code)

#### 3. **Q3P.3 - ต้องแก้ไข**

**ใน Use Case Description คุณเขียน:**
```sql
UPDATE "CUSTOMER_SESSION"
   SET "Used_Sessions" = "Used_Sessions" + 1
   WHERE "Session_Id" = ${SessionId}
   AND "Customer_Username" = ${CustomerUsername}
   AND ("Total_Sessions" - "Used_Sessions") > 0;

INSERT INTO "CUSTOMER_LOG" (
"Customer_Username",
"Log_Type",
"Created_at"
) VALUES (
${CustomerUsername},
'CHECK_IN',
NOW()
);
```

**ปัญหาคือ:**
1. ❌ Table name ผิด: `"CUSTOMER_SESSION"` → ควรเป็น `customer_sessions`
2. ❌ Column names ผิด: ใช้ uppercase snake_case → ควรเป็น lowercase snake_case
3. ❌ Logic ผิด: INSERT log ใหม่ → ควร UPDATE log เดิมที่ PENDING เป็น CONFIRMED
4. ❌ ไม่มี WHERE clause สำหรับ UPDATE log

**Logic ที่ใช้จริง:**
- อัปเดต log เดิม (ที่มี status = 'PENDING') เป็น 'CONFIRMED'
- หัก session ผ่าน transaction

### ✅ SQL ที่ใช้จริง

**Q3P.2 - ตรวจสอบสิทธิ์ (ใน Use Case Layer):**
```go
// ตรวจสอบจากข้อมูลที่ได้จาก Q3P.1
if targetAppointment.UsedSessions >= targetAppointment.TotalSessions {
    return error("No remaining sessions available")
}
```

**Q3P.3 - หักสิทธิ์และบันทึก log:**
```sql
-- Step 1: อัปเดต log status จาก PENDING เป็น CONFIRMED
UPDATE customer_logs
SET status = 'CONFIRMED'
WHERE id = ? AND status = 'PENDING';

-- Step 2: หัก used_sessions (ผ่าน repository method)
-- IncrementUsedSessions(ctx, tx, sessionID)
-- ซึ่งทำ:
UPDATE customer_sessions
SET used_sessions = used_sessions + 1,
    updated_at = NOW()
WHERE id = ?
  AND status = 'ACTIVE'
  AND used_sessions < total_sessions;
```

### ✅ แก้ไข Use Case Description

**Step 5:**
```
5. ระบบดำเนินการตรวจสอบข้อมูลดังนี้

ตรวจสอบความถูกต้อง (Database):
- ตรวจสอบว่า appointment นี้เป็นของ trainer นี้จริงหรือไม่ (จากข้อมูล Q3P.1)
- ตรวจสอบว่ามี pending check-in log อยู่จริงหรือไม่ (checkinStatus = 'PENDING' และ checkinLogId > 0)

ตรวจสอบสิทธิ์คงเหลือในแพ็กเกจของลูกค้า โดยใช้ข้อมูลจาก Q3P.1:
- ตรวจสอบว่า used_sessions < total_sessions
- ถ้าไม่ผ่าน → แสดงข้อความ "No remaining sessions available"
```

**Step 6:**
```
6. เมื่อผ่านเงื่อนไข (5) ระบบดำเนินการหักสิทธิ์และบันทึก log ตาม Query Q3P.3

Q3P.3 (ใน Transaction):
-- อัปเดต check-in log status จาก PENDING เป็น CONFIRMED
UPDATE customer_logs
SET status = 'CONFIRMED'
WHERE id = ${checkin_log_id} AND status = 'PENDING';

-- หัก used_sessions
UPDATE customer_sessions
SET used_sessions = used_sessions + 1,
    updated_at = NOW()
WHERE id = ${session_id}
  AND status = 'ACTIVE'
  AND used_sessions < total_sessions;

หากไม่ผ่านเงื่อนไข (5) ระบบจะแสดงข้อความแจ้งเตือนว่า "No remaining sessions available" หรือ "Appointment not found"
```

---

## สรุปการแก้ไข

### Use Case 4C

1. ✅ Q4C.1 - Duration: ถูกต้องแล้ว
2. ❌ Q4C.2 - Session: **ต้องแก้** → สร้าง PENDING log เท่านั้น (ไม่หัก session)
3. ❌ Step 7: **ต้องแก้** → ข้อความสำหรับ Session ควรบอกว่า "รอ Trainer ยืนยัน"

### Use Case 3P

1. ✅ Q3P.1 - SELECT appointments: ถูกต้องแล้ว
2. ❌ Q3P.2 - **ต้องแก้** → ใช้ข้อมูลจาก Q3P.1 แทนการ query ใหม่ (ไม่ต้องมี COUNT query)
3. ❌ Q3P.3 - **ต้องแก้** → UPDATE log เดิม (PENDING → CONFIRMED) แทน INSERT ใหม่
4. ❌ Table/Column names: **ต้องแก้** → ใช้ lowercase snake_case (`customer_sessions`, `customer_logs`)

---

## คำแนะนำเพิ่มเติม

1. **ใช้ Transaction:** Q3P.3 ควรอยู่ใน transaction เดียวกันเพื่อความปลอดภัย
2. **Error Handling:** ต้องระบุ error cases ที่ชัดเจน (ไม่มี schedule, ไม่มี pending log, session หมด)
3. **Naming Convention:** ใช้ lowercase snake_case สำหรับ table และ column names ตาม MariaDB/MySQL convention

