package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type TopRankRepo struct {
	db *sqlx.DB
}

func NewTopRankRepo(db *sqlx.DB) *TopRankRepo {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &TopRankRepo{db: db}
}

func (r *TopRankRepo) GetTop100() (model.TopRankResponse, error) {
	query := `
		SELECT 
			base_asset
		FROM (
			SELECT 
				base_asset, 
				last_updated,
				market_cap,
				ROW_NUMBER() OVER (PARTITION BY base_asset ORDER BY last_updated DESC) AS row_number
			FROM 
				crypto_binance_defination
		) AS subquery
		WHERE row_number = 1
		ORDER BY market_cap DESC;
		`

	var topList []model.TopRank
	if err := r.db.Select(&topList, query); err != nil {
		return model.TopRankResponse{}, fmt.Errorf("failed to get top 100: %w", err)
	}

	baseAssetList := make([]string, 0, len(topList))
	for _, top := range topList {
		baseAssetList = append(baseAssetList, top.BaseAsset)
	}

	return model.TopRankResponse{TopList: baseAssetList}, nil
}
