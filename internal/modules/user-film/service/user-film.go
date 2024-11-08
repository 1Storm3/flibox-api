package service

import (
	"context"

	"net/http"

	"kbox-api/internal/modules/user-film/dto"
	"kbox-api/internal/modules/user-film/mapper"
	"kbox-api/internal/modules/user-film/repository"
	"kbox-api/shared/httperror"
)

var _ UserFilmServiceInterface = (*UserFilmService)(nil)

type UserFilmServiceInterface interface {
	GetAll(userId string) ([]dto.GetUserFilmResponseDTO, error)
	Add(userId string, filmId string) error
	Delete(userId string, filmId string) error
}

type UserFilmService struct {
	userFilmRepo repository.UserFilmRepositoryInterface
}

func NewUserFilmService(userFilmRepo repository.UserFilmRepositoryInterface) UserFilmServiceInterface {
	return &UserFilmService{
		userFilmRepo: userFilmRepo,
	}
}

func (s *UserFilmService) GetAll(userId string) ([]dto.GetUserFilmResponseDTO, error) {
	result, err := s.userFilmRepo.GetAll(context.Background(), userId)

	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, httperror.New(
			http.StatusNotFound,
			"Избранные фильмы не найдены у этого пользователя",
		)
	}

	var res []dto.GetUserFilmResponseDTO
	for _, film := range result {
		res = append(res, mapper.MapDomainUserFilmToResponseDTO(film))
	}

	return res, nil
}

func (s *UserFilmService) Add(userId string, filmId string) error {
	return s.userFilmRepo.Add(context.Background(), userId, filmId)
}

func (s *UserFilmService) Delete(userId string, filmId string) error {
	return s.userFilmRepo.Delete(context.Background(), userId, filmId)
}
