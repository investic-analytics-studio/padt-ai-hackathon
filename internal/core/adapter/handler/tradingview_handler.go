package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/quantsmithapp/datastation-backend/internal/core/service"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/internal/tradingview"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type TradingViewHandler struct {
	logger       logger.Logger
	client       *tradingview.TradingViewClient
	priceService *service.PriceService
}

// Request structure for multiple symbols
type MultiSymbolRequest struct {
	Symbols  []SymbolRequest `json:"symbols"`
	Interval string          `json:"interval"`
	Bars     string          `json:"bars"`
	From     string          `json:"from,omitempty"`
	To       string          `json:"to,omitempty"`
}

type SymbolRequest struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

type requestPayload struct {
	Columns             []string `json:"columns"`
	Markets             []string `json:"markets"`
	Range               []int    `json:"range"`
	IgnoreUnknownFields bool     `json:"ignore_unknown_fields"`
	Sort                struct {
		SortBy    string `json:"sortBy"`
		SortOrder string `json:"sortOrder"`
	} `json:"sort"`
}

func NewTradingViewHandler(priceService *service.PriceService) *TradingViewHandler {
	// Initialize the client when creating the handler
	client, err := tradingview.NewTradingViewClient("", "")
	if err != nil {
		// Handle error, perhaps panic or return nil
		panic(err)
	}

	return &TradingViewHandler{
		logger:       logger.NewLogger(),
		client:       client,
		priceService: priceService,
	}
}

func (h *TradingViewHandler) GetHistoricalData(c *fiber.Ctx) error {
	symbol := c.Query("symbol", "BTCUSDT")
	exchange := c.Query("exchange", "BINANCE")
	interval := c.Query("interval", "1D")
	bars := c.QueryInt("bars", 100)
	extendedSession := c.QueryBool("extended_session", false)

	var futContract *int
	if futContractStr := c.Query("fut_contract"); futContractStr != "" {
		if val, err := strconv.Atoi(futContractStr); err == nil {
			futContract = &val
		}
	}

	h.logger.Info(fmt.Sprintf(
		"Getting historical data (REST) for symbol: %s, exchange: %s, interval: %s, bars: %d",
		symbol, exchange, interval, bars,
	))

	// Remove the client creation here since we now have it as a field
	data, err := h.client.GetHistoricalData(
		symbol,
		exchange,
		interval,
		bars,
		futContract,
		extendedSession,
	)
	if err != nil {
		h.logger.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

func (h *TradingViewHandler) SearchSymbol(c *fiber.Ctx) error {
	text := c.Query("text", "")
	exchange := c.Query("exchange", "")
	start, _ := strconv.Atoi(c.Query("start", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	params := tradingview.SearchParams{
		Text:          text,
		Exchange:      exchange,
		Start:         start,
		Limit:         limit,
		SearchType:    "crypto",
		Lang:          "en",
		Domain:        "production",
		SortByCountry: "US",
	}

	response, err := h.client.SearchSymbol(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return in the format your frontend expects
	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"symbols":           response.Symbols,
			"symbols_remaining": response.SymbolsRemaining,
		},
	})
}

func (h *TradingViewHandler) GetMultiHistoricalData(c *fiber.Ctx) error {
	var req MultiSymbolRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	attemptsCount := 2
	// Validate and process request parameters
	interval := req.Interval
	if interval == "" {
		interval = "1D" // default interval
	}

	// Convert bars from string to int
	bars, err := strconv.Atoi(req.Bars)
	if err != nil {
		bars = 10 // default bars
	}

	// Parse time range if provided
	var fromTime, toTime time.Time
	var useTimeRange bool

	if req.From != "" && req.To != "" {
		fromTime, err = time.Parse(time.RFC3339, req.From)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid 'from' date format",
			})
		}

		toTime, err = time.Parse(time.RFC3339, req.To)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid 'to' date format",
			})
		}

		useTimeRange = true
	}

	// Create channels for results
	type result struct {
		key  string
		data []tradingview.HistoricalData
	}
	resultChan := make(chan result, len(req.Symbols))

	var wg sync.WaitGroup

	// Process each symbol request concurrently
	for _, symbolReq := range req.Symbols {
		wg.Add(1)
		go func(sym SymbolRequest) {
			defer wg.Done()

			var data []tradingview.HistoricalData

			// Try up to 5 times
			for attempts := 0; attempts < attemptsCount; attempts++ {
				client, err := tradingview.NewTradingViewClient("", "")
				if err != nil {
					h.logger.Error(fmt.Errorf("attempt %d: failed to create client for %s:%s: %w",
						attempts+1, sym.Symbol, sym.Exchange, err))
					time.Sleep(time.Second)
					continue
				}

				// Get historical data
				data, err = client.GetHistoricalData(
					sym.Symbol,
					sym.Exchange,
					interval,
					bars,
					nil,
					false,
				)

				if err == nil && data != nil {
					// Filter data by time range if specified
					if useTimeRange && data != nil {
						filteredData := make([]tradingview.HistoricalData, 0)
						for _, d := range data {
							if (d.DateTime.Equal(fromTime) || d.DateTime.After(fromTime)) &&
								(d.DateTime.Equal(toTime) || d.DateTime.Before(toTime)) {
								filteredData = append(filteredData, d)
							}
						}
						data = filteredData
					}
					break
				}

				h.logger.Error(fmt.Errorf("attempt %d: failed to get data for %s:%s: %w",
					attempts+1, sym.Symbol, sym.Exchange, err))
				time.Sleep(time.Second)
			}

			// last attempt
			if data == nil {
				client, err := tradingview.NewTradingViewClient("", "")
				if err != nil {
					h.logger.Error(fmt.Errorf("attempt %d: failed to create client for %s:%s: %w",
						attemptsCount+1, sym.Symbol, sym.Exchange, err))
					time.Sleep(time.Second)
				}

				log.Printf("Last attempt for %s:%s", fmt.Sprintf("%s.P", sym.Symbol), sym.Exchange)
				// Get historical data
				log.Println(sym.Symbol + ".P")
				data, err = client.GetHistoricalData(
					sym.Symbol+".P",
					sym.Exchange,
					interval,
					bars,
					nil,
					false,
				)

				if err == nil && data != nil {
					// Filter data by time range if specified
					if useTimeRange && data != nil {
						filteredData := make([]tradingview.HistoricalData, 0)
						for _, d := range data {
							if (d.DateTime.Equal(fromTime) || d.DateTime.After(fromTime)) &&
								(d.DateTime.Equal(toTime) || d.DateTime.Before(toTime)) {
								filteredData = append(filteredData, d)
							}
						}
						data = filteredData
					}
				}

				h.logger.Error(fmt.Errorf("attempt %d: failed to get data for %s:%s: %w",
					attemptsCount+1, sym.Symbol, sym.Exchange, err))
				time.Sleep(time.Second * 2)
			}

			resultChan <- result{
				key:  strings.Split(sym.Symbol, ":")[0], // Extract just the ticker part
				data: data,
			}
		}(symbolReq)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make(map[string][]tradingview.HistoricalData)
	for r := range resultChan {
		if r.data != nil {
			results[r.key] = r.data
		}
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data for all symbols",
		})
	}
	// cal price change percentage bar 1 - bar 24, priceResults should in format like {"BTC": 0.01, "ETH": 0.02}
	priceResults := make(map[string]interface{})
	for symbol, result := range results {
		if len(result) > 0 && len(result) >= 24 {
			priceResults[symbol] = (result[len(result)-1].Close - result[len(result)-24].Close) / result[len(result)-24].Close
		}
	}

	// Add safety checks for the symbols loop as well
	for _, symbol := range req.Symbols {
		// Extract just the ticker part from the symbol
		tickerName := strings.Split(symbol.Symbol, ":")[0]
		if _, ok := priceResults[tickerName]; !ok {
			priceResults[tickerName] = nil
		}
	}

	return c.JSON(fiber.Map{
		"data": results,
	})
}

// Get24HoursPriceChange handles requests for 24-hour price changes for multiple symbols
func (h *TradingViewHandler) Get24HoursPriceChangeTradingView(c *fiber.Ctx) error {
	var req model.Get24HoursPriceBySymbolsReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	payload := requestPayload{
		Columns:             []string{"base_currency", "exchange", "currency", "24h_close_change|5"},
		Markets:             []string{"coin"},
		IgnoreUnknownFields: false,
		Range:               []int{0, 150},
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
			"error": "Failed to encode request",
		})
	}

	response, err := http.Post(
		"https://scanner.tradingview.com/coin/scan?label-product=screener-coin",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var scannerResponse model.ScannerResponse

	if err := json.Unmarshal(body, &scannerResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse scanner response",
		})
	}

	priceChangeResults := h.priceService.GetPriceChange(req.Symbols, scannerResponse.Data)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"results": priceChangeResults,
	})
}
