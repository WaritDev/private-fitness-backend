package repositories

import (
	"context"
	"database/sql"
	"time"
)

// CustomerSessionRepository - Repository สำหรับจัดการ customer sessions
type CustomerSessionRepository interface {
	// RegisterCustomerSession - สร้าง User, Customer, CustomerSession, TrainingSchedules, CustomerLog ในครั้งเดียว (Transaction)
	RegisterCustomerSession(ctx context.Context, tx *sql.Tx, params RegisterSessionParams) (int32, error)
}

// RegisterSessionParams - Parameters สำหรับ RegisterCustomerSession
type RegisterSessionParams struct {
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

	// Customer Session table
	TrainerUsername string
	ProductID       int32
	SalesUsername   string
	PurchaseDate    time.Time
	TotalSessions   int32 // Total number of sessions in the package
	UsedSessions    int32
	PricePaid       string // DECIMAL as string
	DiscountAmount  string // DECIMAL as string

	// Training Schedules (multiple)
	Schedules []TrainingScheduleParams
}

// TrainingScheduleParams - Parameters for each training schedule
type TrainingScheduleParams struct {
	TrainerUsername  string
	CustomerUsername string
	SessionID        int32 // Will be set after creating customer_session
	StartTime        time.Time
	EndTime          time.Time // StartTime + 2 hours
	ScheduleType     string    // 'APPOINTMENT'
}
