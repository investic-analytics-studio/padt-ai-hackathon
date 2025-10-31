package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/util"
)

type marketSummaryHandler struct {
	serv domain.MarketSummaryService
}

func NewMarketSummaryHandler(serv domain.MarketSummaryService) *marketSummaryHandler {
	return &marketSummaryHandler{
		serv: serv,
	}
}

func (h *marketSummaryHandler) GetDailyMarketOverview(c *fiber.Ctx) error {
	result, err := h.serv.GetDailyMarketOverview()
	if err != nil {
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetTopGainer(c *fiber.Ctx) error {
	result, err := h.serv.GetTopGainer()
	if err != nil {
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetTopLoser(c *fiber.Ctx) error {
	result, err := h.serv.GetTopLoser()
	if err != nil {
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetAdvancerDeclinerDistributionHist(c *fiber.Ctx) error {
	result, err := h.serv.GetAdvancerDeclinerDistributionHist()
	if err != nil {
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetAdvancerDeclinerDistributionBar(c *fiber.Ctx) error {
	result, err := h.serv.GetAdvancerDeclinerDistributionBar()
	if err != nil {
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetTopTurnoverFloat(c *fiber.Ctx) error {
	result, err := h.serv.GetTopTurnoverFloat()
	if err != nil {
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCMEGoldOI(c *fiber.Ctx) error {
	result, err := h.serv.GetCMEGoldOI()
	if err != nil {
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetStockAlertPlots(c *fiber.Ctx) error {
	result, err := h.serv.GetStockAlertPlots()
	if err != nil {
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetStockAlertStatsDates(c *fiber.Ctx) error {
	result, err := h.serv.GetStockAlertStatsDates()
	if err != nil {
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetStockAlertStats(c *fiber.Ctx) error {
	result, err := h.serv.GetStockAlertStats()
	if err != nil {
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetStockAlertDetections(c *fiber.Ctx) error {
	result, err := h.serv.GetStockAlertDetections()
	if err != nil {
		log.Printf("Error in GetStockAlertDetections handler: %v", err)
		return util.ResponseError(c, err)
	}
	log.Printf("GetStockAlertDetections handler result count: %d", len(result))
	if len(result) == 0 {
		return util.ResponseOK(c, []model.StockAlertDetection{})
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCryptoAlertPlots(c *fiber.Ctx) error {
	result, err := h.serv.GetCryptoAlertPlots()
	if err != nil {
		log.Printf("Error in GetCryptoAlertPlots handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCryptoAlertStatsDates(c *fiber.Ctx) error {
	result, err := h.serv.GetCryptoAlertStatsDates()
	if err != nil {
		log.Printf("Error in GetCryptoAlertStatsDates handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCryptoAlertStats(c *fiber.Ctx) error {
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 50)

	result, err := h.serv.GetCryptoAlertStats(page, limit)
	if err != nil {
		log.Printf("Error in GetCryptoAlertStats handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCryptoAlertDetections(c *fiber.Ctx) error {
	result, err := h.serv.GetCryptoAlertDetections()
	if err != nil {
		log.Printf("Error in GetCryptoAlertDetections handler: %v", err)
		return util.ResponseError(c, err)
	}
	log.Printf("GetCryptoAlertDetections handler result count: %d", len(result))
	if len(result) == 0 {
		return util.ResponseOK(c, []model.CryptoAlertDetection{})
	}
	return util.ResponseOK(c, result)
}

// Add these new methods to the existing file

func (h *marketSummaryHandler) GetCryptoAlertPlots1D(c *fiber.Ctx) error {
	result, err := h.serv.GetCryptoAlertPlots1D()
	if err != nil {
		log.Printf("Error in GetCryptoAlertPlots1D handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCryptoAlertStatsDates1D(c *fiber.Ctx) error {
	result, err := h.serv.GetCryptoAlertStatsDates1D()
	if err != nil {
		log.Printf("Error in GetCryptoAlertStatsDates1D handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCryptoAlertStats1D(c *fiber.Ctx) error {
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 50)

	result, err := h.serv.GetCryptoAlertStats1D(page, limit)
	if err != nil {
		log.Printf("Error in GetCryptoAlertStats1D handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetCryptoAlertDetections1D(c *fiber.Ctx) error {
	result, err := h.serv.GetCryptoAlertDetections1D()
	if err != nil {
		log.Printf("Error in GetCryptoAlertDetections1D handler: %v", err)
		return util.ResponseError(c, err)
	}
	log.Printf("GetCryptoAlertDetections1D handler result count: %d", len(result))
	if len(result) == 0 {
		return util.ResponseOK(c, []model.CryptoAlertDetection1D{})
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetAllCryptoAlertStats(c *fiber.Ctx) error {
	result, err := h.serv.GetAllCryptoAlertStats()
	if err != nil {
		log.Printf("Error in GetAllCryptoAlertStats handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}

func (h *marketSummaryHandler) GetAllCryptoAlertStats1D(c *fiber.Ctx) error {
	result, err := h.serv.GetAllCryptoAlertStats1D()
	if err != nil {
		log.Printf("Error in GetAllCryptoAlertStats1D handler: %v", err)
		return util.ResponseError(c, err)
	}
	return util.ResponseOK(c, result)
}
