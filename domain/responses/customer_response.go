package responses

import "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"

type ListCustomersResponse struct {
	Data []dbmodel.ListCustomersRow `json:"data"`
	Meta PageMeta                   `json:"meta"`
}

type CustomerUpdatedResponse struct {
	Message string `json:"message"`
}

type CustomerDeletedResponse struct {
	Message string `json:"message"`
}
type Customer struct {
	Username                   string     `json:"username"`
	FirstName                  string     `json:"firstName"`
	LastName                   string     `json:"lastName"`
	Gender                     string     `json:"gender"`
	DateOfBirth                string     `json:"dateOfBirth"`
	PhoneNumber                string     `json:"phoneNumber"`
	Gmail                      string     `json:"gmail"`
	IsActive                   NullBool   `json:"isActive"`
	HealthInfo                 NullString `json:"healthInfo"`
	Address                    NullString `json:"address"`
	CompanyName                NullString `json:"companyName"`
	CompanyPosition            NullString `json:"companyPosition"`
	MaritalStatus              NullString `json:"maritalStatus"`
	EmergencyContactName       NullString `json:"emergencyContactName"`
	EmergencyContactRelationship NullString `json:"emergencyContactRelationship"`
	EmergencyContactPhone      NullString `json:"emergencyContactPhone"`
	MarketingSource            NullString `json:"marketingSource"`
}