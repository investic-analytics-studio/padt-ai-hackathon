package repo

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type sentimentAnalysisRepo struct {
	db *sqlx.DB
}

func NewSentimentAnalysisRepo(db *sqlx.DB) port.SentimentAnalysisRepo {
	return &sentimentAnalysisRepo{db: db}
}

func (r *sentimentAnalysisRepo) GetSentimentAnalysisByAuthorList(authorList []string, dateRange string) (model.SentimentAnalysisModel, error) {
	overallSentiment, err := r.getOverallSentiment(authorList, dateRange)
	if err != nil {
		return model.SentimentAnalysisModel{}, fmt.Errorf("failed to get overall sentiment: %w", err)
	}

	tokenSentiment, err := r.getTokenSentiment(authorList, dateRange)
	if err != nil {
		return model.SentimentAnalysisModel{}, fmt.Errorf("failed to get token sentiment: %w", err)
	}

	tweetIds, err := r.getTweetIdsPerSentiment(authorList, dateRange)
	if err != nil {
		return model.SentimentAnalysisModel{}, fmt.Errorf("failed to get tweet IDs: %w", err)
	}

	return model.SentimentAnalysisModel{
		OverallSentimentAnalysis: overallSentiment,
		TokenSentimentAnalysis:   tokenSentiment,
		TwitterIdResult:          tweetIds,
	}, nil
}

func (r *sentimentAnalysisRepo) GetSentimentAnalysisByTier(tier string, dateRange string) (model.SentimentAnalysisModel, error) {
	authorList, err := r.getAuthorListByTier(strings.ToUpper(tier))
	if err != nil {
		return model.SentimentAnalysisModel{}, fmt.Errorf("failed to get author list: %w", err)
	}

	overallSentiment, err := r.getOverallSentiment(authorList, dateRange)
	if err != nil {
		return model.SentimentAnalysisModel{}, fmt.Errorf("failed to get overall sentiment: %w", err)
	}

	tokenSentiment, err := r.getTokenSentiment(authorList, dateRange)
	if err != nil {
		return model.SentimentAnalysisModel{}, fmt.Errorf("failed to get token sentiment: %w", err)
	}

	return model.SentimentAnalysisModel{
		OverallSentimentAnalysis: overallSentiment,
		TokenSentimentAnalysis:   tokenSentiment,
	}, nil
}

func (r *sentimentAnalysisRepo) getAuthorListByTier(tier string) ([]string, error) {
	var authorList []string
	query := `
		SELECT author_username 
		FROM twitter_crypto_author_profile
		WHERE author_tier = $1
		AND is_select = True
	`
	if err := r.db.Select(&authorList, query, tier); err != nil {
		return nil, fmt.Errorf("failed to get author list: %w", err)
	}
	return authorList, nil
}

func (r *sentimentAnalysisRepo) getOverallSentiment(authorList []string, dateRange string) (model.OverallSentimentResult, error) {
	query := fmt.Sprintf(`
		SELECT 
			CASE 
				WHEN tcs.sentiment = 'Neutral' THEN 'neutral'
				WHEN tcs.sentiment = 'Bullish' THEN 'positive'
				WHEN tcs.sentiment = 'Bearish' THEN 'negative'
				ELSE tcs.sentiment
			END as sentiment,
			count(tcs.sentiment) as count
		FROM twitter_crypto_signal tcs
		JOIN twitter_crypto_tweets_foxhole tct ON tcs.tweet_id = tct.id
		WHERE tct.author_username = ANY($1)
		AND tct.tweet_created_at >= CURRENT_DATE - INTERVAL '%s days'
		AND tct.tweet_created_at <= CURRENT_DATE
		AND tcs.sentiment IS NOT null
		AND tcs.sentiment != 'NaN'
		AND tcs.ticker != 'NONE'
		GROUP BY tcs.sentiment
	`, dateRange)

	var results []model.OverallSentimentQueryResult
	if err := r.db.Select(&results, query, pq.Array(authorList)); err != nil {
		return model.OverallSentimentResult{}, fmt.Errorf("failed to get overall sentiment: %w", err)
	}

	var sentiment model.OverallSentimentResult
	for _, result := range results {
		switch result.Sentiment {
		case "positive":
			sentiment.PositiveCount = result.Count
		case "negative":
			sentiment.NegativeCount = result.Count
		case "neutral":
			sentiment.NeutralCount = result.Count
		}
	}
	return sentiment, nil
}

func (r *sentimentAnalysisRepo) getTokenSentiment(authorList []string, dateRange string) ([]model.TokenSentiment, error) {
	query := fmt.Sprintf(`
		SELECT 
			REPLACE(ticker, 'USDT', '') as ticker,
			CASE 
				WHEN tcs.sentiment = 'Neutral' THEN 'neutral'
				WHEN tcs.sentiment = 'Bullish' THEN 'positive'
				WHEN tcs.sentiment = 'Bearish' THEN 'negative'
				ELSE tcs.sentiment
			END as sentiment,
			count(sentiment) as count
		FROM twitter_crypto_signal tcs
		JOIN twitter_crypto_tweets_foxhole tct ON tcs.tweet_id = tct.id
		WHERE tct.author_username = ANY($1)
		AND tct.tweet_created_at >= CURRENT_DATE - INTERVAL '%s days'
		AND tct.tweet_created_at <= CURRENT_DATE
		AND tcs.sentiment IS NOT null
		AND tcs.sentiment != 'NaN'
		AND ticker != 'NONE'
		GROUP BY REPLACE(ticker, 'USDT', ''), tcs.sentiment 
	`, dateRange)

	var results []model.TickerSentimentQueryResult
	if err := r.db.Select(&results, query, pq.Array(authorList)); err != nil {
		return nil, fmt.Errorf("failed to get token sentiment: %w", err)
	}

	tickerMap := make(map[string]*model.TokenSentiment)
	for _, result := range results {
		sentiment, exists := tickerMap[result.Ticker]
		if !exists {
			sentiment = &model.TokenSentiment{
				Ticker: result.Ticker,
			}
			tickerMap[result.Ticker] = sentiment
		}

		switch result.Sentiment {
		case "positive":
			sentiment.PositiveCount = result.Count
		case "negative":
			sentiment.NegativeCount = result.Count
		case "neutral":
			sentiment.NeutralCount = result.Count
		}
	}

	tokenSentiments := make([]model.TokenSentiment, 0, len(tickerMap))
	for _, sentiment := range tickerMap {
		tokenSentiments = append(tokenSentiments, *sentiment)
	}
	return tokenSentiments, nil
}

func (r *sentimentAnalysisRepo) getTweetIdsPerSentiment(authorList []string, dateRange string) ([]model.TickerSentimentTwitterIdResult, error) {
	query := fmt.Sprintf(`
		WITH unique_tweets AS (
			SELECT DISTINCT ON (tcs.tweet_id) 
				tcs.tweet_id, 
				tcs.sentiment,
				REPLACE(tcs.ticker, 'USDT', '') as ticker
			FROM twitter_crypto_signal tcs
			JOIN twitter_crypto_tweets_foxhole tct ON tcs.tweet_id = tct.id
			WHERE 
				tct.author_username = ANY($1)
				AND tct.tweet_created_at >= CURRENT_DATE - INTERVAL '%s days'
				AND tct.tweet_created_at <= CURRENT_DATE
				AND tcs.sentiment IS NOT NULL
				AND tcs.sentiment != 'NaN'
				AND tcs.ticker != 'NONE'
			ORDER BY tcs.tweet_id, tcs.sentiment
		)
		SELECT 
			CASE 
				WHEN sentiment = 'Neutral' THEN 'neutral'
				WHEN sentiment = 'Bullish' THEN 'positive'
				WHEN sentiment = 'Bearish' THEN 'negative'
				ELSE sentiment
			END as sentiment,
			ticker,
			STRING_AGG(tweet_id::TEXT, ', ') AS tweet_ids
		FROM unique_tweets
		GROUP BY sentiment, ticker;
	`, dateRange)

	var results []model.TwitterIdResult
	if err := r.db.Select(&results, query, pq.Array(authorList)); err != nil {
		return []model.TickerSentimentTwitterIdResult{}, fmt.Errorf("failed to get tweet IDs: %w", err)
	}

	var tweetIds []model.TickerSentimentTwitterIdResult
	for _, result := range results {
		tweetIds = append(tweetIds, model.TickerSentimentTwitterIdResult{
			Ticker:    result.Ticker,
			Sentiment: result.Sentiment,
			TweetIds:  strings.Split(result.TweetIds, ","),
		})
	}

	return tweetIds, nil
}
