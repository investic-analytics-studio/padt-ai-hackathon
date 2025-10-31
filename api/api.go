package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	v2 "github.com/quantsmithapp/datastation-backend/api/v2"
)

func InitAPI(app *fiber.App) {
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, X-Knowledge-Base", // Optionally expose any custom headers you may use
		MaxAge:           86400,                              // 24 hours cache for preflight requests
	}))
	app.Use(helmet.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))

	api := app.Group("/api")
	v2.BindApiV2(api)
}
