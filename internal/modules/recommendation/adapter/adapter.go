package adapter

import (
	"context"

	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/pkg/proto/gengrpc"
)

type RecommendService interface {
	CreateRecommendations(params RecommendationsParams) error
	GetFilmNamesForRecommendations(ctx context.Context, userID string) ([]*gengrpc.Film, error)
	GetFilmName(film *model.Film) *gengrpc.Film
	GetUniqueFilmIDsForRecommendations(ctx context.Context, recommendations []string) ([]*int, error)
	AddFilmRecommendations(ctx context.Context, userID string, filmIds []*int) error
}

type ModuleInterface interface {
	Service() (RecommendService, error)
}

type RecommendationsParams struct {
	UserID string
}
