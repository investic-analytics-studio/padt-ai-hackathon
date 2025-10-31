package service

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type JwtService struct {
	firebaseClient *auth.Client
	config         *config.Config
	authRepo       port.AuthRepo
}

type jwtCustomClaims struct {
	jwt.RegisteredClaims
	port.JwtClaims
}

func NewJwtService(firebaseClient *auth.Client, config *config.Config, authRepo port.AuthRepo) port.JwtService {
	return &JwtService{
		firebaseClient: firebaseClient,
		config:         config,
		authRepo:       authRepo,
	}
}

func (s *JwtService) GenerateToken(ctx context.Context, uid string, loginType string) (string, error) {
	claims := jwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		JwtClaims: port.JwtClaims{
			UID:       uid,
			LoginType: loginType,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.Jwt.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}
func (s *JwtService) VerifyCRMToken(tokenString string) (*port.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Jwt.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// check if user exists in database
	exists, err := s.authRepo.CheckCRMUserByUid(claims.UID)
	if err != nil || !exists {
		return nil, fmt.Errorf("user not found or invalid")
	}

	return &claims.JwtClaims, nil
}
func (s *JwtService) VerifyToken(tokenString string) (*port.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Jwt.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// check if user exists in database
	exists, err := s.authRepo.CheckUserByUid(claims.UID)
	if err != nil || !exists {
		return nil, fmt.Errorf("user not found or invalid")
	}

	return &claims.JwtClaims, nil
}

func (s *JwtService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	claims, err := s.VerifyToken(oldToken)
	if err != nil {
		return "", fmt.Errorf("failed to verify token: %v", err)
	}

	// Check if user still exists and is valid
	exists, err := s.authRepo.CheckUserByUid(claims.UID)
	if err != nil || !exists {
		return "", fmt.Errorf("user not found or invalid")
	}

	// Generate new token
	return s.GenerateToken(ctx, claims.UID, claims.LoginType)
}
