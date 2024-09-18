package service

import "context"

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (Film, error)
	Save(film Film) error
}
