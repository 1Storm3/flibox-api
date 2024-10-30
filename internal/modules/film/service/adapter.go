package service

import "context"

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (Film, error)
	Save(film Film) error
	Search(match string, genres []string, page, pageSize int) ([]FilmSearch, int64, error)
}

type FilmServiceI interface {
	GetOne(filmId string) (Film, error)
}
