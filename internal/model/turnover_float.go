package model

import "time"

type TurnoverFloat struct {
	Date          *time.Time `gorm:"column:date" json:"date"`
	Market        string     `gorm:"column:market" json:"market"`
	StockName     string     `gorm:"column:stock_name" json:"stock_name"`
	TurnoverFloat float64    `gorm:"column:turnover_float" json:"turnover_float"`
	Open          float64    `gorm:"column:open" json:"open"`
	High          float64    `gorm:"column:high" json:"high"`
	Low           float64    `gorm:"column:low" json:"low"`
	Close         float64    `gorm:"column:close" json:"close"`
	MarketCap     float64    `gorm:"column:market_cap" json:"market_cap"`
}
