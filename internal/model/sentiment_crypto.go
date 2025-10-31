package model

import (
	"time"
)

/*=====	Original SentimentCrypto Model ============ */
//		type SentimentCrypto struct {
//			Title           string         `json:"title" gorm:"column:title;type:text;not null"`
//			Text            sql.NullString `json:"text" gorm:"column:text;type:text"`
//			SourceName      string         `json:"source_name" gorm:"column:source_name;type:varchar(255);not null"`
//			Date            time.Time      `json:"date" gorm:"column:date;type:datetime;not null"`
//			Topics          string         `json:"topics" gorm:"column:topics;type:varchar(255);not null"`
//			Sentiment       string         `json:"sentiment" gorm:"column:sentiment;type:varchar(50);not null"`
//			Type            string         `json:"type" gorm:"column:type;type:varchar(50);not null"`
//			Tickers         string         `json:"tickers" gorm:"column:tickers;type:text;not null"`
//			LastUpdate      string         `json:"last_update" gorm:"column:lastupdate;type:varchar(50);not null"`
//			NewsID          string         `json:"news_id" gorm:"column:news_id;type:varchar(50);not null"`
//			RankScore       float64        `json:"rank_score" gorm:"column:rank_score;type:decimal(3,2);not null"`
//			PositiveCount   int            `json:"positive_count" gorm:"-"`
//			NegativeCount   int            `json:"negative_count" gorm:"-"`
//			NetSentiment    int            `json:"net_sentiment" gorm:"-"`
//			TotalSentiments int            `json:"total_sentiments" gorm:"-"`
//		}

type SentimentCrypto struct {
	Date      time.Time `json:"date" gorm:"column:date;type:datetime;not null"`
	Sentiment string    `json:"sentiment" gorm:"column:sentiment;type:varchar(50);not null"`
}

type CryptoTicker struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// TableName specifies the table name for this model
func (SentimentCrypto) TableName() string {
	return "crypto_sentiment"
}

// You can add more structs here if needed for sentiment crypto functionality
