package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type SentimentCryptoRepo interface {
	GetSentimentAggregate(ctx context.Context, ticker string) (interface{}, error)
	GetUniqueTickers(ctx context.Context) ([]model.TickerInfo, error)
	GetUniqueTopics(ctx context.Context) ([]string, error)
	GetUniqueSourceNames(ctx context.Context) ([]string, error)
}

type SentimentCryptoService interface {
	GetSentimentAggregate(ctx context.Context, ticker string) (interface{}, error)
	GetUniqueTickers(ctx context.Context) ([]model.TickerInfo, error)
	GetUniqueTopics(ctx context.Context) ([]string, error)
	GetUniqueSourceNames(ctx context.Context) ([]string, error)
}
