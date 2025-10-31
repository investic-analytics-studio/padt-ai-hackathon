package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
	"github.com/quantsmithapp/datastation-backend/pkg/util"
)

type OHLCVDataHandler struct {
	serv   domain.OHLCVDataService
	logger logger.Logger
}

func NewOHLCVDataHandler(serv domain.OHLCVDataService) *OHLCVDataHandler {
	return &OHLCVDataHandler{
		serv:   serv,
		logger: logger.NewLogger(),
	}
}

func (h *OHLCVDataHandler) GetCryptoOHLCV(c *fiber.Ctx) error {
	req, err := parseOHLCVRequest(c)
	if err != nil {
		h.logger.Error(err)
		return util.ResponseError(c, err)
	}

	result, err := h.serv.GetCryptoOHLCV(req)
	if err != nil {
		h.logger.Error(err)
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func (h *OHLCVDataHandler) GetForexOHLCV(c *fiber.Ctx) error {
	req, err := parseOHLCVRequest(c)
	if err != nil {
		h.logger.Error(err)
		return util.ResponseError(c, err)
	}

	result, err := h.serv.GetForexOHLCV(req)
	if err != nil {
		h.logger.Error(err)
		return util.ResponseError(c, err)
	}

	return util.ResponseOK(c, result)
}

func parseOHLCVRequest(c *fiber.Ctx) (model.OHLCVRequest, error) {
	req := model.OHLCVRequest{
		Ticker:    c.Query("ticker"),
		TimeFrame: c.Query("tf"),
		AllPair:   c.QueryBool("all_pair", false),
	}

	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		return req, fmt.Errorf("invalid start_date format: %v", err)
	}
	req.StartDate = startDate

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			// Try parsing with the date-only format if RFC3339 fails
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				return req, fmt.Errorf("invalid end_date format: %v", err)
			}
		}
		req.EndDate = &endDate
	} else {
		now := time.Now()
		req.EndDate = &now
	}

	return req, nil
}
