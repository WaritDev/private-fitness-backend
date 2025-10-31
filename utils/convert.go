package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

// ToYMD: รับได้ทั้ง time.Time, *time.Time, sql.NullTime, *sql.NullTime → คืน "YYYY-MM-DD" หรือ "" ถ้าไม่มีค่า
func ToYMD(d any) string {
	switch v := d.(type) {
	case time.Time:
		return v.Format("2006-01-02")
	case *time.Time:
		if v != nil {
			return v.Format("2006-01-02")
		}
	case sql.NullTime:
		if v.Valid {
			return v.Time.Format("2006-01-02")
		}
	case *sql.NullTime:
		if v != nil && v.Valid {
			return v.Time.Format("2006-01-02")
		}
	}
	return ""
}

// PtrToString: คืน "" ถ้า pointer ว่าง
func PtrToString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// Atoi32: แปลง string → int32 (ตรวจ overflow)
func Atoi32(s string) (int32, error) {
	n64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(n64), nil
}

// Itoa: int32 → string
func Itoa(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

// ParseMoneyInt64: รับ string เลขฐานสิบ (เช่น DECIMAL ที่ sqlc gen เป็น string) → int64
// note: ถ้าคุณเก็บเป็นสตางค์อยู่แล้ว ควรเป็นจำนวนเต็ม (ไม่มีจุดทศนิยม)
// ถ้ามีจุดทศนิยม ให้เขียนตัวแปลงเพิ่มเองตามรูปแบบข้อมูลจริง
func ParseDecimalToInt64(s string, scale int) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	neg := false
	if s[0] == '+' || s[0] == '-' {
		neg = s[0] == '-'
		s = s[1:]
	}
	// ตัดคอมมา (กรณีมีรูปแบบ 1,234.56)
	s = strings.ReplaceAll(s, ",", "")

	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}

	if intPart == "" {
		intPart = "0"
	}
	for _, ch := range intPart {
		if ch < '0' || ch > '9' {
			return 0, errors.New("invalid money format")
		}
	}
	for _, ch := range fracPart {
		if ch < '0' || ch > '9' {
			return 0, errors.New("invalid money format")
		}
	}

	// ปรับความยาวทศนิยมให้ตรง scale
	if len(fracPart) > scale {
		// ตัดทิ้งตำแหน่งเกิน (truncate) หากอยากปัดเศษ ให้เปลี่ยน logic ตรงนี้
		fracPart = fracPart[:scale]
	}
	for len(fracPart) < scale {
		fracPart += "0"
	}

	whole, err := strconv.ParseInt(intPart, 10, 64)
	if err != nil {
		return 0, errors.New("invalid money format")
	}
	frac := int64(0)
	if scale > 0 {
		if fracPart == "" {
			fracPart = "0"
		}
		frac, err = strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return 0, errors.New("invalid money format")
		}
	}

	mul := int64(1)
	for i := 0; i < scale; i++ {
		mul *= 10
	}
	result := whole*mul + frac
	if neg {
		result = -result
	}
	return result, nil
}