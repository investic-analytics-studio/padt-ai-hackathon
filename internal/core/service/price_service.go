package service

import (
	"slices"
	"strings"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type PriceService struct {
}

func NewPriceService() *PriceService {
	return &PriceService{}
}

func (s *PriceService) GetPriceChange(symbols []string, scannerData []model.ScannerData) map[string]interface{} {
	// Create a map to store price changes
	priceChangeResults := make(map[string]interface{})
	foundList := make([]string, 0)
	// Process each requested symbol
	for _, symbol := range symbols {
		// Look for the symbol in scanner data

		found := false
		for _, item := range scannerData {
			// Extract symbol from CRYPTO:XXUSD format
			scannerSymbol := item.D[0].(string)

			if slices.Contains(foundList, scannerSymbol) {
				continue
			}

			if strings.EqualFold(scannerSymbol, symbol) {
				if priceChange, ok := item.D[3].(float64); ok {
					priceChangeResults[symbol] = priceChange / 100 // Convert percentage to decimal
				} else {
					priceChangeResults[symbol] = nil
				}
				foundList = append(foundList, symbol)
				found = true
				break
			}
		}

		// If symbol not found in scanner data
		if !found {
			priceChangeResults[symbol] = nil
		}
	}

	return priceChangeResults
}
