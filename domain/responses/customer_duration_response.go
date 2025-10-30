package responses

import "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"

type ListCustomerDurationsResponse struct {
	Data []dbmodel.ListCustomerDurationsRow `json:"data"`
	Meta PageMeta                           `json:"meta"`
}