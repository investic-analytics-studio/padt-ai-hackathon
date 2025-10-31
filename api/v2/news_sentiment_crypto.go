package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

// func bindNewsSentimentCryptoAPI(router fiber.Router, authMiddleware fiber.Handler)
func bindNewsSentimentCryptoAPI(router fiber.Router) {
	db, err := infra.GetPostgresConnection()
	if err != nil {
		panic(err)
	}

	newsSentimentRepo := repo.NewNewsSentimentCryptoRepo(db)
	newsSentimentService := service.NewNewsSentimentCryptoService(newsSentimentRepo)
	newsSentimentHandler := handler.NewNewsSentimentCryptoHandler(newsSentimentService)
	newsSentiment := router.Group("/news-sentiment-crypto")
	// newsSentiment := router.Group("/news-sentiment-crypto", authMiddleware)
	newsSentiment.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("News Sentiment Crypto API is working!")
	})
	newsSentiment.Get("/news-sentiment", newsSentimentHandler.GetNewsSentiment)
}
