package service

import "context"

type FilmSimilarRepository interface {
	GetAll(ctx context.Context, filmId string) ([]FilmSimilar, error)
	Save(filmId int, similarId int) error
}
