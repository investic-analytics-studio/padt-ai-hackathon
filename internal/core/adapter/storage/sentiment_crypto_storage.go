package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type sentimentCryptoStorage struct {
	db     *sql.DB
	logger logger.Logger
}

func NewSentimentCryptoStorage(db *sql.DB, logger logger.Logger) port.SentimentCryptoRepo {
	return &sentimentCryptoStorage{
		db:     db,
		logger: logger,
	}
}

func (s *sentimentCryptoStorage) GetSentimentAggregate(ctx context.Context, ticker string) (interface{}, error) {
	s.logger.Info("GetSentimentAggregate storage called")

	query := `
		SELECT date, sentiment, COUNT(topics) as count, tickers
		FROM crypto_sentiment
		WHERE tickers != 'T'
	`

	args := make([]interface{}, 0)

	if ticker != "" {
		query += " AND tickers LIKE ?"
		args = append(args, "%"+ticker+"%")
	}

	query += `
		GROUP BY date, sentiment, tickers
		ORDER BY date DESC, tickers
	`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		var date time.Time
		var sentiment string
		var count int
		var tickers string

		err := rows.Scan(&date, &sentiment, &count, &tickers)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}

		row := map[string]interface{}{
			"date":      date,
			"sentiment": sentiment,
			"count":     count,
			"ticker":    tickers,
		}
		result = append(result, row)
	}

	if err = rows.Err(); err != nil {
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Info("GetSentimentAggregate storage successful")
	return result, nil
}

func (s *sentimentCryptoStorage) GetAllSentiments(ctx context.Context) ([]model.TwitterCryptoSentiment, error) {
	// Implement this method
	return nil, nil
}

func (s *sentimentCryptoStorage) GetAllTweets(ctx context.Context) ([]model.TwitterCryptoTweet, error) {
	// Implement this method
	return nil, nil
}

func (s *sentimentCryptoStorage) GetAuthorProfiles(ctx context.Context) ([]model.TwitterCryptoAuthorProfile, error) {
	// Implement this method
	return nil, nil
}

func (s *sentimentCryptoStorage) GetUniqueTickers(ctx context.Context) ([]model.TickerInfo, error) {
	// Implement this method
	return nil, nil
}

func (s *sentimentCryptoStorage) GetUniqueTopics(ctx context.Context) ([]string, error) {
	// Implement this method
	return nil, nil
}

func (s *sentimentCryptoStorage) GetUniqueSourceNames(ctx context.Context) ([]string, error) {
	// Implement this method
	return nil, nil
}
