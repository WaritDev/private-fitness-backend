package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type CustomerDurationRepository interface {
	Create(ctx context.Context, params CreateCustomerDurationParams) error
	GetByUsername(ctx context.Context, username string) ([]CustomerDurationInfo, error)
	GetByID(ctx context.Context, id int32) (CustomerDurationInfo, error)

	List(ctx context.Context, limit, offset int32) ([]dbmodel.ListCustomerDurationsRow, error)
	Count(ctx context.Context) (int64, error)

	GetDurationDaysForDurationID(ctx context.Context, durationID int32) (int32, error)
	UpdateEditableFields(ctx context.Context, p UpdateCustomerDurationEditableFieldsParams) error

	Delete(ctx context.Context, id int32) error
	// RegisterDuration - Use Case 2.1C: สร้าง User, Customer, CustomerDuration ในครั้งเดียว (Transaction)
	RegisterDuration(ctx context.Context, tx *sql.Tx, params RegisterDurationParams) (int32, error)

	// GetCustomerActiveDuration - ดึง Duration packages ที่ยัง ACTIVE ของลูกค้า (คล้าย GetCustomerActiveSessions)
	GetCustomerActiveDuration(ctx context.Context, username string) ([]ActiveDurationPackageInfo, error)

	// RenewDuration - Customer self-purchase duration package (ซื้อเพิ่มเอง)
	RenewDuration(ctx context.Context, params RenewDurationParams) error
}

// RenewDurationParams - Parameters for customer self-purchase duration
// StartDate and DurationDays are auto-calculated in SQL (NOW() and products.duration_days)
type RenewDurationParams struct {
	CustomerUsername string
	ProductID        int32
	PricePaid        string // list_price as DECIMAL string
}

type CreateCustomerDurationParams struct {
	CustomerUsername string
	SalesUsername    *string
	ProductID        int32
	PurchaseDate     time.Time
	StartDate        time.Time
	EndDate          time.Time
	PricePaid        string // DECIMAL as string
	DiscountAmount   string // DECIMAL as string
	Status           string
}

type CustomerDurationInfo struct {
	ID               int32
	CustomerUsername string
	SalesUsername    *string
	ProductID        int32
	PurchaseDate     time.Time
	StartDate        time.Time
	EndDate          time.Time
	PricePaid        string
	DiscountAmount   string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CustomerDurationTableRow struct {
	DurationID        int32
	CustomerUsername  string
	CustomerFirstName string
	CustomerLastName  string
	ProductID         int32
	ProductName       string
	ProductType       string
	ProductCategory   string
	DurationDays      int32
	SalesUsername     string
	PurchaseDate      time.Time
	StartDate         time.Time
	EndDate           time.Time
	PricePaid         string
	DiscountAmount    string
	Status            string
}

type UpdateCustomerDurationEditableFieldsParams struct {
	ID             int32   // cd.id
	StartDateYMD   string  // "YYYY-MM-DD"
	PricePaid      string  // e.g. "5000.00"
	DiscountAmount *string // nil = ไม่อัปเดต/ให้เป็น NULL, ไม่ nil = "0.00", "100.00"
	Status         string  // "ACTIVE" | "EXPIRED" | "FROZEN" | "CANCELLED"
}

// ActiveDurationPackageInfo - Duration package ที่ ACTIVE พร้อมข้อมูลจาก JOIN products
type ActiveDurationPackageInfo struct {
	ID               int32
	CustomerUsername string
	ProductID        int32
	ProductName      string // JOIN from products
	DurationDays     int32  // JOIN from products
	SalesUsername    string
	PurchaseDate     time.Time
	StartDate        time.Time
	EndDate          time.Time
	DaysRemaining    int32  // คำนวณ: DATEDIFF(end_date, CURDATE())
	PricePaid        string // DECIMAL as string
	DiscountAmount   string // DECIMAL as string
	Status           string
	CreatedAt        time.Time
}

// RegisterDurationParams - Parameters สำหรับ RegisterDuration (Use Case 2.1C)
type RegisterDurationParams struct {
	// User table
	Username    string
	Password    string // Already hashed
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	PhoneNumber string
	Gmail       string

	// Customer table
	HealthInfo                   string
	Address                      string
	CompanyName                  string
	CompanyPosition              string
	MaritalStatus                string
	EmergencyContactName         string
	EmergencyContactRelationship string
	EmergencyContactPhone        string
	MarketingSource              string

	// Customer Duration table
	ProductID      int32
	SalesUsername  string
	PurchaseDate   time.Time
	StartDate      time.Time
	EndDate        time.Time
	PricePaid      string // DECIMAL as string
	DiscountAmount string // DECIMAL as string
}
