package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type PerformanceRepo interface {
	AuthorNavRepo(ctx context.Context, period string) ([]model.AuthorNav, error)
	AuthorDetailRepo(ctx context.Context, authorUsername string, period string, start int, limit int) (model.AuthorDetail, error)
	AuthorMultiholdingDetailRepo(ctx context.Context, authorUsername string, period string, holdingPeriod string, start int, limit int) (model.AuthorDetail, error)
	GetAuthorSentimentAnalysis(ctx context.Context, authorUsername string, period string) (bearishTokens []model.SentimentToken, bullishTokens []model.SentimentToken, err error)
	MultiholdingPortNavRepo(ctx context.Context, period string, holdingPeriod string) ([]model.AuthorNav, error)
}

type PerformanceService interface {
	GetAuthorNav(ctx context.Context, period string) ([]model.AuthorNav, error)
	GetAuthorDetail(ctx context.Context, authorUsername string, period string, start int, limit int) (model.AuthorDetail, error)
	GetAuthorMultiholdingDetail(ctx context.Context, authorUsername string, period string, holdingPeriod string, start int, limit int) (model.AuthorDetail, error)
	GetMultiholdingPortNav(ctx context.Context, period string, holdingPeriod string) ([]model.AuthorNav, error)
}
