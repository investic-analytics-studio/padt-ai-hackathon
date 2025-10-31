package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindTelegramAPI(router fiber.Router, config *config.Config) {
	telegramService, err := service.NewTelegramService(config.Telegram.BotToken)
	if err != nil {
		return
	}

	telegramHandler := handler.NewTelegramHandler(telegramService)
	router.Post("/telegram/send-message", telegramHandler.HandleSendMessage)
}
