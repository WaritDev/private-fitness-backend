package responses

import (
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// RegisterCustomerSessionResponse - Response หลังจากลงทะเบียนสำเร็จ
type RegisterCustomerSessionResponse struct {
	Username         string                 `json:"username"`
	SessionID        int32                  `json:"sessionId"`
	TrainerUsername  string                 `json:"trainerUsername"`
	ProductID        int32                  `json:"productId"`
	TotalSessions    int                    `json:"totalSessions"`
	SchedulesCreated int                    `json:"schedulesCreated"`
	CreatedSchedules []TrainingScheduleInfo `json:"createdSchedules"`
	Message          string                 `json:"message"`
}

// TrainingScheduleInfo - ข้อมูลนัดหมายที่สร้างแล้ว
type TrainingScheduleInfo struct {
	ScheduleID int32     `json:"scheduleId"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	DayOfWeek  string    `json:"dayOfWeek"`
}

// CheckPermissionResponse - Response สำหรับตรวจสอบสิทธิ์การจอง
type CheckPermissionResponse struct {
	HasPermission bool   `json:"hasPermission"`
	Message       string `json:"message"`
}

// CustomerSessionPackageResponse - Response สำหรับข้อมูล Session package
type CustomerSessionPackageResponse struct {
	ID                int32     `json:"id"`
	CustomerUsername  string    `json:"customerUsername"`
	TrainerUsername   string    `json:"trainerUsername"`
	ProductID         int32     `json:"productId"`
	ProductName       string    `json:"productName"`
	TotalSessions     int32     `json:"totalSessions"`
	UsedSessions      int32     `json:"usedSessions"`
	SessionsRemaining int32     `json:"sessionsRemaining"`
	PurchaseDate      time.Time `json:"purchaseDate"`
	PricePaid         float64   `json:"pricePaid"`
	DiscountAmount    float64   `json:"discountAmount"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
}

type ListCustomerSessionsResponse struct {
	Data []dbmodel.ListCustomerSessionsRow `json:"data"`
	Meta PageMeta                          `json:"meta"`
}

type CustomerSessionUpdatedResponse struct {
	Message string `json:"message"`
}

type CustomerSessionDeletedResponse struct {
	Message string `json:"message"`
}

type CustomerSession struct {
	ID               string `json:"id"`
	CustomerUsername string `json:"customerUsername"`
	TrainerUsername  string `json:"trainerUsername"`
	ProductID        string `json:"productId"`
	SalesUsername    string `json:"salesUsername"`
	PurchaseDate     string `json:"purchaseDate"`
	TotalSessions    int32  `json:"totalSessions"`
	UsedSessions     int32  `json:"usedSessions"`
	PricePaid        int64  `json:"pricePaid"`
	DiscountAmount   int64  `json:"discountAmount"`
	Status           string `json:"status"`
}

// RenewSessionResponse - Response สำหรับการต่ออายุ/ซื้อเพิ่ม Session Package
type RenewSessionResponse struct {
	ID               int32     `json:"id"`
	CustomerUsername string    `json:"customerUsername"`
	TrainerUsername  string    `json:"trainerUsername"`
	ProductID        int32     `json:"productId"`
	ProductName      string    `json:"productName"`
	TotalSessions    int32     `json:"totalSessions"`
	UsedSessions     int32     `json:"usedSessions"`
	PurchaseDate     time.Time `json:"purchaseDate"`
	PricePaid        float64   `json:"pricePaid"`
	DiscountAmount   float64   `json:"discountAmount"`
	Status           string    `json:"status"`
	Message          string    `json:"message"`
}
