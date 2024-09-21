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
	Save(sequel []service.Sequel) error
}
