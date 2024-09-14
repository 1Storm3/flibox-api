package repository

import "context"

type FilmRepository interface {
	GetFilm(ctx context.Context, filmId string)
}
