package port

import "github.com/quantsmithapp/datastation-backend/internal/model"

type WinRateRepo interface {
	GetWinRate() ([]model.WinRateEntities, error)
}

type WinRateService interface {
	GetWinRate() ([]model.WinRateModel, error)
}
