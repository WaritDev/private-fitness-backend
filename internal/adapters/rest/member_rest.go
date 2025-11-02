package rest

import (
	"fmt"
	"strings"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type MemberHandler struct {
	MemberUC *usecases.MemberUseCase
	AuthUC   *usecases.AuthUseCase
}

func ProvideMemberHandler(
	memberUC *usecases.MemberUseCase,
	authUC *usecases.AuthUseCase,
) *MemberHandler {
	return &MemberHandler{
		MemberUC: memberUC,
		AuthUC:   authUC,
	}
}

// GenerateQRCode - POST /api/member/qrcode
// Use Case 5C: สร้าง QR Code สำหรับ Check-in
func (h *MemberHandler) GenerateQRCode(c *fiber.Ctx) error {
	// Extract and verify token
	token := c.Cookies("pf_auth")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	fmt.Println("token", token)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Authentication required",
			"result":      nil,
		})
	}

	// Verify JWT token
	payload, err := h.AuthUC.VerifyToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Invalid or expired token",
			"result":      nil,
		})
	}

	// Parse request body
	var req requests.GenerateQRCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Validate package type
	if req.PackageType != "DURATION" && req.PackageType != "SESSION" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid package type (must be DURATION or SESSION)",
			"result":      nil,
		})
	}

	// Generate QR Token
	qrToken, err := h.MemberUC.GenerateQRToken(payload.Sub, req.PackageType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	// Build QR Code URL
	// TODO: แก้เป็น domain จริง (ตอนนี้ใช้ localhost)
	baseURL := "http://localhost:8000"
	qrCodeURL := fmt.Sprintf("%s/api/checkin?token=%s", baseURL, qrToken)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusOK,
		"message":     "QR Code generated successfully",
		"result": fiber.Map{
			"qrCodeUrl":   qrCodeURL,
			"token":       qrToken,
			"packageType": req.PackageType,
			"expiresIn":   60, // 60 seconds
		},
	})
}

// CheckIn - GET /api/checkin?token=xxx
// Use Case 5C: สแกน QR Code เพื่อ Check-in
func (h *MemberHandler) CheckIn(c *fiber.Ctx) error {
	// Get token from query string
	tokenString := c.Query("token")
	if tokenString == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing token parameter")
	}

	// Verify QR Token
	payload, err := h.MemberUC.VerifyQRToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired QR code")
	}

	// Perform check-in
	result, err := h.MemberUC.CheckIn(c.Context(), payload.Sub, payload.PackageType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Check-in failed: %s", err.Error()))
	}

	// Return success message (simple HTML for display on scanner device)
	htmlResponse := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Check-in Success</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
		}
		.container {
			text-align: center;
			background: white;
			padding: 40px;
			border-radius: 20px;
			box-shadow: 0 10px 40px rgba(0,0,0,0.2);
		}
		.success {
			font-size: 60px;
			color: #10b981;
			margin-bottom: 20px;
		}
		h1 {
			color: #1f2937;
			margin-bottom: 10px;
		}
		p {
			color: #6b7280;
			font-size: 18px;
		}
		.badge {
			display: inline-block;
			padding: 8px 16px;
			background: #ede9fe;
			color: #7c3aed;
			border-radius: 20px;
			margin-top: 20px;
			font-weight: bold;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="success">✅</div>
		<h1>%s</h1>
		<p>User: <strong>%s</strong></p>
		<div class="badge">%s Package</div>
	</div>
</body>
</html>
	`, result.Message, result.Username, result.PackageType)

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(htmlResponse)
}
