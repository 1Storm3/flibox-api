package app

import (
	"context"
	filmsequelservice "kinopoisk-api/internal/modules/film-sequel/service"
	filmsimilarservice "kinopoisk-api/internal/modules/film-similar/service"
	filmservice "kinopoisk-api/internal/modules/film/service"
	userfilmservice "kinopoisk-api/internal/modules/user-film/service"
	userservice "kinopoisk-api/internal/modules/user/service"
	externalservice "kinopoisk-api/pkg/external-service"
)

type ExternalService interface {
	SetExternalRequest(filmId string) (externalservice.ExternalFilm, error)
}

type UserFilmRepository interface {
	GetAll(ctx context.Context, userId string) ([]userfilmservice.UserFilm, error)
	Add(ctx context.Context, userId string, filmId string) error
	Delete(ctx context.Context, userId string, filmId string) error
}

type UserFilmService interface {
	GetAll(userId string) ([]userfilmservice.UserFilm, error)
	Add(userId string, filmId string) error
	Delete(userId string, filmId string) error
}

type FilmService interface {
	GetOne(filmId string) (filmservice.Film, error)
	Search(match string, genres []string, page, pageSize int) ([]filmservice.FilmSearch, int64, error)
}

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (filmservice.Film, error)
	Save(film filmservice.Film) error
	Search(match string, genres []string, page, pageSize int) ([]filmservice.FilmSearch, int64, error)
}

type FilmSequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]filmsequelservice.FilmSequel, error)
	Save(filmId int, sequelId int) error
}

type FilmSequelService interface {
	GetAll(filmId string) ([]filmservice.Film, error)
}

type UserService interface {
	GetOne(userToken string) (userservice.User, error)
}

type UserRepository interface {
	GetOne(ctx context.Context, userToken string) (userservice.User, error)
}

type FilmSimilarRepository interface {
	GetAll(ctx context.Context, filmId string) ([]filmsimilarservice.FilmSimilar, error)
	Save(filmId int, similarId int) error
}

type FilmSimilarService interface {
	GetAll(filmId string) ([]filmservice.Film, error)
}
