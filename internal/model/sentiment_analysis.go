package model

type OverallSentimentQueryResult struct {
	Sentiment string `json:"sentiment" db:"sentiment"`
	Count     int    `json:"count" db:"count"`
}

type OverallSentimentResult struct {
	PositiveCount int `json:"positive_count"`
	NegativeCount int `json:"negative_count"`
	NeutralCount  int `json:"neutral_count"`
}

type TickerSentimentQueryResult struct {
	Ticker    string `json:"ticker"`
	Sentiment string `json:"sentiment"`
	Count     int    `json:"count"`
}

type TokenSentiment struct {
	Ticker        string `json:"ticker"`
	PositiveCount int    `json:"positive_count"`
	NegativeCount int    `json:"negative_count"`
	NeutralCount  int    `json:"neutral_count"`
}

type SentimentAnalysisModel struct {
	OverallSentimentAnalysis OverallSentimentResult           `json:"overall_sentiment"`
	TokenSentimentAnalysis   []TokenSentiment                 `json:"token_sentiment"`
	TwitterIdResult          []TickerSentimentTwitterIdResult `json:"twitter_id_result"`
}

type SentimentAnalysisByAuthorListRequest struct {
	AuthorList []string `json:"author_list"`
	DateRange  string   `json:"date_range"`
}

type SentimentAnalysisByTierRequest struct {
	Tier      string `json:"tier"`
	DateRange string `json:"date_range"`
}

type TwitterContent struct {
	Text         string `json:"text"`
	MediaURL     string `json:"media_url"`
	TweetCreate  string `json:"tweet_created_at"`
	LikeCount    int    `json:"like_count"`
	RetweetCount int    `json:"retweet_count"`
}

type TwitterIdResult struct {
	Ticker    string `json:"ticker" db:"ticker"`
	Sentiment string `json:"sentiment" db:"sentiment"`
	TweetIds  string `json:"tweet_ids" db:"tweet_ids"`
}

type TickerSentimentTwitterIdResult struct {
	Ticker    string   `json:"ticker" db:"ticker"`
	Sentiment string   `json:"sentiment" db:"sentiment"`
	TweetIds  []string `json:"tweet_ids" db:"tweet_ids"`
}
