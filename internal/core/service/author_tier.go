package service

import (
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type AuthorTierService struct {
	repo port.AuthorTierRepo
}

func NewAuthorTierService(repo port.AuthorTierRepo) *AuthorTierService {
	return &AuthorTierService{repo: repo}
}

func (s *AuthorTierService) GetAuthorsByTier(tier string) ([]string, error) {
	return s.repo.GetAuthorsByTier(tier)
}

func (s *AuthorTierService) GetAllTiers() ([]string, error) {
	return s.repo.GetAllTiers()
}
