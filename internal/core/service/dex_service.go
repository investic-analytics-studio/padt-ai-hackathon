package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

const (
	defaultDexPriority      = 2
	defaultDexExecutionFee  = 0.1
	defaultDexPositionSize  = 0.1
	defaultDexLeverage      = 1
	defaultDexStopLoss      = 0.0
	defaultDexHoldingPeriod = 48
	minDexPositionSize      = 0.0
	maxDexPositionSize      = 1.0
	minDexLeverage          = 1
	maxDexLeverage          = 10
	minDexStopLoss          = 0.0
	maxDexStopLoss          = 100.0
)

type DexService struct {
	repo port.DexRepo
}

func NewDexService(repo port.DexRepo) *DexService {
	return &DexService{repo: repo}
}

func normalizeDexExchange(exchange string) string {
	trimmed := strings.TrimSpace(strings.ToLower(exchange))
	if trimmed == "" {
		return ""
	}
	var builder strings.Builder
	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func (s *DexService) Connect(ctx context.Context, uid string, req model.DexConnectRequest) (string, error) {
	apiKey := strings.TrimSpace(req.APIKey)
	privateKey := strings.TrimSpace(req.PrivateKey)
	tradingAccount := strings.TrimSpace(req.TradingAccountID)
	exchange := normalizeDexExchange(req.Exchange)

	if apiKey == "" || privateKey == "" || tradingAccount == "" || exchange == "" {
		if exchange == "" {
			return "", model.ErrDexInvalidExchange
		}
		return "", model.ErrDexMissingFields
	}

	var walletAddress string
	if req.WalletAddress != nil && strings.TrimSpace(*req.WalletAddress) != "" {
		walletAddress = strings.ToLower(strings.TrimSpace(*req.WalletAddress))
		if !strings.HasPrefix(walletAddress, "0x") || len(walletAddress) != 42 {
			return "", fmt.Errorf("invalid wallet address format")
		}
	} else {
		var err error
		walletAddress, err = deriveWalletAddress(privateKey)
		if err != nil {
			if errors.Is(err, model.ErrDexInvalidKey) {
				return "", model.ErrDexInvalidKey
			}
			return "", err
		}
	}
	fmt.Println("=========11===============")
	exists, err := s.repo.WalletExistsByAddress(ctx, uid, walletAddress, exchange)
	if err != nil {
		return "", err
	}
	if exists {
		return "", model.ErrDexWalletExists
	}
	fmt.Println("=========22===============")
	valid, err := s.ValidateDexCredentials(ctx, exchange, apiKey, privateKey, tradingAccount)
	if err != nil {
		if errors.Is(err, model.ErrDexInvalidExchange) || errors.Is(err, model.ErrDexMissingFields) {
			return "", err
		}
		logger.Errorf("dex connect: credential validation failed exchange=%s err=%v", exchange, err)
		return "", fmt.Errorf("failed to validate api credentials")
	}
	if !valid {
		return "", model.ErrDexInvalidCredentials
	}
	fmt.Println("=========33===============")
	priority := defaultDexPriority

	executionFee := defaultDexExecutionFee

	positionSize := defaultDexPositionSize
	if positionSize <= minDexPositionSize || positionSize > maxDexPositionSize {
		return "", model.ErrDexInvalidPosition
	}

	leverage := defaultDexLeverage
	if leverage < minDexLeverage || leverage > maxDexLeverage {
		return "", model.ErrDexInvalidLeverage
	}

	slPercentage := defaultDexStopLoss
	if slPercentage < minDexStopLoss || slPercentage > maxDexStopLoss {
		return "", model.ErrDexInvalidSL
	}

	defaultHoldingPeriod := defaultDexHoldingPeriod
	record := model.DexWalletRecord{
		WalletAddress:          walletAddress,
		APIKey:                 apiKey,
		PrivateKey:             privateKey,
		TradingAccount:         tradingAccount,
		Priority:               priority,
		ExecutionFee:           executionFee,
		Leverage:               leverage,
		PositionSizePercentage: positionSize,
		WalletName:             req.WalletName,
		SlPercentage:           &slPercentage,
		HoldingHourPeriod:      &defaultHoldingPeriod,
		Exchange:               exchange,
	}

	return s.repo.InsertDexWallet(ctx, uid, record)
}

func (s *DexService) ListWallets(ctx context.Context, uid string, exchange string) ([]model.WalletInfo, error) {
	exchange = normalizeDexExchange(exchange)
	if exchange == "" {
		return nil, model.ErrDexInvalidExchange
	}

	records, err := s.repo.ListDexWallets(ctx, uid, exchange)
	if err != nil {
		return nil, err
	}

	var result []model.WalletInfo
	for _, r := range records {
		authors, err := s.repo.GetSubscribeAuthor(ctx, r.ID)
		if err != nil {
			logger.Errorf("dex list wallets: failed to fetch subscribed authors for wallet %s exchange=%s: %v", r.ID, exchange, err)
			return nil, err
		}
		sl := 0.0
		if r.SlPercentage != nil {
			sl = *r.SlPercentage
		}
		tp := 0.0
		if r.TpPercentage != nil {
			tp = *r.TpPercentage
		}
		holding := defaultDexHoldingPeriod
		if r.HoldingHourPeriod != nil {
			holding = *r.HoldingHourPeriod
		}

		info := model.WalletInfo{
			WalletName:             derefOrDefault(r.WalletName, ""),
			WalletID:               r.ID,
			WalletAddress:          r.WalletAddress,
			WalletType:             r.Exchange,
			Balance:                0,
			PrivyWalletID:          "",
			Priority:               r.Priority,
			Authors:                authors,
			TpPercentage:           tp,
			SlPercentage:           sl,
			HoldingHourPeriod:      holding,
			PositionSizePercentage: r.PositionSizePercentage,
			Leverage:               r.Leverage,
			HyperliquidBasecode:    false,
			CreatedAt:              r.CreatedAt,
			UpdatedAt:              r.UpdatedAt,
			DeletedAt:              r.DeletedAt,
		}
		result = append(result, info)
	}

	return result, nil
}

func (s *DexService) ActiveWallet(ctx context.Context, uid, walletID string) error {
	return s.repo.ActiveDexWallet(ctx, uid, walletID)
}

func (s *DexService) DeactiveWallet(ctx context.Context, uid, walletID string) error {
	return s.repo.DeactiveDexWallet(ctx, uid, walletID)
}

func (s *DexService) UpdatePositionSize(ctx context.Context, uid, walletID string, positionSize float64) error {
	if positionSize <= minDexPositionSize || positionSize > maxDexPositionSize {
		return model.ErrDexInvalidPosition
	}
	return s.repo.UpdateDexWalletPositionSize(ctx, uid, walletID, positionSize)
}

func (s *DexService) UpdateLeverage(ctx context.Context, uid, walletID string, leverage int, exchange string) error {
	if leverage < minDexLeverage || leverage > maxDexLeverage {
		return model.ErrDexInvalidLeverage
	}
	exchange = normalizeDexExchange(exchange)
	if exchange == "" {
		return model.ErrDexInvalidExchange
	}
	return s.repo.UpdateDexWalletLeverage(ctx, uid, walletID, leverage, exchange)
}

func (s *DexService) UpdateSL(ctx context.Context, uid, walletID string, slPercentage float64, exchange string) error {
	if slPercentage < minDexStopLoss || slPercentage > maxDexStopLoss {
		return model.ErrDexInvalidSL
	}
	exchange = normalizeDexExchange(exchange)
	if exchange == "" {
		return model.ErrDexInvalidExchange
	}
	return s.repo.UpdateDexWalletSL(ctx, uid, walletID, slPercentage, exchange)
}

func (s *DexService) UpdateHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error {
	if holdingPeriod < 0 {
		return fmt.Errorf("holding period must be non-negative")
	}
	return s.repo.UpdateDexWalletHoldingPeriod(ctx, uid, walletID, holdingPeriod)
}

func (s *DexService) GetWalletTotalValue(ctx context.Context, uid, walletID, exchange string) (model.DexWalletTotalValue, error) {
	exchange = normalizeDexExchange(exchange)
	if exchange == "" {
		return model.DexWalletTotalValue{}, model.ErrDexInvalidExchange
	}
	return s.repo.GetDexWalletTotalValue(ctx, uid, walletID, exchange)
}

func (s *DexService) ValidateDexCredentials(ctx context.Context, exchange, apiKey, privateKey, tradingAccountID string) (bool, error) {
	exchange = normalizeDexExchange(exchange)
	apiKey = strings.TrimSpace(apiKey)
	privateKey = strings.TrimSpace(privateKey)
	tradingAccountID = strings.TrimSpace(tradingAccountID)
	if exchange == "" {
		return false, model.ErrDexInvalidExchange
	}
	if apiKey == "" || privateKey == "" || tradingAccountID == "" {
		return false, model.ErrDexMissingFields
	}
	return s.repo.ValidateDexCredentials(ctx, exchange, apiKey, privateKey, tradingAccountID)
}

func (s *DexService) UpdateAPICredentials(ctx context.Context, uid, walletID, apiKey, privateKey, tradingAccountID, exchange string) error {
	apiKey = strings.TrimSpace(apiKey)
	privateKey = strings.TrimSpace(privateKey)
	tradingAccountID = strings.TrimSpace(tradingAccountID)
	exchange = normalizeDexExchange(exchange)
	walletID = strings.TrimSpace(walletID)

	if walletID == "" || apiKey == "" || privateKey == "" || tradingAccountID == "" {
		return model.ErrDexMissingFields
	}
	if exchange == "" {
		return model.ErrDexInvalidExchange
	}
	valid, err := s.ValidateDexCredentials(ctx, exchange, apiKey, privateKey, tradingAccountID)
	if err != nil {
		if errors.Is(err, model.ErrDexInvalidExchange) || errors.Is(err, model.ErrDexMissingFields) {
			return err
		}
		logger.Errorf("dex update api credentials: validation failed exchange=%s wallet_id=%s err=%v", exchange, walletID, err)
		return fmt.Errorf("failed to validate api credentials")
	}
	if !valid {
		return model.ErrDexInvalidCredentials
	}
	return s.repo.UpdateDexWalletAPICredentials(ctx, uid, walletID, apiKey, privateKey, tradingAccountID, exchange)
}

func (s *DexService) SubscribeAuthor(ctx context.Context, author string, walletID string) (string, error) {
	return s.repo.SubscribeAuthor(ctx, author, walletID)
}

func (s *DexService) UnsubscribeAuthor(ctx context.Context, author string, walletID string) error {
	return s.repo.UnsubscribeAuthor(ctx, author, walletID)
}

func deriveWalletAddress(privateKey string) (string, error) {
	key := strings.TrimSpace(privateKey)
	if !strings.HasPrefix(key, "0x") {
		return "", model.ErrDexInvalidKey
	}
	key = strings.TrimPrefix(key, "0x")
	if len(key) != 64 {
		return "", model.ErrDexInvalidKey
	}

	ecdsaKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}
	address := crypto.PubkeyToAddress(ecdsaKey.PublicKey)
	return strings.ToLower(address.Hex()), nil
}

func derefOrDefault[T comparable](ptr *T, fallback T) T {
	if ptr == nil {
		return fallback
	}
	return *ptr
}
