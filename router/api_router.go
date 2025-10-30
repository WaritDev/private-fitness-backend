package router

import (
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest"
	"github.com/gofiber/fiber/v2"
)

func RegisterApiRouter(app *fiber.App, handler *rest.Handler) {
	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	})

	app.Get("/api/manager/dashboard", handler.Manager.GetDashboard)

	// api
	apiGroup := app.Group("/api")
	// Auth routes
	authGroup := apiGroup.Group("/auth")
	authGroup.Post("/login", handler.Auth.Login)
	authGroup.Post("/signup", handler.Auth.Signup)
	authGroup.Post("/logout", handler.Auth.Logout)
	authGroup.Get("/logout", handler.Auth.Logout)
	authGroup.Get("/me", handler.Auth.Me)

	products := apiGroup.Group("/products")
	products.Get("/", handler.Product.ListAllProducts)

	// ⚠️ IMPORTANT: Specific routes MUST come BEFORE dynamic routes (/:id)
	// Otherwise /durations will be matched as id="durations"
	products.Get("/durations", handler.Product.ListDurations)
	products.Get("/sessions", handler.Product.ListSessions)

	// Dynamic route must be last
	products.Get("/:id", handler.Product.GetProductByID)

	// User validation routes (for Use Case: กรอกข้อมูลสมาชิก)
	users := apiGroup.Group("/users")
	users.Get("/check-phone", handler.User.CheckPhoneNumber) // Q3S.1
	users.Get("/check-email", handler.User.CheckEmail)       // Q3S.2

	// Protected routes
	// customers := apiGroup.Group("/customer/")
	// customers.
}
