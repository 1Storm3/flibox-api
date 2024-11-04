package service

import (
	"context"
	"net/http"

	"kbox-api/internal/modules/user-film/dto"
	"kbox-api/internal/modules/user-film/mapper"
	"kbox-api/shared/httperror"
)

type UserFilmService struct {
	userFilmRepo UserFilmRepository
}

func NewUserFilmService(userFilmRepo UserFilmRepository) *UserFilmService {
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
