package repo

import (
	"database/sql"
	"fmt"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
)

type timescaleRepo struct {
	db     *sql.DB
	logger logger.Logger
}

func NewTimescaleRepo(db *sql.DB) port.TimescaleRepo {
	return &timescaleRepo{
		db:     db,
		logger: logger.NewLogger(),
	}
}

func (r *timescaleRepo) GetCryptoOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error) {
	return r.getOHLCV("binance", req)
}

func (r *timescaleRepo) GetForexOHLCV(req model.OHLCVRequest) ([]model.OHLCVData, error) {
	return r.getOHLCV("forex", req)
}

func (r *timescaleRepo) getOHLCV(dataType string, req model.OHLCVRequest) ([]model.OHLCVData, error) {
	tableName := fmt.Sprintf("%s_%s", dataType, req.TimeFrame)

	// Log the table name and request details for debugging
	r.logger.Info(fmt.Sprintf("Querying table: %s with request: %+v", tableName, req))

	// Check if the table exists
	var tableExists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", tableName).Scan(&tableExists)
	if err != nil {
		r.logger.Error(fmt.Errorf("error checking if table exists: %v", err))
		return nil, err
	}
	if !tableExists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	// Get the table schema
	rows, err := r.db.Query("SELECT column_name FROM information_schema.columns WHERE table_name = $1", tableName)
	if err != nil {
		r.logger.Error(fmt.Errorf("error getting table schema for %s: %v", tableName, err))
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			r.logger.Error(fmt.Errorf("error scanning column name: %v", err))
			return nil, err
		}
		columns = append(columns, column)
	}
	r.logger.Info(fmt.Sprintf("Table %s columns: %v", tableName, columns))

	// Determine the time column name
	timeColumn := "bucket" // Default to "bucket"
	if contains(columns, "time") {
		timeColumn = "time"
	} else if !contains(columns, "bucket") {
		return nil, fmt.Errorf("neither 'time' nor 'bucket' column found in table %s", tableName)
	}

	// Construct the SQL query
	query := fmt.Sprintf(`
		SELECT %s, ticker, open, high, low, close, volume
		FROM %s
		WHERE %s >= $1
	`, timeColumn, tableName, timeColumn)

	args := []interface{}{req.StartDate}
	argCount := 1

	if req.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND %s <= $%d", timeColumn, argCount)
		args = append(args, req.EndDate)
	}

	if !req.AllPair {
		argCount++
		query += fmt.Sprintf(" AND ticker = $%d", argCount)
		args = append(args, req.Ticker)
	}

	query += fmt.Sprintf(" ORDER BY %s, ticker", timeColumn)

	// Log the full SQL query and arguments
	r.logger.Info(fmt.Sprintf("Executing SQL query: %s with args: %v", query, args))

	// Execute the query
	rows, err = r.db.Query(query, args...)
	if err != nil {
		r.logger.Error(fmt.Errorf("error querying %s: %v", tableName, err))
		return nil, err
	}
	defer rows.Close()

	var result []model.OHLCVData
	for rows.Next() {
		var data model.OHLCVData
		err := rows.Scan(&data.Time, &data.Ticker, &data.Open, &data.High, &data.Low, &data.Close, &data.Volume)
		if err != nil {
			r.logger.Error(fmt.Errorf("error scanning row: %v", err))
			return nil, err
		}
		result = append(result, data)
	}

	// Log the number of results and the first result for debugging
	r.logger.Info(fmt.Sprintf("Query returned %d results", len(result)))
	if len(result) > 0 {
		r.logger.Info(fmt.Sprintf("First result: %+v", result[0]))
	}

	return result, nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
