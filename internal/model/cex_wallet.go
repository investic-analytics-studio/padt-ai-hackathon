package model

import (
	"time"
)

type CexWalletTotalValue struct {
	TotalValue float64 `json:"total_value"`
}
type CexWalletRecord struct {
	ID                     string     `db:"id"`
	CryptoUserID           string     `db:"crypto_user_id"`
	WalletAddress          string     `db:"wallet_address"`
	APIKey                 string     `db:"api_key"`
	APISecret              string     `db:"api_secret"`
	Priority               int        `db:"priority"`
	ExecutionFee           float64    `db:"execution_fee"`
	Exchange               string     `db:"exchange"`
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

type CexWalletInfo struct {
	WalletName             string            `json:"wallet_name"`
	WalletID               string            `json:"wallet_id"`
	WalletType             string            `json:"wallet_type"`
	PrivyWalletID          string            `json:"privy_wallet_id"`
	Priority               int               `json:"priority"`
	Authors                []SubscribeAuthor `json:"authors"`
	TpPercentage           float64           `json:"tp_percentage"`
	SlPercentage           float64           `json:"sl_percentage"`
	HoldingHourPeriod      int               `json:"holding_hour_period"`
	PositionSizePercentage float64           `json:"position_size_percentage"`
	Leverage               int               `json:"leverage"`
	HyperliquidBasecode    bool              `json:"hyperliquid_basecode"`
	CreatedAt              *time.Time        `json:"created_at"`
	UpdatedAt              *time.Time        `json:"updated_at"`
	DeletedAt              *time.Time        `json:"deleted_at"`
}
