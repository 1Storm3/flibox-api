package recommendation

import (
	"context"
	"strconv"

	"github.com/1Storm3/flibox-api/internal/delivery/grpc"
	"github.com/1Storm3/flibox-api/internal/model"
	filmService "github.com/1Storm3/flibox-api/internal/modules/film"
	historyfilms "github.com/1Storm3/flibox-api/internal/modules/history-films"
	"github.com/1Storm3/flibox-api/internal/modules/recommendation/adapter"
	userfilm "github.com/1Storm3/flibox-api/internal/modules/user-film"
	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"github.com/1Storm3/flibox-api/pkg/proto/gengrpc"
)

type Service struct {
	historyFilmsService historyfilms.ServiceInterface
	filmService         filmService.ServiceInterface
	userFilmService     userfilm.ServiceInterface
	grpcClient          grpc.ClientConnInterface
}

func NewRecommendationService(
	grpcClient grpc.ClientConnInterface,
	filmService filmService.ServiceInterface,
	userFilmService userfilm.ServiceInterface,
	historyFilmsService historyfilms.ServiceInterface) *Service {
	return &Service{
		grpcClient:          grpcClient,
		filmService:         filmService,
		userFilmService:     userFilmService,
		historyFilmsService: historyFilmsService,
	}
}

func (s *Service) CreateRecommendations(params adapter.RecommendationsParams) error {
	ctx := context.Background()

	filmNames, err := s.GetFilmNamesForRecommendations(ctx, params.UserID)
	if err != nil {
		return err
	}

	if len(filmNames) == 0 {
		logger.Info("Нет фильмов для рекомендаций")
		return nil
	}

	err = s.userFilmService.DeleteMany(ctx, params.UserID)
	if err != nil {
		return err
	}

	recommendations, err := s.grpcClient.GetRecommendations(ctx, filmNames)
	if err != nil {
		return err
	}

	filmIds, err := s.GetUniqueFilmIDsForRecommendations(ctx, recommendations)
	if err != nil {
		return err
	}

	if len(filmIds) == 0 {
		logger.Info("Нет рекомендаций")
		return nil
	}

	err = s.AddFilmRecommendations(ctx, params.UserID, filmIds)
	if err != nil {
		return err
	}

	logger.Info("Рекомендации созданы")
	return nil
}

func (s *Service) GetFilmNamesForRecommendations(ctx context.Context, userID string) ([]*gengrpc.Film, error) {
	var filmNames []*gengrpc.Film

	favouriteFilms, err := s.userFilmService.GetAll(ctx, userID, model.TypeUserFavourite, 5)
	if err != nil {
		return nil, err
	}

	historyFilms, err := s.historyFilmsService.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, historyFilm := range historyFilms {
		film := &historyFilm.Film
		filmNames = append(filmNames, s.GetFilmName(film))
	}

	for _, favouriteFilm := range favouriteFilms {
		film := &favouriteFilm.Film
		filmNames = append(filmNames, s.GetFilmName(film))
	}

	return filmNames, nil
}

func (s *Service) GetFilmName(film *model.Film) *gengrpc.Film {
	var filmName string
	if film.NameOriginal != nil {
		filmName = *film.NameOriginal
	} else if film.NameRU != nil {
		filmName = *film.NameRU
	}
	return &gengrpc.Film{NameOriginal: filmName}
}

func (s *Service) GetUniqueFilmIDsForRecommendations(ctx context.Context, recommendations []string) ([]*int, error) {
	var filmIds []*int
	seenFilmIDs := make(map[int]struct{})

	for _, film := range recommendations {
		filmExist, err := s.filmService.GetOneByNameRu(ctx, film)
		if err != nil {
			return nil, err
		}

		if filmExist.ID == nil {
			// Запрос во внешний апи
			continue
		}

		if _, exists := seenFilmIDs[*filmExist.ID]; exists {
			continue
		}
		seenFilmIDs[*filmExist.ID] = struct{}{}
		filmIds = append(filmIds, filmExist.ID)
	}

	return filmIds, nil
}

func (s *Service) AddFilmRecommendations(ctx context.Context, userID string, filmIds []*int) error {
	var recommendFilms []userfilm.Params
	for _, id := range filmIds {
		recommendFilms = append(recommendFilms, userfilm.Params{
			UserID: userID,
			FilmID: strconv.Itoa(*id),
			Type:   model.TypeUserRecommend,
		})
	}

	return s.userFilmService.AddMany(ctx, recommendFilms)
}
