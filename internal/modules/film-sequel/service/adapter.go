package service

import (
	"context"

	"kbox-api/internal/model"
)

type FilmSequelRepository interface {
	GetAll(ctx context.Context, filmId string) ([]model.FilmSequel, error)
	Save(filmId int, sequelId int) error
}
