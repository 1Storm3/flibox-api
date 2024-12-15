package userfilm

import (
	"context"
	"net/http"

	"github.com/1Storm3/flibox-api/internal/model"
	"github.com/1Storm3/flibox-api/internal/shared/httperror"
)

var _ ServiceInterface = (*Service)(nil)

type ServiceInterface interface {
	GetAll(ctx context.Context, userId string, typeUserFilm model.TypeUserFilm, limit int) ([]GetUserFilmResponseDTO, error)
	Add(ctx context.Context, params Params) error
	Delete(ctx context.Context, params Params) error
	AddMany(ctx context.Context, params []Params) error
	DeleteMany(ctx context.Context, userID string) error
}

type Service struct {
	repository RepositoryInterface
}

func NewUserFilmService(
	repository RepositoryInterface,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddMany(ctx context.Context, params []Params) error {
	return s.repository.AddMany(ctx, params)
}

func (s *Service) DeleteMany(ctx context.Context, userID string) error {
	return s.repository.DeleteMany(ctx, userID)
}

func (s *Service) GetAll(ctx context.Context, userId string, typeUserFilm model.TypeUserFilm, limit int) ([]GetUserFilmResponseDTO, error) {
	result, err := s.repository.GetAllForRecommend(ctx, userId, typeUserFilm, limit)

	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		if typeUserFilm == model.TypeUserFavourite {
			return []GetUserFilmResponseDTO{}, httperror.New(
				http.StatusNotFound,
				"Избранные фильмы не найдены у этого пользователя",
			)
		} else {
			return []GetUserFilmResponseDTO{},
				httperror.New(
					http.StatusNotFound,
					"Рекомендации не найдены у этого пользователя",
				)
		}
	}

	var res []GetUserFilmResponseDTO
	for _, film := range result {
		res = append(res, MapDomainUserFilmToResponseDTO(film))
	}

	return res, nil
}

func (s *Service) Add(ctx context.Context, params Params) error {
	return s.repository.Add(ctx, params)
}

func (s *Service) Delete(ctx context.Context, params Params) error {
	return s.repository.Delete(ctx, params)
}
