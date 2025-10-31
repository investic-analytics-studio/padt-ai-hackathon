package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type FeaturesRepo struct {
	db *sqlx.DB
}

func NewFeaturesOverviewRepo(db *sqlx.DB) *FeaturesRepo {
	return &FeaturesRepo{db: db}
}
func (r *FeaturesRepo) GetFeature(featureName string) (model.FeaturesEntities, error) {
	query := `
		SELECT 
			feature_name,
			feature_enable
		FROM crypto_features_toggle
		WHERE feature_name = $1
	`

	var featuresEntity model.FeaturesEntities
	if err := r.db.Get(&featuresEntity, query, featureName); err != nil {
		return model.FeaturesEntities{}, fmt.Errorf("failed to get %s feature: %w", featureName, err)
	}
	return featuresEntity, nil
}
