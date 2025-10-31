package model

import "time"

type Author struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
type NotificationGroupList struct {
	GroupID     string    `json:"group_id" db:"group_id"`
	GroupName   string    `json:"group_name" db:"group_name"`
	AuthorNames []Author  `json:"author_names" db:"author_names"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}
