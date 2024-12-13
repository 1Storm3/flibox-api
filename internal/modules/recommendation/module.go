package recommendation

import (
	"kbox-api/internal/delivery/grpc"
	"kbox-api/internal/modules/film"
	historyfilms "kbox-api/internal/modules/history-films"
	"kbox-api/internal/modules/recommendation/adapter"
	userfilm "kbox-api/internal/modules/user-film"
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
