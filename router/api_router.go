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
	users.Get("/check-gmail", handler.User.CheckGmail)       // Q3S.2

	// Trainer routes (for Use Case: กรอกข้อมูลสมัครคอร์ส Sessions)
	trainers := apiGroup.Group("/trainers")
	trainers.Post("/match", handler.Trainer.MatchTrainer) // Q4S - Match trainer by day/time
	trainers.Get("/", handler.Trainer.ListAllTrainers)    // List all active trainers

	// Payment routes (for Use Case: ยืนยันการชำระเงิน)
	payments := apiGroup.Group("/payments")
	payments.Get("/info/:productId", handler.Payment.GetPaymentInfo) // Q5S.1 - Get payment info with QR code

	// Protected routes
	// customers := apiGroup.Group("/customer/")
	// customers.
}
