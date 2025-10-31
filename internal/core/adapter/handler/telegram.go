package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type TelegramHandler struct {
	telegramService port.TelegramService
}

func NewTelegramHandler(telegramService port.TelegramService) *TelegramHandler {
	return &TelegramHandler{
		telegramService: telegramService,
	}
}

// SendMessageRequest represents the request body for sending a message
type SendMessageRequest struct {
	Text   string `json:"text"`
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
}

// HandleSendMessage handles the POST request to send a message
func (h *TelegramHandler) HandleSendMessage(c *fiber.Ctx) error {
	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Text == "" || req.ChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Text and chat_id are required",
		})
	}

	chatID, err := strconv.ParseInt(req.ChatID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat_id format",
		})
	}

	if err := h.telegramService.SendMessage(c.Context(), req.Text, chatID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send message",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
