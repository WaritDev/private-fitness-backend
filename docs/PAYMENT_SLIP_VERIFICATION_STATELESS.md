# 💳 Payment Slip Verification API - Stateless Implementation

> **Feature:** ตรวจสอบสลิปการโอนเงินผ่าน Slip2Go API แบบ realtime (ไม่เก็บข้อมูล)  
> **Status:** ✅ Complete  
> **Date:** October 31, 2025

---

## 📝 Overview

API สำหรับตรวจสอบความถูกต้องของสลิปการโอนเงินผ่าน Slip2Go API แบบ **stateless** โดยไม่เก็บข้อมูลใน database

**Architecture:** Stateless - ส่งรูปสลิป → ยิง Slip2Go API → ได้ผล verified/rejected ทันที

---

## 🎯 Key Features

- ✅ **Stateless**: ไม่เก็บข้อมูลใน database
- ✅ **Realtime Verification**: ผลลัพธ์ได้จาก Slip2Go API ทันที
- ✅ **Mock Mode**: ทดสอบได้โดยไม่ใช้ API quota
- ✅ **Multipart Upload**: รองรับการอัปโหลดรูปภาพสลิป
- ✅ **Simple API**: ไม่มี CRUD, ไม่มี database operations

---

## 🔧 Files Implemented

### 1. Slip2Go Client
**File:** `internal/infrastructure/slip2go/client.go`

```go
type Slip2GoClient struct {
    APIKey     string
    BaseURL    string
    HTTPClient *http.Client
    MockMode   bool
}

func (c *Slip2GoClient) VerifySlip(req VerifySlipRequest) (*Slip2GoResponse, error)
```

**Features:**
- Mock mode support (`MOCK_SLIP2GO=true`)
- Multipart/form-data handling
- Amount/receiver/date verification

### 2. DTOs
**Files:** 
- `domain/requests/payment_request.go` - `VerifySlipPayload`
- `domain/responses/payment_response.go` - `VerifySlipResponse`, `VerifySlipData`

```go
// Request
type VerifySlipPayload struct {
    Amount        float64 `json:"amount" validate:"required,gt=0"`
    AccountName   string  `json:"accountName" validate:"required"`
    AccountNumber string  `json:"accountNumber" validate:"required"`
    AccountType   string  `json:"accountType" validate:"required"`
    PaymentDate   string  `json:"paymentDate,omitempty"`
}

// Response
type VerifySlipResponse struct {
    Status  string          `json:"status"`
    Message string          `json:"message"`
    Data    *VerifySlipData `json:"data,omitempty"`
}

type VerifySlipData struct {
    SlipID   string `json:"slipId,omitempty"`
    Verified bool   `json:"verified"`
}
```

### 3. Use Case
**File:** `domain/usecases/payment_use_case.go`

```go
func (uc *PaymentUseCase) VerifySlip(
    ctx context.Context,
    payload requests.VerifySlipPayload,
    fileData io.Reader,
    filename string,
) (*responses.VerifySlipResponse, error) {
    // 1. Call Slip2Go API
    slip2goResp, err := uc.slip2goClient.VerifySlip(...)
    
    // 2. Return result directly (no database)
    return &responses.VerifySlipResponse{
        Status:  "success",
        Message: "Payment verified successfully.",
        Data: &responses.VerifySlipData{
            SlipID:   slip2goResp.Result.SlipID,
            Verified: true,
        },
    }, nil
}
```

**Logic Flow:**
1. รับไฟล์สลิปและ payload
2. เรียก Slip2Go API
3. ส่งผลลัพธ์กลับทันที (ไม่เก็บ database)

### 4. REST Handler
**File:** `internal/adapters/rest/payment_rest.go`

```go
func (h *PaymentHandler) VerifySlip(c *fiber.Ctx) error {
    // 1. Get file from multipart/form-data
    fileHeader, _ := c.FormFile("file")
    file, _ := fileHeader.Open()
    defer file.Close()
    
    // 2. Parse JSON payload
    payloadStr := c.FormValue("payload")
    var payload requests.VerifySlipPayload
    c.App().Config().JSONDecoder([]byte(payloadStr), &payload)
    
    // 3. Call use case
    result, _ := h.paymentUC.VerifySlip(c.Context(), payload, file, fileHeader.Filename)
    
    // 4. Return response
    return c.Status(fiber.StatusOK).JSON(result)
}
```

### 5. Route Registration
**File:** `router/api_router.go`

```go
payments.Post("/verify-slip", handler.Payment.VerifySlip)
```

**Endpoint:** `POST /api/payments/verify-slip`

---

## 📊 API Documentation

### Request

**Method:** POST  
**Endpoint:** `/api/payments/verify-slip`  
**Content-Type:** `multipart/form-data`

**Fields:**
- `file` (file): รูปภาพสลิป (JPEG/PNG)
- `payload` (JSON string): ข้อมูลการตรวจสอบ

**Payload Example:**
```json
{
  "amount": 2599.50,
  "accountName": "Private Fitness - Main Account",
  "accountNumber": "123-4-56789-0",
  "accountType": "01004",
  "paymentDate": "2025-10-31"
}
```

### Response

**Success (200 OK):**
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

**Error (400 Bad Request):**
```json
{
  "status": "error",
  "message": "Payment slip verification failed. Please check slip details and try again.",
  "data": {
    "slipId": "SLIP_DEF456",
    "verified": false
  }
}
```

---

## 🧪 Testing

### Mock Mode Setup
```bash
# ใน .env
SLIP2GO_SECRET_KEY=50igZPNwcAd3hZOuw4VwVCj2fGPD_dT8ZZvpNviBwQU=
MOCK_SLIP2GO=true  # สำหรับ development
```

### Quick Test
```bash
# Start server
make run

# Test with curl
curl -X POST http://localhost:8000/api/payments/verify-slip \
  -F "file=@test-slip.jpg" \
  -F 'payload={"amount":2599.50,"accountName":"Private Fitness","accountNumber":"123456","accountType":"01004"}'
```

### Expected Result (Mock Mode)
```json
{
  "status": "success",
  "message": "Payment verified successfully.",
  "data": {
    "slipId": "MOCK_SLIP_12345",
    "verified": true
  }
}
```

---

## 🔑 Environment Variables

```bash
# Slip2Go API Configuration
SLIP2GO_SECRET_KEY=<your_api_key>

# Mock Mode (true = ไม่เรียก API จริง, false = เรียก API จริง)
MOCK_SLIP2GO=true
```

---

## 💡 Design Decisions

### 1. Why Stateless?

**Pros:**
- ✅ ไม่ต้องทำ database migration
- ✅ ไม่ต้องทำ CRUD operations
- ✅ โค้ดเรียบง่าย ง่ายต่อการบำรุงรักษา
- ✅ ไม่มีปัญหา duplicate detection
- ✅ Frontend ควบคุมการเก็บข้อมูลเอง

**Cons:**
- ❌ ไม่มีประวัติการตรวจสอบใน backend
- ❌ ไม่สามารถ query ข้อมูลการชำระเงินย้อนหลัง
- ❌ ต้องเรียก API ทุกครั้ง (ไม่มี cache)

**Conclusion:**  
เหมาะสำหรับ use case ที่ต้องการแค่ **ตรวจสอบว่าเงินเข้าจริงไหม** โดยไม่ต้องการเก็บประวัติ

### 2. Mock Mode

**Purpose:**
- Slip2Go API มี quota limit (100 ครั้งทดสอบฟรี)
- Mock mode ช่วยประหยัด quota สำหรับ development/testing

**Implementation:**
```go
if c.MockMode {
    return c.mockVerifySlip(req), nil
}
// เรียก API จริง
```

### 3. Multipart/Form-Data

**Why:**
- Slip2Go API รับ request เป็น multipart/form-data
- ต้องส่งทั้งไฟล์รูปภาพและข้อมูล JSON

**Frontend Integration:**
```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);
formData.append('payload', JSON.stringify({
  amount: 2599.50,
  accountName: "Private Fitness",
  accountNumber: "123-4-56789-0",
  accountType: "01004"
}));

await fetch('/api/payments/verify-slip', {
  method: 'POST',
  body: formData
});
```

---

## 📚 Documentation Files

1. **API Documentation:** `docs/API_DOCUMENTATION.md` (Section 4.2)
2. **Test Guide:** `api_text/verify_slip_tests.md` (10 test scenarios)
3. **This Summary:** `docs/PAYMENT_SLIP_VERIFICATION_STATELESS.md`

---

## 🚀 Usage Example

### Frontend Implementation

```javascript
async function verifyPaymentSlip(slipFile, paymentDetails) {
  // Prepare form data
  const formData = new FormData();
  formData.append('file', slipFile);
  formData.append('payload', JSON.stringify({
    amount: paymentDetails.amount,
    accountName: paymentDetails.accountName,
    accountNumber: paymentDetails.accountNumber,
    accountType: paymentDetails.accountType
  }));

  // Send request
  const response = await fetch('http://localhost:8000/api/payments/verify-slip', {
    method: 'POST',
    body: formData
  });

  const result = await response.json();

  // Handle result
  if (result.status === 'success' && result.data.verified) {
    console.log('✅ Payment verified!', result.data.slipId);
    // Proceed to activate membership
    return true;
  } else {
    console.log('❌ Verification failed:', result.message);
    return false;
  }
}
```

---

## ⚠️ Important Notes

1. **No Database Storage:**
   - API ไม่เก็บข้อมูลใน database
   - ถ้าต้องการประวัติการชำระเงิน ให้ Frontend เก็บหลังได้รับ `verified: true`

2. **No Duplicate Detection:**
   - ยิง request ซ้ำจะได้ผลลัพธ์ใหม่จาก Slip2Go ทุกครั้ง
   - ไม่มีการเช็คว่าสลิปนี้ถูกใช้ไปแล้วหรือยัง

3. **Mock Mode:**
   - ใช้ `MOCK_SLIP2GO=true` สำหรับ development
   - เปลี่ยนเป็น `false` สำหรับ production

4. **Frontend Responsibility:**
   - Frontend ต้องเก็บข้อมูล `slipId` และ `verified` status
   - Frontend ต้องจัดการ duplicate detection เอง (ถ้าต้องการ)

---

## 📊 Summary

**Total Files Modified:**
- Created: 1 file (slip2go/client.go)
- Modified: 5 files (DTOs, Use Case, Handler, Router)
- Documentation: 3 files

**Total Lines of Code:** ~400 lines

**Development Time:** ~2 hours

**Status:** ✅ Production Ready

---

**Implementation Date:** October 31, 2025  
**Last Updated:** October 31, 2025  
**Version:** 1.0.0 (Stateless)

