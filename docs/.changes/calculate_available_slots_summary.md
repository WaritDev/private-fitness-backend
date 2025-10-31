# ✅ calculateAvailableSlots Implementation Summary

**วันที่:** October 31, 2025  
**Status:** ✅ Complete  
**Developer:** Development Team

---

## 📝 สรุปการพัฒนา

### ✅ ฟังก์ชันที่พัฒนาเสร็จแล้ว:

1. **`calculateAvailableSlots`** - คำนวณ Booking Slots ทุก 30 นาที
   - File: `/domain/usecases/booking_use_case.go`
   - Lines: ~76-236
   - Status: ✅ Complete

2. **`getDayOfWeekString`** - Helper function แปลง Weekday → String
   - File: `/domain/usecases/booking_use_case.go`
   - Lines: ~238-252
   - Status: ✅ Complete

3. **`filterCustomerBookings`** - กรองนัดของ customer
   - File: `/domain/usecases/booking_use_case.go`
   - Lines: ~99-130
   - Status: ✅ Enhanced with comments

---

## 🎯 Business Rules ที่ Implement แล้ว

### ✅ 1. นำเวลาทำงานจาก Q2C.2 มาสร้าง slots ทุก 30 นาที
```go
// สร้าง slots ทุก 30 นาที
for slotStart := workStartTime; slotStart.Before(workEndTime); slotStart = slotStart.Add(30 * time.Minute) {
    sessionEnd := slotStart.Add(120 * time.Minute) // 2 hours
    
    // ถ้า session เกินเวลาทำงาน ให้ข้าม
    if sessionEnd.After(workEndTime) {
        continue
    }
    // ...
}
```

**Result:** ✅ สร้าง slots ระยะห่าง 30 นาที จาก availability

---

### ✅ 2. ลบ slots ที่ตรงกับ dayOffs (Q2C.3)
```go
// ตรวจสอบว่า slot นี้ตรงกับวันหยุดหรือไม่
isDayOff := false
checkTime := slotStart
for checkTime.Before(sessionEnd) {
    key := checkTime.Format("2006-01-02 15:04")
    if dayOffMap[key] {
        isDayOff = true
        break
    }
    checkTime = checkTime.Add(time.Minute)
}

if isDayOff {
    continue // ข้าม slot ที่ตรงกับวันหยุด
}
```

**Result:** ✅ ข้าม slots ที่ตรงกับวันหยุด

---

### ✅ 3. ทำเครื่องหมาย slots ที่ถูกจองแล้ว
```go
// ตรวจสอบว่า slot นี้ถูกจองแล้วหรือไม่
isBooked := false
bookedBy := ""

checkTime = slotStart
for checkTime.Before(sessionEnd) {
    key := checkTime.Format("2006-01-02 15:04")
    if appt, exists := appointmentMap[key]; exists {
        isBooked = true
        bookedBy = appt.CustomerUsername
        break
    }
    checkTime = checkTime.Add(time.Minute)
}
```

**Result:** ✅ ทำเครื่องหมาย slots ที่จองแล้ว

---

### ✅ 4. แสดงเฉพาะ slots ของ customer นี้เอง
```go
// กำหนด slot type และ available status
slotType := "available"
available := true

if isBooked {
    available = false
    if customerUsername != "" && bookedBy == customerUsername {
        slotType = "booked"      // จองโดยตัวเอง - แสดงสีเทา "Booked"
    } else {
        slotType = "unavailable" // จองโดยคนอื่น - ไม่แสดงรายละเอียด
        bookedBy = ""            // ซ่อนชื่อคนอื่น
    }
}
```

**Result:** ✅ แสดงเฉพาะนัดของตัวเอง พร้อมข้อความ "Booked"

---

## 🏗️ Architecture & Design

### Data Structures:

#### 1. dayOffMap (O(1) lookup)
```go
dayOffMap := make(map[string]bool)
// Key: "2025-11-06 00:00" (YYYY-MM-DD HH:MM)
// Value: true (is day off)
```

#### 2. appointmentMap (O(1) lookup)
```go
appointmentMap := make(map[string]*repositories.AppointmentInfo)
// Key: "2025-11-01 09:00" (YYYY-MM-DD HH:MM)
// Value: pointer to appointment info
```

#### 3. availabilityByDay (Fast day lookup)
```go
availabilityByDay := make(map[string][]repositories.TrainerAvailabilityInfo)
// Key: "MONDAY", "TUESDAY", etc.
// Value: array of availability info
```

### Algorithm Complexity:

- **Time:** O(D × A × S × 120)
  - D = number of days in calendar range
  - A = availability slots per day
  - S = number of 30-min slots per availability
  - 120 = minutes to check per 2-hour session

- **Space:** O(M + N + R)
  - M = total minutes in all day-offs
  - N = total minutes in all appointments × 120
  - R = total result slots

### Performance Optimizations:

1. ✅ **Hash Maps** - O(1) lookup แทน O(n) loop
2. ✅ **Early Exit** - ข้าม slot ทันทีเมื่อเจอวันหยุด
3. ✅ **Pre-grouping** - จัดกลุ่ม availability ตามวันล่วงหน้า
4. ✅ **Pointer Usage** - ใช้ pointer ใน map เพื่อลด memory copy

---

## 📊 Output Format

### BookingSlot Structure:
```go
type BookingSlot struct {
    ID        int32     `json:"id"`        // Auto-increment ID
    StartTime time.Time `json:"startTime"` // Slot start time
    EndTime   time.Time `json:"endTime"`   // Slot end time (StartTime + 2 hours)
    Available bool      `json:"available"` // true = ว่าง, false = ถูกจอง
    IsBooked  bool      `json:"isBooked"`  // true = จองโดยตัวเอง
    BookedBy  string    `json:"bookedBy"`  // Customer username (ถ้าจองโดยตัวเอง)
    SlotType  string    `json:"slotType"`  // "available", "booked", "unavailable"
}
```

### SlotType Values:

| SlotType | Available | IsBooked | BookedBy | UI Display |
|----------|-----------|----------|----------|------------|
| `available` | `true` | `false` | `""` | 🟢 สีเขียว "Available" |
| `booked` | `false` | `true` | `"cust01"` | 🔘 สีเทา "Booked" |
| `unavailable` | `false` | `false` | `""` | 🔴 สีแดง "Unavailable" |

---

## 🧪 Testing Examples

### Example 1: Basic Available Slots
```json
Input:
- availability: MONDAY 09:00-17:00
- dayOffs: []
- appointments: []
- calendar: 2025-11-03 (Monday)

Output:
[
  {
    "id": 1,
    "startTime": "2025-11-03T09:00:00Z",
    "endTime": "2025-11-03T11:00:00Z",
    "available": true,
    "isBooked": false,
    "bookedBy": "",
    "slotType": "available"
  },
  // ... 09:30, 10:00, ..., 15:00
]
```

### Example 2: With Day Off
```json
Input:
- dayOffs: [{start: "2025-11-03T12:00:00Z", end: "2025-11-03T23:59:59Z"}]

Output:
- Slots 09:00-11:30: "available"
- Slots >= 12:00: ไม่แสดง (ข้าม)
```

### Example 3: Customer's Own Booking
```json
Input:
- appointments: [
    {id: 1, start: "2025-11-03T09:00:00Z", end: "2025-11-03T11:00:00Z", customer: "cust01"}
  ]
- customerUsername: "cust01"

Output:
[
  {
    "id": 1,
    "startTime": "2025-11-03T09:00:00Z",
    "endTime": "2025-11-03T11:00:00Z",
    "available": false,
    "isBooked": true,
    "bookedBy": "cust01",
    "slotType": "booked"  // แสดงสีเทา + "Booked"
  }
]
```

### Example 4: Others' Bookings
```json
Input:
- appointments: [
    {id: 2, start: "2025-11-03T10:00:00Z", end: "2025-11-03T12:00:00Z", customer: "cust02"}
  ]
- customerUsername: "cust01"

Output:
[
  {
    "id": 2,
    "startTime": "2025-11-03T10:00:00Z",
    "endTime": "2025-11-03T12:00:00Z",
    "available": false,
    "isBooked": false,
    "bookedBy": "",  // ซ่อนชื่อคนอื่น
    "slotType": "unavailable"
  }
]
```

---

## 📝 Frontend Integration Guide

### Step 1: Fetch Slots
```javascript
const response = await fetch(
  `http://localhost:8000/api/bookings/slots?trainerUsername=trainer1&customerUsername=cust01&calendarStart=2025-11-01T00:00:00Z&calendarEnd=2025-11-30T23:59:59Z`
);

const data = await response.json();
const slots = data.result.availableSlots;
```

### Step 2: Render Slots by Type
```javascript
slots.forEach(slot => {
  switch(slot.slotType) {
    case 'available':
      renderSlot(slot, 'green', 'Available', true);
      break;
    case 'booked':
      renderSlot(slot, 'gray', 'Booked', false); // แสดงสีเทา
      break;
    case 'unavailable':
      renderSlot(slot, 'red', 'Unavailable', false);
      break;
  }
});
```

### Step 3: Handle Click
```javascript
function handleSlotClick(slot) {
  if (slot.slotType === 'available') {
    bookAppointment(slot.startTime, slot.endTime);
  } else if (slot.slotType === 'booked') {
    showCancelDialog(slot);
  } else {
    showMessage('This slot is not available');
  }
}
```

---

## 📚 Documentation Files

1. **Main Implementation:** `/domain/usecases/booking_use_case.go`
2. **Logic Documentation:** `/docs/CALCULATE_AVAILABLE_SLOTS_LOGIC.md`
3. **API Guide:** `/docs/BOOKING_FLOW_FRONTEND_GUIDE.md`
4. **This Summary:** `/docs/.changes/calculate_available_slots_summary.md`

---

## ✅ Completion Checklist

- [x] Implement calculateAvailableSlots function
- [x] Implement getDayOfWeekString helper
- [x] Update filterCustomerBookings with comments
- [x] Create dayOffMap for fast lookup
- [x] Create appointmentMap for fast lookup
- [x] Create availabilityByDay grouping
- [x] Handle edge case: session exceeds work hours
- [x] Handle edge case: day off overlap
- [x] Handle edge case: appointment overlap
- [x] Implement 30-minute interval logic
- [x] Implement 2-hour session duration
- [x] Implement customer booking filter
- [x] Hide other customers' names
- [x] Show only own bookings as "Booked"
- [x] Write comprehensive documentation
- [x] Create algorithm flow diagram
- [x] Write testing examples
- [x] Write frontend integration guide

---

## 🚀 Next Steps

### For Backend:
```bash
# 1. Build project
cd /Users/pleng/cs-ku/year-3/sa/private-fitness-backend
go build -o tmp/bin/server cmd/app/main.go

# 2. Run tests
go test ./domain/usecases/... -v

# 3. Start server
./tmp/bin/server
```

### For Frontend:
1. เรียก API `/api/bookings/slots` 
2. ใช้ `availableSlots` array จาก response
3. Render slots ตาม `slotType`:
   - `available` → สีเขียว, คลิกได้
   - `booked` → สีเทา, แสดง "Booked"
   - `unavailable` → สีแดง, คลิกไม่ได้

### Testing:
```bash
# Test API
curl 'http://localhost:8000/api/bookings/slots?trainerUsername=trainer1&customerUsername=cust01&calendarStart=2025-11-01T00:00:00Z&calendarEnd=2025-11-30T23:59:59Z'

# Expected: JSON with availableSlots array
```

---

**Status:** ✅ Implementation Complete  
**Ready for:** Testing & Frontend Integration
