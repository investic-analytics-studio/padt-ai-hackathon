package model

type WinRateEntities struct {
	AuthorId         string  `db:"author_id"`
	AuthorUsername   string  `db:"author_username"`
	CryptoWinRate1D  float64 `db:"crypto_winrate_1D"`
	CryptoWinRate3D  float64 `db:"crypto_winrate_3D"`
	CryptoWinRate7D  float64 `db:"crypto_winrate_7D"`
	CryptoWinRate15D float64 `db:"crypto_winrate_15D"`
	CryptoWinRate30D float64 `db:"crypto_winrate_30D"`
	TotalCountSignal float64 `db:"total_count_signals"`
}

type WinRateModel struct {
	AuthorId         string  `json:"author_id"`
	AuthorUsername   string  `json:"author_username"`
	CryptoWinRate1D  float64 `json:"crypto_win_rate_1d"`
	CryptoWinRate3D  float64 `json:"crypto_win_rate_3d"`
	CryptoWinRate7D  float64 `json:"crypto_win_rate_7d"`
	CryptoWinRate15D float64 `json:"crypto_win_rate_15d"`
	CryptoWinRate30D float64 `json:"crypto_win_rate_30d"`
}
