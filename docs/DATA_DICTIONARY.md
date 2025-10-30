# 📊 Data Dictionary & TypeScript Types

> **สำหรับ Frontend Developer (Next.js + TypeScript)**  
> Database Schema → TypeScript Interface Mappings  
> Updated: October 30, 2025

---

## 📋 Table of Contents

1. [Enums & Constants](#1-enums--constants)
2. [Core Entities](#2-core-entities)
3. [API Response Types](#3-api-response-types)
4. [Request Types](#4-request-types)
5. [Utility Types](#5-utility-types)
6. [Type Guards](#6-type-guards)

---

## 1. Enums & Constants

### 1.1 User Roles

**Database:** `ENUM('ADMIN', 'TRAINER', 'SALES', 'CUSTOMER', 'MANAGER')`

```typescript
export enum UserRole {
  ADMIN = 'ADMIN',
  TRAINER = 'TRAINER',
  SALES = 'SALES',
  CUSTOMER = 'CUSTOMER',
  MANAGER = 'MANAGER'
}

// Type guard
export const isCustomer = (role: UserRole): boolean => role === UserRole.CUSTOMER;
export const isTrainer = (role: UserRole): boolean => role === UserRole.TRAINER;
```

---

### 1.2 Gender

**Database:** `ENUM('MALE', 'FEMALE', 'OTHER')`

```typescript
export enum Gender {
  MALE = 'MALE',
  FEMALE = 'FEMALE',
  OTHER = 'OTHER'
}

// Display labels (for UI)
export const GenderLabels: Record<Gender, string> = {
  [Gender.MALE]: 'ชาย',
  [Gender.FEMALE]: 'หญิง',
  [Gender.OTHER]: 'อื่นๆ'
};
```

---

### 1.3 Marital Status

**Database:** `ENUM('SINGLE', 'MARRIED', 'DIVORCED', 'WIDOWED')`

```typescript
export enum MaritalStatus {
  SINGLE = 'SINGLE',
  MARRIED = 'MARRIED',
  DIVORCED = 'DIVORCED',
  WIDOWED = 'WIDOWED'
}

// Display labels
export const MaritalStatusLabels: Record<MaritalStatus, string> = {
  [MaritalStatus.SINGLE]: 'โสด',
  [MaritalStatus.MARRIED]: 'สมรส',
  [MaritalStatus.DIVORCED]: 'หย่าร้าง',
  [MaritalStatus.WIDOWED]: 'หม้าย'
};
```

---

### 1.4 Product Type

**Database:** `ENUM('DURATION', 'SESSION')`

```typescript
export enum ProductType {
  DURATION = 'DURATION', // แพ็กเกจรายเดือน/รายปี
  SESSION = 'SESSION'    // แพ็กเกจครั้ง
}

// Display labels
export const ProductTypeLabels: Record<ProductType, string> = {
  [ProductType.DURATION]: 'แพ็กเกจรายเดือน',
  [ProductType.SESSION]: 'แพ็กเกจครั้ง'
};
```

---

### 1.5 Product Category

**Database:** `ENUM('ECONOMIC', 'BUSINESS', 'FIRST_CLASS')`

```typescript
export enum ProductCategory {
  ECONOMIC = 'ECONOMIC',
  BUSINESS = 'BUSINESS',
  FIRST_CLASS = 'FIRST_CLASS'
}

// Display labels
export const ProductCategoryLabels: Record<ProductCategory, string> = {
  [ProductCategory.ECONOMIC]: 'ประหยัด',
  [ProductCategory.BUSINESS]: 'มาตรฐาน',
  [ProductCategory.FIRST_CLASS]: 'พรีเมียม'
};

// Color mapping (for badges)
export const ProductCategoryColors: Record<ProductCategory, string> = {
  [ProductCategory.ECONOMIC]: 'bg-green-100 text-green-800',
  [ProductCategory.BUSINESS]: 'bg-blue-100 text-blue-800',
  [ProductCategory.FIRST_CLASS]: 'bg-purple-100 text-purple-800'
};
```

---

### 1.6 Package Status

**Database:** `ENUM('ACTIVE', 'EXPIRED', 'CANCELLED', 'COMPLETED')`

```typescript
export enum PackageStatus {
  ACTIVE = 'ACTIVE',       // ใช้งานได้
  EXPIRED = 'EXPIRED',     // หมดอายุ
  CANCELLED = 'CANCELLED', // ยกเลิก
  COMPLETED = 'COMPLETED'  // ใช้ครบแล้ว (สำหรับ SESSION)
}

// Display labels
export const PackageStatusLabels: Record<PackageStatus, string> = {
  [PackageStatus.ACTIVE]: 'ใช้งานได้',
  [PackageStatus.EXPIRED]: 'หมดอายุ',
  [PackageStatus.CANCELLED]: 'ยกเลิก',
  [PackageStatus.COMPLETED]: 'ใช้ครบแล้ว'
};

// Color mapping
export const PackageStatusColors: Record<PackageStatus, string> = {
  [PackageStatus.ACTIVE]: 'bg-green-100 text-green-800',
  [PackageStatus.EXPIRED]: 'bg-gray-100 text-gray-800',
  [PackageStatus.CANCELLED]: 'bg-red-100 text-red-800',
  [PackageStatus.COMPLETED]: 'bg-blue-100 text-blue-800'
};
```

---

### 1.7 Schedule Type

**Database:** `ENUM('APPOINTMENT', 'DAY_OFF')`

```typescript
export enum ScheduleType {
  APPOINTMENT = 'APPOINTMENT', // นัดหมาย
  DAY_OFF = 'DAY_OFF'         // วันหยุด
}
```

---

### 1.8 Day of Week

```typescript
export enum DayOfWeek {
  SUNDAY = 'SUNDAY',
  MONDAY = 'MONDAY',
  TUESDAY = 'TUESDAY',
  WEDNESDAY = 'WEDNESDAY',
  THURSDAY = 'THURSDAY',
  FRIDAY = 'FRIDAY',
  SATURDAY = 'SATURDAY'
}

// Display labels (Thai)
export const DayOfWeekLabels: Record<DayOfWeek, string> = {
  [DayOfWeek.SUNDAY]: 'อาทิตย์',
  [DayOfWeek.MONDAY]: 'จันทร์',
  [DayOfWeek.TUESDAY]: 'อังคาร',
  [DayOfWeek.WEDNESDAY]: 'พุธ',
  [DayOfWeek.THURSDAY]: 'พฤหัสบดี',
  [DayOfWeek.FRIDAY]: 'ศุกร์',
  [DayOfWeek.SATURDAY]: 'เสาร์'
};

// Short labels
export const DayOfWeekShortLabels: Record<DayOfWeek, string> = {
  [DayOfWeek.SUNDAY]: 'อา.',
  [DayOfWeek.MONDAY]: 'จ.',
  [DayOfWeek.TUESDAY]: 'อ.',
  [DayOfWeek.WEDNESDAY]: 'พ.',
  [DayOfWeek.THURSDAY]: 'พฤ.',
  [DayOfWeek.FRIDAY]: 'ศ.',
  [DayOfWeek.SATURDAY]: 'ส.'
};
```

---

## 2. Core Entities

### 2.1 User

**Table:** `users`

```typescript
export interface User {
  username: string;           // VARCHAR(100) PRIMARY KEY
  password?: string;          // VARCHAR(255) - ไม่ส่งมาใน API response
  role: UserRole;             // ENUM
  firstName: string;          // VARCHAR(100)
  lastName: string;           // VARCHAR(100)
  gender: Gender;             // ENUM
  dateOfBirth: string;        // DATE (ISO 8601: "1995-01-15")
  phoneNumber: string;        // VARCHAR(20)
  gmail: string;              // VARCHAR(255)
  specialty: string | null;   // VARCHAR(100) - สำหรับ TRAINER
  isActive: boolean;          // TINYINT(1)
  createdAt: string;          // TIMESTAMP (RFC3339: "2025-10-30T00:00:00Z")
  updatedAt: string;          // TIMESTAMP
}

// Display name helper
export const getUserFullName = (user: User): string => 
  `${user.firstName} ${user.lastName}`;
```

---

### 2.2 Customer

**Table:** `customers`

```typescript
export interface Customer {
  username: string;                      // VARCHAR(100) PRIMARY KEY (FK to users)
  healthInfo: string;                    // TEXT
  address: string;                       // TEXT
  companyName: string;                   // VARCHAR(200)
  companyPosition: string;               // VARCHAR(100)
  maritalStatus: MaritalStatus;          // ENUM
  emergencyContactName: string;          // VARCHAR(255)
  emergencyContactRelationship: string;  // VARCHAR(50)
  emergencyContactPhone: string;         // VARCHAR(20)
  marketingSource: string;               // VARCHAR(100)
}

// Combined type with User info
export interface CustomerWithUser extends Customer {
  user: User;
}
```

---

### 2.3 Product

**Table:** `products`

```typescript
export interface Product {
  id: number;                      // INT AUTO_INCREMENT PRIMARY KEY
  name: string;                    // TEXT
  type: ProductType;               // ENUM('DURATION', 'SESSION')
  category: ProductCategory;       // ENUM
  listPrice: number;               // DECIMAL(10,2) → number
  durationDays: number | null;     // INT (สำหรับ DURATION type)
  sessionAmount: number | null;    // INT (สำหรับ SESSION type)
  isActive: boolean;               // TINYINT(1)
  paymentAccountId: number;        // INT (FK)
  createdAt: string;               // TIMESTAMP
  updatedAt: string;               // TIMESTAMP
}

// Type-specific products
export type DurationProduct = Product & {
  type: ProductType.DURATION;
  durationDays: number;
  sessionAmount: null;
};

export type SessionProduct = Product & {
  type: ProductType.SESSION;
  durationDays: null;
  sessionAmount: number;
};

// Type guards
export const isDurationProduct = (product: Product): product is DurationProduct =>
  product.type === ProductType.DURATION;

export const isSessionProduct = (product: Product): product is SessionProduct =>
  product.type === ProductType.SESSION;
```

---

### 2.4 Customer Session (Session Package)

**Table:** `customer_sessions`

```typescript
export interface CustomerSession {
  id: number;                  // INT AUTO_INCREMENT PRIMARY KEY
  customerUsername: string;    // VARCHAR(100) (FK)
  trainerUsername: string;     // VARCHAR(100) (FK)
  salesUsername: string;       // VARCHAR(100) (FK)
  productId: number;           // INT (FK)
  purchaseDate: string;        // TIMESTAMP
  totalSessions: number;       // INT
  usedSessions: number;        // INT
  pricePaid: number;           // DECIMAL(10,2) → number
  discountAmount: number;      // DECIMAL(10,2) → number
  status: PackageStatus;       // ENUM
  createdAt: string;           // TIMESTAMP
  updatedAt: string;           // TIMESTAMP
}

// With calculated fields (from API)
export interface CustomerSessionWithDetails extends CustomerSession {
  productName: string;         // JOIN from products
  sessionsRemaining: number;   // totalSessions - usedSessions
}

// Helper
export const getSessionsRemaining = (session: CustomerSession): number =>
  session.totalSessions - session.usedSessions;
```

---

### 2.5 Customer Duration (Duration Package)

**Table:** `customer_durations`

```typescript
export interface CustomerDuration {
  id: number;                // INT AUTO_INCREMENT PRIMARY KEY
  customerUsername: string;  // VARCHAR(100) (FK)
  salesUsername: string;     // VARCHAR(100) (FK)
  productId: number;         // INT (FK)
  purchaseDate: string;      // TIMESTAMP
  startDate: string;         // TIMESTAMP
  endDate: string;           // TIMESTAMP
  pricePaid: number;         // DECIMAL(10,2) → number
  discountAmount: number;    // DECIMAL(10,2) → number
  status: PackageStatus;     // ENUM (ACTIVE, EXPIRED, CANCELLED)
  createdAt: string;         // TIMESTAMP
  updatedAt: string;         // TIMESTAMP
}

// Helper functions
export const isDurationActive = (duration: CustomerDuration): boolean => {
  const now = new Date();
  const end = new Date(duration.endDate);
  return duration.status === PackageStatus.ACTIVE && now < end;
};

export const getDaysRemaining = (duration: CustomerDuration): number => {
  const now = new Date();
  const end = new Date(duration.endDate);
  const diff = end.getTime() - now.getTime();
  return Math.ceil(diff / (1000 * 60 * 60 * 24));
};
```

---

### 2.6 Training Schedule (Appointment)

**Table:** `training_schedules`

```typescript
export interface TrainingSchedule {
  id: number;                      // INT AUTO_INCREMENT PRIMARY KEY
  trainerUsername: string;         // VARCHAR(100) (FK)
  customerUsername: string | null; // VARCHAR(100) (FK) - null for DAY_OFF
  sessionId: number | null;        // INT (FK) - null for DAY_OFF
  startTime: string;               // TIMESTAMP (RFC3339)
  endTime: string;                 // TIMESTAMP (RFC3339)
  scheduleType: ScheduleType;      // ENUM
  createdAt: string;               // TIMESTAMP
  updatedAt: string;               // TIMESTAMP
}

// Type-specific schedules
export type Appointment = TrainingSchedule & {
  scheduleType: ScheduleType.APPOINTMENT;
  customerUsername: string;
  sessionId: number;
};

export type DayOff = TrainingSchedule & {
  scheduleType: ScheduleType.DAY_OFF;
  customerUsername: null;
  sessionId: null;
};

// Type guards
export const isAppointment = (schedule: TrainingSchedule): schedule is Appointment =>
  schedule.scheduleType === ScheduleType.APPOINTMENT;

export const isDayOff = (schedule: TrainingSchedule): schedule is DayOff =>
  schedule.scheduleType === ScheduleType.DAY_OFF;

// Helper - Check if appointment is in the past
export const isAppointmentPast = (schedule: TrainingSchedule): boolean => {
  return new Date() > new Date(schedule.startTime);
};
```

---

### 2.7 Payment Account

**Table:** `payment_accounts`

```typescript
export interface PaymentAccount {
  id: number;            // INT AUTO_INCREMENT PRIMARY KEY
  accountName: string;   // VARCHAR(255)
  accountNumber: string; // VARCHAR(50)
  bankName: string;      // VARCHAR(100)
  qrCodeImageUrl: string;// TEXT
  isActive: boolean;     // TINYINT(1)
}
```

---

### 2.8 Trainer Availability

**Table:** `training_availabilities`

```typescript
export interface TrainerAvailability {
  id: number;               // INT AUTO_INCREMENT PRIMARY KEY
  trainerUsername: string;  // VARCHAR(100) (FK)
  dayOfWeek: DayOfWeek;     // ENUM
  startTime: string;        // TIME ("09:00:00")
  endTime: string;          // TIME ("17:00:00")
  createdAt: string;        // TIMESTAMP
  updatedAt: string;        // TIMESTAMP
}
```

---

## 3. API Response Types

### 3.1 Standard API Response

```typescript
export interface ApiResponse<T = any> {
  status: 'success' | 'error';
  status_code: number;
  message: string;
  result: T | null;
}

// Type alias for successful responses
export type ApiSuccessResponse<T> = ApiResponse<T> & {
  status: 'success';
  result: T;
};

// Type alias for error responses
export type ApiErrorResponse = ApiResponse<null> & {
  status: 'error';
  result: null;
};
```

---

### 3.2 Authentication Response

```typescript
// Login response
export interface LoginResponse {
  token: string;
  user: {
    sub: string;      // username
    role: UserRole;
    firstName: string;
    lastName: string;
  };
}

// Me (verify token) response
export interface MeResponse {
  authenticated: boolean;
  user?: {
    sub: string;
    role: UserRole;
    firstName: string;
    lastName: string;
  };
}

// Logout response
export interface LogoutResponse {
  ok: boolean;
}
```

---

### 3.3 Product Response

```typescript
// List products (same as Product entity)
export type ProductListResponse = Product[];

// Get product by ID (same as Product entity)
export type ProductDetailResponse = Product;
```

---

### 3.4 User Validation Response

```typescript
export interface PhoneCheckResponse {
  exists: boolean;    // true = มีคนใช้แล้ว
  available: boolean; // true = ใช้ได้ (inverse of exists)
}

export interface GmailCheckResponse {
  exists: boolean;
  available: boolean;
}
```

---

### 3.5 Payment Info Response

```typescript
export interface PaymentInfoResponse {
  productId: number;
  productName: string;
  productType: ProductType;
  productCategory: ProductCategory;
  listPrice: number;
  discountAmount: number;
  payableAmount: number;          // listPrice - discountAmount
  sessionAmount: number | null;
  durationDays: number | null;
  paymentAccountId: number;
  accountName: string;
  accountNumber: string;
  bankName: string;
  qrCodeUrl: string;
  accountActive: boolean;
}
```

---

### 3.6 Customer Registration Response

```typescript
// Duration registration
export interface CustomerDurationRegistrationResponse {
  username: string;
  durationId: number;
  productId: number;
  salesUsername: string;
  startDate: string;
  endDate: string;
  durationDays: number;
  pricePaid: string;        // DECIMAL string format
  discountAmount: string;   // DECIMAL string format
  message: string;
}

// Session registration
export interface CustomerSessionRegistrationResponse {
  username: string;
  sessionId: number;
  trainerUsername: string;
  productId: number;
  totalSessions: number;
  schedulesCreated: number;
  createdSchedules: Array<{
    scheduleId: number;
    startTime: string;
    endTime: string;
    dayOfWeek: DayOfWeek;
  }>;
  message: string;
}
```

---

### 3.7 Customer Session Response

```typescript
// Check permission
export interface CheckPermissionResponse {
  hasPermission: boolean;
  canBook: boolean;
}

// Active session packages
export interface CustomerSessionPackageResponse {
  id: number;
  customerUsername: string;
  trainerUsername: string;
  productId: number;
  productName: string;
  totalSessions: number;
  usedSessions: number;
  sessionsRemaining: number;
  purchaseDate: string;
  pricePaid: number;
  discountAmount: number;
  status: PackageStatus;
  createdAt: string;
}

export type ActiveSessionPackagesResponse = CustomerSessionPackageResponse[];
```

---

### 3.8 Booking Response

```typescript
// Get booking slots
export interface BookingSlotsResponse {
  trainerUsername: string;
  calendarStart: string;
  calendarEnd: string;
  weeklyAvailability: Array<{
    dayOfWeek: DayOfWeek;
    startTime: string;      // "09:00:00"
    endTime: string;        // "17:00:00"
  }>;
  dayOffSlots: Array<{
    startTime: string;      // RFC3339
    endTime: string;        // RFC3339
  }>;
  bookedAppointments: Array<{
    startTime: string;      // RFC3339
    endTime: string;        // RFC3339
    customerUsername: string;
  }>;
  availableSlots: any[];    // TODO: Backend ยังไม่ implement
  customerBookings: Array<{
    startTime: string;
    endTime: string;
    available: boolean;
    isBooked: boolean;
    bookedBy: string;
    slotType: string;
  }>;
  message: string;
}

// Book appointment
export interface BookAppointmentResponse {
  success: boolean;
  message: string;
  trainerUsername: string;
  customerUsername: string;
  startTime: string;
  endTime: string;
  sessionId: number;
  remainingSession: number;
}

// Cancel appointment
export interface CancelAppointmentResponse {
  success: boolean;
  message: string;
  appointmentId: number;
  customerUsername: string;
  startTime: string;
  endTime: string;
  sessionId: number;
  remainingSessions: number;
}
```

---

## 4. Request Types

### 4.1 Authentication Request

```typescript
export interface LoginRequest {
  username: string;
  password: string;
}

export interface SignupRequest {
  username: string;
  password: string;
  confirmPassword: string;
  firstName: string;
  lastName: string;
  gender: Gender;
  dateOfBirth: string;    // "YYYY-MM-DD"
  phone: string;
  gmail: string;
  // ... other fields based on role
}
```

---

### 4.2 Customer Registration Request

```typescript
// Base customer info
interface BaseCustomerRegistrationRequest {
  // User fields
  username: string;
  password: string;
  confirmPassword: string;
  firstName: string;
  lastName: string;
  gender: Gender;
  dateOfBirth: string;      // "YYYY-MM-DD"
  phone: string;
  gmail: string;
  
  // Customer fields
  healthInfo: string;
  address: string;
  companyName: string;
  companyPosition: string;
  maritalStatus: MaritalStatus;
  emergencyContactName: string;
  emergencyContactRelationship: string;
  emergencyContactPhone: string;
  marketingSource: string;
  
  // Purchase fields
  productId: number;
  salesUsername: string;
  pricePaid: number;
  discountAmount: number;
}

// Duration registration
export interface CustomerDurationRegistrationRequest extends BaseCustomerRegistrationRequest {
  startDate: string;        // "YYYY-MM-DD"
  durationDays: number;
}

// Session registration
export interface ScheduleRequest {
  startTime: string;        // RFC3339: "2025-11-05T10:00:00Z"
  endTime: string;          // RFC3339
  dayOfWeek: DayOfWeek;
}

export interface CustomerSessionRegistrationRequest extends BaseCustomerRegistrationRequest {
  trainerUsername: string;
  totalSessions: number;
  schedules: ScheduleRequest[];
}
```

---

### 4.3 Booking Request

```typescript
export interface BookAppointmentRequest {
  trainerUsername: string;
  customerUsername: string;
  sessionId: number | null;  // null = auto-find ACTIVE session
  startTime: string;         // RFC3339
  endTime: string;           // RFC3339
}

export interface CancelAppointmentRequest {
  customerUsername: string;
}
```

---

## 5. Utility Types

### 5.1 Date/Time Utilities

```typescript
// Date format helpers
export const formatDate = (date: string): string => {
  // "2025-10-30T00:00:00Z" → "30 ต.ค. 2568"
  return new Date(date).toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

export const formatDateTime = (date: string): string => {
  // "2025-10-30T14:30:00Z" → "30 ต.ค. 2568 14:30"
  return new Date(date).toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

export const formatTime = (time: string): string => {
  // "09:00:00" → "09:00"
  return time.substring(0, 5);
};

// Convert Date to RFC3339 string
export const toRFC3339 = (date: Date): string => {
  return date.toISOString();
};

// Convert Date to YYYY-MM-DD string
export const toDateString = (date: Date): string => {
  return date.toISOString().split('T')[0];
};
```

---

### 5.2 Price Utilities

```typescript
// Format price with Thai Baht symbol
export const formatPrice = (price: number): string => {
  return new Intl.NumberFormat('th-TH', {
    style: 'currency',
    currency: 'THB',
    minimumFractionDigits: 2
  }).format(price);
};

// Calculate discount percentage
export const calculateDiscountPercentage = (
  listPrice: number,
  discountAmount: number
): number => {
  return Math.round((discountAmount / listPrice) * 100);
};

// Calculate payable amount
export const calculatePayableAmount = (
  listPrice: number,
  discountAmount: number
): number => {
  return listPrice - discountAmount;
};
```

---

### 5.3 Validation Utilities

```typescript
// Phone number validation (Thai format)
export const isValidThaiPhone = (phone: string): boolean => {
  return /^0[0-9]{9}$/.test(phone);
};

// Email validation
export const isValidEmail = (email: string): boolean => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};

// Password strength check
export interface PasswordStrength {
  isValid: boolean;
  hasMinLength: boolean;    // >= 8 characters
  hasUpperCase: boolean;
  hasLowerCase: boolean;
  hasNumber: boolean;
  hasSpecialChar: boolean;
}

export const checkPasswordStrength = (password: string): PasswordStrength => {
  return {
    isValid: password.length >= 8 &&
             /[A-Z]/.test(password) &&
             /[a-z]/.test(password) &&
             /[0-9]/.test(password),
    hasMinLength: password.length >= 8,
    hasUpperCase: /[A-Z]/.test(password),
    hasLowerCase: /[a-z]/.test(password),
    hasNumber: /[0-9]/.test(password),
    hasSpecialChar: /[!@#$%^&*(),.?":{}|<>]/.test(password)
  };
};

// Username validation (4-30 characters, alphanumeric + underscore)
export const isValidUsername = (username: string): boolean => {
  return /^[a-zA-Z0-9_]{4,30}$/.test(username);
};
```

---

## 6. Type Guards

### 6.1 API Response Type Guards

```typescript
// Check if response is successful
export const isApiSuccess = <T>(
  response: ApiResponse<T>
): response is ApiSuccessResponse<T> => {
  return response.status === 'success' && response.result !== null;
};

// Check if response is error
export const isApiError = (
  response: ApiResponse<any>
): response is ApiErrorResponse => {
  return response.status === 'error';
};
```

---

### 6.2 User Type Guards

```typescript
// Check if user is authenticated
export const isAuthenticated = (meResponse: MeResponse): boolean => {
  return meResponse.authenticated === true;
};

// Check if user has specific role
export const hasRole = (user: User, role: UserRole): boolean => {
  return user.role === role;
};

// Check if user can access customer features
export const canAccessCustomerFeatures = (user: User): boolean => {
  return user.role === UserRole.CUSTOMER;
};

// Check if user can manage bookings
export const canManageBookings = (user: User): boolean => {
  return [UserRole.ADMIN, UserRole.MANAGER, UserRole.TRAINER].includes(user.role);
};
```

---

### 6.3 Package Type Guards

```typescript
// Check if package is active and has remaining sessions/days
export const isPackageUsable = (
  pkg: CustomerSession | CustomerDuration
): boolean => {
  if (pkg.status !== PackageStatus.ACTIVE) return false;
  
  if ('totalSessions' in pkg) {
    // CustomerSession
    return getSessionsRemaining(pkg) > 0;
  } else {
    // CustomerDuration
    return isDurationActive(pkg);
  }
};
```

---

## 7. Sample Usage in Next.js

### 7.1 Fetching Products

```typescript
import { ApiResponse, Product, isApiSuccess } from '@/types';

export async function getProducts(): Promise<Product[]> {
  const response = await fetch('http://localhost:8000/api/products');
  const data: ApiResponse<Product[]> = await response.json();
  
  if (isApiSuccess(data)) {
    return data.result;
  }
  
  throw new Error(data.message);
}
```

---

### 7.2 Login Flow

```typescript
import { 
  ApiResponse, 
  LoginRequest, 
  LoginResponse, 
  isApiSuccess 
} from '@/types';

export async function login(
  credentials: LoginRequest
): Promise<LoginResponse> {
  const response = await fetch('http://localhost:8000/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(credentials)
  });
  
  const data: ApiResponse<LoginResponse> = await response.json();
  
  if (isApiSuccess(data)) {
    // Store token in localStorage
    localStorage.setItem('token', data.result.token);
    localStorage.setItem('user', JSON.stringify(data.result.user));
    return data.result;
  }
  
  throw new Error(data.message);
}
```

---

### 7.3 Type-Safe State Management (with Zustand)

```typescript
import { create } from 'zustand';
import { User, Product, CustomerSessionWithDetails } from '@/types';

interface AppState {
  // User state
  user: User | null;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  
  // Products state
  products: Product[];
  setProducts: (products: Product[]) => void;
  
  // Active sessions state
  activeSessions: CustomerSessionWithDetails[];
  setActiveSessions: (sessions: CustomerSessionWithDetails[]) => void;
}

export const useAppStore = create<AppState>((set) => ({
  user: null,
  isAuthenticated: false,
  setUser: (user) => set({ user, isAuthenticated: user !== null }),
  
  products: [],
  setProducts: (products) => set({ products }),
  
  activeSessions: [],
  setActiveSessions: (activeSessions) => set({ activeSessions })
}));
```

---

### 7.4 Form Validation with Zod

```typescript
import { z } from 'zod';
import { Gender, MaritalStatus } from '@/types';
import { isValidThaiPhone, isValidEmail } from '@/utils/validation';

export const customerRegistrationSchema = z.object({
  username: z.string()
    .min(4, 'Username ต้องมีอย่างน้อย 4 ตัวอักษร')
    .max(30, 'Username ต้องไม่เกิน 30 ตัวอักษร')
    .regex(/^[a-zA-Z0-9_]+$/, 'Username ใช้ได้เฉพาะ a-z, A-Z, 0-9, _'),
  
  password: z.string()
    .min(8, 'Password ต้องมีอย่างน้อย 8 ตัวอักษร'),
  
  confirmPassword: z.string(),
  
  firstName: z.string().min(1, 'กรุณากรอกชื่อ'),
  lastName: z.string().min(1, 'กรุณากรอกนามสกุล'),
  
  gender: z.nativeEnum(Gender),
  
  dateOfBirth: z.string().refine((date) => {
    const birthDate = new Date(date);
    const today = new Date();
    const age = today.getFullYear() - birthDate.getFullYear();
    return age >= 15 && age <= 100;
  }, 'อายุต้องอยู่ระหว่าง 15-100 ปี'),
  
  phone: z.string().refine(isValidThaiPhone, 'รูปแบบเบอร์โทรไม่ถูกต้อง (เช่น 0812345678)'),
  
  gmail: z.string().refine(isValidEmail, 'รูปแบบอีเมลไม่ถูกต้อง'),
  
  healthInfo: z.string().min(1, 'กรุณากรอกข้อมูลสุขภาพ'),
  address: z.string().min(1, 'กรุณากรอกที่อยู่'),
  companyName: z.string().min(1, 'กรุณากรอกชื่อบริษัท'),
  companyPosition: z.string().min(1, 'กรุณากรอกตำแหน่ง'),
  maritalStatus: z.nativeEnum(MaritalStatus),
  
  emergencyContactName: z.string().min(1, 'กรุณากรอกชื่อผู้ติดต่อฉุกเฉิน'),
  emergencyContactRelationship: z.string().min(1, 'กรุณากรอกความสัมพันธ์'),
  emergencyContactPhone: z.string().refine(isValidThaiPhone, 'รูปแบบเบอร์โทรไม่ถูกต้อง'),
  
  marketingSource: z.string().min(1, 'กรุณาระบุแหล่งที่รู้จัก')
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Password ไม่ตรงกัน',
  path: ['confirmPassword']
});

export type CustomerRegistrationFormData = z.infer<typeof customerRegistrationSchema>;
```

---

## 8. Export Barrel

**Create:** `@/types/index.ts`

```typescript
// Enums
export * from './enums';

// Core Entities
export * from './entities/user';
export * from './entities/customer';
export * from './entities/product';
export * from './entities/customerSession';
export * from './entities/customerDuration';
export * from './entities/trainingSchedule';
export * from './entities/paymentAccount';

// API Types
export * from './api/responses';
export * from './api/requests';

// Utils
export * from './utils/validation';
export * from './utils/formatting';
export * from './utils/typeGuards';
```

---

## 📚 Additional Resources

- **API Documentation:** `/docs/API_DOCUMENTATION.md`
- **Database Schema:** `/internal/infrastructure/db/schema/`
- **Seed Data:** `/seeds/seeds.sql`

---

**สร้างโดย:** Backend Team  
**อัปเดตล่าสุด:** 30 ตุลาคม 2568  
**เวอร์ชัน:** 1.0.0

---

**Happy Coding! 🚀**
