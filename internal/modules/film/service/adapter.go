package service

import (
	"context"

	"kbox-api/internal/model"
)

type FilmRepository interface {
	GetOne(ctx context.Context, filmId string) (model.Film, error)
	Save(film model.Film) error
	Search(match string, genres []string, page, pageSize int) ([]model.Film, int64, error)
}
