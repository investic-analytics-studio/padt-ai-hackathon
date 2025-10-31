package model

import (
	"time"
)

type CryptoAlertPlot struct {
	ID                    int64     `gorm:"column:id" json:"id"`
	AverageReturn         float64   `gorm:"column:average_return" json:"average_return"`
	AverageReturnSD1Plus  float64   `gorm:"column:average_return_sd1_plus" json:"average_return_sd1_plus"`
	AverageReturnSD1Minus float64   `gorm:"column:average_return_sd1_minus" json:"average_return_sd1_minus"`
	AverageReturnSD2Plus  float64   `gorm:"column:average_return_sd2_plus" json:"average_return_sd2_plus"`
	AverageReturnSD2Minus float64   `gorm:"column:average_return_sd2_minus" json:"average_return_sd2_minus"`
	ReturnDate            string    `gorm:"column:return_date" json:"return_date"`
	Symbol                string    `gorm:"column:symbol" json:"symbol"`
	Detection             string    `gorm:"column:detection" json:"detection"`
	CreatedAt             time.Time `gorm:"column:created_at" json:"created_at"`
}

type CryptoAlertStats struct {
	ID                int64     `gorm:"column:id" json:"id"`
	Symbol            string    `gorm:"column:symbol" json:"symbol"`
	Detection         string    `gorm:"column:detection" json:"detection"`
	TotalTrade        int       `gorm:"column:total_trade" json:"total_trade"`
	WinCount          int       `gorm:"column:win_count" json:"win_count"`
	LossCount         int       `gorm:"column:loss_count" json:"loss_count"`
	WinRate           float64   `gorm:"column:win_rate" json:"win_rate"`
	AvgPositiveReturn float64   `gorm:"column:avg_positive_return" json:"avg_positive_return"`
	AvgNegativeReturn float64   `gorm:"column:avg_negative_return" json:"avg_negative_return"`
	AvgReturn         float64   `gorm:"column:avg_return" json:"avg_return"`
	RRR               float64   `gorm:"column:RRR" json:"rrr"`
	ExpectedReturn    float64   `gorm:"column:expected_return" json:"expected_return"`
	ReturnDate        string    `gorm:"column:return_date" json:"return_date"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	Date              time.Time `gorm:"column:date" json:"date"`
}

type CryptoAlertDetection struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Symbol    string    `gorm:"column:symbol" json:"symbol"`
	Detection string    `gorm:"column:detection" json:"detection"`
	Close     float64   `gorm:"column:close" json:"close"`
	Date      time.Time `gorm:"column:date" json:"date"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

type CryptoAlertStats1D struct {
	ID                int64     `gorm:"column:id" json:"id"`
	Symbol            string    `gorm:"column:symbol" json:"symbol"`
	Detection         string    `gorm:"column:detection" json:"detection"`
	TotalTrade        int       `gorm:"column:total_trade" json:"total_trade"`
	WinCount          int       `gorm:"column:win_count" json:"win_count"`
	LossCount         int       `gorm:"column:loss_count" json:"loss_count"`
	WinRate           float64   `gorm:"column:win_rate" json:"win_rate"`
	AvgPositiveReturn float64   `gorm:"column:avg_positive_return" json:"avg_positive_return"`
	AvgNegativeReturn float64   `gorm:"column:avg_negative_return" json:"avg_negative_return"`
	AvgReturn         float64   `gorm:"column:avg_return" json:"avg_return"`
	RRR               float64   `gorm:"column:RRR" json:"rrr"`
	ExpectedReturn    float64   `gorm:"column:expected_return" json:"expected_return"`
	ReturnDate        string    `gorm:"column:return_date" json:"return_date"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	Date              time.Time `gorm:"column:date" json:"date"`
}

type CryptoAlertPlot1D struct {
	ID                    int64     `gorm:"column:id" json:"id"`
	AverageReturn         float64   `gorm:"column:average_return" json:"average_return"`
	AverageReturnSD1Plus  float64   `gorm:"column:average_return_sd1_plus" json:"average_return_sd1_plus"`
	AverageReturnSD1Minus float64   `gorm:"column:average_return_sd1_minus" json:"average_return_sd1_minus"`
	AverageReturnSD2Plus  float64   `gorm:"column:average_return_sd2_plus" json:"average_return_sd2_plus"`
	AverageReturnSD2Minus float64   `gorm:"column:average_return_sd2_minus" json:"average_return_sd2_minus"`
	ReturnDate            string    `gorm:"column:return_date" json:"return_date"`
	Symbol                string    `gorm:"column:symbol" json:"symbol"`
	Detection             string    `gorm:"column:detection" json:"detection"`
	CreatedAt             time.Time `gorm:"column:created_at" json:"created_at"`
}

type CryptoAlertDetection1D struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Symbol    string    `gorm:"column:symbol" json:"symbol"`
	Detection string    `gorm:"column:detection" json:"detection"`
	Close     float64   `gorm:"column:close" json:"close"`
	Date      time.Time `gorm:"column:date" json:"date"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
