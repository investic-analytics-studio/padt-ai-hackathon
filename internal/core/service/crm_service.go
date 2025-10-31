package service

import (
	"context"
	"fmt"

	"github.com/quantsmithapp/datastation-backend/internal/core/adapter/repo"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type CryptoCRMService struct {
	repo       *repo.CRMRepo
	jwtService port.JwtService
}

func NewCryptoCRMService(repo *repo.CRMRepo, jwtService port.JwtService) *CryptoCRMService {

	return &CryptoCRMService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *CryptoCRMService) CRMLogin(ctx context.Context, body model.CRMLoginBody) (string, error) {
	uid, err := s.repo.CheckCRMUser(body)
	if err != nil {
		return "", fmt.Errorf("user validation failed")
	}
	if uid == "" {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(ctx, uid, "CRM")
	if err != nil {
		return "", fmt.Errorf("token generation failed: %v", err)
	}

	return token, nil
}

func (s *CryptoCRMService) CheckCRMUserIsExit(ctx context.Context, uid string) (bool, error) {
	return s.repo.CheckCRMUserIsExit(ctx, uid)
}
func (s *CryptoCRMService) ValidateDisplaycode(ctx context.Context, displayCode string) (bool, error) {
	return s.repo.ValidateDisplaycode(ctx, displayCode)
}
func (s *CryptoCRMService) CheckKOLUserIsExit(ctx context.Context, uid string) (bool, error) {
	return s.repo.CheckKOLUserIsExit(ctx, uid)
}
func (s *CryptoCRMService) InsertKolUser(ctx context.Context, kolUser model.KOLUser) (string, error) {
	return s.repo.InsertKolUser(ctx, kolUser)
}
func (s *CryptoCRMService) GetKolReferDetail(ctx context.Context) ([]map[string]interface{}, error) {
	details, err := s.repo.GetKolReferDetails(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, len(details))
	for i, detail := range details {
		results[i] = map[string]interface{}{
			"uid":                 detail.UID,
			"email":               detail.Email,
			"twitter":             detail.Twitter,
			"display_code":        detail.DisplayCode,
			"referral_use_number": detail.ReferralUseNumber,
			"created_at":          detail.CreatedAt,
		}
	}

	return results, nil
}

func (s *CryptoCRMService) GetAllUsers(ctx context.Context) ([]port.UserBasicInfo, error) {
	repoUsers, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]port.UserBasicInfo, len(repoUsers))
	for i, u := range repoUsers {
		users[i] = port.UserBasicInfo{
			UUID:        u.UUID,
			Email:       u.Email,
			TwitterName: u.TwitterName,
		}
	}
	return users, nil
}

func (s *CryptoCRMService) UpdateDisplayCode(ctx context.Context, cryptoUserID, displayCode string) error {
	// First check if user exists
	exists, err := s.repo.IsKolUserExists(ctx, cryptoUserID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with ID %s does not exist", cryptoUserID)
	}

	// Then check if display code exists for another user
	exists, err = s.repo.IsDisplayCodeExists(ctx, displayCode, cryptoUserID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("display_code already exists")
	}
	return s.repo.UpdateDisplayCode(ctx, cryptoUserID, displayCode)
}

func (s *CryptoCRMService) GetCryptoUser(ctx context.Context, page int, order string, search string, isCopytradeApproved *bool) ([]model.CryptoUser, error) {
	cryptoUsers, err := s.repo.GetCryptoUser(ctx, page, order, search, isCopytradeApproved)
	if err != nil {
		return nil, err
	}

	return cryptoUsers, nil
}
func (s *CryptoCRMService) GetRefferalScore(ctx context.Context) ([]model.UserWithReferralScores, error) {
	scores, err := s.repo.GetRefferalScore(ctx)
	if err != nil {
		return nil, err
	}
	return scores, nil
}
func (s *CryptoCRMService) GetUserReferral(ctx context.Context) ([]model.UserWithRefCode, error) {
	refUser, err := s.repo.GetUserReferral(ctx)
	if err != nil {
		return nil, err
	}
	return refUser, nil

}
func (s *CryptoCRMService) UpdateUserApprove(ctx context.Context, uid string, approve bool) error {
	return s.repo.UpdateUserApprove(ctx, uid, approve)
}

// CountCryptoUser returns total rows based on the same filters used by GetCryptoUser
func (s *CryptoCRMService) CountCryptoUser(ctx context.Context, search string, isCopytradeApproved *bool) (int, error) {
	return s.repo.CountCryptoUser(ctx, search, isCopytradeApproved)
}

// GetPrivyUserOverview aggregates user, wallets (with authors), and trade logs
func (s *CryptoCRMService) GetPrivyUserOverview(ctx context.Context, userID string, status *string, order string, limit, offset int) (*model.PrivyUserOverview, error) {
	// 1) User Profile
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2) Wallets
	wallets, err := s.repo.GetPrivyWalletsByUserUUID(ctx, user.UUID)
	if err != nil {
		return nil, err
	}
	// Fill authors for each wallet
	for i := range wallets {
		authors, err := s.repo.GetAuthorsByWalletID(ctx, wallets[i].WalletID)
		if err != nil {
			return nil, err
		}
		wallets[i].Authors = authors
	}

	// 3) Trade logs
	logs, err := s.repo.GetTradeLogsByUserUUID(ctx, user.UUID, status, order, limit, offset)
	if err != nil {
		return nil, err
	}

	out := &model.PrivyUserOverview{
		UUID:        user.UUID,
		Email:       user.Email,
		TwitterName: user.TwitterName,
		WalletCount: len(wallets),
		Wallets:     wallets,
		TradeLogs:   logs,
	}
	return out, nil
}
