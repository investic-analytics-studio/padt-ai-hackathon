package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/quantsmithapp/datastation-backend/internal/constant"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type TwitterCryptoRepo struct {
	db *sqlx.DB
}

func NewTwitterCryptoRepo(db *sqlx.DB) port.TwitterCryptoRepo {
	return &TwitterCryptoRepo{db: db}
}

func (r *TwitterCryptoRepo) GetAllSentiments(ctx context.Context) ([]model.TwitterCryptoSentiment, error) {
	var sentiments []model.TwitterCryptoSentiment
	query := `
		SELECT id, text, author_username, 
		       CASE WHEN response_score IS NULL THEN 0 ELSE response_score END as response_score, 
		       response_sentiment, response_tickers, created_at, updated_at
		FROM public.twitter_crypto_sentiment
	`
	err := r.db.SelectContext(ctx, &sentiments, query)
	return sentiments, err
}

func (r *TwitterCryptoRepo) GetAllTweets(ctx context.Context) ([]model.TwitterCryptoTweet, error) {
	var tweets []model.TwitterCryptoTweet
	query := `
		SELECT id, url, text, source, retweetcount, replycount, likecount, quotecount, 
		       CASE WHEN viewcount = 'NaN' THEN 0 ELSE viewcount END as viewcount,
		       tweet_created_at, bookmarkcount, isreply, conversationid, ispinned, isretweet, isquote, 
		       COALESCE(media_url, '') as media_url, tweet_created_date, author_username, tickers_rule_based, created_at, updated_at
		FROM public.twitter_crypto_tweets_foxhole
	`
	err := r.db.SelectContext(ctx, &tweets, query)
	return tweets, err
}

func (r *TwitterCryptoRepo) GetAuthorProfiles(ctx context.Context) ([]model.TwitterCryptoAuthorProfile, error) {
	var profiles []model.TwitterCryptoAuthorProfile
	query := `
		SELECT author_username, author_url, author_twitterurl, author_name, 
			   author_followers, author_following, created_at, updated_at
		FROM public.twitter_crypto_author_profile
		WHERE is_select = true
		ORDER BY author_username
	`
	err := r.db.SelectContext(ctx, &profiles, query)
	return profiles, err
}

func (r *TwitterCryptoRepo) GetAuthorWinrate(ctx context.Context, selectedWinratePeriod string) ([]model.TwitterCryptoAuthorWinRate, error) {

	var additionalWhereClause string

	if selectedWinratePeriod != "overall" {
		field, exists := constant.WINRATE_FIELD_NAME[selectedWinratePeriod]
		if !exists {
			err := errors.New("invalid winrate period")
			logger.Error(err)
			return nil, err
		}
		additionalWhereClause = fmt.Sprintf(`twitter_crypto_backtesting.%s`, field)
	} else {
		additionalWhereClause = `twitter_crypto_backtesting.crypto_winrate_1d * 0.35 + 
			twitter_crypto_backtesting.crypto_winrate_3d * 0.3 + 
			twitter_crypto_backtesting.crypto_winrate_7d * 0.2 + 
			twitter_crypto_backtesting.crypto_winrate_15d * 0.1 +
			twitter_crypto_backtesting.crypto_winrate_30d * 0.05`
	}

	var authorWinrateQueryBuilder strings.Builder
	authorWinrateQueryBuilder.WriteString(fmt.Sprintf(`
	WITH RankedAuthors AS (
		SELECT 
			twitter_crypto_author_profile.author_username,
			twitter_crypto_author_profile.author_url,
			twitter_crypto_author_profile.author_twitterurl,
			twitter_crypto_author_profile.author_name,
			twitter_crypto_author_profile.author_followers,
			twitter_crypto_author_profile.author_following,
			twitter_crypto_author_profile.created_at,
			twitter_crypto_author_profile.updated_at,
			%s AS winrate,
			ROW_NUMBER() OVER (ORDER BY %s DESC) AS rank
		FROM twitter_crypto_author_profile
		INNER JOIN twitter_crypto_backtesting
		ON twitter_crypto_backtesting.author_id = twitter_crypto_author_profile.author_id
		WHERE twitter_crypto_author_profile.is_select = true AND twitter_crypto_backtesting.total_count_signals >= 10
	)
	SELECT 
		CASE WHEN rank <= 5 THEN 'anonymouse'::TEXT ELSE author_username END AS author_username,
		CASE WHEN rank <= 5 THEN 'anonymouse'::TEXT ELSE author_url END AS author_url,
		CASE WHEN rank <= 5 THEN 'anonymouse'::TEXT ELSE author_twitterurl END AS author_twitterurl,
		CASE WHEN rank <= 5 THEN 'anonymouse'::TEXT ELSE author_name END AS author_name,
		CASE WHEN rank <= 5 THEN 0 ELSE author_followers END AS author_followers,
    	CASE WHEN rank <= 5 THEN 0 ELSE author_following END AS author_following,
		created_at,
		updated_at,
		winrate
	FROM RankedAuthors
	ORDER BY winrate DESC;
	`, additionalWhereClause, additionalWhereClause))

	var profiles []model.TwitterCryptoAuthorWinRate
	err := r.db.SelectContext(ctx, &profiles, authorWinrateQueryBuilder.String())
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return profiles, nil
}

func (r *TwitterCryptoRepo) GetPaginatedTweets(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweet, int, error) {
	var tweets []model.TwitterCryptoTweet
	var total int

	// Validate and sanitize sort parameters
	allowedSortFields := map[string]bool{
		"tweet_created_at": true,
		"likecount":        true,
		"retweetcount":     true,
		"viewcount":        true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "tweet_created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM public.twitter_crypto_tweets_foxhole`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := fmt.Sprintf(`
		SELECT id, url, text, source, retweetcount, replycount, likecount, quotecount, 
		       CASE WHEN viewcount = 'NaN' THEN 0 ELSE viewcount END as viewcount,
		       tweet_created_at, bookmarkcount, isreply, conversationid, ispinned, 
		       isretweet, isquote, COALESCE(media_url, '') as media_url, tweet_created_date, author_username, 
		       tickers_rule_based, created_at, updated_at
		FROM public.twitter_crypto_tweets_foxhole
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortBy, sortOrder)

	err = r.db.SelectContext(ctx, &tweets, query, limit, start)
	return tweets, total, err
}

func (r *TwitterCryptoRepo) GetPaginatedSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoSentiment, int, error) {
	var sentiments []model.TwitterCryptoSentiment
	var total int

	// Validate and sanitize sort parameters
	allowedSortFields := map[string]bool{
		"created_at":     true,
		"response_score": true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM public.twitter_crypto_sentiment`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := fmt.Sprintf(`
		SELECT id, text, author_username, 
		       CASE WHEN response_score IS NULL THEN 0 ELSE response_score END as response_score, 
		       response_sentiment, response_tickers, created_at, updated_at
		FROM public.twitter_crypto_sentiment
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortBy, sortOrder)

	err = r.db.SelectContext(ctx, &sentiments, query, limit, start)
	return sentiments, total, err
}

func (r *TwitterCryptoRepo) GetTweetsWithSentiments(ctx context.Context, start, limit int, sortBy, sortOrder string) ([]model.TwitterCryptoTweetWithSentiment, int, error) {
	var tweets []model.TwitterCryptoTweetWithSentiment
	var totalTweet int

	// Validate and sanitize sort parameters
	allowedSortFields := map[string]bool{
		"tweet_created_at": true,
		"likecount":        true,
		"retweetcount":     true,
		"viewcount":        true,
		"response_score":   true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "t.tweet_created_at"
	} else if sortBy == "response_score" {
		sortBy = "s.response_score"
	} else {
		sortBy = "t." + sortBy
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Get total count
	countTweetQuery := `
		SELECT COUNT(*) 
		FROM public.twitter_crypto_tweets_foxhole t
		LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
	`
	err := r.db.GetContext(ctx, &totalTweet, countTweetQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data with joined sentiment information
	query := fmt.Sprintf(`
		SELECT 
			t.id, t.url, t.text, t.source, t.retweetcount, t.replycount, 
			t.likecount, t.quotecount, 
			CASE WHEN t.viewcount = 'NaN' THEN 0 ELSE t.viewcount END as viewcount,
			t.tweet_created_at, t.bookmarkcount, t.isreply, t.conversationid, 
			t.ispinned, t.isretweet, t.isquote, COALESCE(t.media_url, '') as media_url, t.tweet_created_date, 
			t.author_username, t.tickers_rule_based, t.created_at, t.updated_at,
			s.response_score,
			s.response_sentiment,
			s.response_tickers
		FROM public.twitter_crypto_tweets_foxhole t
		LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortBy, sortOrder)

	err = r.db.SelectContext(ctx, &tweets, query, limit, start)
	return tweets, totalTweet, err
}

func (r *TwitterCryptoRepo) GetTweetsWithSentimentAuthorSignal(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAuthorAndSignal, int, error) {
	var tweets []model.TwitterCryptoTweetWithSentimentAuthorAndSignal
	var total int

	// Validate and sanitize sort parameters
	allowedSortFields := map[string]bool{
		"tweet_created_at": true,
		"likecount":        true,
		"retweetcount":     true,
		"viewcount":        true,
		"response_score":   true,
		"author_followers": true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "t.tweet_created_at"
	} else if sortBy == "response_score" {
		sortBy = "s.response_score"
	} else if sortBy == "author_followers" {
		sortBy = "a.author_followers"
	} else {
		sortBy = "t." + sortBy
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Build WHERE clause
	var whereClauses []string
	var args []interface{}
	argCount := 1

	if len(authors) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("t.author_username = ANY($%d)", argCount))
		args = append(args, pq.Array(authors))
		argCount++
	}

	if fromDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.tweet_created_at >= $%d", argCount))
		args = append(args, fromDate)
		argCount++
	}

	if toDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.tweet_created_at <= $%d", argCount))
		args = append(args, toDate)
		argCount++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count with filters
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM public.twitter_crypto_tweets_foxhole t
		LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
		LEFT JOIN public.twitter_crypto_author_profile a ON t.author_username = a.author_username
		%s
	`, whereClause)

	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// First, get the tweets without signals
	mainQueryWithTweetSentimentAndSignal := fmt.Sprintf(`
		WITH tweet_data AS (
			SELECT DISTINCT
				t.id, t.url, t.text, t.source, t.retweetcount, t.replycount, 
				t.likecount, t.quotecount, 
				CASE 
					WHEN t.viewcount IS NULL THEN 0 
					WHEN t.viewcount = 'NaN' THEN 0 
					ELSE t.viewcount::float 
				END as viewcount,
				t.tweet_created_at, t.bookmarkcount, t.isreply, t.conversationid, 
				t.ispinned, t.isretweet, t.isquote, COALESCE(t.media_url, '') as media_url, t.tweet_created_date, 
				t.author_username, t.tickers_rule_based, t.created_at, t.updated_at,
				s.response_score,
				COALESCE(s.response_sentiment, '') as response_sentiment,
				COALESCE(s.response_tickers, '') as response_tickers,
				a.author_url,
				a.author_twitterurl,
				a.author_name,
				a.author_followers,
				a.author_following,
				t.id as tweet_id
			FROM public.twitter_crypto_tweets_foxhole t
			LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
			LEFT JOIN public.twitter_crypto_author_profile a ON t.author_username = a.author_username
			%s
			ORDER BY %s %s
			LIMIT $%d OFFSET $%d
		)
		SELECT 
			td.*,
			ts.created_at as signal_created_at,
			ts.updated_at as signal_updated_at,
			ts.content as signal_content,
			ts.ticker as signal_ticker,
			COALESCE(ts.action, 'NONE') as signal_action,
			ts.score as signal_score,
			ts.sentiment as signal_sentiment,
			COALESCE(ts.tweet_id, td.id) as tweet_id
		FROM tweet_data td
		LEFT JOIN public.twitter_crypto_signal ts ON td.id = ts.tweet_id
		WHERE ts.ticker != 'NONE' AND ts.ticker != 'USDCUSDT
		'
	`, whereClause, sortBy, sortOrder, argCount, argCount+1)

	// Add pagination parameters to args
	args = append(args, limit, start)

	// Execute query and scan directly into tweets
	err = r.db.SelectContext(ctx, &tweets, mainQueryWithTweetSentimentAndSignal, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error scanning tweets: %v", err)
	}

	return tweets, total, err
}

func (r *TwitterCryptoRepo) GetTweetsWithSentimentsAndAuthor(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time, searchTokenSymbolValue string) ([]model.TwitterCryptoTweetWithSentimentAndAuthor, int, error) {
	var tweets []model.TwitterCryptoTweetWithSentimentAndAuthor
	var totalTweet int

	// Validate and sanitize sort parameters
	allowedSortFields := map[string]bool{
		"tweet_created_at": true,
		"likecount":        true,
		"retweetcount":     true,
		"viewcount":        true,
		"response_score":   true,
		"author_followers": true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "t.tweet_created_at"
	} else if sortBy == "response_score" {
		sortBy = "s.response_score"
	} else if sortBy == "author_followers" {
		sortBy = "a.author_followers"
	} else {
		sortBy = "t." + sortBy
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Build WHERE clause
	var whereClauses []string
	var args []interface{}
	argCount := 1

	if len(authors) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("t.author_username = ANY($%d)", argCount))
		args = append(args, pq.Array(authors))
		argCount++
	}

	if fromDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.tweet_created_at >= $%d", argCount))
		args = append(args, fromDate)
		argCount++
	}

	if toDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.tweet_created_at <= $%d", argCount))
		args = append(args, toDate)
		argCount++
	}

	if searchTokenSymbolValue != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("s.response_tickers ILIKE $%d", argCount))
		args = append(args, "%"+searchTokenSymbolValue+"%")
		argCount++
	}

	whereClauses = append(whereClauses, "a.is_select = true")

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count with filters
	countTweetQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM public.twitter_crypto_tweets_foxhole t
		LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
		LEFT JOIN public.twitter_crypto_author_profile a ON t.author_username = a.author_username
		%s
	`, whereClause)

	if len(args) > 0 {
		err := r.db.GetContext(ctx, &totalTweet, countTweetQuery, args...)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.db.GetContext(ctx, &totalTweet, countTweetQuery)
		if err != nil {
			return nil, 0, err
		}
	}

	// Get paginated data with all filters
	tweetWithGroupPaginatedQuery := fmt.Sprintf(`
		SELECT 
			t.id, t.url, t.text, t.source, t.retweetcount, t.replycount, 
			t.likecount, t.quotecount, 
            CASE 
                WHEN t.viewcount IS NULL THEN 0 
                WHEN t.viewcount = 'NaN' THEN 0 
                ELSE t.viewcount::float 
            END as viewcount, t.tweet_created_at, t.bookmarkcount, t.isreply, t.conversationid, 
			t.ispinned, t.isretweet, t.isquote, COALESCE(t.media_url, '') as media_url, t.tweet_created_date, 
			t.author_username, t.tickers_rule_based, t.created_at, t.updated_at,
			s.response_score,
            COALESCE(s.response_sentiment, '') as response_sentiment,
            COALESCE(s.response_tickers, '') as response_tickers,
			a.author_url,
			a.author_twitterurl,
			a.author_name,
			a.author_followers,
			a.author_following
		FROM public.twitter_crypto_tweets_foxhole t
		LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
		LEFT JOIN public.twitter_crypto_author_profile a ON t.author_username = a.author_username
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, sortOrder, argCount, argCount+1)

	// Add pagination parameters to args
	args = append(args, limit, start)

	err := r.db.SelectContext(ctx, &tweets, tweetWithGroupPaginatedQuery, args...)
	return tweets, totalTweet, err
}

func (r *TwitterCryptoRepo) GetSummaries(ctx context.Context, start, limit int, sortBy, sortOrder string, fromDate, toDate *time.Time) ([]model.TwitterCryptoSummary, int, error) {
	var summaries []model.TwitterCryptoSummary
	var totalSummary int

	// Validate and sanitize sort parameters
	allowedSortFields := map[string]bool{
		"date_time":  true,
		"created_at": true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "date_time"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Build WHERE clause for date filtering
	var whereClauses []string
	var args []interface{}
	argCount := 1

	if fromDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("date_time >= $%d", argCount))
		args = append(args, fromDate)
		argCount++
	}

	if toDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("date_time <= $%d", argCount))
		args = append(args, toDate)
		argCount++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count with filters
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM public.twitter_crypto_key_topic
		%s`, whereClause)

	if len(args) > 0 {
		err := r.db.GetContext(ctx, &totalSummary, countQuery, args...)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.db.GetContext(ctx, &totalSummary, countQuery)
		if err != nil {
			return nil, 0, err
		}
	}

	// Get paginated data with date filters
	topicQuery := fmt.Sprintf(`
		SELECT date_time, response_topic, created_at, updated_at
		FROM public.twitter_crypto_key_topic
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, sortOrder, argCount, argCount+1)

	// Add pagination parameters to args
	args = append(args, limit, start)

	err := r.db.SelectContext(ctx, &summaries, topicQuery, args...)
	return summaries, totalSummary, err
}

func (r *TwitterCryptoRepo) GetBubbleSentiment(ctx context.Context) (model.BubbleSentimentModel, error) {
	var cryptoSentiment []model.CryptoSentimentEntities

	// First check total count to avoid unnecessary query if no data exists
	totalSentimentCountQuery := `
		SELECT COUNT(*) FROM twitter_crypto_sentiment
		WHERE created_at >= NOW() - INTERVAL '24 hours'
	`

	var totalSentimentCount int
	if err := r.db.GetContext(ctx, &totalSentimentCount, totalSentimentCountQuery); err != nil {
		return model.BubbleSentimentModel{}, fmt.Errorf("failed to get total count: %w", err)
	}

	// If no data exists for last 24 hours, return empty model
	if totalSentimentCount == 0 {
		return model.BubbleSentimentModel{
			CryptoSentiment: []model.CryptoSentimentEntities{},
			TotalCount:      0,
		}, nil
	}

	// return column names ['ticker', 'total_count', 'positive_pct', 'neutral_pct', 'negative_pct']
	sentimentQuery := `
		WITH cleaned AS (
			SELECT 
				response_sentiment,
				LOWER(translate(response_tickers, '$[]''', '')) AS response_tickers_clean
			FROM 
				twitter_crypto_sentiment t 
			WHERE 
				created_at >= NOW() - INTERVAL '24 hours'
		),
		unnested AS (
			SELECT 
				response_sentiment,
				TRIM(ticker) AS ticker
			FROM cleaned,
				unnest(string_to_array(response_tickers_clean, ',')) AS ticker
		),
		agg AS (
			SELECT 
				ticker,
				COUNT(*) AS total_count,
				ROUND((COUNT(*) FILTER (WHERE response_sentiment = 'Positive')::numeric / COUNT(*)) * 100, 2) AS positive_pct,
				ROUND((COUNT(*) FILTER (WHERE response_sentiment = 'Neutral')::numeric / COUNT(*)) * 100, 2) AS neutral_pct,
				ROUND((COUNT(*) FILTER (WHERE response_sentiment = 'Negative')::numeric / COUNT(*)) * 100, 2) AS negative_pct
			FROM unnested
			WHERE ticker IS NOT NULL
			AND ticker <> 'nan'
			GROUP BY ticker
		)
		SELECT 
			ticker,
			total_count,
			positive_pct,
			neutral_pct,
			negative_pct,
			GREATEST(positive_pct, neutral_pct, negative_pct) AS max_pct,
			CASE 
				WHEN positive_pct >= neutral_pct AND positive_pct >= negative_pct THEN 'Positive'
				WHEN neutral_pct >= positive_pct AND neutral_pct >= negative_pct THEN 'Neutral'
				ELSE 'Negative'
			END AS max_sentiment_type
		FROM agg
		ORDER BY ticker;
	`

	if err := r.db.SelectContext(ctx, &cryptoSentiment, sentimentQuery); err != nil {
		return model.BubbleSentimentModel{}, fmt.Errorf("failed to get bubble sentiment graph: %w", err)
	}

	return model.BubbleSentimentModel{
		CryptoSentiment: cryptoSentiment,
		TotalCount:      totalSentimentCount,
	}, nil
}

func (r *TwitterCryptoRepo) SearchTokenMentionSymbolsByAuthors(ctx context.Context, symbol string, createdAt string, id string, limit int, selectedTimeRange string, authors []string) ([]model.TokenMentionSymbolAndAuthor, time.Time, string, int, error) {

	timeRangeMap := map[string]string{
		"7d": "7 days",
		"1m": "30 days",
		"2m": "60 days",
		"3m": "90 days",
		"6m": "180 days",
	}

	querySearchTokenMentionsSymbol := fmt.Sprintf(`
		SELECT 
			DISTINCT REPLACE(tcs.ticker, 'USDT', '') as symbol
		FROM twitter_crypto_signal tcs 
		JOIN twitter_crypto_tweets_foxhole tct 
			ON tcs.tweet_id = tct.id
		JOIN twitter_crypto_author_profile tcap 
			ON tct.author_username = tcap.author_username
		WHERE tcap.author_username = ANY($1)
			AND tct.tweet_created_at >= NOW() - INTERVAL '%s'
			AND tcs.ticker != 'NONE'
    		AND tcs.ticker != ''
    		AND tcs.ticker IS NOT NULL
	`, timeRangeMap[selectedTimeRange])

	var tokens []model.TokenMentionSymbolAndAuthor
	var lastCreatedAt time.Time
	var lastId string

	err := r.db.SelectContext(ctx, &tokens, querySearchTokenMentionsSymbol, pq.Array(authors))
	if err != nil {
		return nil, time.Time{}, "", 0, err
	}
	return tokens, lastCreatedAt, lastId, len(tokens), nil
}

func (r *TwitterCryptoRepo) GetTweetsWithSentimentsAndTier(ctx context.Context, start, limit int, sortBy, sortOrder string, authors []string, fromDate, toDate *time.Time) ([]model.TwitterCryptoTweetWithSentimentAndTier, int, error) {
	var tweets []model.TwitterCryptoTweetWithSentimentAndTier
	var totalTweet int

	// Validate and sanitize sort parameters
	allowedSortFields := map[string]bool{
		"tweet_created_at": true,
		"likecount":        true,
		"retweetcount":     true,
		"viewcount":        true,
		"response_score":   true,
		"author_followers": true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "t.tweet_created_at"
	} else if sortBy == "response_score" {
		sortBy = "s.response_score"
	} else if sortBy == "author_followers" {
		sortBy = "a.author_followers"
	} else {
		sortBy = "t." + sortBy
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Build WHERE clause
	var whereClauses []string
	var args []interface{}
	argCount := 1

	if len(authors) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("t.author_username = ANY($%d)", argCount))
		args = append(args, pq.Array(authors))
		argCount++
	}

	if fromDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.tweet_created_at >= $%d", argCount))
		args = append(args, fromDate)
		argCount++
	}

	if toDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.tweet_created_at <= $%d", argCount))
		args = append(args, toDate)
		argCount++
	}

	whereClauses = append(whereClauses, "a.is_select = true")

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count with filters
	countTweetQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM public.twitter_crypto_tweets_foxhole t
		LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
		LEFT JOIN public.twitter_crypto_author_profile a ON t.author_username = a.author_username
		%s
	`, whereClause)

	if len(args) > 0 {
		err := r.db.GetContext(ctx, &totalTweet, countTweetQuery, args...)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.db.GetContext(ctx, &totalTweet, countTweetQuery)
		if err != nil {
			return nil, 0, err
		}
	}

	// Updated query to include signal information using CTE
	tweetWithGroupPaginatedQuery := fmt.Sprintf(`
		WITH tweet_data AS (
			SELECT 
				t.id,
				t.id as tweet_id,
				s.response_score,
				COALESCE(s.response_sentiment, '') as response_sentiment,
				COALESCE(s.response_tickers, '') as response_tickers,
				t.tweet_created_at,
				t.tweet_created_date
			FROM public.twitter_crypto_tweets_foxhole t
			LEFT JOIN public.twitter_crypto_sentiment s ON t.id = s.id
			LEFT JOIN public.twitter_crypto_author_profile a ON t.author_username = a.author_username
			%s
			ORDER BY %s %s
			LIMIT $%d OFFSET $%d
		)
		SELECT 
			td.response_score,
			td.response_sentiment,
			td.response_tickers,
			td.tweet_created_at,
			td.tweet_created_date,
			ts.created_at as signal_created_at,
			ts.updated_at as signal_updated_at,
			ts.content as signal_content,
			ts.ticker as signal_ticker,
			COALESCE(ts.action, 'NONE') as signal_action,
			ts.score as signal_score,
			ts.sentiment as signal_sentiment
		FROM tweet_data td
		LEFT JOIN public.twitter_crypto_signal ts ON td.tweet_id = ts.tweet_id
		WHERE ts.ticker != 'NONE' AND ts.ticker != 'USDCUSDT'
	`, whereClause, sortBy, sortOrder, argCount, argCount+1)

	// Add pagination parameters to args
	args = append(args, limit, start)

	// Execute query and scan into tweets
	err := r.db.SelectContext(ctx, &tweets, tweetWithGroupPaginatedQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error scanning tweets: %v", err)
	}

	return tweets, totalTweet, err
}

