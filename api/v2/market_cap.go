package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindMarketCapAPI(router fiber.Router, authMiddleware fiber.Handler) {

	marketCapRepo := repo.NewMarketCapRepo(infra.PostgresDB)
	marketCapService := service.NewMarketCapService(marketCapRepo)
	marketCapHandler := handler.NewMarketCapHandler(marketCapService)

	// marketCapGroup := router.Group("/market-cap", authMiddleware)
	_ = authMiddleware
	marketCapGroup := router.Group("/market-cap")
	marketCapGroup.Get("/", marketCapHandler.GetMarketCap)
}
