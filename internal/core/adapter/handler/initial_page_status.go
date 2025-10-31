package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
)

type (
	initialPageStatusHandler struct {
		cfg *config.PageShowConfig
	}
)

func NewInitialPageStatus(cfg *config.PageShowConfig) *initialPageStatusHandler {
	return &initialPageStatusHandler{
		cfg: cfg,
	}
}

func (h *initialPageStatusHandler) GetInitialPageStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"crypto_lite_page":     h.cfg.CryptoLitePage,
		"twitter_page":         h.cfg.TwitterPage,
		"sentiment_page":       h.cfg.SentimentPage,
		"stats_page":           h.cfg.StatsPage,
		"sector_page":          h.cfg.SectorPage,
		"overview_market_page": h.cfg.OverviewMarketPage,
		"copy_trade_page":      h.cfg.CopyTradePage,
		"genesis_page":         h.cfg.GenesisPage,
	})
}
