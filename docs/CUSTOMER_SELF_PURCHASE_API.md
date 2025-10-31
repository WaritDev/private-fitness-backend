# Customer Self-Purchase API - Frontend Implementation Guide

## Overview
ลูกค้าสามารถซื้อแพ็กเกจเพิ่มเองในระบบได้ 2 ประเภท:
1. **Duration Package** (ต่ออายุสมาชิก) - เพิ่มวันใช้งาน
2. **Session Package** (เพิ่มครั้งออกกำลังกาย) - เพิ่มจำนวนครั้งเทรนกับเทรนเนอร์

---

## 🎯 Business Logic

### การซื้อแพ็กเกจเพิ่มเอง (Self-Purchase)

#### **ความแตกต่างจากการซื้อผ่าน Sales:**

| Feature | Sales Registration | Customer Self-Purchase |
|---------|-------------------|------------------------|
| **sales_username** | มีค่า (เช่น "sales01") | **NULL** (ซื้อเอง) |
| **discount_amount** | มีส่วนลด (ตาม Sales) | **0** (ไม่มีส่วนลด) |
| **price_paid** | อาจต่ำกว่า list_price | **= list_price** (ราคาเต็ม) |
| **Use Case** | 3S + 4S (Sales ลงทะเบียน) | Customer ซื้อเพิ่มเอง |

#### **Frontend Behavior:**

1. **หน้าแสดงแพ็กเกจ:**
   - แสดงเฉพาะแพ็กเกจที่ `status = 'ACTIVE'`
   - Duration: แสดง `daysRemaining` (วันคงเหลือ)
   - Session: แสดง `totalSessions - usedSessions` (ครั้งคงเหลือ)

2. **หน้าซื้อเพิ่ม:**
   - ไม่ต้องเลือก Start Date (Backend คำนวณเป็น `NOW()` อัตโนมัติ)
   - Duration: Start Date = Today, End Date = Today + duration_days
   - Session: ไม่มี end_date (ใช้ได้เรื่อยๆ จนกว่าจะหมด)

3. **การคำนวณวันที่:**
   - **Duration**: 
     - `start_date` = `NOW()` (TIMESTAMP, วันเวลาปัจจุบัน)
     - `end_date` = `NOW() + duration_days` (TIMESTAMP, auto-calculated by SQL)
     - Backend ดึง `duration_days` จาก `products` table อัตโนมัติ
   - **Session**: ไม่มี end_date (ใช้ได้เรื่อยๆ จนกว่าจะหมดครั้ง)

---

## 📦 API 1: Customer Self-Purchase Duration Package

### Endpoint
```
POST /api/customer-durations/renew
```

### Authentication
**Required**: JWT Token (Customer only)

### Use Case
ลูกค้าต่ออายุแพ็กเกจ Duration เพื่อขยายวันใช้งาน

### Request Body
```json
{
  "productId": 1
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `productId` | integer | Yes | รหัสสินค้าประเภท DURATION |

**Note**: `startDate` is auto-calculated by backend as `NOW()` (TIMESTAMP)

### Request Example
```javascript
async function renewDurationPackage(productId) {
  try {
    const response = await fetch('/api/customer-durations/renew', {
      method: 'POST',
      credentials: 'include', // ส่ง JWT cookie
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        productId: productId
        // startDate is auto-calculated by backend as NOW()
      })
    });

    const data = await response.json();

    if (data.status === 'success') {
      console.log('Duration package renewed:', data.result);
      // {
      //   id: 15,
      //   customerUsername: "cust01",
      //   productId: 1,
      //   productName: "Duration 30 days",
      //   durationDays: 30,
      //   salesUsername: null,        // ← ซื้อเอง
      //   purchaseDate: "2025-10-31T...",
      //   startDate: "2025-11-01T...",
      //   endDate: "2025-12-01T...",  // ← startDate + 30 วัน
      //   daysRemaining: 30,
      //   pricePaid: 2500.00,         // ← ราคาเต็ม (list_price)
      //   discountAmount: 0,          // ← ไม่มีส่วนลด
      //   status: "ACTIVE",
      //   message: "Duration package renewed successfully"
      // }
      
      // Redirect or refresh package list
      window.location.href = '/my-packages';
    } else {
      alert(data.message);
    }
  } catch (error) {
    console.error('Failed to renew duration:', error);
    alert('Failed to renew duration package');
  }
}
```

### Response (Success - 201 Created)
```json
{
  "status": "success",
  "status_code": 201,
  "message": "Duration package renewed successfully",
  "result": {
    "id": 15,
    "customerUsername": "cust01",
    "productId": 1,
    "productName": "Duration 30 days",
    "durationDays": 30,
    "salesUsername": null,
    "purchaseDate": "2025-10-31T14:30:00+07:00",
    "startDate": "2025-11-01T00:00:00+07:00",
    "endDate": "2025-12-01T00:00:00+07:00",
    "daysRemaining": 30,
    "pricePaid": 2500.00,
    "discountAmount": 0,
    "status": "ACTIVE",
    "message": "Duration package renewed successfully"
  }
}
```

### Response (Error - 400 Bad Request)
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Product not found or not a DURATION type",
  "result": null
}
```

### Response (Error - 401 Unauthorized)
```json
{
  "status": "error",
  "status_code": 401,
  "message": "Unauthorized: Please login first",
  "result": null
}
```

### Business Rules

1. **Product Validation:**
   - Product ต้องมีอยู่จริงและ `type = 'DURATION'`
   - Product ต้อง `status = 'ACTIVE'`

2. **Pricing:**
   - `price_paid` = `list_price` (ราคาเต็ม)
   - `discount_amount` = 0 (ไม่มีส่วนลด)

3. **Sales Info:**
   - `sales_username` = **NULL** (บ่งบอกว่าลูกค้าซื้อเอง)

4. **Date Calculation:**
   - `purchase_date` = NOW() (TIMESTAMP)
   - `start_date` = NOW() (TIMESTAMP, auto-calculated by backend)
   - `end_date` = `NOW() + duration_days` (TIMESTAMP, auto-calculated by backend)

5. **Status:**
   - `status` = 'ACTIVE' (พร้อมใช้งานทันที)

### SQL Query (Backend Implementation)
```sql
-- RenewCustomerDuration (INSERT...SELECT pattern)
INSERT INTO customer_durations (
  customer_username,
  product_id,
  sales_username,      -- NULL for self-purchase
  purchase_date,       -- NOW()
  start_date,          -- NOW() (auto-calculated)
  end_date,            -- NOW() + duration_days (auto-calculated from products table)
  price_paid,
  discount_amount,     -- 0 for self-purchase
  status               -- 'ACTIVE'
)
SELECT 
  ?,                   -- customer_username (from JWT)
  ?,                   -- product_id
  NULL,                -- sales_username (self-purchase)
  NOW(),               -- purchase_date
  NOW(),               -- start_date (TIMESTAMP)
  DATE_ADD(NOW(), INTERVAL p.duration_days DAY), -- end_date (TIMESTAMP)
  ?,                   -- price_paid (list_price)
  '0.00',              -- discount_amount
  'ACTIVE'             -- status
FROM products p
WHERE p.id = ?         -- product_id (validate + retrieve duration_days)
  AND p.type = 'DURATION'
  AND p.is_active = 1
  AND p.duration_days IS NOT NULL
LIMIT 1;
```

**Key Features:**
- ✅ Uses `INSERT...SELECT` pattern for atomic validation + insertion
- ✅ Auto-calculates `start_date` as `NOW()` (TIMESTAMP)
- ✅ Auto-calculates `end_date` from `products.duration_days`
- ✅ Validates product in same query (type, is_active, duration_days)
- ✅ Single atomic operation (no separate SELECT needed)

---

## 🎯 API 2: Customer Self-Purchase Session Package

### Endpoint
```
POST /api/customer-sessions/renew
```

### Authentication
**Required**: JWT Token (Customer only)

### Use Case
ลูกค้าซื้อแพ็กเกจ Session เพิ่มเพื่อออกกำลังกายกับเทรนเนอร์

### Request Body
```json
{
  "productId": 4,
  "trainerUsername": "trainer1"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `productId` | integer | Yes | รหัสสินค้าประเภท SESSION |
| `trainerUsername` | string | Yes | Username ของเทรนเนอร์ที่เลือก |

### Request Example
```javascript
async function renewSessionPackage(productId, trainerUsername) {
  try {
    const response = await fetch('/api/customer-sessions/renew', {
      method: 'POST',
      credentials: 'include', // ส่ง JWT cookie
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        productId: productId,
        trainerUsername: trainerUsername
      })
    });

    const data = await response.json();

    if (data.status === 'success') {
      console.log('Session package renewed:', data.result);
      // {
      //   id: 28,
      //   customerUsername: "cust01",
      //   trainerUsername: "trainer1",
      //   productId: 4,
      //   productName: "Session 10 times - Economic",
      //   totalSessions: 10,
      //   usedSessions: 0,            // ← เริ่มต้น 0
      //   purchaseDate: "2025-10-31T...",
      //   pricePaid: 5000.00,         // ← ราคาเต็ม (list_price)
      //   discountAmount: 0,          // ← ไม่มีส่วนลด
      //   status: "ACTIVE",
      //   message: "Session package renewed successfully"
      // }
      
      // Redirect or refresh package list
      window.location.href = '/my-packages';
    } else {
      alert(data.message);
    }
  } catch (error) {
    console.error('Failed to renew session:', error);
    alert('Failed to renew session package');
  }
}
```

### Response (Success - 201 Created)
```json
{
  "status": "success",
  "status_code": 201,
  "message": "Session package renewed successfully",
  "result": {
    "id": 28,
    "customerUsername": "cust01",
    "trainerUsername": "trainer1",
    "productId": 4,
    "productName": "Session 10 times - Economic",
    "totalSessions": 10,
    "usedSessions": 0,
    "purchaseDate": "2025-10-31T14:30:00+07:00",
    "pricePaid": 5000.00,
    "discountAmount": 0,
    "status": "ACTIVE",
    "message": "Session package renewed successfully"
  }
}
```

### Response (Error - 400 Bad Request)
```json
{
  "status": "error",
  "status_code": 400,
  "message": "Product not found or not a SESSION type",
  "result": null
}
```

```json
{
  "status": "error",
  "status_code": 400,
  "message": "Trainer not found",
  "result": null
}
```

### Response (Error - 401 Unauthorized)
```json
{
  "status": "error",
  "status_code": 401,
  "message": "Unauthorized: Please login first",
  "result": null
}
```

### Business Rules

1. **Product Validation:**
   - Product ต้องมีอยู่จริงและ `type = 'SESSION'`
   - Product ต้อง `status = 'ACTIVE'`
   - Product ต้องมี `session_amount > 0`

2. **Trainer Validation:**
   - Trainer ต้องมีอยู่จริง
   - Trainer ต้อง `status = 'ACTIVE'`

3. **Pricing:**
   - `price_paid` = `list_price` (ราคาเต็ม)
   - `discount_amount` = 0 (ไม่มีส่วนลด)

4. **Sales Info:**
   - `sales_username` = **NULL** (บ่งบอกว่าลูกค้าซื้อเอง)

5. **Session Info:**
   - `total_sessions` = `session_amount` (จำนวนครั้งตาม product)
   - `used_sessions` = 0 (เริ่มต้นยังไม่ได้ใช้)

6. **Status:**
   - `status` = 'ACTIVE' (พร้อมใช้งานทันที)

### SQL Query (Backend Implementation)
```sql
-- RenewCustomerSession (ทำไปแล้วใน code)
INSERT INTO customer_sessions (
  customer_username,
  trainer_username,
  product_id,
  sales_username,      -- NULL for self-purchase
  purchase_date,       -- NOW()
  total_sessions,      -- From product.session_amount
  used_sessions,       -- 0 for new package
  price_paid,
  discount_amount,     -- 0 for self-purchase
  status               -- 'ACTIVE'
) VALUES (
  ?, -- customerUsername (from JWT)
  ?, -- trainerUsername
  ?, -- productId
  NULL,
  NOW(),
  ?, -- totalSessions
  0,
  ?, -- list_price
  0,
  'ACTIVE'
);
```

---

## 🖥️ Frontend Implementation

### 1. Display Active Packages Page

```javascript
// หน้าแสดงแพ็กเกจทั้งหมด
async function loadMyPackages() {
  try {
    // ดึงแพ็กเกจ Duration
    const durationsResponse = await fetch('/api/member/check-duration-permission', {
      credentials: 'include'
    });
    const durationsData = await durationsResponse.json();

    // ดึงแพ็กเกจ Session
    const sessionsResponse = await fetch('/api/member/check-session-permission', {
      credentials: 'include'
    });
    const sessionsData = await sessionsResponse.json();

    // Render Duration Packages
    renderDurationPackages(durationsData.activeDurations);

    // Render Session Packages
    renderSessionPackages(durationsData.activeSessions);

  } catch (error) {
    console.error('Failed to load packages:', error);
  }
}

function renderDurationPackages(durations) {
  const container = document.getElementById('duration-packages');
  
  if (durations.length === 0) {
    container.innerHTML = '<p>No active duration packages</p>';
    return;
  }

  let html = '<h3>Duration Packages</h3><div class="packages-grid">';

  durations.forEach(duration => {
    html += `
      <div class="package-card">
        <h4>${duration.productName}</h4>
        <p>Days Remaining: <strong>${duration.daysRemaining}</strong> / ${duration.durationDays} days</p>
        <p>Start: ${formatDate(duration.startDate)}</p>
        <p>End: ${formatDate(duration.endDate)}</p>
        <p>Status: <span class="badge badge-${duration.status.toLowerCase()}">${duration.status}</span></p>
        ${duration.salesUsername ? 
          `<p class="text-muted">Purchased via Sales: ${duration.salesUsername}</p>` : 
          `<p class="text-success">Self-Purchased</p>`
        }
      </div>
    `;
  });

  html += '</div>';
  container.innerHTML = html;
}

function renderSessionPackages(sessions) {
  const container = document.getElementById('session-packages');
  
  if (sessions.length === 0) {
    container.innerHTML = '<p>No active session packages</p>';
    return;
  }

  let html = '<h3>Session Packages</h3><div class="packages-grid">';

  sessions.forEach(session => {
    const remainingSessions = session.totalSessions - session.usedSessions;
    
    html += `
      <div class="package-card">
        <h4>${session.productName}</h4>
        <p>Trainer: <strong>${session.trainerUsername}</strong></p>
        <p>Remaining: <strong>${remainingSessions}</strong> / ${session.totalSessions} sessions</p>
        <p>Purchase: ${formatDate(session.purchaseDate)}</p>
        <p>Status: <span class="badge badge-${session.status.toLowerCase()}">${session.status}</span></p>
        ${session.salesUsername ? 
          `<p class="text-muted">Purchased via Sales: ${session.salesUsername}</p>` : 
          `<p class="text-success">Self-Purchased</p>`
        }
      </div>
    `;
  });

  html += '</div>';
  container.innerHTML = html;
}

function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}
```

### 2. Purchase Duration Package

```javascript
// หน้าซื้อแพ็กเกจ Duration
async function loadDurationProducts() {
  try {
    const response = await fetch('/api/products/durations', {
      credentials: 'include'
    });
    const data = await response.json();

    renderDurationProducts(data.result);
  } catch (error) {
    console.error('Failed to load duration products:', error);
  }
}

function renderDurationProducts(products) {
  const container = document.getElementById('duration-products');
  
  let html = '<h3>Select Duration Package</h3><div class="products-grid">';

  products.forEach(product => {
    // Convert price from cents to THB
    const priceInBaht = product.listPrice / 100;

    html += `
      <div class="product-card">
        <h4>${product.productName}</h4>
        <p>Duration: <strong>${product.durationDays}</strong> days</p>
        <p>Price: <strong>฿${priceInBaht.toLocaleString('th-TH', {
          minimumFractionDigits: 2,
          maximumFractionDigits: 2
        })}</strong></p>
        <button class="btn btn-primary" onclick="buyDuration(${product.id}, '${product.productName}', ${priceInBaht})">
          Buy Now
        </button>
      </div>
    `;
  });

  html += '</div>';
  container.innerHTML = html;
}

function buyDuration(productId, productName, price) {
  // Confirm purchase (no startDate needed - backend auto-calculates)
  if (!confirm(`Confirm purchase?\n\nPackage: ${productName}\nPrice: ฿${price}\n\nStart Date: Today (auto-calculated)\nEnd Date: Today + ${durationDays} days`)) {
    return;
  }

  // Call API (simplified - no startDate parameter)
  renewDurationPackage(productId);
}
```

### 3. Purchase Session Package

```javascript
// หน้าซื้อแพ็กเกจ Session
async function loadSessionProducts() {
  try {
    // ดึงรายการ Products
    const productsResponse = await fetch('/api/products/sessions', {
      credentials: 'include'
    });
    const productsData = await productsResponse.json();

    // ดึงรายการ Trainers
    const trainersResponse = await fetch('/api/trainers', {
      credentials: 'include'
    });
    const trainersData = await trainersResponse.json();

    renderSessionPurchaseForm(productsData.result, trainersData.result);
  } catch (error) {
    console.error('Failed to load session products:', error);
  }
}

function renderSessionPurchaseForm(products, trainers) {
  const container = document.getElementById('session-purchase');
  
  let html = `
    <h3>Buy Session Package</h3>
    <form id="buy-session-form">
      <div class="form-group">
        <label for="sessionProduct">Select Package:</label>
        <select id="sessionProduct" name="productId" class="form-control" required>
          <option value="">-- Select Package --</option>
  `;

  products.forEach(product => {
    const priceInBaht = product.listPrice / 100;
    html += `
      <option value="${product.id}" data-price="${priceInBaht}" data-sessions="${product.sessionAmount}">
        ${product.productName} - ${product.sessionAmount} sessions - ฿${priceInBaht.toLocaleString('th-TH')}
      </option>
    `;
  });

  html += `
        </select>
      </div>

      <div class="form-group">
        <label for="sessionTrainer">Select Trainer:</label>
        <select id="sessionTrainer" name="trainerUsername" class="form-control" required>
          <option value="">-- Select Trainer --</option>
  `;

  trainers.forEach(trainer => {
    html += `
      <option value="${trainer.username}">
        ${trainer.firstName} ${trainer.lastName} (${trainer.username})
      </option>
    `;
  });

  html += `
        </select>
      </div>

      <div id="purchase-summary" class="alert alert-info" style="display:none;">
        <h5>Purchase Summary</h5>
        <p>Package: <strong id="summary-package"></strong></p>
        <p>Sessions: <strong id="summary-sessions"></strong></p>
        <p>Price: <strong id="summary-price"></strong></p>
        <p>Trainer: <strong id="summary-trainer"></strong></p>
      </div>

      <button type="submit" class="btn btn-success">Buy Now</button>
      <button type="button" class="btn btn-secondary" onclick="window.location.href='/my-packages'">Cancel</button>
    </form>
  `;

  container.innerHTML = html;

  // Attach event handlers
  document.getElementById('sessionProduct').addEventListener('change', updatePurchaseSummary);
  document.getElementById('sessionTrainer').addEventListener('change', updatePurchaseSummary);
  document.getElementById('buy-session-form').addEventListener('submit', handleBuySession);
}

function updatePurchaseSummary() {
  const productSelect = document.getElementById('sessionProduct');
  const trainerSelect = document.getElementById('sessionTrainer');
  const summary = document.getElementById('purchase-summary');

  if (productSelect.value && trainerSelect.value) {
    const selectedOption = productSelect.options[productSelect.selectedIndex];
    const price = parseFloat(selectedOption.getAttribute('data-price'));
    const sessions = selectedOption.getAttribute('data-sessions');
    const packageName = selectedOption.text.split(' - ')[0];
    const trainerName = trainerSelect.options[trainerSelect.selectedIndex].text;

    document.getElementById('summary-package').textContent = packageName;
    document.getElementById('summary-sessions').textContent = `${sessions} sessions`;
    document.getElementById('summary-price').textContent = `฿${price.toLocaleString('th-TH', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })}`;
    document.getElementById('summary-trainer').textContent = trainerName;

    summary.style.display = 'block';
  } else {
    summary.style.display = 'none';
  }
}

async function handleBuySession(event) {
  event.preventDefault();

  const productId = parseInt(document.getElementById('sessionProduct').value);
  const trainerUsername = document.getElementById('sessionTrainer').value;

  if (!productId || !trainerUsername) {
    alert('Please select both package and trainer');
    return;
  }

  // Confirm
  if (!confirm('Confirm purchase?')) {
    return;
  }

  // Call API
  await renewSessionPackage(productId, trainerUsername);
}
```

---

## 📊 Data Flow

### Duration Purchase Flow
```
[Customer] → Select Duration Package
    ↓
[Frontend] → POST /api/customer-durations/renew { productId }
    ↓
[Backend] → INSERT...SELECT with JOIN to products
            - Validate Product (type=DURATION, is_active=1, duration_days NOT NULL)
            - Retrieve duration_days from products table
            - Calculate dates:
              • sales_username = NULL
              • discount_amount = 0
              • price_paid = list_price
              • start_date = NOW() (TIMESTAMP)
              • end_date = NOW() + duration_days (TIMESTAMP)
            - Single atomic operation
    ↓
[Frontend] → Show success + Redirect to /my-packages
    ↓
[Frontend] → Display updated active packages
```

### Session Purchase Flow
```
[Customer] → Select Session Package
    ↓
[Customer] → Select Trainer
    ↓
[Frontend] → POST /api/customer-sessions/renew
    ↓
[Backend] → Validate Product (type=SESSION, status=ACTIVE)
    ↓
[Backend] → Validate Trainer (status=ACTIVE)
    ↓
[Backend] → INSERT customer_sessions
            - sales_username = NULL
            - discount_amount = 0
            - price_paid = list_price
            - total_sessions = product.session_amount
            - used_sessions = 0
    ↓
[Frontend] → Show success + Redirect to /my-packages
    ↓
[Frontend] → Display updated active packages
```

---

## ⚠️ Error Handling

### Common Errors

```javascript
async function handlePurchaseError(error) {
  if (error.status_code === 401) {
    // Unauthorized - redirect to login
    alert('Please login first');
    window.location.href = '/login';
  } else if (error.status_code === 400) {
    // Bad request - show error message
    alert(error.message);
  } else if (error.status_code === 404) {
    // Not found
    alert('Product or trainer not found');
  } else {
    // Server error
    alert('Failed to purchase package. Please try again later.');
  }
}
```

### Validation Rules

#### Duration Package:
- ✅ Product ID must exist
- ✅ Product type must be 'DURATION'
- ✅ Product status must be 'ACTIVE'
- ✅ Product duration_days must be NOT NULL and > 0
- ✅ Customer must be logged in (JWT token)
- ❌ ~~Start date validation~~ (removed - auto-calculated by backend)

#### Session Package:
- ✅ Product ID must exist
- ✅ Product type must be 'SESSION'
- ✅ Product status must be 'ACTIVE'
- ✅ Product session_amount > 0
- ✅ Trainer username must exist
- ✅ Trainer status must be 'ACTIVE'
- ✅ Customer must be logged in

---

## 🔍 Database Queries

### Check Active Packages

```sql
-- Get Active Duration Packages
SELECT 
  cd.id,
  cd.customer_username,
  cd.product_id,
  p.product_name,
  cd.duration_days,
  cd.sales_username,
  cd.purchase_date,
  cd.start_date,
  cd.end_date,
  DATEDIFF(cd.end_date, NOW()) AS days_remaining,
  cd.price_paid,
  cd.discount_amount,
  cd.status
FROM customer_durations cd
JOIN products p ON cd.product_id = p.id
WHERE cd.customer_username = ?
  AND cd.status = 'ACTIVE'
ORDER BY cd.end_date DESC;

-- Get Active Session Packages
SELECT 
  cs.id,
  cs.customer_username,
  cs.trainer_username,
  cs.product_id,
  p.product_name,
  cs.total_sessions,
  cs.used_sessions,
  cs.purchase_date,
  cs.price_paid,
  cs.discount_amount,
  cs.status
FROM customer_sessions cs
JOIN products p ON cs.product_id = p.id
WHERE cs.customer_username = ?
  AND cs.status = 'ACTIVE'
  AND cs.used_sessions < cs.total_sessions
ORDER BY cs.purchase_date DESC;
```

---

## 🎨 UI/UX Recommendations

### 1. Package Cards
- แสดงชัดเจนว่าแพ็กเกจไหนซื้อเอง (Self-Purchased)
- แสดงราคาที่จ่ายจริง + ส่วนลด (ถ้ามี)
- Duration: แสดง Progress Bar วันคงเหลือ
- Session: แสดง Progress Bar ครั้งคงเหลือ

### 2. Purchase Flow
- แสดง Summary ก่อน Confirm
- ราคาแสดงเป็น THB (÷100)
- แสดง Terms & Conditions
- Confirmation popup ก่อนจ่ายเงิน

### 3. Success Message
- แสดง Popup หรือ Toast notification
- แสดงรายละเอียดแพ็กเกจที่ซื้อ
- Auto-redirect ไปหน้า My Packages

---

## ✅ Summary

### API Endpoints:

1. ✅ **POST `/api/customer-sessions/renew`** (พร้อมใช้งานแล้ว)
   - Customer self-purchase session package
   - sales_username = NULL
   - discount_amount = 0

2. ✅ **POST `/api/customer-durations/renew`** (พร้อมใช้งานแล้ว)
   - Customer self-purchase duration package
   - sales_username = NULL
   - discount_amount = 0

### Key Features:

- 🔒 JWT Authentication required
- 💰 Full price (no discount for self-purchase)
- 📊 INSERT new package (not UPDATE existing)
- 🎯 Frontend shows only ACTIVE packages
- ➕ Duration: แสดงวันคงเหลือ
- ➕ Session: แสดงครั้งคงเหลือ

### Frontend จะต้อง:
1. Convert price from cents to THB (÷100)
2. Display active packages with progress
3. Provide purchase form with product selection
4. Handle JWT token authentication
5. Show success/error messages

---

**Backend ทั้ง Session และ Duration Renewal API พร้อมใช้งานแล้ว!** 🚀

### ⚠️ Next Steps:

1. **Generate sqlc code:**
   ```bash
   make gen-sqlc
   ```

2. **Generate wire dependencies:**
   ```bash
   make gen-wire
   ```

3. **Build and test:**
   ```bash
   go build -o tmp/bin/server cmd/app/main.go
   ```

4. **Test APIs:**
   ```bash
   # Login first
   curl -X POST http://localhost:8000/api/auth/login \
     -H "Content-Type: application/json" \
     -c cookies.txt \
     -d '{"username": "cust01", "password": "password123"}'

   # Test Duration Renewal (simplified - no startDate)
   curl -X POST http://localhost:8000/api/customer-durations/renew \
     -b cookies.txt \
     -H "Content-Type: application/json" \
     -d '{"productId": 1}'

   # Test Session Renewal
   curl -X POST http://localhost:8000/api/customer-sessions/renew \
     -b cookies.txt \
     -H "Content-Type: application/json" \
     -d '{"productId": 4, "trainerUsername": "trainer1"}'
   ```
