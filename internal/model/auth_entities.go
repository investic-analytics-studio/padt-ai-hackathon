package model

const (
	SignUpMethodEmail  = "email"
	SignUpMethodGoogle = "google"
	SignUpXMethod      = "twitter"
)

type SignUpBody struct {
	MethodSignUp string `json:"method_sign_up" validate:"required"`
	UID          string `json:"uid"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password"`
	TwitterUID   string `json:"twitter_uid"`
	TwitterName  string `json:"twitter_name"`
}

type EmailReq struct {
	Email string `json:"email" validate:"required,email"`
}

type AutoValidateRequest struct {
	UID string `json:"uid" validate:"required"`
}

type GenerateTokenRequest struct {
	UID       string `json:"uid" validate:"required"`
	LoginType string `json:"login_type" validate:"required"`
}

type VerifyTokenRequest struct {
	JWTToken string `json:"jwt_token" validate:"required"`
}

type RefreshTokenRequest struct {
	JWTToken string `json:"jwt_token" validate:"required"`
}

type LoginBody struct {
	UID       string `json:"uid" validate:"required"`
	LoginType string `json:"login_type" validate:"required"`
}

type TwitterUIDReq struct {
	TwitterUID string `json:"twitter_uid" validate:"required"`
}
