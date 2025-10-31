package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type NotificationRepo interface {
	GetNotificationGroupList(ctx context.Context, uid string) ([]model.NotificationGroupList, error)
	UpdateGroupName(ctx context.Context, groupID, newGroupName string) error
	AddGroup(ctx context.Context, uid, GroupName string) error
	CountGroup(ctx context.Context, uid string) (int, error)
	AddAuthor(ctx context.Context, uid, groupID, authorName string) error
	RemoveAuthor(ctx context.Context, uid, GroupAutherID string) error
	UpdateTelegram(ctx context.Context, uid, chatID, userID string) error
	GetTelegram(ctx context.Context, uid string) (model.Telegram, error)
	DeleteGroupAuthor(ctx context.Context, uid, groupID string) error
	DisconnectNotification(ctx context.Context, uid string) error
}

type NotificationService interface {
	GetNotificationGroupList(ctx context.Context, uid string) ([]model.NotificationGroupList, error)
	UpdateGroupName(ctx context.Context, groupID, newGroupName string) error
	AddGroup(ctx context.Context, uid, GroupName string) error
	CountGroup(ctx context.Context, uid string) (int, error)
	AddAuthor(ctx context.Context, uid, groupID, authorName string) error
	RemoveAuthor(ctx context.Context, uid, GroupAutherID string) error
	UpdateTelegram(ctx context.Context, uid, chatID, userID string) error
	GetTelegram(ctx context.Context, uid string) (model.Telegram, error)
	DeleteGroupAuthor(ctx context.Context, uid, groupID string) error
	DisconnectNotification(ctx context.Context, uid string) error
}
