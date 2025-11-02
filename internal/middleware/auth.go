package middleware

import (
	"strings"

	"github.com/WaritDev/private-fitness-backend/domain/services"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates middleware that extracts JWT token and sets username in context
func AuthMiddleware(authUC *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from cookie or Authorization header
		token := c.Cookies("pf_auth")
		if token == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// If no token, return unauthorized
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Authentication required",
				"result":      nil,
			})
		}

		// Verify token
		payload, err := authUC.VerifyToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Invalid or expired token",
				"result":      nil,
			})
		}

		// Set username and role in locals for use in handlers
		c.Locals("username", payload.Sub)
		c.Locals("role", payload.Role)
		c.Locals("firstName", payload.FirstName)
		c.Locals("lastName", payload.LastName)

		// Continue to next handler
		return c.Next()
	}
}
