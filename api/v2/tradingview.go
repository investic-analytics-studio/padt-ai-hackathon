package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindTradingViewAPI(router fiber.Router, authMiddleware fiber.Handler) {
	priceService := service.NewPriceService()
	hdl := handler.NewTradingViewHandler(priceService)

	_ = authMiddleware
	tv := router.Group("/tradingview")
	tv.Get("/historical", hdl.GetHistoricalData)
	tv.Post("/historical/multi", hdl.GetMultiHistoricalData)
	tv.Get("/search", hdl.SearchSymbol)
	tv.Post("/get-24-hour-price-change", hdl.Get24HoursPriceChangeTradingView)
}
