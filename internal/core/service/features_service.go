package service

import (
	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type FeaturesService struct {
	repo *repo.FeaturesRepo
}

func NewFeaturesService(repo *repo.FeaturesRepo) *FeaturesService {
	return &FeaturesService{repo: repo}
}
func (s *FeaturesService) GetFeature(featureName string) (model.FeaturesEntities, error) {
	return s.repo.GetFeature(featureName)
}
