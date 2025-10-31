package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

// DexService defines operations available for DEX wallets with exchange awareness.
type DexService interface {
	Connect(ctx context.Context, uid string, req model.DexConnectRequest) (string, error)
	ListWallets(ctx context.Context, uid string, exchange string) ([]model.WalletInfo, error)
	ActiveWallet(ctx context.Context, uid, walletID string) error
	DeactiveWallet(ctx context.Context, uid, walletID string) error
	UpdatePositionSize(ctx context.Context, uid, walletID string, positionSize float64) error
	UpdateLeverage(ctx context.Context, uid, walletID string, leverage int, exchange string) error
	UpdateHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error
	UpdateSL(ctx context.Context, uid, walletID string, slPercentage float64, exchange string) error
	SubscribeAuthor(ctx context.Context, author string, walletID string) (string, error)
	UnsubscribeAuthor(ctx context.Context, author string, walletID string) error
	GetWalletTotalValue(ctx context.Context, uid, walletID, exchange string) (model.DexWalletTotalValue, error)
	ValidateDexCredentials(ctx context.Context, exchange, apiKey, privateKey, tradingAccountID string) (bool, error)
	UpdateAPICredentials(ctx context.Context, uid, walletID, apiKey, privateKey, tradingAccountID, exchange string) error
}

// DexRepo captures data access patterns required by the DEX service layer.
type DexRepo interface {
	InsertDexWallet(ctx context.Context, uid string, record model.DexWalletRecord) (string, error)
	ListDexWallets(ctx context.Context, uid string, exchange string) ([]model.DexWalletRecord, error)
	ActiveDexWallet(ctx context.Context, uid, walletID string) error
	DeactiveDexWallet(ctx context.Context, uid, walletID string) error
	GetSubscribeAuthor(ctx context.Context, walletID string) ([]model.SubscribeAuthor, error)
	WalletExistsByAddress(ctx context.Context, uid string, walletAddress string, exchange string) (bool, error)
	UpdateDexWalletPositionSize(ctx context.Context, uid, walletID string, positionSize float64) error
	UpdateDexWalletLeverage(ctx context.Context, uid, walletID string, leverage int, exchange string) error
	UpdateDexWalletHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error
	UpdateDexWalletSL(ctx context.Context, uid, walletID string, sl float64, exchange string) error
	SubscribeAuthor(ctx context.Context, author string, walletID string) (string, error)
	UnsubscribeAuthor(ctx context.Context, author string, walletID string) error
	GetDexWalletTotalValue(ctx context.Context, uid, walletID string, exchange string) (model.DexWalletTotalValue, error)
	ValidateDexCredentials(ctx context.Context, exchange, apiKey, privateKey, tradingAccountID string) (bool, error)
	UpdateDexWalletAPICredentials(ctx context.Context, uid, walletID, apiKey, privateKey, tradingAccountID, exchange string) error
}
