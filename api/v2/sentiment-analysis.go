package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindSentimentAnalysisRouter(router fiber.Router, authMiddleware fiber.Handler) {
	db, err := infra.GetPostgresConnection()
	if err != nil {
		panic(err)
	}
	sentimentAnalysisRepo := repo.NewSentimentAnalysisRepo(db)
	sentimentAnalysisService := service.NewSentimentAnalysisService(sentimentAnalysisRepo)
	sentimentAnalysisHandler := handler.NewSentimentAnalysisHandler(sentimentAnalysisService)
	// sentimentAnalysisRouter := router.Group("/sentiment-analysis", authMiddleware)
	_ = authMiddleware
	sentimentAnalysisRouter := router.Group("/sentiment-analysis")
	sentimentAnalysisRouter.Post("/author-list", sentimentAnalysisHandler.GetSentimentAnalysisByAuthorList)
	sentimentAnalysisRouter.Post("/tier", sentimentAnalysisHandler.GetSentimentAnalysisByTier)
}
