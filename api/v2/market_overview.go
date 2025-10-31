package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

// func bindMarketOverviewAPI(router fiber.Router, authMiddleware fiber.Handler)
func bindMarketOverviewAPI(router fiber.Router) {
	marketOverviewRepo := repo.NewMarketOverviewRepo(infra.PostgresDB)
	marketOverviewService := service.NewMarketOverviewService(marketOverviewRepo)
	marketOverviewHandler := handler.NewMarketOverviewHandler(marketOverviewService)
	// marketOverviewGroup := router.Group("/market-overview", authMiddleware)
	marketOverviewGroup := router.Group("/market-overview")
	marketOverviewGroup.Get("/sentiment-market-overview", marketOverviewHandler.GetMarketOverview)
	marketOverviewGroup.Get("/table", marketOverviewHandler.GetMarketOverviewTable)
	marketOverviewGroup.Get("/token-detail", marketOverviewHandler.GetTokenDetail)
}
