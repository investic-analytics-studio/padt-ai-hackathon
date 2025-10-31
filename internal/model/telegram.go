package model

type Telegram struct {
	ChatID string `db:"telegram_chat_id" json:"telegram_chat_id"`
	UserID string `db:"telegram_user_id" json:"telegram_user_id"`
}
