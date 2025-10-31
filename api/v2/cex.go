package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindCexAPI(router fiber.Router, authMiddleware fiber.Handler) {
	cexRepo := repo.NewCexRepo(infra.CryptoDB)
	cexService := service.NewCexService(cexRepo)
	cexHandler := handler.NewCexHandler(cexService)

	router.Post("/cex/add-wallet", authMiddleware, cexHandler.Connect)
	router.Get("/cex/wallet-info", authMiddleware, cexHandler.ListWallets)
	router.Get("/cex/wallet-total-value", authMiddleware, cexHandler.GetWalletTotalValue)
	router.Post("/cex/subscribe-author", authMiddleware, cexHandler.SubscribeAuthor)
	router.Post("/cex/unsubscribe-author", authMiddleware, cexHandler.UnsubscribeAuthor)
	router.Post("/cex/active-wallet", authMiddleware, cexHandler.ActiveWallet)
	router.Post("/cex/deactive-wallet", authMiddleware, cexHandler.DeactiveWallet)
	router.Post("/cex/update-position-size", authMiddleware, cexHandler.UpdatePositionSize)
	router.Post("/cex/update-holding-period", authMiddleware, cexHandler.UpdateHoldingPeriod)
	router.Post("/cex/update-api-key", authMiddleware, cexHandler.UpdateAPIKey)
	router.Post("/cex/update-sl", authMiddleware, cexHandler.UpdateSL)
}
