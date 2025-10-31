package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type ExtractTagHandler struct {
	service port.ExtractTagService
}

func NewExtractTagHandler(service port.ExtractTagService) *ExtractTagHandler {
	return &ExtractTagHandler{service: service}
}

func (h *ExtractTagHandler) GetAllTags(c *fiber.Ctx) error {
	tags, err := h.service.GetAllTags()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	listTags, err := h.service.GetUniqueTags()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tags":      tags.CoinWithTags,
		"list_tags": listTags,
	})
}
