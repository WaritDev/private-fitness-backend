package router

import (
	"github.com/WaritDev/private-fitness-backend/internal/adapters/rest/controller"
	"github.com/WaritDev/private-fitness-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterApiRouter(app *fiber.App, controller *controller.Controller) {
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
	staffs.Get("/", controller.Staff.ListStaffs)
	staffs.Post("/create", controller.Staff.CreateStaff)
	staffs.Post("/:username/update", controller.Staff.UpdateStaff)
	staffs.Delete("/:username", controller.Staff.DeleteStaff)
	staffs.Get("/:username", controller.Staff.GetStaffByUsername)

	// Manager routes
	managers := apiGroup.Group("/manager")
	managers.Get("/dashboard", controller.Manager.GetDashboard)

	// Customer routes
	customers := apiGroup.Group("/customers")
	customers.Get("/", controller.Customer.ListCustomers)
	customers.Post("/:username/update", controller.Customer.UpdateCustomer)
	customers.Delete("/:username", controller.Customer.DeleteCustomer)
	customers.Get("/:username", controller.Customer.GetByUsername)

	// Customer Duration routes
	durations := apiGroup.Group("/customer-durations")
	durations.Post("/renew", controller.CustomerDuration.RenewDuration) // ✅ ต่ออายุ Duration (ลูกค้าซื้อเอง)
	durations.Get("/", controller.CustomerDuration.ListDurations)
	durations.Post("/:id/update", controller.CustomerDuration.Update)
	durations.Delete("/:id", controller.CustomerDuration.Delete)
	durations.Get("/:id", controller.CustomerDuration.GetByID)

	// Customer Session routes
	sessions := apiGroup.Group("/customer-sessions")
	sessions.Post("/renew", controller.CustomerSession.RenewSession) // ✅ ต่ออายุ Session (ลูกค้าซื้อเอง)
	sessions.Get("/", controller.CustomerSession.ListSessions)
	sessions.Post("/:id/update", controller.CustomerSession.Update)
	sessions.Delete("/:id", controller.CustomerSession.Delete)
	sessions.Get("/:id", controller.CustomerSession.GetByID)

	// Customer Log routes
	customerLogs := apiGroup.Group("/customer-logs")
	customerLogs.Get("/", controller.CustomerLog.List)
	customerLogs.Post("/:id/update", controller.CustomerLog.Update)
	customerLogs.Delete("/:id", controller.CustomerLog.Delete)
	customerLogs.Get("/:id", controller.CustomerLog.GetByID)

	// Auth routes
	authGroup := apiGroup.Group("/auth")
	authGroup.Post("/login", controller.Auth.Login)
	// NOTE: /signup removed - registration now requires Sales pre-entry (Use Case 3S + 4S) before customer completes (Use Case 2.2C)
	authGroup.Post("/logout", controller.Auth.Logout)
	authGroup.Get("/logout", controller.Auth.Logout)
	authGroup.Get("/me", controller.Auth.Me)

	products := apiGroup.Group("/products")
	products.Get("/", controller.Product.ListAllProducts)

	// ⚠️ IMPORTANT: Specific routes MUST come BEFORE dynamic routes (/:id)
	// Otherwise /durations will be matched as id="durations"
	products.Get("/durations", controller.Product.ListDurations)
	products.Get("/sessions", controller.Product.ListSessions)
	products.Get("/", controller.Product.List)
	products.Post("/create", controller.Product.Create)
	products.Post("/:id/update", controller.Product.Update)
	products.Delete("/:id", controller.Product.Delete)
	products.Get("/:id", controller.Product.GetByID)

	// Dynamic route must be last
	products.Get("/:id", controller.Product.GetProductByID)

	// User validation routes (for Use Case: กรอกข้อมูลสมาชิก)
	users := apiGroup.Group("/users")
	users.Get("/check-phone", controller.User.CheckPhoneNumber)
	users.Get("/check-gmail", controller.User.CheckGmail)
	users.Get("/check-username", controller.User.CheckUsername)

	// Trainer routes (for Use Case: กรอกข้อมูลสมัครคอร์ส Sessions)
	trainers := apiGroup.Group("/trainers")
	trainers.Post("/match", controller.Trainer.MatchTrainer)
	trainers.Get("/", controller.Trainer.ListAllTrainers) // List all active trainers

	// Use Case 1P: Manage Working Hours (Trainer must be logged in)
	// Apply authentication middleware to working hours routes
	workingHoursGroup := trainers.Group("/working-hours", middleware.AuthMiddleware(controller.Auth.UC))
	workingHoursGroup.Get("/", controller.Trainer.GetWorkingHours)         // Q1P.1: Get trainer's working hours
	workingHoursGroup.Post("/", controller.Trainer.AddWorkingTime)         // Q1P.2 + Q1P.3: Add new working time
	workingHoursGroup.Put("/:id", controller.Trainer.UpdateWorkingTime)    // Q1P.4: Update working time
	workingHoursGroup.Delete("/:id", controller.Trainer.DeleteWorkingTime) // Q1P.5: Delete working time

	// Use Case 3P: Manage Day-Offs (Trainer must be logged in)
	// Apply authentication middleware to day-offs routes
	dayOffsGroup := trainers.Group("/day-offs", middleware.AuthMiddleware(controller.Auth.UC))
	dayOffsGroup.Get("/", controller.Trainer.GetDayOffs)         // Q3P.1: Get trainer's day-offs
	dayOffsGroup.Post("/", controller.Trainer.AddDayOff)         // Q3P.2 + Q3P.3 + Q3P.4: Add new day-off
	dayOffsGroup.Delete("/:id", controller.Trainer.DeleteDayOff) // Q3P.5: Delete day-off

	// Use Case: Trainer Calendar & Check-in Confirmation (Trainer must be logged in)
	calendarGroup := trainers.Group("/calendar", middleware.AuthMiddleware(controller.Auth.UC))
	calendarGroup.Get("/", controller.Trainer.GetCalendar) // Get appointments with pending check-ins

	checkinGroup := trainers.Group("/checkin", middleware.AuthMiddleware(controller.Auth.UC))
	checkinGroup.Post("/", controller.Trainer.ConfirmCheckIn) // Trainer confirm check-in and deduct session

	// Payment routes (for Use Case: ยืนยันการชำระเงิน)
	payments := apiGroup.Group("/payments")
	payments.Get("/info/:productId", controller.Payment.GetPaymentInfo)
	payments.Post("/verify-slip", controller.Payment.VerifySlip) // Payment slip verification with Slip2Go API
	payments.Get("/", controller.Payment.List)
	payments.Post("/create", controller.Payment.Create)
	payments.Post("/:id/update", controller.Payment.Update)
	payments.Delete("/:id", controller.Payment.Delete)
	payments.Get("/:id", controller.Payment.GetByID)

	// Customer registration routes (after Sales pre-entry)
	customers.Post("/sessions/register", controller.CustomerSession.Register)
	customers.Post("/durations/register", controller.CustomerDuration.Register)
	customers.Get("/sessions/check-permission", controller.CustomerSession.CheckPermission)     // Check booking permission
	customers.Get("/sessions/active/:username", controller.CustomerSession.GetActiveSessions)   // Get active session packages
	customers.Get("/durations/active/:username", controller.CustomerDuration.GetActiveDuration) // Get active duration packages

	// Booking routes (Use Case 3C: จองเวลาออกกำลังกายกับเทรนเนอร์)
	bookings := apiGroup.Group("/bookings")
	bookings.Get("/slots", controller.Booking.GetBookingSlots)           // Get available booking slots
	bookings.Post("/book", controller.Booking.BookAppointment)           // Q3C: Book appointment with trainer
	bookings.Delete("/cancel/:id", controller.Booking.CancelAppointment) // Cancel appointment

	// Member routes (Use Case 5C: สแกนเข้า Fitness)
	member := apiGroup.Group("/member")
	member.Post("/qrcode", controller.Member.GenerateQRCode) // Generate QR Code for check-in

	// Check-in route (public - accessed by QR scanner)
	apiGroup.Get("/checkin", controller.Member.CheckIn) // Check-in via QR Code scan

	// Protected routes
	// customers := apiGroup.Group("/customer/")
	// customers.
}
