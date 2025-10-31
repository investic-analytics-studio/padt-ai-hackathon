package domain

import "time"

type CryptoUserRefcode struct {
	CryptoUserID    string    `json:"crypto_user_id" db:"crypto_user_id"`
	Refcode         string    `json:"refcode" db:"refcode"`
	CryptoRefUserID *string   `json:"crypto_ref_user_id,omitempty" db:"crypto_ref_user_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
