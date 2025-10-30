package requests

// GenerateQRCodeRequest - Use Case 5C: Request สำหรับสร้าง QR Code
type GenerateQRCodeRequest struct {
	PackageType string `json:"packageType"` // "DURATION" หรือ "SESSION"
}

// CheckInRequest - Use Case 5C: Request สำหรับ Check-in (จาก query string)
type CheckInRequest struct {
	Token string `json:"token"` // JWT token ที่ embedded ใน QR Code
}
