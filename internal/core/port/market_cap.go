package port

import "github.com/quantsmithapp/datastation-backend/internal/model"

type MarketCapRepo interface {
	GetMarketCap() (model.MarketCapResponse, error)
}

type MarketCapService interface {
	GetMarketCap() (model.MarketCapResponse, error)
}
