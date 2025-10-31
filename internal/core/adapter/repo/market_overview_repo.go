package repo

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type MarketOverviewRepo struct {
	db *sqlx.DB
}

func NewMarketOverviewRepo(db *sqlx.DB) port.MarketOverviewRepo {
	return &MarketOverviewRepo{db: db}
}

func (r *MarketOverviewRepo) GetMarketOverviewTable(tier string, timeframe string, additionalTiers []string) ([]model.MarketOverviewTable, error) {
	timeframe = strings.ToUpper(timeframe)

	// Add tier filter if not "all"
	var tierFilter string

	if len(additionalTiers) > 0 {
		tierFilter = fmt.Sprintf("AND  tcap.author_tier IN ('%s')", strings.Join(additionalTiers, "', '"))
	} else if tier != "ALL" {
		tierFilter = fmt.Sprintf("AND  tcap.author_tier = '%s'", tier)
	} else {
		tierFilter = ""
	}

	query := fmt.Sprintf(`
SELECT 
    signal.ticker AS ticker,
    COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) AS long_count,
    COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) AS short_count,
    COUNT(signal.sentiment) AS total_signals, -- Total count of LONG + SHORT signals
    CASE 
        WHEN COUNT(CASE WHEN signal.sentiment IN ('Bullish', 'Bearish') THEN 1 END) = 0 
        THEN NULL 
        ELSE ROUND(
            (COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) - 
             COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END)) * 1.0 / 
            (COUNT(CASE WHEN signal.sentiment IN ('Bullish', 'Bearish') THEN 1 END)), 2
        ) 
    END AS long_short_ratio,
    CONCAT(
        COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END), '/', 
        COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END)
    ) AS long_short_str,
    CASE 
        WHEN COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) > COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) 
        THEN 'LONG'
        WHEN COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) > COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) 
        THEN 'SHORT'
        ELSE 'NEUTRAL'
    END AS sentiment
FROM 
    twitter_crypto_signal AS signal
JOIN 
    twitter_crypto_tweets_foxhole AS tweets 
    ON signal.tweet_id = tweets.id
JOIN 
    twitter_crypto_author_profile tcap 
    ON tweets.author_username = tcap.author_username 
WHERE
    signal.ticker != 'NONE'
    AND signal.ticker IS NOT NULL
    AND signal.sentiment  IN ('Bullish', 'Bearish')
    AND tweets.tweet_created_date >= CURRENT_DATE - INTERVAL '%s days'
    AND tcap.is_select = true
    %s
GROUP BY 
    signal.ticker;
	`, timeframe, tierFilter)

	var marketOverviewTable []model.MarketOverviewTable
	if err := r.db.Select(&marketOverviewTable, query); err != nil {
		return nil, fmt.Errorf("failed to get market overview table: %w", err)
	}
	return marketOverviewTable, nil
}

func (r *MarketOverviewRepo) GetMarketOverview(days int, authorTier string) (model.MarketOverview, error) {
	query := `
		SELECT 
			COUNT(*) as total_tweet,
			COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) as bullish_tweet,
			COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) as bearish_tweet,
			CASE 
				WHEN COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) > COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) THEN 'BULLISH'
				WHEN COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) > COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) THEN 'BEARISH'
				ELSE 'NEUTRAL'
			END as bias,
			CASE 
				WHEN COUNT(CASE WHEN signal.sentiment IN ('Bullish', 'Bearish') THEN 1 END) = 0 THEN 0
				WHEN COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) > COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) THEN 
					ROUND(COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) * 100.0 / 
						(COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) + COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END)))
				WHEN COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) > COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) THEN 
					ROUND(COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END) * 100.0 / 
						(COUNT(CASE WHEN signal.sentiment = 'Bullish' THEN 1 END) + COUNT(CASE WHEN signal.sentiment = 'Bearish' THEN 1 END)))
				ELSE 50
			END as bias_percentage
		FROM twitter_crypto_signal signal
		JOIN twitter_crypto_tweets_foxhole tweets ON signal.tweet_id = tweets.id
		JOIN twitter_crypto_author_profile tcap 
	on tweets.author_username = tcap.author_username 
		WHERE tweets.tweet_created_date >= CURRENT_DATE - INTERVAL '%d days'
		%s
		AND tcap.is_select = true
		AND signal.ticker != 'NONE'
	`

	// Add filters if not "all"
	var filters []string

	// Add author tier filter
	if authorTier != "" && authorTier != "all" {
		filters = append(filters, fmt.Sprintf("AND tcap.author_tier = '%s'", authorTier))
	}

	// Combine all filters
	filterStr := strings.Join(filters, " ")

	// Format the query with days and filters
	query = fmt.Sprintf(query, days, filterStr)

	var result model.MarketOverview
	err := r.db.Get(&result, query)
	if err != nil {
		return model.MarketOverview{}, fmt.Errorf("failed to get market overview: %w", err)
	}

	return result, nil
}

func (r *MarketOverviewRepo) GetNewsOverview(timeframe string) ([]model.NewsOverview, error) {

	query := fmt.Sprintf(`
WITH exploded AS (
    SELECT 
        id,
        TRIM(REPLACE(coin, '$', '')) AS coin, 
        sentiment
    FROM cryptonewsapi_crypto_news_sentiment,
    LATERAL UNNEST(STRING_TO_ARRAY(tickers, ',')) AS coin
    WHERE sentiment IN ('Positive', 'Negative')
    AND date >= CURRENT_DATE - interval '%s days'
),
aggregated AS (
    SELECT 
        coin,
        SUM(CASE WHEN sentiment = 'Positive' THEN 1 ELSE 0 END) AS pos_count,
        SUM(CASE WHEN sentiment = 'Negative' THEN 1 ELSE 0 END) AS neg_count
    FROM exploded
    GROUP BY coin
)
SELECT 
    coin,
    pos_count,
    neg_count,
    pos_count + neg_count as total_counts,
    CASE 
        WHEN (pos_count + neg_count) > 0 
        THEN ROUND(1.0 * (pos_count - neg_count) / (pos_count + neg_count), 4)
        ELSE NULL 
    END AS positive_ratio,
    CASE 
        WHEN (pos_count + neg_count) > 0 
        THEN 
            CASE 
                WHEN ROUND(1.0 * (pos_count - neg_count) / (pos_count + neg_count), 4) > 0 THEN 'LONG'
                WHEN ROUND(1.0 * (pos_count - neg_count) / (pos_count + neg_count), 4) < 0 THEN 'SHORT'
                ELSE 'NEUTRAL'
            END
        ELSE 'NEUTRAL'
    END AS sentiment
FROM aggregated;
`, timeframe)

	var newsOverview []model.NewsOverview
	if err := r.db.Select(&newsOverview, query); err != nil {
		return nil, fmt.Errorf("failed to get news overview: %w", err)
	}

	return newsOverview, nil
}

func (r *MarketOverviewRepo) GetTokenDetailMarketOverview(token string, days int) ([]model.TokenDetailMarketOverview, error) {
	// Try both with and without USDT suffix
	tokenVariations := []string{token}
	if !strings.HasSuffix(strings.ToUpper(token), "USDT") {
		tokenVariations = append(tokenVariations, token+"USDT")
	}

	query := `
	WITH signal_data AS (
		SELECT 
			ts.tweet_id,
			ts.created_at as signal_created_at,
			ts.updated_at as signal_updated_at,
		COALESCE(ts.content, '') as signal_content,
		COALESCE(ts.ticker, '') as signal_ticker,
		COALESCE(ts.action, 'NONE') as signal_action,
		COALESCE(ts.score, 0) as signal_score,
		COALESCE(ts.sentiment, '') as signal_sentiment
		FROM twitter_crypto_signal ts
		WHERE ts.ticker != 'NONE' 
		AND ts.ticker != 'USDCUSDT'
		AND (LOWER(ts.ticker) = LOWER($1) OR LOWER(ts.ticker) = LOWER($2))
	)
	SELECT 
		t.id,
		t.text,
		t.url,
		t.source,
		t.retweetcount,
		t.replycount,
		t.likecount,
		t.quotecount,
		CASE 
			WHEN t.viewcount IS NULL THEN 0 
			WHEN t.viewcount = 'NaN' THEN 0 
			ELSE t.viewcount::float 
		END as viewcount,
		t.tweet_created_at,
		t.bookmarkcount,
		t.isreply,
		t.conversationid,
		t.ispinned,
		t.isretweet,
		t.isquote,
		COALESCE(t.media_url, '') AS media_url,
		t.tweet_created_date,
		t.author_username,
		t.tickers_rule_based,
		t.created_at,
		t.updated_at,
		CASE 
			WHEN s.response_score IS NULL THEN 0 
			ELSE s.response_score 
		END as response_score,
		COALESCE(s.response_sentiment, '') as response_sentiment,
		COALESCE(s.response_tickers, '') as response_tickers,
		COALESCE(a.author_url, '') AS author_url,
		COALESCE(a.author_twitterurl, '') AS author_twitterurl,
		COALESCE(a.author_name, '') AS author_name,
		a.author_followers,
		a.author_following,
		sd.signal_created_at,
		sd.signal_updated_at,
		COALESCE(sd.signal_content, '') AS signal_content,
		COALESCE(sd.signal_ticker, '') AS signal_ticker,
		COALESCE(sd.signal_action, 'NONE') AS signal_action,
		COALESCE(sd.signal_score, 0) AS signal_score,
		COALESCE(sd.signal_sentiment, '') AS signal_sentiment,
		t.id as tweet_id
	FROM signal_data sd
	JOIN twitter_crypto_tweets_foxhole t ON sd.tweet_id = t.id
	LEFT JOIN twitter_crypto_sentiment s ON t.id = s.id
	LEFT JOIN twitter_crypto_author_profile a ON t.author_username = a.author_username
	WHERE t.tweet_created_date >= CURRENT_DATE - INTERVAL '%d days'
	ORDER BY t.tweet_created_at DESC
	`
	// Format the query with the days parameter
	formattedQuery := fmt.Sprintf(query, days)

	// fmt.Printf("Token detail query %s params: %s, %s", formattedQuery, tokenVariations[0], tokenVariations[len(tokenVariations)-1])

	var results []model.TokenDetailMarketOverview

	// Try with both token variations
	err := r.db.Select(&results, formattedQuery, tokenVariations[0], tokenVariations[len(tokenVariations)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to get token detail market overview: %w", err)
	}

	// If no results, try a more lenient query without the date filter
	if len(results) == 0 {
		lenientQuery := strings.Replace(formattedQuery,
			"WHERE t.tweet_created_date >= CURRENT_DATE - INTERVAL '%d days'",
			"", 1)

		err = r.db.Select(&results, lenientQuery, tokenVariations[0], tokenVariations[len(tokenVariations)-1])
		if err != nil {
			return nil, fmt.Errorf("failed to get token detail with lenient query: %w", err)
		}

		// If still no results, log this for debugging
		if len(results) == 0 {
			fmt.Printf("No results found for token %s or %s with any date range\n",
				tokenVariations[0], tokenVariations[len(tokenVariations)-1])
		}
	}

	return results, nil
}
