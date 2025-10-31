package port

import "context"

// TelegramService defines the interface for Telegram bot operations
type TelegramService interface {
	// Start starts the Telegram bot
	Start(ctx context.Context) error
	// SendMessage sends a message to a specific chat
	SendMessage(ctx context.Context, text string, chatID int64) error
}
