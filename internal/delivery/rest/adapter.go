package rest

import (
	"kinopoisk-api/internal/service"
)

type FilmService interface {
	GetOne(filmId string) (service.Film, error)
}

type UserService interface {
	GetOne(userToken string) (service.User, error)
}

type FilmSequelService interface {
	GetAll(filmId string) ([]service.Film, error)
}

type FilmSimilarService interface {
	GetAll(filmId string) ([]service.Film, error)
}
