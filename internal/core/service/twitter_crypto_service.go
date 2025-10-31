package service

import (
	"context"
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type TwitterCryptoService struct {
	repo port.TwitterCryptoRepo
}

func NewTwitterCryptoService(repo port.TwitterCryptoRepo) port.TwitterCryptoService {
	return &TwitterCryptoService{repo: repo}
}

func (s *TwitterCryptoService) GetAllSentiments(ctx context.Context) ([]model.TwitterCryptoSentiment, error) {
	return s.repo.GetAllSentiments(ctx)
}

func (s *TwitterCryptoService) GetAllTweets(ctx context.Context) ([]model.TwitterCryptoTweet, error) {
	return s.repo.GetAllTweets(ctx)
}

func (s *TwitterCryptoService) GetAuthorProfiles(ctx context.Context) ([]model.TwitterCryptoAuthorProfile, error) {
	return s.repo.GetAuthorProfiles(ctx)
}

func (s *TwitterCryptoService) GetAuthorWinrate(ctx context.Context, selectedWinratePeriod string) ([]model.TwitterCryptoAuthorWinRate, error) {
	return s.repo.GetAuthorWinrate(ctx, selectedWinratePeriod)
}

func (s *TwitterCryptoService) GetPaginatedTweets(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweet, int, error) {
	return s.repo.GetPaginatedTweets(ctx, start, limit, sortBy, sortOrder)
}

func (s *TwitterCryptoService) GetPaginatedSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoSentiment, int, error) {
	return s.repo.GetPaginatedSentiments(ctx, start, limit, sortBy, sortOrder)
}

func (s *TwitterCryptoService) GetTweetsWithSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweetWithSentiment, int, error) {
	return s.repo.GetTweetsWithSentiments(ctx, start, limit, sortBy, sortOrder)
}

func (s *TwitterCryptoService) GetTweetsWithSentimentsAndAuthor(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time, searchTokenSymbolValue string) ([]model.TwitterCryptoTweetWithSentimentAndAuthor, int, error) {
	return s.repo.GetTweetsWithSentimentsAndAuthor(ctx, start, limit, sortBy, sortOrder, authors, fromDate, toDate, searchTokenSymbolValue)
}
func (s *TwitterCryptoService) GetTweetsWithSentimentAuthorSignal(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAuthorAndSignal, int, error) {
	return s.repo.GetTweetsWithSentimentAuthorSignal(ctx, start, limit, sortBy, sortOrder, authors, fromDate, toDate)
}
func (s *TwitterCryptoService) GetTweetsWithSentimentsAndTier(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAndTier, int, error) {
	tweets, total, err := s.repo.GetTweetsWithSentimentsAndTier(ctx, start, limit, sortBy, sortOrder, authors, fromDate, toDate)
	if err != nil {
		return nil, 0, err
	}

	return tweets, total, nil
}

func (s *TwitterCryptoService) GetSummaries(ctx context.Context, start, limit int, sortBy, sortOrder string, fromDate, toDate *time.Time) ([]model.TwitterCryptoSummary, int, error) {
	return s.repo.GetSummaries(ctx, start, limit, sortBy, sortOrder, fromDate, toDate)
}
func (s *TwitterCryptoService) GetBubbleSentiment(ctx context.Context) (model.BubbleSentimentModel, error) {
	return s.repo.GetBubbleSentiment(ctx)
}

func (s *TwitterCryptoService) SearchTokenMentionSymbolsByAuthors(ctx context.Context, symbol string, createdAt string, id string, limit int, selectedTimeRange string, authors []string) ([]model.TokenMentionSymbolAndAuthor, time.Time, string, int, error) {
	return s.repo.SearchTokenMentionSymbolsByAuthors(ctx, symbol, createdAt, id, limit, selectedTimeRange, authors)
}
