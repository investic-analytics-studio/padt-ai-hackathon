package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindCryptoNotificationAPI(router fiber.Router, authMiddleware fiber.Handler, config *config.Config) {
	notificationRepo := repo.NewCryptoNotificationRepo(infra.CryptoDB)
	notificationService := service.NewCryptoNotificationService(notificationRepo)
	notificationHandler := handler.NewCryptoNotificationHandler(notificationService, config)
	router.Get("/notification/get-group", authMiddleware, notificationHandler.GetNotificationGroupList)
	router.Post("/notification/update-group-name", authMiddleware, notificationHandler.UpdateGroupName)
	router.Post("/notification/add-group", authMiddleware, notificationHandler.AddGroup)
	router.Get("/notification/count-group", authMiddleware, notificationHandler.CountGroup)
	router.Post("/notification/add-author", authMiddleware, notificationHandler.AddAuthor)
	router.Post("/notification/remove-author", authMiddleware, notificationHandler.RemoveAuthor)
	router.Post("/notification/update-telegram", authMiddleware, notificationHandler.UpdateTelegram)
	router.Get("/notification/get-telegram", authMiddleware, notificationHandler.GetTelegram)
	router.Post("/notification/delete-group-author", authMiddleware, notificationHandler.DeleteGroupAuthor)
	router.Get("/notification/disconnect", authMiddleware, notificationHandler.DisconnectNotification)
}
