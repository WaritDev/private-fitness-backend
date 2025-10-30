package requests

import "time"

// MatchTrainerRequest - Request สำหรับจับคู่เทรนเนอร์
type MatchTrainerRequest struct {
	DayOfWeek string    `json:"dayOfWeek" validate:"required,oneof=MONDAY TUESDAY WEDNESDAY THURSDAY FRIDAY SATURDAY SUNDAY"` // วันในสัปดาห์
	StartTime time.Time `json:"startTime" validate:"required"`                                                                // เวลาเริ่มต้น (ISO 8601)
	EndTime   time.Time `json:"endTime" validate:"required"`                                                                  // เวลาสิ้นสุด (ISO 8601)
}

// Use Case 1P: Manage Working Hours

// AddWorkingTimeRequest - Request สำหรับเพิ่มเวลาทำงานของ Trainer
type AddWorkingTimeRequest struct {
	DayOfWeek string `json:"dayOfWeek" validate:"required,oneof=MONDAY TUESDAY WEDNESDAY THURSDAY FRIDAY SATURDAY SUNDAY"` // วันในสัปดาห์ที่พร้อมให้บริการ
	StartTime string `json:"startTime" validate:"required"`                                                                // เวลาเริ่มต้น (HH:MM format)
	EndTime   string `json:"endTime" validate:"required"`                                                                  // เวลาสิ้นสุด (HH:MM format)
}

// UpdateWorkingTimeRequest - Request สำหรับแก้ไขเวลาทำงานของ Trainer (Q1P.4)
type UpdateWorkingTimeRequest struct {
	DayOfWeek string `json:"dayOfWeek" validate:"required,oneof=MONDAY TUESDAY WEDNESDAY THURSDAY FRIDAY SATURDAY SUNDAY"` // วันในสัปดาห์ที่พร้อมให้บริการ
	StartTime string `json:"startTime" validate:"required"`                                                                // เวลาเริ่มต้น (HH:MM format)
	EndTime   string `json:"endTime" validate:"required"`                                                                  // เวลาสิ้นสุด (HH:MM format)
}
