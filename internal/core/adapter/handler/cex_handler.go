package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type CexHandler struct {
	service port.CexService
}

func NewCexHandler(service port.CexService) *CexHandler {
	return &CexHandler{service: service}
}

type cexWalletIDRequest struct {
	WalletID string `json:"wallet_id"`
}

type cexAuthorRequest struct {
	Author   string `json:"author"`
	WalletID string `json:"wallet_id"`
}

// Connect godoc
// @Summary      Connect CEX wallet
// @Description  Link a CEX account to the current authenticated user
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload  body      model.CexConnectRequest true "CEX connect payload"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /cex/add-wallet [post]
// @Security     BearerAuth
func (h *CexHandler) Connect(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		logger.Warn("cex connect: missing uid in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.CexConnectRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Errorf("cex connect: body parse error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	walletID, err := h.service.Connect(c.UserContext(), uid, req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrCexWalletExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, model.ErrCexMissingFields),
			errors.Is(err, model.ErrCexMissingCredentials),
			errors.Is(err, model.ErrCexInvalidExchange),
			errors.Is(err, model.ErrCexInvalidPosition),
			errors.Is(err, model.ErrCexInvalidLeverage),
			errors.Is(err, model.ErrCexInvalidSL),
			errors.Is(err, model.ErrCexInvalidCredentials):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Errorf("cex connect: failed to connect wallet for uid=%s: %v", uid, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect CEX wallet"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"wallet_id": walletID})
}

// ListWallets godoc
// @Summary      List CEX wallets
// @Description  Retrieve CEX wallets for the current authenticated user filtered by exchange
// @Tags         copytrade/cex
// @Produce      json
// @Param        exchange query    string true "Exchange identifier (e.g. bybit)"
// @Success      200      {array}  model.CexWalletInfo
// @Failure      400      {object} map[string]string
// @Failure      401      {object} map[string]string
// @Failure      500      {object} map[string]string
// @Router       /cex/wallet-info [get]
// @Security     BearerAuth
func (h *CexHandler) ListWallets(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	exchange := strings.TrimSpace(c.Query("exchange"))
	wallets, err := h.service.ListWallets(c.UserContext(), uid, exchange)
	if err != nil {
		if errors.Is(err, model.ErrCexInvalidExchange) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("cex list wallets: uid=%s exchange=%s err=%v", uid, exchange, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve wallets"})
	}

	return c.Status(fiber.StatusOK).JSON(wallets)
}

// GetWalletTotalValue godoc
// @Summary      Get wallet total value
// @Description  Retrieve the total value of a CEX wallet for the current authenticated user
// @Tags         copytrade/cex
// @Produce      json
// @Param        wallet_id query    string true "CEX wallet ID"
// @Param        exchange  query    string true "Exchange identifier (e.g. binance-th)"
// @Success      200       {object}  model.CexWalletTotalValue
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /cex/wallet-total-value [get]
// @Security     BearerAuth
func (h *CexHandler) GetWalletTotalValue(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)

	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	walletID := strings.TrimSpace(c.Query("wallet_id"))
	if walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	exchange := strings.TrimSpace(c.Query("exchange"))
	totalValue, err := h.service.GetWalletTotalValue(c.UserContext(), uid, walletID, exchange)
	if err != nil {
		if errors.Is(err, model.ErrCexInvalidExchange) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("cex wallet total value: uid=%s wallet_id=%s exchange=%s err=%v", uid, walletID, exchange, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve wallet total value"})
	}

	return c.Status(fiber.StatusOK).JSON(totalValue)
}

// SubscribeAuthor godoc
// @Summary      Subscribe author
// @Description  Subscribe an author to a CEX wallet
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      cexAuthorRequest true "Subscribe author payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/subscribe-author [post]
// @Security     BearerAuth
func (h *CexHandler) SubscribeAuthor(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req cexAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	author := strings.TrimSpace(req.Author)
	walletID := strings.TrimSpace(req.WalletID)
	if author == "" || walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "author and wallet_id are required"})
	}

	subscribeID, err := h.service.SubscribeAuthor(c.UserContext(), author, walletID)
	if err != nil {
		logger.Errorf("cex subscribe author: uid=%s wallet_id=%s author=%s err=%v", uid, walletID, author, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to subscribe author"})
	}
	if strings.TrimSpace(subscribeID) == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get subscribe id"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"subscribe_id": subscribeID})
}

// UnsubscribeAuthor godoc
// @Summary      Unsubscribe author
// @Description  Unsubscribe an author from a CEX wallet
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      cexAuthorRequest true "Unsubscribe author payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/unsubscribe-author [post]
// @Security     BearerAuth
func (h *CexHandler) UnsubscribeAuthor(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req cexAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	author := strings.TrimSpace(req.Author)
	walletID := strings.TrimSpace(req.WalletID)
	if author == "" || walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "author and wallet_id are required"})
	}

	if err := h.service.UnsubscribeAuthor(c.UserContext(), author, walletID); err != nil {
		logger.Errorf("cex unsubscribe author: uid=%s wallet_id=%s author=%s err=%v", uid, walletID, author, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unsubscribe author"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "author has been unsubscribed"})
}

// ActiveWallet godoc
// @Summary      Activate CEX wallet
// @Description  Activate the selected CEX wallet for the current user
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      cexWalletIDRequest true "CEX wallet id"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/active-wallet [post]
// @Security     BearerAuth
func (h *CexHandler) ActiveWallet(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req cexWalletIDRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.WalletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.ActiveWallet(c.UserContext(), uid, req.WalletID); err != nil {
		logger.Errorf("cex active wallet: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to activate wallet"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "wallet activated"})
}

// DeactiveWallet godoc
// @Summary      Deactivate CEX wallet
// @Description  Deactivate the selected CEX wallet for the current user
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      cexWalletIDRequest true "CEX wallet id"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/deactive-wallet [post]
// @Security     BearerAuth
func (h *CexHandler) DeactiveWallet(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req cexWalletIDRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.WalletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.DeactiveWallet(c.UserContext(), uid, req.WalletID); err != nil {
		logger.Errorf("cex deactive wallet: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deactivate wallet"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "wallet deactivated"})
}

// UpdatePositionSize godoc
// @Summary      Update position size percentage
// @Description  Update the position size percentage of a CEX wallet
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      model.CexUpdatePositionSizeRequest true "Position size payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/update-position-size [post]
// @Security     BearerAuth
func (h *CexHandler) UpdatePositionSize(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.CexUpdatePositionSizeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.WalletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.UpdatePositionSize(c.UserContext(), uid, req.WalletID, req.PositionSizePercentage); err != nil {
		if errors.Is(err, model.ErrCexInvalidPosition) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("cex update position size: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update position size"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "position size updated"})
}

// UpdateHoldingPeriod godoc
// @Summary      Update holding period
// @Description  Update the holding period (in hours) configured for a CEX wallet
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      model.CexUpdateHoldingPeriodRequest true "Holding period payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/update-holding-period [post]
// @Security     BearerAuth
func (h *CexHandler) UpdateHoldingPeriod(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var req model.CexUpdateHoldingPeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.WalletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.UpdateHoldingPeriod(c.UserContext(), uid, req.WalletID, req.HoldingPeriod); err != nil {
		if errors.Is(err, model.ErrCexInvalidLeverage) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("cex holding period: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update holding period"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "holding period updated"})
}

func (h *CexHandler) UpdateLeverage(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.CexUpdateLeverageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.WalletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.UpdateLeverage(c.UserContext(), uid, req.WalletID, req.Leverage); err != nil {
		if errors.Is(err, model.ErrCexInvalidLeverage) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("cex update leverage: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update leverage"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "leverage updated"})
}

// UpdateSL godoc
// @Summary      Update stop loss
// @Description  Update the stop loss percentage of a CEX wallet
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      model.CexUpdateSLRequest true "Stop loss payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/update-sl [post]
// @Security     BearerAuth
func (h *CexHandler) UpdateSL(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.CexUpdateSLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.WalletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.UpdateSL(c.UserContext(), uid, req.WalletID, req.SlPercentage, req.Exchange); err != nil {
		if errors.Is(err, model.ErrCexInvalidSL) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("cex update sl: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update sl percentage"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "sl percentage updated"})
}

// UpdateAPIKey godoc
// @Summary      Update API credentials
// @Description  Update the API key and secret of a CEX wallet
// @Tags         copytrade/cex
// @Accept       json
// @Produce      json
// @Param        payload body      model.CexUpdateAPIKeyRequest true "API credential payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cex/update-api-key [post]
// @Security     BearerAuth
func (h *CexHandler) UpdateAPIKey(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.CexUpdateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}
	if strings.TrimSpace(req.Exchange) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "exchange is required"})
	}

	if err := h.service.UpdateAPICredentials(c.UserContext(), uid, req.WalletID, req.APIKey, req.APISecret, req.Exchange); err != nil {
		switch {
		case errors.Is(err, model.ErrCexMissingCredentials),
			errors.Is(err, model.ErrCexInvalidExchange),
			errors.Is(err, model.ErrCexInvalidCredentials):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Errorf("cex update api key: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update API credentials"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "api credentials updated"})
}
