package responses

import "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"

type ListCustomerLogsResponse struct {
	Data []dbmodel.ListCustomerLogsRow `json:"data"`
	Meta PageMeta                      `json:"meta"`
}

type CustomerLogUpdatedResponse struct {
	Message string `json:"message"`
}

type CustomerLogDeletedResponse struct {
	Message string `json:"message"`
}