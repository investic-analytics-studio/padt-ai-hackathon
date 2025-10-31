package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type MarketCapRepo struct {
	db *sqlx.DB
}

func NewMarketCapRepo(db *sqlx.DB) port.MarketCapRepo {
	return &MarketCapRepo{db: db}
}

func (r *MarketCapRepo) GetMarketCap() (model.MarketCapResponse, error) {
	query := `
	SELECT 
		base_asset as "base_asset",
		fully_diluted_market_cap as "fully_diluted_market_cap",
		market_cap as "market_cap"
	FROM (
		SELECT 
			base_asset, 
			last_updated,
			market_cap,
			fully_diluted_market_cap,
			ROW_NUMBER() OVER (PARTITION BY base_asset ORDER BY last_updated DESC) AS row_number
		FROM 
			crypto_binance_defination
	) AS subquery
	WHERE row_number = 1
	ORDER BY market_cap DESC;
	`

	var marketCap []model.MarketCapEntities
	if err := r.db.Select(&marketCap, query); err != nil {
		return model.MarketCapResponse{}, fmt.Errorf("failed to get market cap: %w", err)
	}
	return model.MarketCapResponse{AllMarketCap: marketCap}, nil
}
