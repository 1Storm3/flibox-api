package recommendation

import (
	"context"
	"strconv"

	"kbox-api/internal/delivery/grpc"
	"kbox-api/internal/model"
	filmService "kbox-api/internal/modules/film"
	historyfilms "kbox-api/internal/modules/history-films"
	"kbox-api/internal/modules/recommendation/adapter"
	userfilm "kbox-api/internal/modules/user-film"
	"kbox-api/internal/shared/logger"
	"kbox-api/pkg/proto/gengrpc"
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

	var filmNames []*gengrpc.Film

	favouriteFilms, err := s.userFilmService.GetAll(ctx, params.UserID, model.TypeUserFavourite, 5)
	if err != nil {
		return err
	}

	historyFilms, err := s.historyFilmsService.GetAll(ctx, params.UserID)
	if err != nil {
		return err
	}

	for _, film := range historyFilms {
		var filmName string
		if film.Film.NameOriginal != nil {
			filmName = *film.Film.NameOriginal
		} else if film.Film.NameRU != nil {
			filmName = *film.Film.NameRU
		}
		filmNames = append(filmNames, &gengrpc.Film{
			NameOriginal: filmName,
		})
	}

	for _, film := range favouriteFilms {
		var filmName string
		if film.Film.NameOriginal != nil {
			filmName = *film.Film.NameOriginal
		} else if film.Film.NameRU != nil {
			filmName = *film.Film.NameRU
		}
		filmNames = append(filmNames, &gengrpc.Film{
			NameOriginal: filmName,
		})
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
	var filmIds []*int
	seenFilmIDs := make(map[int]struct{})

	for _, film := range recommendations {
		filmExist, err := s.filmService.GetOneByNameRu(ctx, film)
		if err != nil {
			return err
		}
		if filmExist.ID == nil {
			// TODO: Интеграция с внешним API для добавления фильмов
			continue
		}

		if _, exists := seenFilmIDs[*filmExist.ID]; exists {
			continue
		}
		seenFilmIDs[*filmExist.ID] = struct{}{}
		filmIds = append(filmIds, filmExist.ID)
	}

	var recommendFilms []userfilm.Params
	for _, id := range filmIds {
		recommendFilms = append(recommendFilms, userfilm.Params{
			UserID: params.UserID,
			FilmID: strconv.Itoa(*id),
			Type:   model.TypeUserRecommend,
		})
	}

	if len(recommendFilms) == 0 {
		logger.Info("Нет рекомендаций")
		return nil
	}

	err = s.userFilmService.AddMany(ctx, recommendFilms)
	if err != nil {
		return err
	}

	logger.Info("Рекомендации созданы")

	return nil
}
