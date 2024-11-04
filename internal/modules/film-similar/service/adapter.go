package service

import (
	"context"

	"kbox-api/internal/model"
)

type FilmSimilarRepository interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSimilar, error)
	Save(filmId int, similarId int) error
}
