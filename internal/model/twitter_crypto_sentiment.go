package model

import (
	"encoding/json"
	"math"
	"time"
)

type TwitterCryptoSentiment struct {
	ID                string    `json:"id" db:"id"`
	Text              string    `json:"text" db:"text"`
	AuthorUsername    string    `json:"author_username" db:"author_username"`
	ResponseScore     *int      `json:"response_score" db:"response_score"`
	ResponseSentiment string    `json:"response_sentiment" db:"response_sentiment"`
	ResponseTickers   string    `json:"response_tickers" db:"response_tickers"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type TwitterCryptoTweet struct {
	ID               string          `json:"id" db:"id"`
	URL              string          `json:"url" db:"url"`
	Text             string          `json:"text" db:"text"`
	Source           string          `json:"source" db:"source"`
	RetweetCount     int             `json:"retweetcount" db:"retweetcount"`
	ReplyCount       int             `json:"replycount" db:"replycount"`
	LikeCount        int             `json:"likecount" db:"likecount"`
	QuoteCount       int             `json:"quotecount" db:"quotecount"`
	ViewCount        NullableFloat64 `json:"viewcount" db:"viewcount"`
	TweetCreatedAt   time.Time       `json:"tweet_created_at" db:"tweet_created_at"`
	BookmarkCount    int             `json:"bookmarkcount" db:"bookmarkcount"`
	IsReply          bool            `json:"isreply" db:"isreply"`
	Conversation     string          `json:"conversation" db:"conversationid"`
	IsPinned         bool            `json:"ispinned" db:"ispinned"`
	IsRetweet        bool            `json:"isretweet" db:"isretweet"`
	IsQuote          bool            `json:"isquote" db:"isquote"`
	MediaURL         string          `json:"media_url" db:"media_url"`
	TweetCreatedDate time.Time       `json:"tweet_created_date" db:"tweet_created_date"`
	AuthorUsername   string          `json:"author_username" db:"author_username"`
	TickersRuleBased string          `json:"tickers_rule_based" db:"tickers_rule_based"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

type SignalActionType string

const (
	SignalActionLong    SignalActionType = "LONG"
	SignalActionShort   SignalActionType = "SHORT"
	SignalActionNatural SignalActionType = "NONE"
)

type TwitterCryptoSignal struct {
	SignalCreatedAt *time.Time       `json:"signal_created_at" db:"signal_created_at"`
	SignalUpdatedAt *time.Time       `json:"signal_updated_at" db:"signal_updated_at"`
	SignalContent   *string          `json:"signal_content" db:"signal_content"`
	SignalTicker    *string          `json:"signal_ticker" db:"signal_ticker"`
	SignalAction    SignalActionType `json:"signal_action" db:"signal_action"`
	SignalScore     *float64         `json:"signal_score" db:"signal_score"`
	SignalSentiment *string          `json:"signal_sentiment" db:"signal_sentiment"`
}

type NullableFloat64 float64

func (f NullableFloat64) MarshalJSON() ([]byte, error) {
	if math.IsNaN(float64(f)) {
		return []byte("null"), nil
	}
	return json.Marshal(float64(f))
}

type TwitterCryptoAuthorProfile struct {
	AuthorUsername   string    `json:"author_username" db:"author_username"`
	AuthorURL        string    `json:"author_url" db:"author_url"`
	AuthorTwitterURL string    `json:"author_twitterurl" db:"author_twitterurl"`
	AuthorName       string    `json:"author_name" db:"author_name"`
	AuthorFollowers  int       `json:"author_followers" db:"author_followers"`
	AuthorFollowing  int       `json:"author_following" db:"author_following"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type TwitterCryptoTweetWithSentiment struct {
	TwitterCryptoTweet
	ResponseScore     *int   `json:"response_score,omitempty" db:"response_score"`
	ResponseSentiment string `json:"response_sentiment,omitempty" db:"response_sentiment"`
	ResponseTickers   string `json:"response_tickers,omitempty" db:"response_tickers"`
}

type TwitterCryptoTweetWithSentimentAndAuthor struct {
	TwitterCryptoTweet
	ResponseScore     *int   `json:"response_score,omitempty" db:"response_score"`
	ResponseSentiment string `json:"response_sentiment,omitempty" db:"response_sentiment"`
	ResponseTickers   string `json:"response_tickers,omitempty" db:"response_tickers"`
	AuthorURL         string `json:"author_url" db:"author_url"`
	AuthorTwitterURL  string `json:"author_twitter_url" db:"author_twitterurl"`
	AuthorName        string `json:"author_name" db:"author_name"`
	AuthorFollowers   int    `json:"author_followers" db:"author_followers"`
	AuthorFollowing   int    `json:"author_following" db:"author_following"`
}

type TwitterCryptoTweetWithSentimentAndTier struct {
	ResponseScore     *int      `json:"response_score,omitempty" db:"response_score"`
	ResponseSentiment string    `json:"response_sentiment,omitempty" db:"response_sentiment"`
	ResponseTickers   string    `json:"response_tickers,omitempty" db:"response_tickers"`
	TweetCreatedAt    time.Time `json:"tweet_created_at" db:"tweet_created_at"`
	TweetCreatedDate  time.Time `json:"tweet_created_date" db:"tweet_created_date"`
	TwitterCryptoSignal
}

type TwitterCryptoAuthorWinRate struct {
	AuthorUsername   string    `json:"author_username" db:"author_username"`
	AuthorURL        string    `json:"author_url" db:"author_url"`
	AuthorTwitterURL string    `json:"author_twitterurl" db:"author_twitterurl"`
	AuthorName       string    `json:"author_name" db:"author_name"`
	AuthorFollowers  int       `json:"author_followers" db:"author_followers"`
	AuthorFollowing  int       `json:"author_following" db:"author_following"`
	WinRate          float32   `json:"winrate" db:"winrate"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type TwitterCryptoTweetWithSentimentAuthorAndSignal struct {
	TwitterCryptoTweet
	ResponseScore     *float64 `json:"response_score" db:"response_score"`
	ResponseSentiment string   `json:"response_sentiment" db:"response_sentiment"`
	ResponseTickers   string   `json:"response_tickers" db:"response_tickers"`
	AuthorURL         string   `json:"author_url" db:"author_url"`
	AuthorTwitterURL  string   `json:"author_twitter_url" db:"author_twitterurl"`
	AuthorName        string   `json:"author_name" db:"author_name"`
	AuthorFollowers   int      `json:"author_followers" db:"author_followers"`
	AuthorFollowing   int      `json:"author_following" db:"author_following"`
	TweetID           string   `json:"tweet_id" db:"tweet_id"`
	TwitterCryptoSignal
}

type TokenMentionSymbolWithAuthorResponse struct {
	Tokens     []TokenMentionSymbolAndAuthor `json:"tokens" db:"tokens"`
	NextCursor NextCursor                    `json:"next_cursor" db:"next_cursor"`
	TotalCount int                           `json:"total_count" db:"total_count"`
}

type TokenMentionSymbolAndAuthor struct {
	Symbol string `json:"symbol" db:"symbol"`
}

type NextCursor struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ID        string    `json:"id" db:"id"`
}
