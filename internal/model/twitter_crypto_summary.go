package model

import (
	"time"
)

type TwitterCryptoSummary struct {
	DateTime      time.Time `json:"date_time" db:"date_time"`
	ResponseTopic string    `json:"response_topic" db:"response_topic"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
