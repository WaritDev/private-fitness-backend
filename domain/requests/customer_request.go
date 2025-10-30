package requests

type ListCustomersRequest struct {
	Page  int32 `query:"page"`
	Limit int32 `query:"limit"`
}

type UpdateCustomerRequest struct {
	NewPassword        *string `json:"newPassword,omitempty"`
	ConfirmNewPassword *string `json:"confirmNewPassword,omitempty"`

	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Gender      string  `json:"gender"`       // MALE|FEMALE|OTHER
	DateOfBirth string  `json:"dateOfBirth"`  // YYYY-MM-DD
	PhoneNumber string  `json:"phoneNumber"`
	Gmail       string  `json:"gmail"`
	IsActive    bool    `json:"isActive"`

	HealthInfo                 *string `json:"healthInfo,omitempty"`
	Address                    *string `json:"address,omitempty"`
	CompanyName                *string `json:"companyName,omitempty"`
	CompanyPosition            *string `json:"companyPosition,omitempty"`
	MaritalStatus              *string `json:"maritalStatus,omitempty"`
	EmergencyContactName       *string `json:"emergencyContactName,omitempty"`
	EmergencyContactRelationship *string `json:"emergencyContactRelationship,omitempty"`
	EmergencyContactPhone      *string `json:"emergencyContactPhone,omitempty"`
	MarketingSource            *string `json:"marketingSource,omitempty"`
}