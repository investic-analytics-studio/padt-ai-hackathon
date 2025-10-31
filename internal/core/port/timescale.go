package port

import (
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type TimescaleRepo interface {
	GetCryptoOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error)
	GetForexOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error)
}
