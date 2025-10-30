package responses

import "time"

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
