package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindAuth(router fiber.Router, authRepo port.AuthRepo, authMiddleware fiber.Handler, config *config.Config) {
	jwtService := service.NewJwtService(infra.FirebaseClient, config, authRepo)
	authService := service.NewAuthService(infra.FirebaseClient, authRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService, jwtService)
	auth := router.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/exist-email", authHandler.ExistEmail)
	auth.Get("/get-me", authMiddleware, authHandler.GetMe)
	auth.Post("/auto-validate-email-in-firebase", authHandler.AutoValidateEmailInFirebase)
	auth.Patch("/all-user-auto-validate-email-in-firebase", authHandler.AllUserAutoValidateEmailInFirebase)
	auth.Post("/login", authHandler.Login)
	auth.Post("/exist-twitter-uid", authHandler.ExistTwitterUID)

}
