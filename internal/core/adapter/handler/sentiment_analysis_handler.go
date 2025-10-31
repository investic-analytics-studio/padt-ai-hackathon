package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type sentimentAnalysisHandler struct {
	service port.SentimentAnalysisService
}

func NewSentimentAnalysisHandler(service port.SentimentAnalysisService) *sentimentAnalysisHandler {
	return &sentimentAnalysisHandler{service: service}
}

func (h *sentimentAnalysisHandler) GetSentimentAnalysisByAuthorList(c *fiber.Ctx) error {
	var body model.SentimentAnalysisByAuthorListRequest
	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sentimentAnalysis, err := h.service.GetSentimentAnalysisByAuthorList(body.AuthorList, body.DateRange)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(sentimentAnalysis)
}

func (h *sentimentAnalysisHandler) GetSentimentAnalysisByTier(c *fiber.Ctx) error {
	var body model.SentimentAnalysisByTierRequest
	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sentimentAnalysis, err := h.service.GetSentimentAnalysisByTier(body.Tier, body.DateRange)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(sentimentAnalysis)
}
