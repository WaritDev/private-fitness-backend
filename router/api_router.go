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
	staffs.Get("/:username", handler.Staff.GetStaffByUsername)

	// Manager routes
	managers := apiGroup.Group("/manager")
	managers.Get("/dashboard", handler.Manager.GetDashboard)

	// Customer routes
	customers := apiGroup.Group("/customers")
	customers.Get("/", handler.Customer.ListCustomers)
	customers.Post("/:username/update", handler.Customer.UpdateCustomer)
	customers.Delete("/:username", handler.Customer.DeleteCustomer)
	customers.Get("/:username", handler.Customer.GetByUsername)

	// Customer Duration routes
	durations := apiGroup.Group("/customer-durations")
	durations.Get("/", handler.CustomerDuration.ListDurations)
	durations.Post("/:id/update", handler.CustomerDuration.Update)
	durations.Delete("/:id", handler.CustomerDuration.Delete)
	durations.Get("/:id", handler.CustomerDuration.GetByID)

	// Customer Session routes
	sessions := apiGroup.Group("/customer-sessions")
	sessions.Get("/", handler.CustomerSession.ListSessions)
	sessions.Post("/:id/update", handler.CustomerSession.Update)
	sessions.Delete("/:id", handler.CustomerSession.Delete)
	sessions.Get("/:id", handler.CustomerSession.GetByID)

	// Customer Log routes
	customerLogs := apiGroup.Group("/customer-logs")
	customerLogs.Get("/", handler.CustomerLog.List)
	customerLogs.Post("/:id/update", handler.CustomerLog.Update)
	customerLogs.Delete("/:id", handler.CustomerLog.Delete)

	// Auth routes
	authGroup := apiGroup.Group("/auth")
	authGroup.Post("/login", handler.Auth.Login)
	// NOTE: /signup removed - registration now requires Sales pre-entry (Use Case 3S + 4S) before customer completes (Use Case 2.2C)
	authGroup.Post("/logout", handler.Auth.Logout)
	authGroup.Get("/logout", handler.Auth.Logout)
	authGroup.Get("/me", handler.Auth.Me)

	products := apiGroup.Group("/products")
	products.Get("/", handler.Product.ListAllProducts)

	// ⚠️ IMPORTANT: Specific routes MUST come BEFORE dynamic routes (/:id)
	// Otherwise /durations will be matched as id="durations"
	products.Get("/durations", handler.Product.ListDurations)
	products.Get("/sessions", handler.Product.ListSessions)
	products.Get("/", handler.Product.List)
	products.Post("/create", handler.Product.Create)
	products.Post("/:id/update", handler.Product.Update)
	products.Delete("/:id", handler.Product.Delete)

	// Dynamic route must be last
	products.Get("/:id", handler.Product.GetProductByID)

	// User validation routes (for Use Case: กรอกข้อมูลสมาชิก)
	users := apiGroup.Group("/users")
	users.Get("/check-phone", handler.User.CheckPhoneNumber)
	users.Get("/check-gmail", handler.User.CheckGmail)

	// Trainer routes (for Use Case: กรอกข้อมูลสมัครคอร์ส Sessions)
	trainers := apiGroup.Group("/trainers")
	trainers.Post("/match", handler.Trainer.MatchTrainer)
	trainers.Get("/", handler.Trainer.ListAllTrainers) // List all active trainers

	// Use Case 1P: Manage Working Hours (Trainer must be logged in)
	trainers.Get("/working-hours", handler.Trainer.GetWorkingHours)          // Q1P.1: Get trainer's working hours
	trainers.Post("/working-hours", handler.Trainer.AddWorkingTime)          // Q1P.2 + Q1P.3: Add new working time
	trainers.Put("/working-hours/:id", handler.Trainer.UpdateWorkingTime)    // Q1P.4: Update working time
	trainers.Delete("/working-hours/:id", handler.Trainer.DeleteWorkingTime) // Q1P.5: Delete working time

	// Use Case 3P: Manage Day-Offs (Trainer must be logged in)
	trainers.Get("/day-offs", handler.Trainer.GetDayOffs)          // Q3P.1: Get trainer's day-offs
	trainers.Post("/day-offs", handler.Trainer.AddDayOff)          // Q3P.2 + Q3P.3 + Q3P.4: Add new day-off
	trainers.Delete("/day-offs/:id", handler.Trainer.DeleteDayOff) // Q3P.5: Delete day-off

	// Payment routes (for Use Case: ยืนยันการชำระเงิน)
	payments := apiGroup.Group("/payments")
	payments.Get("/info/:productId", handler.Payment.GetPaymentInfo)
	payments.Get("/", handler.Payment.List)
	payments.Post("/create", handler.Payment.Create)
	payments.Post("/:id/update", handler.Payment.Update)
	payments.Delete("/:id", handler.Payment.Delete)

	// Customer registration routes (after Sales pre-entry)
	customers.Post("/sessions/register", handler.CustomerSession.Register)
	customers.Post("/durations/register", handler.CustomerDuration.Register)
	customers.Get("/sessions/check-permission", handler.CustomerSession.CheckPermission)     // Check booking permission
	customers.Get("/sessions/active/:username", handler.CustomerSession.GetActiveSessions)   // Get active session packages
	customers.Get("/durations/active/:username", handler.CustomerDuration.GetActiveDuration) // Get active duration packages

	// Booking routes (Use Case 3C: จองเวลาออกกำลังกายกับเทรนเนอร์)
	bookings := apiGroup.Group("/bookings")
	bookings.Get("/slots", handler.Booking.GetBookingSlots)           // Get available booking slots
	bookings.Post("/book", handler.Booking.BookAppointment)           // Q3C: Book appointment with trainer
	bookings.Delete("/cancel/:id", handler.Booking.CancelAppointment) // Cancel appointment

	// Member routes (Use Case 5C: สแกนเข้า Fitness)
	member := apiGroup.Group("/member")
	member.Post("/qrcode", handler.Member.GenerateQRCode) // Generate QR Code for check-in

	// Check-in route (public - accessed by QR scanner)
	apiGroup.Get("/checkin", handler.Member.CheckIn) // Check-in via QR Code scan

	// Protected routes
	// customers := apiGroup.Group("/customer/")
	// customers.
}
