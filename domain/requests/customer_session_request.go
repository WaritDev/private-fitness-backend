package requests

import "time"

// RegisterCustomerSessionRequest - รวมข้อมูลจาก Use Case 3S + 4S + 2.2C
type RegisterCustomerSessionRequest struct {
	// จาก Use Case 2.2C (Customer creates account)
	Username        string `json:"username" validate:"required,min=4,max=30"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`

	// จาก Use Case 3S (กรอกข้อมูลสมาชิก - Customer Info)
	FirstName                    string `json:"firstName" validate:"required"`
	LastName                     string `json:"lastName" validate:"required"`
	Gender                       string `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	DateOfBirth                  string `json:"dateOfBirth" validate:"required"` // Format: YYYY-MM-DD
	PhoneNumber                  string `json:"phoneNumber" validate:"required"`
	Gmail                        string `json:"gmail" validate:"required,email"`
	HealthInfo                   string `json:"healthInfo"`
	Address                      string `json:"address" validate:"required"`
	CompanyName                  string `json:"companyName"`
	CompanyPosition              string `json:"companyPosition"`
	MaritalStatus                string `json:"maritalStatus" validate:"required,oneof=SINGLE MARRIED DIVORCED WIDOWED"`
	EmergencyContactName         string `json:"emergencyContactName" validate:"required"`
	EmergencyContactRelationship string `json:"emergencyContactRelationship" validate:"required"`
	EmergencyContactPhone        string `json:"emergencyContactPhone" validate:"required"`
	MarketingSource              string `json:"marketingSource"`

	// จาก Use Case 4S (กรอกข้อมูลสมัครคอร์ส Sessions - Session Details)
	TrainerUsername string  `json:"trainerUsername" validate:"required"`
	ProductID       int32   `json:"productId" validate:"required"`
	SalesUsername   string  `json:"salesUsername" validate:"required"`
	PricePaid       float64 `json:"pricePaid" validate:"required,min=0"`
	DiscountAmount  float64 `json:"discountAmount" validate:"min=0"`

	// Session schedules for first week (จาก 4S - days & times selected)
	SessionSchedules []SessionScheduleInput `json:"sessionSchedules" validate:"required,min=1"`
}

// SessionScheduleInput - ข้อมูลนัดหมายแต่ละครั้ง
type SessionScheduleInput struct {
	DayOfWeek string    `json:"dayOfWeek" validate:"required,oneof=MONDAY TUESDAY WEDNESDAY THURSDAY FRIDAY SATURDAY SUNDAY"`
	StartTime time.Time `json:"startTime" validate:"required"` // ISO 8601
}
