package service

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

//	type PerformanceService struct {
//		logger logger.Logger
//		repo   port.PerformanceRepo
//	}
type PerformanceService struct {
	repo *repo.PerformanceRepo
}

func NewPerformanceService(repo *repo.PerformanceRepo) *PerformanceService {
	return &PerformanceService{repo: repo}
}
func (s *PerformanceService) GetMultiholdingPortNav(ctx context.Context, period string, holdingPeriod string) ([]model.AuthorNav, error) {
	result, err := s.repo.MultiholdingPortNavRepo(ctx, period, holdingPeriod)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *PerformanceService) GetAuthorNav(ctx context.Context, period string) ([]model.AuthorNav, error) {
	result, err := s.repo.AuthorNavRepo(ctx, period)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PerformanceService) GetAuthorDetail(ctx context.Context, authorUsername string, period string, start int, limit int) (model.AuthorDetail, error) {
	// Get the basic author detail data
	result, err := s.repo.AuthorDetailRepo(ctx, authorUsername, period, start, limit)
	if err != nil {
		return model.AuthorDetail{}, err
	}

	// Get sentiment analysis data for the period
	bearishTokens, bullishTokens, err := s.repo.GetAuthorSentimentAnalysis(ctx, authorUsername, period)
	if err != nil {
		// Log warning but continue with empty sentiment data
		bearishTokens = []model.SentimentToken{}
		bullishTokens = []model.SentimentToken{}
	}

	// Ensure we always have initialized slices (never null)
	if bearishTokens == nil {
		bearishTokens = []model.SentimentToken{}
	}
	if bullishTokens == nil {
		bullishTokens = []model.SentimentToken{}
	}

	// Add sentiment data to the result
	result.BearishTokens = bearishTokens
	result.BullishTokens = bullishTokens

	return result, nil
}

func (s *PerformanceService) GetAuthorMultiholdingDetail(ctx context.Context, authorUsername string, period string, holdingPeriod string, start int, limit int) (model.AuthorDetail, error) {
	// Get the basic author detail data with multiholding nav
	result, err := s.repo.AuthorMultiholdingDetailRepo(ctx, authorUsername, period, holdingPeriod, start, limit)
	if err != nil {
		return model.AuthorDetail{}, err
	}

	// Get sentiment analysis data for the period
	bearishTokens, bullishTokens, err := s.repo.GetAuthorSentimentAnalysis(ctx, authorUsername, period)
	if err != nil {
		// Log warning but continue with empty sentiment data
		bearishTokens = []model.SentimentToken{}
		bullishTokens = []model.SentimentToken{}
	}

	// Ensure we always have initialized slices (never null)
	if bearishTokens == nil {
		bearishTokens = []model.SentimentToken{}
	}
	if bullishTokens == nil {
		bullishTokens = []model.SentimentToken{}
	}

	// Add sentiment data to the result
	result.BearishTokens = bearishTokens
	result.BullishTokens = bullishTokens

	return result, nil
}
