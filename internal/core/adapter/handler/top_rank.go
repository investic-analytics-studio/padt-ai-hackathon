package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type TopRankHandler struct {
	service port.TopRankService
}

func NewTopRankHandler(service port.TopRankService) *TopRankHandler {
	return &TopRankHandler{service: service}
}

func (h *TopRankHandler) GetTop100(c *fiber.Ctx) error {
	top100, err := h.service.GetTop100()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(top100)
}
