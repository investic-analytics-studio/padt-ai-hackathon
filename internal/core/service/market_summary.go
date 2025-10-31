package service

import (
	"reflect"
	"strings"
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
	"go.uber.org/zap"
)

type marketSummaryService struct {
	logger logger.Logger
	repo   port.MarketSummaryRepo
}

func NewMarketSummaryService(repo port.MarketSummaryRepo, logger logger.Logger) domain.MarketSummaryService {
	return &marketSummaryService{
		logger: logger,
		repo:   repo,
	}
}

func (s *marketSummaryService) GetDailyMarketOverview() ([]domain.DailyMarketOverviewData, error) {
	data, err := s.repo.GetDailyMarketOverview()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	var result = make([]domain.DailyMarketOverviewData, 0)
	for _, item := range data {
		appear := false
		for i, r := range result {
			if r.Market == item.Market {
				result[i].Data = append(r.Data, item)
				appear = true
				break
			}
		}

		if !appear {
			result = append(result, domain.DailyMarketOverviewData{
				Market: item.Market,
				Data:   []model.DailyMarketOverview{item},
			})
		}
	}

	return result, nil
}

func (s *marketSummaryService) GetTopGainer() ([]model.TopGainer, error) {
	result, err := s.repo.GetTopGainer()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return result, nil
}

func (s *marketSummaryService) GetTopLoser() ([]model.TopLoser, error) {
	result, err := s.repo.GetTopLoser()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return result, nil
}

func (s *marketSummaryService) GetAdvancerDeclinerDistributionHist() ([]domain.AdvancerDeclinerDistributionHist, error) {
	data, err := s.repo.GetAdvancerDeclinerDistribution()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	result := make([]domain.AdvancerDeclinerDistributionHist, 0)
	reflected := reflect.ValueOf(*data)
	for i := 0; i < reflected.NumField(); i++ {
		val := reflected.Field(i).Int()
		jsonTag := reflected.Type().Field(i).Tag.Get("json")
		item := domain.AdvancerDeclinerDistributionHist{
			Count: val,
			Index: i,
		}

		switch strings.Split(jsonTag, "_")[0] {
		case "up":
			item.Col = "up"
		case "down":
			item.Col = "down"
		case "even":
			item.Col = "even"
		}

		result = append(result, item)
	}

	return result, nil
}

func (s *marketSummaryService) GetAdvancerDeclinerDistributionBar() (*domain.AdvancerDeclinerDistributionBar, error) {
	data, err := s.repo.GetAdvancerDeclinerDistribution()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	var sumDecline int64 = 0
	var sumAdvanc int64 = 0
	var sumNeutral int64 = 0
	reflected := reflect.ValueOf(*data)
	for i := 0; i < reflected.NumField(); i++ {
		val := reflected.Field(i).Int()
		jsonTag := reflected.Type().Field(i).Tag.Get("json")

		switch strings.Split(jsonTag, "_")[0] {
		case "up":
			sumAdvanc += val
		case "down":
			sumDecline += val
		case "even":
			sumNeutral += val
		}
	}

	total := sumAdvanc + sumDecline + sumNeutral
	result := new(domain.AdvancerDeclinerDistributionBar)
	result.Advanc = (float64(sumAdvanc) / float64(total)) * 100
	result.Decline = (float64(sumDecline) / float64(total)) * 100
	result.Neutral = (float64(sumNeutral) / float64(total)) * 100
	return result, nil
}

func (s *marketSummaryService) GetTopTurnoverFloat() ([]model.TurnoverFloat, error) {
	result, err := s.repo.GetTopTurnoverFloat()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return result, nil
}

func (s *marketSummaryService) GetCMEGoldOI() ([]model.CMEGoldOI, error) {
	result, err := s.repo.GetCMEGoldOI()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetStockAlertPlots() ([]model.StockAlertPlot, error) {
	result, err := s.repo.GetStockAlertPlots()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetStockAlertStatsDates() ([]time.Time, error) {
	result, err := s.repo.GetStockAlertStatsDates()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetStockAlertStats() ([]model.StockAlertStats, error) {
	result, err := s.repo.GetStockAlertStats()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetStockAlertDetections() ([]model.StockAlertDetection, error) {
	result, err := s.repo.GetStockAlertDetections()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetStockAlertDetections service result count", zap.Int("count", len(result)))
	if len(result) == 0 {
		s.logger.Info("No stock alert detections found")
	}
	return result, nil
}

func (s *marketSummaryService) GetCryptoAlertPlots() ([]model.CryptoAlertPlot, error) {
	result, err := s.repo.GetCryptoAlertPlots()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetCryptoAlertStatsDates() ([]time.Time, error) {
	result, err := s.repo.GetCryptoAlertStatsDates()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetCryptoAlertDetections() ([]model.CryptoAlertDetection, error) {
	result, err := s.repo.GetCryptoAlertDetections()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetCryptoAlertDetections service result count", zap.Int("count", len(result)))
	if len(result) == 0 {
		s.logger.Info("No crypto alert detections found")
	}
	return result, nil
}

// Add these new methods to the existing file

func (s *marketSummaryService) GetCryptoAlertPlots1D() ([]model.CryptoAlertPlot1D, error) {
	result, err := s.repo.GetCryptoAlertPlots1D()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetCryptoAlertStatsDates1D() ([]time.Time, error) {
	result, err := s.repo.GetCryptoAlertStatsDates1D()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *marketSummaryService) GetCryptoAlertDetections1D() ([]model.CryptoAlertDetection1D, error) {
	result, err := s.repo.GetCryptoAlertDetections1D()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetCryptoAlertDetections1D service result count", zap.Int("count", len(result)))
	if len(result) == 0 {
		s.logger.Info("No crypto alert detections 1D found")
	}
	return result, nil
}

func (s *marketSummaryService) GetCryptoAlertStats(page int, limit int) (*domain.PaginatedResponse, error) {
	result, total, err := s.repo.GetCryptoAlertStats(page, limit)
	if err != nil {
		return nil, err
	}
	hasNextPage := (int64(page+1) * int64(limit)) < total
	return &domain.PaginatedResponse{
		Data:        result,
		Total:       total,
		HasNextPage: hasNextPage,
	}, nil
}

func (s *marketSummaryService) GetCryptoAlertStats1D(page int, limit int) (*domain.PaginatedResponse, error) {
	result, total, err := s.repo.GetCryptoAlertStats1D(page, limit)
	if err != nil {
		return nil, err
	}
	hasNextPage := (int64(page+1) * int64(limit)) < total
	return &domain.PaginatedResponse{
		Data:        result,
		Total:       total,
		HasNextPage: hasNextPage,
	}, nil
}

func (s *marketSummaryService) GetAllCryptoAlertStats() ([]model.CryptoAlertStats, error) {
	result, err := s.repo.GetAllCryptoAlertStats()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetAllCryptoAlertStats service result count", zap.Int("count", len(result)))
	return result, nil
}

func (s *marketSummaryService) GetAllCryptoAlertStats1D() ([]model.CryptoAlertStats1D, error) {
	result, err := s.repo.GetAllCryptoAlertStats1D()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetAllCryptoAlertStats1D service result count", zap.Int("count", len(result)))
	return result, nil
}
