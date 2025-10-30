package requests

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	// User info
	Username    string `json:"username" validate:"required,min=4,max=30,alphanum"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"firstName" validate:"required,min=1"`
	LastName    string `json:"lastName" validate:"required,min=1"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"dateOfBirth"`
	PhoneNumber string `json:"phone" validate:"required,min=1"`
	Email       string `json:"email" validate:"required,email"`

	// Customer-specific info
	MarketingSource              string `json:"marketingSource"`
	EmergencyContactPhone        string `json:"emergencyContactPhone"`
	EmergencyContactRelationship string `json:"emergencyContactRelationship"`
	EmergencyContactName         string `json:"emergencyContactName"`
	MaritalStatus                string `json:"maritalStatus"`
	CompanyPosition              string `json:"companyPosition"`
	CompanyName                  string `json:"companyName"`
	Address                      string `json:"address"`
	HealthInfo                   string `json:"healthInfo"`
}
