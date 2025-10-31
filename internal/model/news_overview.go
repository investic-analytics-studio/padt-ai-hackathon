package model

type NewsOverview struct {
	Coin          string  `db:"coin"`
	PosCount      int     `db:"pos_count"`
	NegCount      int     `db:"neg_count"`
	TotalCounts   int     `db:"total_counts"`
	Sentiment     string  `db:"sentiment"`
	PositiveRatio float64 `db:"positive_ratio"`
}
