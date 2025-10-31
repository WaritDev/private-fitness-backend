package requests

type ListCustomerDurationsRequest struct {
	Page  int32 `query:"page"`
	Limit int32 `query:"limit"`
}

type UpdateCustomerDurationRequest struct {
	StartDate      string  `json:"startDate"`      // YYYY-MM-DD
	PricePaid      float64 `json:"pricePaid"`      // >= 0
	DiscountAmount float64 `json:"discountAmount"` // >= 0
	Status         string  `json:"status"`         // ACTIVE|EXPIRED|FROZEN|CANCELLED
}

// RenewDurationRequest - Customer Self-Purchase Duration Package
// StartDate is auto-calculated as NOW() in SQL
type RenewDurationRequest struct {
	ProductID int32 `json:"productId" validate:"required,gt=0"`
}

// RegisterCustomerDurationRequest - Use Case 2.1C: ลงทะเบียนผู้ใช้งานสำหรับแพ็กเกจ Duration
// Combines data from Use Case 3S (customer info) and Use Case 2.1C (account creation + duration package)
type RegisterCustomerDurationRequest struct {
	// Account Creation (Use Case 2.1C)
	Username        string `json:"username" validate:"required,min=4,max=30,alphanum"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`

	// Customer Information (Use Case 3S)
	FirstName                    string `json:"firstName" validate:"required"`
	LastName                     string `json:"lastName" validate:"required"`
	Gender                       string `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	DateOfBirth                  string `json:"dateOfBirth" validate:"required"` // YYYY-MM-DD format
	PhoneNumber                  string `json:"phoneNumber" validate:"required"`
	Gmail                        string `json:"gmail" validate:"required,email"`
	HealthInfo                   string `json:"healthInfo"`
	Address                      string `json:"address"`
	CompanyName                  string `json:"companyName"`
	CompanyPosition              string `json:"companyPosition"`
	MaritalStatus                string `json:"maritalStatus" validate:"required,oneof=SINGLE MARRIED DIVORCED WIDOWED"`
	EmergencyContactName         string `json:"emergencyContactName" validate:"required"`
	EmergencyContactRelationship string `json:"emergencyContactRelationship" validate:"required"`
	EmergencyContactPhone        string `json:"emergencyContactPhone" validate:"required"`
	MarketingSource              string `json:"marketingSource"`

	// Duration Package Details (Use Case 2.1C)
	ProductID      int32   `json:"productId" validate:"required,gt=0"`
	SalesUsername  string  `json:"salesUsername" validate:"required"`
	StartDate      string  `json:"startDate" validate:"required"` // YYYY-MM-DD format
	DurationDays   int32   `json:"durationDays" validate:"required,gt=0"`
	PricePaid      float64 `json:"pricePaid" validate:"required,gte=0"`
	DiscountAmount float64 `json:"discountAmount" validate:"gte=0"`
}
