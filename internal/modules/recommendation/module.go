package recommendation

import (
	"github.com/1Storm3/flibox-api/internal/delivery/grpc"
	"github.com/1Storm3/flibox-api/internal/modules/film"
	historyfilms "github.com/1Storm3/flibox-api/internal/modules/history-films"
	"github.com/1Storm3/flibox-api/internal/modules/recommendation/adapter"
	userfilm "github.com/1Storm3/flibox-api/internal/modules/user-film"
)

type Module struct {
	recommendationService adapter.RecommendService
	filmModule            film.ModuleInterface
	historyFilmsModule    historyfilms.ModuleInterface
	userFilmModule        userfilm.ModuleInterface
	grpcClient            grpc.ClientConnInterface
}

func NewRecommendationModule(
	filmModule film.ModuleInterface,
	historyFilmsModule historyfilms.ModuleInterface,
	userFilmModule userfilm.ModuleInterface,
	grpcClient grpc.ClientConnInterface,
) *Module {
	return &Module{
		filmModule:         filmModule,
		historyFilmsModule: historyFilmsModule,
		userFilmModule:     userFilmModule,
		grpcClient:         grpcClient,
	}
}

func (m *Module) Service() (adapter.RecommendService, error) {
	if m.recommendationService == nil {
		filmService, err := m.filmModule.Service()
		if err != nil {
			return nil, err
		}
		historyFilmsService, err := m.historyFilmsModule.Service()
		if err != nil {
			return nil, err
		}
		userFilmService, err := m.userFilmModule.Service()
		if err != nil {
			return nil, err
		}
		grpcClient := m.grpcClient
		m.recommendationService = NewRecommendationService(
			grpcClient, filmService, userFilmService, historyFilmsService,
		)
	}
	return m.recommendationService, nil
}
