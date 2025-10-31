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
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type numericField struct {
	value float64
	set   bool
}

func (n *numericField) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*n = numericField{}
		return nil
	}

	if len(trimmed) > 0 && trimmed[0] == '"' {
		unquoted, err := strconv.Unquote(trimmed)
		if err != nil {
			return err
		}
		trimmed = strings.TrimSpace(unquoted)
	}

	val, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return err
	}

	n.value = val
	n.set = true
	return nil
}

func (n numericField) Float64() (float64, bool) {
	return n.value, n.set
}

type DexRepo struct {
	db *sqlx.DB
}

func NewDexRepo(db *sqlx.DB) *DexRepo {
	return &DexRepo{db: db}
}

func (r *DexRepo) WalletExistsByAddress(ctx context.Context, uid string, walletAddress string, exchange string) (bool, error) {
	query := `
		SELECT 1
		FROM crypto_copytrade_wallet_dex
		WHERE crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $1)
		AND wallet_address = $2
		AND exchange = $3
	`
	var exists int
	if err := r.db.QueryRowContext(ctx, query, uid, walletAddress, exchange).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.Errorf("failed to check dex wallet existence for exchange=%s: %v", exchange, err)
		return false, err
	}
	return true, nil
}

func (r *DexRepo) InsertDexWallet(ctx context.Context, uid string, record model.DexWalletRecord) (string, error) {
	query := `
		INSERT INTO crypto_copytrade_wallet_dex (
			crypto_user_id,
			wallet_address,
			api_key,
			private_key,
			trading_account,
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
			$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		RETURNING id
	`

	var walletName interface{}
	if record.WalletName != nil {
		walletName = *record.WalletName
	} else {
		walletName = nil
	}

	var walletID string
	err := r.db.QueryRowContext(
		ctx,
		query,
		uid,
		record.WalletAddress,
		record.APIKey,
		record.PrivateKey,
		record.TradingAccount,
		record.Priority,
		record.Exchange,
		record.ExecutionFee,
		record.Leverage,
		record.PositionSizePercentage,
		walletName,
		record.SlPercentage,
		record.HoldingHourPeriod,
	).Scan(&walletID)
	if err != nil {
		logger.Errorf("failed to insert dex wallet for exchange=%s: %v", record.Exchange, err)
		return "", err
	}
	return walletID, nil
}

func (r *DexRepo) ListDexWallets(ctx context.Context, uid string, exchange string) ([]model.DexWalletRecord, error) {
	query := `
		SELECT
			id,
			crypto_user_id,
			wallet_address,
			api_key,
			private_key,
			trading_account,
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
		FROM crypto_copytrade_wallet_dex
		WHERE crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $1)
		AND exchange = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query, uid, exchange)
	if err != nil {
		logger.Errorf("failed to list dex wallets for exchange=%s: %v", exchange, err)
		return nil, err
	}
	defer rows.Close()

	var wallets []model.DexWalletRecord
	for rows.Next() {
		var wallet model.DexWalletRecord
		if err := rows.StructScan(&wallet); err != nil {
			logger.Errorf("failed to scan dex wallet for exchange=%s: %v", exchange, err)
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	if err := rows.Err(); err != nil {
		logger.Errorf("rows error listing dex wallets for exchange=%s: %v", exchange, err)
		return nil, err
	}
	return wallets, nil
}

func (r *DexRepo) ActiveDexWallet(ctx context.Context, uid, walletID string) error {
	query := `
		UPDATE crypto_copytrade_wallet_dex
		SET deleted_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $2)
	`
	result, err := r.db.ExecContext(ctx, query, walletID, uid)
	if err != nil {
		logger.Errorf("failed to activate dex wallet id=%s: %v", walletID, err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("failed to get rows affected activating dex wallet id=%s: %v", walletID, err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no wallet activated (not found or not owned by user)")
	}
	return nil
}

func (r *DexRepo) DeactiveDexWallet(ctx context.Context, uid, walletID string) error {
	query := `
		UPDATE crypto_copytrade_wallet_dex
		SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $2)
	`
	result, err := r.db.ExecContext(ctx, query, walletID, uid)
	if err != nil {
		logger.Errorf("failed to deactivate dex wallet id=%s: %v", walletID, err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("failed to get rows affected deactivating dex wallet id=%s: %v", walletID, err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no wallet deactivated (not found or not owned by user)")
	}
	return nil
}

func (r *DexRepo) UpdateDexWalletPositionSize(ctx context.Context, uid, walletID string, positionSize float64) error {
	query := `
		UPDATE crypto_copytrade_wallet_dex
		SET position_size_percentage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
	`
	result, err := r.db.ExecContext(ctx, query, positionSize, walletID, uid)
	if err != nil {
		logger.Errorf("failed to update dex position size wallet_id=%s: %v", walletID, err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating dex position size wallet_id=%s: %v", walletID, err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}

func (r *DexRepo) UpdateDexWalletLeverage(ctx context.Context, uid, walletID string, leverage int, exchange string) error {
	query := `
		UPDATE crypto_copytrade_wallet_dex
		SET leverage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
		AND exchange = $4
	`
	result, err := r.db.ExecContext(ctx, query, leverage, walletID, uid, exchange)
	if err != nil {
		logger.Errorf("failed to update dex leverage wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating dex leverage wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}

func (r *DexRepo) UpdateDexWalletSL(ctx context.Context, uid, walletID string, sl float64, exchange string) error {
	reqBody := map[string]interface{}{
		"account_id": walletID,
		"exchange":   exchange,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	query := `
		UPDATE crypto_copytrade_wallet_dex
		SET sl_percentage = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
		AND exchange = $4
	`
	result, err := r.db.ExecContext(ctx, query, sl, walletID, uid, exchange)
	if err != nil {
		logger.Errorf("failed to update dex sl percentage wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating dex sl percentage wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}

	cfg := config.GetConfig()
	url := cfg.CryptoTradingBot.BaseURL + fmt.Sprintf("/%s/update-sl", exchange)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		logger.Errorf("failed to create HTTP request for dex sl update wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return fmt.Errorf("failed to create update request")
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-API-Token", cfg.CryptoTradingBot.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to send dex sl update request wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return fmt.Errorf("failed to send update request")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Infof("Successfully updated TP/SL for wallet %s exchange=%s", walletID, exchange)
		return nil
	case http.StatusNotFound:
		logger.Warnf("Wallet %s (exchange=%s) not found in external service, but database update was successful: %s", walletID, exchange, string(body))
		return nil
	case http.StatusBadRequest:
		logger.Errorf("Bad request to external service for wallet %s exchange=%s: %s", walletID, exchange, string(body))
		return fmt.Errorf("invalid request to external service: %s", string(body))
	case http.StatusUnauthorized:
		logger.Errorf("Unauthorized request to external service for wallet %s exchange=%s: %s", walletID, exchange, string(body))
		return fmt.Errorf("authentication failed with external service")
	case http.StatusInternalServerError:
		logger.Errorf("External service internal error for wallet %s exchange=%s: %s", walletID, exchange, string(body))
		return fmt.Errorf("external service error: %s", string(body))
	default:
		logger.Errorf("Unexpected response from external service for wallet %s exchange=%s - status %d: %s", walletID, exchange, resp.StatusCode, string(body))
		return fmt.Errorf("external service returned status %d", resp.StatusCode)
	}
}

func (r *DexRepo) GetDexWalletTotalValue(ctx context.Context, uid, walletID string, exchange string) (model.DexWalletTotalValue, error) {
	var totalValue model.DexWalletTotalValue

	cfg := config.GetConfig()
	baseURL := strings.TrimSuffix(cfg.CryptoTradingBot.BaseURL, "/")
	url := fmt.Sprintf("%s/%s/account-info", baseURL, exchange)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Errorf("failed to create HTTP request for dex account info wallet_id=%s exchange=%s: %v", walletID, exchange, err)
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
		logger.Errorf("failed to send dex account info request wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return totalValue, fmt.Errorf("failed to send account info request")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("unexpected response fetching dex account info wallet_id=%s exchange=%s status=%d: %s", walletID, exchange, resp.StatusCode, string(body))
		return totalValue, fmt.Errorf("external service returned status %d", resp.StatusCode)
	}

	var accountInfo struct {
		Status           string       `json:"status"`
		TotalBalanceUSDT numericField `json:"total_balance_usdt"`
		TotalValueUSDT   numericField `json:"total_value_usdt"`
		Data             struct {
			TotalBalanceUSDT numericField `json:"total_balance_usdt"`
			TotalValueUSDT   numericField `json:"total_value_usdt"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &accountInfo); err != nil {
		logger.Errorf("failed to decode dex account info response wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return totalValue, fmt.Errorf("failed to parse account info response")
	}

	if accountInfo.Status != "" && accountInfo.Status != "ok" {
		logger.Errorf("external service returned non-ok status for dex wallet %s exchange=%s: %s", walletID, exchange, accountInfo.Status)
		return totalValue, fmt.Errorf("external service returned status %s", accountInfo.Status)
	}

	if val, ok := accountInfo.TotalBalanceUSDT.Float64(); ok {
		totalValue.TotalValue = val
	} else if val, ok := accountInfo.TotalValueUSDT.Float64(); ok {
		totalValue.TotalValue = val
	} else if val, ok := accountInfo.Data.TotalBalanceUSDT.Float64(); ok {
		totalValue.TotalValue = val
	} else if val, ok := accountInfo.Data.TotalValueUSDT.Float64(); ok {
		totalValue.TotalValue = val
	}
	return totalValue, nil
}

func (r *DexRepo) ValidateDexCredentials(ctx context.Context, exchange, apiKey, privateKey, tradingAccountID string) (bool, error) {
	cfg := config.GetConfig()
	baseURL := strings.TrimSuffix(cfg.CryptoTradingBot.BaseURL, "/")
	url := fmt.Sprintf("%s/%s/connect", baseURL, exchange)

	payload := map[string]string{
		"api_key":            apiKey,
		"private_key":        privateKey,
		"trading_account_id": tradingAccountID,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("failed to marshal dex credential payload exchange=%s: %v", exchange, err)
		return false, fmt.Errorf("failed to prepare credential request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		logger.Errorf("failed to create dex credential validation request exchange=%s: %v", exchange, err)
		return false, fmt.Errorf("failed to create credential validation request")
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", cfg.CryptoTradingBot.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to send dex credential validation request exchange=%s: %v", exchange, err)
		return false, fmt.Errorf("failed to send credential validation request")
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
			logger.Warnf("dex credential validation rejected exchange=%s status=%d body=%s", exchange, resp.StatusCode, string(respBody))
			return false, nil
		}
		logger.Errorf("dex credential validation failed exchange=%s status=%d body=%s", exchange, resp.StatusCode, string(respBody))
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

	logger.Errorf("unexpected dex credential validation response exchange=%s: %s", exchange, string(respBody))
	return false, fmt.Errorf("failed to parse credential validation response")
}

func (r *DexRepo) UpdateDexWalletAPICredentials(ctx context.Context, uid, walletID, apiKey, privateKey, tradingAccountID, exchange string) error {
	query := `
		UPDATE crypto_copytrade_wallet_dex
		SET api_key = $1,
		    private_key = $2,
		    trading_account = $3,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $5)
		AND exchange = $6
	`
	result, err := r.db.ExecContext(ctx, query, apiKey, privateKey, tradingAccountID, walletID, uid, exchange)
	if err != nil {
		logger.Errorf("failed to update dex api credentials wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("failed to get rows affected updating dex api credentials wallet_id=%s exchange=%s: %v", walletID, exchange, err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}

func (r *DexRepo) UpdateDexWalletHoldingPeriod(ctx context.Context, uid, walletID string, holdingPeriod int) error {
	query := `
		UPDATE crypto_copytrade_wallet_dex
		SET holding_hour_period = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND crypto_user_id = (SELECT id FROM crypto_user WHERE uuid = $3)
	`
	result, err := r.db.ExecContext(ctx, query, holdingPeriod, walletID, uid)
	if err != nil {
		logger.Errorf("failed to update dex holding period wallet_id=%s: %v", walletID, err)
		return err
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		logger.Errorf("failed to get rows affected updating dex holding period wallet_id=%s: %v", walletID, err)
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no wallet updated (not found or not owned by user)")
	}
	return nil
}

func (r *DexRepo) GetSubscribeAuthor(ctx context.Context, walletID string) ([]model.SubscribeAuthor, error) {
	query := `
		SELECT id, author_username
		FROM crypto_copytrade_authors_privy
		WHERE crypto_user_wallet_id_privy = $1
	`
	rows, err := r.db.QueryContext(ctx, query, walletID)
	if err != nil {
		logger.Errorf("failed to query subscribed authors for dex wallet %s: %v", walletID, err)
		return nil, err
	}
	defer rows.Close()

	var authors []model.SubscribeAuthor
	for rows.Next() {
		var author model.SubscribeAuthor
		if err := rows.Scan(&author.ID, &author.AuthorUsername); err != nil {
			logger.Errorf("failed to scan subscribed author for dex wallet %s: %v", walletID, err)
			return nil, err
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		logger.Errorf("row iteration error while fetching authors for dex wallet %s: %v", walletID, err)
		return nil, err
	}

	return authors, nil
}

func (r *DexRepo) SubscribeAuthor(ctx context.Context, author string, walletID string) (string, error) {
	query := `
		INSERT INTO crypto_copytrade_authors_privy (crypto_user_wallet_id_privy, author_username)
		VALUES ($1, $2)
		RETURNING id
	`
	var subscribeID string
	err := r.db.QueryRowContext(ctx, query, walletID, author).Scan(&subscribeID)
	if err != nil {
		logger.Errorf("failed to insert and get id for dex wallet %s author=%s: %v", walletID, author, err)
		return "", err
	}
	return subscribeID, nil
}

func (r *DexRepo) UnsubscribeAuthor(ctx context.Context, author string, walletID string) error {
	query := `
		DELETE FROM crypto_copytrade_authors_privy
		WHERE crypto_user_wallet_id_privy = $1
		AND author_username = $2
		RETURNING id
	`
	var subscribeID string
	if err := r.db.QueryRowContext(ctx, query, walletID, author).Scan(&subscribeID); err != nil {
		logger.Errorf("failed to delete and get id for dex wallet %s author=%s: %v", walletID, author, err)
		return err
	}
	return nil
}
