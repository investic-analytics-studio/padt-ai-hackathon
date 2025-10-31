package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type winRateHandler struct {
	winRateService port.WinRateService
}

func NewWinRateHandler(winRateService port.WinRateService) *winRateHandler {
	return &winRateHandler{winRateService: winRateService}
}

func (h *winRateHandler) GetWinRate(c *fiber.Ctx) error {
	winRate, err := h.winRateService.GetWinRate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(winRate)
}
