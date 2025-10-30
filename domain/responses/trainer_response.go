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
