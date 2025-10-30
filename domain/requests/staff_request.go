package requests

type ListStaffsRequest struct {
	Page  int32 `query:"page"`
	Limit int32 `query:"limit"`
}

type CreateStaffRequest struct {
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	ConfirmPass   string  `json:"confirmPassword"`
	Role          string  `json:"role"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	Gender        string  `json:"gender"`
	DateOfBirth   string  `json:"dateOfBirth"`
	PhoneNumber   string  `json:"phoneNumber"`
	Gmail         string  `json:"gmail"`
	Specialty     *string `json:"specialty,omitempty"`
}

type UpdateStaffRequest struct {
	NewPassword        *string `json:"newPassword,omitempty"`
	ConfirmNewPassword *string `json:"confirmNewPassword,omitempty"`
	Role               string  `json:"role"`
	FirstName          string  `json:"firstName"`
	LastName           string  `json:"lastName"`
	Gender             string  `json:"gender"`
	DateOfBirth        string  `json:"dateOfBirth"`
	PhoneNumber        string  `json:"phoneNumber"`
	Gmail              string  `json:"gmail"`
	Specialty          *string `json:"specialty,omitempty"`
	IsActive           bool    `json:"isActive"`
}