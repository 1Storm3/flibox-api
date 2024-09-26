package rest

import (
	"kinopoisk-api/internal/service"
)

type FilmService interface {
	GetOne(filmId string) (service.Film, error) // domain model, httperror
}

type SequelService interface {
	GetAll(filmId string) ([]service.Sequel, error)
}

type UserService interface {
	GetOne(userToken string) (service.User, error)
}
