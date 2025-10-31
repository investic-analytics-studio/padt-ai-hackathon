package model

import (
	"time"
)

type MarketOverviewTable struct {
	Ticker          string  `db:"ticker" json:"ticker"`
	LongCount       int     `db:"long_count" json:"long_count"`
	ShortCount      int     `db:"short_count" json:"short_count"`
	TotalSignals    int     `db:"total_signals" json:"total_signals"`
	LongShortRatio  float64 `db:"long_short_ratio" json:"long_short_ratio"`
	LongShortString string  `db:"long_short_str" json:"long_short_str"`
	Sentiment       string  `db:"sentiment" json:"sentiment"`
}

type MarketOverviewTableResponse struct {
	MarketOverviewTable
	AlphaSentiment string `json:"alpha_sentiment,omitempty"`
}

type MarketOverview struct {
	TotalTweet     int    `db:"total_tweet" json:"total_tweet"`
	BullishTweet   int    `db:"bullish_tweet" json:"bullish_tweet"`
	BearishTweet   int    `db:"bearish_tweet" json:"bearish_tweet"`
	Bias           string `db:"bias" json:"bias"`
	BiasPercentage int    `db:"bias_percentage" json:"bias_percentage"`
}

type EnrichedMarketOverview struct {
	Ticker            string   `json:"ticker"`
	LongCount         int      `json:"long_count"`
	ShortCount        int      `json:"short_count"`
	TotalSignals      int      `json:"total_signals"`
	LongShortRatio    float64  `json:"long_short_ratio"`
	LongShortString   string   `json:"long_short_str"`
	XSentiment        string   `json:"x_sentiment"`
	AlphaSentiment    string   `json:"alpha_sentiment"`
	PriceClose        *float64 `json:"price_close,omitempty"`
	PriceChange24h    *float64 `json:"price_change_24h,omitempty"`
	MarketCap         *float64 `json:"market_cap,omitempty"`
	Volume24h         *float64 `json:"volume_24h,omitempty"`
	NewsOverviewRatio *float32 `json:"news_overview_ratio,omitempty"`
	NewsSentiment     *string  `json:"news_sentiment,omitempty"`
}

type TokenDetailMarketOverview struct {
	ID                string          `json:"id" db:"id"`
	Text              string          `json:"text" db:"text"`
	URL               string          `json:"url" db:"url"`
	Source            string          `json:"source" db:"source"`
	RetweetCount      int             `json:"retweetcount" db:"retweetcount"`
	ReplyCount        int             `json:"replycount" db:"replycount"`
	LikeCount         int             `json:"likecount" db:"likecount"`
	QuoteCount        int             `json:"quotecount" db:"quotecount"`
	ViewCount         NullableFloat64 `json:"viewcount" db:"viewcount"`
	TweetCreatedAt    time.Time       `json:"tweet_created_at" db:"tweet_created_at"`
	BookmarkCount     int             `json:"bookmarkcount" db:"bookmarkcount"`
	IsReply           bool            `json:"isreply" db:"isreply"`
	Conversation      string          `json:"conversation" db:"conversationid"`
	IsPinned          bool            `json:"ispinned" db:"ispinned"`
	IsRetweet         bool            `json:"isretweet" db:"isretweet"`
	IsQuote           bool            `json:"isquote" db:"isquote"`
	MediaURL          string          `json:"media_url" db:"media_url"`
	TweetCreatedDate  time.Time       `json:"tweet_created_date" db:"tweet_created_date"`
	AuthorUsername    string          `json:"author_username" db:"author_username"`
	TickersRuleBased  string          `json:"tickers_rule_based" db:"tickers_rule_based"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	ResponseScore     float64         `json:"response_score" db:"response_score"`
	ResponseSentiment string          `json:"response_sentiment" db:"response_sentiment"`
	ResponseTickers   string          `json:"response_tickers" db:"response_tickers"`
	AuthorURL         string          `json:"author_url" db:"author_url"`
	AuthorTwitterURL  string          `json:"author_twitter_url" db:"author_twitterurl"`
	AuthorName        string          `json:"author_name" db:"author_name"`
	AuthorFollowers   int             `json:"author_followers" db:"author_followers"`
	AuthorFollowing   int             `json:"author_following" db:"author_following"`
	SignalCreatedAt   time.Time       `json:"signal_created_at" db:"signal_created_at"`
	SignalUpdatedAt   time.Time       `json:"signal_updated_at" db:"signal_updated_at"`
	SignalContent     string          `json:"signal_content" db:"signal_content"`
	SignalTicker      string          `json:"signal_ticker" db:"signal_ticker"`
	SignalAction      string          `json:"signal_action" db:"signal_action"`
	SignalScore       float64         `json:"signal_score" db:"signal_score"`
	SignalSentiment   string          `json:"signal_sentiment" db:"signal_sentiment"`
	TweetID           string          `json:"tweet_id" db:"tweet_id"`
}
