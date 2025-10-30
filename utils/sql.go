package utils

import (
	"database/sql"
	"time"
)

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

func NT(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}