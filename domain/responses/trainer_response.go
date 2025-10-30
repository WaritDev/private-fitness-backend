package responses

import "time"

// TrainerMatchResponse - Response สำหรับผลลัพธ์การจับคู่เทรนเนอร์
type TrainerMatchResponse struct {
	TrainerUsername string    `json:"trainerUsername"` // Username ของเทรนเนอร์
	TrainerName     string    `json:"trainerName"`     // ชื่อเต็มของเทรนเนอร์
	DayOfWeek       string    `json:"dayOfWeek"`       // วันในสัปดาห์ที่จับคู่
	StartTime       time.Time `json:"startTime"`       // เวลาเริ่มต้น
	EndTime         time.Time `json:"endTime"`         // เวลาสิ้นสุด
	Appointments    int64     `json:"appointments"`    // จำนวนนัดหมายที่เทรนเนอร์มีในวันนั้น
}

// TrainerListResponse - Response สำหรับรายชื่อเทรนเนอร์
type TrainerListResponse struct {
	Username string `json:"username"` // Username
	Name     string `json:"name"`     // ชื่อเต็ม
}

// Use Case 1P: Manage Working Hours

// TrainerAvailabilityResponse - Response สำหรับแสดงเวลาทำงานของ Trainer
type TrainerAvailabilityResponse struct {
	AvailabilityID int32  `json:"availabilityId"` // ID ของเวลาทำงาน
	DayOfWeek      string `json:"dayOfWeek"`      // วันในสัปดาห์
	StartTime      string `json:"startTime"`      // เวลาเริ่มต้น (HH:MM format)
	EndTime        string `json:"endTime"`        // เวลาสิ้นสุด (HH:MM format)
}

// WorkingHoursResponse - Response สำหรับรายการเวลาทำงานทั้งหมด
type WorkingHoursResponse struct {
	Status       string                        `json:"status"`       // สถานะ "success"
	Message      string                        `json:"message"`      // ข้อความ
	WorkingHours []TrainerAvailabilityResponse `json:"workingHours"` // รายการเวลาทำงาน
}

// AddWorkingTimeResponse - Response สำหรับการเพิ่มเวลาทำงาน
type AddWorkingTimeResponse struct {
	Status  string `json:"status"`  // สถานะ "success" หรือ "error"
	Message string `json:"message"` // ข้อความ
}

// UpdateWorkingTimeResponse - Response สำหรับการแก้ไขเวลาทำงาน (Q1P.4)
type UpdateWorkingTimeResponse struct {
	Status  string `json:"status"`  // สถานะ "success" หรือ "error"
	Message string `json:"message"` // ข้อความ
}

// DeleteWorkingTimeResponse - Response สำหรับการลบเวลาทำงาน (Q1P.5)
type DeleteWorkingTimeResponse struct {
	Status  string `json:"status"`  // สถานะ "success" หรือ "error"
	Message string `json:"message"` // ข้อความ
}

// Use Case 3P: Manage Day-Offs

// DayOffResponse - Response สำหรับแสดงวันหยุดของ Trainer
type DayOffResponse struct {
	ScheduleID int32     `json:"scheduleId"` // ID ของวันหยุด
	StartTime  time.Time `json:"startTime"`  // เวลาเริ่มต้น (วันหยุด 00:00:00)
	EndTime    time.Time `json:"endTime"`    // เวลาสิ้นสุด (วันหยุด 23:59:59)
}

// DayOffsListResponse - Response สำหรับรายการวันหยุดทั้งหมด
type DayOffsListResponse struct {
	Status  string           `json:"status"`  // สถานะ "success"
	Message string           `json:"message"` // ข้อความ
	DayOffs []DayOffResponse `json:"dayOffs"` // รายการวันหยุด
}

// AddDayOffResponse - Response สำหรับการเพิ่มวันหยุด
type AddDayOffResponse struct {
	Status  string `json:"status"`  // สถานะ "success" หรือ "error"
	Message string `json:"message"` // ข้อความ
}

// DeleteDayOffResponse - Response สำหรับการลบวันหยุด
type DeleteDayOffResponse struct {
	Status  string `json:"status"`  // สถานะ "success" หรือ "error"
	Message string `json:"message"` // ข้อความ
}
