package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type NewsSentimentCryptoHandler struct {
	service port.NewsSentimentCryptoService
}

func NewNewsSentimentCryptoHandler(service port.NewsSentimentCryptoService) *NewsSentimentCryptoHandler {
	return &NewsSentimentCryptoHandler{service: service}
}

func (h *NewsSentimentCryptoHandler) GetNewsSentiment(c *fiber.Ctx) error {
	newsSentiment, err := h.service.GetNewsSentiment(c.Context(), time.Now().UTC())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"bubble_sentiment": newsSentiment,
	})
}
