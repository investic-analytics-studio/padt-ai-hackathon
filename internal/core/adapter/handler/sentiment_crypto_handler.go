package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
	"go.uber.org/zap"
)

type SentimentCryptoHandler struct {
	service domain.SentimentCryptoService
	logger  logger.Logger
}

func NewSentimentCryptoHandler(service domain.SentimentCryptoService, logger logger.Logger) *SentimentCryptoHandler {
	return &SentimentCryptoHandler{
		service: service,
		logger:  logger,
	}
}
func (h *SentimentCryptoHandler) GetSentimentAggregateCountRow(c *fiber.Ctx) error {
	h.logger.Info("GetSentimentAggregateCountRow handler called")
	ctx := context.Background()

	ticker := c.Query("ticker", "")
	timeRange := c.Query("range", "all")
	sourceNames := c.Query("sources", "")
	topics := c.Query("topics", "")

	h.logger.Info("Request parameters",
		zap.String("ticker", ticker),
		zap.String("time_range", timeRange),
		zap.String("sources", sourceNames),
		zap.String("topics", topics),
	)

	totalRow, err := h.service.GetSentimentAggregateCountRow(ctx, ticker, timeRange, sourceNames, topics)

	if err != nil {
		h.logger.Error(err, zap.String("method", "GetSentimentAggregateCountRow"))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve sentiment aggregate data",
		})
	}

	h.logger.Info("GetSentimentAggregateCountRow handler successful", zap.Int("data_count_row", totalRow))

	return c.JSON(fiber.Map{
		"totalRow": totalRow,
	})
}
func (h *SentimentCryptoHandler) GetSentimentAggregate(c *fiber.Ctx) error {
	h.logger.Info("GetSentimentAggregate handler called")
	ctx := context.Background()

	ticker := c.Query("ticker", "")
	timeRange := c.Query("range", "all")
	sourceNames := c.Query("sources", "")
	topics := c.Query("topics", "")
	limit := c.Query("limit", "")
	offset := c.Query("offset", "")

	h.logger.Info("Request parameters",
		zap.String("ticker", ticker),
		zap.String("time_range", timeRange),
		zap.String("sources", sourceNames),
		zap.String("topics", topics),
		zap.String("limit", limit),
		zap.String("offset", offset))

	data, err := h.service.GetSentimentAggregate(ctx, ticker, timeRange, sourceNames, topics, limit, offset)

	if err != nil {
		h.logger.Error(err, zap.String("method", "GetSentimentAggregate"))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve sentiment aggregate data",
		})
	}

	if len(data) == 0 {
		h.logger.Info("No data found for the given parameters")
		return c.JSON([]model.SentimentCrypto{})
	}

	h.logger.Info("GetSentimentAggregate handler successful", zap.Int("data_count", len(data)))

	return c.JSON(data)
}

func (h *SentimentCryptoHandler) GetUniqueTickers(c *fiber.Ctx) error {
	h.logger.Info("GetUniqueTickers handler called")
	ctx := context.Background()

	tickers, err := h.service.GetUniqueTickers(ctx)
	if err != nil {
		h.logger.Error(err, zap.String("method", "GetUniqueTickers"))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve unique tickers",
		})
	}

	h.logger.Info("GetUniqueTickers handler successful", zap.Int("tickers_count", len(tickers)))
	return c.JSON(tickers)
}

func (h *SentimentCryptoHandler) GetUniqueTopics(c *fiber.Ctx) error {
	h.logger.Info("GetUniqueTopics handler called")
	ctx := context.Background()

	topics, err := h.service.GetUniqueTopics(ctx)
	if err != nil {
		h.logger.Error(err, zap.String("method", "GetUniqueTopics"))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve unique topics",
		})
	}

	h.logger.Info("GetUniqueTopics handler successful", zap.Int("topics_count", len(topics)))
	return c.JSON(topics)
}

func (h *SentimentCryptoHandler) GetUniqueSourceNames(c *fiber.Ctx) error {
	h.logger.Info("GetUniqueSourceNames handler called")
	ctx := context.Background()

	sourceNames, err := h.service.GetUniqueSourceNames(ctx)
	if err != nil {
		h.logger.Error(err, zap.String("method", "GetUniqueSourceNames"))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve unique source names",
		})
	}

	h.logger.Info("GetUniqueSourceNames handler successful", zap.Int("source_names_count", len(sourceNames)))
	return c.JSON(sourceNames)
}
