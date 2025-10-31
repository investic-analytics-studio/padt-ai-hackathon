package port

import "github.com/quantsmithapp/datastation-backend/internal/model"

type MarketOverviewRepo interface {
	GetMarketOverview(days int, tier string) (model.MarketOverview, error)
	GetMarketOverviewTable(tier string, timeframe string, additionalTiers []string) ([]model.MarketOverviewTable, error)
	GetNewsOverview(timeframe string) ([]model.NewsOverview, error)
	GetTokenDetailMarketOverview(token string, days int) ([]model.TokenDetailMarketOverview, error)
}

type MarketOverviewService interface {
	GetMarketOverview(days int, tier string) (model.MarketOverview, error)
	GetMarketOverviewTable(tier string, timeframe string) ([]model.MarketOverviewTableResponse, error)
	GetEnrichedMarketOverviewTable(tier string, timeframe string, scannerResult map[string]interface{}) ([]model.EnrichedMarketOverview, error)
	GetTokenDetailMarketOverview(token string, days int) ([]model.TokenDetailMarketOverview, error)
}
