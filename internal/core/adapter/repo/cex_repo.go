package repo

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type CexRepo struct {
	db *sqlx.DB
}

func NewCexRepo(db *sqlx.DB) *CexRepo {
	return &CexRepo{db: db}
}

func (r *CexRepo) WalletExistsByAddress(ctx context.Context, uid string, walletAddress string, exchange string) (bool, error) {
	query := `
		SELECT 1
		FROM crypto_copytrade_wallet_cex
		WHERE crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $1)
		AND wallet_address = $2 
		AND exchange = $3
	`
	var exists int
	if err := r.db.QueryRowContext(ctx, query, uid, walletAddress, exchange).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.Errorf("failed to check cex wallet existence: %v", err)
		return false, err
	}
	return true, nil
}

func (r *CexRepo) InsertCexWallet(ctx context.Context, uid string, record model.CexWalletRecord) (string, error) {
	query := `
		INSERT INTO crypto_copytrade_wallet_cex (
			crypto_user_id,
			wallet_address,
			api_key,
			api_secret,
			priority,
			exchange,
			execution_fee,
			leverage,
			position_size_percentage,
			wallet_name,
			sl_percentage,
			holding_hour_period
		)
		VALUES (
			(SELECT id FROM crypto_user WHERE uuid = $1),
			$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
		RETURNING id
	`

	var walletName interface{}
	if record.WalletName != nil {
		walletName = *record.WalletName
	} else {
		walletName = nil
	}

	var holdingPeriod interface{}
	if record.HoldingHourPeriod != nil {
		holdingPeriod = *record.HoldingHourPeriod
	} else {
		holdingPeriod = nil
	}

	var walletID string
	err := r.db.QueryRowContext(
		ctx,
		query,
		uid,
		record.WalletAddress,
		record.APIKey,
		record.APISecret,
		record.Priority,
		record.Exchange,
		record.ExecutionFee,
		record.Leverage,
		record.PositionSizePercentage,
		walletName,
		record.SlPercentage,
		holdingPeriod,
	).Scan(&walletID)
	if err != nil {
		logger.Errorf("failed to insert cex wallet: %v", err)
		return "", err
	}
	return walletID, nil
}

func (r *CexRepo) ListCexWallets(ctx context.Context, uid string, exchange string) ([]model.CexWalletRecord, error) {
	query := `
		SELECT
			id,
			crypto_user_id,
			wallet_address,
			api_key,
			api_secret,
			priority,
			exchange,
			execution_fee,
			leverage,
			position_size_percentage,
			wallet_name,
			tp_percentage,
			sl_percentage,
			holding_hour_period,
			deleted_at,
			created_at,
			updated_at
		FROM crypto_copytrade_wallet_cex
		WHERE crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $1)
		AND exchange = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query, uid, exchange)
	if err != nil {
		logger.Errorf("failed to list CEX wallets: %v", err)
		return nil, err
	}
	defer rows.Close()

	var wallets []model.CexWalletRecord
	for rows.Next() {
		var wallet model.CexWalletRecord
		if err := rows.StructScan(&wallet); err != nil {
			logger.Errorf("failed to scan CEX wallet: %v", err)
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	if err := rows.Err(); err != nil {
		logger.Errorf("rows error listing CEX wallets: %v", err)
		return nil, err
	}
	return wallets, nil
}

func (r *CexRepo) ActiveCexWallet(ctx context.Context, uid, walletID string) error {
	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET deleted_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $2)
		
	`
	result, err := r.db.ExecContext(ctx, query, walletID, uid)
	if err != nil {
		logger.Errorf("failed to activate CEX wallet: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("failed to get rows affected activating CEX wallet: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no wallet activated (not found or not owned by user)")
	}
	return nil
}

func (r *CexRepo) DeactiveCexWallet(ctx context.Context, uid, walletID string) error {
	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $2)
	`
	result, err := r.db.ExecContext(ctx, query, walletID, uid)
	if err != nil {
		logger.Errorf("failed to deactivate CEX wallet: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("failed to get rows affected deactivating CEX wallet: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no wallet deactivated (not found or not owned by user)")
	}
	return nil
}

func (r *CexRepo) UpdateCexWalletPositionSize(ctx context.Context, uid, walletID string, positionSize float64) error {
	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET position_size_percentage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
	`
	result, err := r.db.ExecContext(ctx, query, positionSize, walletID, uid)
	if err != nil {
		logger.Errorf("failed to update CEX position size: %v", err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating CEX position size: %v", err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}

func (r *CexRepo) UpdateCexWalletHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error {
	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET holding_hour_period = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
	`
	result, err := r.db.ExecContext(ctx, query, holdingPeriod, walletID, uid)
	if err != nil {
		logger.Errorf("failed to update CEX holding period: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("failed to get rows affected updating CEX holding period: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}

func (r *CexRepo) GetCexWalletTotalValue(ctx context.Context, uid, walletID string, exchange string) (model.CexWalletTotalValue, error) {
	var totalValue model.CexWalletTotalValue
	cfg := config.GetConfig()
	baseURL := strings.TrimSuffix(cfg.CryptoTradingBot.BaseURL, "/")
	url := fmt.Sprintf("%s/%s/account-info", baseURL, exchange)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Errorf("failed to create HTTP request: %v", err)
		return totalValue, fmt.Errorf("failed to create account info request")
	}
	query := req.URL.Query()
	query.Set("account_id", walletID)
	req.URL.RawQuery = query.Encode()
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-API-Token", cfg.CryptoTradingBot.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to send account info request: %v", err)
		return totalValue, fmt.Errorf("failed to send account info request")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("unexpected response fetching account info for wallet %s - status %d: %s", walletID, resp.StatusCode, string(body))
		return totalValue, fmt.Errorf("external service returned status %d", resp.StatusCode)
	}

	var accountInfo struct {
		Status           string  `json:"status"`
		TotalBalanceUSDT float64 `json:"total_balance_usdt"`
	}
	if err := json.Unmarshal(body, &accountInfo); err != nil {
		logger.Errorf("failed to decode account info response for wallet %s: %v", walletID, err)
		return totalValue, fmt.Errorf("failed to parse account info response")
	}

	if accountInfo.Status != "" && accountInfo.Status != "ok" {
		logger.Errorf("external service returned non-ok status for wallet %s: %s", walletID, accountInfo.Status)
		return totalValue, fmt.Errorf("external service returned status %s", accountInfo.Status)
	}

	totalValue.TotalValue = accountInfo.TotalBalanceUSDT
	return totalValue, nil

}

func (r *CexRepo) UpdateCexWalletLeverage(ctx context.Context, uid, walletID string, leverage int) error {
	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET leverage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
	`
	result, err := r.db.ExecContext(ctx, query, leverage, walletID, uid)
	if err != nil {
		logger.Errorf("failed to update CEX leverage: %v", err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating CEX leverage: %v", err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}
func (r *CexRepo) UpdateCexWalletSL(ctx context.Context, uid, walletID string, sl float64, exchange string) error {

	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET sl_percentage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
	`
	result, err := r.db.ExecContext(ctx, query, sl, walletID, uid)
	if err != nil {
		logger.Errorf("failed to update CEX sl percentage: %v", err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating CEX sl percentage: %v", err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}

	/*============== send request to trigger api service for pending new sl ===========*/
	reqBody := map[string]interface{}{
		"account_id": walletID,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("failed to marshal request body: %v", err)
		return fmt.Errorf("failed to prepare update request")
	}
	cfg := config.GetConfig()
	url := cfg.CryptoTradingBot.BaseURL + "/" + exchange + "/update-sl"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		logger.Errorf("failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create update request")
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-API-Token", cfg.CryptoTradingBot.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to send update request: %v", err)
		return fmt.Errorf("failed to send update request")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Infof("Successfully updated SL for wallet %s", walletID)
		return nil
	case http.StatusNotFound:
		// Wallet not found in external service - log warning but don't fail
		logger.Warnf("Wallet %s not found in external service, but database update was successful: %s", walletID, string(body))
		return nil
	case http.StatusBadRequest:
		logger.Errorf("Bad request to external service for wallet %s: %s", walletID, string(body))
		return fmt.Errorf("invalid request to external service: %s", string(body))
	case http.StatusUnauthorized:
		logger.Errorf("Unauthorized request to external service for wallet %s: %s", walletID, string(body))
		return fmt.Errorf("authentication failed with external service")
	case http.StatusInternalServerError:
		logger.Errorf("External service internal error for wallet %s: %s", walletID, string(body))
		return fmt.Errorf("external service error: %s", string(body))
	default:
		logger.Errorf("Unexpected response from external service for wallet %s - status %d: %s", walletID, resp.StatusCode, string(body))
		return fmt.Errorf("external service returned status %d", resp.StatusCode)
	}
	return nil

}

func (r *CexRepo) UpdateCexWalletTP(ctx context.Context, uid, walletID string, tp float64, exchange string) error {

	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET tp_percentage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
	`
	result, err := r.db.ExecContext(ctx, query, tp, walletID, uid)
	if err != nil {
		logger.Errorf("failed to update CEX sl percentage: %v", err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating CEX tp percentage: %v", err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil

}

func (r *CexRepo) GetSubscribeAuthor(ctx context.Context, walletID string) ([]model.SubscribeAuthor, error) {
	query := `
		SELECT id, author_username
		FROM crypto_copytrade_authors_privy
		WHERE crypto_user_wallet_id_privy = $1
	`
	rows, err := r.db.QueryContext(ctx, query, walletID)
	if err != nil {
		logger.Errorf("failed to query subscribed authors for CEX wallet %s: %v", walletID, err)
		return nil, err
	}
	defer rows.Close()

	var authors []model.SubscribeAuthor
	for rows.Next() {
		var author model.SubscribeAuthor
		if err := rows.Scan(&author.ID, &author.AuthorUsername); err != nil {
			logger.Errorf("failed to scan subscribed author for CEX wallet %s: %v", walletID, err)
			return nil, err
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		logger.Errorf("row iteration error while fetching authors for CEX wallet %s: %v", walletID, err)
		return nil, err
	}

	return authors, nil
}
func (s *CexRepo) UnsubscribeAuthor(ctx context.Context, author string, walletID string) error {
	query := `
		DELETE FROM crypto_copytrade_authors_privy
		WHERE crypto_user_wallet_id_privy = $1
		AND author_username = $2
		RETURNING id
	`
	var subscribeId string
	err := s.db.QueryRowContext(ctx, query, walletID, author).Scan(&subscribeId)
	if err != nil {
		logger.Errorf("failed to delete and get id: %v", err)
		return err
	}
	return nil
}

func (r *CexRepo) UpdateCexWalletAPICredentials(ctx context.Context, uid, walletID string, apiKey string, apiSecret string, exchange string) error {
	query := `
		UPDATE crypto_copytrade_wallet_cex
		SET api_key = $1, api_secret = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $4)
		AND exchange = $5
	`
	result, err := r.db.ExecContext(ctx, query, apiKey, apiSecret, walletID, uid, exchange)
	if err != nil {
		logger.Errorf("failed to update CEX API credentials: %v", err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating CEX API credentials: %v", err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}

func (r *CexRepo) ValidateCexCredentials(ctx context.Context, exchange string, apiKey string, apiSecret string) (bool, error) {
	cfg := config.GetConfig()
	baseURL := strings.TrimSuffix(cfg.CryptoTradingBot.BaseURL, "/")
	url := fmt.Sprintf("%s/%s/connect", baseURL, exchange)

	payload := map[string]string{
		"api_key":    apiKey,
		"api_secret": apiSecret,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("failed to marshal credential payload: %v", err)
		return false, fmt.Errorf("failed to prepare credential request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		logger.Errorf("failed to create credential validation request: %v", err)
		return false, fmt.Errorf("failed to create credential validation request")
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", cfg.CryptoTradingBot.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to send credential validation request: %v", err)
		return false, fmt.Errorf("failed to send credential validation request")
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("credential validation failed for exchange=%s status=%d body=%s", exchange, resp.StatusCode, string(respBody))
		return false, fmt.Errorf("external service returned status %d", resp.StatusCode)
	}

	var ok bool
	if err := json.Unmarshal(respBody, &ok); err == nil {
		return ok, nil
	}

	var wrapped struct {
		Success bool   `json:"success"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal(respBody, &wrapped); err == nil {
		if wrapped.Success || strings.EqualFold(wrapped.Status, "ok") {
			return true, nil
		}
		return false, nil
	}

	logger.Errorf("unexpected credential validation response for exchange=%s: %s", exchange, string(respBody))
	return false, fmt.Errorf("failed to parse credential validation response")
}
func (s *CexRepo) SubscribeAuthor(ctx context.Context, author string, walletID string) (string, error) {
	query := `
		INSERT INTO crypto_copytrade_authors_privy (crypto_user_wallet_id_privy, author_username)
		VALUES ($1, $2)
		RETURNING id
	`
	var subscribeId string
	err := s.db.QueryRowContext(ctx, query, walletID, author).Scan(&subscribeId)
	if err != nil {
		logger.Errorf("failed to insert and get id: %v", err)
		return "", err
	}
	return subscribeId, nil
}
