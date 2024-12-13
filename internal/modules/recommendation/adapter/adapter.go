package adapter

type RecommendService interface {
	CreateRecommendations(params RecommendationsParams) error
}

type ModuleInterface interface {
	Service() (RecommendService, error)
}

type RecommendationsParams struct {
	UserID string
}
