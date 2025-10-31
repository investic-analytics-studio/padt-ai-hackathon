package model

type ExtractTagEntities struct {
	BaseAsset  string `json:"base_asset" db:"base_asset"`
	BinanceTag string `json:"binance_tag" db:"binance_tag"`
}

type ExtractTagResponse struct {
	CoinWithTags map[string][]string `json:"coin_with_tags"`
}
