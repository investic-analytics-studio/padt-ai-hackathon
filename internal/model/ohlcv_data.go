package model

import (
	"time"
)

type OHLCVData struct {
	Time   time.Time `json:"time" gorm:"column:bucket;primaryKey"`
	Ticker string    `json:"ticker" gorm:"column:ticker;primaryKey"`
	Open   float64   `json:"open" gorm:"column:open"`
	High   float64   `json:"high" gorm:"column:high"`
	Low    float64   `json:"low" gorm:"column:low"`
	Close  float64   `json:"close" gorm:"column:close"`
	Volume float64   `json:"volume" gorm:"column:volume"`
}

type OHLCVRequest struct {
	Ticker    string     `json:"ticker"`
	TimeFrame string     `json:"tf"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	AllPair   bool       `json:"all_pair"` // If true, return data for all symbols; if false, only return data for the specified ticker
}
