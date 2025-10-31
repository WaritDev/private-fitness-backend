package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// CustomerSessionRepository - Repository สำหรับจัดการ customer sessions
type CustomerSessionRepository interface {
	// RegisterCustomerSession - สร้าง User, Customer, CustomerSession, TrainingSchedules, CustomerLog ในครั้งเดียว (Transaction)
	RegisterCustomerSession(ctx context.Context, tx *sql.Tx, params RegisterSessionParams) (int32, error)

	// CheckBookingPermission - ตรวจสอบว่า Customer มีสิทธิ์จองหรือไม่ (มี Session package ACTIVE และยังมีสิทธิ์คงเหลือ)
	CheckBookingPermission(ctx context.Context, customerUsername string) (bool, error)

	// IncrementUsedSessions - Q3C.6: อัปเดตจำนวนครั้งที่ใช้ไป (used_sessions + 1)
	IncrementUsedSessions(ctx context.Context, tx *sql.Tx, sessionID int32) error

	// GetActiveSessionByCustomer - หา Session package ACTIVE ของ Customer
	GetActiveSessionByCustomer(ctx context.Context, customerUsername string) (*ActiveSessionInfo, error)

	// DecrementUsedSessions - ยกเลิกการจอง: ลดจำนวนครั้งที่ใช้ไป (used_sessions - 1)
	DecrementUsedSessions(ctx context.Context, tx *sql.Tx, sessionID int32) error

	// IncrementUsedSessionsByUsername - Q5C.2: อัปเดตจำนวนครั้งที่ใช้ไปสำหรับ Check-in (Use Case 5C)
	IncrementUsedSessionsByUsername(ctx context.Context, customerUsername string) error

	// GetCustomerActiveSessions - ดึงข้อมูล Session packages ที่ยัง ACTIVE ของลูกค้า
	GetCustomerActiveSessions(ctx context.Context, customerUsername string) ([]ActiveSessionPackageInfo, error)

	List(ctx context.Context) ([]dbmodel.ListCustomerSessionsRow, error)

	UpdateEditableFields(ctx context.Context, p UpdateCustomerSessionEditableFieldsParams) error
	CheckTrainerExists(ctx context.Context, username string) (int64, error)

	Delete(ctx context.Context, id int32) error

	GetByID(ctx context.Context, id int32) (dbmodel.GetCustomerSessionByIDRow, error)
}

// ActiveSessionInfo - ข้อมูล Session package ที่ active
type ActiveSessionInfo struct {
	ID               int32
	CustomerUsername string
	TrainerUsername  string
	TotalSessions    int32
	UsedSessions     int32
}

// ActiveSessionPackageInfo - ข้อมูล Session package ที่ ACTIVE พร้อมรายละเอียดจาก products
type ActiveSessionPackageInfo struct {
	ID                int32
	CustomerUsername  string
	TrainerUsername   string
	ProductID         int32
	ProductName       string
	TotalSessions     int32
	UsedSessions      int32
	SessionsRemaining int32
	PurchaseDate      time.Time
	PricePaid         string
	DiscountAmount    string
	Status            string
	CreatedAt         time.Time
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

type UpdateCustomerSessionEditableFieldsParams struct {
	ID              int32
	TrainerUsername string
	PricePaid       string
	DiscountAmount  sql.NullString
	Status          string
}
