package utils

import "database/sql"

func ToNullString(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *p, Valid: true}
}

func NS(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func NI32(v sql.NullInt32) int32 {
	if v.Valid {
		return v.Int32
	}
	return 0
}

func NB(v sql.NullBool) bool {
	if v.Valid {
		return v.Bool
	}
	return false
}