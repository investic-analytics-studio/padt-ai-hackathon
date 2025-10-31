package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type MarketCapHandler struct {
	service port.MarketCapService
}

func NewMarketCapHandler(service port.MarketCapService) *MarketCapHandler {
	return &MarketCapHandler{service: service}
}

func (h *MarketCapHandler) GetMarketCap(c *fiber.Ctx) error {
	marketCap, err := h.service.GetMarketCap()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(marketCap)
}
