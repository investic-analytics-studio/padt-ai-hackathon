package model

import "time"

type NewsSentimentEntities struct {
	Ticker           string    `json:"ticker" db:"coin"`
	Date             time.Time `json:"date" db:"date"`
	Negative         float64   `json:"negative" db:"negative"`
	Neutral          float64   `json:"neutral" db:"neutral"`
	Positive         float64   `json:"positive" db:"positive"`
	TotalMentions    float64   `json:"total_mentions" db:"total_mentions"`
	NetSentiment     string    `json:"net_sentiment" db:"net_sentiment"`
	BullishRatio     float64   `json:"bullish_ratio" db:"bullish_ratio"`
	LastBullishRatio float64   `json:"last_bullish_ratio" db:"last_bullish_ratio"`
	DiffBullishRatio float64   `json:"diff_bullish_ratio" db:"diff_bullish_ratio"`
	Bullishcutoff    float64   `json:"bullish_cutoff" db:"bullish_cutoff"`
	BearishCutoff    float64   `json:"bearish_cutoff" db:"bearish_cutoff"`
	Decision         string    `json:"decision" db:"decision"`
}
