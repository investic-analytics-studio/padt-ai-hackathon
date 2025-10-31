package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

func bindFeatures(router fiber.Router, authMiddleware fiber.Handler) {
	featuresRepo := repo.NewFeaturesOverviewRepo(infra.CryptoDB)
	featuresService := service.NewFeaturesService(featuresRepo)
	featuresHandler := handler.NewFeaturesHandler(featuresService)

	features := router.Group("/features")
	features.Get("/:featureName", featuresHandler.GetFeatureHandle)
}
