package model

import "time"

type Wallet struct {
	WalletName    string  `json:"wallet_name" example:"Wallet Name"`
	WalletAddress string  `json:"wallet_address" example:"0x1234567890123456789012345678901234567890"`
	PrivyWalletID string  `json:"privy_wallet_id" example:"309339ae6e80xxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
	Priority      int     `json:"priority" example:"2"`
	Fee           float32 `json:"fee" example:"0.1"`
	WalletType    string  `json:"wallet_type" example:"copytrade or genesis"`
}

type UserWallet struct {
	WalletName             string     `json:"wallet_name" example:"Wallet Name"`
	WalletID               string     `json:"wallet_id" example:"1234567890"`
	WalletType             string     `json:"wallet_type" example:"copytrade or genesis"`
	Balance                float64    `json:"balance" example:"1000"`
	WalletAddress          string     `json:"wallet_address" example:"0x1234567890123456789012345678901234567890"`
	PrivyWalletID          string     `json:"privy_wallet_id" example:"309339ae6e80xxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
	Priority               int        `json:"priority" example:"2"`
	ExecutionFee           float32    `json:"execution_fee" example:"0.1"`
	TpPercentage           float64    `json:"tp_percentage" example:"10"`
	SlPercentage           float64    `json:"sl_percentage" example:"10"`
	HoldingHourPeriod      int        `json:"holding_hour_period" example:"24"`
	PositionSizePercentage float64    `json:"position_size_percentage" example:"0.2"`
	Leverage               int        `json:"leverage" example:"1"`
	HyperliquidBasecode    bool       `json:"hyperliquid_basecode" example:"true"`
	CreatedAt              *time.Time `json:"created_at" example:"2025-07-08T08:20:49+07:00"`
	UpdatedAt              *time.Time `json:"updated_at" example:"2025-07-08T08:20:49+07:00"`
	DeletedAt              *time.Time `json:"deleted_at" example:"2025-07-08T08:20:49+07:00"`
}

type SubscribeAuthor struct {
	ID             string `json:"id"`
	AuthorUsername string `json:"author_username"`
}

type WalletInfo struct {
	WalletName             string            `json:"wallet_name"`
	WalletID               string            `json:"wallet_id"`
	WalletAddress          string            `json:"wallet_address"`
	WalletType             string            `json:"wallet_type"`
	Balance                float64           `json:"balance"`
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

type PrivyWallet struct {
	ID                string   `json:"id"`
	Address           string   `json:"address"`
	ChainType         string   `json:"chain_type"`
	PolicyIDs         []string `json:"policy_ids"`
	AdditionalSigners []string `json:"additional_signers"`
	OwnerID           string   `json:"owner_id"`
	CreatedAt         int64    `json:"created_at"`  // epoch millis
	ExportedAt        *int64   `json:"exported_at"` // null or epoch millis
	ImportedAt        *int64   `json:"imported_at"` // null or epoch millis
}
