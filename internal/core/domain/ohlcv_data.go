package domain

import (
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type OHLCVDataService interface {
	GetCryptoOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error)
	GetForexOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error)
}
