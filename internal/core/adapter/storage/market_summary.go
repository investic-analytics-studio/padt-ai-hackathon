package storage

import (
	"bytes"
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/parquet-go/parquet-go"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

// TODO: turn to sql query
type marketSummaryStorage struct {
	storage *storage.Client
}

func NewMarketSummaryStorage(storage *storage.Client) port.MarketSummaryRepo {
	return &marketSummaryStorage{storage: storage}
}

func (s *marketSummaryStorage) GetDailyMarketOverview() ([]model.DailyMarketOverview, error) {
	file, err := s.storage.Bucket("ds-caching-bucket").Object("pg1_daily_overview.parquet").NewReader(context.TODO())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	readerAt := bytes.NewReader(buf)
	data := parquet.NewReader(readerAt)
	defer data.Close()

	result := make([]model.DailyMarketOverview, 0)
	for {
		var row model.DailyMarketOverview
		if err := data.Read(&row); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func (s *marketSummaryStorage) GetTopGainer() ([]model.TopGainer, error) {
	file, err := s.storage.Bucket("ds-caching-bucket").Object("pg1_top_gainer.parquet").NewReader(context.TODO())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	readerAt := bytes.NewReader(buf)
	data := parquet.NewReader(readerAt)
	defer data.Close()

	result := make([]model.TopGainer, 0)
	for {
		var row model.TopGainer
		if err := data.Read(&row); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func (s *marketSummaryStorage) GetTopLoser() ([]model.TopLoser, error) {
	file, err := s.storage.Bucket("ds-caching-bucket").Object("pg1_top_loser.parquet").NewReader(context.TODO())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	readerAt := bytes.NewReader(buf)
	data := parquet.NewReader(readerAt)
	defer data.Close()

	result := make([]model.TopLoser, 0)
	for {
		var row model.TopLoser
		if err := data.Read(&row); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func (s *marketSummaryStorage) GetAdvancerDeclinerDistribution() (*model.AdvancerDeclinerDistribution, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetTopTurnoverFloat() ([]model.TurnoverFloat, error) { return nil, nil }

func (s *marketSummaryStorage) GetCMEGoldOI() ([]model.CMEGoldOI, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetStockAlertPlots() ([]model.StockAlertPlot, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetStockAlertStatsDates() ([]time.Time, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetStockAlertStats() ([]model.StockAlertStats, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetStockAlertDetections() ([]model.StockAlertDetection, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetCryptoAlertDetections1D() ([]model.CryptoAlertDetection1D, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetCryptoAlertStatsDates1D() ([]time.Time, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetCryptoAlertDetections() ([]model.CryptoAlertDetection, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetCryptoAlertPlots() ([]model.CryptoAlertPlot, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetCryptoAlertPlots1D() ([]model.CryptoAlertPlot1D, error) {
	return nil, nil
}

func (s *marketSummaryStorage) GetCryptoAlertStatsDates() ([]time.Time, error) {
	return nil, nil
}
func (s *marketSummaryStorage) GetCryptoAlertStats1D(page int, limit int) ([]model.CryptoAlertStats1D, int64, error) {
	return nil, 0, nil
}
func (s *marketSummaryStorage) GetCryptoAlertStats(page int, limit int) ([]model.CryptoAlertStats, int64, error) {
	return nil, 0, nil
}
func (s *marketSummaryStorage) GetAllCryptoAlertStats() ([]model.CryptoAlertStats, error) {
	return nil, nil
}
func (s *marketSummaryStorage) GetAllCryptoAlertStats1D() ([]model.CryptoAlertStats1D, error) {
	return nil, nil
}
