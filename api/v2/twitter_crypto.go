package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindTwitterCryptoAPI(router fiber.Router, authMiddleware fiber.Handler) {
	db, err := infra.GetPostgresConnection()
	if err != nil {
		panic(err)
	}

	twitterCryptoRepo := repo.NewTwitterCryptoRepo(db)
	authorTierRepo := repo.NewAuthorTierRepo(db)
	twitterCryptoService := service.NewTwitterCryptoService(twitterCryptoRepo)
	authorTierService := service.NewAuthorTierService(authorTierRepo)
	twitterCryptoHandler := handler.NewTwitterCryptoHandler(twitterCryptoService, authorTierService)
	// twitterCrypto := router.Group("/twitter-crypto", authMiddleware)
	twitterCrypto := router.Group("/twitter-crypto")
	twitterCrypto.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Twitter Crypto API is working!")
	})
	twitterCrypto.Get("/sentiments", twitterCryptoHandler.GetAllSentiments)
	twitterCrypto.Get("/tweets", twitterCryptoHandler.GetAllTweets)
	twitterCrypto.Get("/author-profiles", twitterCryptoHandler.GetAuthorProfiles)
	twitterCrypto.Get("/author-winrate", authMiddleware, twitterCryptoHandler.GetAuthorWinrate)
	twitterCrypto.Get("/paginated-tweets", twitterCryptoHandler.GetPaginatedTweets)
	twitterCrypto.Get("/paginated-sentiments", twitterCryptoHandler.GetPaginatedSentiments)
	twitterCrypto.Get("/tweets-with-sentiments", twitterCryptoHandler.GetTweetsWithSentiments)
	twitterCrypto.Get("/tweets-with-sentiment-author-signal", twitterCryptoHandler.GetTweetsWithSentimentAuthorSignal)
	twitterCrypto.Get("/tweets-with-sentiments-and-author", twitterCryptoHandler.GetTweetsWithSentimentsAndAuthor)
	twitterCrypto.Get("/tweets-summaries-1h", twitterCryptoHandler.GetSummaries)
	twitterCrypto.Get("/bubble-chart-data", twitterCryptoHandler.GetBubbleSentiment)
	twitterCrypto.Get("/search-token-mention-symbols-by-authors", twitterCryptoHandler.SearchTokenMentionSymbolsByAuthors)
	twitterCrypto.Get("/tiers", twitterCryptoHandler.GetAllTiers)
	twitterCrypto.Get("/tweets-with-sentiments-and-tier", twitterCryptoHandler.GetTweetsWithSentimentsAndTier)
}
