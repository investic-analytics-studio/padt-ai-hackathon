package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type FeaturesHandler struct {
	service port.FeaturesService
}
type GetFeatureResponse struct {
	FeatureName   string `json:"feature_name" example:"refcode"`
	FeatureEnable bool   `json:"feature_enable" example:"true"`
}

func NewFeaturesHandler(service port.FeaturesService) *FeaturesHandler {
	return &FeaturesHandler{service: service}
}

// GetFeatureHandle godoc
// @Summary      Get feature by name
// @Description  Retrieves the status of a feature by its name
// @Tags         features
// @Accept       json
// @Produce      json
// @Param        featureName path string true "Feature Name" default(refcode)
// @Success      200 {object} GetFeatureResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /features/{featureName} [get]
func (h *FeaturesHandler) GetFeatureHandle(c *fiber.Ctx) error {
	featureName := c.Params("featureName")
	if featureName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing feature name",
		})
	}

	feature, err := h.service.GetFeature(featureName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := GetFeatureResponse{
		FeatureName:   feature.FeatureName,
		FeatureEnable: feature.FeatureEnable,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
