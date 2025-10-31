package service

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type sentimentCryptoService struct {
	repo   domain.SentimentCryptoRepository
	logger logger.Logger
}

func NewSentimentCryptoService(repo domain.SentimentCryptoRepository, logger logger.Logger) domain.SentimentCryptoService {
	return &sentimentCryptoService{
		repo:   repo,
		logger: logger,
	}
}
func (s *sentimentCryptoService) GetSentimentAggregateCountRow(ctx context.Context, ticker, timeRange, sourceNames, topics string) (int, error) {
	s.logger.Info("GetSentimentAggregateCountRow service called")
	totalRow, err := s.repo.GetSentimentAggregateCountRow(ctx, ticker, timeRange, sourceNames, topics)
	if err != nil {
		s.logger.Error(err)
		return 0, err
	}
	s.logger.Info("GetSentimentAggregateCountRow service successful")
	return totalRow, nil
}
func (s *sentimentCryptoService) GetSentimentAggregate(ctx context.Context, ticker, timeRange, sourceNames, topics, limit, offset string) ([]model.SentimentCrypto, error) {
	s.logger.Info("GetSentimentAggregate service called")
	data, err := s.repo.GetSentimentAggregate(ctx, ticker, timeRange, sourceNames, topics, limit, offset)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetSentimentAggregate service successful")
	return data, nil
}

func (s *sentimentCryptoService) GetUniqueTickers(ctx context.Context) ([]model.TickerInfo, error) {
	s.logger.Info("GetUniqueTickers service called")
	tickers, err := s.repo.GetUniqueTickers(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetUniqueTickers service successful")
	return tickers, nil
}

func (s *sentimentCryptoService) GetUniqueTopics(ctx context.Context) ([]string, error) {
	s.logger.Info("GetUniqueTopics service called")
	topics, err := s.repo.GetUniqueTopics(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetUniqueTopics service successful")
	return topics, nil
}

func (s *sentimentCryptoService) GetUniqueSourceNames(ctx context.Context) ([]string, error) {
	s.logger.Info("GetUniqueSourceNames service called")
	sourceNames, err := s.repo.GetUniqueSourceNames(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Info("GetUniqueSourceNames service successful")
	return sourceNames, nil
}
