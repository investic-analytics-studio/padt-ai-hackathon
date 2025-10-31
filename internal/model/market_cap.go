package model

type MarketCapResponse struct {
	AllMarketCap []MarketCapEntities `json:"all_market_cap"`
}

type MarketCapEntities struct {
	BaseAsset             string  `json:"base_asset" db:"base_asset"`
	MarketCap             float64 `json:"market_cap" db:"market_cap"`
	FullyDilutedMarketCap float64 `json:"fully_diluted_market_cap" db:"fully_diluted_market_cap"`
}
