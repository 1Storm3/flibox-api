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

type FilmSequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]service.FilmSequel, error)
	Save(filmId int, sequelId int) error
}

type FilmSequelService interface {
	GetAll(filmId string) ([]service.Film, error)
}

type UserService interface {
	GetOne(userToken string) (service.User, error)
}

type UserRepository interface {
	GetOne(ctx context.Context, userToken string) (service.User, error)
}

type FilmSimilarRepository interface {
	GetAll(ctx context.Context, filmId string) ([]service.FilmSimilar, error)
	Save(filmId int, similarId int) error
}

type FilmSimilarService interface {
	GetAll(filmId string) ([]service.Film, error)
}
