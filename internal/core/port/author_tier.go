package port

type AuthorTierRepo interface {
	GetAuthorsByTier(tier string) ([]string, error)
	GetAllTiers() ([]string, error)
}

type AuthorTierService interface {
	GetAuthorsByTier(tier string) ([]string, error)
	GetAllTiers() ([]string, error)
}
