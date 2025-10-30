package responses

// GenerateQRCodeResponse - Use Case 5C: Response หลังสร้าง QR Code
type GenerateQRCodeResponse struct {
	Status      string `json:"status"`      // "success"
	QRCodeURL   string `json:"qrCodeUrl"`   // URL สำหรับฝังใน QR Code
	Token       string `json:"token"`       // JWT token (ให้ frontend ใช้ generate QR)
	PackageType string `json:"packageType"` // "DURATION" หรือ "SESSION"
	ExpiresIn   int64  `json:"expiresIn"`   // เวลาหมดอายุ (seconds) - 60 seconds
	Message     string `json:"message"`     // คำอธิบาย
}

// CheckInResponse - Use Case 5C: Response หลัง Check-in สำเร็จ
type CheckInResponse struct {
	Status      string `json:"status"`      // "success"
	Message     string `json:"message"`     // "Welcome, [FirstName]!"
	Username    string `json:"username"`    // username ของลูกค้า
	FirstName   string `json:"firstName"`   // ชื่อลูกค้า
	PackageType string `json:"packageType"` // "DURATION" หรือ "SESSION"
	LogID       int32  `json:"logId"`       // ID ของ log ที่สร้าง
}
