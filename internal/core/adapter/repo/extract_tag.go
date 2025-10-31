package repo

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type ExtractTagRepo struct {
	db *sqlx.DB
}

func NewExtractTagRepo(db *sqlx.DB) *ExtractTagRepo {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &ExtractTagRepo{db: db}
}

func (r *ExtractTagRepo) GetAllTags() ([]model.ExtractTagEntities, error) {
	//=========================== temp old category ==============//
	// query := `
	//     SELECT
	//         base_asset,
	//         COALESCE(
	//             CASE
	//                 WHEN perpetual_binance_tag NOT IN ('[]', '{}') THEN perpetual_binance_tag
	//                 ELSE binance_spot_tags
	//             END, ''
	//         ) as binance_tag
	//     FROM crypto_binance_defination`
	query := `
		SELECT 
            base_asset,
            base_asset_tag as binance_tag
        FROM crypto_binance_category cbc
	`
	var tags []model.ExtractTagEntities
	if err := r.db.Select(&tags, query); err != nil {
		return nil, fmt.Errorf("failed to get all tags: %w", err)
	}
	return tags, nil
}

func (r *ExtractTagRepo) GetUniqueTags() ([]string, error) {
	query := `
        SELECT DISTINCT
            COALESCE(
                CASE 
                    WHEN perpetual_binance_tag NOT IN ('[]', '{}') THEN perpetual_binance_tag
                    ELSE binance_spot_tags
                END, ''
            ) as binance_tag
        FROM crypto_binance_defination`

	var rawTags []string
	if err := r.db.Select(&rawTags, query); err != nil {
		return nil, fmt.Errorf("failed to get unique tags: %w", err)
	}

	var tags []string
	for _, rawTag := range rawTags {
		// Remove curly braces and split by comma
		trimmedTag := strings.Trim(rawTag, "{}")
		if trimmedTag != "" {
			splitTags := strings.Split(trimmedTag, ",")
			tags = append(tags, splitTags...)
		}
	}

	return tags, nil
}
