package infra

import (
	"log"
	"os"
	"time"

	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

var AppDB *gorm.DB

func InitAppDB() {
	var err error
	cfg := config.GetConfig()

	StockDB, err = NewGormDB(cfg.StockDatabase)
	if err != nil {
		logger.Fatal(err)
	}
	AnalyticDB, err = NewGormDB(cfg.AnalyticDatabase)
	if err != nil {
		logger.Fatal(err)
	}
}

var StockDB *gorm.DB
var AnalyticDB *gorm.DB

func InitStockDB() {
	var err error
	cfg := config.GetConfig()

	StockDB, err = NewGormDB(cfg.StockDatabase)
	if err != nil {
		logger.Fatal(err)
	}
}

func InitAnalyticDB() {
	var err error
	cfg := config.GetConfig()
	AnalyticDB, err = NewGormDB(cfg.AnalyticDatabase)
	if err != nil {
		// Handle error, perhaps log it and exit
		panic("failed to connect to analytic database")
	}
}

func NewGormDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	newLogger := gormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlog.Config{
			SlowThreshold:             time.Second,    // Slow SQL threshold
			LogLevel:                  gormlog.Silent, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,           // Don't include params in the SQL log
			Colorful:                  true,           // Disable color
		},
	)

	return gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: newLogger,
	})
}
