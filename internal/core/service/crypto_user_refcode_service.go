package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/domain"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type CryptoUserRefcodeService struct {
	repo *repo.CryptoUserRefcodeRepo
}

func NewCryptoUserRefcodeService(repo *repo.CryptoUserRefcodeRepo) *CryptoUserRefcodeService {
	return &CryptoUserRefcodeService{repo: repo}
}

func generateRefcode() (string, error) {
	const maxRetries = 3
	var refcode string
	var err error

	for i := 0; i < maxRetries; i++ {
		// Generate 5 random bytes
		bytes := make([]byte, 5)
		if _, err = rand.Read(bytes); err != nil {
			continue
		}

		// Encode to base32 and take first 8 characters
		refcode = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
		refcode = strings.ToUpper(refcode[:8])

		// If we got here, we have a valid refcode
		return refcode, nil
	}

	return "", &repo.ErrRefcodeGeneration{
		Attempts: maxRetries,
		Err:      fmt.Errorf("failed to generate refcode after %d attempts", maxRetries),
	}
}
func (s *CryptoUserRefcodeService) GenerateRefcodeBynumRequest(ctx context.Context, cryptoUserID string, refcodeNum int) ([]string, error) {
	const maxRetries = 10
	refcodes := make([]string, 0, refcodeNum)
	attempts := 0
	for len(refcodes) < refcodeNum && attempts < maxRetries*refcodeNum {
		attempts++

		refcode, err := generateRefcode()
		if err != nil {
			continue
		}

		available, err := s.CheckRefcodeNotDuplicate(ctx, refcode)
		if err != nil {
			return nil, err
		}

		if !available {
			continue
		}

		err = s.repo.Create(ctx, cryptoUserID, refcode)
		if err != nil {
			return nil, err
		}

		refcodes = append(refcodes, refcode)

	}
	if len(refcodes) < refcodeNum {
		return nil, &repo.ErrRefcodeGeneration{
			Attempts: attempts,
			Err:      fmt.Errorf("only generated %d/%d refcodes", len(refcodes), refcodeNum),
		}
	}

	return refcodes, nil
}
func (s *CryptoUserRefcodeService) GenerateAndSaveRefcode(ctx context.Context, cryptoUserID string) ([]string, error) {
	const maxRetries = 10
	const refcodeNum = 5

	refcodes := make([]string, 0, refcodeNum)
	attempts := 0

	refcodeCount, err := s.repo.CheckRefcodeCountByUserID(ctx, cryptoUserID)
	if err != nil {
		return nil, err
	}

	if refcodeCount >= refcodeNum {
		return nil, &repo.ErrRefcodeGeneration{
			Attempts: attempts,
			Err:      fmt.Errorf("already have %d refcodes", refcodeCount),
		}
	}

	for len(refcodes) < refcodeNum && attempts < maxRetries*refcodeNum {
		attempts++

		refcode, err := generateRefcode()
		if err != nil {
			continue
		}

		available, err := s.CheckRefcodeNotDuplicate(ctx, refcode)
		if err != nil {
			return nil, err
		}

		if !available {
			continue
		}

		err = s.repo.Create(ctx, cryptoUserID, refcode)
		if err != nil {
			return nil, err
		}

		refcodes = append(refcodes, refcode)
	}

	if len(refcodes) < refcodeNum {
		return nil, &repo.ErrRefcodeGeneration{
			Attempts: attempts,
			Err:      fmt.Errorf("only generated %d/%d refcodes", len(refcodes), refcodeNum),
		}
	}

	return refcodes, nil

	// for i := 0; i < maxRetries; i++ {
	// 	refcode, err := generateRefcode()
	// 	if err != nil {
	// 		continue
	// 	}

	// 	available, err := s.CheckRefcodeNotDuplicate(ctx, refcode)
	// 	if err != nil {
	// 		return "", &repo.ErrDatabaseOperation{
	// 			Operation: "check refcode not duplicate",
	// 			Err:       err,
	// 		}
	// 	}

	// 	if available {
	// 		// Save to database
	// 		err = s.repo.Create(ctx, cryptoUserID, refcode)
	// 		if err != nil {
	// 			return "", &repo.ErrDatabaseOperation{
	// 				Operation: "create crypto user refcode",
	// 				Err:       err,
	// 			}
	// 		}
	// 		return refcode, nil
	// 	}
	// }
	// return "", &repo.ErrRefcodeGeneration{
	// 	Attempts: maxRetries,
	// 	Err:      fmt.Errorf("failed to generate refcode after %d attempts", maxRetries),
	// }
}

func (s *CryptoUserRefcodeService) GetByCryptoUserID(ctx context.Context, cryptoUserID string) ([]*domain.CryptoUserRefcode, error) {
	return s.repo.GetByCryptoUserID(ctx, cryptoUserID)
}
func (s *CryptoUserRefcodeService) GetRefferalScore(ctx context.Context) ([]*model.RefferalScore, error) {
	return s.repo.GetRefferalScore(ctx)
}
func (s *CryptoUserRefcodeService) CheckUserIDExists(ctx context.Context, userID string) (bool, error) {
	exists, err := s.repo.CheckUserIDExists(ctx, userID)
	if err != nil {
		return false, &repo.ErrDatabaseOperation{
			Operation: "check user ID existence",
			Err:       err,
		}
	}
	return exists, nil
}

func (s *CryptoUserRefcodeService) CheckRefcodeAvailable(ctx context.Context, refcode string) (bool, error) {
	available, err := s.repo.CheckRefcodeAvailable(ctx, refcode)
	if err != nil {
		return false, &repo.ErrDatabaseOperation{
			Operation: "check refcode availability",
			Err:       err,
		}
	}
	return available, nil
}
func (s *CryptoUserRefcodeService) CheckKolcode(ctx context.Context, refcode string, userID string) (model.KolCodeEntry, error) {
	kolCode, err := s.repo.CheckKolcode(ctx, refcode, userID)
	if err != nil {
		return model.KolCodeEntry{}, err
	}
	return kolCode, nil

}

// InsertKolcode(ctx context.Context, kolUserID, refcode, userID string) error
func (s *CryptoUserRefcodeService) InsertKolcode(ctx context.Context, kolUserID string, refcode string, userID string) error {
	err := s.repo.InsertKolcode(ctx, kolUserID, refcode, userID)
	if err != nil {
		return err
	}
	return nil

}
func (s *CryptoUserRefcodeService) CheckAndUpdateRefcode(ctx context.Context, refcode string, userID string) (bool, error) {
	updated, err := s.repo.CheckAndUpdateRefcode(ctx, refcode, userID)
	if err != nil {
		return false, err
	}
	return updated, nil
}

func (s *CryptoUserRefcodeService) CheckRefcodeNotDuplicate(ctx context.Context, refcode string) (bool, error) {
	return s.repo.CheckRefcodeNotDuplicate(ctx, refcode)
}

func (s *CryptoUserRefcodeService) CheckRefcodeCountByUserID(ctx context.Context, cryptoUserID string) (int, error) {
	return s.repo.CheckRefcodeCountByUserID(ctx, cryptoUserID)
}

func (s *CryptoUserRefcodeService) GetCryptoKolCode(ctx context.Context, cryptoUserID string) (model.KolUsed, error) {
	return s.repo.GetCryptoKolCode(ctx, cryptoUserID)

}

func (s *CryptoUserRefcodeService) GetRefferalScoreRanking(ctx context.Context, offsetDays int) ([]*model.RefferalScoreRanking, error) {
	return s.repo.GetRefferalScoreRanking(ctx, offsetDays)
}

func (s *CryptoUserRefcodeService) CheckXUserIsExit(ctx context.Context, twitterName string) (model.CheckXUser, error) {
	return s.repo.CheckXUserIsExit(ctx, twitterName)
}
