package port

import (
	"context"

	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type UserInfo struct {
	Email       string
	TwitterUID  string
	TwitterName string
}

type AuthRepo interface {
	Create(body *model.SignUpBody, methodSignUp string) error
	ExistEmail(email string) (bool, error)
	GetUserByEmail(email string) (*model.User, error)
	CheckUserByUid(uid string) (bool, error)
	CheckCRMUserByUid(uid string) (bool, error)
	GetUserInfo(uid string) (UserInfo, error)
	GetCRMUserInfo(uid string) (model.CRMUser, error)
	ExistTwitterUID(uid string) (bool, error)
}

type AuthService interface {
	SignUp(body model.SignUpBody) error
	ExistEmail(email string) (bool, error)
	GetUserByEmail(email string) (*model.User, error)
	AutoValidateEmailInFirebase(ctx context.Context, uid string) error
	AllUserAutoValidateEmailInFirebase(ctx context.Context) error
	Login(ctx context.Context, body model.LoginBody) (string, error)
	ExistTwitterUID(uid string) (bool, error)
	// GenerateWalletToken(ctx context.Context, walletAddress string) (string, error)
}
