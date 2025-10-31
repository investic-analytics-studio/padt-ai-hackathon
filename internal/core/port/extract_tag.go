package port

import (
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type ExtractTagRepo interface {
	GetAllTags() ([]model.ExtractTagEntities, error)
	GetUniqueTags() ([]string, error)
}

type ExtractTagService interface {
	GetAllTags() (model.ExtractTagResponse, error)
	GetUniqueTags() ([]string, error)
}
