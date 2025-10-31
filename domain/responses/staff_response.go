package responses

import "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"

type ListStaffsResponse struct {
	Data []dbmodel.ListStaffsRow `json:"data"`
	Meta PageMeta                `json:"meta"`
}

type PageMeta struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int32 `json:"total_pages"`
}

type StaffCreatedResponse struct {
	Message string `json:"message"`
}

type StaffUpdatedResponse struct {
	Message string `json:"message"`
}

type StaffDeletedResponse struct {
	Message string `json:"message"`
}

type NullString struct {
    String string `json:"String"`
    Valid  bool   `json:"Valid"`
}

type NullBool struct {
    Bool  bool `json:"Bool"`
    Valid bool `json:"Valid"`
}

type Staff struct {
    Username    string     `json:"username"`
    Role        string     `json:"role"`
    FirstName   string     `json:"firstName"`
    LastName    string     `json:"lastName"`
    Gender      string     `json:"gender"`      // MALE|FEMALE|OTHER
    DateOfBirth string     `json:"dateOfBirth"` // YYYY-MM-DD
    PhoneNumber string     `json:"phoneNumber"`
    Gmail       string     `json:"gmail"`
    Specialty   NullString `json:"specialty"`
    IsActive    NullBool   `json:"isActive"`
}