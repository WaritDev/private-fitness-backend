package main

import (
	"fmt"

	"github.com/WaritDev/private-fitness-backend/config"
	"github.com/WaritDev/private-fitness-backend/internal/wire"
	"github.com/WaritDev/private-fitness-backend/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize controller
	controller := wire.InitializeController()

	// Initialize Fiber server
	cfg := config.ProvideConfig()

	// Fiber app
	app := fiber.New()

	// Cors config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://private-fitness-frontend.vercel.app", // 👈 ต้องระบุ origin ชัดเจน
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // 👈 ต้องเปิดเพื่อให้ส่ง cookie ได้
	}))

	router.RegisterApiRouter(app, controller)

	if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		panic(err)
	}
}
