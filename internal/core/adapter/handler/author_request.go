package handler

type AuthorRequest struct {
	Author   string `json:"author"`
	WalletID string `json:"wallet_id"`
}

type UnAuthorRequest struct {
	Author   string `json:"author"`
	WalletID string `json:"wallet_id"`
}
