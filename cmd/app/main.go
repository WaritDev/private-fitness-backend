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
	// Initialize handler
	handler := wire.InitializeHandler()

	// Initialize Fiber server
	cfg := config.ProvideConfig()

	// Fiber app
	app := fiber.New()

	// Cors config
	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	router.RegisterApiRouter(app, handler)

	if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		panic(err)
	}
}
