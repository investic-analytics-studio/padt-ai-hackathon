package model

import (
	"errors"
	"time"
)

var (
	ErrDexWalletExists       = errors.New("dex wallet already connected")
	ErrDexInvalidKey         = errors.New("private key must be prefixed with 0x and contain 64 hex characters")
	ErrDexInvalidPosition    = errors.New("position size percentage must be greater than 0 and at most 1")
	ErrDexInvalidLeverage    = errors.New("leverage must be between 1 and 10")
	ErrDexInvalidSL          = errors.New("sl percentage must be between 0 and 100")
	ErrDexInvalidCredentials = errors.New("invalid dex credentials")
	ErrDexInvalidExchange    = errors.New("exchange must be provided")
	ErrDexMissingFields      = errors.New("api_key, private_key, trading_account_id, and exchange are required")
)

type DexConnectRequest struct {
	APIKey           string  `json:"api_key" example:"33rADOAxPstLKpzcJmDPgASDLjs"`
	PrivateKey       string  `json:"private_key" example:"0xa5e9ed193e183f0bdd72712a89df949e8a8de3b815caa570c44de8697a93b6fs"`
	TradingAccountID string  `json:"trading_account_id" example:"516866395334867"`
	Exchange         string  `json:"exchange" example:"dydx"`
	WalletAddress    *string `json:"wallet_address,omitempty" example:"0xa5e9ed193e183f0bdd72712a89df949e8a8de3b"`
	WalletName       *string `json:"wallet_name,omitempty" example:"DEX Main"`
}

type DexUpdatePositionSizeRequest struct {
	WalletID               string  `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	PositionSizePercentage float64 `json:"position_size_percentage" example:"0.2"`
}

// DexUpdateHoldingPeriodRequest represents a request to update the holding period.
type DexUpdateHoldingPeriodRequest struct {
	WalletID      string `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	HoldingPeriod int    `json:"holding_period" example:"48"`
}

type DexUpdateLeverageRequest struct {
	WalletID string `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	Exchange string `json:"exchange" example:"dydx"`
	Leverage int    `json:"leverage" example:"3"`
}

type DexUpdateSLRequest struct {
	WalletID     string  `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	Exchange     string  `json:"exchange" example:"dydx"`
	SlPercentage float64 `json:"sl_percentage" example:"25"`
}

type DexUpdateAPICredentialsRequest struct {
	WalletID         string `json:"wallet_id" example:"e50b0c09-18c5-4ff0-a832-54473e1b739e"`
	Exchange         string `json:"exchange" example:"dydx"`
	APIKey           string `json:"api_key" example:"dex_api_key"`
	PrivateKey       string `json:"private_key" example:"0xa5e9ed193e183f0bdd72712a89df949e8a8de3b815caa570c44de8697a93b6fs"`
	TradingAccountID string `json:"trading_account" example:"516866395334867"`
}

// DexWalletRecord references the shared persistence representation used for DEX wallets.
type DexWalletRecord struct {
	ID                     string     `db:"id"`
	CryptoUserID           string     `db:"crypto_user_id"`
	WalletAddress          string     `db:"wallet_address"`
	APIKey                 string     `db:"api_key"`
	PrivateKey             string     `db:"private_key"`
	TradingAccount         string     `db:"trading_account"`
	Priority               int        `db:"priority"`
	Exchange               string     `db:"exchange"`
	ExecutionFee           float64    `db:"execution_fee"`
	Leverage               int        `db:"leverage"`
	PositionSizePercentage float64    `db:"position_size_percentage"`
	WalletName             *string    `db:"wallet_name"`
	TpPercentage           *float64   `db:"tp_percentage"`
	SlPercentage           *float64   `db:"sl_percentage"`
	HoldingHourPeriod      *int       `db:"holding_hour_period"`
	DeletedAt              *time.Time `db:"deleted_at"`
	CreatedAt              *time.Time `db:"created_at"`
	UpdatedAt              *time.Time `db:"updated_at"`
}

// DexWalletTotalValue represents the aggregated USD value of a DEX wallet.
type DexWalletTotalValue struct {
	TotalValue float64 `json:"total_value"`
}
