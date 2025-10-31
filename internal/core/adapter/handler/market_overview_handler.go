package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/pkg/util"
)

type MarketOverviewHandler struct {
	marketOverviewService port.MarketOverviewService
}

func NewMarketOverviewHandler(marketOverviewService port.MarketOverviewService) *MarketOverviewHandler {
	return &MarketOverviewHandler{
		marketOverviewService: marketOverviewService,
	}
}

// GetMarketOverviewHandle godoc
// @Summary      Get merket overview data
// @Description  Retrieves the market data
// @Tags         market-overview
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /market-overview/sentiment-market-overview [get]
func (h *MarketOverviewHandler) GetMarketOverviewTable(c *fiber.Ctx) error {
	tier := c.Query("tier", "ALL")
	timeframe := c.Query("timeframe", "30")

	tier = strings.ToUpper(tier)
	// validate tier
	if tier != "ALL" && tier != "SSS" && tier != "S" && tier != "A" && tier != "B" && tier != "C" && tier != "D" && tier != "OTHER" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("invalid tier: %s. Must be one of: all, s, a, b, c, d, other", tier),
		})
	}

	// validate timeframe should be int
	_, err := strconv.Atoi(timeframe)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("invalid timeframe: %s. Must be an integer", timeframe),
		})
	}

	payload := requestPayload{
		Columns:             []string{"base_currency", "exchange", "currency", "close", "24h_close_change|5", "market_cap_calc", "24h_vol_cmc"},
		Markets:             []string{"coin"},
		IgnoreUnknownFields: false,
		Range:               []int{0, 300},
		Sort: struct {
			SortBy    string `json:"sortBy"`
			SortOrder string `json:"sortOrder"`
		}{
			SortBy:    "crypto_total_rank",
			SortOrder: "asc",
		},
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to encode request: %w", err),
		})
	}

	response, err := http.Post(
		"https://scanner.tradingview.com/coin/scan?label-product=screener-coin",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to make request: %w", err),
		})
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to read response body: %w", err),
		})
	}

	var scannerResult map[string]interface{}
	if err := json.Unmarshal(body, &scannerResult); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to unmarshal response: %w", err),
		})
	}

	enrichedData, err := h.marketOverviewService.GetEnrichedMarketOverviewTable(tier, timeframe, scannerResult)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(enrichedData)
}

func (h *MarketOverviewHandler) GetMarketOverview(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)
	tier := c.Query("tier", "all")
	marketOverview, err := h.marketOverviewService.GetMarketOverview(days, tier)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return util.ResponseOK(c, marketOverview)
}

// GetTokenDetail handles requests for token-specific market overview data
func (h *MarketOverviewHandler) GetTokenDetail(c *fiber.Ctx) error {
	// Get query parameters
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "token parameter is required",
		})
	}

	// Parse days parameter with default value of 7
	days := 7
	daysStr := c.Query("days")
	if daysStr != "" {
		var err error
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid days parameter, must be a number",
			})
		}

		// Validate days range
		if days < 1 || days > 365 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "days parameter must be between 1 and 365",
			})
		}
	}

	// Call the service to get token detail data
	tokenDetails, err := h.marketOverviewService.GetTokenDetailMarketOverview(token, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get token detail: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tokenDetails,
	})
}
