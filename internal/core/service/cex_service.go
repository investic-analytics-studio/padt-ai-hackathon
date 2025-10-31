package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

const (
	defaultCexPriority      = 2
	defaultCexExecutionFee  = 0.1
	defaultCexPositionSize  = 0.1
	defaultCexLeverage      = 1
	defaultCexStopLoss      = 0.0
	defaultCexHoldingPeriod = 48
	minCexLeverage          = 1
	maxCexLeverage          = 100
	positionSizeLowerBound  = 0.0
	positionSizeUpperBound  = 1.0
	stopLossLowerBound      = 0.0
	stopLossUpperBound      = 100.0
)

type CexService struct {
	repo port.CexRepo
}

func NewCexService(repo port.CexRepo) *CexService {
	return &CexService{repo: repo}
}

func normalizeExchangeName(exchange string) string {
	trimmed := strings.TrimSpace(strings.ToLower(exchange))
	if trimmed == "" {
		return ""
	}
	normalized := strings.Builder{}
	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			normalized.WriteRune(r)
		}
	}
	return normalized.String()
}

func (s *CexService) Connect(ctx context.Context, uid string, req model.CexConnectRequest) (string, error) {
	apiKey := strings.TrimSpace(req.APIKey)
	apiSecret := strings.TrimSpace(req.APISecret)
	exchange := normalizeExchangeName(req.Exchange)

	if exchange == "" {
		return "", model.ErrCexInvalidExchange
	}
	
	valid, err := s.repo.ValidateCexCredentials(ctx, exchange, apiKey, apiSecret)
	if err != nil {
		logger.Errorf("cex connect: credential validation failed exchange=%s err=%v", exchange, err)
		return "", fmt.Errorf("failed to validate api credentials")
	}
	if !valid {
		return "", model.ErrCexInvalidCredentials
	}
	
	walletAddress := strings.ToLower(uid + "." + exchange)

	if walletAddress == "" || apiKey == "" || apiSecret == "" {
		return "", model.ErrCexMissingFields
	}

	exists, err := s.repo.WalletExistsByAddress(ctx, uid, walletAddress, exchange)
	if err != nil {
		return "", err
	}
	if exists {
		return "", model.ErrCexWalletExists
	}

	priority := defaultCexPriority
	executionFee := defaultCexExecutionFee
	positionSize := defaultCexPositionSize

	if positionSize <= positionSizeLowerBound || positionSize > positionSizeUpperBound {
		return "", model.ErrCexInvalidPosition
	}

	leverage := defaultCexLeverage
	if leverage < minCexLeverage || leverage > maxCexLeverage {
		return "", model.ErrCexInvalidLeverage
	}
	sl := defaultCexStopLoss
	if sl < stopLossLowerBound || sl > stopLossUpperBound {
		return "", model.ErrCexInvalidSL
	}

	defaultHoldingPeriod := defaultCexHoldingPeriod
	record := model.CexWalletRecord{
		WalletAddress:          walletAddress,
		APIKey:                 apiKey,
		APISecret:              apiSecret,
		Priority:               priority,
		Exchange:               exchange,
		ExecutionFee:           executionFee,
		PositionSizePercentage: positionSize,
		Leverage:               leverage,
		WalletName:             req.WalletName,
		SlPercentage:           &sl,
		HoldingHourPeriod:      &defaultHoldingPeriod,
	}

	return s.repo.InsertCexWallet(ctx, uid, record)
}

func (s *CexService) ListWallets(ctx context.Context, uid string, exchange string) ([]model.CexWalletInfo, error) {
	exchange = normalizeExchangeName(exchange)
	if exchange == "" {
		return nil, model.ErrCexInvalidExchange
	}

	records, err := s.repo.ListCexWallets(ctx, uid, exchange)
	if err != nil {
		return nil, err
	}

	var result []model.CexWalletInfo
	for _, r := range records {
		authors, err := s.repo.GetSubscribeAuthor(ctx, r.ID)
		if err != nil {
			logger.Errorf("cex service: failed to fetch subscribed authors for wallet %s: %v", r.ID, err)
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

		holding := defaultCexHoldingPeriod
		if r.HoldingHourPeriod != nil {
			holding = *r.HoldingHourPeriod
		}

		walletName := ""
		if r.WalletName != nil {
			walletName = *r.WalletName
		}

		info := model.CexWalletInfo{
			WalletName:             walletName,
			WalletID:               r.ID,
			WalletType:             r.Exchange,
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

func (s *CexService) ActiveWallet(ctx context.Context, uid, walletID string) error {
	return s.repo.ActiveCexWallet(ctx, uid, walletID)
}

func (s *CexService) DeactiveWallet(ctx context.Context, uid, walletID string) error {
	return s.repo.DeactiveCexWallet(ctx, uid, walletID)
}

func (s *CexService) UpdateHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error {
	return s.repo.UpdateCexWalletHoldingPeriod(ctx, uid, walletID, holdingPeriod)
}

func (s *CexService) UpdatePositionSize(ctx context.Context, uid, walletID string, positionSize float64) error {
	if positionSize <= positionSizeLowerBound || positionSize > positionSizeUpperBound {
		return model.ErrCexInvalidPosition
	}
	return s.repo.UpdateCexWalletPositionSize(ctx, uid, walletID, positionSize)
}

func (s *CexService) UpdateLeverage(ctx context.Context, uid, walletID string, leverage int) error {
	if leverage < minCexLeverage || leverage > maxCexLeverage {
		return model.ErrCexInvalidLeverage
	}
	return s.repo.UpdateCexWalletLeverage(ctx, uid, walletID, leverage)
}

func (s *CexService) UpdateAPICredentials(ctx context.Context, uid, walletID string, apiKey string, apiSecret string, exchange string) error {
	apiKey = strings.TrimSpace(apiKey)
	apiSecret = strings.TrimSpace(apiSecret)
	exchange = normalizeExchangeName(exchange)
	if walletID == "" || apiKey == "" || apiSecret == "" {
		return model.ErrCexMissingCredentials
	}
	if exchange == "" {
		return model.ErrCexInvalidExchange
	}
	if exchange == "binance-th" {
		valid, err := s.repo.ValidateCexCredentials(ctx, exchange, apiKey, apiSecret)
		if err != nil {
			logger.Errorf("cex update api key: credential validation failed exchange=%s err=%v", exchange, err)
			return fmt.Errorf("failed to validate api credentials")
		}
		if !valid {
			return model.ErrCexInvalidCredentials
		}
	}
	return s.repo.UpdateCexWalletAPICredentials(ctx, uid, walletID, apiKey, apiSecret, exchange)
}

func (s *CexService) UpdateSL(ctx context.Context, uid, walletID string, sl float64, exchange string) error {
	if sl < stopLossLowerBound || sl > stopLossUpperBound {
		return model.ErrCexInvalidSL
	}
	exchange = normalizeExchangeName(exchange)
	if exchange == "" {
		return model.ErrCexInvalidExchange
	}
	return s.repo.UpdateCexWalletSL(ctx, uid, walletID, sl, exchange)
}

func (s *CexService) SubscribeAuthor(ctx context.Context, author, walletID string) (string, error) {
	return s.repo.SubscribeAuthor(ctx, author, walletID)
}

func (s *CexService) UnsubscribeAuthor(ctx context.Context, author, walletID string) error {
	return s.repo.UnsubscribeAuthor(ctx, author, walletID)
}

func (s *CexService) GetWalletTotalValue(ctx context.Context, uid, walletID, exchange string) (model.CexWalletTotalValue, error) {
	exchange = normalizeExchangeName(exchange)
	if exchange == "" {
		return model.CexWalletTotalValue{}, model.ErrCexInvalidExchange
	}
	return s.repo.GetCexWalletTotalValue(ctx, uid, walletID, exchange)
}
