package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthorTierRepo struct {
	db *sqlx.DB
}

func NewAuthorTierRepo(db *sqlx.DB) *AuthorTierRepo {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &AuthorTierRepo{db: db}
}

func (r *AuthorTierRepo) GetAuthorsByTier(tier string) ([]string, error) {
	query := `
	SELECT author_username
	FROM twitter_crypto_author_profile
	WHERE author_tier = $1
	AND is_select = true
	`
	var authorList []string
	if err := r.db.Select(&authorList, query, tier); err != nil {
		return nil, fmt.Errorf("failed to get authors with tier: %w", err)
	}
	return authorList, nil
}

func (r *AuthorTierRepo) GetAllTiers() ([]string, error) {
	allTierQuery := `
	SELECT DISTINCT author_tier
	FROM twitter_crypto_author_profile
	WHERE is_select = true
	`
	var tierList []string
	if err := r.db.Select(&tierList, allTierQuery); err != nil {
		return nil, fmt.Errorf("failed to get all tiers: %w", err)
	}
	return tierList, nil
}

func (r *AuthorTierRepo) GetAuthorNameMap() (map[string]string, error) {
	query := `
    SELECT author_username, author_name
    FROM twitter_crypto_author_profile
    WHERE is_select = true
    `
	var rows []struct {
		AuthorUsername string `db:"author_username"`
		AuthorName     string `db:"author_name"`
	}

	if err := r.db.Select(&rows, query); err != nil {
		return nil, fmt.Errorf("failed to get author name map: %w", err)
	}

	authorMap := make(map[string]string)
	for _, row := range rows {
		authorMap[row.AuthorUsername] = row.AuthorName
	}

	return authorMap, nil
}
