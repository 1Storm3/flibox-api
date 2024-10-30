package service

import (
	"context"
	"kinopoisk-api/internal/modules/film/service"
)

type UserFilm struct {
	UserId string       `json:"userId" gorm:"column:user_id"`
	FilmId int          `json:"filmId" gorm:"column:film_id"`
	Film   service.Film `gorm:"foreignKey:FilmId;references:ID"`
}

type UserFilmService struct {
	userFilmRepo UserFilmRepository
}

func NewUserFilmService(userFilmRepo UserFilmRepository) *UserFilmService {
	return &UserFilmService{
		userFilmRepo: userFilmRepo,
	}
}

func (s *UserFilmService) GetAll(userId string) ([]UserFilm, error) {
	return s.userFilmRepo.GetAll(context.Background(), userId)
}

func (s *UserFilmService) Add(userId string, filmId string) error {
	return s.userFilmRepo.Add(context.Background(), userId, filmId)
}

func (s *UserFilmService) Delete(userId string, filmId string) error {
	return s.userFilmRepo.Delete(context.Background(), userId, filmId)
}
