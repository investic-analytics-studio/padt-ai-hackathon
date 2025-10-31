package service

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type CryptoNotificationService struct {
	repo port.NotificationRepo
}

func NewCryptoNotificationService(repo port.NotificationRepo) *CryptoNotificationService {
	return &CryptoNotificationService{repo: repo}
}

func (s *CryptoNotificationService) GetNotificationGroupList(ctx context.Context, uid string) ([]model.NotificationGroupList, error) {
	return s.repo.GetNotificationGroupList(ctx, uid)
}

func (s *CryptoNotificationService) UpdateGroupName(ctx context.Context, groupID, newGroupName string) error {
	return s.repo.UpdateGroupName(ctx, groupID, newGroupName)
}
func (s *CryptoNotificationService) AddGroup(ctx context.Context, uid, groupName string) error {
	return s.repo.AddGroup(ctx, uid, groupName)
}
func (s *CryptoNotificationService) CountGroup(ctx context.Context, uid string) (int, error) {
	return s.repo.CountGroup(ctx, uid)
}
func (s *CryptoNotificationService) AddAuthor(ctx context.Context, uid, groupID, authorName string) error {
	return s.repo.AddAuthor(ctx, uid, groupID, authorName)
}
func (s *CryptoNotificationService) RemoveAuthor(ctx context.Context, uid, groupAutherID string) error {
	return s.repo.RemoveAuthor(ctx, uid, groupAutherID)
}
func (s *CryptoNotificationService) UpdateTelegram(ctx context.Context, uid, chatID, userID string) error {
	return s.repo.UpdateTelegram(ctx, uid, chatID, userID)
}
func (s *CryptoNotificationService) GetTelegram(ctx context.Context, uid string) (model.Telegram, error) {
	return s.repo.GetTelegram(ctx, uid)
}
func (s *CryptoNotificationService) DeleteGroupAuthor(ctx context.Context, uid, groupID string) error {
	return s.repo.DeleteGroupAuthor(ctx, uid, groupID)
}
func (s *CryptoNotificationService) DisconnectNotification(ctx context.Context, uid string) error {
	return s.repo.DisconnectNotification(ctx, uid)
}
