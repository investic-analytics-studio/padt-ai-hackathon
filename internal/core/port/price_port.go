package port

import "github.com/quantsmithapp/datastation-backend/internal/model"

type PricePort interface {
	GetPriceChange(symbols []string, scannerData []model.ScannerData) map[string]interface{}
}
