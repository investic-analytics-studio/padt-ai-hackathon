package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type CryptoUserRefcodeRepo struct {
	db *sqlx.DB
}

func NewCryptoUserRefcodeRepo(db *sqlx.DB) *CryptoUserRefcodeRepo {
	return &CryptoUserRefcodeRepo{db: db}
}
func nullToEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func (r *CryptoUserRefcodeRepo) Create(ctx context.Context, cryptoUserID, refcode string) error {
	query := `
		INSERT INTO crypto_user_refcode (crypto_user_id, refcode)
		VALUES ($1, $2)
		RETURNING refcode
	`

	var returnedRefcode string
	err := r.db.GetContext(ctx, &returnedRefcode, query, cryptoUserID, refcode)
	if err != nil {
		return &ErrDatabaseOperation{
			Operation: "create crypto user refcode",
			Err:       err,
		}
	}

	return nil
}

func (r *CryptoUserRefcodeRepo) GetByCryptoUserID(ctx context.Context, cryptoUserID string) ([]*domain.CryptoUserRefcode, error) {
	query := `
		SELECT cur.crypto_user_id, cur.refcode, cur.crypto_ref_user_id, cur.created_at, cur.updated_at
		FROM crypto_user_refcode cur
		WHERE cur.crypto_user_id = $1
		AND NOT EXISTS (
			SELECT 1
			FROM crypto_user_kol_code kukc
			WHERE kukc.refcode = cur.refcode
		);
	`

	var refcodes []*domain.CryptoUserRefcode
	err := r.db.SelectContext(ctx, &refcodes, query, cryptoUserID)
	if err != nil {
		return nil, &ErrDatabaseOperation{
			Operation: "get crypto user refcodes",
			Err:       err,
		}
	}

	return refcodes, nil
}
func (r *CryptoUserRefcodeRepo) GetRefferalScore(ctx context.Context) ([]*model.RefferalScore, error) {
	if r.db == nil {
		return nil, &ErrInvalidOperation{
			Operation: "get referral score",
			Reason:    "db connection is nil",
		}
	}
	type refferalScoreRaw struct {
		// CryptoUserID      string         `db:"crypto_user_id"`
		CryptoUserEmail   sql.NullString `db:"email"`
		CryptoUserTwitter sql.NullString `db:"twitter_name"`
		TotalPoint        string         `db:"total_points"`
		Date              time.Time      `db:"date"`
	}
	excludedEmails := []string{"contact@9catdigital.com", "vithan.m@padt.ai", "testrefprod01@mail.com", "vithan.m@investicstudio.com", "peerapong.n@investicstudio.com", "padtaicrypto@gmail.com"}
	excludedTwitterNames := []string{"PADT_ai"}

	query := `SELECT email, twitter_name, total_points, date
	FROM (
		SELECT 
		    crs.crypto_user_id,
			cu.email, 
			cu.twitter_name, 
			crs.total_points, 
			crs.date,
			ROW_NUMBER() OVER (PARTITION BY crs.crypto_user_id ORDER BY crs.date DESC) AS rn
		FROM crypto_refferal_score crs
		INNER JOIN crypto_user cu ON cu.uuid = crs.crypto_user_id
		WHERE (cu.email NOT IN (?) OR cu.email IS NULL)
		AND (cu.twitter_name NOT IN (?) OR cu.twitter_name IS NULL)
	) AS ranked
	WHERE rn = 1
	ORDER BY total_points DESC
	`

	query, args, err := sqlx.In(query, excludedEmails, excludedTwitterNames)
	if err != nil {
		return nil, &ErrInvalidOperation{
			Operation: "get referral score",
			Reason:    fmt.Sprintf("failed to execute query or scan results: %v", err),
		}
	}
	query = r.db.Rebind(query)
	var rawScores []refferalScoreRaw
	err = r.db.SelectContext(ctx, &rawScores, query, args...)
	if err != nil {
		return nil, &ErrInvalidOperation{
			Operation: "get referral score",
			Reason:    fmt.Sprintf("failed to execute query or scan results: %v", err),
		}
	}
	refScores := make([]*model.RefferalScore, 0, len(rawScores))
	for _, raw := range rawScores {
		refScores = append(refScores, &model.RefferalScore{
			CryptoUserEmail:   nullToEmpty(raw.CryptoUserEmail),
			CryptoUserTwitter: nullToEmpty(raw.CryptoUserTwitter),
			TotalPoint:        raw.TotalPoint,
			Date:              raw.Date,
		})
	}

	return refScores, nil
}

func (r *CryptoUserRefcodeRepo) GetCryptoKolCode(ctx context.Context, userID string) (model.KolUsed, error) {
	query := `
		SELECT display_code, count(cur.crypto_ref_user_id) AS used_num
		FROM crypto_user_kol_code 
		LEFT JOIN crypto_user_refcode cur on cur.refcode = crypto_user_kol_code.refcode
		WHERE crypto_user_kol_code.crypto_user_id = $1
		GROUP BY display_code 
	`

	var kolCode model.KolUsed

	err := r.db.GetContext(ctx, &kolCode, query, userID)
	if err != nil {
		return model.KolUsed{}, &ErrInvalidOperation{
			Operation: "get kolcode",
			Reason:    "kolcode is not available",
		}

	}

	return kolCode, nil

}
func (r *CryptoUserRefcodeRepo) CheckUserIDExists(ctx context.Context, userID string) (bool, error) {
	query := `
		SELECT 
		EXISTS (
			SELECT 1 
			FROM crypto_user_refcode 
			WHERE crypto_user_id = $1 
			OR crypto_ref_user_id = $1
		)
		OR EXISTS (
			SELECT 1 
			FROM crypto_user_kol_code 
			WHERE crypto_user_id = $1
		);

	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, userID)
	if err != nil {
		fmt.Print("Error: " + err.Error())
		return false, &ErrDatabaseOperation{
			Operation: "check user ID existence",
			Err:       err,
		}
	}

	return exists, nil
}

func (r *CryptoUserRefcodeRepo) CheckRefcodeAvailable(ctx context.Context, refcode string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM crypto_user_refcode 
			WHERE UPPER(refcode) = UPPER($1)
			AND (crypto_ref_user_id IS NULL OR crypto_ref_user_id = '')
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, refcode)
	if err != nil {

		return false, &ErrDatabaseOperation{
			Operation: "check refcode availability",
			Err:       err,
		}
	}

	return exists, nil
}

func (r *CryptoUserRefcodeRepo) CheckRefcodeNotDuplicate(ctx context.Context, refcode string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM crypto_user_refcode 
		WHERE refcode = $1
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, refcode)
	if err != nil {
		return false, &ErrDatabaseOperation{
			Operation: "check refcode not duplicate",
			Err:       err,
		}
	}

	// If count is 0, the refcode is available
	return count == 0, nil
}

func (r *CryptoUserRefcodeRepo) UpdateRefcodeUser(ctx context.Context, refcode string, userID string) error {
	query := `
		UPDATE crypto_user_refcode 
		SET crypto_ref_user_id = $1
		WHERE UPPER(refcode) = UPPER($2) 
		AND (crypto_ref_user_id IS NULL OR crypto_ref_user_id = '')
		RETURNING refcode
	`

	var returnedRefcode string
	err := r.db.GetContext(ctx, &returnedRefcode, query, userID, refcode)
	if err != nil {
		return &ErrDatabaseOperation{
			Operation: "update refcode user",
			Err:       err,
		}
	}

	return nil
}
func (r *CryptoUserRefcodeRepo) CheckKolcode(ctx context.Context, refcode string, userID string) (model.KolCodeEntry, error) {
	var affCode model.KolCodeEntry
	query := `
		SELECT crypto_user_id, display_code, refcode 
		FROM crypto_user_kol_code  
		WHERE UPPER(display_code) = UPPER($1)
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &affCode, query, refcode)
	if err != nil {
		return model.KolCodeEntry{}, &ErrInvalidOperation{
			Operation: "use refcode",
			Reason:    "refcode is not available",
		}

	}

	return affCode, nil
}

func (r *CryptoUserRefcodeRepo) InsertKolcode(ctx context.Context, kolUserID, refcode, userID string) error {
	const checkQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM crypto_user_refcode
			WHERE crypto_user_id = $1 AND refcode = $2 AND crypto_ref_user_id = $3
		)
	`
	var exists bool
	err := r.db.GetContext(ctx, &exists, checkQuery, kolUserID, refcode, userID)
	if err != nil {
		return &ErrDatabaseOperation{
			Operation: "check kol code existence",
			Err:       err,
		}
	}

	if !exists {
		const insertQuery = `
			INSERT INTO crypto_user_refcode (crypto_user_id, refcode, crypto_ref_user_id)
			VALUES ($1, $2, $3)
		`
		_, err := r.db.ExecContext(ctx, insertQuery, kolUserID, refcode, userID)
		if err != nil {
			return &ErrDatabaseOperation{
				Operation: "insert kol code",
				Err:       err,
			}
		}
	}

	return nil
}

func (r *CryptoUserRefcodeRepo) CheckAndUpdateRefcode(ctx context.Context, refcode string, userID string) (bool, error) {
	// First check if refcode is available
	isAvailable, err := r.CheckRefcodeAvailable(ctx, refcode)
	if err != nil {
		return false, err
	}

	// If not available, return false with appropriate error
	if !isAvailable {
		return false, &ErrInvalidOperation{
			Operation: "use refcode",
			Reason:    "refcode is not available",
		}
	}

	// Try to update it
	err = r.UpdateRefcodeUser(ctx, refcode, userID)
	if err != nil {
		return false, &ErrDatabaseOperation{
			Operation: "update refcode user",
			Err:       err,
		}
	}

	// Return true only when update is successful
	return true, nil
}

func (r *CryptoUserRefcodeRepo) CheckRefcodeCountByUserID(ctx context.Context, cryptoUserID string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM crypto_user_refcode 
		WHERE crypto_user_id = $1
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, cryptoUserID)
	if err != nil {
		return 0, &ErrDatabaseOperation{
			Operation: "check refcode count by user ID",
			Err:       err,
		}
	}

	return count, nil
}
func (r *CryptoUserRefcodeRepo) CheckXUserIsExit(ctx context.Context, twitterName string) (model.CheckXUser, error) {
	var checkXUser model.CheckXUser
	fmt.Println("twitterName: " + twitterName)
	query := `
		SELECT 
		EXISTS (
			SELECT 1 
			FROM crypto_user 
			WHERE twitter_name = $1 
		) AS is_user_exit,
		EXISTS (
			SELECT 1 
			FROM crypto_user 
			WHERE twitter_name = $1 AND 
			telegram_chat_id IS NOT NULL
	) AS is_user_telegram_exit;
	`

	err := r.db.GetContext(ctx, &checkXUser, query, twitterName)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return model.CheckXUser{}, &ErrDatabaseOperation{
			Operation: "check x user is exit",
			Err:       err,
		}
	}

	return checkXUser, nil
}
func (r *CryptoUserRefcodeRepo) GetRefferalScoreRanking(ctx context.Context, offsetDays int) ([]*model.RefferalScoreRanking, error) {
	excludedEmails := []interface{}{"contact@9catdigital.com", "vithan.m@padt.ai", "testrefprod01@mail.com", "vithan.m@investicstudio.com", "peerapong.n@investicstudio.com", "padtaicrypto@gmail.com"}
	excludedTwitterNames := []interface{}{"PADT_ai"}

	baseQuery := `
		WITH daily_ranks AS (
			SELECT
				crs.date,
				crs.crypto_user_id,
				COALESCE(cu.email, '') as email,
				COALESCE(cu.twitter_name, '') as twitter_name,
				crs.total_points AS total_point,
				RANK() OVER (
					PARTITION BY crs.date
					ORDER BY crs.total_points DESC
				) AS rk
			FROM crypto_refferal_score crs
			JOIN crypto_user cu ON cu.uuid = crs.crypto_user_id
			WHERE (cu.email IS NULL OR cu.email NOT IN (?))
			AND (cu.twitter_name IS NULL OR cu.twitter_name NOT IN (?))
		),
		ranked_with_prev AS (
			SELECT
				date,
				crypto_user_id,
				email,
				twitter_name,
				total_point,
				rk AS rank,
				LAG(rk, ?) OVER (
					PARTITION BY crypto_user_id
					ORDER BY date
				) AS prev_rank
			FROM daily_ranks
		)
		SELECT
			date,
			email,
			twitter_name,
			total_point,
			rank,
			prev_rank - rank AS rank_change
		FROM ranked_with_prev
		WHERE date = (SELECT MAX(date) FROM crypto_refferal_score)
		ORDER BY rank;
	`

	// Use sqlx.In to expand the NOT IN clauses
	query, args, err := sqlx.In(baseQuery, excludedEmails, excludedTwitterNames, offsetDays)
	if err != nil {
		return nil, &ErrDatabaseOperation{
			Operation: "prepare referral score ranking query",
			Err:       err,
		}
	}

	// Rebind for target driver (e.g., Postgres uses $1, MySQL uses ?)
	query = r.db.Rebind(query)

	var rankings []*model.RefferalScoreRanking
	err = r.db.SelectContext(ctx, &rankings, query, args...)
	if err != nil {
		return nil, &ErrDatabaseOperation{
			Operation: "get referral score ranking",
			Err:       err,
		}
	}

	return rankings, nil
}
