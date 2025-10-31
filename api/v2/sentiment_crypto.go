package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

func bindSentimentCryptoAPI(router fiber.Router, authMiddleware fiber.Handler) {
	logger := logger.NewLogger()
	sentimentRepo := repo.NewSentimentCryptoRepo(infra.StockDB, logger)
	sentimentService := service.NewSentimentCryptoService(sentimentRepo, logger)
	sentimentHandler := handler.NewSentimentCryptoHandler(sentimentService, logger)

	// sentimentCrypto := router.Group("/sentiment-crypto", authMiddleware)
	sentimentCrypto := router.Group("/sentiment-crypto")
	sentimentCrypto.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Sentiment Crypto API is working!")
	})
	sentimentCrypto.Get("/aggregateTotal", sentimentHandler.GetSentimentAggregateCountRow)
	sentimentCrypto.Get("/aggregate", sentimentHandler.GetSentimentAggregate)
	sentimentCrypto.Get("/unique-tickers", sentimentHandler.GetUniqueTickers)
	sentimentCrypto.Get("/unique-topics", authMiddleware, sentimentHandler.GetUniqueTopics)
	sentimentCrypto.Get("/unique-source-names", authMiddleware, sentimentHandler.GetUniqueSourceNames)
}
