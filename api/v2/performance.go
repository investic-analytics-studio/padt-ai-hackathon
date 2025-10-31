package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

// type GetCryptoRefUserResponse map[string]string

func bindPerformanceAPI(router fiber.Router, authMiddleware fiber.Handler) {
	db, err := infra.GetPostgresConnection()
	if err != nil {
		panic(err)
	}
	authorTierRepo := repo.NewAuthorTierRepo(db)
	performanceRepo := repo.NewPerformanceRepo(infra.CryptoDB, db, authorTierRepo)
	performanceService := service.NewPerformanceService(performanceRepo)
	performanceHandler := handler.NewPerformanceHandler(performanceService)
	router.Get("/performance/get-author-nav", performanceHandler.GetAuthorNavHandle)
	router.Get("/performance/get-multiholding-port-nav", performanceHandler.GetMultiholdingPortNavHandle)
	router.Get("/performance/get-author-detail", performanceHandler.GetAuthorDetailHandle)
	router.Get("/performance/get-author-multiholding-detail", performanceHandler.GetAuthorMultiholdingDetailHandle)
}
