package port

import (
	"context"
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type TwitterCryptoRepo interface {
	GetAllSentiments(ctx context.Context) ([]model.TwitterCryptoSentiment, error)
	GetAllTweets(ctx context.Context) ([]model.TwitterCryptoTweet, error)
	GetAuthorProfiles(ctx context.Context) ([]model.TwitterCryptoAuthorProfile, error)
	GetAuthorWinrate(ctx context.Context, selectedWinratePeriod string) ([]model.TwitterCryptoAuthorWinRate, error)
	GetPaginatedTweets(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweet, int, error)
	GetPaginatedSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoSentiment, int, error)
	GetTweetsWithSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweetWithSentiment, int, error)
	GetTweetsWithSentimentAuthorSignal(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAuthorAndSignal, int, error)
	GetTweetsWithSentimentsAndAuthor(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time, searchTokenSymbolValue string) ([]model.TwitterCryptoTweetWithSentimentAndAuthor, int, error)
	GetTweetsWithSentimentsAndTier(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAndTier, int, error)
	GetSummaries(ctx context.Context, start, limit int, sortBy, sortOrder string, fromDate, toDate *time.Time) ([]model.TwitterCryptoSummary, int, error)
	GetBubbleSentiment(ctx context.Context) (model.BubbleSentimentModel, error)
	SearchTokenMentionSymbolsByAuthors(ctx context.Context, symbol string, createdAt string, id string, limit int, selectedTimeRange string, authors []string) ([]model.TokenMentionSymbolAndAuthor, time.Time, string, int, error)
}

type TwitterCryptoService interface {
	GetAllSentiments(ctx context.Context) ([]model.TwitterCryptoSentiment, error)
	GetAllTweets(ctx context.Context) ([]model.TwitterCryptoTweet, error)
	GetAuthorProfiles(ctx context.Context) ([]model.TwitterCryptoAuthorProfile, error)
	GetAuthorWinrate(ctx context.Context, selectedWinratePeriod string) ([]model.TwitterCryptoAuthorWinRate, error)
	GetPaginatedTweets(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweet, int, error)
	GetPaginatedSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoSentiment, int, error)
	GetTweetsWithSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweetWithSentiment, int, error)
	GetTweetsWithSentimentsAndAuthor(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time, searchTokenSymbolValue string) ([]model.TwitterCryptoTweetWithSentimentAndAuthor, int, error)
	GetTweetsWithSentimentAuthorSignal(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAuthorAndSignal, int, error)
	GetTweetsWithSentimentsAndTier(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAndTier, int, error)
	GetSummaries(ctx context.Context, start, limit int, sortBy, sortOrder string, fromDate, toDate *time.Time) ([]model.TwitterCryptoSummary, int, error)
	GetBubbleSentiment(ctx context.Context) (model.BubbleSentimentModel, error)
	SearchTokenMentionSymbolsByAuthors(ctx context.Context, symbol string, createdAt string, id string, limit int, selectedTimeRange string, authors []string) ([]model.TokenMentionSymbolAndAuthor, time.Time, string, int, error)
}
