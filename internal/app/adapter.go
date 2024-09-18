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
