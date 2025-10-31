package domain

import (
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type MarketSummaryService interface {
	GetDailyMarketOverview() ([]DailyMarketOverviewData, error)
	GetTopGainer() ([]model.TopGainer, error)
	GetTopLoser() ([]model.TopLoser, error)
	GetAdvancerDeclinerDistributionHist() ([]AdvancerDeclinerDistributionHist, error)
	GetAdvancerDeclinerDistributionBar() (*AdvancerDeclinerDistributionBar, error)
	GetTopTurnoverFloat() ([]model.TurnoverFloat, error)
	GetCMEGoldOI() ([]model.CMEGoldOI, error)
	GetStockAlertPlots() ([]model.StockAlertPlot, error)
	GetStockAlertStatsDates() ([]time.Time, error)
	GetStockAlertStats() ([]model.StockAlertStats, error)
	GetStockAlertDetections() ([]model.StockAlertDetection, error)
	GetCryptoAlertPlots() ([]model.CryptoAlertPlot, error)
	GetCryptoAlertStatsDates() ([]time.Time, error)
	GetCryptoAlertDetections() ([]model.CryptoAlertDetection, error)
	GetCryptoAlertPlots1D() ([]model.CryptoAlertPlot1D, error)
	GetCryptoAlertStatsDates1D() ([]time.Time, error)
	GetCryptoAlertDetections1D() ([]model.CryptoAlertDetection1D, error)
	GetCryptoAlertStats(page int, limit int) (*PaginatedResponse, error)
	GetCryptoAlertStats1D(page int, limit int) (*PaginatedResponse, error)
	GetAllCryptoAlertStats() ([]model.CryptoAlertStats, error)
	GetAllCryptoAlertStats1D() ([]model.CryptoAlertStats1D, error)
}

type DailyMarketOverviewData struct {
	Market string                      `json:"market"`
	Data   []model.DailyMarketOverview `json:"data"`
}

type AdvancerDeclinerDistributionHist struct {
	Index int    `json:"index"`
	Count int64  `json:"count"`
	Col   string `json:"col"`
}

type AdvancerDeclinerDistributionBar struct {
	Date    string  `json:"date"`
	Decline float64 `json:"decline"`
	Neutral float64 `json:"neutral"`
	Advanc  float64 `json:"advanc"`
}

type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	Total       int64       `json:"total"`
	HasNextPage bool        `json:"hasNextPage"`
}
