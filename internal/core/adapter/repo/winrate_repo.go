package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type winRateRepo struct {
	db *sqlx.DB
}

func NewWinRateRepo(db *sqlx.DB) *winRateRepo {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &winRateRepo{db: db}
}

func (r *winRateRepo) GetWinRate() ([]model.WinRateEntities, error) {
	// query := `
	// SELECT author_id, author_username, "crypto_winrate_1D", "crypto_winrate_3D", "crypto_winrate_7D", "crypto_winrate_15D", "crypto_winrate_30D", "total_count_signals"
	// FROM twitter_crypto_backtesting
	// `
	// var winRateResult []model.WinRateEntities
	// if err := r.db.Select(&winRateResult, query); err != nil {
	// 	return nil, fmt.Errorf("failed to get win rate: %w", err)
	// }
	var winRateResult []model.WinRateEntities
	return winRateResult, nil
}
