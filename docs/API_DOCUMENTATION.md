# 🚀 Private Fitness API Documentation

> **สำหรับ Frontend Developer**  
> Base URL: `http://localhost:8000`  
> Updated: October 30, 2025

---

## 📋 Table of Contents

1. [Authentication APIs](#1-authentication-apis)
2. [Product APIs](#2-product-apis)
3. [User Validation APIs](#3-user-validation-apis)
4. [Payment APIs](#4-payment-apis)
5. [Customer Registration APIs](#5-customer-registration-apis)
6. [Booking APIs](#6-booking-apis)
7. [Response Format](#7-response-format)
8. [Error Codes](#8-error-codes)

---

## 1. Authentication APIs

### 1.1 Login

**Endpoint:** `POST /api/auth/login`

**Description:** เข้าสู่ระบบด้วย username และ password (Use Case 0S: เข้าสู่ระบบ)

**Business Logic:**
1. ตรวจสอบ username และ password (bcrypt hash)
2. ตรวจสอบสถานะบัญชี (is_active = true)
3. อัปเดต `updated_at` เพื่อ track last login time (Q0S.2)
4. สร้าง JWT token (7 days expiry)

**Request Body:**
```json
{
  "username": "cust01",
  "password": "Password123!"
}
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Login successful",
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "sub": "cust01",
      "role": "CUSTOMER",
      "firstName": "John",
      "lastName": "Doe"
    }
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "invalid credentials",
  "result": null
}
```

**Cookie Set:**
- Cookie name: `pf_auth`
- HTTPOnly: true
- SameSite: Lax
- Max-Age: 604800 (7 days)

**Usage Example:**
```javascript
const response = await fetch('http://localhost:8000/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include', // ต้องมีเพื่อรับ cookie
  body: JSON.stringify({
    username: 'cust01',
    password: 'Password123!'
  })
});

const data = await response.json();
if (data.status === 'success') {
  localStorage.setItem('token', data.result.token);
  localStorage.setItem('user', JSON.stringify(data.result.user));
}
```

---

### 1.2 Get Current User (Me)

**Endpoint:** `GET /api/auth/me`

**Description:** ดึงข้อมูลผู้ใช้ปัจจุบันจาก JWT token (รองรับทั้ง cookie และ Authorization header)

**Request Headers:**
```
Authorization: Bearer {token}
// หรือส่ง cookie pf_auth
```

**Success Response (200 OK - Authenticated):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "User retrieved successfully",
  "result": {
    "authenticated": true,
    "user": {
      "sub": "cust01",
      "role": "CUSTOMER",
      "firstName": "John",
      "lastName": "Doe"
    }
  }
}
```

**Response (200 OK - Not Authenticated):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "User not authenticated",
  "result": {
    "authenticated": false
  }
}
```

**Usage Example:**
```javascript
const response = await fetch('http://localhost:8000/api/auth/me', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`
  },
  credentials: 'include'
});

const data = await response.json();
if (data.result.authenticated) {
  console.log('User:', data.result.user);
}
```

---

### 1.3 Logout

**Endpoint:** `POST /api/auth/logout` หรือ `GET /api/auth/logout`

**Description:** ออกจากระบบ (ลบ cookie)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Logged out successfully",
  "result": {
    "ok": true
  }
}
```

**Cookie Cleared:**
- Cookie `pf_auth` จะถูกลบ (MaxAge=-1)

**Usage Example:**
```javascript
await fetch('http://localhost:8000/api/auth/logout', {
  method: 'POST',
  credentials: 'include'
});

localStorage.removeItem('token');
localStorage.removeItem('user');
```

---

## 2. Product APIs

### 2.1 List All Products

**Endpoint:** `GET /api/products`

**Description:** ดึงรายการสินค้า/แพ็กเกจทั้งหมดที่ active

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Products retrieved successfully",
  "result": [
    {
      "id": 1,
      "name": "Monthly Gym Access - Basic",
      "type": "DURATION",
      "category": "ECONOMIC",
      "listPrice": 1200.00,
      "durationDays": 30,
      "sessionAmount": null,
      "isActive": true,
      "paymentAccountId": 1,
      "createdAt": "2025-07-02T00:00:00Z",
      "updatedAt": "2025-10-30T00:00:00Z"
    },
    {
      "id": 10,
      "name": "Yoga Sessions - 5 Pack",
      "type": "SESSION",
      "category": "ECONOMIC",
      "listPrice": 1500.00,
      "durationDays": null,
      "sessionAmount": 5,
      "isActive": true,
      "paymentAccountId": 1,
      "createdAt": "2025-09-01T00:00:00Z",
      "updatedAt": "2025-10-30T00:00:00Z"
    }
  ]
}
```

**Field Descriptions:**
- `type`: `"DURATION"` (รายเดือน/รายปี) หรือ `"SESSION"` (แพ็กเกจครั้ง)
- `category`: `"ECONOMIC"`, `"BUSINESS"`, `"FIRST_CLASS"`
- `durationDays`: จำนวนวันที่ใช้ได้ (สำหรับ DURATION)
- `sessionAmount`: จำนวนครั้งที่ใช้ได้ (สำหรับ SESSION)

---

### 2.2 Get Product By ID

**Endpoint:** `GET /api/products/:id`

**Description:** ดึงข้อมูลสินค้า/แพ็กเกจตาม ID

**Path Parameters:**
- `id` (integer): Product ID

**Example:** `GET /api/products/10`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Product retrieved successfully",
  "result": {
    "id": 10,
    "name": "Yoga Sessions - 5 Pack",
    "type": "SESSION",
    "category": "ECONOMIC",
    "listPrice": 1500.00,
    "durationDays": null,
    "sessionAmount": 5,
    "isActive": true,
    "paymentAccountId": 1,
    "createdAt": "2025-09-01T00:00:00Z",
    "updatedAt": "2025-10-30T00:00:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "status": "error",
  "status_code": 404,
  "message": "product not found",
  "result": null
}
```

---

### 2.3 List Duration Products

**Endpoint:** `GET /api/products/durations`

**Description:** ดึงเฉพาะแพ็กเกจ **DURATION** (รายเดือน/รายปี)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Duration products retrieved successfully",
  "result": [
    {
      "id": 1,
      "name": "Monthly Gym Access - Basic",
      "type": "DURATION",
      "category": "ECONOMIC",
      "listPrice": 1200.00,
      "durationDays": 30,
      "sessionAmount": null,
      "isActive": true,
      "paymentAccountId": 1,
      "createdAt": "2025-07-02T00:00:00Z",
      "updatedAt": "2025-10-30T00:00:00Z"
    }
  ]
}
```

**Use Case:** ใช้แสดงตัวเลือกแพ็กเกจรายเดือน/รายปีในหน้าสมัครสมาชิก

---

### 2.4 List Session Products

**Endpoint:** `GET /api/products/sessions`

**Description:** ดึงเฉพาะแพ็กเกจ **SESSION** (แพ็กเกจครั้ง)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Session products retrieved successfully",
  "result": [
    {
      "id": 10,
      "name": "Yoga Sessions - 5 Pack",
      "type": "SESSION",
      "category": "ECONOMIC",
      "listPrice": 1500.00,
      "durationDays": null,
      "sessionAmount": 5,
      "isActive": true,
      "paymentAccountId": 1,
      "createdAt": "2025-09-01T00:00:00Z",
      "updatedAt": "2025-10-30T00:00:00Z"
    }
  ]
}
```

**Use Case:** ใช้แสดงตัวเลือกแพ็กเกจ Personal Training ในหน้าสมัครคอร์ส Sessions

---

## 3. User Validation APIs

### 3.1 Check Phone Number

**Endpoint:** `GET /api/users/check-phone`

**Description:** ตรวจสอบว่าเบอร์โทรซ้ำหรือไม่ (Use Case Q3S.1)

**Query Parameters:**
- `phone` (string): เบอร์โทรศัพท์ที่ต้องการตรวจสอบ

**Example:** `GET /api/users/check-phone?phone=0811111001`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Phone number check completed",
  "result": {
    "exists": true,
    "available": false
  }
}
```

**Field Descriptions:**
- `exists`: `true` = มีคนใช้เบอร์นี้แล้ว, `false` = ยังไม่มีใครใช้
- `available`: `true` = ใช้ได้, `false` = ซ้ำ (เป็น inverse ของ exists)

**Usage Example:**
```javascript
const phone = '0811111001';
const response = await fetch(`http://localhost:8000/api/users/check-phone?phone=${phone}`);
const data = await response.json();

if (data.result.exists) {
  alert('เบอร์โทรนี้ถูกใช้งานแล้ว');
}
```

---

### 3.2 Check Gmail

**Endpoint:** `GET /api/users/check-gmail`

**Description:** ตรวจสอบว่าอีเมลซ้ำหรือไม่ (Use Case Q3S.2)

**Query Parameters:**
- `gmail` (string): อีเมลที่ต้องการตรวจสอบ

**Example:** `GET /api/users/check-gmail?gmail=john.doe@gmail.com`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Gmail check completed",
  "result": {
    "exists": false,
    "available": true
  }
}
```

**Usage Example:**
```javascript
const gmail = 'john.doe@gmail.com';
const response = await fetch(`http://localhost:8000/api/users/check-gmail?gmail=${encodeURIComponent(gmail)}`);
const data = await response.json();

if (data.result.exists) {
  alert('อีเมลนี้ถูกใช้งานแล้ว');
}
```

---

## 4. Payment APIs

### 4.1 Get Payment Info

**Endpoint:** `GET /api/payments/info/:productId`

**Description:** ดึงข้อมูลชำระเงินสำหรับสินค้า/แพ็กเกจ (Use Case 5S: ยืนยันการชำระเงิน)

**Path Parameters:**
- `productId` (integer): Product ID

**Query Parameters:**
- `discount` (float, optional): จำนวนส่วนลด (default: 0)

**Example:** `GET /api/payments/info/11?discount=200.50`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Payment info retrieved successfully",
  "result": {
    "productId": 11,
    "productName": "Yoga Sessions - 10 Pack",
    "productType": "SESSION",
    "productCategory": "ECONOMIC",
    "listPrice": 2800.00,
    "discountAmount": 200.50,
    "payableAmount": 2599.50,
    "sessionAmount": 10,
    "durationDays": null,
    "paymentAccountId": 1,
    "accountName": "Private Fitness - Main Account",
    "accountNumber": "123-4-56789-0",
    "bankName": "Bangkok Bank",
    "qrCodeUrl": "https://example.com/qr/main.png",
    "accountActive": true
  }
}
```

**Field Descriptions:**
- `listPrice`: ราคาปกติ
- `discountAmount`: ส่วนลด
- `payableAmount`: ยอดที่ต้องชำระจริง (listPrice - discountAmount)
- `sessionAmount`: จำนวนครั้ง (สำหรับ SESSION)
- `durationDays`: จำนวนวัน (สำหรับ DURATION)
- `qrCodeUrl`: URL รูป QR Code สำหรับสแกนจ่ายเงิน

**Usage Example:**
```javascript
const productId = 11;
const discount = 200.50;
const response = await fetch(`http://localhost:8000/api/payments/info/${productId}?discount=${discount}`);
const data = await response.json();

// แสดงข้อมูลการชำระเงิน
console.log('ราคาปกติ:', data.result.listPrice);
console.log('ส่วนลด:', data.result.discountAmount);
console.log('ยอดชำระ:', data.result.payableAmount);
console.log('บัญชีธนาคาร:', data.result.accountNumber);
console.log('QR Code:', data.result.qrCodeUrl);
```

---

## 5. Customer Registration APIs

### 5.1 Register Customer Duration (รายเดือน/รายปี)

**Endpoint:** `POST /api/customers/durations/register`

**Description:** ลงทะเบียนลูกค้าใหม่สำหรับแพ็กเกจ Duration (Use Case 2.1C)

**Request Body:**
```json
{
  "username": "cust09",
  "password": "SecurePass123!",
  "confirmPassword": "SecurePass123!",
  "firstName": "สมชาย",
  "lastName": "ใจดี",
  "gender": "MALE",
  "dateOfBirth": "1995-01-15",
  "phone": "0899999999",
  "gmail": "somchai@gmail.com",
  "healthInfo": "ไม่มีโรคประจำตัว",
  "address": "123 ถนนสุขุมวิท กรุงเทพฯ 10110",
  "companyName": "บริษัท ABC จำกัด",
  "companyPosition": "Software Engineer",
  "maritalStatus": "SINGLE",
  "emergencyContactName": "สมหญิง ใจดี",
  "emergencyContactRelationship": "แม่",
  "emergencyContactPhone": "0888888888",
  "marketingSource": "Facebook Ads",
  "productId": 1,
  "salesUsername": "sales1",
  "startDate": "2025-11-01",
  "durationDays": 30,
  "pricePaid": 1200.00,
  "discountAmount": 0.00
}
```

**Required Fields:**
- `username`: ชื่อผู้ใช้ (4-30 ตัวอักษร, ไม่ซ้ำ)
- `password`: รหัสผ่าน (min 8 ตัวอักษร)
- `confirmPassword`: ยืนยันรหัสผ่าน (ต้องตรงกับ password)
- All other fields as shown above

**Field Validations:**
- `gender`: `"MALE"`, `"FEMALE"`, `"OTHER"`
- `maritalStatus`: `"SINGLE"`, `"MARRIED"`, `"DIVORCED"`, `"WIDOWED"`
- `dateOfBirth`: Format `YYYY-MM-DD`
- `startDate`: Format `YYYY-MM-DD`
- `gmail`: Must be valid email format

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Customer duration registered successfully",
  "result": {
    "username": "cust09",
    "durationId": 109,
    "productId": 1,
    "salesUsername": "sales1",
    "startDate": "2025-11-01T00:00:00Z",
    "endDate": "2025-12-01T00:00:00Z",
    "durationDays": 30,
    "pricePaid": "1200.00",
    "discountAmount": "0.00",
    "message": "Customer duration registered successfully"
  }
}
```

**Error Responses:**

**400 Bad Request - Username Exists:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "USERNAME_ALREADY_EXISTS",
  "result": null
}
```

**400 Bad Request - Password Mismatch:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Passwords do not match",
  "result": null
}
```

**Transaction Details:**
1. สร้าง User (table: users)
2. สร้าง Customer (table: customers)
3. สร้าง CustomerDuration (table: customer_durations)

**Usage Example:**
```javascript
const registerData = {
  username: 'cust09',
  password: 'SecurePass123!',
  confirmPassword: 'SecurePass123!',
  firstName: 'สมชาย',
  lastName: 'ใจดี',
  gender: 'MALE',
  dateOfBirth: '1995-01-15',
  phone: '0899999999',
  gmail: 'somchai@gmail.com',
  healthInfo: 'ไม่มีโรคประจำตัว',
  address: '123 ถนนสุขุมวิท กรุงเทพฯ 10110',
  companyName: 'บริษัท ABC จำกัด',
  companyPosition: 'Software Engineer',
  maritalStatus: 'SINGLE',
  emergencyContactName: 'สมหญิง ใจดี',
  emergencyContactRelationship: 'แม่',
  emergencyContactPhone: '0888888888',
  marketingSource: 'Facebook Ads',
  productId: 1,
  salesUsername: 'sales1',
  startDate: '2025-11-01',
  durationDays: 30,
  pricePaid: 1200.00,
  discountAmount: 0.00
};

const response = await fetch('http://localhost:8000/api/customers/durations/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(registerData)
});

const data = await response.json();
if (data.status === 'success') {
  alert('ลงทะเบียนสำเร็จ!');
}
```

---

### 5.2 Register Customer Session (แพ็กเกจครั้ง)

**Endpoint:** `POST /api/customers/sessions/register`

**Description:** ลงทะเบียนลูกค้าใหม่สำหรับแพ็กเกจ Sessions (Use Case 2.2C)

**Request Body:**
```json
{
  "username": "cust10",
  "password": "SecurePass123!",
  "confirmPassword": "SecurePass123!",
  "firstName": "สมหญิง",
  "lastName": "รักสุขภาพ",
  "gender": "FEMALE",
  "dateOfBirth": "1998-05-20",
  "phone": "0877777777",
  "gmail": "somying@gmail.com",
  "healthInfo": "แพ้อาหารทะเล",
  "address": "456 ถนนพระราม 4 กรุงเทพฯ 10330",
  "companyName": "บริษัท XYZ จำกัด",
  "companyPosition": "Marketing Manager",
  "maritalStatus": "SINGLE",
  "emergencyContactName": "สมชาย รักสุขภาพ",
  "emergencyContactRelationship": "พี่ชาย",
  "emergencyContactPhone": "0866666666",
  "marketingSource": "Google Search",
  "productId": 11,
  "trainerUsername": "trainer1",
  "salesUsername": "sales1",
  "totalSessions": 10,
  "pricePaid": 2800.00,
  "discountAmount": 0.00,
  "schedules": [
    {
      "startTime": "2025-11-05T10:00:00Z",
      "endTime": "2025-11-05T11:00:00Z",
      "dayOfWeek": "TUESDAY"
    },
    {
      "startTime": "2025-11-07T10:00:00Z",
      "endTime": "2025-11-07T11:00:00Z",
      "dayOfWeek": "THURSDAY"
    }
  ]
}
```

**Required Fields:**
- All fields from Duration registration +
- `trainerUsername`: เทรนเนอร์ที่เลือก
- `totalSessions`: จำนวนครั้งทั้งหมด
- `schedules`: รายการนัดหมาย (array)

**Schedule Object:**
- `startTime`: วันเวลาเริ่ม (RFC3339 format: `2025-11-05T10:00:00Z`)
- `endTime`: วันเวลาสิ้นสุด (RFC3339 format)
- `dayOfWeek`: วันในสัปดาห์ (`"MONDAY"`, `"TUESDAY"`, etc.)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Customer session registered successfully",
  "result": {
    "username": "cust10",
    "sessionId": 1009,
    "trainerUsername": "trainer1",
    "productId": 11,
    "totalSessions": 10,
    "schedulesCreated": 2,
    "createdSchedules": [
      {
        "scheduleId": 5020,
        "startTime": "2025-11-05T10:00:00Z",
        "endTime": "2025-11-05T11:00:00Z",
        "dayOfWeek": "TUESDAY"
      },
      {
        "scheduleId": 5021,
        "startTime": "2025-11-07T10:00:00Z",
        "endTime": "2025-11-07T11:00:00Z",
        "dayOfWeek": "THURSDAY"
      }
    ],
    "message": "Customer session registered successfully"
  }
}
```

**Transaction Details:**
1. สร้าง User (table: users)
2. สร้าง Customer (table: customers)
3. สร้าง CustomerSession (table: customer_sessions)
4. สร้าง TrainingSchedules หลายรายการ (table: training_schedules)
5. สร้าง CustomerLog (table: customer_logs, log_type: 'BOOK_SESSION')

---

### 5.3 Check Booking Permission

**Endpoint:** `GET /api/customers/sessions/check-permission`

**Description:** ตรวจสอบว่าลูกค้ามีสิทธิ์จองนัดหรือไม่ (ต้องมี Session package ACTIVE และยังมีสิทธิ์คงเหลือ)

**Query Parameters:**
- `username` (string): ชื่อผู้ใช้ลูกค้า

**Example:** `GET /api/customers/sessions/check-permission?username=cust01`

**Success Response (200 OK - Has Permission):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Permission check completed",
  "result": {
    "hasPermission": true,
    "canBook": true
  }
}
```

**Response (200 OK - No Permission):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Permission check completed",
  "result": {
    "hasPermission": false,
    "canBook": false
  }
}
```

**Use Case:** เรียกก่อนแสดงหน้าจองนัด เพื่อตรวจสอบว่ามีสิทธิ์หรือไม่

---

### 5.4 Get Active Session Packages

**Endpoint:** `GET /api/customers/sessions/active/:username`

**Description:** ดึงข้อมูล Session packages ที่ยัง ACTIVE ของลูกค้า (แสดงจำนวน sessions คงเหลือ)

**Path Parameters:**
- `username` (string): ชื่อผู้ใช้ลูกค้า

**Example:** `GET /api/customers/sessions/active/cust01`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Active sessions retrieved successfully",
  "result": [
    {
      "id": 1001,
      "customerUsername": "cust01",
      "trainerUsername": "trainer1",
      "productId": 11,
      "productName": "Yoga Sessions - 10 Pack",
      "totalSessions": 10,
      "usedSessions": 4,
      "sessionsRemaining": 6,
      "purchaseDate": "2025-10-05T00:00:00Z",
      "pricePaid": 2800.00,
      "discountAmount": 0.00,
      "status": "ACTIVE",
      "createdAt": "2025-10-05T00:00:00Z"
    }
  ]
}
```

**Response (200 OK - No Active Sessions):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Active sessions retrieved successfully",
  "result": []
}
```

**Field Descriptions:**
- `totalSessions`: จำนวนครั้งทั้งหมดที่ซื้อ
- `usedSessions`: จำนวนครั้งที่ใช้ไปแล้ว
- `sessionsRemaining`: จำนวนครั้งคงเหลือ (totalSessions - usedSessions)

**Use Case:** แสดงข้อมูลแพ็กเกจในหน้าโปรไฟล์ หรือก่อนจองนัด

---

### 5.5 Get Active Duration Packages

**Endpoint:** `GET /api/customers/durations/active/:username`

**Description:** ดึงข้อมูล Duration packages ที่ยัง ACTIVE ของลูกค้า (แสดงจำนวนวันคงเหลือ)

**Path Parameters:**
- `username` (string): ชื่อผู้ใช้ลูกค้า

**Example:** `GET /api/customers/durations/active/cust01`

**Success Response (200 OK):**
```json
{
  "status": "OK",
  "status_code": 200,
  "message": "Active duration packages retrieved successfully",
  "result": [
    {
      "id": 2001,
      "customerUsername": "cust01",
      "productId": 1,
      "productName": "1 Month Gym Pass",
      "durationDays": 30,
      "salesUsername": "sales01",
      "purchaseDate": "2025-10-01T00:00:00Z",
      "startDate": "2025-10-01T00:00:00Z",
      "endDate": "2025-10-31T00:00:00Z",
      "daysRemaining": 15,
      "pricePaid": 1500.00,
      "discountAmount": 0.00,
      "status": "ACTIVE",
      "createdAt": "2025-10-01T00:00:00Z"
    }
  ]
}
```

**Response (200 OK - No Active Durations):**
```json
{
  "status": "OK",
  "status_code": 200,
  "message": "Active duration packages retrieved successfully",
  "result": []
}
```

**Field Descriptions:**
- `durationDays`: จำนวนวันทั้งหมดที่ซื้อ
- `startDate`: วันเริ่มต้นแพ็กเกจ
- `endDate`: วันสิ้นสุดแพ็กเกจ
- `daysRemaining`: จำนวนวันคงเหลือ (คำนวณจาก DATEDIFF(endDate, CURDATE()))

**Business Logic:**
- JOIN กับตาราง `products` เพื่อดึง `product_name` และ `duration_days`
- คำนวณ `daysRemaining` ด้วย SQL: `DATEDIFF(end_date, CURDATE())`
- กรองเฉพาะ `status = 'ACTIVE'`
- เรียงลำดับตาม `created_at DESC`

**Use Case:** แสดงข้อมูลแพ็กเกจในหน้าโปรไฟล์ หรือเช็คว่ามีสิทธิ์เข้าใช้งานฟิตเนสหรือไม่

**Usage Example:**
```javascript
const response = await fetch(`http://localhost:8000/api/customers/durations/active/cust01`, {
  method: 'GET',
  headers: { 'Content-Type': 'application/json' }
});

const data = await response.json();
if (data.status === 'OK' && data.result.length > 0) {
  const activeDuration = data.result[0];
  console.log(`Days remaining: ${activeDuration.daysRemaining}`);
  console.log(`Package: ${activeDuration.productName}`);
}
```

---

## 6. Booking APIs

### 6.1 Get Booking Slots

**Endpoint:** `GET /api/bookings/slots`

**Description:** ดึงช่วงเวลาว่างสำหรับจองนัดกับเทรนเนอร์ (Use Case 3C: Q3C.3)

**Query Parameters:**
- `trainerUsername` (string, required): ชื่อ username ของเทรนเนอร์
- `customerUsername` (string, optional): ชื่อ username ของลูกค้า
- `calendarStart` (string, required): วันที่เริ่มต้น (RFC3339: `2025-11-01T00:00:00Z`)
- `calendarEnd` (string, required): วันที่สิ้นสุด (RFC3339: `2025-11-30T23:59:59Z`)

**Example:** 
```
GET /api/bookings/slots?trainerUsername=trainer1&customerUsername=cust01&calendarStart=2025-11-01T00:00:00Z&calendarEnd=2025-11-30T23:59:59Z
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Booking slots retrieved successfully",
  "result": {
    "trainerUsername": "trainer1",
    "calendarStart": "2025-11-01T00:00:00Z",
    "calendarEnd": "2025-11-30T23:59:59Z",
    "weeklyAvailability": [
      {
        "dayOfWeek": "MONDAY",
        "startTime": "09:00:00",
        "endTime": "17:00:00"
      },
      {
        "dayOfWeek": "TUESDAY",
        "startTime": "09:00:00",
        "endTime": "17:00:00"
      }
    ],
    "dayOffSlots": [
      {
        "startTime": "2025-11-06T00:00:00Z",
        "endTime": "2025-11-06T23:59:59Z"
      }
    ],
    "bookedAppointments": [
      {
        "startTime": "2025-11-01T09:00:00Z",
        "endTime": "2025-11-01T10:00:00Z",
        "customerUsername": "cust01"
      }
    ],
    "availableSlots": [],
    "customerBookings": [
      {
        "startTime": "2025-11-01T09:00:00Z",
        "endTime": "2025-11-01T10:00:00Z",
        "available": false,
        "isBooked": true,
        "bookedBy": "cust01",
        "slotType": "booked"
      }
    ],
    "message": "Booking slots retrieved successfully"
  }
}
```

**Field Descriptions:**
- `weeklyAvailability`: เวลาทำงานประจำสัปดาห์ของเทรนเนอร์
- `dayOffSlots`: วันหยุด/ช่วงเวลาที่ไม่รับนัด
- `bookedAppointments`: นัดที่ถูกจองแล้วทั้งหมด
- `availableSlots`: ช่วงเวลาว่าง (TODO: Backend ยังไม่ได้ implement การคำนวณ)
- `customerBookings`: นัดของลูกค้าที่ระบุ (ถ้ามี customerUsername)

**Use Case:** แสดงปฏิทินจองนัด โดยนำข้อมูลไป render ใน calendar component

---

### 6.2 Book Appointment

**Endpoint:** `POST /api/bookings/book`

**Description:** จองนัดหมายกับเทรนเนอร์ (Use Case 3C: Q3C.6)

**Request Body:**
```json
{
  "trainerUsername": "trainer1",
  "customerUsername": "cust01",
  "sessionId": null,
  "startTime": "2025-11-10T10:00:00Z",
  "endTime": "2025-11-10T11:00:00Z"
}
```

**Field Descriptions:**
- `trainerUsername`: ชื่อ username ของเทรนเนอร์
- `customerUsername`: ชื่อ username ของลูกค้า
- `sessionId`: ID ของ session package (ถ้าเป็น `null` จะหา ACTIVE session อัตโนมัติ)
- `startTime`: วันเวลาเริ่ม (RFC3339 format)
- `endTime`: วันเวลาสิ้นสุด (RFC3339 format)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Appointment booked successfully",
  "result": {
    "success": true,
    "message": "Appointment booked successfully",
    "trainerUsername": "trainer1",
    "customerUsername": "cust01",
    "startTime": "2025-11-10T10:00:00Z",
    "endTime": "2025-11-10T11:00:00Z",
    "sessionId": 1001,
    "remainingSession": 5
  }
}
```

**Error Responses:**

**400 - No Active Session:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Customer does not have an active session package or no sessions remaining",
  "result": {
    "success": false,
    "message": "Customer does not have an active session package or no sessions remaining"
  }
}
```

**400 - Time Slot Not Available:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Time slot is not available. Found 1 overlapping appointment(s)",
  "result": {
    "success": false,
    "message": "Time slot is not available. Found 1 overlapping appointment(s)"
  }
}
```

**Transaction Details:**
1. ตรวจสอบ ACTIVE session (auto-find ถ้า sessionId = null)
2. ตรวจสอบช่วงเวลาว่าง (CheckTimeSlotAvailability)
3. สร้าง TrainingSchedule (INSERT)
4. อัปเดต used_sessions + 1 (UPDATE customer_sessions)
5. บันทึก log (INSERT customer_logs, log_type: 'BOOK_SESSION')

**Usage Example:**
```javascript
const bookingData = {
  trainerUsername: 'trainer1',
  customerUsername: 'cust01',
  sessionId: null, // auto-find ACTIVE session
  startTime: '2025-11-10T10:00:00Z',
  endTime: '2025-11-10T11:00:00Z'
};

const response = await fetch('http://localhost:8000/api/bookings/book', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify(bookingData)
});

const data = await response.json();
if (data.result.success) {
  alert(`จองนัดสำเร็จ! คงเหลือ ${data.result.remainingSession} ครั้ง`);
}
```

---

### 6.3 Cancel Appointment

**Endpoint:** `DELETE /api/bookings/cancel/:id`

**Description:** ยกเลิกการจองนัดหมาย (คืนสิทธิ์ 1 ครั้ง)

**Path Parameters:**
- `id` (integer): Appointment ID (จาก training_schedules.id)

**Request Body:**
```json
{
  "customerUsername": "cust01"
}
```

**Example:** `DELETE /api/bookings/cancel/5007`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Appointment canceled successfully",
  "result": {
    "success": true,
    "message": "Appointment canceled successfully",
    "appointmentId": 5007,
    "customerUsername": "cust01",
    "startTime": "2025-11-01T09:00:00Z",
    "endTime": "2025-11-01T10:00:00Z",
    "sessionId": 1001,
    "remainingSessions": 7
  }
}
```

**Error Responses:**

**400 - Appointment Not Found:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Appointment not found",
  "result": {
    "success": false,
    "message": "Appointment not found"
  }
}
```

**400 - Unauthorized:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "You are not authorized to cancel this appointment",
  "result": {
    "success": false,
    "message": "You are not authorized to cancel this appointment"
  }
}
```

**400 - Past Appointment:**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Cannot cancel past appointments",
  "result": {
    "success": false,
    "message": "Cannot cancel past appointments"
  }
}
```

**Business Rules:**
1. ✅ ต้องเป็นเจ้าของนัดหมาย (appointment.customerUsername == request.customerUsername)
2. ✅ ต้องเป็นนัดที่ยังไม่ผ่านไป (time.Now().Before(appointment.StartTime))
3. ✅ คืนสิทธิ์ 1 ครั้ง (used_sessions - 1)

**Transaction Details:**
1. ตรวจสอบ appointment exists
2. ตรวจสอบ ownership
3. ตรวจสอบเวลา (ไม่ให้ยกเลิกนัดที่ผ่านไปแล้ว)
4. ลบ TrainingSchedule (DELETE)
5. ลด used_sessions - 1 (UPDATE customer_sessions)
6. บันทึก log (INSERT customer_logs, log_type: 'CANCEL_SESSION')

**Usage Example:**
```javascript
const appointmentId = 5007;
const response = await fetch(`http://localhost:8000/api/bookings/cancel/${appointmentId}`, {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    customerUsername: 'cust01'
  })
});

const data = await response.json();
if (data.result.success) {
  alert(`ยกเลิกนัดสำเร็จ! คงเหลือ ${data.result.remainingSessions} ครั้ง`);
}
```

---

## 7. Response Format

### Standard Response Structure

ทุก API จะใช้ response format เดียวกัน:

```json
{
  "status": "success" | "error",
  "status_code": 200 | 400 | 404 | 500,
  "message": "Human readable message",
  "result": {} | [] | null
}
```

### HTTP Status Codes

| Code | Status | Description |
|------|--------|-------------|
| 200 | OK | Request สำเร็จ |
| 201 | Created | สร้างข้อมูลสำเร็จ (ไม่ค่อยใช้ในระบบนี้) |
| 400 | Bad Request | Request ไม่ถูกต้อง, validation error |
| 404 | Not Found | ไม่พบข้อมูลที่ขอ |
| 500 | Internal Server Error | Server error |

---

## 8. Error Codes

### Common Error Messages

| Error Message | Description | Solution |
|--------------|-------------|----------|
| `invalid credentials` | Username/password ไม่ถูกต้อง | ตรวจสอบ username และ password |
| `USERNAME_ALREADY_EXISTS` | Username ซ้ำ | เปลี่ยน username ใหม่ |
| `Passwords do not match` | Password และ confirmPassword ไม่ตรงกัน | ตรวจสอบ password |
| `product not found` | ไม่พบ product ที่ระบุ | ตรวจสอบ productId |
| `Appointment not found` | ไม่พบนัดหมาย | ตรวจสอบ appointment ID |
| `You are not authorized to cancel this appointment` | ไม่ใช่เจ้าของนัดหมาย | ใช้ customerUsername ที่ถูกต้อง |
| `Cannot cancel past appointments` | ไม่สามารถยกเลิกนัดที่ผ่านไปแล้ว | ยกเลิกก่อนเวลานัด |
| `Time slot is not available` | ช่วงเวลาถูกจองแล้ว | เลือกช่วงเวลาอื่น |
| `Customer does not have an active session package` | ไม่มี session package หรือหมดสิทธิ์แล้ว | ต้องซื้อแพ็กเกจใหม่ |

---

## 9. Authentication Flow

### 9.1 Login → Store Token

```javascript
// 1. Login
const loginResponse = await fetch('http://localhost:8000/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ username: 'cust01', password: 'Password123!' })
});

const loginData = await loginResponse.json();

// 2. Store token และ user info
if (loginData.status === 'success') {
  localStorage.setItem('token', loginData.result.token);
  localStorage.setItem('user', JSON.stringify(loginData.result.user));
}
```

### 9.2 Authenticated Request

```javascript
// ใช้ token ในทุก request ที่ต้องการ authentication
const token = localStorage.getItem('token');

const response = await fetch('http://localhost:8000/api/bookings/book', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  credentials: 'include',
  body: JSON.stringify({...})
});
```

### 9.3 Check Authentication Status

```javascript
// ตรวจสอบว่า user ยัง login อยู่หรือไม่
const checkAuth = async () => {
  const token = localStorage.getItem('token');
  if (!token) return false;

  const response = await fetch('http://localhost:8000/api/auth/me', {
    headers: { 'Authorization': `Bearer ${token}` },
    credentials: 'include'
  });

  const data = await response.json();
  return data.result?.authenticated ?? false;
};
```

### 9.4 Logout

```javascript
// Logout
await fetch('http://localhost:8000/api/auth/logout', {
  method: 'POST',
  credentials: 'include'
});

// Clear local storage
localStorage.removeItem('token');
localStorage.removeItem('user');
```

---

## 10. Common Use Cases

### Use Case 1: ซื้อแพ็กเกจ Duration (รายเดือน)

```javascript
// 1. ดึงรายการ Duration products
const productsResponse = await fetch('http://localhost:8000/api/products/durations');
const products = await productsResponse.json();

// 2. เลือก product แล้วดูข้อมูลการชำระเงิน
const productId = 1;
const discount = 100;
const paymentResponse = await fetch(`http://localhost:8000/api/payments/info/${productId}?discount=${discount}`);
const paymentInfo = await paymentResponse.json();

// 3. กรอกข้อมูลและลงทะเบียน
const registerResponse = await fetch('http://localhost:8000/api/customers/durations/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'newuser',
    password: 'SecurePass123!',
    confirmPassword: 'SecurePass123!',
    // ... ข้อมูลอื่นๆ
    productId: productId,
    pricePaid: paymentInfo.result.payableAmount,
    discountAmount: discount
  })
});
```

### Use Case 2: ซื้อแพ็กเกจ Sessions และจองนัด

```javascript
// 1. ดึงรายการ Session products
const sessionsResponse = await fetch('http://localhost:8000/api/products/sessions');
const sessions = await sessionsResponse.json();

// 2. ลงทะเบียนพร้อมนัดหมาย
const registerResponse = await fetch('http://localhost:8000/api/customers/sessions/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'newuser',
    // ... ข้อมูลอื่นๆ
    productId: 11,
    trainerUsername: 'trainer1',
    schedules: [
      {
        startTime: '2025-11-05T10:00:00Z',
        endTime: '2025-11-05T11:00:00Z',
        dayOfWeek: 'TUESDAY'
      }
    ]
  })
});
```

### Use Case 3: ตรวจสอบและจองนัดเพิ่ม

```javascript
// 1. Login
await login('cust01', 'Password123!');

// 2. ตรวจสอบสิทธิ์
const permissionResponse = await fetch('http://localhost:8000/api/customers/sessions/check-permission?username=cust01');
const permission = await permissionResponse.json();

if (!permission.result.hasPermission) {
  alert('ไม่มีสิทธิ์จองนัด กรุณาซื้อแพ็กเกจ');
  return;
}

// 3. ดูแพ็กเกจที่มี
const packagesResponse = await fetch('http://localhost:8000/api/customers/sessions/active/cust01');
const packages = await packagesResponse.json();
console.log('Sessions คงเหลือ:', packages.result[0].sessionsRemaining);

// 4. ดูช่วงเวลาว่าง
const slotsResponse = await fetch('http://localhost:8000/api/bookings/slots?trainerUsername=trainer1&customerUsername=cust01&calendarStart=2025-11-01T00:00:00Z&calendarEnd=2025-11-30T23:59:59Z');
const slots = await slotsResponse.json();

// 5. จองนัด
const bookResponse = await fetch('http://localhost:8000/api/bookings/book', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    trainerUsername: 'trainer1',
    customerUsername: 'cust01',
    sessionId: null,
    startTime: '2025-11-10T10:00:00Z',
    endTime: '2025-11-10T11:00:00Z'
  })
});
```

---

## 11. Date/Time Format

### RFC3339 Format

API ใช้ **RFC3339** format สำหรับ date/time:

```
2025-11-01T10:00:00Z
```

**JavaScript Conversion:**
```javascript
// Date object → RFC3339 string
const date = new Date('2025-11-01T10:00:00');
const rfc3339 = date.toISOString(); // "2025-11-01T10:00:00.000Z"

// RFC3339 string → Date object
const dateObj = new Date('2025-11-01T10:00:00Z');
```

**Date-only Format:**

สำหรับ `dateOfBirth` และ `startDate` ใช้ format:
```
2025-11-01
```

---

## 12. Testing with cURL

### Login
```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"username":"cust01","password":"Password123!"}'
```

### Get Products
```bash
curl http://localhost:8000/api/products
```

### Check Phone
```bash
curl "http://localhost:8000/api/users/check-phone?phone=0811111001"
```

### Book Appointment
```bash
curl -X POST http://localhost:8000/api/bookings/book \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "trainerUsername":"trainer1",
    "customerUsername":"cust01",
    "sessionId":null,
    "startTime":"2025-11-10T10:00:00Z",
    "endTime":"2025-11-10T11:00:00Z"
  }'
```

### Cancel Appointment
```bash
curl -X DELETE http://localhost:8000/api/bookings/cancel/5007 \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"customerUsername":"cust01"}'
```

---

## 📞 Support

หากมีข้อสงสัยหรือพบปัญหา:
1. ตรวจสอบ error message ใน response
2. ดู HTTP status code
3. ตรวจสอบ request body format
4. ตรวจสอบ authentication token

**Happy Coding! 🚀**
