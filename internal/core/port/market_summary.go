package port

import (
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type MarketSummaryRepo interface {
	GetDailyMarketOverview() ([]model.DailyMarketOverview, error)
	GetTopGainer() ([]model.TopGainer, error)
	GetTopLoser() ([]model.TopLoser, error)
	GetAdvancerDeclinerDistribution() (*model.AdvancerDeclinerDistribution, error)
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
	// Updated method signatures for paginated endpoints
	GetCryptoAlertStats(page int, limit int) ([]model.CryptoAlertStats, int64, error)
	GetCryptoAlertStats1D(page int, limit int) ([]model.CryptoAlertStats1D, int64, error)
	GetAllCryptoAlertStats() ([]model.CryptoAlertStats, error)
	GetAllCryptoAlertStats1D() ([]model.CryptoAlertStats1D, error)
}

// SELECT t1.database_lastupdate, t1.Stock_name, t1.TurnoverFloat,  t1.Open, t1.High, t1.Low, t1.Close, t1.market, t1.MarketCap
// FROM fundamental_dataframe_EODHD t1
//          INNER JOIN (
//     SELECT Stock_name, MAX(Date) AS Last_Date
//     FROM fundamental_dataframe_EODHD
//     GROUP BY Stock_name
// ) t2 ON t1.Stock_name = t2.Stock_name AND t1.Date = t2.Last_Date WHERE t1.TurnoverFloat IS NOT NULL ORDER BY t1.TurnoverFloat DESC;
