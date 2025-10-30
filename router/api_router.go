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

	
	// api 
	apiGroup := app.Group("/api")

	// Staff routes
	staffs := apiGroup.Group("/staffs")
	staffs.Get("/", handler.Staff.ListStaffs)
	staffs.Post("/create", handler.Staff.CreateStaff)
	staffs.Post("/:username/update", handler.Staff.UpdateStaff)
	staffs.Delete("/:username", handler.Staff.DeleteStaff)

	// Manager routes
	managers := apiGroup.Group("/manager")
	managers.Get("/dashboard", handler.Manager.GetDashboard)
	
	// Auth routes
	authGroup := apiGroup.Group("/auth")
	authGroup.Post("/login", handler.Auth.Login)
	authGroup.Post("/signup", handler.Auth.Signup)
	authGroup.Post("/logout", handler.Auth.Logout)
	authGroup.Get("/logout", handler.Auth.Logout)
	authGroup.Get("/me", handler.Auth.Me)

	products := apiGroup.Group("/products")
	products.Get("/", handler.Product.ListAllProducts)
	products.Get("/:id", handler.Product.GetProductByID)

	durations := products.Group("/durations")
	durations.Get("/", handler.Product.ListDurations)

	sessions := products.Group("/sessions")
	sessions.Get("/", handler.Product.ListSessions)

	// Protected routes
	// customers := apiGroup.Group("/customer/")
	// customers.
}
