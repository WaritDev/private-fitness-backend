# Use Case 1P: Manage Working Hours - Frontend Implementation Guide

## Overview
Personal Trainer สามารถจัดการเวลาทำงานประจำสัปดาห์ได้ 3 รูปแบบ:
1. **Normal Case**: เพิ่มเวลาทำงานใหม่
2. **Alternative Case A**: แก้ไขเวลาทำงาน
3. **Alternative Case B**: ลบเวลาทำงาน

---

## 🔐 Authentication Requirements
**ทุก API ต้องมี JWT Token ใน Cookie หรือ Authorization Header**
```javascript
// JWT Token จะมี username ของ Trainer ที่ Login อยู่
headers: {
  'Authorization': 'Bearer <JWT_TOKEN>'
}
```

---

## 📋 Normal Case: เพิ่มเวลาทำงาน

### Flow Diagram
```
[Trainer] → Click "Working Hours" Menu
    ↓
[GET /api/trainers/working-hours] → Q1P.1: ดึงข้อมูลเวลาทำงานทั้งหมด
    ↓
[Display] → แสดงตารางเวลาทำงานปัจจุบัน
    ↓
[Trainer] → Click "Add Working Time" Button
    ↓
[Form] → กรอก Day_of_Week, Start_Time, End_Time
    ↓
[POST /api/trainers/working-hours] → Q1P.2 + Q1P.3: Validate + Insert
    ↓
[Success] → Show popup 4 วินาที + Redirect + Refresh ด้วย Q1P.1
```

### Step-by-Step Implementation

#### **Step 1: Load Working Hours Page**
```javascript
// GET /api/trainers/working-hours
async function loadWorkingHours() {
  try {
    const response = await fetch('/api/trainers/working-hours', {
      method: 'GET',
      credentials: 'include', // ส่ง cookie JWT token
      headers: {
        'Content-Type': 'application/json'
      }
    });

    const data = await response.json();
    
    if (data.status === 'success') {
      // Q1P.1 Response
      const workingHours = data.workingHours;
      /*
      workingHours = [
        {
          availabilityId: 1,
          dayOfWeek: "MONDAY",
          startTime: "09:00",
          endTime: "12:00"
        },
        {
          availabilityId: 2,
          dayOfWeek: "MONDAY",
          startTime: "14:00",
          endTime: "18:00"
        },
        ...
      ]
      */
      
      // Render ตารางเวลาทำงาน
      renderWorkingHoursTable(workingHours);
    }
  } catch (error) {
    console.error('Failed to load working hours:', error);
    showErrorMessage('Failed to load working hours');
  }
}
```

#### **Step 2: Render Table**
```javascript
function renderWorkingHoursTable(workingHours) {
  // จัดกลุ่มตาม Day_of_Week
  const groupedByDay = {
    'MONDAY': [],
    'TUESDAY': [],
    'WEDNESDAY': [],
    'THURSDAY': [],
    'FRIDAY': [],
    'SATURDAY': [],
    'SUNDAY': []
  };

  workingHours.forEach(item => {
    groupedByDay[item.dayOfWeek].push(item);
  });

  // สร้างตาราง HTML
  let tableHTML = `
    <table class="table">
      <thead>
        <tr>
          <th>Day</th>
          <th>Start Time</th>
          <th>End Time</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
  `;

  for (const [day, slots] of Object.entries(groupedByDay)) {
    if (slots.length === 0) {
      // ไม่มีเวลาทำงานในวันนี้
      tableHTML += `
        <tr>
          <td>${formatDayName(day)}</td>
          <td colspan="3" class="text-muted">No working hours</td>
        </tr>
      `;
    } else {
      slots.forEach((slot, index) => {
        tableHTML += `
          <tr>
            ${index === 0 ? `<td rowspan="${slots.length}">${formatDayName(day)}</td>` : ''}
            <td>${slot.startTime}</td>
            <td>${slot.endTime}</td>
            <td>
              <button class="btn btn-sm btn-primary" onclick="editWorkingTime(${slot.availabilityId})">Edit</button>
              <button class="btn btn-sm btn-danger" onclick="deleteWorkingTime(${slot.availabilityId})">Delete</button>
            </td>
          </tr>
        `;
      });
    }
  }

  tableHTML += `
      </tbody>
    </table>
    <button class="btn btn-success" onclick="showAddWorkingTimeForm()">Add Working Time</button>
  `;

  document.getElementById('working-hours-container').innerHTML = tableHTML;
}

function formatDayName(day) {
  const dayNames = {
    'MONDAY': 'จันทร์',
    'TUESDAY': 'อังคาร',
    'WEDNESDAY': 'พุธ',
    'THURSDAY': 'พฤหัสบดี',
    'FRIDAY': 'ศุกร์',
    'SATURDAY': 'เสาร์',
    'SUNDAY': 'อาทิตย์'
  };
  return dayNames[day] || day;
}
```

#### **Step 3: Show Add Form**
```javascript
function showAddWorkingTimeForm() {
  const formHTML = `
    <h3>Add Working Time</h3>
    <form id="add-working-time-form">
      <div class="form-group">
        <label for="dayOfWeek">Day of Week:</label>
        <select id="dayOfWeek" name="dayOfWeek" class="form-control" required>
          <option value="">-- Select Day --</option>
          <option value="MONDAY">Monday (จันทร์)</option>
          <option value="TUESDAY">Tuesday (อังคาร)</option>
          <option value="WEDNESDAY">Wednesday (พุธ)</option>
          <option value="THURSDAY">Thursday (พฤหัสบดี)</option>
          <option value="FRIDAY">Friday (ศุกร์)</option>
          <option value="SATURDAY">Saturday (เสาร์)</option>
          <option value="SUNDAY">Sunday (อาทิตย์)</option>
        </select>
      </div>

      <div class="form-group">
        <label for="startTime">Start Time:</label>
        <input type="time" id="startTime" name="startTime" class="form-control" required />
      </div>

      <div class="form-group">
        <label for="endTime">End Time:</label>
        <input type="time" id="endTime" name="endTime" class="form-control" required />
      </div>

      <div id="error-message" class="text-danger"></div>

      <button type="submit" class="btn btn-success">Save</button>
      <button type="button" class="btn btn-secondary" onclick="cancelAddWorkingTime()">Cancel</button>
    </form>
  `;

  document.getElementById('working-hours-container').innerHTML = formHTML;

  // Attach submit handler
  document.getElementById('add-working-time-form').addEventListener('submit', handleAddWorkingTime);
}

function cancelAddWorkingTime() {
  // กลับไปหน้าตาราง
  loadWorkingHours();
}
```

#### **Step 4: Submit Add Working Time**
```javascript
async function handleAddWorkingTime(event) {
  event.preventDefault();

  // ดึงค่าจากฟอร์ม
  const dayOfWeek = document.getElementById('dayOfWeek').value;
  const startTime = document.getElementById('startTime').value; // "09:00"
  const endTime = document.getElementById('endTime').value;     // "12:00"

  // Validation ฝั่ง Client
  if (!dayOfWeek || !startTime || !endTime) {
    showFormError('All fields are required');
    return;
  }

  if (endTime <= startTime) {
    showFormError('End time must be after start time');
    return;
  }

  try {
    // Q1P.2 + Q1P.3: POST /api/trainers/working-hours
    const response = await fetch('/api/trainers/working-hours', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        dayOfWeek: dayOfWeek,
        startTime: startTime,
        endTime: endTime
      })
    });

    const data = await response.json();

    if (data.status === 'success') {
      // Step 9: แสดง popup 4 วินาที
      showSuccessPopup('Working time added successfully');

      // Redirect + Refresh ด้วย Q1P.1
      setTimeout(() => {
        loadWorkingHours();
      }, 4000);

    } else if (data.status === 'error') {
      // Validation error จาก backend (Q1P.2)
      showFormError(data.message);
    }
  } catch (error) {
    console.error('Failed to add working time:', error);
    showFormError('Failed to add working time. Please try again.');
  }
}

function showFormError(message) {
  document.getElementById('error-message').textContent = message;
}

function showSuccessPopup(message) {
  // แสดง popup หรือ toast notification
  const popup = document.createElement('div');
  popup.className = 'alert alert-success popup';
  popup.textContent = message;
  popup.style.position = 'fixed';
  popup.style.top = '20px';
  popup.style.right = '20px';
  popup.style.zIndex = '9999';
  
  document.body.appendChild(popup);

  // ลบหลัง 4 วินาที
  setTimeout(() => {
    popup.remove();
  }, 4000);
}
```

---

## 📝 Alternative Case A: แก้ไขเวลาทำงาน

### Flow Diagram
```
[Trainer] → Click "Working Hours" Menu
    ↓
[GET /api/trainers/working-hours] → Q1P.A1: ดึงข้อมูลเวลาทำงานทั้งหมด
    ↓
[Display] → แสดงตารางเวลาทำงาน
    ↓
[Trainer] → Click "Edit" Button บนแถวที่ต้องการ
    ↓
[GET /api/trainers/working-hours] → Q1P.A2: ดึงข้อมูลเดิมมาแสดงในฟอร์ม (Optional: ใช้ข้อมูลจากตารางได้)
    ↓
[Form] → แก้ไข Day_of_Week, Start_Time, End_Time
    ↓
[PUT /api/trainers/working-hours/:id] → Q1P.A3 + Q1P.A4: Validate + Update
    ↓
[Success] → Show popup 4 วินาที + Redirect + Refresh ด้วย Q1P.1
```

### Implementation

#### **Step 1: Edit Button Click**
```javascript
function editWorkingTime(availabilityId) {
  // ดึงข้อมูลจาก cache (จากตารางที่แสดงอยู่)
  const workingHoursCache = getWorkingHoursCache(); // จาก Step Q1P.1
  const workingTime = workingHoursCache.find(item => item.availabilityId === availabilityId);

  if (!workingTime) {
    showErrorMessage('Working time not found');
    return;
  }

  // แสดงฟอร์มแก้ไข (Pre-fill ข้อมูลเดิม)
  showEditWorkingTimeForm(workingTime);
}

function showEditWorkingTimeForm(workingTime) {
  const formHTML = `
    <h3>Edit Working Time</h3>
    <form id="edit-working-time-form" data-id="${workingTime.availabilityId}">
      <div class="form-group">
        <label for="editDayOfWeek">Day of Week:</label>
        <select id="editDayOfWeek" name="dayOfWeek" class="form-control" required>
          <option value="MONDAY" ${workingTime.dayOfWeek === 'MONDAY' ? 'selected' : ''}>Monday (จันทร์)</option>
          <option value="TUESDAY" ${workingTime.dayOfWeek === 'TUESDAY' ? 'selected' : ''}>Tuesday (อังคาร)</option>
          <option value="WEDNESDAY" ${workingTime.dayOfWeek === 'WEDNESDAY' ? 'selected' : ''}>Wednesday (พุธ)</option>
          <option value="THURSDAY" ${workingTime.dayOfWeek === 'THURSDAY' ? 'selected' : ''}>Thursday (พฤหัสบดี)</option>
          <option value="FRIDAY" ${workingTime.dayOfWeek === 'FRIDAY' ? 'selected' : ''}>Friday (ศุกร์)</option>
          <option value="SATURDAY" ${workingTime.dayOfWeek === 'SATURDAY' ? 'selected' : ''}>Saturday (เสาร์)</option>
          <option value="SUNDAY" ${workingTime.dayOfWeek === 'SUNDAY' ? 'selected' : ''}>Sunday (อาทิตย์)</option>
        </select>
      </div>

      <div class="form-group">
        <label for="editStartTime">Start Time:</label>
        <input type="time" id="editStartTime" name="startTime" class="form-control" value="${workingTime.startTime}" required />
      </div>

      <div class="form-group">
        <label for="editEndTime">End Time:</label>
        <input type="time" id="editEndTime" name="endTime" class="form-control" value="${workingTime.endTime}" required />
      </div>

      <div id="edit-error-message" class="text-danger"></div>

      <button type="submit" class="btn btn-success">Save</button>
      <button type="button" class="btn btn-secondary" onclick="loadWorkingHours()">Cancel</button>
    </form>
  `;

  document.getElementById('working-hours-container').innerHTML = formHTML;

  // Attach submit handler
  document.getElementById('edit-working-time-form').addEventListener('submit', handleUpdateWorkingTime);
}
```

#### **Step 2: Submit Update**
```javascript
async function handleUpdateWorkingTime(event) {
  event.preventDefault();

  const form = event.target;
  const availabilityId = form.getAttribute('data-id');

  // ดึงค่าจากฟอร์ม
  const dayOfWeek = document.getElementById('editDayOfWeek').value;
  const startTime = document.getElementById('editStartTime').value;
  const endTime = document.getElementById('editEndTime').value;

  // Validation ฝั่ง Client
  if (!dayOfWeek || !startTime || !endTime) {
    showEditFormError('All fields are required');
    return;
  }

  if (endTime <= startTime) {
    showEditFormError('End time must be after start time');
    return;
  }

  try {
    // Q1P.A3 + Q1P.A4: PUT /api/trainers/working-hours/:id
    const response = await fetch(`/api/trainers/working-hours/${availabilityId}`, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        dayOfWeek: dayOfWeek,
        startTime: startTime,
        endTime: endTime
      })
    });

    const data = await response.json();

    if (data.status === 'success') {
      // Step 9: แสดง popup 4 วินาที
      showSuccessPopup('Working time updated successfully');

      // Redirect + Refresh ด้วย Q1P.1
      setTimeout(() => {
        loadWorkingHours();
      }, 4000);

    } else if (data.status === 'error') {
      // Validation error (Q1P.A3)
      showEditFormError(data.message);
    }
  } catch (error) {
    console.error('Failed to update working time:', error);
    showEditFormError('Failed to update working time. Please try again.');
  }
}

function showEditFormError(message) {
  document.getElementById('edit-error-message').textContent = message;
}
```

---

## 🗑️ Alternative Case B: ลบเวลาทำงาน

### Flow Diagram
```
[Trainer] → Click "Working Hours" Menu
    ↓
[GET /api/trainers/working-hours] → Q1P.B1: ดึงข้อมูลเวลาทำงานทั้งหมด
    ↓
[Display] → แสดงตารางเวลาทำงาน
    ↓
[Trainer] → Click "Delete" Button บนแถวที่ต้องการลบ
    ↓
[Popup] → "Confirm deletion of this working time?"
    ↓
[Trainer] → Click "Confirm"
    ↓
[DELETE /api/trainers/working-hours/:id] → Q1P.B2: ลบข้อมูล
    ↓
[Success] → Show popup 4 วินาที + Refresh ด้วย Q1P.1
```

### Implementation

```javascript
function deleteWorkingTime(availabilityId) {
  // Step 5: แสดง confirmation popup
  const confirmed = confirm('Confirm deletion of this working time?');

  if (!confirmed) {
    return; // ยกเลิก
  }

  // Step 7: เรียก API ลบข้อมูล
  performDeleteWorkingTime(availabilityId);
}

async function performDeleteWorkingTime(availabilityId) {
  try {
    // Q1P.B2: DELETE /api/trainers/working-hours/:id
    const response = await fetch(`/api/trainers/working-hours/${availabilityId}`, {
      method: 'DELETE',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      }
    });

    const data = await response.json();

    if (data.status === 'success') {
      // Step 8: แสดง popup 4 วินาที
      showSuccessPopup('Working time deleted successfully');

      // Refresh หน้า Working Hours ด้วย Q1P.1
      setTimeout(() => {
        loadWorkingHours();
      }, 4000);

    } else if (data.status === 'error') {
      showErrorMessage(data.message);
    }
  } catch (error) {
    console.error('Failed to delete working time:', error);
    showErrorMessage('Failed to delete working time. Please try again.');
  }
}
```

---

## 🔧 API Reference

### 1. GET /api/trainers/working-hours
**Query Q1P.1: ดึงข้อมูลเวลาทำงานทั้งหมดของเทรนเนอร์**

**Request:**
```http
GET /api/trainers/working-hours HTTP/1.1
Cookie: jwt_token=<TOKEN>
```

**Response (Success):**
```json
{
  "status": "success",
  "message": "Working hours retrieved successfully",
  "workingHours": [
    {
      "availabilityId": 1,
      "dayOfWeek": "MONDAY",
      "startTime": "09:00",
      "endTime": "12:00"
    },
    {
      "availabilityId": 2,
      "dayOfWeek": "MONDAY",
      "startTime": "14:00",
      "endTime": "18:00"
    },
    {
      "availabilityId": 3,
      "dayOfWeek": "TUESDAY",
      "startTime": "10:00",
      "endTime": "16:00"
    }
  ]
}
```

---

### 2. POST /api/trainers/working-hours
**Query Q1P.2 + Q1P.3: ตรวจสอบทับซ้อน + เพิ่มเวลาทำงานใหม่**

**Request:**
```http
POST /api/trainers/working-hours HTTP/1.1
Cookie: jwt_token=<TOKEN>
Content-Type: application/json

{
  "dayOfWeek": "MONDAY",
  "startTime": "09:00",
  "endTime": "12:00"
}
```

**Response (Success):**
```json
{
  "status": "success",
  "message": "Working time added successfully"
}
```

**Response (Validation Error - Overlap):**
```json
{
  "status": "error",
  "message": "Working time overlaps with existing schedule"
}
```

**Response (Validation Error - Time Range):**
```json
{
  "status": "error",
  "message": "End time must be after start time"
}
```

**Response (Validation Error - Format):**
```json
{
  "status": "error",
  "message": "Invalid start time format. Expected HH:MM"
}
```

---

### 3. PUT /api/trainers/working-hours/:id
**Query Q1P.A3 + Q1P.A4: ตรวจสอบทับซ้อน + อัปเดตเวลาทำงาน**

**Request:**
```http
PUT /api/trainers/working-hours/5 HTTP/1.1
Cookie: jwt_token=<TOKEN>
Content-Type: application/json

{
  "dayOfWeek": "TUESDAY",
  "startTime": "10:00",
  "endTime": "14:00"
}
```

**Response (Success):**
```json
{
  "status": "success",
  "message": "Working time updated successfully"
}
```

**Response (Not Found):**
```json
{
  "status": "error",
  "message": "Working hour not found or does not belong to you"
}
```

**Response (Validation Error):**
```json
{
  "status": "error",
  "message": "End time must be after start time"
}
```

---

### 4. DELETE /api/trainers/working-hours/:id
**Query Q1P.B2: ลบเวลาทำงาน**

**Request:**
```http
DELETE /api/trainers/working-hours/5 HTTP/1.1
Cookie: jwt_token=<TOKEN>
```

**Response (Success):**
```json
{
  "status": "success",
  "message": "Working time deleted successfully"
}
```

**Response (Not Found):**
```json
{
  "status": "error",
  "message": "Working hour not found or does not belong to you"
}
```

---

## ⚠️ Error Handling

### Client-Side Validation
```javascript
function validateWorkingTimeForm(dayOfWeek, startTime, endTime) {
  const errors = [];

  // 1. ตรวจสอบค่าว่าง
  if (!dayOfWeek) errors.push('Day of week is required');
  if (!startTime) errors.push('Start time is required');
  if (!endTime) errors.push('End time is required');

  // 2. ตรวจสอบ Day_of_Week
  const validDays = ['MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY'];
  if (dayOfWeek && !validDays.includes(dayOfWeek)) {
    errors.push('Invalid day of week');
  }

  // 3. ตรวจสอบรูปแบบเวลา (HH:MM)
  const timeRegex = /^([01]\d|2[0-3]):([0-5]\d)$/;
  if (startTime && !timeRegex.test(startTime)) {
    errors.push('Invalid start time format. Expected HH:MM');
  }
  if (endTime && !timeRegex.test(endTime)) {
    errors.push('Invalid end time format. Expected HH:MM');
  }

  // 4. ตรวจสอบ End_Time > Start_Time
  if (startTime && endTime && endTime <= startTime) {
    errors.push('End time must be after start time');
  }

  return errors;
}
```

### Server-Side Validation
Backend จะตรวจสอบเพิ่มเติม:
- ✅ Q1P.2: ช่วงเวลาทำงานห้ามซ้อนทับกับช่วงเวลาที่มีอยู่แล้ว
- ✅ Ownership: เฉพาะ Trainer ที่ Login เท่านั้นที่แก้ไข/ลบเวลาทำงานของตนเอง

---

## 📊 Data Dictionary

### Working Hours Fields

| Field Name | Type | Format | Required | Description |
|------------|------|--------|----------|-------------|
| `availabilityId` | Integer | - | Yes (Response) | รหัสเวลาทำงาน (Primary Key) |
| `dayOfWeek` | String | ENUM | Yes | วันในสัปดาห์: MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, SUNDAY |
| `startTime` | String | HH:MM | Yes | เวลาเริ่มต้น เช่น "09:00" |
| `endTime` | String | HH:MM | Yes | เวลาสิ้นสุด เช่น "18:00" |

### Validation Rules

| Rule | Description |
|------|-------------|
| `dayOfWeek` | ห้ามค่าว่าง, ต้องเป็น MONDAY/TUESDAY/WEDNESDAY/THURSDAY/FRIDAY/SATURDAY/SUNDAY |
| `startTime` | ห้ามค่าว่าง, ต้องเป็นรูปแบบ HH:MM (00:00 - 23:59) |
| `endTime` | ห้ามค่าว่าง, ต้องเป็นรูปแบบ HH:MM (00:00 - 23:59) |
| `endTime > startTime` | เวลาสิ้นสุดต้องอยู่หลังเวลาเริ่มต้น |
| `No Overlap` | ช่วงเวลาทำงานห้ามซ้อนทับกับช่วงเวลาที่มีอยู่แล้วในวันเดียวกัน |

---

## 🎨 UI/UX Recommendations

### 1. Working Hours Table
- แสดงข้อมูลเป็นตาราง โดยจัดกลุ่มตาม Day_of_Week
- เรียงลำดับวัน: Monday → Sunday
- เรียงเวลาใน 1 วัน: Start_Time จากน้อยไปมาก
- แสดงปุ่ม "Edit" และ "Delete" ในแต่ละแถว

### 2. Add/Edit Form
- ใช้ Dropdown สำหรับ Day_of_Week
- ใช้ `<input type="time">` สำหรับ Start_Time และ End_Time
- แสดง error message สีแดงใต้ฟอร์ม
- ปุ่ม "Save" และ "Cancel"

### 3. Delete Confirmation
- ใช้ Modal popup หรือ `confirm()` dialog
- ข้อความ: "Confirm deletion of this working time?"
- ปุ่ม "Confirm" และ "Cancel"

### 4. Success Popup
- แสดง Toast notification มุมขวาบน
- Auto-hide หลัง 4 วินาที
- สีเขียว สำหรับ success
- ข้อความ: "Working time added/updated/deleted successfully"

---

## 🧪 Testing Checklist

### Normal Case (Add)
- [ ] แสดงตารางเวลาทำงานได้ถูกต้อง
- [ ] คลิก "Add Working Time" แสดงฟอร์ม
- [ ] Validate ค่าว่าง (Client-side)
- [ ] Validate End_Time > Start_Time (Client-side)
- [ ] Validate ทับซ้อน (Server-side)
- [ ] บันทึกสำเร็จ → แสดง popup 4 วินาที
- [ ] Redirect + Refresh ตารางได้ถูกต้อง

### Alternative Case A (Edit)
- [ ] คลิก "Edit" แสดงฟอร์มพร้อมข้อมูลเดิม
- [ ] แก้ไขและบันทึกได้ถูกต้อง
- [ ] Validate เช่นเดียวกับ Add
- [ ] แสดง popup "updated successfully"

### Alternative Case B (Delete)
- [ ] คลิก "Delete" แสดง confirmation
- [ ] คลิก "Confirm" ลบข้อมูลได้
- [ ] แสดง popup "deleted successfully"
- [ ] Refresh ตารางอัตโนมัติ

---

## 📌 Notes

1. **JWT Token**: ต้องมี JWT token ใน cookie หรือ Authorization header ทุก request
2. **Trainer Username**: ดึงจาก JWT claims ใน backend (ไม่ต้องส่งใน request body)
3. **Time Format**: ใช้ HH:MM (24-hour format)
4. **Day Order**: Frontend ควรจัดเรียงตาม MONDAY → SUNDAY
5. **Overlap Check**: Backend จะตรวจสอบทับซ้อนให้ (ใช้ SQL OVERLAPS clause)
6. **Ownership**: Backend จะตรวจสอบว่า Trainer แก้ไข/ลบเฉพาะเวลาทำงานของตนเอง

---

## 🚀 Quick Start Example

```javascript
// ตัวอย่างการใช้งานทั้ง 3 Use Case

// 1. Load Working Hours (เมื่อเปิดหน้า)
window.addEventListener('DOMContentLoaded', () => {
  loadWorkingHours();
});

// 2. Add Working Time
document.getElementById('add-btn').addEventListener('click', () => {
  showAddWorkingTimeForm();
});

// 3. Edit Working Time
function editWorkingTime(id) {
  const workingTime = getWorkingTimeById(id);
  showEditWorkingTimeForm(workingTime);
}

// 4. Delete Working Time
function deleteWorkingTime(id) {
  if (confirm('Confirm deletion of this working time?')) {
    performDeleteWorkingTime(id);
  }
}
```

---

## ✅ Summary

API สำหรับ Use Case 1P พัฒนาครบทั้ง 4 endpoints:
1. ✅ GET `/api/trainers/working-hours` - ดึงข้อมูลเวลาทำงาน (Q1P.1)
2. ✅ POST `/api/trainers/working-hours` - เพิ่มเวลาทำงานใหม่ (Q1P.2 + Q1P.3)
3. ✅ PUT `/api/trainers/working-hours/:id` - แก้ไขเวลาทำงาน (Q1P.A3 + Q1P.A4)
4. ✅ DELETE `/api/trainers/working-hours/:id` - ลบเวลาทำงาน (Q1P.B2)

**Backend พร้อมใช้งานแล้ว!** Frontend สามารถเริ่มพัฒนาตามเอกสารนี้ได้เลยครับ 🎉
