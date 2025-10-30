package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/golang-jwt/jwt/v5"
)

type MemberUseCase struct {
	logRepo     repositories.CustomerLogRepository
	sessionRepo repositories.CustomerSessionRepository
	userRepo    repositories.UserRepo
	authRepo    repositories.AuthRepo
	db          *sql.DB
}

func ProvideMemberUseCase(
	logRepo repositories.CustomerLogRepository,
	sessionRepo repositories.CustomerSessionRepository,
	userRepo repositories.UserRepo,
	authRepo repositories.AuthRepo,
	db *sql.DB,
) *MemberUseCase {
	return &MemberUseCase{
		logRepo:     logRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		authRepo:    authRepo,
		db:          db,
	}
}

// QRTokenPayload - Payload สำหรับ QR Code JWT Token
type QRTokenPayload struct {
	Sub         string `json:"sub"`         // username
	PackageType string `json:"packageType"` // "DURATION" หรือ "SESSION"
	jwt.RegisteredClaims
}

// getJWTSecret - ดึง JWT secret จาก environment variable
func (u *MemberUseCase) getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	return secret
}

// GenerateQRToken - Use Case 5C: สร้าง JWT token สำหรับ QR Code (หมดอายุใน 60 วินาที)
func (u *MemberUseCase) GenerateQRToken(username, packageType string) (string, error) {
	// Validate package type
	if packageType != "DURATION" && packageType != "SESSION" {
		return "", fmt.Errorf("invalid package type: %s (must be DURATION or SESSION)", packageType)
	}

	// Create token with 60 seconds expiry
	expirationTime := time.Now().Add(60 * time.Second)
	claims := &QRTokenPayload{
		Sub:         username,
		PackageType: packageType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.getJWTSecret()))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// VerifyQRToken - Use Case 5C: Verify JWT token จาก QR Code
func (u *MemberUseCase) VerifyQRToken(tokenString string) (*QRTokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &QRTokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.getJWTSecret()), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*QRTokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// CheckIn - Use Case 5C: บันทึกการเข้าใช้งานฟิตเนส (Q5C.1 และ Q5C.2)
func (u *MemberUseCase) CheckIn(ctx context.Context, username, packageType string) (*responses.CheckInResponse, error) {
	// ดึงข้อมูลผู้ใช้
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Begin transaction
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Q5C.1: สร้าง customer log (CHECK_IN)
	err = u.logRepo.CreateCustomerLog(ctx, tx, username, "CHECK_IN")
	if err != nil {
		return nil, fmt.Errorf("failed to create check-in log: %w", err)
	}

	// Q5C.2: ถ้าเป็น SESSION package ให้อัปเดต used_sessions
	if packageType == "SESSION" {
		err = u.sessionRepo.IncrementUsedSessionsByUsername(ctx, username)
		if err != nil {
			// Log warning แต่ไม่ fail การ check-in
			fmt.Printf("Warning: failed to increment used sessions for %s: %v\n", username, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Build response
	return &responses.CheckInResponse{
		Status:      "success",
		Message:     fmt.Sprintf("Welcome, %s!", user.FirstName),
		Username:    username,
		FirstName:   user.FirstName,
		PackageType: packageType,
		LogID:       0, // TODO: get actual log ID if needed
	}, nil
}
