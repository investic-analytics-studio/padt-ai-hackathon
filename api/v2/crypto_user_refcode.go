package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

// type GetCryptoRefUserResponse map[string]string

func bindCryptoUserRefcodeAPI(router fiber.Router, authMiddleware fiber.Handler) {
	refcodeRepo := repo.NewCryptoUserRefcodeRepo(infra.CryptoDB)
	refcodeService := service.NewCryptoUserRefcodeService(refcodeRepo)
	authRepo := repo.NewAuthRepo(infra.CryptoDB)
	refcodeHandler := handler.NewCryptoUserRefcodeHandler(refcodeService, authRepo)

	router.Post("/generate-refcode", authMiddleware, refcodeHandler.GenerateRefcodeRequest)
	router.Get("/check-user-id", authMiddleware, refcodeHandler.CheckUserIDExists)
	router.Post("/check-refcode", authMiddleware, refcodeHandler.CheckAndUpdateRefcode, refcodeHandler.CheckAndInsertKolcode)
	router.Get("/get-crypto-ref-user", authMiddleware, refcodeHandler.GetCryptoRefUser)
	router.Get("/get-kolcode", authMiddleware, refcodeHandler.GetCryptoKolCode)
	router.Post("/generate-refcode-bynum", authMiddleware, refcodeHandler.GenerateRefcodeBynumRequest)
	router.Get("/get-refferal-score", refcodeHandler.GetRefferalScore)
	router.Post("/get-refferal-score-ranking", refcodeHandler.GetRefferalScoreRanking)
	router.Post("/check-xuser", refcodeHandler.CheckXUserIsExit)
}
