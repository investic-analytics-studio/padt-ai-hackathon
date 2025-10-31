package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type DexHandler struct {
	service port.DexService
}

func NewDexHandler(service port.DexService) *DexHandler {
	return &DexHandler{service: service}
}

type dexWalletRequest struct {
	WalletID string `json:"wallet_id"`
}

// Connect godoc
// @Summary      Connect DEX wallet
// @Description  Link a DEX account to the current authenticated user
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload  body      model.DexConnectRequest true "DEX connect payload"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /dex/add-wallet [post]
// @Security     BearerAuth
func (h *DexHandler) Connect(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		logger.Warn("dex connect: missing uid in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.DexConnectRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Errorf("dex connect: body parse error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	walletID, err := h.service.Connect(c.UserContext(), uid, req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrDexWalletExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, model.ErrDexMissingFields),
			errors.Is(err, model.ErrDexInvalidExchange),
			errors.Is(err, model.ErrDexInvalidKey),
			errors.Is(err, model.ErrDexInvalidPosition),
			errors.Is(err, model.ErrDexInvalidLeverage),
			errors.Is(err, model.ErrDexInvalidSL):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Errorf("dex connect: failed to connect wallet for uid=%s: %v", uid, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect DEX wallet"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"wallet_id": walletID})
}

// ListWallets godoc
// @Summary      List DEX wallets
// @Description  Retrieve DEX wallets for the current authenticated user filtered by exchange
// @Tags         copytrade/dex
// @Produce      json
// @Param        exchange query    string true "Exchange identifier (e.g. dydx)"
// @Success      200      {array}  model.WalletInfo
// @Failure      400      {object} map[string]string
// @Failure      401      {object} map[string]string
// @Failure      500      {object} map[string]string
// @Router       /dex/wallet-info [get]
// @Security     BearerAuth
func (h *DexHandler) ListWallets(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	exchange := strings.TrimSpace(c.Query("exchange"))
	wallets, err := h.service.ListWallets(c.UserContext(), uid, exchange)
	if err != nil {
		if errors.Is(err, model.ErrDexInvalidExchange) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("dex list wallets: uid=%s exchange=%s err=%v", uid, exchange, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve wallets"})
	}

	return c.Status(fiber.StatusOK).JSON(wallets)
}

// GetWalletTotalValue godoc
// @Summary      Get wallet total value
// @Description  Retrieve the total value of a DEX wallet for the current authenticated user
// @Tags         copytrade/dex
// @Produce      json
// @Param        wallet_id query    string true "DEX wallet ID"
// @Param        exchange  query    string true "Exchange identifier (e.g. dydx)"
// @Success      200       {object}  model.DexWalletTotalValue
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /dex/wallet-total-value [get]
// @Security     BearerAuth
func (h *DexHandler) GetWalletTotalValue(c *fiber.Ctx) error {
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
		if errors.Is(err, model.ErrDexInvalidExchange) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("dex wallet total value: uid=%s wallet_id=%s exchange=%s err=%v", uid, walletID, exchange, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve wallet total value"})
	}

	return c.Status(fiber.StatusOK).JSON(totalValue)
}

// ActiveWallet godoc
// @Summary      Activate DEX wallet
// @Description  Set the selected DEX wallet as active for the current user
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload body      dexWalletRequest true "DEX wallet activation payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/active-wallet [post]
// @Security     BearerAuth
func (h *DexHandler) ActiveWallet(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req dexWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.ActiveWallet(c.UserContext(), uid, req.WalletID); err != nil {
		logger.Errorf("dex active wallet: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to activate wallet"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "wallet activated"})
}

// DeactiveWallet godoc
// @Summary      Deactivate DEX wallet
// @Description  Mark the selected DEX wallet as inactive for the current user
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload body      dexWalletRequest true "DEX wallet deactivation payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/deactive-wallet [post]
// @Security     BearerAuth
func (h *DexHandler) DeactiveWallet(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req dexWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.DeactiveWallet(c.UserContext(), uid, req.WalletID); err != nil {
		logger.Errorf("dex deactive wallet: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deactivate wallet"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "wallet deactivated"})
}

// UpdatePositionSize godoc
// @Summary      Update position size
// @Description  Update the position size percentage of a DEX wallet
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload body      model.DexUpdatePositionSizeRequest true "Position size payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/update-position-size [post]
// @Security     BearerAuth
func (h *DexHandler) UpdatePositionSize(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.DexUpdatePositionSizeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.UpdatePositionSize(c.UserContext(), uid, req.WalletID, req.PositionSizePercentage); err != nil {
		if errors.Is(err, model.ErrDexInvalidPosition) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Errorf("dex update position size: uid=%s wallet_id=%s err=%v", uid, req.WalletID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update position size"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "position size updated"})
}

// UpdateLeverage godoc
// @Summary      Update leverage
// @Description  Update the leverage of a DEX wallet
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload body      model.DexUpdateLeverageRequest true "Leverage payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/update-leverage [post]
// @Security     BearerAuth
func (h *DexHandler) UpdateLeverage(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.DexUpdateLeverageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.WalletID) == "" || strings.TrimSpace(req.Exchange) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id and exchange are required"})
	}

	if err := h.service.UpdateLeverage(c.UserContext(), uid, req.WalletID, req.Leverage, req.Exchange); err != nil {
		switch {
		case errors.Is(err, model.ErrDexInvalidLeverage),
			errors.Is(err, model.ErrDexInvalidExchange):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Errorf("dex update leverage: uid=%s wallet_id=%s exchange=%s err=%v", uid, req.WalletID, req.Exchange, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update leverage"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "leverage updated"})
}

// UpdateAPICredentials godoc
// @Summary      Update DEX API credentials
// @Description  Update the API key, private key, and trading account of a DEX wallet
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload body      model.DexUpdateAPICredentialsRequest true "DEX API credential payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/update-api-credentials [post]
// @Security     BearerAuth
func (h *DexHandler) UpdateAPICredentials(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.DexUpdateAPICredentialsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}
	if strings.TrimSpace(req.Exchange) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "exchange is required"})
	}
	if strings.TrimSpace(req.APIKey) == "" || strings.TrimSpace(req.PrivateKey) == "" || strings.TrimSpace(req.TradingAccountID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "api_key, private_key, and trading_account are required"})
	}

	if err := h.service.UpdateAPICredentials(c.UserContext(), uid, req.WalletID, req.APIKey, req.PrivateKey, req.TradingAccountID, req.Exchange); err != nil {
		switch {
		case errors.Is(err, model.ErrDexMissingFields),
			errors.Is(err, model.ErrDexInvalidExchange),
			errors.Is(err, model.ErrDexInvalidCredentials):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Errorf("dex update api credentials: uid=%s wallet_id=%s exchange=%s err=%v", uid, req.WalletID, req.Exchange, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update API credentials"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "api credentials updated"})
}

// UpdateSL godoc
// @Summary      Update stop loss
// @Description  Update the stop loss percentage of a DEX wallet
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload body      model.DexUpdateSLRequest true "Stop loss payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/update-sl [post]
// @Security     BearerAuth
func (h *DexHandler) UpdateSL(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.DexUpdateSLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.WalletID) == "" || strings.TrimSpace(req.Exchange) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id and exchange are required"})
	}

	if err := h.service.UpdateSL(c.UserContext(), uid, req.WalletID, req.SlPercentage, req.Exchange); err != nil {
		switch {
		case errors.Is(err, model.ErrDexInvalidSL),
			errors.Is(err, model.ErrDexInvalidExchange):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Errorf("dex update sl: uid=%s wallet_id=%s exchange=%s err=%v", uid, req.WalletID, req.Exchange, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update sl percentage"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "sl percentage updated"})
}

// UpdateHoldingPeriod godoc
// @Summary      Update holding period
// @Description  Update the holding period (in hours) configured for a DEX wallet
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        payload body      model.DexUpdateHoldingPeriodRequest true "Holding period payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/update-holding-period [post]
// @Security     BearerAuth
func (h *DexHandler) UpdateHoldingPeriod(c *fiber.Ctx) error {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req model.DexUpdateHoldingPeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id is required"})
	}

	if err := h.service.UpdateHoldingPeriod(c.UserContext(), uid, req.WalletID, req.HoldingPeriod); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update holding period"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "holding period updated"})
}

// SubscribeAuthor godoc
// @Summary      Subscribe author
// @Description  Subscribe the specified author for copy-trading on a DEX wallet
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        author body      AuthorRequest true "Author payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/subscribe-author [post]
// @Security     BearerAuth
func (h *DexHandler) SubscribeAuthor(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	if uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req AuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.Author) == "" || strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "author and wallet id are required"})
	}

	subscribeID, err := h.service.SubscribeAuthor(c.UserContext(), req.Author, req.WalletID)
	if err != nil {
		logger.Errorf("dex subscribe author: uid=%s wallet_id=%s author=%s err=%v", uid, req.WalletID, req.Author, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to subscribe author"})
	}
	if subscribeID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get inserted subscribe id"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"subscribe_id": subscribeID})
}

// UnsubscribeAuthor godoc
// @Summary      Unsubscribe author
// @Description  Remove the author subscription from a DEX wallet
// @Tags         copytrade/dex
// @Accept       json
// @Produce      json
// @Param        author body      UnAuthorRequest true "Author payload"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /dex/unsubscribe-author [post]
// @Security     BearerAuth
func (h *DexHandler) UnsubscribeAuthor(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	if uid == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req UnAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if strings.TrimSpace(req.Author) == "" || strings.TrimSpace(req.WalletID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "author and wallet id are required"})
	}

	if err := h.service.UnsubscribeAuthor(c.UserContext(), req.Author, req.WalletID); err != nil {
		logger.Errorf("dex unsubscribe author: uid=%s wallet_id=%s author=%s err=%v", uid, req.WalletID, req.Author, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unsubscribe author"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "author has been unsubscribed"})
}
