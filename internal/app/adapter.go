package app

import (
	"context"
	"kinopoisk-api/internal/service"
)

type FilmService interface {
	GetOne(filmId string) (service.Film, error)
}

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (service.Film, error)
	Save(film service.Film) error
}

type SequelService interface {
	GetAll(filmId string) ([]service.Sequel, error)
}

type SequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]service.Sequel, error)
	Save(filmId string, sequel []service.Sequel) error
}

type FilmSequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]service.Sequel, error)
	Save(sequel []service.Sequel) error
}

type FilmSequelService interface {
	GetAll(filmId string) ([]service.Sequel, error)
}

type UserService interface {
	GetOne(userToken string) (service.User, error)
}

type UserRepository interface {
	GetOne(ctx context.Context, userToken string) (service.User, error)
}
