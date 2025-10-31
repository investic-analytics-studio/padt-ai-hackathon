package model

import "time"

type DailyMarketOverview struct {
	Timestamp *time.Time `gorm:"column:Date" parquet:"Date" json:"timestamp"`
	StockName string     `gorm:"column:Stock_name" parquet:"Stock_name" json:"stock_name"`
	Market    string     `gorm:"column:market" parquet:"market" json:"market"`
	PctChange float64    `gorm:"column:pct_change" parquet:"pct_change" json:"pct_change"`
	VolChange float64    `gorm:"column:vol_change" parquet:"vol_change" json:"vol_change"`
}
