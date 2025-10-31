package service

import (
	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type OHLCVDataService struct {
	logger logger.Logger
	repo   port.TimescaleRepo
}

func NewOHLCVDataService(repo port.TimescaleRepo, logger logger.Logger) domain.OHLCVDataService {
	return &OHLCVDataService{
		logger: logger,
		repo:   repo,
	}
}

func (s *OHLCVDataService) GetCryptoOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error) {
	result, err := s.repo.GetCryptoOHLCV(req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}

func (s *OHLCVDataService) GetForexOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error) {
	result, err := s.repo.GetForexOHLCV(req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return result, nil
}
