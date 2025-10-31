package model

import (
	"time"
)

type KolCodeEntry struct {
	CryptoUserID string `db:"crypto_user_id"`
	DisplayCode  string `db:"display_code"`
	Refcode      string `db:"refcode"`
}

type KolUsed struct {
	DisplayCode string `db:"display_code"`
	Used        int    `db:"used_num"`
}

type RefferalScore struct {
	// CryptoUserID      string    `db:"crypto_user_id"`
	CryptoUserEmail   string    `db:"email"`
	CryptoUserTwitter string    `db:"twitter_name"`
	TotalPoint        string    `db:"total_points"`
	Date              time.Time `db:"date"`
}

type RefferalScoreRanking struct {
	Date time.Time `db:"date"`
	// CryptoUserID string    `db:"crypto_user_id"`
	TotalPoint        int    `db:"total_point"`
	Rank              int    `db:"rank"`
	RankChange        *int   `db:"rank_change"`
	CryptoUserEmail   string `db:"email"`
	CryptoUserTwitter string `db:"twitter_name"`
}

type CRMUser struct {
	ID       string `db:"id"`
	Username string `db:"username"`
}
type KOLUser struct {
	ID          string `db:"id"`
	DisplayCode string `db:"display_code"`
}
