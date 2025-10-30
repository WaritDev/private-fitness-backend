package utils

import (
	"database/sql"
	"fmt"

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

func CoalesceTrueBool(p *bool) sql.NullBool {
	if p == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *p, Valid: true}
}

func CoalesceBool(b *bool) bool {
    if b == nil {
        return true
    }
    return *b
}


func NullInt32FromPtr(p *int32) sql.NullInt32 {
	if p == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *p, Valid: true}
}

func Decimal2(f float64) string {
	return fmt.Sprintf("%.2f", f)
}