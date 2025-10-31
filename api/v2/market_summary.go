package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

func bindMarketSummaryAPI(router fiber.Router, authMiddleware fiber.Handler) {
	logger := logger.NewLogger()
	repo := repo.NewMarketSummaryRepo(infra.StockDB)
	serv := service.NewMarketSummaryService(repo, logger)
	hdl := handler.NewMarketSummaryHandler(serv)

	router.Group("/market-summary", authMiddleware).
		Get("/daily-overview", hdl.GetDailyMarketOverview).
		Get("/top-gainer", hdl.GetTopGainer).
		Get("/top-loser", hdl.GetTopLoser).
		Get("/add-hist", hdl.GetAdvancerDeclinerDistributionHist).
		Get("/add-bar", hdl.GetAdvancerDeclinerDistributionBar).
		Get("/top-turnover", hdl.GetTopTurnoverFloat).
		Get("/cme-gold-oi", hdl.GetCMEGoldOI).
		Get("/stock-alert-plots", hdl.GetStockAlertPlots).
		Get("/stock-alert-stats", hdl.GetStockAlertStats).
		Get("/stock-alert-detections", hdl.GetStockAlertDetections).
		Get("/crypto-alert-plots-1d", hdl.GetCryptoAlertPlots1D).
		Get("/crypto-alert-stats-dates-1d", hdl.GetCryptoAlertStatsDates1D).
		Get("/crypto-alert-stats-1d", hdl.GetCryptoAlertStats1D).
		Get("/crypto-alert-detections-1d", hdl.GetCryptoAlertDetections1D).
		Get("/crypto-alert-plots", hdl.GetCryptoAlertPlots).
		Get("/crypto-alert-stats", hdl.GetCryptoAlertStats).
		Get("/crypto-alert-detections", hdl.GetCryptoAlertDetections).
		Get("/all-crypto-alert-stats", hdl.GetAllCryptoAlertStats).
		Get("/all-crypto-alert-stats-1d", hdl.GetAllCryptoAlertStats1D)
}
