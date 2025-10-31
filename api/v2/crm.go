package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

// type GetCryptoRefUserResponse map[string]string

func bindCryptoCRMAPI(router fiber.Router, authRepo port.AuthRepo, AuthCRMMiddleware fiber.Handler, config *config.Config) {
	jwtService := service.NewJwtService(infra.FirebaseClient, config, authRepo)
	crmRepo := repo.NewCRMRepo(infra.CryptoDB)
	crmService := service.NewCryptoCRMService(crmRepo, jwtService)
	crmHandler := handler.NewCryptoCRMHandler(crmService)

	router.Post("/crm/login", crmHandler.Login)
	router.Post("/crm/new-kol", AuthCRMMiddleware, crmHandler.NewKolUser)
	router.Get("/crm/kol_refer_detail", AuthCRMMiddleware, crmHandler.KolReferDetail)
	router.Get("/crm/get_users", AuthCRMMiddleware, crmHandler.GetAllUsers)
	router.Post("/crm/update_display_code", AuthCRMMiddleware, crmHandler.UpdateDisplayCode)

	router.Get("/crm/get-crypto-user", AuthCRMMiddleware, crmHandler.GetCryptoUser)
	router.Post("/crm/update-crypto-user-approve", AuthCRMMiddleware, crmHandler.UpdateCryptoUserApprove)
	router.Get("/crm/get_refferal_score", AuthCRMMiddleware, crmHandler.GetRefferalScore)
	router.Get("/crm/get_user_referral", AuthCRMMiddleware, crmHandler.GetUserReferral)
	router.Get("/crm/privy-user-overview", AuthCRMMiddleware, crmHandler.GetPrivyUserOverview)

}
