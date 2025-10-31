package service

import (
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type winRateService struct {
	winRateRepo port.WinRateRepo
}

func NewWinRateService(winRateRepo port.WinRateRepo) *winRateService {
	return &winRateService{winRateRepo: winRateRepo}
}

func (s *winRateService) GetWinRate() ([]model.WinRateModel, error) {
	result, err := s.winRateRepo.GetWinRate()
	if err != nil {
		return nil, err
	}

	winRateModel := make([]model.WinRateModel, 0)
	for _, v := range result {
		winRateModel = append(winRateModel, model.WinRateModel{
			AuthorId:         v.AuthorId,
			AuthorUsername:   v.AuthorUsername,
			CryptoWinRate1D:  v.CryptoWinRate1D,
			CryptoWinRate3D:  v.CryptoWinRate3D,
			CryptoWinRate7D:  v.CryptoWinRate7D,
			CryptoWinRate15D: v.CryptoWinRate15D,
			CryptoWinRate30D: v.CryptoWinRate30D,
		})
	}
	return winRateModel, nil
}
