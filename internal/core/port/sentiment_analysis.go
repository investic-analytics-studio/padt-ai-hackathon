package port

import "github.com/quantsmithapp/datastation-backend/internal/model"

type SentimentAnalysisRepo interface {
	GetSentimentAnalysisByAuthorList(authorList []string, dateRange string) (model.SentimentAnalysisModel, error)
	GetSentimentAnalysisByTier(tier string, dateRange string) (model.SentimentAnalysisModel, error)
}

type SentimentAnalysisService interface {
	GetSentimentAnalysisByAuthorList(authorList []string, dateRange string) (model.SentimentAnalysisModel, error)
	GetSentimentAnalysisByTier(tier string, dateRange string) (model.SentimentAnalysisModel, error)
}
