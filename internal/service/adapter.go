package service

import "context"

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (Film, error)
}
