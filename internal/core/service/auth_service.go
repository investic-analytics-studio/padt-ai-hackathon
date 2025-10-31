package service

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
	"google.golang.org/api/iterator"
)

type AuthService struct {
	firebaseClient *auth.Client
	authRepo       port.AuthRepo
	jwtService     port.JwtService
	// logger logger.Logger
}

func NewAuthService(firebaseClient *auth.Client, authRepo port.AuthRepo, jwtService port.JwtService) port.AuthService {
	return &AuthService{
		firebaseClient: firebaseClient,
		authRepo:       authRepo,
		jwtService:     jwtService,
	}
}

func (s *AuthService) SignUp(body model.SignUpBody) error {

	if body.MethodSignUp == model.SignUpMethodEmail {
		data := (&auth.UserToCreate{}).
			Email(body.Email).
			Password(body.Password)

		newCtx := context.Background()
		result, err := s.firebaseClient.CreateUser(newCtx, data)
		body.UID = result.UID
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}

		if !result.EmailVerified {
			if err := s.AutoValidateEmailInFirebase(newCtx, result.UID); err != nil {
				return err
			}
		}

	} else if body.MethodSignUp == model.SignUpMethodGoogle {
		fmt.Println("google sign up")
	} else if body.MethodSignUp == model.SignUpXMethod {
		fmt.Println("twitter sign up")
	} else {
		return fmt.Errorf("invalid sign up method")
	}

	userData := model.SignUpBody{
		UID:         body.UID,
		Email:       body.Email,
		TwitterUID:  body.TwitterUID,
		TwitterName: body.TwitterName,
	}

	if err := s.authRepo.Create(&userData, body.MethodSignUp); err != nil {
		if err := s.firebaseClient.DeleteUser(context.Background(), body.UID); err != nil {
			return fmt.Errorf("failed to delete user: %v", err)
		}
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (s *AuthService) ExistEmail(email string) (bool, error) {
	return s.authRepo.ExistEmail(email)
}

func (s *AuthService) GetUserByEmail(email string) (*model.User, error) {
	return s.authRepo.GetUserByEmail(email)
}

func (s *AuthService) AutoValidateEmailInFirebase(ctx context.Context, uid string) error {
	user, err := s.firebaseClient.GetUser(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if !user.EmailVerified {
		updatedUser, err := s.firebaseClient.UpdateUser(ctx, uid, (&auth.UserToUpdate{}).EmailVerified(true))
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
		if !updatedUser.EmailVerified {
			return fmt.Errorf("failed to verify email for user: %s", uid)
		}
	}

	return nil
}

func (s *AuthService) AllUserAutoValidateEmailInFirebase(ctx context.Context) error {
	iter := s.firebaseClient.Users(ctx, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error listing users: %v", err)
		}

		if !user.EmailVerified {
			if err := s.AutoValidateEmailInFirebase(ctx, user.UID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *AuthService) Login(ctx context.Context, body model.LoginBody) (string, error) {
	// check firebase uid is valid
	_, err := s.firebaseClient.GetUser(ctx, body.UID)
	if err != nil {
		return "", fmt.Errorf("credential is invalid")
	}

	// check if user exists in database
	exists, err := s.authRepo.CheckUserByUid(body.UID)
	if err != nil {
		return "", fmt.Errorf("credential is invalid")
	}

	if !exists {
		return "", fmt.Errorf("credential is invalid")
	}

	// generate JWT token
	token, err := s.jwtService.GenerateToken(ctx, body.UID, body.LoginType)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func (s *AuthService) ExistTwitterUID(uid string) (bool, error) {
	return s.authRepo.ExistTwitterUID(uid)
}

// func listAllUsers(client *auth.Client) ([]map[string]interface{}, error) {
// 	var users []map[string]interface{}

// 	iter := client.Users(context.Background(), "")
// 	for {
// 		userRecord, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			return nil, fmt.Errorf("error iterating users: %v", err)
// 		}

// 		// Add user details to the list
// 		users = append(users, map[string]interface{}{
// 			"uid":           userRecord.UID,
// 			"emailVerified": userRecord.EmailVerified,
// 			"email":         userRecord.Email,
// 		})
// 	}

// 	return users, nil
// }

// func (s *AuthService) GenerateWalletToken(ctx context.Context, walletAddress string) (string, error) {
// 	// Create a custom user with the wallet address as UID
// 	data := (&auth.UserToCreate{}).
// 		UID(walletAddress)

// 	// Create user in Firebase if not exists
// 	_, err := s.firebaseClient.CreateUser(ctx, data)
// 	if err != nil {
// 		// If user already exists, ignore the error
// 		if !strings.Contains(err.Error(), "already exists") {
// 			return "", fmt.Errorf("failed to create user: %v", err)
// 		}
// 	}

// 	// Generate custom token for the wallet address
// 	token, err := s.firebaseClient.CustomToken(ctx, walletAddress)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate custom token: %v", err)
// 	}

// 	// // Create user in database if not exists
// 	// userData := model.SignUpBody{
// 	// 	UID:          walletAddress,
// 	// 	MethodSignUp: "wallet",
// 	// }

// 	// // Try to create user in database, ignore if already exists
// 	// _ = s.repo.Create(&userData)

// 	return token, nil
// }
