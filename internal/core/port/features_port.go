package port

import (
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type FeaturesService interface {
	GetFeature(string) (model.FeaturesEntities, error)
}
type FeaturesRepo interface {
	GetFeature(string) (model.FeaturesEntities, error)
}
