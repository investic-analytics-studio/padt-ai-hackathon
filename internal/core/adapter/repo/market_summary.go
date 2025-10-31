package repo

import (
	"log"
	"time"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"gorm.io/gorm"
)

type marketSummaryRepo struct {
	db *gorm.DB
}

func NewMarketSummaryRepo(db *gorm.DB) port.MarketSummaryRepo {
	return &marketSummaryRepo{db: db}
}

func (s *marketSummaryRepo) GetDailyMarketOverview() ([]model.DailyMarketOverview, error) {
	tx := s.db.Session(&gorm.Session{})
	statement := `
	SELECT Date, Stock_name, market, pct_change, vol_change
	FROM technical_dataframe_EODHD ps
	WHERE Date >= DATE_SUB(CURDATE(), INTERVAL 5 DAY)
	`
	result := make([]model.DailyMarketOverview, 0)
	if err := tx.Raw(statement).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (s *marketSummaryRepo) GetTopGainer() ([]model.TopGainer, error) {
	tx := s.db.Session(&gorm.Session{})
	statement := `
	SELECT Stock_name, pct_change, market, Close 
	FROM top_gainer2_EODHD 
	ORDER BY pct_change DESC
	`

	result := make([]model.TopGainer, 0)
	if err := tx.Raw(statement).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (s *marketSummaryRepo) GetTopLoser() ([]model.TopLoser, error) {
	tx := s.db.Session(&gorm.Session{})
	statement := `
	SELECT Stock_name, pct_change, market, Close 
	FROM top_loser2_EODHD 
	ORDER BY pct_change ASC
	`

	result := make([]model.TopLoser, 0)
	if err := tx.Raw(statement).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil

}

func (s *marketSummaryRepo) GetAdvancerDeclinerDistribution() (*model.AdvancerDeclinerDistribution, error) {
	tx := s.db.Session(&gorm.Session{})
	statement := "SELECT `-15_percent` as down_15_pct, `-10-15_percent` as down_10_15_pct, `-6-10_percent` as down_6_10_pct,`-4-6_percent` as down_4_6_pct, `-2-4_percent` as down_2_4_pct,`-0-2_percent` as down_0_2_pct,`0_percent` as even_0_pct,`+0-2_percent` as up_0_2_pct,`+2-4_percent` as up_2_4_pct,`+4-6_percent` as up_4_6_pct,`+6-10_percent` as up_6_10_pct,`+10-15_percent` as up_10_15_pct,`+15_percent` as up_15_pct FROM breadth_dataframe_EODHD ORDER BY Date DESC LIMIT 1;"

	result := model.AdvancerDeclinerDistribution{}
	if err := tx.Raw(statement).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *marketSummaryRepo) GetTopTurnoverFloat() ([]model.TurnoverFloat, error) {
	tx := r.db.Session(&gorm.Session{})
	statement := `
		SELECT 
			t1.Date as date, 
			t1.Stock_name as stock_name, 
			t1.TurnoverFloat as turnover_float,  
			t1.Open as open, 
			t1.High as high, 
			t1.Low as low, 
			t1.Close as close, 
			t1.market as market, 
			t1.MarketCap as market_cap
		FROM fundamental_dataframe_EODHD t1
		INNER JOIN (
			SELECT Stock_name, MAX(Date) AS Last_Date
			FROM fundamental_dataframe_EODHD
			GROUP BY Stock_name
		) t2 ON t1.Stock_name = t2.Stock_name AND t1.Date = t2.Last_Date WHERE t1.TurnoverFloat IS NOT NULL ORDER BY t1.TurnoverFloat DESC LIMIT 20;
	`

	turnoverFloats := make([]model.TurnoverFloat, 0)
	if err := tx.Raw(statement).Scan(&turnoverFloats).Error; err != nil {
		return nil, err
	}

	return turnoverFloats, nil
}

func (r *marketSummaryRepo) GetCMEGoldOI() ([]model.CMEGoldOI, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.CMEGoldOI
	if err := tx.Table("cme_gold_oi").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *marketSummaryRepo) GetStockAlertPlots() ([]model.StockAlertPlot, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.StockAlertPlot
	err := tx.Raw("SELECT average_return_sd2_minus,average_return_sd2_plus,average_return_sd1_minus,average_return_sd1_plus,average_return,symbol,detection,return_date FROM tbl_stock_alert_plots").Scan(&result).Error
	return result, err
}

func (r *marketSummaryRepo) GetStockAlertStatsDates() ([]time.Time, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []time.Time
	err := tx.Raw("SELECT DISTINCT date FROM tbl_stock_alert_stats").Scan(&result).Error
	return result, err
}

func (r *marketSummaryRepo) GetStockAlertStats() ([]model.StockAlertStats, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.StockAlertStats
	err := tx.Table("tbl_stock_alert_stats").Find(&result).Error
	return result, err
}

func (r *marketSummaryRepo) GetStockAlertDetections() ([]model.StockAlertDetection, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.StockAlertDetection
	query := `
		SELECT * FROM tbl_stock_alert_detections 
		WHERE date >= DATE_SUB(CURDATE(), INTERVAL 90 DAY)
		ORDER BY date DESC
	`
	err := tx.Raw(query).Scan(&result).Error
	if err != nil {
		log.Printf("Error in GetStockAlertDetections repo: %v", err)
		return nil, err
	}
	log.Printf("GetStockAlertDetections repo result count: %d", len(result))
	return result, nil
}

func (r *marketSummaryRepo) GetCryptoAlertPlots() ([]model.CryptoAlertPlot, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.CryptoAlertPlot
	err := tx.Table("tbl_crypto_alert_plots").Find(&result).Error
	return result, err
}

func (r *marketSummaryRepo) GetCryptoAlertStatsDates() ([]time.Time, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []time.Time
	err := tx.Raw("SELECT DISTINCT date FROM tbl_crypto_alert_stats").Scan(&result).Error
	return result, err
}

func (r *marketSummaryRepo) GetCryptoAlertStats(page int, limit int) ([]model.CryptoAlertStats, int64, error) {
	offset := page * limit
	var total int64
	r.db.Table("tbl_crypto_alert_stats").Where("date >= DATE_SUB(CURDATE(), INTERVAL 5 DAY)").Count(&total)

	query := `
		SELECT * FROM tbl_crypto_alert_stats 
		WHERE date >= DATE_SUB(CURDATE(), INTERVAL 5 DAY)
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`
	var result []model.CryptoAlertStats
	err := r.db.Raw(query, limit, offset).Scan(&result).Error
	if err != nil {
		log.Printf("Error in GetCryptoAlertStats repo: %v", err)
		return nil, 0, err
	}
	log.Printf("GetCryptoAlertStats repo result count: %d", len(result))
	return result, total, nil
}

func (r *marketSummaryRepo) GetCryptoAlertDetections() ([]model.CryptoAlertDetection, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.CryptoAlertDetection
	query := `
		SELECT * FROM tbl_crypto_alert_detections 
		WHERE date >= DATE_SUB(CURDATE(), INTERVAL 5 DAY)
		ORDER BY date DESC
	`
	err := tx.Raw(query).Scan(&result).Error
	if err != nil {
		log.Printf("Error in GetCryptoAlertDetections repo: %v", err)
		return nil, err
	}
	log.Printf("GetCryptoAlertDetections repo result count: %d", len(result))
	return result, nil
}

func (r *marketSummaryRepo) GetCryptoAlertPlots1D() ([]model.CryptoAlertPlot1D, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.CryptoAlertPlot1D
	err := tx.Table("tbl_crypto_alert_plots_1d").Find(&result).Error
	return result, err
}

func (r *marketSummaryRepo) GetCryptoAlertStatsDates1D() ([]time.Time, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []time.Time
	err := tx.Raw("SELECT DISTINCT date FROM tbl_crypto_alert_stats_1d").Scan(&result).Error
	return result, err
}

func (r *marketSummaryRepo) GetCryptoAlertStats1D(page int, limit int) ([]model.CryptoAlertStats1D, int64, error) {
	offset := page * limit
	var total int64
	r.db.Table("tbl_crypto_alert_stats_1d").Where("date >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH)").Count(&total)

	query := `
		SELECT * FROM tbl_crypto_alert_stats_1d 
		WHERE date >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH)
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`
	var result []model.CryptoAlertStats1D
	err := r.db.Raw(query, limit, offset).Scan(&result).Error
	if err != nil {
		log.Printf("Error in GetCryptoAlertStats1D repo: %v", err)
		return nil, 0, err
	}
	log.Printf("GetCryptoAlertStats1D repo result count: %d", len(result))
	return result, total, nil
}

func (r *marketSummaryRepo) GetCryptoAlertDetections1D() ([]model.CryptoAlertDetection1D, error) {
	tx := r.db.Session(&gorm.Session{})
	var result []model.CryptoAlertDetection1D
	query := `
		SELECT * FROM tbl_crypto_alert_detections_1d 
		WHERE date >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH)
		ORDER BY date DESC
	`
	err := tx.Raw(query).Scan(&result).Error
	if err != nil {
		log.Printf("Error in GetCryptoAlertDetections1D repo: %v", err)
		return nil, err
	}
	log.Printf("GetCryptoAlertDetections1D repo result count: %d", len(result))
	return result, nil
}

func (r *marketSummaryRepo) GetAllCryptoAlertStats() ([]model.CryptoAlertStats, error) {
	query := `
		SELECT * FROM tbl_crypto_alert_stats 
		WHERE date >= DATE_SUB(CURDATE(), INTERVAL 5 DAY)
		ORDER BY date DESC
	`
	var result []model.CryptoAlertStats
	err := r.db.Raw(query).Scan(&result).Error
	if err != nil {
		log.Printf("Error in GetAllCryptoAlertStats repo: %v", err)
		return nil, err
	}
	log.Printf("GetAllCryptoAlertStats repo result count: %d", len(result))
	return result, nil
}

func (r *marketSummaryRepo) GetAllCryptoAlertStats1D() ([]model.CryptoAlertStats1D, error) {
	query := `
		SELECT * FROM tbl_crypto_alert_stats_1d 
		WHERE date >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH)
		ORDER BY date DESC
	`
	var result []model.CryptoAlertStats1D
	err := r.db.Raw(query).Scan(&result).Error
	if err != nil {
		log.Printf("Error in GetAllCryptoAlertStats1D repo: %v", err)
		return nil, err
	}
	log.Printf("GetAllCryptoAlertStats1D repo result count: %d", len(result))
	return result, nil
}
