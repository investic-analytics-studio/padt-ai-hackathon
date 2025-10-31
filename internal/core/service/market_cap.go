package service

import (
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type MarketCapService struct {
	repo port.MarketCapRepo
}

func NewMarketCapService(repo port.MarketCapRepo) port.MarketCapService {
	return &MarketCapService{repo: repo}
}

func (s *MarketCapService) GetMarketCap() (model.MarketCapResponse, error) {
	return s.repo.GetMarketCap()
}
