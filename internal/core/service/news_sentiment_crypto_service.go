package service

import (
	"context"
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type NewsSentimentCryptoService struct {
	repo port.NewsSentimentCryptoRepo
}

func NewNewsSentimentCryptoService(repo port.NewsSentimentCryptoRepo) *NewsSentimentCryptoService {
	return &NewsSentimentCryptoService{repo: repo}
}

func (s *NewsSentimentCryptoService) GetNewsSentiment(ctx context.Context, date time.Time) ([]model.NewsSentimentEntities, error) {
	return s.repo.GetNewsSentiment(ctx, date)
}
