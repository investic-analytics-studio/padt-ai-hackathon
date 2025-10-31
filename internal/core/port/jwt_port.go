package port

import "context"

type JwtService interface {
	GenerateToken(ctx context.Context, uid string, loginType string) (string, error)
	VerifyToken(tokenString string) (*JwtClaims, error)
	VerifyCRMToken(tokenString string) (*JwtClaims, error)
	RefreshToken(ctx context.Context, oldToken string) (string, error)
}

type JwtClaims struct {
	UID       string `json:"uid"`
	LoginType string `json:"login_type"`
}
