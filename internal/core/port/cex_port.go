package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

// CexService exposes the business operations available for CEX wallets.
type CexService interface {
	Connect(ctx context.Context, uid string, req model.CexConnectRequest) (string, error)
	ListWallets(ctx context.Context, uid string, exchange string) ([]model.CexWalletInfo, error)
	ActiveWallet(ctx context.Context, uid, walletID string) error
	DeactiveWallet(ctx context.Context, uid, walletID string) error
	UpdateHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error
	UpdatePositionSize(ctx context.Context, uid, walletID string, positionSize float64) error
	UpdateLeverage(ctx context.Context, uid, walletID string, leverage int) error
	UpdateSL(ctx context.Context, uid, walletID string, sl float64, exchange string) error
	UpdateAPICredentials(ctx context.Context, uid, walletID string, apiKey string, apiSecret string, exchange string) error
	SubscribeAuthor(ctx context.Context, author, walletID string) (string, error)
	UnsubscribeAuthor(ctx context.Context, author, walletID string) error
	GetWalletTotalValue(ctx context.Context, uid, walletID, exchange string) (model.CexWalletTotalValue, error)
}

// CexRepo describes the storage interactions required by the service layer.
type CexRepo interface {
	WalletExistsByAddress(ctx context.Context, uid string, walletAddress string, exchange string) (bool, error)
	InsertCexWallet(ctx context.Context, uid string, record model.CexWalletRecord) (string, error)
	ListCexWallets(ctx context.Context, uid string, exchange string) ([]model.CexWalletRecord, error)
	ActiveCexWallet(ctx context.Context, uid, walletID string) error
	DeactiveCexWallet(ctx context.Context, uid, walletID string) error
	UpdateCexWalletHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error
	UpdateCexWalletPositionSize(ctx context.Context, uid, walletID string, positionSize float64) error
	UpdateCexWalletLeverage(ctx context.Context, uid, walletID string, leverage int) error
	UpdateCexWalletSL(ctx context.Context, uid, walletID string, sl float64, exchange string) error
	UpdateCexWalletAPICredentials(ctx context.Context, uid, walletID string, apiKey string, apiSecret string, exchange string) error
	ValidateCexCredentials(ctx context.Context, exchange string, apiKey string, apiSecret string) (bool, error)
	GetSubscribeAuthor(ctx context.Context, walletID string) ([]model.SubscribeAuthor, error)
	GetCexWalletTotalValue(ctx context.Context, uid, walletID string, exchange string) (model.CexWalletTotalValue, error)
	SubscribeAuthor(ctx context.Context, author, walletID string) (string, error)
	UnsubscribeAuthor(ctx context.Context, author, walletID string) error
}
