package controller

import (
	"os"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	UC *services.AuthService
}

func ProvideAuthController(uc *services.AuthService) *AuthController {
	return &AuthController{UC: uc}
}

// POST /api/auth/login
func (h *AuthController) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Username and password are required",
			"result":      nil,
		})
	}

	// Authenticate user
	result, err := h.UC.Login(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Invalid username or password",
				"result":      nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	// Set HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "pf_auth",
		Value:    result.Token,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   os.Getenv("NODE_ENV") == "production",
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Login successful",
		"result": fiber.Map{
			"token": result.Token,
			"user":  result.User,
		},
	})
}

// POST /api/auth/signup
func (h *AuthController) Signup(c *fiber.Ctx) error {
	var req requests.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Register user
	if err := h.UC.Signup(c.Context(), req); err != nil {
		statusCode := fiber.StatusInternalServerError
		if strings.Contains(err.Error(), "username already exists") {
			statusCode = fiber.StatusConflict
		} else if strings.Contains(err.Error(), "invalid username") || strings.Contains(err.Error(), "weak password") {
			statusCode = fiber.StatusBadRequest
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":      "Error",
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Created",
		"status_code": fiber.StatusCreated,
		"message":     "User registered successfully",
		"result":      fiber.Map{"ok": true},
	})
}

// POST/GET /api/auth/logout
func (h *AuthController) Logout(c *fiber.Ctx) error {
	// Clear cookie
	c.Cookie(&fiber.Cookie{
		Name:     "pf_auth",
		Value:    "",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   os.Getenv("NODE_ENV") == "production",
		Path:     "/",
		MaxAge:   -1, // Delete cookie
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Logged out successfully",
		"result":      fiber.Map{"ok": true},
	})
}

// GET /api/auth/me
func (h *AuthController) Me(c *fiber.Ctx) error {
	// Extract token from cookie or Authorization header
	token := c.Cookies("pf_auth")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "OK",
			"status_code": fiber.StatusOK,
			"message":     "User not authenticated",
			"result": fiber.Map{
				"authenticated": false,
			},
		})
	}

	// Verify token
	payload, err := h.UC.VerifyToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "OK",
			"status_code": fiber.StatusOK,
			"message":     "User not authenticated",
			"result": fiber.Map{
				"authenticated": false,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "User retrieved successfully",
		"result": fiber.Map{
			"authenticated": true,
			"user": fiber.Map{
				"sub":       payload.Sub,
				"role":      payload.Role,
				"firstName": payload.FirstName,
				"lastName":  payload.LastName,
			},
		},
	})
}
