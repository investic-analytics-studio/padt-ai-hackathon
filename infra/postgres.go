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
	postgresOnce sync.Once
	PostgresDB   *sqlx.DB
)

// InitPostgresDB initializes the PostgreSQL database connection
func InitPostgresDB() error {
	var err error
	postgresOnce.Do(func() {
		PostgresDB, err = connectPostgresDB()
	})
	return err
}

// GetPostgresConnection returns the PostgreSQL database connection
func GetPostgresConnection() (*sqlx.DB, error) {
	if PostgresDB == nil {
		return nil, fmt.Errorf("PostgreSQL connection not initialized")
	}
	return PostgresDB, nil
}

func connectPostgresDB() (*sqlx.DB, error) {
	cfg := config.GetConfig()

	dbUser := cfg.PostgresDatabase.User
	dbPass := cfg.PostgresDatabase.Password
	dbName := cfg.PostgresDatabase.DBName
	dbHost := cfg.PostgresDatabase.Host
	dbPort := cfg.PostgresDatabase.Port
	dbUnix := cfg.PostgresDatabase.Unix
	var dsn string
	if dbUnix != "" {
		// Use Unix socket connection
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=%s",
			dbUser, dbPass, dbName, dbUnix)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName)
	}

	logger.Info(fmt.Sprintf("Connecting to PostgreSQL with config: dsn=%s", dsn))

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error(err)
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
