package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	_ "github.com/quantsmithapp/datastation-backend/internal/model" // Required for Swagger to find model definitions
)

type PerformanceHandler struct {
	service port.PerformanceService
}

func NewPerformanceHandler(service port.PerformanceService) *PerformanceHandler {
	return &PerformanceHandler{service: service}
}

// GetAuthorNavHandle godoc
// @Summary      Get author navigation performance data
// @Description  Get author performance navigation data for all authors in specified period
// @Tags         Performance
// @Accept       json
// @Produce      json
// @Param        period query string true "Period" enums(7D,1M,ALL) example("7D")
// @Success      200 {array} model.AuthorNav
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /performance/get-author-nav [get]
// @Security     BearerAuth
func (h *PerformanceHandler) GetAuthorNavHandle(c *fiber.Ctx) error {
	period := c.Query("period", "")
	if period == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Period is required",
		})
	}
	authorNav, err := h.service.GetAuthorNav(c.UserContext(), period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get NAV",
		})
	}
	return c.JSON(authorNav)
}

// GetMultiholdingPortNavHandle godoc
// @Summary      Get multiholding portfolio NAV performance data with different holding period
// @Description  Get multiholding portfolio NAV performance data with different holding period
// @Tags         Performance
// @Accept       json
// @Produce      json
// @Param        period query string false "NAV Period" enums(7D,1M,ALL) default(7D) example("7D")
// @Param        holding_period query string false "Holding Period (Hours)" enums(24,48,72,96,120,144,168) default(24) example("24")
// @Success      200 {array} model.AuthorNav
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /performance/get-multiholding-port-nav [get]
// @Security     BearerAuth
func (h *PerformanceHandler) GetMultiholdingPortNavHandle(c *fiber.Ctx) error {
	period := c.Query("period", "7D")
	holdingPeriod := c.Query("holding_period", "24")
	multiholdingPortNav, err := h.service.GetMultiholdingPortNav(c.UserContext(), period, holdingPeriod)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get NAV",
		})
	}
	return c.JSON(multiholdingPortNav)
}

// GetAuthorDetailHandle godoc
// @Summary      Get author performance details with merged timeline and sentiment analysis
// @Description  Get author detail with merged timeline of tweets and signals, including pagination metadata and sentiment analysis data (bearishTokens and bullishTokens) categorized by the specified period
// @Tags         Performance
// @Accept       json
// @Produce      json
// @Param        author_username query string true "Author Username" example("0xkyle__")
// @Param        period query string true "Period" enums(7D,1M,ALL) example("7D")
// @Param        start query int false "Start offset for pagination" default(0)
// @Param        limit query int false "Limit for pagination (max 100)" default(20)
// @Success      200 {object} model.AuthorDetail
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /performance/get-author-detail [get]
// @Security     BearerAuth
func (h *PerformanceHandler) GetAuthorDetailHandle(c *fiber.Ctx) error {
	authorUsername := c.Query("author_username", "")
	if authorUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Author Username is required",
		})
	}

	period := c.Query("period", "")
	if period == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Period is required",
		})
	}

	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 20)

	authorDetail, err := h.service.GetAuthorDetail(c.UserContext(), authorUsername, period, start, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get author detail",
		})
	}
	return c.JSON(authorDetail)
}

// GetAuthorMultiholdingDetailHandle godoc
// @Summary      Get author performance details with multiholding portfolio NAV data
// @Description  Get author detail with merged timeline of tweets and signals, including pagination metadata and sentiment analysis data (bearishTokens and bullishTokens) categorized by the specified period, using multiholding portfolio NAV data based on the holding period
// @Tags         Performance
// @Accept       json
// @Produce      json
// @Param        author_username query string true "Author Username" example("0xkyle__")
// @Param        period query string true "Period" enums(7D,1M,ALL) example("7D")
// @Param        holding_period query string true "Holding Period (Hours)" enums(24,48,72,96,120,144,168) example("24")
// @Param        start query int false "Start offset for pagination" default(0)
// @Param        limit query int false "Limit for pagination (max 100)" default(20)
// @Success      200 {object} model.AuthorDetail
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /performance/get-author-multiholding-detail [get]
// @Security     BearerAuth
func (h *PerformanceHandler) GetAuthorMultiholdingDetailHandle(c *fiber.Ctx) error {
	authorUsername := c.Query("author_username", "")
	if authorUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Author Username is required",
		})
	}

	period := c.Query("period", "")
	if period == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Period is required",
		})
	}

	holdingPeriod := c.Query("holding_period", "")
	if holdingPeriod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Holding Period is required",
		})
	}

	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 20)

	authorDetail, err := h.service.GetAuthorMultiholdingDetail(c.UserContext(), authorUsername, period, holdingPeriod, start, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get author multiholding detail",
		})
	}
	return c.JSON(authorDetail)
}
