package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type UserBasicInfo struct {
	UUID        string
	Email       string
	TwitterName string
}
type CRMRepo interface {
	CheckCRMUser(login model.CRMLoginBody) (string, error)
	CheckCRMUserIsExit(ctx context.Context, uid string) (bool, error)
	ValidateDisplaycode(ctx context.Context, displayCode string) (bool, error)
	CheckKOLUserIsExit(ctx context.Context, uid string) (bool, error)
	InsertKolUser(ctx context.Context, kolUser model.KOLUser) (string, error)
	IsDisplayCodeExists(ctx context.Context, displayCode, excludeUserID string) (bool, error)
	UpdateDisplayCode(ctx context.Context, cryptoUserID, displayCode string) error
	IsKolUserExists(ctx context.Context, cryptoUserID string) (bool, error)
	GetCryptoUser(ctx context.Context, page int, order string, search string, isCopytradeApproved *bool) ([]model.CryptoUser, error)
	CountCryptoUser(ctx context.Context, search string, isCopytradeApproved *bool) (int, error)
	GetRefferalScore(ctx context.Context) ([]model.UserWithReferralScores, error)
	GetUserReferral(ctx context.Context) ([]model.UserWithRefCode, error)
	UpdateUserApprove(ctx context.Context, uid string, approve bool) error

	// New: CRM details by user UUID (Privy copytrade)
	GetUserByID(ctx context.Context, userID string) (*model.CryptoUser, error)
	GetPrivyWalletsByUserUUID(ctx context.Context, uuid string) ([]model.WalletInfo, error)
	GetAuthorsByWalletID(ctx context.Context, walletID string) ([]model.SubscribeAuthor, error)
	GetTradeLogsByUserUUID(ctx context.Context, uuid string, status *string, order string, limit, offset int) ([]model.TradeLog, error)
}

type CRMService interface {
	CRMLogin(ctx context.Context, body model.CRMLoginBody) (string, error)
	CheckCRMUserIsExit(ctx context.Context, uid string) (bool, error)
	ValidateDisplaycode(ctx context.Context, displayCode string) (bool, error)
	CheckKOLUserIsExit(ctx context.Context, uid string) (bool, error)
	InsertKolUser(ctx context.Context, kolUser model.KOLUser) (string, error)
	GetKolReferDetail(ctx context.Context) ([]map[string]interface{}, error)
	GetAllUsers(ctx context.Context) ([]UserBasicInfo, error)
	UpdateDisplayCode(ctx context.Context, cryptoUserID, displayCode string) error
	GetCryptoUser(ctx context.Context, page int, order string, search string, isCopytradeApproved *bool) ([]model.CryptoUser, error)
	CountCryptoUser(ctx context.Context, search string, isCopytradeApproved *bool) (int, error)
	GetRefferalScore(ctx context.Context) ([]model.UserWithReferralScores, error)
	GetUserReferral(ctx context.Context) ([]model.UserWithRefCode, error)
	UpdateUserApprove(ctx context.Context, uid string, approve bool) error

	// New aggregation: Privy user overview
	GetPrivyUserOverview(ctx context.Context, userID string, status *string, order string, limit, offset int) (*model.PrivyUserOverview, error)
}
