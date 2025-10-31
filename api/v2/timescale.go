package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

func bindTimescaleAPI(router fiber.Router, authMiddleware fiber.Handler) {
	logger := logger.NewLogger()
	db, err := infra.GetTimescaleDBConnection()
	if err != nil {
		logger.Fatal(err)
	}
	repo := repo.NewTimescaleRepo(db)
	serv := service.NewOHLCVDataService(repo, logger)
	hdl := handler.NewOHLCVDataHandler(serv)

	timescaleGroup := router.Group("/timescale", authMiddleware)
	timescaleGroup.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Timescale API is working!")
	})
	timescaleGroup.Get("/crypto", hdl.GetCryptoOHLCV)
	timescaleGroup.Get("/forex", hdl.GetForexOHLCV)
}
