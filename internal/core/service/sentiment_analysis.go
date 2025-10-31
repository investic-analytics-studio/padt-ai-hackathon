package service

import (
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type sentimentAnalysisService struct {
	repo port.SentimentAnalysisRepo
}

func NewSentimentAnalysisService(repo port.SentimentAnalysisRepo) port.SentimentAnalysisService {
	return &sentimentAnalysisService{repo: repo}
}

func (s *sentimentAnalysisService) GetSentimentAnalysisByAuthorList(authorList []string, dateRange string) (model.SentimentAnalysisModel, error) {
	return s.repo.GetSentimentAnalysisByAuthorList(authorList, dateRange)
}

func (s *sentimentAnalysisService) GetSentimentAnalysisByTier(tier string, dateRange string) (model.SentimentAnalysisModel, error) {
	return s.repo.GetSentimentAnalysisByTier(tier, dateRange)
}
