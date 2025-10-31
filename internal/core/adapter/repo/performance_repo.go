package repo

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type PerformanceRepo struct {
	cryptoDB       *sqlx.DB // For crypto_author_nav
	postgresDB     *sqlx.DB // For twitter_crypto_* tables
	authorTierRepo *AuthorTierRepo
}

func NewPerformanceRepo(cryptoDB *sqlx.DB, postgresDB *sqlx.DB, authorTierRepo *AuthorTierRepo) *PerformanceRepo {
	return &PerformanceRepo{
		cryptoDB:       cryptoDB,
		postgresDB:     postgresDB,
		authorTierRepo: authorTierRepo,
	}
}

func (r *PerformanceRepo) MultiholdingPortNavRepo(ctx context.Context, period string, holdingPeriod string) ([]model.AuthorNav, error) {
	var (
		data  []model.AuthorNav
		since time.Time
		rows  []struct {
			AuthorUsername string    `db:"author_username"`
			Datetime       time.Time `db:"day"`
			Nav            float64   `db:"nav"`
		}
	)

	// Build query based on period and holdingPeriod
	query, args, err := r.buildMultiholdingNavQuery(period, holdingPeriod, &since)
	if err != nil {
		return nil, err
	}

	// Execute query
	if err := r.cryptoDB.SelectContext(ctx, &rows, query, args...); err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Group NAV data by author
	navMap := r.groupNavDataByAuthor(rows)

	// Normalize NAV values to percentage (base 100)
	r.normalizeNavValues(navMap)

	// Calculate performance metrics and build result
	authorNames, err := r.authorTierRepo.GetAuthorNameMap()
	if err != nil {
		return nil, err
	}

	data = r.buildAuthorNavData(navMap, authorNames)

	// Sort by ROI and assign ranks
	r.sortAndRankAuthors(&data)

	return data, nil
}

// buildMultiholdingNavQuery creates the appropriate SQL query based on the period and holdingPeriod
func (r *PerformanceRepo) buildMultiholdingNavQuery(period string, holdingPeriod string, since *time.Time) (string, []interface{}, error) {
	now := time.Now()
	var query string
	var args []interface{}

	// Validate holdingPeriod
	validHoldingPeriods := map[string]bool{
		"24": true, "48": true, "72": true, "96": true,
		"120": true, "144": true, "168": true,
	}

	if !validHoldingPeriods[holdingPeriod] {
		return "", nil, errors.New("invalid holding period")
	}

	// Determine the nav column to use based on holdingPeriod
	navColumn := fmt.Sprintf("nav_%s", holdingPeriod)

	switch period {
	case "7D":
		*since = now.AddDate(0, 0, -7)
		query = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					` + navColumn + ` AS nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_port_nav
				WHERE datetime >= $1
			) t
			WHERE t.rn = 1
			ORDER BY t.author_username, t.day ASC;
		`
		args = []interface{}{*since}
	case "1M":
		*since = now.AddDate(0, -1, 0)
		query = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					` + navColumn + ` AS nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_port_nav
				WHERE datetime >= $1
			) t
			WHERE t.rn = 1
			ORDER BY t.author_username, t.day ASC;
		`
		args = []interface{}{*since}
	case "ALL":
		query = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE_TRUNC('month', datetime) AS month,
					DATE(datetime) AS day,
					` + navColumn + ` AS nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE_TRUNC('month', datetime) ORDER BY datetime ASC) AS rn
				FROM crypto_author_port_nav
			) t
			WHERE t.rn = 1
			ORDER BY t.author_username, t.day ASC;
		`
		args = []interface{}{}
	default:
		return "", nil, errors.New("invalid period")
	}

	return query, args, nil
}

func (r *PerformanceRepo) AuthorNavRepo(ctx context.Context, period string) ([]model.AuthorNav, error) {
	var (
		data  []model.AuthorNav
		since time.Time
		rows  []struct {
			AuthorUsername string    `db:"author_username"`
			Datetime       time.Time `db:"day"`
			Nav            float64   `db:"nav"`
		}
	)

	// Build query based on period
	query, args, err := r.buildNavQuery(period, &since)
	if err != nil {
		return nil, err
	}

	// Execute query
	if err := r.cryptoDB.SelectContext(ctx, &rows, query, args...); err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Group NAV data by author
	navMap := r.groupNavDataByAuthor(rows)

	// Normalize NAV values to percentage (base 100)
	r.normalizeNavValues(navMap)

	// Calculate performance metrics and build result
	authorNames, err := r.authorTierRepo.GetAuthorNameMap()
	if err != nil {
		return nil, err
	}

	data = r.buildAuthorNavData(navMap, authorNames)

	// Sort by ROI and assign ranks
	r.sortAndRankAuthors(&data)

	return data, nil
}

func (r *PerformanceRepo) AuthorDetailRepo(ctx context.Context, authorUsername string, period string, start int, limit int) (model.AuthorDetail, error) {
	var (
		data    model.AuthorDetail
		since   time.Time
		navRows []struct {
			AuthorUsername string    `db:"author_username"`
			Datetime       time.Time `db:"day"`
			Nav            float64   `db:"nav"`
		}
	)

	// Validate input
	if authorUsername == "" {
		return data, errors.New("author username cannot be empty")
	}

	// Validate pagination parameters
	if start < 0 {
		start = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}

	fmt.Printf("Using pagination parameters - start: %d, limit: %d\n", start, limit)

	now := time.Now()
	var navQuery string
	var navArgs []interface{}

	// 1. Get NAV data based on period
	switch period {
	case "7D":
		since = now.AddDate(0, 0, -7)
		navQuery = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_nav
				WHERE datetime >= $1 AND author_username = $2
			) t
			WHERE t.rn = 1
			ORDER BY t.day ASC;
		`
		navArgs = []interface{}{since, authorUsername}
	case "1M":
		since = now.AddDate(0, -1, 0)
		navQuery = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_nav
				WHERE datetime >= $1 AND author_username = $2
			) t
			WHERE t.rn = 1
			ORDER BY t.day ASC;
		`
		navArgs = []interface{}{since, authorUsername}
	case "ALL":
		navQuery = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE_TRUNC('month', datetime) AS month,
					DATE(datetime) AS day,
					nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE_TRUNC('month', datetime) ORDER BY datetime ASC) AS rn
				FROM crypto_author_nav
				WHERE author_username = $1
			) t
			WHERE t.rn = 1
			ORDER BY t.day ASC;
		`
		navArgs = []interface{}{authorUsername}
	default:
		return data, errors.New("invalid period")
	}

	if err := r.cryptoDB.SelectContext(ctx, &navRows, navQuery, navArgs...); err != nil {
		fmt.Printf("Error fetching NAV data for username %s: %v\n", authorUsername, err)
		return data, fmt.Errorf("failed to get NAV data for username %s: %w", authorUsername, err)
	}

	if len(navRows) == 0 {
		return data, fmt.Errorf("no NAV data found for author: %s", authorUsername)
	}

	// Process NAV data
	var navs []model.Nav
	for _, row := range navRows {
		navs = append(navs, model.Nav{
			Datetime: row.Datetime,
			Nav:      row.Nav,
		})
	}

	if len(navs) == 0 || navs[0].Nav == 0 {
		return data, fmt.Errorf("invalid NAV data for author: %s", authorUsername)
	}

	// Normalize NAV to percentage (base 100)
	firstNav := navs[0].Nav
	for i, n := range navs {
		navs[i].Nav = 100.0 * n.Nav / firstNav
		if math.IsNaN(navs[i].Nav) || math.IsInf(navs[i].Nav, 0) {
			navs[i].Nav = 100.0 // fallback safe value
		}
	}

	// Calculate ROI and drawdowns
	var roi float64
	if len(navs) > 0 {
		roi = navs[len(navs)-1].Nav - 100
		if math.IsNaN(roi) || math.IsInf(roi, 0) {
			roi = 0
		}
	}

	var peak, maxDrawdown, currentDrawdown float64
	if len(navs) > 0 {
		peak = navs[0].Nav
		for _, nav := range navs {
			if nav.Nav > peak {
				peak = nav.Nav
			}
			drawdown := (peak - nav.Nav) / peak
			if math.IsNaN(drawdown) || math.IsInf(drawdown, 0) {
				drawdown = 0
			}
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
			}
		}
		currentDrawdown = (peak - navs[len(navs)-1].Nav) / peak
		if math.IsNaN(currentDrawdown) || math.IsInf(currentDrawdown, 0) {
			currentDrawdown = 0
		}
		if math.IsNaN(maxDrawdown) || math.IsInf(maxDrawdown, 0) {
			maxDrawdown = 0
		}
	}

	startNav := navs[0].Nav
	endNav := navs[len(navs)-1].Nav

	// 2. Get author profile
	var profile model.AuthorProfile
	profileQuery := `
		SELECT 
			author_username,
			COALESCE(author_url, '') as author_url,
			COALESCE(author_twitterurl, '') as author_twitterurl,
			COALESCE(author_name, '') as author_name,
			COALESCE(author_followers, 0) as author_followers,
			COALESCE(author_following, 0) as author_following,
			COALESCE(author_tier, '') as author_tier,
			COALESCE(created_at, NOW()) as created_at,
			COALESCE(updated_at, NOW()) as updated_at
		FROM twitter_crypto_author_profile 
		WHERE author_username = $1
	`

	fmt.Printf("Executing profile query for: %s\n", authorUsername)
	if err := r.postgresDB.GetContext(ctx, &profile, profileQuery, authorUsername); err != nil {
		fmt.Printf("Warning: Could not fetch profile for %s: %v\n", authorUsername, err)
		// Continue without profile data
		profile = model.AuthorProfile{
			AuthorUsername: authorUsername,
			AuthorName:     authorUsername,
		}
	} else {
		fmt.Printf("Successfully fetched profile for %s: %+v\n", authorUsername, profile)
	}

	// 3. Get recent tweets with pagination
	var tweets []model.AuthorTweet
	tweetsQuery := `
		SELECT 
			id, url, text, source,
			retweetcount, replycount, likecount, quotecount, viewcount,
			tweet_created_at, bookmarkcount,
			isreply, conversationid, ispinned, isretweet, isquote,
			COALESCE(media_url, '') as media_url,
			tweet_created_date,
			COALESCE(tickers_rule_based, '') as tickers_rule_based
		FROM twitter_crypto_tweets_foxhole 
		WHERE author_username = $1
		ORDER BY tweet_created_at DESC
		LIMIT $2 OFFSET $3
	`

	fmt.Printf("Executing tweets query for: %s (limit: %d, offset: %d)\n", authorUsername, limit, start)
	if err := r.postgresDB.SelectContext(ctx, &tweets, tweetsQuery, authorUsername, limit, start); err != nil {
		fmt.Printf("Warning: Could not fetch tweets for %s: %v\n", authorUsername, err)
		tweets = []model.AuthorTweet{} // Empty array if no tweets found
	} else {
		fmt.Printf("Successfully fetched %d tweets for %s\n", len(tweets), authorUsername)
		// Sanitize tweet data to prevent NaN/Infinite values
		for i := range tweets {
			if tweets[i].ViewCount != nil {
				if math.IsNaN(*tweets[i].ViewCount) || math.IsInf(*tweets[i].ViewCount, 0) {
					tweets[i].ViewCount = nil // Set to nil if invalid
				}
			}
		}
	}

	// 4. Get recent signals with pagination
	var signals []model.AuthorSignal
	signalsQuery := `
		SELECT 
			s.tweet_id, s.content, s.ticker, s.action, s.score, s.sentiment, s.prompt_version,
			s.created_at, s.updated_at
		FROM twitter_crypto_signal s
		INNER JOIN twitter_crypto_tweets_foxhole t ON s.tweet_id = t.id
		WHERE t.author_username = $1
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3
	`

	fmt.Printf("Executing signals query for: %s (limit: %d, offset: %d)\n", authorUsername, limit, start)
	if err := r.postgresDB.SelectContext(ctx, &signals, signalsQuery, authorUsername, limit, start); err != nil {
		fmt.Printf("Warning: Could not fetch signals for %s: %v\n", authorUsername, err)
		signals = []model.AuthorSignal{} // Empty array if no signals found
	} else {
		fmt.Printf("Successfully fetched %d signals for %s\n", len(signals), authorUsername)
	}

	// 5. Get total tweet count
	var totalTweets int
	totalTweetsQuery := `
		SELECT COUNT(*) 
		FROM twitter_crypto_tweets_foxhole 
		WHERE author_username = $1
	`

	fmt.Printf("Executing total tweets count query for: %s\n", authorUsername)
	if err := r.postgresDB.GetContext(ctx, &totalTweets, totalTweetsQuery, authorUsername); err != nil {
		fmt.Printf("Warning: Could not fetch total tweets count for %s: %v\n", authorUsername, err)
		totalTweets = 0
	} else {
		fmt.Printf("Total tweets count for %s: %d\n", authorUsername, totalTweets)
	}

	// 6. Get total signal count
	var totalSignals int
	totalSignalsQuery := `
		SELECT COUNT(*) 
		FROM twitter_crypto_signal s
		INNER JOIN twitter_crypto_tweets_foxhole t ON s.tweet_id = t.id
		WHERE t.author_username = $1
	`

	fmt.Printf("Executing total signals count query for: %s\n", authorUsername)
	if err := r.postgresDB.GetContext(ctx, &totalSignals, totalSignalsQuery, authorUsername); err != nil {
		fmt.Printf("Warning: Could not fetch total signals count for %s: %v\n", authorUsername, err)
		totalSignals = 0
	} else {
		fmt.Printf("Total signals count for %s: %d\n", authorUsername, totalSignals)
	}

	// Get author display name from authorTierRepo
	authorNames, err := r.authorTierRepo.GetAuthorNameMap()
	if err != nil {
		fmt.Printf("Warning: Could not fetch author names: %v\n", err)
	}

	authorDisplayName := authorUsername
	if authorNames != nil {
		if name, exists := authorNames[authorUsername]; exists && name != "" {
			authorDisplayName = name
			fmt.Printf("Found author display name: %s -> %s\n", authorUsername, name)
		} else {
			fmt.Printf("Author %s not found in authorNames map, using username as display name\n", authorUsername)
			// Try to get name from profile if available
			if profile.AuthorName != "" && profile.AuthorName != authorUsername {
				authorDisplayName = profile.AuthorName
				fmt.Printf("Using profile author name: %s\n", authorDisplayName)
			}
		}
	}

	fmt.Printf("Final author display name: %s\n", authorDisplayName)
	fmt.Printf("Profile data: %+v\n", profile)
	fmt.Printf("Tweets count: %d (Total: %d)\n", len(tweets), totalTweets)
	fmt.Printf("Signals count: %d (Total: %d)\n", len(signals), totalSignals)

	// Create a map to quickly find signals by tweet ID
	signalMap := make(map[string]*model.AuthorSignal)
	for i := range signals {
		signalMap[signals[i].TweetID] = &signals[i]
	}

	// Merge tweets with their signals
	var mergedTimeline []model.AuthorTweetWithSignals
	for _, tweet := range tweets {
		timelineItem := model.AuthorTweetWithSignals{
			// Copy all tweet fields
			ID:               tweet.ID,
			URL:              tweet.URL,
			Text:             tweet.Text,
			Source:           tweet.Source,
			RetweetCount:     tweet.RetweetCount,
			ReplyCount:       tweet.ReplyCount,
			LikeCount:        tweet.LikeCount,
			QuoteCount:       tweet.QuoteCount,
			ViewCount:        tweet.ViewCount,
			TweetCreatedAt:   tweet.TweetCreatedAt,
			BookmarkCount:    tweet.BookmarkCount,
			IsReply:          tweet.IsReply,
			ConversationID:   tweet.ConversationID,
			IsPinned:         tweet.IsPinned,
			IsRetweet:        tweet.IsRetweet,
			IsQuote:          tweet.IsQuote,
			MediaURL:         tweet.MediaURL,
			TweetCreatedDate: tweet.TweetCreatedDate,
			TickersRuleBased: tweet.TickersRuleBased,
		}

		// Add signal data if it exists for this tweet
		if signal, exists := signalMap[tweet.ID]; exists {
			timelineItem.SignalTicker = &signal.Ticker
			timelineItem.SignalAction = &signal.Action
			timelineItem.SignalScore = &signal.Score
			timelineItem.SignalSentiment = &signal.Sentiment
			timelineItem.SignalPromptVersion = &signal.PromptVersion
			timelineItem.SignalCreatedAt = &signal.CreatedAt
			timelineItem.SignalUpdatedAt = &signal.UpdatedAt
		}

		mergedTimeline = append(mergedTimeline, timelineItem)
	}

	fmt.Printf("Merged timeline count: %d\n", len(mergedTimeline))

	// Calculate pagination info
	actualEnd := start + len(mergedTimeline)
	fmt.Printf("Pagination info - Start: %d, End: %d, Limit: %d, Total: %d\n", start, actualEnd, limit, totalTweets)

	// Combine all data
	data = model.AuthorDetail{
		AuthorName:      authorDisplayName,
		WeightNav:       navs,
		ROI:             roi,
		StartNav:        startNav,
		EndNav:          endNav,
		Drawdown:        currentDrawdown,
		MaximumDrawdown: maxDrawdown,
		Profile:         profile,
		RecentTimeline:  mergedTimeline,
		TotalTimeline:   totalTweets, // Use total tweets as the timeline count
		Start:           start,
		End:             actualEnd,
		Limit:           limit,
		BearishTokens:   make([]model.SentimentToken, 0), // Will be populated by service layer
		BullishTokens:   make([]model.SentimentToken, 0), // Will be populated by service layer
	}

	// Final sanitization check for all floating point values
	if math.IsNaN(data.ROI) || math.IsInf(data.ROI, 0) {
		data.ROI = 0.0
	}
	if math.IsNaN(data.StartNav) || math.IsInf(data.StartNav, 0) {
		data.StartNav = 100.0
	}
	if math.IsNaN(data.EndNav) || math.IsInf(data.EndNav, 0) {
		data.EndNav = 100.0
	}
	if math.IsNaN(data.Drawdown) || math.IsInf(data.Drawdown, 0) {
		data.Drawdown = 0.0
	}
	if math.IsNaN(data.MaximumDrawdown) || math.IsInf(data.MaximumDrawdown, 0) {
		data.MaximumDrawdown = 0.0
	}

	fmt.Printf("Final sanitized data - ROI: %f, StartNav: %f, EndNav: %f, Drawdown: %f, MaxDrawdown: %f\n",
		data.ROI, data.StartNav, data.EndNav, data.Drawdown, data.MaximumDrawdown)

	return data, nil
}

func (r *PerformanceRepo) AuthorMultiholdingDetailRepo(ctx context.Context, authorUsername string, period string, holdingPeriod string, start int, limit int) (model.AuthorDetail, error) {
	var (
		data    model.AuthorDetail
		since   time.Time
		navRows []struct {
			AuthorUsername string    `db:"author_username"`
			Datetime       time.Time `db:"day"`
			Nav            float64   `db:"nav"`
		}
	)

	// Validate input
	if authorUsername == "" {
		return data, errors.New("author username cannot be empty")
	}

	// Validate holding period
	validHoldingPeriods := map[string]bool{
		"24": true, "48": true, "72": true, "96": true,
		"120": true, "144": true, "168": true,
	}

	if !validHoldingPeriods[holdingPeriod] {
		return data, errors.New("invalid holding period")
	}

	// Validate pagination parameters
	if start < 0 {
		start = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}

	fmt.Printf("Using pagination parameters - start: %d, limit: %d\n", start, limit)

	now := time.Now()
	var navQuery string
	var navArgs []interface{}

	// Determine the nav column to use based on holdingPeriod
	navColumn := fmt.Sprintf("nav_%s", holdingPeriod)

	// 1. Get NAV data based on period using multiholding nav column
	switch period {
	case "7D":
		since = now.AddDate(0, 0, -7)
		navQuery = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					` + navColumn + ` AS nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_port_nav
				WHERE datetime >= $1 AND author_username = $2
			) t
			WHERE t.rn = 1
			ORDER BY t.day ASC;
		`
		navArgs = []interface{}{since, authorUsername}
	case "1M":
		since = now.AddDate(0, -1, 0)
		navQuery = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					` + navColumn + ` AS nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_port_nav
				WHERE datetime >= $1 AND author_username = $2
			) t
			WHERE t.rn = 1
			ORDER BY t.day ASC;
		`
		navArgs = []interface{}{since, authorUsername}
	case "ALL":
		navQuery = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE_TRUNC('month', datetime) AS month,
					DATE(datetime) AS day,
					` + navColumn + ` AS nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE_TRUNC('month', datetime) ORDER BY datetime ASC) AS rn
				FROM crypto_author_port_nav
				WHERE author_username = $1
			) t
			WHERE t.rn = 1
			ORDER BY t.day ASC;
		`
		navArgs = []interface{}{authorUsername}
	default:
		return data, errors.New("invalid period")
	}

	if err := r.cryptoDB.SelectContext(ctx, &navRows, navQuery, navArgs...); err != nil {
		fmt.Printf("Error fetching multiholding NAV data for username %s: %v\n", authorUsername, err)
		return data, fmt.Errorf("failed to get multiholding NAV data for username %s: %w", authorUsername, err)
	}

	if len(navRows) == 0 {
		return data, fmt.Errorf("no multiholding NAV data found for author: %s", authorUsername)
	}

	// Process NAV data
	var navs []model.Nav
	for _, row := range navRows {
		navs = append(navs, model.Nav{
			Datetime: row.Datetime,
			Nav:      row.Nav,
		})
	}

	if len(navs) == 0 || navs[0].Nav == 0 {
		return data, fmt.Errorf("invalid multiholding NAV data for author: %s", authorUsername)
	}

	// Normalize NAV to percentage (base 100)
	firstNav := navs[0].Nav
	for i, n := range navs {
		navs[i].Nav = 100.0 * n.Nav / firstNav
		if math.IsNaN(navs[i].Nav) || math.IsInf(navs[i].Nav, 0) {
			navs[i].Nav = 100.0 // fallback safe value
		}
	}

	// Calculate ROI and drawdowns
	var roi float64
	if len(navs) > 0 {
		roi = navs[len(navs)-1].Nav - 100
		if math.IsNaN(roi) || math.IsInf(roi, 0) {
			roi = 0
		}
	}

	var peak, maxDrawdown, currentDrawdown float64
	if len(navs) > 0 {
		peak = navs[0].Nav
		for _, nav := range navs {
			if nav.Nav > peak {
				peak = nav.Nav
			}
			drawdown := (peak - nav.Nav) / peak
			if math.IsNaN(drawdown) || math.IsInf(drawdown, 0) {
				drawdown = 0
			}
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
			}
		}
		currentDrawdown = (peak - navs[len(navs)-1].Nav) / peak
		if math.IsNaN(currentDrawdown) || math.IsInf(currentDrawdown, 0) {
			currentDrawdown = 0
		}
		if math.IsNaN(maxDrawdown) || math.IsInf(maxDrawdown, 0) {
			maxDrawdown = 0
		}
	}

	startNav := navs[0].Nav
	endNav := navs[len(navs)-1].Nav

	// 2. Get author profile (same as original method)
	var profile model.AuthorProfile
	profileQuery := `
		SELECT 
			author_username,
			COALESCE(author_url, '') as author_url,
			COALESCE(author_twitterurl, '') as author_twitterurl,
			COALESCE(author_name, '') as author_name,
			COALESCE(author_followers, 0) as author_followers,
			COALESCE(author_following, 0) as author_following,
			COALESCE(author_tier, '') as author_tier,
			COALESCE(created_at, NOW()) as created_at,
			COALESCE(updated_at, NOW()) as updated_at
		FROM twitter_crypto_author_profile 
		WHERE author_username = $1
	`

	fmt.Printf("Executing profile query for: %s\n", authorUsername)
	if err := r.postgresDB.GetContext(ctx, &profile, profileQuery, authorUsername); err != nil {
		fmt.Printf("Warning: Could not fetch profile for %s: %v\n", authorUsername, err)
		// Continue without profile data
		profile = model.AuthorProfile{
			AuthorUsername: authorUsername,
			AuthorName:     authorUsername,
		}
	} else {
		fmt.Printf("Successfully fetched profile for %s: %+v\n", authorUsername, profile)
	}

	// 3. Get recent tweets with pagination (same as original method)
	var tweets []model.AuthorTweet
	tweetsQuery := `
		SELECT 
			id, url, text, source,
			retweetcount, replycount, likecount, quotecount, viewcount,
			tweet_created_at, bookmarkcount,
			isreply, conversationid, ispinned, isretweet, isquote,
			COALESCE(media_url, '') as media_url,
			tweet_created_date,
			COALESCE(tickers_rule_based, '') as tickers_rule_based
		FROM twitter_crypto_tweets_foxhole 
		WHERE author_username = $1
		ORDER BY tweet_created_at DESC
		LIMIT $2 OFFSET $3
	`

	fmt.Printf("Executing tweets query for: %s (limit: %d, offset: %d)\n", authorUsername, limit, start)
	if err := r.postgresDB.SelectContext(ctx, &tweets, tweetsQuery, authorUsername, limit, start); err != nil {
		fmt.Printf("Warning: Could not fetch tweets for %s: %v\n", authorUsername, err)
		tweets = []model.AuthorTweet{} // Empty array if no tweets found
	} else {
		fmt.Printf("Successfully fetched %d tweets for %s\n", len(tweets), authorUsername)
		// Sanitize tweet data to prevent NaN/Infinite values
		for i := range tweets {
			if tweets[i].ViewCount != nil {
				if math.IsNaN(*tweets[i].ViewCount) || math.IsInf(*tweets[i].ViewCount, 0) {
					tweets[i].ViewCount = nil // Set to nil if invalid
				}
			}
		}
	}

	// 4. Get recent signals with pagination (same as original method)
	var signals []model.AuthorSignal
	signalsQuery := `
		SELECT 
			s.tweet_id, s.content, s.ticker, s.action, s.score, s.sentiment, s.prompt_version,
			s.created_at, s.updated_at
		FROM twitter_crypto_signal s
		INNER JOIN twitter_crypto_tweets_foxhole t ON s.tweet_id = t.id
		WHERE t.author_username = $1
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3
	`

	fmt.Printf("Executing signals query for: %s (limit: %d, offset: %d)\n", authorUsername, limit, start)
	if err := r.postgresDB.SelectContext(ctx, &signals, signalsQuery, authorUsername, limit, start); err != nil {
		fmt.Printf("Warning: Could not fetch signals for %s: %v\n", authorUsername, err)
		signals = []model.AuthorSignal{} // Empty array if no signals found
	} else {
		fmt.Printf("Successfully fetched %d signals for %s\n", len(signals), authorUsername)
	}

	// 5. Get total tweet count (same as original method)
	var totalTweets int
	totalTweetsQuery := `
		SELECT COUNT(*) 
		FROM twitter_crypto_tweets_foxhole 
		WHERE author_username = $1
	`

	fmt.Printf("Executing total tweets count query for: %s\n", authorUsername)
	if err := r.postgresDB.GetContext(ctx, &totalTweets, totalTweetsQuery, authorUsername); err != nil {
		fmt.Printf("Warning: Could not fetch total tweets count for %s: %v\n", authorUsername, err)
		totalTweets = 0
	} else {
		fmt.Printf("Total tweets count for %s: %d\n", authorUsername, totalTweets)
	}

	// 6. Get total signal count (same as original method)
	var totalSignals int
	totalSignalsQuery := `
		SELECT COUNT(*) 
		FROM twitter_crypto_signal s
		INNER JOIN twitter_crypto_tweets_foxhole t ON s.tweet_id = t.id
		WHERE t.author_username = $1
	`

	fmt.Printf("Executing total signals count query for: %s\n", authorUsername)
	if err := r.postgresDB.GetContext(ctx, &totalSignals, totalSignalsQuery, authorUsername); err != nil {
		fmt.Printf("Warning: Could not fetch total signals count for %s: %v\n", authorUsername, err)
		totalSignals = 0
	} else {
		fmt.Printf("Total signals count for %s: %d\n", authorUsername, totalSignals)
	}

	// Get author display name from authorTierRepo (same as original method)
	authorNames, err := r.authorTierRepo.GetAuthorNameMap()
	if err != nil {
		fmt.Printf("Warning: Could not fetch author names: %v\n", err)
	}

	authorDisplayName := authorUsername
	if authorNames != nil {
		if name, exists := authorNames[authorUsername]; exists && name != "" {
			authorDisplayName = name
			fmt.Printf("Found author display name: %s -> %s\n", authorUsername, name)
		} else {
			fmt.Printf("Author %s not found in authorNames map, using username as display name\n", authorUsername)
			// Try to get name from profile if available
			if profile.AuthorName != "" && profile.AuthorName != authorUsername {
				authorDisplayName = profile.AuthorName
				fmt.Printf("Using profile author name: %s\n", authorDisplayName)
			}
		}
	}

	fmt.Printf("Final author display name: %s\n", authorDisplayName)
	fmt.Printf("Profile data: %+v\n", profile)
	fmt.Printf("Tweets count: %d (Total: %d)\n", len(tweets), totalTweets)
	fmt.Printf("Signals count: %d (Total: %d)\n", len(signals), totalSignals)

	// Create a map to quickly find signals by tweet ID (same as original method)
	signalMap := make(map[string]*model.AuthorSignal)
	for i := range signals {
		signalMap[signals[i].TweetID] = &signals[i]
	}

	// Merge tweets with their signals (same as original method)
	var mergedTimeline []model.AuthorTweetWithSignals
	for _, tweet := range tweets {
		timelineItem := model.AuthorTweetWithSignals{
			// Copy all tweet fields
			ID:               tweet.ID,
			URL:              tweet.URL,
			Text:             tweet.Text,
			Source:           tweet.Source,
			RetweetCount:     tweet.RetweetCount,
			ReplyCount:       tweet.ReplyCount,
			LikeCount:        tweet.LikeCount,
			QuoteCount:       tweet.QuoteCount,
			ViewCount:        tweet.ViewCount,
			TweetCreatedAt:   tweet.TweetCreatedAt,
			BookmarkCount:    tweet.BookmarkCount,
			IsReply:          tweet.IsReply,
			ConversationID:   tweet.ConversationID,
			IsPinned:         tweet.IsPinned,
			IsRetweet:        tweet.IsRetweet,
			IsQuote:          tweet.IsQuote,
			MediaURL:         tweet.MediaURL,
			TweetCreatedDate: tweet.TweetCreatedDate,
			TickersRuleBased: tweet.TickersRuleBased,
		}

		// Add signal data if it exists for this tweet
		if signal, exists := signalMap[tweet.ID]; exists {
			timelineItem.SignalTicker = &signal.Ticker
			timelineItem.SignalAction = &signal.Action
			timelineItem.SignalScore = &signal.Score
			timelineItem.SignalSentiment = &signal.Sentiment
			timelineItem.SignalPromptVersion = &signal.PromptVersion
			timelineItem.SignalCreatedAt = &signal.CreatedAt
			timelineItem.SignalUpdatedAt = &signal.UpdatedAt
		}

		mergedTimeline = append(mergedTimeline, timelineItem)
	}

	fmt.Printf("Merged timeline count: %d\n", len(mergedTimeline))

	// Calculate pagination info
	actualEnd := start + len(mergedTimeline)
	fmt.Printf("Pagination info - Start: %d, End: %d, Limit: %d, Total: %d\n", start, actualEnd, limit, totalTweets)

	// Combine all data
	data = model.AuthorDetail{
		AuthorName:      authorDisplayName,
		WeightNav:       navs,
		ROI:             roi,
		StartNav:        startNav,
		EndNav:          endNav,
		Drawdown:        currentDrawdown,
		MaximumDrawdown: maxDrawdown,
		Profile:         profile,
		RecentTimeline:  mergedTimeline,
		TotalTimeline:   totalTweets, // Use total tweets as the timeline count
		Start:           start,
		End:             actualEnd,
		Limit:           limit,
		BearishTokens:   make([]model.SentimentToken, 0), // Will be populated by service layer
		BullishTokens:   make([]model.SentimentToken, 0), // Will be populated by service layer
	}

	// Final sanitization check for all floating point values
	if math.IsNaN(data.ROI) || math.IsInf(data.ROI, 0) {
		data.ROI = 0.0
	}
	if math.IsNaN(data.StartNav) || math.IsInf(data.StartNav, 0) {
		data.StartNav = 100.0
	}
	if math.IsNaN(data.EndNav) || math.IsInf(data.EndNav, 0) {
		data.EndNav = 100.0
	}
	if math.IsNaN(data.Drawdown) || math.IsInf(data.Drawdown, 0) {
		data.Drawdown = 0.0
	}
	if math.IsNaN(data.MaximumDrawdown) || math.IsInf(data.MaximumDrawdown, 0) {
		data.MaximumDrawdown = 0.0
	}

	fmt.Printf("Final sanitized multiholding data - ROI: %f, StartNav: %f, EndNav: %f, Drawdown: %f, MaxDrawdown: %f\n",
		data.ROI, data.StartNav, data.EndNav, data.Drawdown, data.MaximumDrawdown)

	return data, nil
}

func (r *PerformanceRepo) GetAuthorSentimentAnalysis(ctx context.Context, authorUsername string, period string) (bearishTokens []model.SentimentToken, bullishTokens []model.SentimentToken, err error) {
	// Initialize empty slices to avoid null in JSON response
	bearishTokens = []model.SentimentToken{}
	bullishTokens = []model.SentimentToken{}

	// Validate input
	if authorUsername == "" {
		return bearishTokens, bullishTokens, errors.New("author username cannot be empty")
	}

	now := time.Now()
	var since time.Time
	var query string
	var args []interface{}

	// Determine time range based on period
	switch period {
	case "7D":
		since = now.AddDate(0, 0, -7)
		query = `
			SELECT s.ticker, s.sentiment, COUNT(*) as count
			FROM twitter_crypto_signal s
			INNER JOIN twitter_crypto_tweets_foxhole t ON s.tweet_id = t.id
			WHERE t.author_username = $1 
			AND t.tweet_created_at >= $2
			AND s.ticker != ''
			AND s.ticker != 'NONE'
			AND s.sentiment IN ('Bearish', 'Bullish')
			GROUP BY s.ticker, s.sentiment
			ORDER BY s.ticker, s.sentiment;
		`
		args = []interface{}{authorUsername, since}
	case "1M":
		since = now.AddDate(0, -1, 0)
		query = `
			SELECT s.ticker, s.sentiment, COUNT(*) as count
			FROM twitter_crypto_signal s
			INNER JOIN twitter_crypto_tweets_foxhole t ON s.tweet_id = t.id
			WHERE t.author_username = $1 
			AND t.tweet_created_at >= $2
			AND s.ticker != ''
			AND s.ticker != 'NONE'
			AND s.sentiment IN ('Bearish', 'Bullish')
			GROUP BY s.ticker, s.sentiment
			ORDER BY s.ticker, s.sentiment;
		`
		args = []interface{}{authorUsername, since}
	case "ALL":
		query = `
			SELECT s.ticker, s.sentiment, COUNT(*) as count
			FROM twitter_crypto_signal s
			INNER JOIN twitter_crypto_tweets_foxhole t ON s.tweet_id = t.id
			WHERE t.author_username = $1
			AND s.ticker != ''
			AND s.ticker != 'NONE'
			AND s.sentiment IN ('Bearish', 'Bullish')
			GROUP BY s.ticker, s.sentiment
			ORDER BY s.ticker, s.sentiment;
		`
		args = []interface{}{authorUsername}
	default:
		return bearishTokens, bullishTokens, errors.New("invalid period")
	}

	fmt.Printf("Executing sentiment analysis query for: %s, period: %s\n", authorUsername, period)

	// Execute query to get sentiment data
	var sentimentData []struct {
		Ticker    string `db:"ticker"`
		Sentiment string `db:"sentiment"`
		Count     int    `db:"count"`
	}

	if err := r.postgresDB.SelectContext(ctx, &sentimentData, query, args...); err != nil {
		fmt.Printf("Error fetching sentiment data for username %s: %v\n", authorUsername, err)
		return bearishTokens, bullishTokens, fmt.Errorf("failed to get sentiment data for username %s: %w", authorUsername, err)
	}

	fmt.Printf("Found %d sentiment records for %s\n", len(sentimentData), authorUsername)

	// Create maps to aggregate sentiment data by ticker
	bearishMap := make(map[string]int)
	bullishMap := make(map[string]int)

	// Process and categorize sentiment data
	for _, data := range sentimentData {
		ticker := data.Ticker
		sentiment := data.Sentiment
		count := data.Count

		switch sentiment {
		case "Bearish":
			bearishMap[ticker] += count
		case "Bullish":
			bullishMap[ticker] += count
		}
	}

	// Convert maps to slices
	for ticker, count := range bearishMap {
		bearishTokens = append(bearishTokens, model.SentimentToken{
			Ticker:    ticker,
			Count:     count,
			Sentiment: "Bearish",
		})
	}

	for ticker, count := range bullishMap {
		bullishTokens = append(bullishTokens, model.SentimentToken{
			Ticker:    ticker,
			Count:     count,
			Sentiment: "Bullish",
		})
	}

	// Sort by count (descending)
	sort.Slice(bearishTokens, func(i, j int) bool {
		return bearishTokens[i].Count > bearishTokens[j].Count
	})

	sort.Slice(bullishTokens, func(i, j int) bool {
		return bullishTokens[i].Count > bullishTokens[j].Count
	})

	fmt.Printf("Categorized sentiment data - Bearish tokens: %d, Bullish tokens: %d\n", len(bearishTokens), len(bullishTokens))

	// Ensure we always return initialized slices (empty arrays instead of null)
	if bearishTokens == nil {
		bearishTokens = []model.SentimentToken{}
	}
	if bullishTokens == nil {
		bullishTokens = []model.SentimentToken{}
	}

	return bearishTokens, bullishTokens, nil
}

// buildNavQuery creates the appropriate SQL query based on the period
func (r *PerformanceRepo) buildNavQuery(period string, since *time.Time) (string, []interface{}, error) {
	now := time.Now()
	var query string
	var args []interface{}

	switch period {
	case "7D":
		*since = now.AddDate(0, 0, -7)
		query = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_nav
				WHERE datetime >= $1
			) t
			WHERE t.rn = 1
			ORDER BY t.author_username, t.day ASC;
		`
		args = []interface{}{*since}
	case "1M":
		*since = now.AddDate(0, -1, 0)
		query = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE(datetime) AS day,
					nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE(datetime) ORDER BY datetime DESC) AS rn
				FROM crypto_author_nav
				WHERE datetime >= $1
			) t
			WHERE t.rn = 1
			ORDER BY t.author_username, t.day ASC;
		`
		args = []interface{}{*since}
	case "ALL":
		query = `
			SELECT t.author_username, t.day, t.nav 
			FROM (
				SELECT 
					author_username,
					DATE_TRUNC('month', datetime) AS month,
					DATE(datetime) AS day,
					nav,
					ROW_NUMBER() OVER (PARTITION BY author_username, DATE_TRUNC('month', datetime) ORDER BY datetime ASC) AS rn
				FROM crypto_author_nav
			) t
			WHERE t.rn = 1
			ORDER BY t.author_username, t.day ASC;
		`
		args = []interface{}{}
	default:
		return "", nil, errors.New("invalid period")
	}

	return query, args, nil
}

// groupNavDataByAuthor organizes NAV data by author username
func (r *PerformanceRepo) groupNavDataByAuthor(rows []struct {
	AuthorUsername string    `db:"author_username"`
	Datetime       time.Time `db:"day"`
	Nav            float64   `db:"nav"`
}) map[string][]model.Nav {
	navMap := make(map[string][]model.Nav)
	for _, row := range rows {
		navMap[row.AuthorUsername] = append(navMap[row.AuthorUsername], model.Nav{
			Datetime: row.Datetime,
			Nav:      row.Nav,
		})
	}
	return navMap
}

// normalizeNavValues converts NAV values to percentage with base 100
func (r *PerformanceRepo) normalizeNavValues(navMap map[string][]model.Nav) {
	for author, navs := range navMap {
		if len(navs) == 0 || navs[0].Nav == 0 {
			// Prevent division by zero: set NAV = 100
			for i := range navs {
				navs[i].Nav = 100.0
			}
			navMap[author] = navs
			continue
		}

		firstNav := navs[0].Nav
		for i, n := range navs {
			if firstNav != 0 {
				navs[i].Nav = 100.0 * n.Nav / firstNav
				if math.IsNaN(navs[i].Nav) || math.IsInf(navs[i].Nav, 0) {
					navs[i].Nav = 100.0 // fallback safe value
				}
			} else {
				navs[i].Nav = 100.0
			}
		}
		navMap[author] = navs
	}
}

// buildAuthorNavData calculates performance metrics and builds AuthorNav objects
func (r *PerformanceRepo) buildAuthorNavData(navMap map[string][]model.Nav, authorNames map[string]string) []model.AuthorNav {
	var data []model.AuthorNav

	for author, navs := range navMap {
		if len(navs) == 0 {
			continue
		}

		// Calculate ROI
		roi := navs[len(navs)-1].Nav - 100
		if math.IsNaN(roi) || math.IsInf(roi, 0) {
			roi = 0
		}

		// Skip authors with non-positive ROI
		// if roi <= 0 {
		// 	continue
		// }

		// Calculate drawdowns
		_, maxDrawdown, currentDrawdown := r.calculateDrawdowns(navs)

		// Create AuthorNav object
		authorName := authorNames[author]
		data = append(data, model.AuthorNav{
			AuthorUsername:  author,
			AuthorName:      authorName,
			WeightNav:       navs,
			ROI:             roi,
			StartNav:        navs[0].Nav,
			EndNav:          navs[len(navs)-1].Nav,
			Drawdown:        currentDrawdown,
			MaximumDrawdown: maxDrawdown,
		})
	}

	return data
}

// calculateDrawdowns computes maximum and current drawdowns
func (r *PerformanceRepo) calculateDrawdowns(navs []model.Nav) (peak, maxDrawdown, currentDrawdown float64) {
	if len(navs) == 0 {
		return 0, 0, 0
	}

	peak = navs[0].Nav
	for _, nav := range navs {
		if nav.Nav > peak {
			peak = nav.Nav
		}
		drawdown := (peak - nav.Nav) / peak
		if math.IsNaN(drawdown) || math.IsInf(drawdown, 0) {
			drawdown = 0
		}
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	currentDrawdown = (peak - navs[len(navs)-1].Nav) / peak
	if math.IsNaN(currentDrawdown) || math.IsInf(currentDrawdown, 0) {
		currentDrawdown = 0
	}
	if math.IsNaN(maxDrawdown) || math.IsInf(maxDrawdown, 0) {
		maxDrawdown = 0
	}

	return peak, maxDrawdown, currentDrawdown
}

// sortAndRankAuthors sorts authors by ROI and assigns ranks
func (r *PerformanceRepo) sortAndRankAuthors(data *[]model.AuthorNav) {
	sort.Slice(*data, func(i, j int) bool {
		return (*data)[i].ROI > (*data)[j].ROI
	})

	for i := range *data {
		(*data)[i].Rank = i + 1
	}

	// Limit to top 50 authors
	if len(*data) > 100 {
		*data = (*data)[:100]
	}
}
