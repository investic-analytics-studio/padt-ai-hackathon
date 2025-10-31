package service

import (
	"strings"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type ExtractTagService struct {
	repo port.ExtractTagRepo
}

func NewExtractTagService(repo port.ExtractTagRepo) *ExtractTagService {
	return &ExtractTagService{repo: repo}
}

func (s *ExtractTagService) GetAllTags() (model.ExtractTagResponse, error) {
	tags, err := s.repo.GetAllTags()
	if err != nil {
		return model.ExtractTagResponse{}, err
	}

	coinWithTags := make(map[string][]string)
	for _, tag := range tags {
		coinWithTags[tag.BaseAsset] = getTags(tag.BinanceTag)
	}

	return model.ExtractTagResponse{CoinWithTags: coinWithTags}, nil
}

func (s *ExtractTagService) GetUniqueTags() ([]string, error) {
	rawTags, err := s.repo.GetUniqueTags()
	if err != nil {
		return nil, err
	}

	uniqueTags := make(map[string]struct{})
	for _, tag := range rawTags {
		upperTag := strings.ToUpper(tag)
		uniqueTags[upperTag] = struct{}{}
	}

	var uniqueTagsList []string
	for tag := range uniqueTags {
		uniqueTagsList = append(uniqueTagsList, tag)
	}

	return uniqueTagsList, nil
}

func getTags(tagString string) []string {
	// Remove the curly braces
	tagString = strings.Trim(tagString, "{}")
	// Split the string by commas
	tags := strings.Split(tagString, ",")
	// uppercase the tags
	for i, tag := range tags {
		tags[i] = strings.ToUpper(tag)
	}
	return tags
}
