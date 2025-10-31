package port

import (
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type TopRankRepo interface {
	GetTop100() (model.TopRankResponse, error)
}

type TopRankService interface {
	GetTop100() (model.TopRankResponse, error)
}
