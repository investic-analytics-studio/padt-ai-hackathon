package v2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/handler"
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
)

// func bindGetTagAPI(router fiber.Router, authMiddleware fiber.Handler)
func bindGetTagAPI(router fiber.Router) {

	extractTagRepo := repo.NewExtractTagRepo(infra.PostgresDB)
	extractTagService := service.NewExtractTagService(extractTagRepo)
	extractTagHandler := handler.NewExtractTagHandler(extractTagService)

	// extractTagGroup := router.Group("/extract-tags", authMiddleware)
	extractTagGroup := router.Group("/extract-tags")
	extractTagGroup.Get("/", extractTagHandler.GetAllTags)
}
