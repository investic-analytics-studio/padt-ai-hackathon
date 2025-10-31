package model

type TopLoser struct {
	Close     float64 `gorm:"column:Close" parquet:"Close" json:"close"`
	StockName string  `gorm:"column:Stock_name" parquet:"Stock_name" json:"stock_name"`
	Market    string  `gorm:"column:market" parquet:"market" json:"market"`
	PctChange float64 `gorm:"column:pct_change" parquet:"pct_change" json:"pct_change"`
}
