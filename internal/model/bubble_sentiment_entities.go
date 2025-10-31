package model

type CryptoSentimentEntities struct {
	Ticker           string  `json:"ticker" db:"ticker"`
	TotalCount       int     `json:"total_count" db:"total_count"`
	PositivePct      float64 `json:"positive_pct" db:"positive_pct"`
	NeutralPct       float64 `json:"neutral_pct" db:"neutral_pct"`
	NegativePct      float64 `json:"negative_pct" db:"negative_pct"`
	MaxPct           float64 `json:"max_pct" db:"max_pct"`
	MaxSentimentType string  `json:"max_sentiment_type" db:"max_sentiment_type"`
}

type BubbleSentimentModel struct {
	CryptoSentiment []CryptoSentimentEntities `json:"crypto_sentiment"`
	TotalCount      int                       `json:"total_count"`
}
