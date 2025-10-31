package model

import (
	"time"
)

type AuthorNav struct {
	AuthorUsername  string  `json:"authorUsername"`
	AuthorName      string  `json:"authorName"`
	WeightNav       []Nav   `json:"WeightNav"`
	ROI             float64 `json:"roi"`
	Rank            int     `json:"rank"`
	StartNav        float64 `json:"startNav"`
	EndNav          float64 `json:"endNav"`
	Drawdown        float64 `json:"drawdown"`
	MaximumDrawdown float64 `json:"maximumDrawdown"`
}
type Nav struct {
	Datetime time.Time `json:"datetime"`
	Nav      float64   `json:"nav"`
}

type AuthorProfile struct {
	AuthorUsername   string    `json:"author_username" db:"author_username"`
	AuthorURL        string    `json:"author_url" db:"author_url"`
	AuthorTwitterURL string    `json:"author_twitter_url" db:"author_twitterurl"`
	AuthorName       string    `json:"author_name" db:"author_name"`
	AuthorFollowers  int       `json:"author_followers" db:"author_followers"`
	AuthorFollowing  int       `json:"author_following" db:"author_following"`
	AuthorTier       string    `json:"author_tier" db:"author_tier"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type AuthorTweet struct {
	ID               string    `json:"id" db:"id"`
	URL              string    `json:"url" db:"url"`
	Text             string    `json:"text" db:"text"`
	Source           string    `json:"source" db:"source"`
	RetweetCount     int       `json:"retweet_count" db:"retweetcount"`
	ReplyCount       int       `json:"reply_count" db:"replycount"`
	LikeCount        int       `json:"like_count" db:"likecount"`
	QuoteCount       int       `json:"quote_count" db:"quotecount"`
	ViewCount        *float64  `json:"view_count" db:"viewcount"`
	TweetCreatedAt   time.Time `json:"tweet_created_at" db:"tweet_created_at"`
	BookmarkCount    int       `json:"bookmark_count" db:"bookmarkcount"`
	IsReply          bool      `json:"is_reply" db:"isreply"`
	ConversationID   string    `json:"conversation_id" db:"conversationid"`
	IsPinned         bool      `json:"is_pinned" db:"ispinned"`
	IsRetweet        bool      `json:"is_retweet" db:"isretweet"`
	IsQuote          bool      `json:"is_quote" db:"isquote"`
	MediaURL         string    `json:"media_url" db:"media_url"`
	TweetCreatedDate time.Time `json:"tweet_created_date" db:"tweet_created_date"`
	TickersRuleBased string    `json:"tickers_rule_based" db:"tickers_rule_based"`
}

type AuthorSignal struct {
	TweetID       string    `json:"tweet_id" db:"tweet_id"`
	Content       string    `json:"content" db:"content"`
	Ticker        string    `json:"ticker" db:"ticker"`
	Action        string    `json:"action" db:"action"`
	Score         int       `json:"score" db:"score"`
	Sentiment     string    `json:"sentiment" db:"sentiment"`
	PromptVersion string    `json:"prompt_version" db:"prompt_version"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// AuthorTweetWithSignals represents a tweet with its associated signals
type AuthorTweetWithSignals struct {
	// Tweet data
	ID               string    `json:"id" db:"id"`
	URL              string    `json:"url" db:"url"`
	Text             string    `json:"text" db:"text"`
	Source           string    `json:"source" db:"source"`
	RetweetCount     int       `json:"retweet_count" db:"retweetcount"`
	ReplyCount       int       `json:"reply_count" db:"replycount"`
	LikeCount        int       `json:"like_count" db:"likecount"`
	QuoteCount       int       `json:"quote_count" db:"quotecount"`
	ViewCount        *float64  `json:"view_count" db:"viewcount"`
	TweetCreatedAt   time.Time `json:"tweet_created_at" db:"tweet_created_at"`
	BookmarkCount    int       `json:"bookmark_count" db:"bookmarkcount"`
	IsReply          bool      `json:"is_reply" db:"isreply"`
	ConversationID   string    `json:"conversation_id" db:"conversationid"`
	IsPinned         bool      `json:"is_pinned" db:"ispinned"`
	IsRetweet        bool      `json:"is_retweet" db:"isretweet"`
	IsQuote          bool      `json:"is_quote" db:"isquote"`
	MediaURL         string    `json:"media_url" db:"media_url"`
	TweetCreatedDate time.Time `json:"tweet_created_date" db:"tweet_created_date"`
	TickersRuleBased string    `json:"tickers_rule_based" db:"tickers_rule_based"`

	// Signal data (optional - may be null if no signal exists)
	SignalTicker        *string    `json:"signal_ticker,omitempty" db:"signal_ticker"`
	SignalAction        *string    `json:"signal_action,omitempty" db:"signal_action"`
	SignalScore         *int       `json:"signal_score,omitempty" db:"signal_score"`
	SignalSentiment     *string    `json:"signal_sentiment,omitempty" db:"signal_sentiment"`
	SignalPromptVersion *string    `json:"signal_prompt_version,omitempty" db:"signal_prompt_version"`
	SignalCreatedAt     *time.Time `json:"signal_created_at,omitempty" db:"signal_created_at"`
	SignalUpdatedAt     *time.Time `json:"signal_updated_at,omitempty" db:"signal_updated_at"`
}

// SentimentToken represents sentiment analysis data for a specific token in performance context
type SentimentToken struct {
	Ticker    string `json:"ticker"`
	Count     int    `json:"count"`
	Sentiment string `json:"sentiment"`
}

type AuthorDetail struct {
	AuthorName      string                   `json:"authorName"`
	WeightNav       []Nav                    `json:"WeightNav"`
	ROI             float64                  `json:"roi"`
	StartNav        float64                  `json:"startNav"`
	EndNav          float64                  `json:"endNav"`
	Drawdown        float64                  `json:"drawdown"`
	MaximumDrawdown float64                  `json:"maximumDrawdown"`
	Profile         AuthorProfile            `json:"profile"`
	RecentTimeline  []AuthorTweetWithSignals `json:"recentTimeline"`
	TotalTimeline   int                      `json:"totalTimeline"`
	Start           int                      `json:"start"`
	End             int                      `json:"end"`
	Limit           int                      `json:"limit"`
	BearishTokens   []SentimentToken         `json:"bearishTokens"`
	BullishTokens   []SentimentToken         `json:"bullishTokens"`
}
