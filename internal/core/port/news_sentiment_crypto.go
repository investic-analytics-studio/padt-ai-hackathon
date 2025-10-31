package port

import (
	"context"
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type NewsSentimentCryptoRepo interface {
	GetNewsSentiment(ctx context.Context, date time.Time) ([]model.NewsSentimentEntities, error)
}

type NewsSentimentCryptoService interface {
	GetNewsSentiment(ctx context.Context, date time.Time) ([]model.NewsSentimentEntities, error)
}
