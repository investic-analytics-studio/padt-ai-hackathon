package service

import (
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type TopRankService struct {
	repo port.TopRankRepo
}

func NewTopRankService(repo port.TopRankRepo) *TopRankService {
	return &TopRankService{repo: repo}
}

func (s *TopRankService) GetTop100() (model.TopRankResponse, error) {
	return s.repo.GetTop100()
}
