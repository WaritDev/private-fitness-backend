package psql

import (
	"fmt"
	"os"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRepository struct{}

func ProvideAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

type jwtClaims struct {
	Sub       string `json:"sub"`
	Role      string `json:"role"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	jwt.RegisteredClaims
}

// SignJWT creates a JWT token with 7 days expiration
func (r *AuthRepository) SignJWT(payload repositories.JWTPayload) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	claims := jwtClaims{
		Sub:       payload.Sub,
		Role:      payload.Role,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyJWT verifies and decodes a JWT token
func (r *AuthRepository) VerifyJWT(tokenString string) (*repositories.JWTPayload, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &repositories.JWTPayload{
			Sub:       claims.Sub,
			Role:      claims.Role,
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
