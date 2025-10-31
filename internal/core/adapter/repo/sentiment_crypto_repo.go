package repo

import (
	"context"
	"sort"
	"strings"

	"github.com/quantsmithapp/datastation-backend/internal/constant"
	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type sentimentCryptoRepo struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewSentimentCryptoRepo(db *gorm.DB, logger logger.Logger) domain.SentimentCryptoRepository {
	return &sentimentCryptoRepo{
		db:     db,
		logger: logger,
	}
}
func (r *sentimentCryptoRepo) GetSentimentAggregateCountRow(ctx context.Context, ticker, timeRange, sourceNames, topics string) (int, error) {
	r.logger.Info("GetSentimentAggregateCountRow repo called")
	tx := r.db.Session(&gorm.Session{})

	var totalSentimentRowQueryBuilder strings.Builder
	args := make([]interface{}, 0)

	totalSentimentRowQueryBuilder.WriteString(`
		SELECT COUNT(*) 
		FROM (
			SELECT 1 
			FROM crypto_sentiment 
			WHERE tickers != 'T'
	`)

	// Apply filters
	addSentimentTickerFilter(&totalSentimentRowQueryBuilder, &args, "tickers LIKE ?", ticker)
	addSentimentTimeFilter(&totalSentimentRowQueryBuilder, timeRange)
	addSentimentStringFilter(&totalSentimentRowQueryBuilder, &args, "source_name LIKE ?", sourceNames)
	addSentimentStringFilter(&totalSentimentRowQueryBuilder, &args, "topics LIKE ?", topics)

	// Finalize query
	totalSentimentRowQueryBuilder.WriteString(`
			GROUP BY date, sentiment
			ORDER BY date DESC
		) AS count_table
	`)

	// Get Total Count for all rows
	var totalSentimentRow int
	if err := tx.Raw(totalSentimentRowQueryBuilder.String(), args...).Scan(&totalSentimentRow).Error; err != nil {
		r.logger.Error(err, zap.String("query", "repo GetSentimentAggregateCountRow count"))
		return 0, err
	}
	return totalSentimentRow, nil
}
func (r *sentimentCryptoRepo) GetSentimentAggregate(ctx context.Context, ticker, timeRange, sourceNames, topics, limit, offset string) ([]model.SentimentCrypto, error) {
	r.logger.Info("GetSentimentAggregate repo called")
	tx := r.db.Session(&gorm.Session{})

	var sentimentQueryBuilder strings.Builder
	args := make([]interface{}, 0)

	// Base query
	sentimentQueryBuilder.WriteString(`
		SELECT date, sentiment, COUNT(*) as count
		FROM crypto_sentiment
		WHERE tickers != 'T'
	`)

	// Apply filters
	addSentimentTickerFilter(&sentimentQueryBuilder, &args, "tickers LIKE ?", ticker)
	addSentimentTimeFilter(&sentimentQueryBuilder, timeRange)
	addSentimentStringFilter(&sentimentQueryBuilder, &args, "source_name LIKE ?", sourceNames)
	addSentimentStringFilter(&sentimentQueryBuilder, &args, "topics LIKE ?", topics)

	// Finalize query with grouping and ordering
	sentimentQueryBuilder.WriteString(`
		GROUP BY date, sentiment
		ORDER BY date DESC
	`)

	// Apply pagination
	sentimentQueryBuilder.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, limit, offset)
	// Execute query
	var result []model.SentimentCrypto
	err := tx.Raw(sentimentQueryBuilder.String(), args...).Scan(&result).Error
	if err != nil {
		r.logger.Error(err, zap.String("query", "repo GetSentimentAggregate"))
		return nil, err
	}

	r.logger.Info("GetSentimentAggregate repo successful", zap.Int("results_count", len(result)))
	return result, nil
}

func (r *sentimentCryptoRepo) GetUniqueTickers(ctx context.Context) ([]model.TickerInfo, error) {
	r.logger.Info("GetUniqueTickers repo called")
	var rawTickers []string
	err := r.db.Table("crypto_sentiment").
		Select("DISTINCT SUBSTRING_INDEX(SUBSTRING_INDEX(tickers, ',', n.n), ',', -1) as ticker").
		Joins("JOIN (SELECT 1 as n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4) n ON CHAR_LENGTH(tickers) - CHAR_LENGTH(REPLACE(tickers, ',', '')) >= n.n - 1").
		Where("tickers != ''").
		Pluck("ticker", &rawTickers).Error
	if err != nil {
		r.logger.Error(err, zap.String("query", "repo GetUniqueTickers"))
		return nil, err
	}

	// Process tickers and create map for deduplication
	tickerMap := make(map[string]model.TickerInfo)
	for _, ticker := range rawTickers {
		// Remove surrounding quotes, spaces, and square brackets
		ticker = strings.Trim(ticker, "' []")
		// Remove leading $ if present
		ticker = strings.TrimPrefix(ticker, "$")
		// Convert to uppercase
		ticker = strings.ToUpper(ticker)
		// Remove any trailing numbers or apostrophes
		ticker = strings.TrimRight(ticker, "0123456789'")

		// Convert to lowercase for comparison with CRYPTO_TICKERS
		tickerLower := strings.ToLower(ticker)

		// Only add if it exists in CRYPTO_TICKERS constant
		if name, exists := constant.CRYPTO_TICKERS[tickerLower]; exists && ticker != "" {
			tickerMap[ticker] = model.TickerInfo{
				Symbol: ticker,
				Name:   name,
			}
		}
	}

	// Convert map to slice
	uniqueTickers := make([]model.TickerInfo, 0, len(tickerMap))
	for _, info := range tickerMap {
		uniqueTickers = append(uniqueTickers, info)
	}

	// Sort the tickers alphabetically by symbol
	sort.Slice(uniqueTickers, func(i, j int) bool {
		return uniqueTickers[i].Symbol < uniqueTickers[j].Symbol
	})

	r.logger.Info("GetUniqueTickers repo successful", zap.Int("tickers_count", len(uniqueTickers)))
	return uniqueTickers, nil
}

func (r *sentimentCryptoRepo) GetUniqueTopics(ctx context.Context) ([]string, error) {
	r.logger.Info("GetUniqueTopics repo called")
	var rawTopics []string
	err := r.db.Table("crypto_sentiment").
		Select("DISTINCT topics").
		Where("topics != ''").
		Pluck("topics", &rawTopics).Error
	if err != nil {
		r.logger.Error(err, zap.String("query", "repo GetUniqueTopics"))
		return nil, err
	}

	// Process topics to remove duplicates and standardize format
	topicMap := make(map[string]bool)
	for _, topic := range rawTopics {
		// Split the topic string in case it contains multiple topics
		splitTopics := strings.Split(topic, ",")
		for _, t := range splitTopics {
			// Remove surrounding quotes, spaces, and square brackets
			t = strings.Trim(t, "' []")
			// Convert to title case (capitalize first letter of each word)
			t = strings.Title(strings.ToLower(t))
			// Remove any trailing numbers or apostrophes
			t = strings.TrimRight(t, "0123456789'")
			// Add to map (this automatically removes duplicates)
			if t != "" {
				topicMap[t] = true
			}
		}
	}

	// Convert map keys to slice
	uniqueTopics := make([]string, 0, len(topicMap))
	for topic := range topicMap {
		uniqueTopics = append(uniqueTopics, topic)
	}

	// Sort the topics alphabetically
	sort.Strings(uniqueTopics)

	r.logger.Info("GetUniqueTopics repo successful", zap.Int("topics_count", len(uniqueTopics)))
	return uniqueTopics, nil
}

func (r *sentimentCryptoRepo) GetUniqueSourceNames(ctx context.Context) ([]string, error) {
	r.logger.Info("GetUniqueSourceNames repo called")
	var sourceNames []string
	err := r.db.Table("crypto_sentiment").
		Select("DISTINCT source_name").
		Where("source_name != ''").
		Pluck("source_name", &sourceNames).Error
	if err != nil {
		r.logger.Error(err, zap.String("query", "repo GetUniqueSourceNames"))
		return nil, err
	}

	// Process source names to remove duplicates and standardize format
	sourceNameMap := make(map[string]bool)
	for _, name := range sourceNames {
		// Remove surrounding spaces
		name = strings.TrimSpace(name)
		// Convert to title case (capitalize first letter of each word)
		name = strings.Title(strings.ToLower(name))
		// Add to map (this automatically removes duplicates)
		if name != "" {
			sourceNameMap[name] = true
		}
	}

	// Convert map keys to slice
	uniqueSourceNames := make([]string, 0, len(sourceNameMap))
	for name := range sourceNameMap {
		uniqueSourceNames = append(uniqueSourceNames, name)
	}

	// Sort the source names alphabetically
	sort.Strings(uniqueSourceNames)

	r.logger.Info("GetUniqueSourceNames repo successful", zap.Int("source_names_count", len(uniqueSourceNames)))
	return uniqueSourceNames, nil
}

// Helper to add a single filter if the value is not empty
func addSentimentTickerFilter(builder *strings.Builder, args *[]interface{}, condition, value string) {
	if value != "" {
		builder.WriteString(" AND " + condition)
		*args = append(*args, "%"+value+"%")
	}
}

// Helper to add time range filter
func addSentimentTimeFilter(builder *strings.Builder, timeRange string) {
	if filter := constant.GetTimeFilter(timeRange); filter != "" {
		builder.WriteString(filter)
	}
}

// Helper to add multiple OR conditions for filters like sourceNames and topics
func addSentimentStringFilter(builder *strings.Builder, args *[]interface{}, condition, values string) {
	if values == "" {
		return
	}

	items := strings.Split(values, ",")
	builder.WriteString(" AND (")

	for i, item := range items {
		if i > 0 {
			builder.WriteString(" OR ")
		}
		builder.WriteString(condition)
		*args = append(*args, "%"+strings.TrimSpace(item)+"%")
	}
	builder.WriteString(")")
}
