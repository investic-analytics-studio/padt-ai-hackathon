package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
	"github.com/quantsmithapp/datastation-backend/internal/middleware"
)

func BindApiV2(router fiber.Router) {
	v2 := router.Group("/v2")

	config := config.GetConfig()
	authRepo := repo.NewAuthRepo(infra.CryptoDB)
	jwtService := service.NewJwtService(infra.FirebaseClient, &config, authRepo)
	authMiddleware := middleware.AuthMiddleware(authRepo, jwtService)
	authCRMMiddleware := middleware.AuthCRMMiddleware(authRepo, jwtService)
	bindInitPageStatus(v2)

	bindMarketSummaryAPI(v2, authMiddleware)
	bindTwitterCryptoAPI(v2, authMiddleware)
	bindSentimentCryptoAPI(v2, authMiddleware)
	bindTradingViewAPI(v2, authMiddleware)
	bindAuth(v2, authRepo, authMiddleware, &config)
	bindGetTagAPI(v2)
	bindTopRankAPI(v2)
	bindWinRateAPI(v2, authMiddleware)
	bindMarketCapAPI(v2, authMiddleware)
	bindNewsSentimentCryptoAPI(v2)
	bindSentimentAnalysisRouter(v2, authMiddleware)
	bindMarketOverviewAPI(v2)
	bindCryptoUserRefcodeAPI(v2, authMiddleware)
	bindFeatures(v2, authMiddleware)
	bindCryptoCRMAPI(v2, authRepo, authCRMMiddleware, &config)
	bindCryptoNotificationAPI(v2, authMiddleware, &config)
	bindTelegramAPI(v2, &config)
	bindDexAPI(v2, authMiddleware)
	bindCexAPI(v2, authMiddleware)
	bindPerformanceAPI(v2, authMiddleware)
}
