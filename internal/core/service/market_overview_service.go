package service

import (
	"strings"
	"sync"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type MarketOverviewService struct {
	marketOverviewRepo port.MarketOverviewRepo
}

func NewMarketOverviewService(marketOverviewRepo port.MarketOverviewRepo) *MarketOverviewService {
	return &MarketOverviewService{
		marketOverviewRepo: marketOverviewRepo,
	}
}

func (s *MarketOverviewService) GetMarketOverview(days int, tier string) (model.MarketOverview, error) {
	return s.marketOverviewRepo.GetMarketOverview(days, tier)
}

func (s *MarketOverviewService) GetMarketOverviewTable(tier string, timeframe string) ([]model.MarketOverviewTableResponse, error) {
	marketOverviewTable, err := s.marketOverviewRepo.GetMarketOverviewTable(tier, timeframe, []string{})
	if err != nil {
		return nil, err
	}

	var marketOverviewTableForS []model.MarketOverviewTable
	// Get tier A data and combine with current tier data
	getMarketOverviewTableForS, err := s.marketOverviewRepo.GetMarketOverviewTable("", timeframe, []string{"S", "SSS"})
	if err != nil {
		return nil, err
	}
	marketOverviewTableForS = getMarketOverviewTableForS

	// Initialize the slice with the correct length
	marketOverviewTableForAWithAlphaSentiment := make([]model.MarketOverviewTableResponse, len(marketOverviewTable))

	// Create a map for quick lookup of tier A sentiment
	tierSSentimentMap := make(map[string]string)
	for _, tableS := range marketOverviewTableForS {
		tierSSentimentMap[tableS.Ticker] = tableS.Sentiment
	}

	// add alpha sentiment to marketOverviewTableForA
	for i, table := range marketOverviewTable {
		marketOverviewTableForAWithAlphaSentiment[i].Ticker = table.Ticker
		marketOverviewTableForAWithAlphaSentiment[i].LongCount = table.LongCount
		marketOverviewTableForAWithAlphaSentiment[i].ShortCount = table.ShortCount
		marketOverviewTableForAWithAlphaSentiment[i].TotalSignals = table.TotalSignals
		marketOverviewTableForAWithAlphaSentiment[i].LongShortRatio = table.LongShortRatio
		marketOverviewTableForAWithAlphaSentiment[i].LongShortString = table.LongShortString
		marketOverviewTableForAWithAlphaSentiment[i].Sentiment = table.Sentiment

		if sentiment, exists := tierSSentimentMap[table.Ticker]; exists {
			marketOverviewTableForAWithAlphaSentiment[i].AlphaSentiment = sentiment
		} else {
			marketOverviewTableForAWithAlphaSentiment[i].AlphaSentiment = "NEUTRAL"
		}
	}

	return marketOverviewTableForAWithAlphaSentiment, nil
}

func (s *MarketOverviewService) GetEnrichedMarketOverviewTable(tier string, timeframe string, scannerResult map[string]interface{}) ([]model.EnrichedMarketOverview, error) {
	marketOverviewTable, err := s.GetMarketOverviewTable(tier, timeframe)
	if err != nil {
		return nil, err
	}

	// Create channels for processing
	enrichedDataChan := make(chan model.EnrichedMarketOverview, len(marketOverviewTable))
	var wg sync.WaitGroup

	// Process each market overview entry
	for _, entry := range marketOverviewTable {
		wg.Add(1)
		go func(entry model.MarketOverviewTableResponse) {
			defer wg.Done()

			enriched := model.EnrichedMarketOverview{
				Ticker:          strings.TrimSuffix(entry.Ticker, "USDT"),
				LongCount:       entry.LongCount,
				ShortCount:      entry.ShortCount,
				TotalSignals:    entry.TotalSignals,
				LongShortRatio:  entry.LongShortRatio,
				LongShortString: entry.LongShortString,
				XSentiment:      entry.Sentiment,
				AlphaSentiment:  entry.AlphaSentiment,
				PriceClose:      nil,
				PriceChange24h:  nil,
				MarketCap:       nil,
				Volume24h:       nil,
			}

			// Find matching scanner data
			if scannerData, ok := scannerResult["data"].([]interface{}); ok {
				for _, item := range scannerData {
					if data, ok := item.(map[string]interface{}); ok {
						if d, ok := data["d"].([]interface{}); ok && len(d) >= 7 {
							symbol, ok := d[0].(string)
							if !ok {
								continue
							}

							if strings.EqualFold(symbol, enriched.Ticker) {
								if priceClose, ok := d[3].(float64); ok {
									enriched.PriceClose = &priceClose
								}

								if priceChange24h, ok := d[4].(float64); ok {
									enriched.PriceChange24h = &priceChange24h
								}

								if marketCap, ok := d[5].(float64); ok {
									enriched.MarketCap = &marketCap
								}

								if volume24h, ok := d[6].(float64); ok {
									enriched.Volume24h = &volume24h
								}
								break
							}
						}
					}
				}
			}

			enrichedDataChan <- enriched
		}(entry)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(enrichedDataChan)
	}()

	// Get news overview
	newsOverview, err := s.marketOverviewRepo.GetNewsOverview(timeframe)
	if err != nil {
		return nil, err
	}

	// Create a map for quick news lookup
	newsMap := make(map[string]model.NewsOverview)
	for _, news := range newsOverview {
		newsMap[news.Coin] = news
	}

	// Collect results and combine with news data
	enrichedData := make([]model.EnrichedMarketOverview, 0, len(marketOverviewTable))
	for data := range enrichedDataChan {
		// Add news overview if available
		if news, exists := newsMap[data.Ticker]; exists {
			newsOverviewRatio := float32(news.PositiveRatio)
			data.NewsOverviewRatio = &newsOverviewRatio
			data.NewsSentiment = &news.Sentiment
		}
		enrichedData = append(enrichedData, data)
	}

	return enrichedData, nil
}

func (s *MarketOverviewService) GetTokenDetailMarketOverview(token string, days int) ([]model.TokenDetailMarketOverview, error) {
	return s.marketOverviewRepo.GetTokenDetailMarketOverview(token, days)
}
