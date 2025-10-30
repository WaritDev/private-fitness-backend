package utils

import (
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

func StrPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func EnumMarital(p *string) dbmodel.CustomersMaritalStatus {
	if p == nil {
		return dbmodel.CustomersMaritalStatus("")
	}
	return dbmodel.CustomersMaritalStatus(*p)
}