package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

var TimescaleDB *sql.DB

func InitTimescaleDB() {
	var err error
	cfg := config.GetConfig()
	TimescaleDB, err = NewTimescaleDBConnection(cfg.TimescaleDatabase)
	if err != nil {
		logger.Fatal(err)
	}
}

func NewTimescaleDBConnection(cfg config.TimescaleConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

func GetTimescaleDBConnection() (*sql.DB, error) {
	if TimescaleDB == nil {
		return nil, fmt.Errorf("TimescaleDB connection not initialized")
	}
	return TimescaleDB, nil
}
