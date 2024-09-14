package service

import "context"

type FilmRepository interface {
	GetFilms(ctx context.Context, filmId string)
}
