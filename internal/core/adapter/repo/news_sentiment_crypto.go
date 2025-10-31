package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type NewsSentimentCryptoRepo struct {
	db *sqlx.DB
}

func NewNewsSentimentCryptoRepo(db *sqlx.DB) port.NewsSentimentCryptoRepo {
	return &NewsSentimentCryptoRepo{db: db}
}

func (r *NewsSentimentCryptoRepo) GetNewsSentiment(ctx context.Context, date time.Time) ([]model.NewsSentimentEntities, error) {
	var data []model.NewsSentimentEntities
	sentimentQuery := `
WITH sentiment_counts AS (
    SELECT 
        ticker,
        date_trunc('hour', date) AS truncated_hour,
        sentiment,
        COUNT(*) as sentiment_count
    FROM (
        SELECT 
            UNNEST(string_to_array(REPLACE(tickers, '$', ''), ',')) AS ticker,
            date,
            LOWER(sentiment) as sentiment
        FROM cryptonewsapi_crypto_news_sentiment
        WHERE date >= NOW() + interval '7 hours' - INTERVAL '30 days'  -- Adjust time range as needed
    ) AS expanded_tickers
    GROUP BY ticker, date_trunc('hour', date), sentiment
),
time_series AS (
    SELECT 
        ticker,
        generate_series(
            date_trunc('hour', NOW() + interval '7 hours' - INTERVAL '30 days'),
            date_trunc('hour', NOW() + interval '7 hours'),
            '24 hour'
        ) AS hour_timestamp
    FROM (
        SELECT DISTINCT ticker 
        FROM sentiment_counts
    ) t
),
rolling_counts AS (
    SELECT 
        ts.ticker,
        ts.hour_timestamp::date as date,
    	hour_timestamp AS original_datetime,
        COALESCE(SUM(CASE WHEN sc.sentiment = 'positive' THEN sc.sentiment_count END), 0) as positive,
        COALESCE(SUM(CASE WHEN sc.sentiment = 'negative' THEN sc.sentiment_count END), 0) as negative,
        COALESCE(SUM(CASE WHEN sc.sentiment = 'neutral' THEN sc.sentiment_count END), 0) as neutral
    FROM time_series ts
    LEFT JOIN sentiment_counts sc 
        ON ts.ticker = sc.ticker 
        AND sc.truncated_hour >= ts.hour_timestamp - INTERVAL '24 hours'
        AND sc.truncated_hour <= ts.hour_timestamp
    GROUP BY ts.ticker, ts.hour_timestamp
),
sentiment_metrics AS (
	SELECT
		p.ticker as coin,
		p.date,
		p.original_datetime,
		p.positive,
		p.negative,
		p.neutral,
		(p.positive + p.negative + p.neutral) AS total_mentions,
		(p.positive - p.negative) AS net_sentiment,
		CAST(p.positive - p.negative AS FLOAT) / NULLIF(CAST(p.positive + p.negative + p.neutral AS FLOAT), 0) AS bullish_ratio
	FROM
		rolling_counts p
),
final_decisions AS (
		SELECT
			sm.*,
			COALESCE(lag.bullish_ratio, 0) AS last_bullish_ratio,
			CASE 
            -- When both are positive, take absolute difference if current is higher
            WHEN sm.bullish_ratio > 0 AND COALESCE(lag.bullish_ratio, 0) > 0 AND sm.bullish_ratio > lag.bullish_ratio THEN
                ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0))
            -- When both are positive but current is lower, take negative difference
            WHEN sm.bullish_ratio > 0 AND COALESCE(lag.bullish_ratio, 0) > 0 AND sm.bullish_ratio <= lag.bullish_ratio THEN
                ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0))
            -- When both are negative and current is more negative, take negative difference
            WHEN sm.bullish_ratio < 0 AND COALESCE(lag.bullish_ratio, 0) < 0 AND sm.bullish_ratio < lag.bullish_ratio THEN
                -(ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)))
            -- When both are negative but current is less negative, take positive difference
            WHEN sm.bullish_ratio < 0 AND COALESCE(lag.bullish_ratio, 0) < 0 AND sm.bullish_ratio >= lag.bullish_ratio THEN
                - ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0))
            -- Default case: regular difference for mixed positive/negative scenarios
            ELSE 
                sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)
        END AS diff_bullish_ratio,
        CASE
            WHEN COUNT(*) OVER (PARTITION BY sm.coin ORDER BY sm.date ROWS BETWEEN 29 PRECEDING AND CURRENT ROW) >= 3
            THEN AVG(
                CASE 
                    WHEN sm.bullish_ratio > 0 AND COALESCE(lag.bullish_ratio, 0) > 0 AND sm.bullish_ratio > lag.bullish_ratio THEN
                        ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0))
                    WHEN sm.bullish_ratio > 0 AND COALESCE(lag.bullish_ratio, 0) > 0 AND sm.bullish_ratio <= lag.bullish_ratio THEN
                        ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0))
                    WHEN sm.bullish_ratio < 0 AND COALESCE(lag.bullish_ratio, 0) < 0 AND sm.bullish_ratio < lag.bullish_ratio THEN
                        -(ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)))
                    WHEN sm.bullish_ratio < 0 AND COALESCE(lag.bullish_ratio, 0) < 0 AND sm.bullish_ratio >= lag.bullish_ratio THEN
                        -(ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)))
                    ELSE 
                        sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)
                END
            ) OVER (PARTITION BY sm.coin ORDER BY sm.date ROWS BETWEEN 29 PRECEDING AND CURRENT ROW)
        END AS rolling_30d_diff_bullish_ratio_mean,
        CASE
            WHEN COUNT(*) OVER (PARTITION BY sm.coin ORDER BY sm.date ROWS BETWEEN 29 PRECEDING AND CURRENT ROW) >= 3
            THEN STDDEV(
                CASE 
                    WHEN sm.bullish_ratio > 0 AND COALESCE(lag.bullish_ratio, 0) > 0 AND sm.bullish_ratio > lag.bullish_ratio THEN
                        ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0))
                    WHEN sm.bullish_ratio > 0 AND COALESCE(lag.bullish_ratio, 0) > 0 AND sm.bullish_ratio <= lag.bullish_ratio THEN
                        ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0))
                    WHEN sm.bullish_ratio < 0 AND COALESCE(lag.bullish_ratio, 0) < 0 AND sm.bullish_ratio < lag.bullish_ratio THEN
                        -(ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)))
                    WHEN sm.bullish_ratio < 0 AND COALESCE(lag.bullish_ratio, 0) < 0 AND sm.bullish_ratio >= lag.bullish_ratio THEN
                        -(ABS(sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)))
                    ELSE 
                        sm.bullish_ratio - COALESCE(lag.bullish_ratio, 0)
                END
            ) OVER (PARTITION BY sm.coin ORDER BY sm.date ROWS BETWEEN 29 PRECEDING AND CURRENT ROW)
        END AS rolling_30d_diff_bullish_ratio_std
		FROM
			sentiment_metrics sm
		LEFT JOIN sentiment_metrics lag
			ON sm.coin  = lag.coin 
			AND sm.date = lag.date + INTERVAL '1 day'
),
	max_dates_per_coin AS (
		SELECT
			coin,
			MAX(date) AS max_date
		FROM
			final_decisions
		GROUP BY
			coin
)
SELECT
		fd.coin,
		fd.date,
		fd.total_mentions,
		fd.net_sentiment,
		fd.positive,
		fd.negative,
		fd.neutral,
		fd.bullish_ratio,
		fd.last_bullish_ratio,
		fd.diff_bullish_ratio,
		GREATEST(fd.rolling_30d_diff_bullish_ratio_mean + fd.rolling_30d_diff_bullish_ratio_std, 0.01) AS bullish_cutoff,
		LEAST(fd.rolling_30d_diff_bullish_ratio_mean - fd.rolling_30d_diff_bullish_ratio_std, -0.01) AS bearish_cutoff,
		CASE
			WHEN fd.diff_bullish_ratio >= GREATEST(fd.rolling_30d_diff_bullish_ratio_mean + fd.rolling_30d_diff_bullish_ratio_std, 0.01) 
				AND fd.diff_bullish_ratio > 0 
			THEN 'Bullish'
			WHEN fd.diff_bullish_ratio <= LEAST(fd.rolling_30d_diff_bullish_ratio_mean - fd.rolling_30d_diff_bullish_ratio_std, -0.01) 
				AND fd.diff_bullish_ratio < 0 
			THEN 'Bearish'
			ELSE 'Holding'
		END AS decision
	FROM
		final_decisions fd
	JOIN
		max_dates_per_coin mdc ON fd.coin = mdc.coin AND fd.date = mdc.max_date
	WHERE
		fd.total_mentions >= 5
		AND fd.original_datetime > NOW() + interval '7 hours' - INTERVAL '1 day'
		AND fd.original_datetime <= NOW() + interval '7 hours';
	`

	if err := r.db.SelectContext(ctx, &data, sentimentQuery); err != nil {
		return nil, fmt.Errorf("failed to get news sentiment: %w", err)
	}

	return data, nil
}
