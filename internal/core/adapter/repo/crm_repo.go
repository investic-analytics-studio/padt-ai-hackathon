package repo

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type CRMRepo struct {
	db *sqlx.DB
}

type KolReferDetail struct {
	UID               string    `db:"crypto_user_id"`
	Email             string    `db:"email"`
	Twitter           string    `db:"twitter"`
	DisplayCode       string    `db:"display_code"`
	ReferralUseNumber int       `db:"referral_use_number"`
	CreatedAt         time.Time `db:"created_at"`
}

// UserBasicInfo represents the basic user info for CRM user listing
type UserBasicInfo struct {
	UUID        string `db:"uuid"`
	Email       string `db:"email"`
	TwitterName string `db:"twitter_name"`
}

func NewCRMRepo(db *sqlx.DB) *CRMRepo {
	return &CRMRepo{db: db}
}

func hashMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func (r *CRMRepo) CheckCRMUser(login model.CRMLoginBody) (string, error) {
	hashedPassword := hashMD5(login.Password)
	var uid string

	query := `
		SELECT id
		FROM crypto_crm_user
		WHERE username = $1 AND password = $2
		LIMIT 1
	`

	err := r.db.Get(&uid, query, login.Username, hashedPassword)
	if err != nil {
		return "", err
	}

	return uid, nil
}

func (r *CRMRepo) CheckCRMUserIsExit(ctx context.Context, uid string) (bool, error) {

	query := `
		SELECT count(uuid)
		FROM crypto_user
		WHERE uuid = $1
	`
	var count int
	err := r.db.Get(&count, query, uid)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *CRMRepo) CheckKOLUserIsExit(ctx context.Context, uid string) (bool, error) {
	query := `
		SELECT count(crypto_user_id)
		FROM crypto_user_kol_code
		WHERE crypto_user_id = $1
	`
	var count int
	err := r.db.Get(&count, query, uid)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CRMRepo) ValidateDisplaycode(ctx context.Context, displayCode string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM crypto_user_kol_code 
			WHERE display_code = $1 OR refcode = $1
		) 
		OR EXISTS (
			SELECT 1 FROM crypto_user_refcode 
			WHERE refcode = $1
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, displayCode)
	if err != nil {
		return false, err
	}

	if exists {
		return false, fmt.Errorf("duplicate code detected, use a different one")
	}

	return true, nil
}

func (r *CRMRepo) InsertKolUser(ctx context.Context, kolUser model.KOLUser) (string, error) {
	query := `
		INSERT INTO crypto_user_kol_code (crypto_user_id, display_code, refcode)
		VALUES (:crypto_user_id, :display_code, :refcode)
		RETURNING crypto_user_id
	`

	params := map[string]interface{}{
		"crypto_user_id": kolUser.ID,
		"display_code":   kolUser.DisplayCode,
		"refcode":        kolUser.DisplayCode,
	}

	var insertedID string
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &insertedID, params); err != nil {
		return "", err
	}

	return insertedID, nil
}

func (r *CRMRepo) GetKolReferDetails(ctx context.Context) ([]KolReferDetail, error) {
	query := `
		SELECT 
			COALESCE(kol.crypto_user_id, '') as crypto_user_id,
			COALESCE(u.email, '') as email,
			COALESCE(u.twitter_name, '') as twitter,
			COALESCE(kol.display_code, '') as display_code,  
			COUNT(cur.crypto_ref_user_id) AS referral_use_number,
			kol.created_at
		FROM crypto_user_kol_code kol
		LEFT JOIN crypto_user u
			ON u.uuid = kol.crypto_user_id
		LEFT JOIN crypto_user_refcode cur
			ON u.uuid = cur.crypto_user_id AND  kol.refcode = cur.refcode 
		GROUP BY
			kol.crypto_user_id,
			kol.display_code,
			kol.refcode,
			u.email,
			u.twitter_name,
			kol.created_at
		
	`

	var results []KolReferDetail
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No KOL referral details found")
		return results, nil
	}

	return results, nil
}

// func (r *CRMRepo) GetKolReferDetails(ctx context.Context, page int) ([]KolReferDetail, error) {
// 	query := `
// 		SELECT
// 			COALESCE(kol.crypto_user_id, '') as crypto_user_id,
// 			COALESCE(u.email, '') as email,
// 			COALESCE(u.twitter_name, '') as twitter,
// 			COALESCE(kol.display_code, '') as display_code,
// 			COUNT(cur.crypto_ref_user_id) AS referral_use_number,
// 			kol.created_at
// 		FROM crypto_user_kol_code kol
// 		LEFT JOIN crypto_user u
// 			ON u.uuid = kol.crypto_user_id
// 		LEFT JOIN crypto_user_refcode cur
// 			ON u.uuid = cur.crypto_user_id
// 		GROUP BY
// 			kol.crypto_user_id,
// 			kol.display_code,
// 			kol.refcode,
// 			u.email,
// 			u.twitter_name,
// 			kol.created_at
// 		LIMIT 10 OFFSET ($1 - 1) * 10;
// 	`

// 	var results []KolReferDetail
// 	err := r.db.SelectContext(ctx, &results, query, page)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %w", err)
// 	}

// 	if len(results) == 0 {
// 		fmt.Println("No KOL referral details found for page:", page)
// 		return results, nil
// 	}

// 	return results, nil
// }

// GetAllUsers returns all users with uuid, email, and twitter_uid
func (r *CRMRepo) GetAllUsers(ctx context.Context) ([]UserBasicInfo, error) {
	query := `
		SELECT 
			COALESCE(uuid, ''::text) AS uuid, 
			COALESCE(email, ''::text) AS email, 
			COALESCE(twitter_name, ''::text) AS twitter_name
		FROM crypto_user;
	`
	var users []UserBasicInfo
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

// Check if crypto_user_id exists in database
func (r *CRMRepo) IsKolUserExists(ctx context.Context, cryptoUserID string) (bool, error) {
	query := `SELECT COUNT(1) FROM crypto_user_kol_code WHERE crypto_user_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, query, cryptoUserID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Check if a display_code exists for another user
func (r *CRMRepo) IsDisplayCodeExists(ctx context.Context, displayCode, excludeUserID string) (bool, error) {
	query := `SELECT COUNT(1) FROM crypto_user_kol_code WHERE display_code = $1 AND crypto_user_id != $2`
	var count int
	err := r.db.GetContext(ctx, &count, query, displayCode, excludeUserID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update display_code for a user
func (r *CRMRepo) UpdateDisplayCode(ctx context.Context, cryptoUserID, displayCode string) error {
	query := `UPDATE crypto_user_kol_code SET display_code = $1, updated_at = NOW() WHERE crypto_user_id = $2`
	_, err := r.db.ExecContext(ctx, query, displayCode, cryptoUserID)
	return err
}

func (r *CRMRepo) GetRefferalScore(ctx context.Context) ([]model.UserWithReferralScores, error) {
	query := `
	SELECT
		cu.uuid, cu.email, cu.twitter_name, -- User details
		crs.crypto_user_id, crs.date, crs.direct_points, -- Score details
		crs.indirect_points, crs.total_points,
		crs.created_at, crs.updated_at
	FROM crypto_refferal_score crs
	INNER JOIN crypto_user cu ON cu.uuid = crs.crypto_user_id
	ORDER BY cu.uuid, crs.date; -- Ordering helps if you want scores in a particular sequence
	`
	var flatData []model.FlatReferralData
	err := r.db.SelectContext(ctx, &flatData, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get referral scores: %w", err)
	}

	userScoresMap := make(map[string]*model.UserWithReferralScores)

	for _, item := range flatData {
		if _, ok := userScoresMap[item.UserUUID]; !ok {
			userScoresMap[item.UserUUID] = &model.UserWithReferralScores{
				UserUUID:    item.UserUUID,
				Email:       item.Email,
				TwitterName: item.TwitterName,
				Scores:      make([]model.CryptoReferralScore, 0),
			}
		}
		score := model.CryptoReferralScore{
			Date:           item.Date,
			DirectPoints:   item.DirectPoints,
			IndirectPoints: item.IndirectPoints,
			TotalPoints:    item.TotalPoints,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
		userScoresMap[item.UserUUID].Scores = append(userScoresMap[item.UserUUID].Scores, score)
	}

	// Convert map to slice
	result := make([]model.UserWithReferralScores, 0, len(userScoresMap))
	for _, userScore := range userScoresMap {
		result = append(result, *userScore)
	}

	return result, nil
}
func (r *CRMRepo) GetKolCode(ctx context.Context, uuid string) ([]model.KolCode, error) {
	query := `
	SELECT 
		crypto_user_kol_code.display_code AS display_code, 
		COALESCE(NULLIF(crypto_user.email, ''), crypto_user.twitter_name) AS crypto_ref_user
	FROM crypto_user_kol_code
	LEFT JOIN crypto_user_refcode ON crypto_user_kol_code.refcode = crypto_user_refcode.refcode 
	LEFT JOIN crypto_user ON (
		crypto_user.uuid = crypto_user_refcode.crypto_ref_user_id
	)
	WHERE crypto_user_kol_code.crypto_user_id = $1 
	GROUP BY crypto_user_kol_code.display_code, crypto_ref_user
	`
	var kolcode []model.KolCode
	err := r.db.SelectContext(ctx, &kolcode, query, uuid)
	if err != nil {
		fmt.Println("Failed to get kolcode:", err)
		return nil, fmt.Errorf("failed to get kolcode: %w", err)
	}
	return kolcode, nil
}
func (r *CRMRepo) GetRefCode(ctx context.Context, uuid string) ([]model.RefCode, error) {
	query := `
	SELECT
		crypto_user_refcode.refcode, 
		COALESCE(NULLIF(crypto_user.email, ''), crypto_user.twitter_name) AS crypto_ref_user
	FROM crypto_user_refcode 
	LEFT JOIN crypto_user ON crypto_user.uuid = crypto_user_refcode.crypto_ref_user_id
	WHERE crypto_user_refcode.refcode NOT IN (
			SELECT refcode FROM crypto_user_kol_code
		) AND 
		crypto_user_refcode.crypto_user_id = $1 
	`
	var refcode []model.RefCode
	err := r.db.SelectContext(ctx, &refcode, query, uuid)
	if err != nil {
		fmt.Println("Failed to get refcode:", err)
		return nil, fmt.Errorf("failed to get refcode: %w", err)
	}
	return refcode, nil
}
func (r *CRMRepo) GetUserReferral(ctx context.Context) ([]model.UserWithRefCode, error) {
	query := `
    SELECT
        crypto_user.uuid, crypto_user.email, crypto_user.twitter_name
    FROM crypto_user 
    INNER JOIN crypto_user_refcode ON crypto_user.uuid = crypto_user_refcode.crypto_user_id 
    AND crypto_user_refcode.crypto_ref_user_id IS NOT NULL
    `
	var flatData []model.FlatUserWithRefCode
	err := r.db.SelectContext(ctx, &flatData, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get referral scores: %w", err)
	}

	// 1. รวบรวม UUID ทั้งหมด
	userUUIDs := make([]string, 0, len(flatData))
	for _, item := range flatData {
		userUUIDs = append(userUUIDs, item.UserUUID)
	}

	// 2. Batch query RefCode
	refQuery := `
    SELECT
        crypto_user_refcode.refcode, 
        COALESCE(NULLIF(crypto_user.email, ''), crypto_user.twitter_name) AS crypto_ref_user,
        crypto_user_refcode.crypto_user_id
    FROM crypto_user_refcode 
    LEFT JOIN crypto_user ON crypto_user.uuid = crypto_user_refcode.crypto_ref_user_id
    WHERE crypto_user_refcode.refcode NOT IN (
        SELECT refcode FROM crypto_user_kol_code
    ) AND crypto_user_refcode.crypto_user_id = ANY($1)
    `
	var allRefCodes []struct {
		RefCode       *string `db:"refcode"`
		CryptoRefUser *string `db:"crypto_ref_user"`
		CryptoUserID  string  `db:"crypto_user_id"`
	}
	err = r.db.SelectContext(ctx, &allRefCodes, refQuery, pq.Array(userUUIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to batch get refcodes: %w", err)
	}
	refMap := make(map[string][]model.RefCode)
	for _, rc := range allRefCodes {
		refMap[rc.CryptoUserID] = append(refMap[rc.CryptoUserID], model.RefCode{
			RefCode:       rc.RefCode,
			CryptoRefUser: rc.CryptoRefUser,
		})
	}

	// 3. Batch query KolCode
	kolQuery := `
    SELECT 
        crypto_user_kol_code.display_code AS display_code, 
        COALESCE(NULLIF(crypto_user.email, ''), crypto_user.twitter_name) AS crypto_ref_user,
        crypto_user_kol_code.crypto_user_id
    FROM crypto_user_kol_code
    LEFT JOIN crypto_user_refcode ON crypto_user_kol_code.refcode = crypto_user_refcode.refcode 
    LEFT JOIN crypto_user ON (crypto_user.uuid = crypto_user_refcode.crypto_ref_user_id)
    WHERE crypto_user_kol_code.crypto_user_id = ANY($1)
    GROUP BY crypto_user_kol_code.display_code, crypto_ref_user, crypto_user_kol_code.crypto_user_id
    `
	var allKolCodes []struct {
		DisplayCode   *string `db:"display_code"`
		CryptoRefUser *string `db:"crypto_ref_user"`
		CryptoUserID  string  `db:"crypto_user_id"`
	}
	err = r.db.SelectContext(ctx, &allKolCodes, kolQuery, pq.Array(userUUIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to batch get kolcodes: %w", err)
	}
	kolMap := make(map[string][]model.KolCode)
	for _, kc := range allKolCodes {
		kolMap[kc.CryptoUserID] = append(kolMap[kc.CryptoUserID], model.KolCode{
			KolCode:       kc.DisplayCode,
			CryptoRefUser: kc.CryptoRefUser,
		})
	}

	// 4. ประกอบผลลัพธ์
	userWithRefcode := make(map[string]*model.UserWithRefCode)
	for _, item := range flatData {
		userWithRefcode[item.UserUUID] = &model.UserWithRefCode{
			UserUUID:    item.UserUUID,
			Email:       item.Email,
			TwitterName: item.TwitterName,
			RefCode:     refMap[item.UserUUID],
			KolCode:     kolMap[item.UserUUID],
		}
	}

	result := make([]model.UserWithRefCode, 0, len(userWithRefcode))
	for _, userCode := range userWithRefcode {
		result = append(result, *userCode)
	}
	return result, nil
}

func (r *CRMRepo) GetCryptoUser(ctx context.Context, page int, order string, search string, isCopytradeApproved *bool) ([]model.CryptoUser, error) {
	// sanitize order
	dir := "DESC"
	if order == "asc" || order == "ASC" {
		dir = "ASC"
	}

	// compute offset
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 200

	base := `
    SELECT
        uuid, email, last_update, created_at, updated_at,
        binding_wallet, twitter_uid, twitter_name, id, address,
        hashed_rt, method_login, stripe_customer_id,
        telegram_chat_id, telegram_user_id, is_copytrade_approved, waiting_list_timestamp
    FROM crypto_user`

	where := ""
	args := []interface{}{}
	if search != "" {
		where = "\n    WHERE (twitter_name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')"
		args = append(args, search)
	}

	if isCopytradeApproved != nil {
		if where == "" {
			where = fmt.Sprintf("\n    WHERE is_copytrade_approved = $%d", len(args)+1)
		} else {
			where += fmt.Sprintf(" AND is_copytrade_approved = $%d", len(args)+1)
		}
		args = append(args, *isCopytradeApproved)
	}

	query := fmt.Sprintf("%s%s\n    ORDER BY waiting_list_timestamp %s NULLS LAST\n    LIMIT 200 OFFSET $%d;", base, where, dir, len(args)+1)

	args = append(args, offset)

	var users []model.CryptoUser
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get crypto users: %w", err)
	}
	return users, nil
}

// CountCryptoUser returns total number of users matching the same filters as GetCryptoUser
func (r *CRMRepo) CountCryptoUser(ctx context.Context, search string, isCopytradeApproved *bool) (int, error) {
	base := `
    SELECT COUNT(*)
    FROM crypto_user`

	where := ""
	args := []interface{}{}

	if search != "" {
		where = "\n    WHERE (twitter_name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')"
		args = append(args, search)
	}

	if isCopytradeApproved != nil {
		if where == "" {
			where = fmt.Sprintf("\n    WHERE is_copytrade_approved = $%d", len(args)+1)
		} else {
			where += fmt.Sprintf(" AND is_copytrade_approved = $%d", len(args)+1)
		}
		args = append(args, *isCopytradeApproved)
	}

	query := fmt.Sprintf("%s%s;", base, where)

	var total int
	if err := r.db.GetContext(ctx, &total, query, args...); err != nil {
		return 0, fmt.Errorf("failed to count crypto users: %w", err)
	}
	return total, nil
}

func (r *CRMRepo) UpdateUserApprove(ctx context.Context, uid string, approve bool) error {
	query := `
	UPDATE crypto_user
	SET is_copytrade_approved = $1
	WHERE uuid = $2
	`
	_, err := r.db.ExecContext(ctx, query, approve, uid)
	if err != nil {
		return fmt.Errorf("failed to update user approve: %w", err)
	}
	return nil
}

// GetUserByID returns a crypto_user by its internal id; accepts UUID value for backward compatibility
func (r *CRMRepo) GetUserByID(ctx context.Context, userID string) (*model.CryptoUser, error) {
	query := `
        SELECT 
            uuid, email, last_update, created_at, updated_at, binding_wallet,
            twitter_uid, twitter_name, id, address, hashed_rt, method_login,
            stripe_customer_id, telegram_chat_id, telegram_user_id,
            is_copytrade_approved, waiting_list_timestamp
        FROM crypto_user
        WHERE id::text = $1 OR uuid::text = $1
        LIMIT 1;
    `
	var u model.CryptoUser
	if err := r.db.GetContext(ctx, &u, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}

// GetPrivyWalletsByUserUUID returns wallets (privy copytrade/genesis) for a given user uuid
func (r *CRMRepo) GetPrivyWalletsByUserUUID(ctx context.Context, uuid string) ([]model.WalletInfo, error) {
	query := `
        SELECT 
            id AS wallet_id,
            wallet_address,
            wallet_id AS privy_wallet_id,
            priority,
            COALESCE(wallet_name, '') AS wallet_name,
            COALESCE(wallet_type, '') AS wallet_type,
            COALESCE(tp_percentage, 0) AS tp_percentage,
            COALESCE(sl_percentage, 0) AS sl_percentage,
            COALESCE(holding_hour_period, 0) AS holding_hour_period,
            COALESCE(position_size_percentage, 0) AS position_size_percentage,
            COALESCE(leverage, 1) AS leverage,
            COALESCE(hyperliquid_basecode, false) AS hyperliquid_basecode,
            created_at, updated_at, deleted_at
        FROM crypto_copytrade_wallet_privy
        WHERE crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $1)
        ORDER BY priority ASC, created_at ASC;
    `

	// Scan into a temporary structure that matches WalletInfo except Authors/Balance
	type walletRow struct {
		WalletID               string     `db:"wallet_id"`
		WalletAddress          string     `db:"wallet_address"`
		PrivyWalletID          string     `db:"privy_wallet_id"`
		Priority               int        `db:"priority"`
		WalletName             string     `db:"wallet_name"`
		WalletType             string     `db:"wallet_type"`
		TpPercentage           float64    `db:"tp_percentage"`
		SlPercentage           float64    `db:"sl_percentage"`
		HoldingHourPeriod      int        `db:"holding_hour_period"`
		PositionSizePercentage float64    `db:"position_size_percentage"`
		Leverage               int        `db:"leverage"`
		HyperliquidBasecode    bool       `db:"hyperliquid_basecode"`
		CreatedAt              *time.Time `db:"created_at"`
		UpdatedAt              *time.Time `db:"updated_at"`
		DeletedAt              *time.Time `db:"deleted_at"`
	}

	var rows []walletRow
	if err := r.db.SelectContext(ctx, &rows, query, uuid); err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	wallets := make([]model.WalletInfo, 0, len(rows))
	for _, w := range rows {
		wallets = append(wallets, model.WalletInfo{
			WalletName:             w.WalletName,
			WalletID:               w.WalletID,
			WalletAddress:          w.WalletAddress,
			WalletType:             w.WalletType,
			Balance:                0, // not calculated in CRM view
			PrivyWalletID:          w.PrivyWalletID,
			Priority:               w.Priority,
			Authors:                nil, // filled by GetAuthorsByWalletID
			TpPercentage:           w.TpPercentage,
			SlPercentage:           w.SlPercentage,
			HoldingHourPeriod:      w.HoldingHourPeriod,
			PositionSizePercentage: w.PositionSizePercentage,
			Leverage:               w.Leverage,
			HyperliquidBasecode:    w.HyperliquidBasecode,
			CreatedAt:              w.CreatedAt,
			UpdatedAt:              w.UpdatedAt,
			DeletedAt:              w.DeletedAt,
		})
	}
	return wallets, nil
}

// GetAuthorsByWalletID returns subscriptions for a wallet (by internal wallet id)
func (r *CRMRepo) GetAuthorsByWalletID(ctx context.Context, walletID string) ([]model.SubscribeAuthor, error) {
	query := `
        SELECT id, author_username
        FROM crypto_copytrade_authors_privy
        WHERE crypto_user_wallet_id_privy = $1
        ORDER BY created_at ASC
    `
	var authors []model.SubscribeAuthor
	if err := r.db.SelectContext(ctx, &authors, query, walletID); err != nil {
		return nil, fmt.Errorf("failed to get authors: %w", err)
	}
	return authors, nil
}

// GetTradeLogsByUserUUID returns trade logs for all privy wallets of the user
func (r *CRMRepo) GetTradeLogsByUserUUID(ctx context.Context, uuid string, status *string, order string, limit, offset int) ([]model.TradeLog, error) {
	dir := "DESC"
	if order == "asc" || order == "ASC" {
		dir = "ASC"
	}
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	base := `
        SELECT id, source, account_id, wallet_id, symbol, side,
               base_size, usdc_value, price, leverage, event, status,
               executed_at, created_at
        FROM trade_logs
        WHERE wallet_id IN (
            SELECT wallet_id
            FROM crypto_copytrade_wallet_privy
            WHERE crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $1)
        )`

	args := []interface{}{uuid}
	where := ""
	if status != nil && *status != "" {
		where = " AND status = $2"
		args = append(args, *status)
	}

	qry := fmt.Sprintf("%s%s ORDER BY executed_at %s LIMIT $%d OFFSET $%d", base, where, dir, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	var logs []model.TradeLog
	if err := r.db.SelectContext(ctx, &logs, qry, args...); err != nil {
		return nil, fmt.Errorf("failed to get trade logs: %w", err)
	}
	return logs, nil
}
