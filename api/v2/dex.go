package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindDexAPI(router fiber.Router, authMiddleware fiber.Handler) {
	dexRepo := repo.NewDexRepo(infra.CryptoDB)
	dexService := service.NewDexService(dexRepo)
	dexHandler := handler.NewDexHandler(dexService)

	router.Post("/dex/add-wallet", authMiddleware, dexHandler.Connect)
	router.Get("/dex/wallet-info", authMiddleware, dexHandler.ListWallets)
	router.Get("/dex/wallet-total-value", authMiddleware, dexHandler.GetWalletTotalValue)
	router.Post("/dex/active-wallet", authMiddleware, dexHandler.ActiveWallet)
	router.Post("/dex/deactive-wallet", authMiddleware, dexHandler.DeactiveWallet)
	router.Post("/dex/update-position-size", authMiddleware, dexHandler.UpdatePositionSize)
	router.Post("/dex/update-leverage", authMiddleware, dexHandler.UpdateLeverage)
	router.Post("/dex/update-api-credentials", authMiddleware, dexHandler.UpdateAPICredentials)
	router.Post("/dex/update-sl", authMiddleware, dexHandler.UpdateSL)
	router.Post("/dex/subscribe-author", authMiddleware, dexHandler.SubscribeAuthor)
	router.Post("/dex/unsubscribe-author", authMiddleware, dexHandler.UnsubscribeAuthor)
}
