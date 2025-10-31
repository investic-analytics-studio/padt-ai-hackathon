package domain

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type SentimentCryptoService interface {
	GetSentimentAggregateCountRow(ctx context.Context, ticker, timeRange, sourceNames, topics string) (int, error)
	GetSentimentAggregate(ctx context.Context, ticker, timeRange, sourceNames, topics, limit, offset string) ([]model.SentimentCrypto, error)
	GetUniqueTickers(ctx context.Context) ([]model.TickerInfo, error)
	GetUniqueTopics(ctx context.Context) ([]string, error)
	GetUniqueSourceNames(ctx context.Context) ([]string, error)
}

type SentimentCryptoRepository interface {
	GetSentimentAggregateCountRow(ctx context.Context, ticker, timeRange, sourceNames, topics string) (int, error)
	GetSentimentAggregate(ctx context.Context, ticker, timeRange, sourceNames, topics, limit, offset string) ([]model.SentimentCrypto, error)
	GetUniqueTickers(ctx context.Context) ([]model.TickerInfo, error)
	GetUniqueTopics(ctx context.Context) ([]string, error)
	GetUniqueSourceNames(ctx context.Context) ([]string, error)
}

// SentimentAggregate struct can be removed if it's not used elsewhere
