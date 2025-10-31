package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type CryptoUserRefcodeRepo interface {
	GenerateAndSaveRefcode(ctx context.Context, cryptoUserID string) ([]string, error)
	CheckUserIDExists(ctx context.Context, cryptoUserID string) (bool, error)
	CheckAndUpdateRefcode(ctx context.Context, refcode, cryptoUserID string) (bool, error)
	CheckKolcode(ctx context.Context, refcode, cryptoUserID string) (model.KolCodeEntry, error)
	InsertKolcode(ctx context.Context, kolUserID, refcode, userID string) error
	GetByCryptoUserID(ctx context.Context, cryptoUserID string) ([]*domain.CryptoUserRefcode, error)
	CheckRefcodeCountByUserID(ctx context.Context, cryptoUserID string) (int, error)
	GetCryptoKolCode(ctx context.Context, userID string) (model.KolUsed, error)
	GenerateRefcodeBynumRequest(ctx context.Context, cryptoUserID string, genNum int) ([]string, error)
	GetRefferalScore(ctx context.Context) ([]*model.RefferalScore, error)
	CheckXUserIsExit(ctx context.Context, twitterName string) (model.CheckXUser, error)
}

type CryptoUserRefcodeService interface {
	GenerateAndSaveRefcode(ctx context.Context, cryptoUserID string) ([]string, error)
	CheckUserIDExists(ctx context.Context, cryptoUserID string) (bool, error)
	CheckAndUpdateRefcode(ctx context.Context, refcode, cryptoUserID string) (bool, error)
	CheckKolcode(ctx context.Context, refcode, cryptoUserID string) (model.KolCodeEntry, error)
	InsertKolcode(ctx context.Context, kolUserID, refcode, userID string) error
	GetByCryptoUserID(ctx context.Context, cryptoUserID string) ([]*domain.CryptoUserRefcode, error)
	CheckRefcodeCountByUserID(ctx context.Context, cryptoUserID string) (int, error)
	GetCryptoKolCode(ctx context.Context, userID string) (model.KolUsed, error)
	GenerateRefcodeBynumRequest(ctx context.Context, cryptoUserID string, genNum int) ([]string, error)
	GetRefferalScore(ctx context.Context) ([]*model.RefferalScore, error)
	GetRefferalScoreRanking(ctx context.Context, offsetDays int) ([]*model.RefferalScoreRanking, error)
	CheckXUserIsExit(ctx context.Context, twitterName string) (model.CheckXUser, error)
}
