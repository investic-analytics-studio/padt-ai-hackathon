package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

// func bindTopRankAPI(router fiber.Router, authMiddleware fiber.Handler)
func bindTopRankAPI(router fiber.Router) {

	topRankRepo := repo.NewTopRankRepo(infra.PostgresDB)
	topRankService := service.NewTopRankService(topRankRepo)
	topRankHandler := handler.NewTopRankHandler(topRankService)

	// topRankGroup := router.Group("/top-rank", authMiddleware)
	topRankGroup := router.Group("/top-rank")
	topRankGroup.Get("/", topRankHandler.GetTop100)
}
