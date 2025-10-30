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