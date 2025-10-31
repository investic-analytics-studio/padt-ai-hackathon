package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindWinRateAPI(router fiber.Router, authMiddleware fiber.Handler) {
	winRateRepo := repo.NewWinRateRepo(infra.PostgresDB)
	winRateService := service.NewWinRateService(winRateRepo)
	winRateHandler := handler.NewWinRateHandler(winRateService)

	winRateGroup := router.Group("/win-rate", authMiddleware)
	winRateGroup.Get("/", winRateHandler.GetWinRate)
}
