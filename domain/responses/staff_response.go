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