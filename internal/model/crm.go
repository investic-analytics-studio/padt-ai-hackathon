package model

import "time"

type CRMLoginBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CryptoUser struct {
	UUID                 string     `gorm:"column:uuid;type:varchar(255);not null" json:"uuid" db:"uuid"`
	Email                *string    `gorm:"column:email;type:varchar(50)" json:"email" db:"email"`
	LastUpdate           string     `gorm:"column:last_update;type:varchar(50);not null" json:"last_update" db:"last_update"`
	CreatedAt            time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at" db:"updated_at"`
	BindingWallet        *string    `gorm:"column:binding_wallet;type:varchar(255)" json:"binding_wallet" db:"binding_wallet"`
	TwitterUID           *string    `gorm:"column:twitter_uid;type:varchar(255)" json:"twitter_uid" db:"twitter_uid"`
	TwitterName          *string    `gorm:"column:twitter_name;type:varchar(255)" json:"twitter_name" db:"twitter_name"`
	ID                   string     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id" db:"id"`
	Address              *string    `gorm:"column:address;type:varchar(255)" json:"address" db:"address"`
	HashedRT             *string    `gorm:"column:hashed_rt;type:varchar(255)" json:"hashed_rt" db:"hashed_rt"`
	MethodLogin          string     `gorm:"column:method_login;type:varchar(20);not null" json:"method_login" db:"method_login"`
	StripeCustomerID     *string    `gorm:"column:stripe_customer_id;type:varchar(64)" json:"stripe_customer_id" db:"stripe_customer_id"`
	TelegramChatID       *string    `gorm:"column:telegram_chat_id;type:varchar(20)" json:"telegram_chat_id" db:"telegram_chat_id"`
	TelegramUserID       *string    `gorm:"column:telegram_user_id;type:varchar(20)" json:"telegram_user_id" db:"telegram_user_id"`
	IsCopytradeApproved  bool       `gorm:"column:is_copytrade_approved;type:boolean" json:"is_copytrade_approved" db:"is_copytrade_approved"`
	WaitinglistTimestamp *time.Time `gorm:"column:waiting_list_timestamp;type:timestamp" json:"waiting_list_timestamp" db:"waiting_list_timestamp"`
}

// CryptoReferralScore represents a single referral score entry.
// CryptoUserID was re-added as it's crucial for grouping.
type CryptoReferralScore struct {
	Date           time.Time `gorm:"column:date;type:timestamp;not null" json:"date" db:"date"`
	DirectPoints   *int      `gorm:"column:direct_points;type:int4" json:"direct_points" db:"direct_points"`
	IndirectPoints *int      `gorm:"column:indirect_points;type:int4" json:"indirect_points" db:"indirect_points"`
	TotalPoints    *int      `gorm:"column:total_points;type:int4" json:"total_points" db:"total_points"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at" db:"updated_at"`
}

// UserWithReferralScores represents a user along with all their referral scores.
// This replaces the previous CryptoUserReferralScore which had an embedded score.
type UserWithReferralScores struct {
	UserUUID    string                `json:"user_uuid" db:"uuid"`
	Email       *string               `json:"email" db:"email"`
	TwitterName *string               `json:"twitter_name" db:"twitter_name"`
	Scores      []CryptoReferralScore `json:"scores"`
}

// FlatReferralData is a temporary struct to help scan the flat results from the SQL JOIN query
// before grouping them into the UserWithReferralScores structure.
// The db tags match the column names (or aliases) from the SQL query.
type FlatReferralData struct {
	UserUUID       string    `db:"uuid"`            // from crypto_user table (cu.uuid)
	Email          *string   `db:"email"`           // from crypto_user table (cu.email)
	TwitterName    *string   `db:"twitter_name"`    // from crypto_user table (cu.twitter_name)
	CryptoUserID   string    `db:"crypto_user_id"`  // from crypto_refferal_score table
	Date           time.Time `db:"date"`            // from crypto_refferal_score table
	DirectPoints   *int      `db:"direct_points"`   // from crypto_refferal_score table
	IndirectPoints *int      `db:"indirect_points"` // from crypto_refferal_score table
	TotalPoints    *int      `db:"total_points"`    // from crypto_refferal_score table
	CreatedAt      time.Time `db:"created_at"`      // from crypto_refferal_score table
	UpdatedAt      time.Time `db:"updated_at"`      // from crypto_refferal_score table
}
type FlatUserWithRefCode struct {
	UserUUID    string  `db:"uuid"`
	Email       *string `db:"email"`
	TwitterName *string `db:"twitter_name"`
	RefCode     *string `json:"refcode"`
	RefUser     *string `json:"refuser"`
	KolCode     *string `json:"kolcode"`
	KolUser     *string `json:"koluser"`
}
type UserWithRefCode struct {
	UserUUID    string    `db:"uuid"`
	Email       *string   `db:"email"`
	TwitterName *string   `db:"twitter_name"`
	RefCode     []RefCode `json:"refcode"`
	KolCode     []KolCode `json:"kolcode"`
}
type RefCode struct {
	RefCode       *string `db:"refcode"`
	CryptoRefUser *string `db:"crypto_ref_user"`
}
type KolCode struct {
	KolCode       *string `db:"display_code"`
	CryptoRefUser *string `db:"crypto_ref_user"`
}

// TradeLog represents a row from table trade_logs
type TradeLog struct {
	ID         int64     `db:"id" json:"id"`
	Source     string    `db:"source" json:"source"`
	AccountID  string    `db:"account_id" json:"account_id"`
	WalletID   *string   `db:"wallet_id" json:"wallet_id"`
	Symbol     string    `db:"symbol" json:"symbol"`
	Side       string    `db:"side" json:"side"`
	BaseSize   *float64  `db:"base_size" json:"base_size"`
	UsdcValue  *float64  `db:"usdc_value" json:"usdc_value"`
	Price      *float64  `db:"price" json:"price"`
	Leverage   *int      `db:"leverage" json:"leverage"`
	Event      string    `db:"event" json:"event"`
	Status     string    `db:"status" json:"status"`
	ExecutedAt time.Time `db:"executed_at" json:"executed_at"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// PrivyUserOverview aggregates user profile, wallets and trade logs for CRM view
type PrivyUserOverview struct {
	UUID        string       `json:"uuid"`
	Email       *string      `json:"email"`
	TwitterName *string      `json:"twitter_name"`
	WalletCount int          `json:"wallet_count"`
	Wallets     []WalletInfo `json:"wallets"`
	TradeLogs   []TradeLog   `json:"trade_logs"`
}
