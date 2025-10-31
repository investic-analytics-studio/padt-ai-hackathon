package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
)

func bindInitPageStatus(router fiber.Router) {
	cfg := config.GetConfig().PageShow
	initialPageStatusHandler := handler.NewInitialPageStatus(&cfg)
	router.Get("/initial-page-status", initialPageStatusHandler.GetInitialPageStatus)

}
