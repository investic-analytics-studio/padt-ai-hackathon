package model

import "errors"

var (
	ErrCexWalletExists       = errors.New("cex wallet already connected")
	ErrCexMissingFields      = errors.New("wallet_address, api_key, api_secret, and exchange are required")
	ErrCexInvalidPosition    = errors.New("position size percentage must be greater than 0 and at most 1")
	ErrCexInvalidLeverage    = errors.New("leverage must be between 1 and 100")
	ErrCexInvalidSL          = errors.New("sl percentage must be between 0 and 100")
	ErrCexInvalidExchange    = errors.New("exchange must be provided")
	ErrCexMissingCredentials = errors.New("api_key and api_secret are required")
	ErrCexInvalidCredentials = errors.New("invalid api credentials")
)

// CexConnectRequest captures the request payload to connect a CEX wallet.
type CexConnectRequest struct {
	// WalletAddress          string   `json:"wallet_address" example:"copytrade-primary"`
	APIKey     string  `json:"api_key" example:"cex_api_key"`
	APISecret  string  `json:"api_secret" example:"cex_api_secret"`
	Exchange   string  `json:"exchange" example:"binance-th"`
	WalletName *string `json:"wallet_name,omitempty" example:"CEX Main"`
	// Priority               *int     `json:"priority,omitempty" example:"2"`
	// ExecutionFee           *float64 `json:"execution_fee,omitempty" example:"0.1"`
	// PositionSizePercentage *float64 `json:"position_size_percentage,omitempty" example:"0.15"`
	// Leverage               *int     `json:"leverage,omitempty" example:"5"`
	// SlPercentage           *float64 `json:"sl_percentage,omitempty" example:"30"`
}

// CexUpdateAPIKeyRequest represents a request to update API credentials.
type CexUpdateAPIKeyRequest struct {
	WalletID  string `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	APIKey    string `json:"api_key" example:"cex_api_key"`
	APISecret string `json:"api_secret" example:"cex_api_secret"`
	Exchange  string `json:"exchange" example:"binance-th"`
}

// CexUpdatePositionSizeRequest represents a request to update the position size percentage.
type CexUpdatePositionSizeRequest struct {
	WalletID               string  `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	PositionSizePercentage float64 `json:"position_size_percentage" example:"0.2"`
}

// CexUpdateLeverageRequest represents a request to update wallet leverage.
type CexUpdateLeverageRequest struct {
	WalletID string `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	Leverage int    `json:"leverage" example:"10"`
}

// CexUpdateHoldingPeriodRequest represents a request to update wallet leverage.
type CexUpdateHoldingPeriodRequest struct {
	WalletID      string `json:"wallet_id"`
	HoldingPeriod int    `json:"holding_period"`
}

// CexUpdateSLRequest represents a request to update stop-loss percentage.
type CexUpdateSLRequest struct {
	WalletID     string  `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	SlPercentage float64 `json:"sl_percentage" example:"25"`
	Exchange     string  `json:"exchange" example:"binance-th"`
}
