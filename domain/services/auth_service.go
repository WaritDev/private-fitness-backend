package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

type AuthService struct {
	authRepo repositories.AuthRepo
	userRepo repositories.UserRepo
}

func ProvideAuthService(authRepo repositories.AuthRepo, userRepo repositories.UserRepo) *AuthService {
	return &AuthService{authRepo: authRepo, userRepo: userRepo}
}

// Login authenticates a user and returns JWT token with user info
func (u *AuthService) Login(ctx context.Context, req requests.LoginRequest) (responses.LoginResponse, error) {
	// Get user from database
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return responses.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return responses.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	// Q0S.2: Update updated_at timestamp to track last login time
	if err := u.userRepo.UpdateLoginTimestamp(ctx, user.Username); err != nil {
		// Log error but don't fail login
		fmt.Printf("Warning: failed to update login timestamp for user %s: %v\n", user.Username, err)
	}

	// Generate JWT token
	payload := repositories.JWTPayload{
		Sub:       user.Username,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	token, err := u.authRepo.SignJWT(payload)
	if err != nil {
		return responses.LoginResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return responses.LoginResponse{
		Token: token,
		User: responses.UserInfo{
			Sub:       user.Username,
			Role:      string(user.Role),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive.Bool,
		},
	}, nil
}

// Signup registers a new customer user
func (u *AuthService) Signup(ctx context.Context, req requests.SignupRequest) error {
	// Basic validation (complex validation done on frontend)
	if len(req.Username) < 4 || len(req.Username) > 30 {
		return fmt.Errorf("username must be between 4 and 30 characters")
	}

	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Check if username exists
	exists, err := u.userRepo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Prepare user params
	userParams := repositories.CreateUserParams{
		Username:    req.Username,
		Password:    string(hashedPassword),
		Role:        "CUSTOMER",
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      stringPtr(req.Gender),
		DateOfBirth: stringPtr(req.DateOfBirth),
		PhoneNumber: stringPtr(req.PhoneNumber),
		Gmail:       stringPtr(req.Gmail),
	}

	// Prepare customer params
	customerParams := repositories.CreateCustomerParams{
		Username:                     req.Username,
		HealthInfo:                   stringPtr(req.HealthInfo),
		Address:                      stringPtr(req.Address),
		CompanyName:                  stringPtr(req.CompanyName),
		CompanyPosition:              stringPtr(req.CompanyPosition),
		MaritalStatus:                stringPtr(req.MaritalStatus),
		EmergencyContactName:         stringPtr(req.EmergencyContactName),
		EmergencyContactRelationship: stringPtr(req.EmergencyContactRelationship),
		EmergencyContactPhone:        stringPtr(req.EmergencyContactPhone),
		MarketingSource:              stringPtr(req.MarketingSource),
	}

	// Create user and customer in a transaction
	if err := u.userRepo.CreateUserWithCustomer(ctx, userParams, customerParams); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// VerifyToken verifies a JWT token and returns the payload
func (u *AuthService) VerifyToken(ctx context.Context, token string) (*repositories.JWTPayload, error) {
	payload, err := u.authRepo.VerifyJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return payload, nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
