package requests

import "time"

// MatchTrainerRequest - Request สำหรับจับคู่เทรนเนอร์
type MatchTrainerRequest struct {
	DayOfWeek string    `json:"dayOfWeek" validate:"required,oneof=MONDAY TUESDAY WEDNESDAY THURSDAY FRIDAY SATURDAY SUNDAY"` // วันในสัปดาห์
	StartTime time.Time `json:"startTime" validate:"required"`                                                                // เวลาเริ่มต้น (ISO 8601)
	EndTime   time.Time `json:"endTime" validate:"required"`                                                                  // เวลาสิ้นสุด (ISO 8601)
}
