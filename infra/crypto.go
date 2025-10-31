package infra

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

var (
	cryptoOnce sync.Once
	CryptoDB   *sqlx.DB
)

// InitPostgresDB initializes the PostgreSQL database connection
func InitCryptoDB() error {
	var err error
	cryptoOnce.Do(func() {
		CryptoDB, err = connectCryptoDB()
	})
	return err
}

// GetPostgresConnection returns the PostgreSQL database connection
func GetCryptoConnection() (*sqlx.DB, error) {
	if CryptoDB == nil {
		return nil, fmt.Errorf("PostgreSQL connection not initialized")
	}
	return CryptoDB, nil
}

func connectCryptoDB() (*sqlx.DB, error) {
	cfg := config.GetConfig()

	dbUser := cfg.CryptoDatabase.User
	dbPass := cfg.CryptoDatabase.Password
	dbName := cfg.CryptoDatabase.DBName
	dbHost := cfg.CryptoDatabase.Host
	dbPort := cfg.CryptoDatabase.Port
	dbUnix := cfg.CryptoDatabase.Unix
	var dsn string
	if dbUnix != "" {
		// Use Unix socket connection
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=%s",
			dbUser, dbPass, dbName, dbUnix)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName)
	}

	logger.Info(fmt.Sprintf("Connecting to Crypto PostgreSQL with config: dsn=%s", dsn))

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error(err)
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}
