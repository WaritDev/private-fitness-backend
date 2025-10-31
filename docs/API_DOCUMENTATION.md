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
7. [Trainer / Working Hours APIs](#7-trainer--working-hours-apis)
8. [Trainer / Day-Offs Management APIs](#8-trainer--day-offs-management-apis)
9. [Error Codes](#9-error-codes)

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

### 4.2 Verify Payment Slip (ตรวจสอบสลิปการโอนเงิน)

**Endpoint:** `POST /api/payments/verify-slip`

**Description:** ตรวจสอบความถูกต้องของสลิปการโอนเงินผ่าน Slip2Go API แบบ realtime (ไม่เก็บข้อมูลใน database)  
**Use Case:** ยืนยันการชำระเงิน (Payment Verification)

**Request Format:** `multipart/form-data`

**Request Fields:**
- `file` (file): ไฟล์รูปภาพสลิปการโอนเงิน (image/jpeg, image/png)
- `payload` (JSON string): ข้อมูลการตรวจสอบ

**Payload Structure:**
```json
{
  "amount": 2599.50,
  "accountName": "Private Fitness - Main Account",
  "accountNumber": "123-4-56789-0",
  "accountType": "01004",
  "paymentDate": "2025-10-31"
}
```

**Payload Fields:**
- `amount` (float, required): จำนวนเงินที่ต้องการตรวจสอบ
- `accountName` (string, required): ชื่อบัญชีปลายทาง
- `accountNumber` (string, required): เลขบัญชีปลายทาง
- `accountType` (string, required): รหัสธนาคาร (e.g., "01004" for SCB)
- `paymentDate` (string, optional): วันที่โอนเงิน (YYYY-MM-DD format)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Payment verified successfully.",
  "data": {
    "slipId": "SLIP_ABC123XYZ",
    "verified": true
  }
}
```

**Error Response - Verification Failed (400 Bad Request):**
```json
{
  "status": "error",
  "message": "Payment slip verification failed. Please check slip details and try again.",
  "data": {
    "slipId": "SLIP_DEF456GHI",
    "verified": false
  }
}
```

**Frontend Integration Example:**
```javascript
// Step 1: Prepare form data
const formData = new FormData();

// Add slip image file
const fileInput = document.getElementById('slipFile');
formData.append('file', fileInput.files[0]);

// Add JSON payload
const payload = {
  amount: 2599.50,
  accountName: 'Private Fitness - Main Account',
  accountNumber: '123-4-56789-0',
  accountType: '01004', // SCB bank code
  paymentDate: '2025-10-31'
};
formData.append('payload', JSON.stringify(payload));

// Step 2: Send request
const response = await fetch('http://localhost:8000/api/payments/verify-slip', {
  method: 'POST',
  body: formData
  // Note: Do NOT set Content-Type header, browser will set it automatically with boundary
});

const result = await response.json();

// Step 3: Handle response
if (result.status === 'success' && result.data.verified) {
  alert('✅ Payment verified successfully!');
  console.log('Slip ID:', result.data.slipId);
  // Proceed with membership activation or next step
} else {
  alert('❌ Verification failed: ' + result.message);
  // Show error to user, allow retry
}
```

**Business Logic:**
1. รับไฟล์สลิปและข้อมูลการชำระเงิน
2. เรียก Slip2Go API เพื่อตรวจสอบสลิป:
   - ตรวจสอบจำนวนเงิน
   - ตรวจสอบบัญชีปลายทาง
   - ตรวจสอบวันที่โอนเงิน (ถ้ามี)
3. ส่งผลลัพธ์กลับไปยัง Frontend ทันที (ไม่เก็บข้อมูล)

**Mock Mode for Development:**
ตั้งค่า environment variable `MOCK_SLIP2GO=true` เพื่อใช้โหมดทดสอบ (ไม่ใช้ API จริง):
```bash
# ใน .env
MOCK_SLIP2GO=true
```

เมื่อเปิดโหมด Mock:
- ระบบจะ return ผลตรวจสอบ verified=true ทันที
- ไม่มีการเรียก Slip2Go API จริง
- ประหยัด API quota (Slip2Go มี 100 ครั้งทดสอบฟรี)
- เหมาะสำหรับ development และ testing

**ข้อควรระวัง:**
- API นี้เป็น **stateless** (ไม่เก็บข้อมูลใน database)
- ผลลัพธ์ได้จาก Slip2Go API แบบ realtime
- ไฟล์รูปสลิปต้องเป็น image format (JPEG, PNG)
- ต้องระบุ `accountType` ให้ถูกต้องตามรหัสธนาคาร (เช่น SCB = "01004")
- ควรแสดง loading indicator ระหว่างรอผลตรวจสอบ (อาจใช้เวลา 3-5 วินาที)
- หากต้องการเก็บประวัติการชำระเงิน ให้ Frontend เก็บไว้หลังได้รับ verified=true

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

## 6. Member / Check-in APIs

### 6.1 Generate QR Code

**Endpoint:** `POST /api/member/qrcode`

**Description:** สร้าง QR Code สำหรับ Check-in เข้าฟิตเนส (Use Case 5C)

**Authentication:** Required (JWT token)

**Request Body:**
```json
{
  "packageType": "DURATION"
}
```

**Request Fields:**
- `packageType` (string, required): ประเภทแพ็กเกจ - `"DURATION"` หรือ `"SESSION"`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "status_code": 200,
  "message": "QR Code generated successfully",
  "result": {
    "qrCodeUrl": "http://localhost:8000/api/checkin?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "packageType": "DURATION",
    "expiresIn": 60
  }
}
```

**Response Fields:**
- `qrCodeUrl` (string): URL สำหรับฝังใน QR Code (เปิดได้ที่เครื่องสแกน)
- `token` (string): JWT token (หมดอายุใน 60 วินาที)
- `packageType` (string): ประเภทแพ็กเกจที่เลือก
- `expiresIn` (number): เวลาหมดอายุ (วินาที)

**Error Response (400 Bad Request):**
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Invalid package type (must be DURATION or SESSION)",
  "result": null
}
```

**Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "status_code": 401,
  "message": "Authentication required",
  "result": null
}
```

**Business Logic:**
1. ตรวจสอบ JWT token ของผู้ใช้
2. Validate `packageType` (ต้องเป็น "DURATION" หรือ "SESSION")
3. สร้าง QR token พิเศษที่หมดอายุใน **60 วินาที**
4. Token มี payload: `{ sub: username, packageType: "DURATION/SESSION", exp: timestamp }`
5. Frontend นำ `qrCodeUrl` ไปสร้าง QR Code ด้วย library เช่น `react-qr-code`

**Use Case:** ลูกค้ากด Toggle เลือก Duration/Session → เรียก API นี้ → แสดง QR Code

**Usage Example:**
```javascript
const response = await fetch('http://localhost:8000/api/member/qrcode', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${userToken}`
  },
  body: JSON.stringify({
    packageType: 'DURATION'
  })
});

const data = await response.json();
if (data.status === 'success') {
  // ใช้ react-qr-code สร้าง QR
  const qrUrl = data.result.qrCodeUrl;
  console.log('QR Code URL:', qrUrl);
  console.log('Expires in:', data.result.expiresIn, 'seconds');
}
```

---

### 6.2 Check-in via QR Code

**Endpoint:** `GET /api/checkin?token={qr_token}`

**Description:** สแกน QR Code เพื่อบันทึกการเข้าใช้งานฟิตเนส (Use Case 5C: Q5C.1 และ Q5C.2)

**Authentication:** None (Public endpoint - accessed by QR scanner)

**Query Parameters:**
- `token` (string, required): JWT token จาก QR Code

**Example:** 
```
GET /api/checkin?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200 OK - HTML):**
```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Check-in Success</title>
</head>
<body>
  <div class="container">
    <div class="success">✅</div>
    <h1>Welcome, John!</h1>
    <p>User: <strong>cust01</strong></p>
    <div class="badge">DURATION Package</div>
  </div>
</body>
</html>
```

**Error Response (400 Bad Request):**
```
Missing token parameter
```

**Error Response (401 Unauthorized):**
```
Invalid or expired QR code
```

**Error Response (500 Internal Server Error):**
```
Check-in failed: user not found
```

**Business Logic:**
1. รับ `token` จาก query string
2. Verify JWT token (ตรวจสอบ signature และ expiry)
3. Extract `username` และ `packageType` จาก token payload
4. **Q5C.1**: สร้าง log ใน `customer_logs` (log_type = 'CHECK_IN')
5. **Q5C.2**: ถ้าเป็น SESSION package → อัปเดต `used_sessions + 1`
6. ดึงชื่อผู้ใช้ (first_name) จาก database
7. แสดง HTML response "Welcome, [FirstName]!"

**Database Operations:**

**Q5C.1 - สร้าง Log (ทั้ง DURATION และ SESSION):**
```sql
INSERT INTO customer_logs (
  customer_username,
  log_type,
  created_at
) VALUES (
  'cust01',
  'CHECK_IN',
  NOW()
);
```

**Q5C.2 - Update Used Sessions (เฉพาะ SESSION):**
```sql
UPDATE customer_sessions
SET used_sessions = used_sessions + 1,
    updated_at = NOW()
WHERE customer_username = 'cust01'
  AND status = 'ACTIVE'
  AND used_sessions < total_sessions
ORDER BY created_at DESC
LIMIT 1;
```

**Use Case Flow:**
1. ลูกค้าเปิดแอพ → คลิก "Member" → เลือก Duration/Session
2. Frontend เรียก `POST /api/member/qrcode` → ได้ URL
3. แสดง QR Code บนมือถือ
4. ลูกค้ายื่นมือถือให้เครื่องสแกนหน้าร้าน
5. เครื่องสแกนเปิด browser → เข้า URL: `GET /api/checkin?token=xxx`
6. Backend บันทึก log + อัปเดต sessions (ถ้าเป็น SESSION)
7. แสดงหน้าจอ "Welcome, John!" + เปิดประตู

**Security:**
- Token หมดอายุใน **60 วินาที** เท่านั้น
- ใช้ได้ครั้งเดียว (แต่ระบบไม่ enforce - อาจมี race condition)
- ไม่ต้อง authenticate เพราะ token เป็น proof แล้ว

**Frontend Integration:**
```javascript
// 1. Generate QR Code
import QRCode from 'react-qr-code';

function MemberCard({ user, packageType }) {
  const [qrUrl, setQrUrl] = useState('');
  const [expiresAt, setExpiresAt] = useState(null);

  const handleGenerateQR = async () => {
    const res = await fetch('/api/member/qrcode', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${userToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ packageType })
    });

    const data = await res.json();
    if (data.status === 'success') {
      setQrUrl(data.result.qrCodeUrl);
      setExpiresAt(Date.now() + data.result.expiresIn * 1000);
    }
  };

  return (
    <div>
      <button onClick={handleGenerateQR}>
        Generate QR Code
      </button>
      {qrUrl && (
        <div>
          <QRCode value={qrUrl} size={256} />
          <p>Expires in: {Math.ceil((expiresAt - Date.now()) / 1000)}s</p>
        </div>
      )}
    </div>
  );
}
```

---

## 7. Booking APIs

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

## 7. Trainer / Working Hours APIs

### 7.1 Get Working Hours (Trainer)

**Endpoint:** `GET /api/trainers/working-hours`

**Description:** ดึงรายการเวลาทำงานทั้งหมดของเทรนเนอร์ (Use Case 1P Step 2: Q1P.1)

**Authorization:** Required (JWT Token - Trainer role)

**Business Logic:**
1. ดึง `trainerUsername` จาก JWT token
2. Query เวลาทำงานทั้งหมดจากตาราง `training_availabilities` (Q1P.1)
3. เรียงลำดับตาม day_of_week (MONDAY → SUNDAY) และ start_time
4. แสดงผลในรูปแบบตาราง

**SQL Query (Q1P.1):**
```sql
SELECT
  id,
  trainer_username,
  day_of_week,
  start_time,
  end_time
FROM training_availabilities
WHERE trainer_username = ?
ORDER BY FIELD(day_of_week, 'MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY'),
         start_time ASC;
```

**Success Response (200 OK):**
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
      "dayOfWeek": "WEDNESDAY",
      "startTime": "10:00",
      "endTime": "16:00"
    }
  ]
}
```

**Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "Unauthorized - JWT token required"
}
```

**Usage Example (React + TypeScript):**
```typescript
// Get working hours
const getWorkingHours = async () => {
  const token = localStorage.getItem('token');
  
  const response = await fetch('http://localhost:8000/api/trainers/working-hours', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });

  const data = await response.json();
  if (data.status === 'success') {
    console.log('Working Hours:', data.workingHours);
    // Display in table format
  }
};
```

---

### 7.2 Add Working Time (Trainer)

**Endpoint:** `POST /api/trainers/working-hours`

**Description:** เพิ่มเวลาทำงานใหม่ของเทรนเนอร์ (Use Case 1P Step 6-9: Q1P.2 + Q1P.3)

**Authorization:** Required (JWT Token - Trainer role)

**Business Logic:**
1. ดึง `trainerUsername` จาก JWT token
2. **Step 7: Validation**
   - ตรวจสอบค่าว่าง: `dayOfWeek`, `startTime`, `endTime` (Model)
   - ตรวจสอบรูปแบบเวลา: HH:MM format (Model)
   - ตรวจสอบ `endTime` > `startTime` (Model)
   - ตรวจสอบเวลาทับซ้อน: Q1P.2 ต้อง return `overlapped_count = 0` (Database)
3. **Step 8:** หากผ่าน validation → บันทึกข้อมูลด้วย Q1P.3
4. **Step 9:** แสดง success message และ redirect กลับหน้า Working Hours

**SQL Queries:**

**Q1P.2: Check Time Overlap**
```sql
SELECT COUNT(id) AS overlapped_count
FROM training_availabilities
WHERE trainer_username = ?
  AND day_of_week = ?
  AND (? < end_time AND ? > start_time);
```

**Q1P.3: Insert New Working Time**
```sql
INSERT INTO training_availabilities (
  trainer_username,
  day_of_week,
  start_time,
  end_time
) VALUES (?, ?, ?, ?);
```

**Request Body:**
```json
{
  "dayOfWeek": "TUESDAY",
  "startTime": "09:00",
  "endTime": "13:00"
}
```

**Request Validation:**
- `dayOfWeek`: Required, must be one of: `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`, `SUNDAY`
- `startTime`: Required, format `HH:MM` (24-hour)
- `endTime`: Required, format `HH:MM` (24-hour), must be after `startTime`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Working time added successfully"
}
```

**Error Responses:**

**400 Bad Request - Invalid Time Format:**
```json
{
  "status": "error",
  "message": "Invalid start time format. Expected HH:MM"
}
```

**400 Bad Request - End Time Before Start Time:**
```json
{
  "status": "error",
  "message": "End time must be after start time"
}
```

**400 Bad Request - Time Overlap:**
```json
{
  "status": "error",
  "message": "Working time overlaps with existing schedule"
}
```

**400 Bad Request - Invalid Day:**
```json
{
  "status": "error",
  "message": "Invalid dayOfWeek. Must be MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, or SUNDAY"
}
```

**Usage Example (React + TypeScript):**
```typescript
interface AddWorkingTimeRequest {
  dayOfWeek: 'MONDAY' | 'TUESDAY' | 'WEDNESDAY' | 'THURSDAY' | 'FRIDAY' | 'SATURDAY' | 'SUNDAY';
  startTime: string; // HH:MM format
  endTime: string;   // HH:MM format
}

// Add working time
const addWorkingTime = async (data: AddWorkingTimeRequest) => {
  const token = localStorage.getItem('token');
  
  const response = await fetch('http://localhost:8000/api/trainers/working-hours', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  const result = await response.json();
  
  if (result.status === 'success') {
    // Show success popup
    alert('Working time added successfully');
    // Redirect back to Working Hours page
    window.location.href = '/trainer/working-hours';
  } else {
    // Show error popup
    alert(result.message);
  }
};

// Example usage in form submit
const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault();
  
  const formData: AddWorkingTimeRequest = {
    dayOfWeek: 'TUESDAY',
    startTime: '09:00',
    endTime: '13:00'
  };
  
  addWorkingTime(formData);
};
```

**Frontend Workflow (Use Case 1P):**
```
1. Trainer clicks "Working Hours" menu
   ↓
2. System redirects to Working Hours page
   → Call GET /api/trainers/working-hours (Q1P.1)
   ↓
3. Display working hours table
   ↓
4. Trainer clicks "Add Working Time" button
   ↓
5. System shows Add Working Time form
   - Day_of_Week dropdown (Monday-Sunday)
   - Start_Time input (HH:MM)
   - End_Time input (HH:MM)
   ↓
6. Trainer fills form and clicks "Save"
   → Call POST /api/trainers/working-hours (Q1P.2 + Q1P.3)
   ↓
7. System validates:
   - Empty fields → Show error in form
   - Invalid time format → Show error popup
   - End_Time <= Start_Time → Show error popup
   - Time overlap → Show error popup
   ↓
8. If validation passes:
   → Show success popup "Working time added successfully"
   → Redirect back to Working Hours page
   → Auto refresh with GET /api/trainers/working-hours
```

**Security Considerations:**
- ✅ JWT authentication required (Trainer role)
- ✅ Username extracted from JWT token (prevent spoofing)
- ✅ Server-side validation for all fields
- ✅ Time overlap check at database level
- ✅ SQL injection prevention (parameterized queries via sqlc)

---

### 7.3 Update Working Time (Trainer)

**Endpoint:** `PUT /api/trainers/working-hours/:id`

**Description:** แก้ไขเวลาทำงานที่มีอยู่แล้ว (Q1P.4)

**Authorization:** Required (JWT Token - Trainer role)

**URL Parameters:**
- `id` (integer, required) - ID ของเวลาทำงานที่ต้องการแก้ไข

**Business Logic:**
1. ตรวจสอบว่า working hour นี้เป็นของ trainer ที่ login หรือไม่ (ownership validation)
2. Validate time format และ endTime > startTime
3. อัพเดทข้อมูลด้วย Q1P.4
4. Return success response

**SQL Query (Q1P.4):**
```sql
UPDATE training_availabilities
SET
  day_of_week = ?,
  start_time = ?,
  end_time = ?
WHERE id = ?;
```

**Request Body:**
```json
{
  "dayOfWeek": "TUESDAY",
  "startTime": "10:00",
  "endTime": "14:00"
}
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Working time updated successfully"
}
```

**Error Responses:**

**400 Bad Request - Not Owner:**
```json
{
  "status": "error",
  "message": "Working hour not found or does not belong to you"
}
```

**400 Bad Request - Invalid Time:**
```json
{
  "status": "error",
  "message": "End time must be after start time"
}
```

**Usage Example:**
```typescript
const updateWorkingTime = async (id: number, data: UpdateWorkingTimeRequest) => {
  const token = localStorage.getItem('token');
  
  const response = await fetch(`http://localhost:8000/api/trainers/working-hours/${id}`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  const result = await response.json();
  
  if (result.status === 'success') {
    alert('Working time updated successfully');
    // Refresh working hours list
  } else {
    alert(result.message);
  }
};
```

---

### 7.4 Delete Working Time (Trainer)

**Endpoint:** `DELETE /api/trainers/working-hours/:id`

**Description:** ลบเวลาทำงาน (Q1P.5)

**Authorization:** Required (JWT Token - Trainer role)

**URL Parameters:**
- `id` (integer, required) - ID ของเวลาทำงานที่ต้องการลบ

**Business Logic:**
1. ตรวจสอบว่า working hour นี้เป็นของ trainer ที่ login หรือไม่ (ownership validation)
2. ลบข้อมูลด้วย Q1P.5
3. Return success response

**SQL Query (Q1P.5):**
```sql
DELETE FROM training_availabilities
WHERE id = ?;
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Working time deleted successfully"
}
```

**Error Response (404 Not Found):**
```json
{
  "status": "error",
  "message": "Working hour not found or does not belong to you"
}
```

**Usage Example:**
```typescript
const deleteWorkingTime = async (id: number) => {
  if (!confirm('Are you sure you want to delete this working time?')) {
    return;
  }

  const token = localStorage.getItem('token');
  
  const response = await fetch(`http://localhost:8000/api/trainers/working-hours/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  const result = await response.json();
  
  if (result.status === 'success') {
    alert('Working time deleted successfully');
    // Refresh working hours list
  } else {
    alert(result.message);
  }
};
```

**Frontend Integration Example:**
```typescript
// Complete Working Hours Management Component
interface WorkingHour {
  availabilityId: number;
  dayOfWeek: string;
  startTime: string;
  endTime: string;
}

const WorkingHoursPage: React.FC = () => {
  const [workingHours, setWorkingHours] = useState<WorkingHour[]>([]);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState({ dayOfWeek: '', startTime: '', endTime: '' });

  // Fetch working hours
  const fetchWorkingHours = async () => {
    const token = localStorage.getItem('token');
    const response = await fetch('http://localhost:8000/api/trainers/working-hours', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await response.json();
    if (data.status === 'success') {
      setWorkingHours(data.workingHours);
    }
  };

  // Update working time
  const handleUpdate = async (id: number) => {
    const token = localStorage.getItem('token');
    const response = await fetch(`http://localhost:8000/api/trainers/working-hours/${id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(editForm)
    });

    const result = await response.json();
    if (result.status === 'success') {
      alert('Updated successfully');
      setEditingId(null);
      fetchWorkingHours();
    } else {
      alert(result.message);
    }
  };

  // Delete working time
  const handleDelete = async (id: number) => {
    if (!confirm('Delete this working time?')) return;
    
    const token = localStorage.getItem('token');
    const response = await fetch(`http://localhost:8000/api/trainers/working-hours/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    });

    const result = await response.json();
    if (result.status === 'success') {
      alert('Deleted successfully');
      fetchWorkingHours();
    } else {
      alert(result.message);
    }
  };

  return (
    <div>
      <h1>Working Hours Management</h1>
      <table>
        <thead>
          <tr>
            <th>Day</th>
            <th>Start Time</th>
            <th>End Time</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {workingHours.map((hour) => (
            <tr key={hour.availabilityId}>
              {editingId === hour.availabilityId ? (
                <>
                  <td>
                    <select value={editForm.dayOfWeek} onChange={(e) => setEditForm({...editForm, dayOfWeek: e.target.value})}>
                      <option value="MONDAY">Monday</option>
                      <option value="TUESDAY">Tuesday</option>
                      {/* ... other days */}
                    </select>
                  </td>
                  <td><input type="time" value={editForm.startTime} onChange={(e) => setEditForm({...editForm, startTime: e.target.value})} /></td>
                  <td><input type="time" value={editForm.endTime} onChange={(e) => setEditForm({...editForm, endTime: e.target.value})} /></td>
                  <td>
                    <button onClick={() => handleUpdate(hour.availabilityId)}>Save</button>
                    <button onClick={() => setEditingId(null)}>Cancel</button>
                  </td>
                </>
              ) : (
                <>
                  <td>{hour.dayOfWeek}</td>
                  <td>{hour.startTime}</td>
                  <td>{hour.endTime}</td>
                  <td>
                    <button onClick={() => {
                      setEditingId(hour.availabilityId);
                      setEditForm({
                        dayOfWeek: hour.dayOfWeek,
                        startTime: hour.startTime,
                        endTime: hour.endTime
                      });
                    }}>Edit</button>
                    <button onClick={() => handleDelete(hour.availabilityId)}>Delete</button>
                  </td>
                </>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
```

---

## 8. Trainer / Day-Offs Management APIs

### 8.1 Get Day-Offs (Trainer)

**Endpoint:** `GET /api/trainers/day-offs`

**Description:** ดึงรายการวันหยุดทั้งหมดของเทรนเนอร์ (Use Case 3P: Q3P.1)

**Authorization:** Required (JWT Token - Trainer role)

**Business Logic:**
1. ดึง `trainerUsername` จาก JWT token
2. Query วันหยุดทั้งหมดจากตาราง `training_schedules` ที่ `schedule_type = 'DAY_OFF'` (Q3P.1)
3. เรียงลำดับจากวันที่ล่าสุด (ORDER BY start_time DESC)
4. แสดงผลในรูปแบบตาราง (ช่วงเวลาเป็น full day: 00:00:00 - 23:59:59)

**SQL Query (Q3P.1):**
```sql
SELECT id, trainer_username, start_time, end_time
FROM training_schedules
WHERE trainer_username = ? AND schedule_type = 'DAY_OFF'
ORDER BY start_time DESC;
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Day-offs retrieved successfully",
  "dayOffs": [
    {
      "scheduleId": 10,
      "startTime": "2025-12-25T00:00:00Z",
      "endTime": "2025-12-25T23:59:59Z"
    },
    {
      "scheduleId": 8,
      "startTime": "2025-12-01T00:00:00Z",
      "endTime": "2025-12-01T23:59:59Z"
    }
  ]
}
```

**Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "Unauthorized - JWT token required"
}
```

**Usage Example (React + TypeScript):**
```typescript
// Get day-offs
const getDayOffs = async () => {
  const token = localStorage.getItem('token');
  
  const response = await fetch('http://localhost:8000/api/trainers/day-offs', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });

  const data = await response.json();
  if (data.status === 'success') {
    console.log('Day-Offs:', data.dayOffs);
    // Display in calendar format
  }
};
```

---

### 8.2 Add Day-Off (Trainer)

**Endpoint:** `POST /api/trainers/day-offs`

**Description:** เพิ่มวันหยุดใหม่ของเทรนเนอร์ พร้อม validation (Use Case 3P: Q3P.2 + Q3P.3 + Q3P.4)

**Authorization:** Required (JWT Token - Trainer role)

**Business Logic Flow:**
1. ดึง `trainerUsername` จาก JWT token
2. **Input Validation:**
   - Parse `dayOffDate` (YYYY-MM-DD format)
   - Convert to full day range:
     - `NewStartTime` = Day_Off_Date 00:00:00
     - `NewEndTime` = Day_Off_Date 23:59:59
3. **Duplicate Validation (Q3P.2):**
   - ตรวจสอบว่ามีวันหยุดในวันเดียวกันแล้วหรือไม่
   - ถ้า `duplicate_count > 0` → Return error "Day-off already exists for this date"
4. **Appointment Overlap Validation (Q3P.3):**
   - ตรวจสอบว่ามีนัดหมายที่ทับกับวันหยุดหรือไม่
   - ถ้า `overlapped_count > 0` → Return error "Cannot create day-off: There are existing appointments on this date"
5. **Insert Day-Off (Q3P.4):**
   - บันทึกวันหยุดใหม่ด้วย `schedule_type = 'DAY_OFF'`
   - Return success message

**SQL Queries:**

**Q3P.2: Check Duplicate Day-Off**
```sql
SELECT COUNT(id) AS duplicate_count
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'DAY_OFF'
  AND DATE(start_time) = ?;
```

**Q3P.3: Check Appointment Overlap**
```sql
SELECT COUNT(id) AS overlapped_count
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'APPOINTMENT'
  AND start_time < ?
  AND end_time > ?;
```

**Q3P.4: Insert Day-Off**
```sql
INSERT INTO training_schedules (
  trainer_username,
  start_time,
  end_time,
  schedule_type
) VALUES (?, ?, ?, 'DAY_OFF');
```

**Request Body:**
```json
{
  "dayOffDate": "2025-12-25"
}
```

**Request Validation:**
- `dayOffDate`: Required, format `YYYY-MM-DD`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Day off created successfully"
}
```

**Error Responses:**

**400 Bad Request - Invalid Date Format:**
```json
{
  "status": "error",
  "message": "Invalid date format. Expected YYYY-MM-DD"
}
```

**400 Bad Request - Duplicate Day-Off:**
```json
{
  "status": "error",
  "message": "Day-off already exists for this date"
}
```

**400 Bad Request - Appointment Overlap:**
```json
{
  "status": "error",
  "message": "Cannot create day-off: There are existing appointments on this date"
}
```

**400 Bad Request - Missing Required Field:**
```json
{
  "status": "error",
  "message": "dayOffDate is required"
}
```

**Usage Example (React + TypeScript):**
```typescript
interface AddDayOffRequest {
  dayOffDate: string; // YYYY-MM-DD format
}

// Add day-off
const addDayOff = async (date: string) => {
  const token = localStorage.getItem('token');
  
  const response = await fetch('http://localhost:8000/api/trainers/day-offs', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      dayOffDate: date // "2025-12-25"
    })
  });

  const data = await response.json();
  
  if (data.status === 'success') {
    alert('วันหยุดถูกเพิ่มเรียบร้อยแล้ว');
    // Refresh day-offs list
    await getDayOffs();
  } else {
    alert(data.message); // Show validation error
  }
};

// Example usage:
// addDayOff('2025-12-25');
```

---

### 8.3 Delete Day-Off (Trainer)

**Endpoint:** `DELETE /api/trainers/day-offs/:id`

**Description:** ลบวันหยุดของเทรนเนอร์ (Use Case 3P: Q3P.5)

**Authorization:** Required (JWT Token - Trainer role)

**Business Logic:**
1. ดึง `trainerUsername` จาก JWT token
2. ดึง `scheduleId` จาก URL parameter
3. **Ownership Validation:**
   - ตรวจสอบว่า day-off นี้เป็นของ trainer ที่ login อยู่
   - ถ้าไม่ใช่ → Return error "Day-off not found or does not belong to you"
4. **Delete Day-Off (Q3P.5):**
   - ลบ record ที่มี `id = scheduleId` และ `schedule_type = 'DAY_OFF'`
   - Return success message

**SQL Query (Q3P.5):**
```sql
DELETE FROM training_schedules
WHERE id = ? AND schedule_type = 'DAY_OFF';
```

**URL Parameters:**
- `id` (required): Schedule ID ของวันหยุดที่ต้องการลบ

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Day-off deleted successfully"
}
```

**Error Responses:**

**400 Bad Request - Invalid ID:**
```json
{
  "status": "error",
  "message": "Invalid ID parameter"
}
```

**404 Not Found - Day-Off Not Found:**
```json
{
  "status": "error",
  "message": "Day-off not found or does not belong to you"
}
```

**Usage Example (React + TypeScript):**
```typescript
// Delete day-off
const deleteDayOff = async (scheduleId: number) => {
  const token = localStorage.getItem('token');
  
  if (!confirm('คุณต้องการลบวันหยุดนี้หรือไม่?')) {
    return;
  }
  
  const response = await fetch(`http://localhost:8000/api/trainers/day-offs/${scheduleId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });

  const data = await response.json();
  
  if (data.status === 'success') {
    alert('วันหยุดถูกลบเรียบร้อยแล้ว');
    // Refresh day-offs list
    await getDayOffs();
  } else {
    alert(data.message); // Show error
  }
};

// Example usage:
// deleteDayOff(10);
```

---

### 8.4 Use Case 3P Complete Flow

**Frontend Implementation Example:**

```typescript
import { useState, useEffect } from 'react';

interface DayOff {
  scheduleId: number;
  startTime: string;
  endTime: string;
}

const TrainerDayOffsPage = () => {
  const [dayOffs, setDayOffs] = useState<DayOff[]>([]);
  const [selectedDate, setSelectedDate] = useState('');
  const [loading, setLoading] = useState(false);
  const token = localStorage.getItem('token');

  // Fetch day-offs on component mount
  useEffect(() => {
    fetchDayOffs();
  }, []);

  const fetchDayOffs = async () => {
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8000/api/trainers/day-offs', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      const data = await response.json();
      if (data.status === 'success') {
        setDayOffs(data.dayOffs);
      }
    } catch (error) {
      console.error('Error fetching day-offs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddDayOff = async () => {
    if (!selectedDate) {
      alert('กรุณาเลือกวันที่');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch('http://localhost:8000/api/trainers/day-offs', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ dayOffDate: selectedDate })
      });
      
      const data = await response.json();
      
      if (data.status === 'success') {
        alert('เพิ่มวันหยุดสำเร็จ');
        setSelectedDate('');
        await fetchDayOffs();
      } else {
        alert(data.message);
      }
    } catch (error) {
      console.error('Error adding day-off:', error);
      alert('เกิดข้อผิดพลาดในการเพิ่มวันหยุด');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteDayOff = async (scheduleId: number) => {
    if (!confirm('คุณต้องการลบวันหยุดนี้หรือไม่?')) {
      return;
    }

    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8000/api/trainers/day-offs/${scheduleId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      
      const data = await response.json();
      
      if (data.status === 'success') {
        alert('ลบวันหยุดสำเร็จ');
        await fetchDayOffs();
      } else {
        alert(data.message);
      }
    } catch (error) {
      console.error('Error deleting day-off:', error);
      alert('เกิดข้อผิดพลาดในการลบวันหยุด');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1>จัดการวันหยุด (Day-Offs Management)</h1>
      
      {/* Add Day-Off Form */}
      <div>
        <h2>เพิ่มวันหยุด</h2>
        <input
          type="date"
          value={selectedDate}
          onChange={(e) => setSelectedDate(e.target.value)}
          min={new Date().toISOString().split('T')[0]}
        />
        <button onClick={handleAddDayOff} disabled={loading}>
          เพิ่มวันหยุด
        </button>
      </div>

      {/* Day-Offs List */}
      <div>
        <h2>รายการวันหยุด</h2>
        {loading ? (
          <p>Loading...</p>
        ) : dayOffs.length === 0 ? (
          <p>ไม่มีวันหยุด</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>วันที่</th>
                <th>เวลาเริ่มต้น</th>
                <th>เวลาสิ้นสุด</th>
                <th>จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {dayOffs.map((dayOff) => (
                <tr key={dayOff.scheduleId}>
                  <td>
                    {new Date(dayOff.startTime).toLocaleDateString('th-TH')}
                  </td>
                  <td>
                    {new Date(dayOff.startTime).toLocaleTimeString('th-TH')}
                  </td>
                  <td>
                    {new Date(dayOff.endTime).toLocaleTimeString('th-TH')}
                  </td>
                  <td>
                    <button 
                      onClick={() => handleDeleteDayOff(dayOff.scheduleId)}
                      disabled={loading}
                    >
                      ลบ
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default TrainerDayOffsPage;
```

---

## 9. Error Codes

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
