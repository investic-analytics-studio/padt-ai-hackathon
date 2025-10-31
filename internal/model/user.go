package model

import "time"

type User struct {
	UID       string     `json:"uid"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdateAt  *time.Time `json:"last_update"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type CheckXUser struct {
	IsUserExit         bool `db:"is_user_exit"`
	IsUserTelegramExit bool `db:"is_user_telegram_exit"`
}
type XUser struct {
	TwitterName string `json:"twitter_name" validate:"required" example:"NoMoonNoBuy"`
}

func (User) TableName() string {
	return "crypto_users"
}
